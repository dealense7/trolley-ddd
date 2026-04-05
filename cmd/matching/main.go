package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/dealense7/go-rates-ddd/internal/common/cfg"
	"github.com/dealense7/go-rates-ddd/internal/common/logger"
	"github.com/dealense7/go-rates-ddd/internal/domain/product"
	"github.com/dealense7/go-rates-ddd/internal/infrastructure/elastic"
	"github.com/dealense7/go-rates-ddd/internal/infrastructure/embedder"
	"github.com/dealense7/go-rates-ddd/internal/infrastructure/persistence/mysql"
	"go.uber.org/fx"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// ── Tuning knobs ──────────────────────────────────────────────────────────────

const (
	// alpha: weight of text score vs cosine in the final blend.
	// Cross-script pairs naturally score near 0 on text → cosine carries them.
	alpha = 0.65

	// exclusiveK: how aggressively high-IDF exclusive tokens penalize a match.
	// 0   = no penalty (pure coverage product).
	// 0.75 = default — tune up if you still see flavor variants matching,
	//        tune down if legitimate cross-language matches start dropping.
	exclusiveK = 0.75

	minFinalScore             = 0.65
	elasticCandidateThreshold = 0.75
	maxCandidates             = 10
	maxGoroutines             = 30
)

// ── Regexp ────────────────────────────────────────────────────────────────────

var (
	tokenRe         = regexp.MustCompile(`[\p{L}\p{N}]{2,}`)
	quantityTokenRe = regexp.MustCompile(`(?i)^\d+([.,]\d+)?\s*\p{L}{0,3}$`)

	// Multipack: "4 X 62 G", "3x39g", "2 × 50 ml" → extract the multiplier (4, 3, 2).
	// Must be checked BEFORE the generic qtyRe so it isn't partially consumed.
	multipackRe = regexp.MustCompile(`(?i)(\d+)\s*[xX×]\s*\d+(?:[.,]\d+)?\s*(?:ml|мл|მლ|cl|dl|л|l|ლ|g|г|kg|кг|oz|lb)\.?`)

	// Extended to cover count units (pz, uds, capsule, etc.).
	qtyRe = regexp.MustCompile(
		`(?i)(\d+(?:[.,]\d+)?)\s*` +
			`(ml|мл|მლ|cl|dl|л|l|ლ|fl\.?\s*oz|oz|kg|кг|g|г|lb|` +
			`pz|pcs|pc|ud|uds|uc|unid|unidades|` +
			`capsule|capsulas|capsules|caps|` +
			`unit|units|stk|st|szt|ks|ც|ცალი)\.?`,
	)
)

// ── Quantity ──────────────────────────────────────────────────────────────────

type qtyKind int8

const (
	qtyVolume qtyKind = iota // ml
	qtyWeight                // g
	qtyCount                 // pieces, uds, pz, capsules …
	qtyPack                  // the N in "N × size" — a count of units
)

type canonicalQty struct {
	value float64
	kind  qtyKind
}

var (
	toML = map[string]float64{
		"ml": 1, "мл": 1, "მლ": 1,
		"cl": 10,
		"dl": 100,
		"l":  1000, "л": 1000, "ლ": 1000,
		"fl oz": 29.5735, "fl.oz": 29.5735, "floz": 29.5735,
		"oz": 29.5735,
	}
	toG = map[string]float64{
		"g": 1, "г": 1,
		"kg": 1000, "кг": 1000,
		"lb": 453.592,
	}
	// Count units all normalise to 1 (the value IS the count).
	toCount = map[string]struct{}{
		"pz": {}, "pcs": {}, "pc": {},
		"ud": {}, "uds": {}, "uc": {}, "unid": {}, "unidades": {},
		"capsule": {}, "capsulas": {}, "capsules": {}, "caps": {},
		"unit": {}, "units": {},
		"stk": {}, "st": {}, "szt": {}, "ks": {},
		"ც": {}, "ცალი": {},
	}
)

func parseQty(name string) *canonicalQty {
	lower := strings.ToLower(name)

	// Multi-pack check first: "4 × 62 G" → kind=qtyPack, value=4.
	if m := multipackRe.FindStringSubmatch(lower); m != nil {
		if n, err := strconv.ParseFloat(m[1], 64); err == nil && n > 0 {
			return &canonicalQty{n, qtyPack}
		}
	}

	m := qtyRe.FindStringSubmatch(lower)
	if m == nil {
		return nil
	}
	val, err := strconv.ParseFloat(strings.ReplaceAll(m[1], ",", "."), 64)
	if err != nil || val <= 0 {
		return nil
	}
	unit := strings.TrimSpace(m[2])

	if mult, ok := toML[unit]; ok {
		return &canonicalQty{math.Round(val*mult*100) / 100, qtyVolume}
	}
	if mult, ok := toG[unit]; ok {
		return &canonicalQty{math.Round(val*mult*100) / 100, qtyWeight}
	}
	if _, ok := toCount[unit]; ok {
		return &canonicalQty{val, qtyCount}
	}
	return nil
}

