package main

import (
	"camControl/internal/config"
	"camControl/internal/storage"
	"context"
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
)

func main() {
	slog.Info("Start")
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Error loading config: ", err)
	}
	slog.Info("Debug conf", "conf", cfg)

	storage, err := storage.NewStorage(&cfg.Db)
	if err != nil {
		slog.Error("Error creating storage: ", err)
		os.Exit(1)
	}
	defer storage.Close()

	testCam, _ := storage.GetCameraByID(context.Background(), 1)
	slog.Info("Test | Camera", "result cam ip", testCam.IP)

	e := echo.New()
	e.HideBanner = true

}
