package handler

import (
	"github.com/erodriguezg/meet/pkg/application/http/rest"
	"github.com/erodriguezg/meet/pkg/application/http/security"
	"github.com/erodriguezg/meet/pkg/core/dto"
	"github.com/erodriguezg/meet/pkg/core/service"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type chBnkAcctFiberHandler struct {
	chiliBankService service.ChiliBankAccountService
	securityService  security.HttpSecurityService
	log              *zap.Logger
}

func NewChiliBankAccountFiberHandler(
	chiliBankService service.ChiliBankAccountService,
	securityService security.HttpSecurityService,
	log *zap.Logger,
) FiberHandler {
	return &chBnkAcctFiberHandler{
		chiliBankService,
		securityService,
		log,
	}
}

// RegisterRoutes implements FiberHandler.
func (port *chBnkAcctFiberHandler) RegisterRoutes(fiberRouter *fiber.Router) {
	router := *fiberRouter
	group := router.Group("/chili-bank")
	group.Get("/banks", port.getBanks)
	group.Get("/account-types", port.getAccountTypes)
	group.Get("/:modelNickname/accounts", port.getModelAccounts)
	group.Post("/:modelNickname/accounts", port.saveModelAccount)
	group.Delete("/:modelNickname/:accountId", port.deleteModelAccount)
}

// ShowAccount godoc
// @Summary      Get Banks
// @Description  Get Chili Banks
// @Tags         Chili Bank
// @Accept       json
// @Produce      json
// @Success      200  {object}  rest.ApiResponse[[]string]
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/chili-bank/banks [get]
func (port *chBnkAcctFiberHandler) getBanks(c *fiber.Ctx) error {
	port.log.Debug("-> getBanks")
	banks, err := port.chiliBankService.GetBanks()
	if err != nil {
		return err
	}
	return c.JSON(rest.ApiOkArray[string](banks))
}

// ShowAccount godoc
// @Summary      Get Account Types
// @Description  Get accounts types for chili banks
// @Tags         Chili Bank
// @Accept       json
// @Produce      json
// @Success      200  {object}  rest.ApiResponse[[]string]
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/chili-bank/account-types [get]
func (port *chBnkAcctFiberHandler) getAccountTypes(c *fiber.Ctx) error {
	port.log.Debug("-> getAccountTypes")
	accountTypes, err := port.chiliBankService.GetAccountTypes()
	if err != nil {
		return err
	}
	return c.JSON(rest.ApiOkArray[string](accountTypes))
}

// ShowAccount godoc
// @Summary      Get Model Accounts
// @Description  Get Model Chili Banks Accounts
// @Tags         Chili Bank
// @Accept       json
// @Produce      json
// @Param        modelNickname path string true "model nickname"
// @Success      200  {object}  rest.ApiResponse[[]dto.ChiliBankAccountDTO]
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/chili-bank/{modelNickname}/accounts [get]
func (port *chBnkAcctFiberHandler) getModelAccounts(c *fiber.Ctx) error {
	modelNicknameParam := c.Params("modelNickname")
	_, err := security.MustHavePermissionToEditModel(port.securityService, c, modelNicknameParam)
	if err != nil {
		return err
	}
	accounts, err := port.chiliBankService.GetAccounts(modelNicknameParam)
	if err != nil {
		return err
	}
	return c.JSON(rest.ApiOkArray(accounts))
}

// ShowAccount godoc
// @Summary      Save Model Account
// @Description  Insert or update a chili bank account for the model
// @Tags         Chili Bank
// @Accept       json
// @Produce      json
// @Param        modelNickname path string true "model nickname"
// @Param        data body dto.ChiliBankAccountDTO true "The chili bank account data"
// @Success      200  {object}  rest.ApiResponse[dto.ChiliBankAccountDTO]
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/chili-bank/{modelNickname}/accounts [post]
func (port *chBnkAcctFiberHandler) saveModelAccount(c *fiber.Ctx) error {
	var payload dto.ChiliBankAccountDTO
	err := c.BodyParser(&payload)
	if err != nil {
		return err
	}
	modelNicknameParam := c.Params("modelNickname")
	_, err = security.MustHavePermissionToEditModel(port.securityService, c, modelNicknameParam)
	if err != nil {
		return err
	}
	accountOut, err := port.chiliBankService.Save(modelNicknameParam, payload)
	if err != nil {
		return err
	}
	return c.JSON(rest.ApiOk(&accountOut))
}

// ShowAccount godoc
// @Summary      Delete Model Account
// @Description  Delete a chili bank account of the model
// @Tags         Chili Bank
// @Accept       json
// @Produce      json
// @Param        modelNickname path string true "model nickname"
// @Param        accountId path string true "chili bank account id"
// @Param        data body dto.ChiliBankAccountDTO true "The chili bank account data"
// @Success      200  {object}  rest.ApiResponse[string]
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/chili-bank/{modelNickname}/{accountId} [delete]
func (port *chBnkAcctFiberHandler) deleteModelAccount(c *fiber.Ctx) error {
	modelNicknameParam := c.Params("modelNickname")
	accountIdParam := c.Params("accountId")
	_, err := security.MustHavePermissionToEditModel(port.securityService, c, modelNicknameParam)
	if err != nil {
		return err
	}
	err = port.chiliBankService.Delete(modelNicknameParam, accountIdParam)
	if err != nil {
		return err
	}
	return c.JSON(rest.ApiOkEmpty())
}
