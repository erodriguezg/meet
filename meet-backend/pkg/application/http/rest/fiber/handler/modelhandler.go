package handler

import (
	"fmt"

	"github.com/erodriguezg/meet/pkg/application/http/rest"
	"github.com/erodriguezg/meet/pkg/application/http/rest/fiber/fiberidentity"
	"github.com/erodriguezg/meet/pkg/application/http/security"
	"github.com/erodriguezg/meet/pkg/core/domain"
	"github.com/erodriguezg/meet/pkg/core/dto"
	"github.com/erodriguezg/meet/pkg/core/service"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type modelFiberHandler struct {
	modelService    service.ModelService
	securityService security.HttpSecurityService
	log             *zap.Logger
}

func NewModelFiberHandler(
	modelService service.ModelService,
	securityService security.HttpSecurityService,
	log *zap.Logger,
) FiberHandler {
	return &modelFiberHandler{modelService, securityService, log}
}

func (port *modelFiberHandler) RegisterRoutes(fiberRouter *fiber.Router) {
	router := *fiberRouter
	group := router.Group("/model")
	group.Post("/register", port.registerModel)
	group.Get("/:nickname", port.findModelByNickName)
	group.Post("/:nickname/prepare-profile-img", port.prepareUploadProfileImg)
	group.Post("/search", port.searchModels)
}

// ShowAccount godoc
// @Summary      Register Model
// @Description  Register a Model
// @Tags         Model
// @Accept       json
// @Produce      json
// @Param        data body dto.ModelRegisterDto true "The registration data"
// @Success      200  {object}  rest.ApiResponse[string]
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/model/register [post]
func (port *modelFiberHandler) registerModel(c *fiber.Ctx) error {

	identity, err := security.MustHavePermission(port.securityService, c, domain.PermissionCodeRegisterModel)
	if err != nil {
		return err
	}

	var payload dto.ModelRegisterDto
	err = c.BodyParser(&payload)
	if err != nil {
		return err
	}

	if payload.PersonId != identity.PersonId {
		return fiberidentity.NewAccessDeniedError(fmt.Errorf("session person id not match"))
	}

	port.log.Debug("-> registerModel", zap.Any("payload", &payload))
	err = port.modelService.RegisterModel(payload)
	if err != nil {
		return err
	}
	return c.JSON(rest.ApiOkEmpty())
}

// ShowAccount godoc
// @Summary      Find Model By NickName
// @Description  Find a model by nickname and get public data
// @Tags         Model
// @Accept       json
// @Produce      json
// @Param        nickname   path     string  false  "nickname of the model"
// @Success      200  {object}  rest.ApiResponse[domain.Model]
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/model/{nickname} [get]
func (port *modelFiberHandler) findModelByNickName(c *fiber.Ctx) error {
	nickNameParam := c.Params("nickname")
	port.log.Debug("-> findModelByNickName", zap.String("nickname", nickNameParam))
	model, err := port.modelService.FindModelByNickName(nickNameParam)
	if err != nil {
		return err
	}
	return c.JSON(rest.ApiOk(model))
}

// ShowAccount godoc
// @Summary      Prepare Upload Profile Image
// @Description  Prepare for upload profile image
// @Tags         Model
// @Accept       json
// @Produce      json
// @Param        nickname   path     string  false  "nickname of the model"
// @Success      200  {object}  rest.ApiResponse[[]dto.ResourceUploadUrlDto]
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/model/{nickname}/prepare-profile-img [post]
func (port *modelFiberHandler) prepareUploadProfileImg(c *fiber.Ctx) error {
	nickNameParam := c.Params("nickname")

	_, err := security.MustHavePermissionToEditModel(port.securityService, c, nickNameParam)
	if err != nil {
		return err
	}

	port.log.Debug("-> prepareUploadProfileImg", zap.String("nickname", nickNameParam))
	uploadsDto, err := port.modelService.PrepareUploadUrlForProfileImage(nickNameParam)
	if err != nil {
		return err
	}
	return c.JSON(rest.ApiOkArray(uploadsDto))
}

// ShowAccount godoc
// @Summary      Search Models
// @Description  Search Models by filters
// @Tags         Model
// @Accept       json
// @Produce      json
// @Param        first  query     string  true  "first position page"
// @Param        last   query     string  true  "last position page"
// @Param        data   body     domain.FilterSearchModel  true  "filters for search"
// @Success      200  {object}  rest.ApiResponse[domain.SearchModelResponse]
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/model/search [post]
func (port *modelFiberHandler) searchModels(c *fiber.Ctx) error {
	first := c.QueryInt("first")
	last := c.QueryInt("last")
	var filters domain.FilterSearchModel
	err := c.BodyParser(&filters)
	if err != nil {
		return err
	}
	searchResponse, err := port.modelService.SearchModels(filters, first, last)
	if err != nil {
		return err
	}
	return c.JSON(rest.ApiOk(searchResponse))
}
