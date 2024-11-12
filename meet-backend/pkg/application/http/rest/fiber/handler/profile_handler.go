package handler

import (
	"github.com/erodriguezg/meet/pkg/application/http/rest"
	"github.com/erodriguezg/meet/pkg/core/service"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type profileFiberHandler struct {
	profileService service.ProfileService
	log            *zap.Logger
}

func NewProfileFiberHandler(
	profileService service.ProfileService,
	log *zap.Logger) FiberHandler {
	return &profileFiberHandler{profileService, log}
}

// RegisterRoutes implements FiberHandler
func (port *profileFiberHandler) RegisterRoutes(fiberRouter *fiber.Router) {
	router := *fiberRouter
	group := router.Group("/profile")
	group.Get("/all", port.findAllProfiles)
}

// privates

// ShowAccount godoc
// @Summary      Find All Profiles
// @Description  Get all profiles
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Success      200  {object}  rest.ApiResponse[[]domain.Profile]
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/profile/all [get]
func (port *profileFiberHandler) findAllProfiles(c *fiber.Ctx) error {
	port.log.Debug("-> findAllProfiles")
	profiles, err := port.profileService.FindAll()
	if err != nil {
		return err
	}
	return c.JSON(rest.ApiOk(&profiles))
}
