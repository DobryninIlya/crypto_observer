package server

import (
	"context"
	coingecko "cryptoObserver/internal/app/coingeko"
	"cryptoObserver/internal/app/migrations"
	"cryptoObserver/internal/app/store/sqlstore"
	worker "cryptoObserver/internal/app/workers"
	"database/sql"
	"github.com/sirupsen/logrus"
	"time"
)

func Start(ctx context.Context, config *Config) (*App, error) {
	db, err := newDB(config.GetDBConnectionString())
	if err != nil {
		return nil, err
	}
	store := sqlstore.New(db)
	logger := logrus.New()
	migrations.MakeMigrations(db, logger)
	cryptoAPI := coingecko.NewCoinGeckoClient(config.CryptoAPI.Token)
	pool := worker.NewWorkerPool(ctx, cryptoAPI, store, config.WorkerPool.Size, time.Duration(config.WorkerPool.UpdateTime)*time.Second, logger)
	defer pool.Start()
	srv := newApp(ctx, store, *config, logger, pool)
	return srv, nil
}

func newDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
