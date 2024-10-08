package addons

import "github.com/gofiber/fiber/v2"

type FiberErrorHandler interface {
	CustomFiberErrorHandler(ctx *fiber.Ctx, err error) error
}
