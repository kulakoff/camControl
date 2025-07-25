package endpoint

import (
	"camControl/internal/models"
	"camControl/internal/service"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"strconv"
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

func (h *PTZHandler) GetPresets(c echo.Context) error {
	slog.Info("PTZHandler | GetPresets")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid ID")
	}

	presets, _ := h.service.GetPresets(c.Request().Context(), uint(id))

	return c.JSON(http.StatusOK, presets)
}

func (h *PTZHandler) Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, "pong")
}

func (h *PTZHandler) RegisterRoutes(e *echo.Echo) {
	g := e.Group("/ptz")
	g.POST("/move", h.MoveCamera)
	g.GET("/presets/:id", h.GetPresets)
	g.GET("/ping", h.Ping)
}
