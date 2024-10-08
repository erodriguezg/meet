package handler

import "github.com/gofiber/fiber/v2"

type FiberHandler interface {
	RegisterRoutes(fiberGroup *fiber.Router)
}
