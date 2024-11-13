package handler

import (
	"strconv"

	"github.com/erodriguezg/meet/pkg/application/http/rest"
	"github.com/erodriguezg/meet/pkg/application/http/security"
	"github.com/erodriguezg/meet/pkg/core/service"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/go-playground/validator/v10"
)

type packFiberHandler struct {
	packService     service.PackService
	securityService security.HttpSecurityService
	validate        *validator.Validate
	log             *zap.Logger
}

type PrepareUploadPackItemDto struct {
	ModelNickName string `json:"modelNickName"`
	PackNumber    int    `json:"packNumber"`
	TypeCode      string `json:"typeCode"`
	IsPublic      bool   `json:"isPublic"`
}

type PackDto struct {
	ModelNickName string `json:"modelNickName"`
	PackNumber    int    `json:"packNumber"`
}

type EditPackTitleDto struct {
	Title string `json:"title" validate:"required,max=30"`
}

type EditPackDescriptionDto struct {
	Description string `json:"description" validate:"max=280"`
}

func NewPackFiberHandler(
	packService service.PackService,
	securityService security.HttpSecurityService,
	validate *validator.Validate,
	log *zap.Logger,
) FiberHandler {
	return &packFiberHandler{packService, securityService, validate, log}
}

func (port *packFiberHandler) RegisterRoutes(fiberRouter *fiber.Router) {
	router := *fiberRouter
	group := router.Group("/pack")
	group.Put("/:modelNickName/new", port.createNewPack)
	group.Delete("/:modelNickName/:packNumber/:packItem", port.deletePackItem)
	group.Delete("/:modelNickName/:packNumber", port.deletePack)
	group.Post("/prepare-upload-item", port.prepareUploadForPackItem)
	group.Post("/ready-to-publish", port.readyToPublishPack)
	group.Post("/publish", port.publishPack)
	group.Post("/:modelNickName/:packNumber/title", port.editPackTitle)
	group.Post("/:modelNickName/:packNumber/description", port.editPackDescription)
	group.Get("/:modelNickName/:packNumber/info", port.getPackInfo)
	group.Get("/:modelNickName/:packNumber/items", port.getItemsFromPack)
	group.Get("/:modelNickName", port.getPacksFromModel)
}

// ShowAccount godoc
// @Summary      Create New Pack
// @Description  Create a new pack for model
// @Tags         Pack
// @Accept       json
// @Produce      json
// @Param        modelNickName   path     string  true  "model nickname"
// @Success      200  {object}  rest.ApiResponse[dto.PackDto]
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/pack/{modelNickName}/new [put]
func (port *packFiberHandler) createNewPack(c *fiber.Ctx) error {
	modelNickNameParam := c.Params("modelNickName")
	port.log.Debug("-> createNewPack", zap.String("modelNickName", modelNickNameParam))
	newPack, err := port.packService.CreateNewPack(modelNickNameParam)
	if err != nil {
		return err
	}
	return c.JSON(rest.ApiOk(newPack))
}

// ShowAccount godoc
// @Summary      Delete Pack Item
// @Description  Delete a item of one Pack
// @Tags         Pack
// @Accept       json
// @Produce      json
// @Param        modelNickName   path     string  true  "model nickname"
// @Param        packNumber   path     int  true  "pack number"
// @Param        packItem   path     int  true  "item number on the pack"
// @Success      200  {object}  rest.ApiResponse[string]
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/pack/{modelNickName}/{packNumber}/{packItem} [delete]
func (port *packFiberHandler) deletePackItem(c *fiber.Ctx) error {
	modelNickNameParam := c.Params("modelNickName")
	packNumber, err := strconv.Atoi(c.Params("packNumber"))
	if err != nil {
		return err
	}
	packItem, err := strconv.Atoi(c.Params("packItem"))
	if err != nil {
		return err
	}
	port.log.Debug("-> deletePackItem.",
		zap.String("modelNickName", modelNickNameParam),
		zap.Int("packNumber", packNumber),
		zap.Int("packItem", packItem))

	err = port.packService.DeletePackItem(modelNickNameParam, packNumber, packItem)
	if err != nil {
		return err
	}
	return c.JSON(rest.ApiOkEmpty())
}

