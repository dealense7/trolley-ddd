package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/dealense7/go-rates-ddd/internal/common/cfg"
	"github.com/dealense7/go-rates-ddd/internal/common/logger"
	"github.com/dealense7/go-rates-ddd/internal/domain/product"
	"github.com/dealense7/go-rates-ddd/internal/infrastructure/elastic"
	"github.com/dealense7/go-rates-ddd/internal/infrastructure/embedder"
	"github.com/dealense7/go-rates-ddd/internal/infrastructure/persistence/mysql"
	"github.com/dealense7/go-rates-ddd/scripts/seed"
	"go.uber.org/fx"
)

type SeedersParams struct {
	fx.In

	LoggerFactory logger.Factory
	Seeders       []seed.Seeder `group:"seeders"`
}

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
		fx.Invoke(func(shutdown fx.Shutdowner) { _ = shutdown.Shutdown() }),
	)

	app.Run()
}

func runEmbeddings(db *mysql.DB, e *elastic.Client, em *embedder.Client) {
	ctx := context.Background()
	e.CreateIndex(ctx)
	type itemStruct struct {
		product.Image
		RawName  string `db:"raw_name"`
		BranchId int64  `db:"branch_id"`
	}
	var items []itemStruct

	query := "SELECT sp.raw_name, sp.branch_id, pi.* FROM product_images as pi join products.scraped_products sp on pi.product_id = sp.id where pi.has_embeddings = false"

	err := db.SelectContext(ctx, &items, query)
	if err != nil {
		fmt.Println(err.Error())
	}

	maxAsyncJobs := 30
	var wg sync.WaitGroup

	guard := make(chan struct{}, maxAsyncJobs)

	for _, item := range items {
		go func() {
			wg.Add(1)
			guard <- struct{}{}
			defer wg.Done()
			defer func() { <-guard }()

			embeddings, _ := em.EmbedFused(item.ImageURL(), item.RawName)
			ed := elastic.Product{
				ProductID: strconv.FormatInt(item.ProductId, 10),
				BranchId:  strconv.FormatInt(item.BranchId, 10),
				Embedding: embeddings,
			}
			err = e.IndexProduct(ctx, ed)
			if err != nil {
				log.Println(err.Error())
			}
			query := `UPDATE product_images SET has_embeddings = true where product_id = :id`

			db.NamedExecContext(ctx, query, map[string]interface{}{"id": item.ProductId})
		}()
	}
	wg.Wait()

}
