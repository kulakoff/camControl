package endpoint

import (
	"camControl/internal/models"
	"camControl/internal/service"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
)

type PTZHandler struct {
	service service.PTZService
}

func NewPTZHandler(s service.PTZService) *PTZHandler {
	return &PTZHandler{service: s}
}

func (h *PTZHandler) MoveCamera(c echo.Context) error {
	slog.Info("PTZHandler | MoveCamera")
	//TODO: add check request params
	var req models.PTZRequest
	err := c.Bind(&req)
	if err != nil {
		slog.Error("MoveCamera, bind request", "err", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	slog.Info("PTZHandler | MoveCamera", "req", req)

	if err := h.service.MoveCamera(c.Request().Context(), &req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "success"})
}

func (h *PTZHandler) MoveCameraSimple(c echo.Context) error {
	slog.Info("PTZHandler | MoveCameraSimple")
	//TODO: add check request params
	var req models.PTZRequest
	err := c.Bind(&req)
	if err != nil {
		slog.Error("MoveCamera, bind request", "err", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	slog.Info("MoveCameraSimple", "data", req)

	return c.NoContent(http.StatusNoContent)
}

func (h *PTZHandler) Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, "pong")
}

func (h *PTZHandler) RegisterRoutes(e *echo.Echo) {
	g := e.Group("/ptz")
	g.POST("/move", h.MoveCamera)
	g.POST("/move/simple", h.MoveCameraSimple)
	g.GET("/ping", h.Ping)
}
