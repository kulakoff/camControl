package main

import (
	"camControl/internal/config"
	"camControl/internal/endpoint"
	"camControl/internal/repository"
	"camControl/internal/service"
	"camControl/internal/storage"
	"github.com/labstack/echo/v4/middleware"
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
)

func main() {
	// TODO: config logger level from ENV
	// config logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	logger.Info("App started")

	// load config
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Error loading config: ", err)
	}
	slog.Info("Debug conf", "conf", cfg)

	// init postgres storage
	camStorage, err := storage.NewPSQLStorage(&cfg.Db, logger)
	if err != nil {
		slog.Error("Error creating storage: ", err)
		os.Exit(1)
	}
	defer camStorage.Close()

	camRepo := repository.New(camStorage.DB)
	camService := service.New(camRepo)

	e := echo.New()
	e.HideBanner = false

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	cameraEndpoint := endpoint.New(camService)
	apiV1 := e.Group("/api/v1")
	cameraEndpoint.RegisterRoutes(apiV1)

	go e.Logger.Fatal(e.Start(cfg.Server.Port))
}