// ShowAccount godoc
// @Summary      Delete Pack
// @Description  Delete one pack
// @Tags         Pack
// @Accept       json
// @Produce      json
// @Param        modelNickName   path     string  true  "model nickname"
// @Param        packNumber   path     int  true  "pack number"
// @Success      200  {object}  rest.ApiResponse[string]
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/pack/{modelNickName}/{packNumber} [delete]
func (port *packFiberHandler) deletePack(c *fiber.Ctx) error {
	modelNickNameParam := c.Params("modelNickName")
	packNumber, err := strconv.Atoi(c.Params("packNumber"))
	if err != nil {
		return err
	}
	port.log.Debug("-> deletePack",
		zap.String("modelNickName", modelNickNameParam),
		zap.Int("packNumber", packNumber))
	err = port.packService.DeletePack(modelNickNameParam, packNumber)
	if err != nil {
		return err
	}
	return c.JSON(rest.ApiOkEmpty())
}

// ShowAccount godoc
// @Summary      Prepare Upload Pack Item
// @Description  Prepare the upload for an item of one pack
// @Tags         Pack
// @Accept       json
// @Produce      json
// @Param        data body PrepareUploadPackItemDto true "Payload Data"
// @Success      200  {object}  rest.ApiResponse[[]dto.ResourceUploadUrlDto]
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/pack/prepare-upload-item [post]
func (port *packFiberHandler) prepareUploadForPackItem(c *fiber.Ctx) error {
	var payload PrepareUploadPackItemDto
	err := c.BodyParser(&payload)
	if err != nil {
		return err
	}
	port.log.Debug("-> prepareUploadForPackItem", zap.Any("payload", payload))
	uploadResources, err := port.packService.PrepareUploadForPackItem(payload.ModelNickName, payload.PackNumber, payload.TypeCode, payload.IsPublic)
	if err != nil {
		return err
	}
	return c.JSON(rest.ApiOkArray(uploadResources))
}

// ShowAccount godoc
// @Summary      Ready To Publish Pack
// @Description  Ready to publish the pack for moderators / admin revision
// @Tags         Pack
// @Accept       json
// @Produce      json
// @Param        data body PackDto true "Payload Data"
// @Success      200  {object}  rest.ApiResponse[string]
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/pack/ready-to-publish [post]
func (port *packFiberHandler) readyToPublishPack(c *fiber.Ctx) error {
	var payload PackDto
	err := c.BodyParser(&payload)
	if err != nil {
		return err
	}

	port.log.Debug("-> readyToPublishPack", zap.Any("payload", payload))

	err = port.packService.ReadyToPublishPack(payload.ModelNickName, payload.PackNumber)
	if err != nil {
		return err
	}

	return c.JSON(rest.ApiOkEmpty())
}

// ShowAccount godoc
// @Summary      Publish Pack
// @Description  Publish the pack ready for consumers
// @Tags         Pack
// @Accept       json
// @Produce      json
// @Param        data body PackDto true "Payload Data"
// @Success      200  {object}  rest.ApiResponse[string]
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/pack/publish [post]
func (port *packFiberHandler) publishPack(c *fiber.Ctx) error {
	var payload PackDto
	err := c.BodyParser(&payload)
	if err != nil {
		return err
	}

	port.log.Debug("-> publishPack", zap.Any("payload", payload))

	err = port.packService.PublishPack(payload.ModelNickName, payload.PackNumber)
	if err != nil {
		return err
	}

	return c.JSON(rest.ApiOkEmpty())
}

// ShowAccount godoc
// @Summary      Get Pack Info
// @Description  Get information about the pack
// @Tags         Pack
// @Accept       json
// @Produce      json
// @Param        modelNickName   path     string  true  "model nickname"
// @Param        packNumber   path     int  true  "pack number"
// @Success      200  {object}  rest.ApiResponse[[]dto.PackInfoDto]
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/pack/{modelNickName}/{packNumber}/info [get]
func (port *packFiberHandler) getPackInfo(c *fiber.Ctx) error {
	modelNickNameParam := c.Params("modelNickName")
	packNumberParam := c.Params("packNumber")

	packNumber, err := strconv.Atoi(packNumberParam)
	if err != nil {
		return err
	}

	var personIdRequester *string
	identity, err := port.securityService.GetIdentity(c)
	if err != nil {
		personIdRequester = nil
	} else {
		personIdRequester = &identity.PersonId
	}

	port.log.Debug("-> getPackInfo",
		zap.String("modelNickName", modelNickNameParam),
		zap.Int("packNumber", packNumber),
		zap.Any("personIdRequester", personIdRequester))

	info, err := port.packService.GetPackInfo(modelNickNameParam, packNumber, personIdRequester)
	if err != nil {
		return nil
	}

	return c.JSON(rest.ApiOk(info))
}