// qtyMultiplier returns 0.25 when sizes explicitly differ, 1.0 otherwise.
// Pack counts and plain unit counts are treated as the same kind so
// "4 × 62 G" (pack=4) vs "2 Uds" (count=2) correctly fires the penalty.
func qtyMultiplier(nameA, nameB string) float64 {
	qa, qb := parseQty(nameA), parseQty(nameB)
	if qa == nil || qb == nil {
		return 1.0
	}

	// Treat qtyPack and qtyCount as interchangeable (both mean "N units").
	aIsCount := qa.kind == qtyCount || qa.kind == qtyPack
	bIsCount := qb.kind == qtyCount || qb.kind == qtyPack
	comparable := (qa.kind == qb.kind) || (aIsCount && bIsCount)
	if !comparable {
		return 1.0
	}

	ratio := qa.value / qb.value
	if ratio > 1 {
		ratio = 1 / ratio
	}
	if ratio >= 0.99 {
		return 1.0
	}
	return 0.25
}

// ── Diacritic stripping ───────────────────────────────────────────────────────

func stripDiacritics(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	out, _, err := transform.String(t, s)
	if err != nil {
		return s
	}
	return out
}

// ── Tokeniser ─────────────────────────────────────────────────────────────────

func tokenize(name string) []string {
	lower := stripDiacritics(strings.ToLower(name))
	raw := tokenRe.FindAllString(lower, -1)
	out := make([]string, 0, len(raw))
	seen := make(map[string]bool, len(raw))
	for _, tok := range raw {
		if !quantityTokenRe.MatchString(tok) && !seen[tok] {
			out = append(out, tok)
			seen[tok] = true
		}
	}
	return out
}

// ── IDF ───────────────────────────────────────────────────────────────────────

type idfIndex struct {
	scores map[string]float64
	n      int
}

func buildIDF(allNames []string) *idfIndex {
	df := make(map[string]int, 50_000)
	for _, name := range allNames {
		for _, tok := range tokenize(name) {
			df[tok]++
		}
	}
	n := len(allNames)
	scores := make(map[string]float64, len(df))
	for tok, count := range df {
		scores[tok] = math.Log(float64(n+1)/float64(count+1)) + 1.0
	}
	return &idfIndex{scores: scores, n: n}
}

func (idx *idfIndex) weight(tok string) float64 {
	if w, ok := idx.scores[tok]; ok {
		return w
	}
	return math.Log(float64(idx.n+1)) + 1.0
}

func (idx *idfIndex) tokenWeights(name string) map[string]float64 {
	toks := tokenize(name)
	w := make(map[string]float64, len(toks))
	for _, tok := range toks {
		w[tok] = idx.weight(tok)
	}
	return w
}

// ── Text score ────────────────────────────────────────────────────────────────
//
// Two components multiplied together:
//
//  1. Coverage product (recall_A × recall_B):
//     "what fraction of each name's IDF weight is explained by the other?"
//     Symmetric: both sides must agree. Better than Jaccard for detecting when
//     one name has extra tokens the other doesn't — e.g. "Ultra Paradise" vs
//     "Ultra White" each have 1 unexplained token → each recall ≈ 0.75 →
//     product ≈ 0.56 (vs Jaccard ≈ 0.60).
//
//  2. Exclusive IDF penalty:
//     The single highest-IDF token present in only one name signals a variant
//     (rosé, honey, white, pearl, integral, paradise …). Penalty scales with
//     how large that IDF is relative to the shared IDF weight.
//     At exclusiveK=0 this is disabled; tune up to tighten, down to loosen.
//
// Cross-script behaviour is emergent: Georgian vs Latin tokens never overlap
// → shared=0 → textScore=0 → final ≈ (1-alpha)×cosine < threshold → rejected.

func (idx *idfIndex) textScore(a, b map[string]float64) float64 {
	var totalA, totalB, shared float64
	var maxExclA, maxExclB float64

	for tok, wa := range a {
		totalA += wa
		if _, ok := b[tok]; ok {
			shared += wa
		} else if wa > maxExclA {
			maxExclA = wa
		}
	}
	for tok, wb := range b {
		totalB += wb
		if _, ok := a[tok]; !ok && wb > maxExclB {
			maxExclB = wb
		}
	}

	if totalA == 0 || totalB == 0 || shared == 0 {
		return 0
	}

	// Component 1: coverage product.
	score := (shared / totalA) * (shared / totalB)

	// Component 2: exclusive IDF penalty.
	if exclusiveK > 0 {
		maxExcl := math.Max(maxExclA, maxExclB)
		if maxExcl > 0 {
			score *= 1.0 / (1.0 + exclusiveK*(maxExcl/shared))
		}
	}

	return score
}

// ── Combined score ────────────────────────────────────────────────────────────

func (idx *idfIndex) score(cosine float64, nameA, nameB string) float64 {
	tokA := idx.tokenWeights(nameA)
	tokB := idx.tokenWeights(nameB)
	text := idx.textScore(tokA, tokB)
	qty := qtyMultiplier(nameA, nameB)
	return (alpha*text + (1-alpha)*cosine) * qty
}

