package endpoint

import (
	"camControl/internal/service"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"strconv"
)

type Endpoint interface {
	GetCameraInfo(c echo.Context) error
	Status(c echo.Context) error
	RegisterRoutes(group *echo.Group)
}
type endpoint struct {
	service service.CameraService
}

func New(service service.CameraService) Endpoint {
	return &endpoint{service: service}
}

func (e *endpoint) Status(с echo.Context) error {
	return с.JSON(http.StatusOK, "Pong")
}

func (e *endpoint) GetCameraInfo(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	slog.Info("GetCameraInfo", "id", id)
	camera, err := e.service.GetCamera(uint(id))
	if err != nil {
		return err
	}
	slog.Info("GetCameraInfo", "cam", camera)
	return c.JSON(http.StatusOK, camera)
}

func (e *endpoint) RegisterRoutes(g *echo.Group) {
	g.GET("/ping", e.Status)
	g.GET("/camera/:id", e.GetCameraInfo)
}