// ShowAccount godoc
// @Summary      Get Items From Pack
// @Description  Get all active items from the pack
// @Tags         Pack
// @Accept       json
// @Produce      json
// @Param        modelNickName   path     string  true  "model nickname"
// @Param        packNumber   path     int  true  "pack number"
// @Success      200  {object}  rest.ApiResponse[[]dto.PackItemDto]
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/pack/{modelNickName}/{packNumber}/items [get]
func (port *packFiberHandler) getItemsFromPack(c *fiber.Ctx) error {
	modelNickNameParam := c.Params("modelNickName")
	packNumberParam := c.Params("packNumber")

	packNumber, err := strconv.Atoi(packNumberParam)
	if err != nil {
		return err
	}

	var personIdRequester *string
	identity, err := port.securityService.GetIdentity(c)
	if err != nil {
		personIdRequester = nil
	} else {
		personIdRequester = &identity.PersonId
	}

	port.log.Debug("-> getItemsFromPack",
		zap.String("modelNickName", modelNickNameParam),
		zap.Int("packNumber", packNumber),
		zap.Any("personIdRequester", personIdRequester))

	items, err := port.packService.GetItemsFromPack(modelNickNameParam, packNumber, personIdRequester)
	if err != nil {
		return nil
	}

	return c.JSON(rest.ApiOkArray(items))
}

// ShowAccount godoc
// @Summary      Get Packs From Model
// @Description  Get all the packs from one model
// @Tags         Pack
// @Accept       json
// @Produce      json
// @Param        modelNickName   path     string  true  "model nickname"
// @Success      200  {object}  rest.ApiResponse[[]dto.PackDto]
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/pack/{modelNickName} [get]
func (port *packFiberHandler) getPacksFromModel(c *fiber.Ctx) error {
	modelNickNameParam := c.Params("modelNickName")

	var personIdRequester *string
	identity, err := port.securityService.GetIdentity(c)
	if err != nil {
		personIdRequester = nil
	} else {
		personIdRequester = &identity.PersonId
	}

	port.log.Debug("-> getPacksFromModel",
		zap.String("modelNickName", modelNickNameParam),
		zap.Any("personIdRequester", personIdRequester))

	items, err := port.packService.GetPacksFromModel(modelNickNameParam, personIdRequester)
	if err != nil {
		return nil
	}

	return c.JSON(rest.ApiOkArray(items))
}

// ShowAccount godoc
// @Summary      Edit Pack Title
// @Description  Edit the pack title by the model or admin
// @Tags         Pack
// @Accept       json
// @Produce      json
// @Param        data body EditPackTitleDto true "Payload Data"
// @Success      200  {object}  rest.ApiResponse[string]
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/pack/{modelNickName}/{packNumber}/title [post]
func (port *packFiberHandler) editPackTitle(c *fiber.Ctx) error {
	modelNickNameParam := c.Params("modelNickName")
	packNumberParam := c.Params("packNumber")
	var payload EditPackTitleDto
	err := c.BodyParser(&payload)
	if err != nil {
		return err
	}
	packNumber, err := strconv.Atoi(packNumberParam)
	if err != nil {
		return err
	}

	err = port.validate.Struct(payload)
	if err != nil {
		return err
	}

	port.log.Debug("-> editPackTitle", zap.Any("payload", payload))
	err = port.packService.EditPackTitle(modelNickNameParam, packNumber, payload.Title)
	if err != nil {
		return err
	}
	return c.JSON(rest.ApiOkEmpty())
}

// ShowAccount godoc
// @Summary      Edit Pack Description
// @Description  Edit the pack description by the model or admin
// @Tags         Pack
// @Accept       json
// @Produce      json
// @Param        data body EditPackDescriptionDto true "Payload Data"
// @Success      200  {object}  rest.ApiResponse[string]
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/pack/{modelNickName}/{packNumber}/description [post]
func (port *packFiberHandler) editPackDescription(c *fiber.Ctx) error {
	modelNickNameParam := c.Params("modelNickName")
	packNumberParam := c.Params("packNumber")
	var payload EditPackDescriptionDto
	err := c.BodyParser(&payload)
	if err != nil {
		return err
	}
	packNumber, err := strconv.Atoi(packNumberParam)
	if err != nil {
		return err
	}

	err = port.validate.Struct(payload)
	if err != nil {
		return err
	}

	port.log.Debug("-> editPackDescription", zap.Any("payload", payload))
	err = port.packService.EditPackDescription(modelNickNameParam, packNumber, payload.Description)
	if err != nil {
		return err
	}
	return c.JSON(rest.ApiOkEmpty())
}
