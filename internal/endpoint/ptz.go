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
	cameraId, err := strconv.Atoi(c.Param("cameraId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid ID")
	}
	slog.Info("PTZHandler | GetPresets from camera ID: " + strconv.Itoa(cameraId))

	presets, _ := h.service.GetPresets(c.Request().Context(), uint(cameraId))

	return c.JSON(http.StatusOK, presets)
}

func (h *PTZHandler) GoToPreset(c echo.Context) error {
	slog.Info("GoToPreset")
	var req models.PTZRequestPreset
	err := c.Bind(&req)
	if err != nil {
		slog.Error("GoToPreset, bind request err", "err", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	slog.Info("PTZHandler | GoToPreset", "req", req)
	h.service.GoToPreset(c.Request().Context(), uint(req.CameraID), req.PresetToken)
	return c.NoContent(http.StatusNoContent)
}

func (h *PTZHandler) Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, "pong")
}

func (h *PTZHandler) RegisterRoutes(e *echo.Echo) {
	g := e.Group("/ptz")
	g.POST("/move", h.MoveCamera)
	g.GET("/presets/:cameraId", h.GetPresets) // get presets by cameraID
	g.POST("/preset", h.GoToPreset)           // go to PTZ preset
	// TODO: create New PTZ preset
	// TODO: delete PTZ preset
	// TODO: update PTZ preset
	g.GET("/ping", h.Ping)
}
