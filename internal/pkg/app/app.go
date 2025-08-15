package app

import (
	"camControl/internal/app/camControl/config"
	camLog "camControl/internal/app/camControl/custom_logger"
	"camControl/internal/app/camControl/endpoint"
	"camControl/internal/app/camControl/repository"
	"camControl/internal/app/camControl/service"
	"camControl/internal/app/camControl/storage"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log/slog"
	"os"
)

type App struct {
	echo       *echo.Echo
	cfg        *config.Config
	logger     *slog.Logger
	storage    *storage.PSQLStorage
	camRepo    repository.CameraRepository
	ptzService service.PTZService
	ptzHandler *endpoint.PTZHandler
}

func New() (*App, error) {
	a := &App{}

	// load config
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Error loading config: ", err)
	}
	a.cfg = cfg

	// config logger
	a.logger = camLog.New(cfg.LogLevel)
	a.logger.Info("App started", "logLevel", cfg.LogLevel)

	// init postgres storage
	a.storage, err = storage.NewPSQLStorage(&a.cfg.Db, a.logger)
	if err != nil {
		a.logger.Error("Error creating storage: ", err)
		os.Exit(1)
	}
	defer a.storage.Close()

	// layer 01
	a.camRepo = repository.NewCameraRepository(a.storage.DB, a.logger)

	// layer 02
	a.ptzService = service.NewPTZService(a.camRepo, a.logger)

	// layer 03
	a.ptzHandler = endpoint.NewPTZHandler(a.ptzService, a.logger)

	a.echo = echo.New()
	a.echo.HideBanner = false
	a.echo.Use(middleware.Logger())

	a.echo.Use(middleware.Recover())

	a.ptzHandler.RegisterRoutes(a.echo)

	return a, nil
}

func (a *App) Start() error {
	err := a.echo.Start(a.cfg.Server.Port)
	if err != nil {
		return err
	}
	return nil
}
