package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gesangwidigdo/url-shortener/internal/platform/config"
	"github.com/gesangwidigdo/url-shortener/internal/platform/db"
	"github.com/gesangwidigdo/url-shortener/internal/platform/httpapi"
	"github.com/gesangwidigdo/url-shortener/internal/platform/logger"
	"github.com/gesangwidigdo/url-shortener/internal/platform/middleware"
	"github.com/gesangwidigdo/url-shortener/internal/platform/router"
	"github.com/gesangwidigdo/url-shortener/internal/url"
	"github.com/joho/godotenv"
)

func main() {
	if err := run(); err != nil {
		slog.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run() error {
	_ = godotenv.Load()
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	log := logger.New(cfg.Env)
	slog.SetDefault(log)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := db.NewPostgres(ctx, db.Config{
		URL:      cfg.DB.URL,
		MaxConns: cfg.DB.MaxConns,
	})
	if err != nil {
		return err
	}
	defer pool.Close()

	urlRepo := url.NewRepository(pool)
	urlSvc := url.NewService(urlRepo, cfg.BaseURL)
	urlHandler := url.NewHandler(urlSvc)

	handler := router.New(
		[]router.Middleware{
			middleware.RequestID(),
			middleware.Logging(log),
			middleware.Recovery(log),
		},
		urlHandler,
	)

	srv := httpapi.New(httpapi.Config{
		Addr:            cfg.HTTP.Addr,
		ReadTimeout:     cfg.HTTP.ReadTimeout,
		WriteTimeout:    cfg.HTTP.WriteTimeout,
		ShutdownTimeout: cfg.HTTP.ShutdownTimeout,
	}, handler, log)
	return srv.Run(ctx)
}
