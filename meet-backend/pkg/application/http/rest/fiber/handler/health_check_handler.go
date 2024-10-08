package handler

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type healthcheckHandler struct {
	log *zap.Logger
}

func NewHealthCheckHandler(log *zap.Logger) FiberHandler {
	return &healthcheckHandler{log}
}

// RegisterRoutes implements FiberHandler
func (port *healthcheckHandler) RegisterRoutes(fiberRouter *fiber.Router) {
	router := *fiberRouter
	router.Get("/health-check", port.healthCheck)
}

// ShowAccount godoc
// @Summary      Health Check
// @Description  Health Check Api Method
// @Tags         Health Check
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]bool
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/health-check [get]
func (port *healthcheckHandler) healthCheck(c *fiber.Ctx) error {
	return c.JSON(map[string]bool{"ok": true})
}
