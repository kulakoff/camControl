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

/**
TODO:
	1: monitoring
		- add monitoring in config (prometheus)
		- if monitoring enabled check available IP camera from monitoring before request,
	2: add Prometheus metrics per camera all PTZ requests
*/

func main() {
	// load config
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Error loading config: ", err)
	}

	// config logger
	// TODO: config logger level from ENV
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	logger.Info("App started")

	// init postgres storage
	camStorage, err := storage.NewPSQLStorage(&cfg.Db, logger)
	if err != nil {
		logger.Error("Error creating storage: ", err)
		os.Exit(1)
	}
	defer camStorage.Close()

	// layer 01
	camRepo := repository.NewCameraRepository(camStorage.DB, logger)
	// layer 02
	ptzService := service.NewPTZService(camRepo, logger)
	// layer 03
	ptzHandler := endpoint.NewPTZHandler(ptzService, logger)

	e := echo.New()
	e.HideBanner = false
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	ptzHandler.RegisterRoutes(e)

	go e.Logger.Fatal(e.Start(cfg.Server.Port))
}
