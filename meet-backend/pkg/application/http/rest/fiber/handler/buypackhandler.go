package handler

import (
	"fmt"

	"github.com/erodriguezg/meet/pkg/application/http/rest"
	"github.com/erodriguezg/meet/pkg/application/http/rest/fiber/fiberidentity"
	"github.com/erodriguezg/meet/pkg/application/http/security"
	"github.com/erodriguezg/meet/pkg/core/service"
	fiber "github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type buyPackHandler struct {
	buyPackService  service.BuyPackService
	securityService security.HttpSecurityService
	log             *zap.Logger
}

type BuyPackDetailsRequest struct {
	ModelNickName string `json:"modelNickName"`
	PackNumber    int    `json:"packNumber"`
}

type BuyPackCreateOrderRequest struct {
	PersonId      string `json:"personId"`
	ModelNickName string `json:"modelNickName"`
	PackNumber    int    `json:"packNumber"`
}

type BuyPackCreateOrderResponse struct {
	OrderId string `json:"orderId"`
}

type BuyPackCapturePaymentRequest struct {
	OrderId string `json:"orderId"`
}

func NewBuyPackHandler(
	buyPackService service.BuyPackService,
	securityService security.HttpSecurityService,
	log *zap.Logger,
) FiberHandler {
	return &buyPackHandler{
		buyPackService,
		securityService,
		log,
	}
}

func (port *buyPackHandler) RegisterRoutes(fiberRouter *fiber.Router) {
	router := *fiberRouter
	group := router.Group("/buy-pack")
	group.Get("/info", port.getPaymentClientData)
	group.Post("/details", port.getPackBuyDetails)
	group.Post("/create-order", port.createBuyPackOrder)
	group.Post("/capture-payment", port.capturePackPayment)
}

// privates

// ShowAccount godoc
// @Summary      Get Payment Client Data
// @Description  Get config client data for payment system
// @Tags         BuyPack
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/buy-pack/info [get]
func (port *buyPackHandler) getPaymentClientData(c *fiber.Ctx) error {
	identity, err := port.securityService.GetIdentity(c)
	if err != nil {
		return err
	}
	if identity == nil {
		return fmt.Errorf("identity is nil")
	}
	clientData, err := port.buyPackService.GetPaymentClientData()
	if err != nil {
		return err
	}
	return c.JSON(rest.ApiOk(&clientData))
}

// ShowAccount godoc
// @Summary      Get Buy Pack Details
// @Description  Get info required for buy the pack
// @Tags         BuyPack
// @Accept       json
// @Produce      json
// @Param        data body BuyPackDetailsRequest true "details buy pack dto"
// @Success      200  {object}  dto.PackBuyDetailDto
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/buy-pack/detail [post]
func (port *buyPackHandler) getPackBuyDetails(c *fiber.Ctx) error {
	identity, err := port.securityService.GetIdentity(c)
	if err != nil {
		return err
	}
	if identity == nil {
		return fmt.Errorf("identity is nil")
	}

	var payload BuyPackDetailsRequest
	err = c.BodyParser(&payload)
	if err != nil {
		return err
	}
	buyPackDto, err := port.buyPackService.GetPackBuyDetails(payload.ModelNickName, payload.PackNumber)
	if err != nil {
		return err
	}
	return c.JSON(rest.ApiOk(&buyPackDto))
}

// ShowAccount godoc
// @Summary      Create Buy Pack Order
// @Description  Generate a new order for buy a pack
// @Tags         BuyPack
// @Accept       json
// @Produce      json
// @Param        data body BuyPackCreateOrderRequest true "Create Order Data"
// @Success      200  {object}  BuyPackCreateOrderResponse
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/buy-pack/create-order [post]
func (port *buyPackHandler) createBuyPackOrder(c *fiber.Ctx) error {
	identity, err := port.securityService.GetIdentity(c)
	if err != nil {
		return err
	}
	if identity == nil {
		return fmt.Errorf("identity is nil")
	}
	var payload BuyPackCreateOrderRequest
	err = c.BodyParser(&payload)
	if err != nil {
		return err
	}

	if identity.PersonId != payload.PersonId {
		return fiberidentity.NewAccessDeniedError(fmt.Errorf("incompatible personId with session data"))
	}

	orderId, err := port.buyPackService.CreateBuyPackOrder(payload.PersonId, payload.ModelNickName, payload.PackNumber)
	if err != nil {
		return err
	}
	responsePayload := BuyPackCreateOrderResponse{
		OrderId: orderId,
	}
	return c.JSON(rest.ApiOk(&responsePayload))
}

// ShowAccount godoc
// @Summary      Capture Pack Payment
// @Description  Capture a Payment to a Pack
// @Tags         BuyPack
// @Accept       json
// @Produce      json
// @Param        data body BuyPackCapturePaymentRequest true "Capture Payment Data"
// @Success      200  {object}  string
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/buy-pack/capture-payment [post]
func (port *buyPackHandler) capturePackPayment(c *fiber.Ctx) error {
	identity, err := port.securityService.GetIdentity(c)
	if err != nil {
		return err
	}
	if identity == nil {
		return fmt.Errorf("identity is nil")
	}

	var payload BuyPackCapturePaymentRequest
	err = c.BodyParser(&payload)
	if err != nil {
		return err
	}

	err = port.buyPackService.CapturePackPayment(payload.OrderId)
	if err != nil {
		return err
	}
	return c.JSON(rest.ApiOkEmpty())
}
