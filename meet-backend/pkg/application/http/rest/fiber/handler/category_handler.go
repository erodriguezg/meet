package handler

import (
	"github.com/erodriguezg/meet/pkg/application/http/security"
	"github.com/erodriguezg/meet/pkg/core/service"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type categoryFiberHandler struct {
	categoryService service.CategoryService
	securityService security.HttpSecurityService
	log             *zap.Logger
}

func NewCategoryFiberHandler(
	categoryService service.CategoryService,
	securityService security.HttpSecurityService,
	log *zap.Logger,
) FiberHandler {
	return &categoryFiberHandler{
		categoryService,
		securityService,
		log,
	}
}

// RegisterRoutes implements FiberHandler.
func (port *categoryFiberHandler) RegisterRoutes(fiberRouter *fiber.Router) {
	router := *fiberRouter
	group := router.Group("/categories")
	group.Get("/all-tree", port.getAllCategoriesTree)
	group.Post("/save", port.saveCategory)
	group.Delete("/delete", port.deleteCategory)
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
func (port *categoryFiberHandler) getAllCategoriesTree(c *fiber.Ctx) error {
	panic("unimplemented")
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
func (port *categoryFiberHandler) saveCategory(c *fiber.Ctx) error {
	panic("unimplemented")
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
func (port *categoryFiberHandler) deleteCategory(c *fiber.Ctx) error {
	panic("unimplemented")
}
