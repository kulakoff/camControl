package endpoint

import (
	"camControl/internal/models"
	"camControl/internal/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

type PTZHandler struct {
	service service.PTZService
}

func NewPTZHandler(s service.PTZService) *PTZHandler {
	return &PTZHandler{service: s}
}

func (h *PTZHandler) MoveCamera(c echo.Context) error {
	//TODO: add check request params
	var req models.PTZRequest
	err := c.Bind(&req)
	if err != nil {
		return err
	}

	if err := h.service.MoveCamera(c.Request().Context(), &req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "success"})
}

func (h *PTZHandler) RegisterRoutes(e *echo.Echo) {
	g := e.Group("/ptz")
	g.POST("/move", h.MoveCamera)
}
