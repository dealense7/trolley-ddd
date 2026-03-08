package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/dealense7/go-rates-ddd/internal/common/cfg"
	"github.com/dealense7/go-rates-ddd/internal/common/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	*sqlx.DB
}

func NewDB(cfg *cfg.Config, loggerFactory logger.Factory) (*DB, error) {

	log := loggerFactory.For(logger.General)

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
	)

	db, err := sqlx.Connect(cfg.DB.Driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection
	if err := db.PingContext(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info("database connection established")

	return &DB{DB: db}, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}
