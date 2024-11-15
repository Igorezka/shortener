package main

import (
	"github.com/Igorezka/shortener/internal/app/config"
	"github.com/Igorezka/shortener/internal/app/http-server/router"
	"github.com/Igorezka/shortener/internal/app/logger"
	"github.com/Igorezka/shortener/internal/app/storage"
	"github.com/Igorezka/shortener/internal/app/storage/memory"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	cfg := config.New()
	log, err := logger.New(cfg.LogLevel)
	if err != nil {
		panic(err)
	}
	defer func(log *zap.Logger) {
		_ = log.Sync()
	}(log)
	store := storage.New(memory.New())

	log.Info(
		"starting server",
		zap.String("Address", cfg.RunAddr),
		zap.String("Base URL", cfg.BaseURL),
	)
	err = http.ListenAndServe(cfg.RunAddr, router.New(cfg, store, log))
	if err != nil {
		panic(err)
	}
}
