package handler

import (
	"strconv"

	"github.com/erodriguezg/meet/pkg/application/http/rest"
	"github.com/erodriguezg/meet/pkg/application/http/security"
	"github.com/erodriguezg/meet/pkg/core/dto"
	"github.com/erodriguezg/meet/pkg/core/service"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type packPMFiberHandler struct {
	packPaymentMethodService service.PackPaymentMethodService
	securityService          security.HttpSecurityService
	log                      *zap.Logger
}

func NewPackPaymentMethodFiberHandler(
	packPaymentMethodService service.PackPaymentMethodService,
	securityService security.HttpSecurityService,
	log *zap.Logger,
) FiberHandler {
	return &packPMFiberHandler{
		packPaymentMethodService,
		securityService,
		log,
	}
}

func (port *packPMFiberHandler) RegisterRoutes(fiberRouter *fiber.Router) {
	router := *fiberRouter
	group := router.Group("/pack-payment-methods")
	group.Get("/:modelNickname/:packNumber", port.getPackPaymentsMethods)
	group.Post("/:modelNickname/:packNumber", port.savePackPaymentMethods)
}

// privates

// ShowAccount godoc
// @Summary      Gets Pack Payment Methods
// @Description  Gets the payments methods for one pack
// @Tags         Pack Payment Methods
// @Accept       json
// @Produce      json
// @Param        modelNickname path string true "model nickname"
// @Param        packNumber path int true "pack number"
// @Success      200  {object}  rest.ApiResponse[dto.PackPaymentMethodDTO]
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/pack-payment-methods/{modelNickname}/{packNumber} [get]
func (port *packPMFiberHandler) getPackPaymentsMethods(c *fiber.Ctx) error {
	modelNicknameParam := c.Params("modelNickname")
	packNumber, err := strconv.Atoi(c.Params("packNumber"))
	if err != nil {
		return err
	}

	_, err = security.MustHavePermissionToEditModel(port.securityService, c, modelNicknameParam)
	if err != nil {
		return err
	}

	port.log.Debug("-> getPackPaymentsMethods", zap.String("modelNickName", modelNicknameParam), zap.Int("packNumber", packNumber))
	existingPackPM, err := port.packPaymentMethodService.GetFromPack(modelNicknameParam, packNumber)
	if err != nil {
		return err
	}

	return c.JSON(rest.ApiOk(&existingPackPM))
}

// ShowAccount godoc
// @Summary      Save Pack Payment Methods
// @Description  Insert or update the payment methods for one pack
// @Tags         Pack Payment Methods
// @Accept       json
// @Produce      json
// @Param        modelNickname path string true "model nickname"
// @Param        data body dto.PackPaymentMethodDTO true "The payment methods data"
// @Success      200  {object}  rest.ApiResponse[dto.PackPaymentMethodDTO]
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/pack-payment-methods/{modelNickname}/{packNumber} [post]
func (port *packPMFiberHandler) savePackPaymentMethods(c *fiber.Ctx) error {
	modelNicknameParam := c.Params("modelNickname")
	packNumber, err := strconv.Atoi(c.Params("packNumber"))
	if err != nil {
		return err
	}

	_, err = security.MustHavePermissionToEditModel(port.securityService, c, modelNicknameParam)
	if err != nil {
		return err
	}

	var payload dto.PackPaymentMethodDTO
	err = c.BodyParser(&payload)
	if err != nil {
		return err
	}

	port.log.Debug("-> savePackPaymentMethods", zap.String("modelNickName", modelNicknameParam), zap.Int("packNumber", packNumber), zap.Any("payload", payload))
	editedPackPM, err := port.packPaymentMethodService.Save(modelNicknameParam, packNumber, payload)
	if err != nil {
		return err
	}

	return c.JSON(rest.ApiOk(&editedPackPM))
}
