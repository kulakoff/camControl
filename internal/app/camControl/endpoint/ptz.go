package endpoint

import (
	"camControl/internal/app/camControl/models"
	"camControl/internal/app/camControl/service"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"strconv"
)

type PTZHandler struct {
	service service.PTZService
	logger  slog.Logger
}

func New(s service.PTZService, logger *slog.Logger) *PTZHandler {
	return &PTZHandler{service: s, logger: *logger}
}

func (h *PTZHandler) MoveCamera(c echo.Context) error {
	h.logger.Debug("PTZService | MoveCamera")
	//TODO: add check request params
	var req models.PTZRequest
	err := c.Bind(&req)
	if err != nil {
		h.logger.Error("MoveCamera, bind request", "err", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	h.logger.Debug("PTZHandler | MoveCamera", "req", req)

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
	h.logger.Debug("PTZHandler | GetPresets from camera ID: " + strconv.Itoa(cameraId))

	presets, _ := h.service.GetPresets(c.Request().Context(), uint(cameraId))

	return c.JSON(http.StatusOK, presets)
}

func (h *PTZHandler) GoToPreset(c echo.Context) error {
	h.logger.Debug("GoToPreset")
	var req models.PTZRequestPreset
	err := c.Bind(&req)
	if err != nil {
		slog.Error("GoToPreset, bind request err", "err", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	h.logger.Debug("PTZHandler | GoToPreset", "req", req)
	h.service.GoToPreset(c.Request().Context(), uint(req.CameraID), req.PresetToken)
	return c.NoContent(http.StatusNoContent)
}

func (h *PTZHandler) SetPreset(c echo.Context) error {
	h.logger.Debug("SetPreset")
	var req models.PTZRequestPreset
	err := c.Bind(&req)
	if err != nil {
		slog.Error("SetPreset, bind request err", "err", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	h.logger.Debug("PTZHandler | SetPreset", "req", req)
	err = h.service.SetPreset(c.Request().Context(), uint(req.CameraID), req.PresetToken)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *PTZHandler) RemovePTZPreset(c echo.Context) error {
	h.logger.Debug("RemovePTZPreset")
	cameraId, err := strconv.Atoi(c.Param("cameraId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid cameraId")
	}
	presetToken, err := strconv.Atoi(c.Param("presetToken"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "presetToken")
	}

	h.logger.Debug("RemovePTZPreset data", "cameraId", cameraId, "presetToken", presetToken)
	h.service.RemovePreset(c.Request().Context(), uint(cameraId), strconv.Itoa(presetToken))
	return c.NoContent(http.StatusNoContent)
}

func (h *PTZHandler) Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, "pong")
}

// RegisterRoutes implement register API routes
func (h *PTZHandler) RegisterRoutes(e *echo.Echo) {
	g := e.Group("/ptz")
	g.GET("/ping", h.Ping)
	g.POST("/move", h.MoveCamera)
	g.GET("/preset/:cameraId", h.GetPresets) // get presets by cameraID
	g.POST("/preset", h.GoToPreset)          // go to PTZ preset
	g.POST("/preset/set", h.SetPreset)
	g.DELETE("/preset/:cameraId/:presetToken", h.RemovePTZPreset)
}
