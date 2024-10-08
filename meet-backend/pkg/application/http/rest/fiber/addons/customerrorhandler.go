package addons

import (
	"github.com/erodriguezg/meet/pkg/application/http/rest"
	"github.com/erodriguezg/meet/pkg/application/http/rest/fiber/fiberidentity"
	"github.com/erodriguezg/meet/pkg/core/exception"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type customFiberErrorHandler struct {
	log *zap.Logger
}

func NewCustomFiberErrorHandler(log *zap.Logger) FiberErrorHandler {
	return &customFiberErrorHandler{log}
}

func (port *customFiberErrorHandler) CustomFiberErrorHandler(ctx *fiber.Ctx, err error) error {
	if businessException, ok := err.(*exception.BusinessException); ok {
		port.log.Debug("business exception: ", zap.Error(businessException))
		apiResponse := rest.ApiBusinessException(businessException)
		statusCode := fiber.StatusOK
		return ctx.Status(statusCode).JSON(apiResponse)
	} else if accessDeniedError, ok := err.(*fiberidentity.FiberAccessDeniedError); ok {
		port.log.Debug("access denied error: ", zap.Error(accessDeniedError))
		apiResponse := rest.ApiAccessDenied()
		statusCode := fiber.StatusUnauthorized
		return ctx.Status(statusCode).JSON(apiResponse)
	} else {
		port.log.Error("api error: ", zap.Error(err))
		return fiber.DefaultErrorHandler(ctx, err)
	}
}
