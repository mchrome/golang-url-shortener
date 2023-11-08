package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/mchrome/url-compression-api/internal/app/apiserver"
	save "github.com/mchrome/url-compression-api/internal/app/apiserver/handlers/url"
	"github.com/mchrome/url-compression-api/internal/app/lib/logger/sl"
	storage "github.com/mchrome/url-compression-api/internal/app/store"
	"github.com/mchrome/url-compression-api/internal/config"
)

func main() {
	// load config
	cfg := config.MustLoad()

	// logger init
	log, err := setupLogger(cfg.Env)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	log.Info("starting api server")
	log.Debug("debug messages are enabled")

	// storage init
	storage, err := storage.NewStorage(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	fmt.Println(storage)

	// router init
	router := chi.NewRouter()

	// middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)

	// define api endpoints
	router.Post("/url", save.New(log, storage))

	log.Info("starting server", slog.String("address", cfg.Address))

	server := apiserver.New()
	if err := server.Start(); err != nil {
		log.Error("can't start api server")
		return
	}
}

func setupLogger(env string) (*slog.Logger, error) {

	var logger *slog.Logger

	switch env {
	case "local":
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "dev":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		return nil, errors.New(fmt.Sprintf("unknown env value: %s", env))
	}

	return logger, nil

}
