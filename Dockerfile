# ── Stage 1: build ────────────────────────────────────────────
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /bin/scraper ./cmd/scraper
RUN go build -o /bin/seeder  ./cmd/seeder

# Scraper
FROM alpine:3.19 AS scraper

# ca-certificates: needed for HTTPS requests in the scraper
# tzdata: time zones in logs
RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /bin/scraper /scraper

ENTRYPOINT ["/scraper"]

# Seeder
FROM alpine:3.19 AS seeder

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /bin/seeder /seeder

ENTRYPOINT ["/seeder"]