// ── Pair dedup key ────────────────────────────────────────────────────────────

func pairKey(a, b int64) string {
	if a > b {
		a, b = b, a
	}
	return strconv.FormatInt(a, 10) + "_" + strconv.FormatInt(b, 10)
}

// ── Main ──────────────────────────────────────────────────────────────────────

func main() {
	app := fx.New(
		fx.Provide(
			cfg.NewConfig,
			logger.NewFactory,
			mysql.NewDB,
			elastic.New,
			embedder.New,
		),
		fx.Invoke(runEmbeddings),
		fx.Invoke(func(s fx.Shutdowner) { _ = s.Shutdown() }),
	)
	app.Run()
}

func runEmbeddings(db *mysql.DB, e *elastic.Client, l logger.Factory, em *embedder.Client) {
	log := l.For(logger.General)
	ctx := context.Background()
	start := time.Now()

	e.CreateIndex(ctx)

	// ── Single query for both IDF and name lookup ─────────────────────────────
	var catalog []struct {
		ID      int64  `db:"id"`
		RawName string `db:"raw_name"`
	}
	if err := db.SelectContext(ctx, &catalog, `
		SELECT id, raw_name
		FROM scraped_products
		WHERE raw_name IS NOT NULL AND raw_name != ''
	`); err != nil {
		log.Error("load catalog: " + err.Error())
		return
	}

	nameByID := make(map[int64]string, len(catalog))
	allNames := make([]string, 0, len(catalog))
	for _, row := range catalog {
		nameByID[row.ID] = row.RawName
		allNames = append(allNames, row.RawName)
	}

	idf := buildIDF(allNames)
	log.Info(fmt.Sprintf("IDF built from %d products", len(allNames)))

	// ── Load items with images ────────────────────────────────────────────────
	type itemStruct struct {
		product.Image
		RawName  string `db:"raw_name"`
		BranchId int64  `db:"branch_id"`
	}
	var items []itemStruct
	if err := db.SelectContext(ctx, &items, `
		SELECT sp.raw_name, sp.branch_id, pi.*
		FROM product_images pi
		JOIN scraped_products sp ON pi.product_id = sp.id
	`); err != nil {
		log.Error("load items: " + err.Error())
		return
	}
	log.Info(fmt.Sprintf("processing %d items", len(items)))

	var (
		wg        sync.WaitGroup
		guard     = make(chan struct{}, maxGoroutines)
		seenPairs sync.Map
	)

	for _, item := range items {
		item := item
		wg.Add(1)
		go func() {
			guard <- struct{}{}
			defer wg.Done()
			defer func() { <-guard }()

			single, err := e.GetSingle(ctx, item.ProductId)
			if err != nil || single == nil {
				return
			}

			candidates, err := e.FindSimilar(
				ctx,
				single.Embeddings,
				item.BranchId,
				maxCandidates,
				elasticCandidateThreshold,
			)
			if err != nil {
				return
			}

			selfID := fmt.Sprintf("%d", item.ProductId)

			for _, candidate := range candidates {
				if candidate.ProductID == selfID {
					continue
				}

				similarId, err := strconv.ParseInt(candidate.ProductID, 10, 64)
				if err != nil {
					continue
				}

				candidateName, ok := nameByID[similarId]
				if !ok {
					continue
				}

				finalScore := idf.score(candidate.Score, item.RawName, candidateName)
				if finalScore < minFinalScore {
					continue
				}

				// In-memory dedup: blocks the mirror goroutine (B→A) before DB round-trip.
				key := pairKey(item.ProductId, similarId)
				if _, loaded := seenPairs.LoadOrStore(key, struct{}{}); loaded {
					continue
				}

				// DB dedup: handles pairs persisted in previous runs.
				var existingID int64
				err = db.GetContext(ctx, &existingID, `
					SELECT id FROM product_matches
					WHERE (scraped_product_id = ? AND similar_scraped_product_id = ?)
					   OR (scraped_product_id = ? AND similar_scraped_product_id = ?)
					LIMIT 1
				`, item.ProductId, similarId, similarId, item.ProductId)
				if !errors.Is(err, sql.ErrNoRows) {
					continue
				}

				match := product.NewMatch(
					item.ProductId,
					similarId,
					fmt.Sprintf("%.4f", finalScore),
					product.MatchTypeML,
				)
				_, err = db.NamedExecContext(ctx, `
					INSERT INTO product_matches
						(scraped_product_id, similar_scraped_product_id, match_type, confidence_score)
					VALUES
						(:scraped_product_id, :similar_scraped_product_id, :match_type, :confidence_score)
				`, match)
				if err != nil {
					log.Error("insert match: " + err.Error())
				}
			}
		}()
	}

	wg.Wait()
	log.Info(fmt.Sprintf("done in %s", time.Since(start).Round(time.Millisecond)))
}
