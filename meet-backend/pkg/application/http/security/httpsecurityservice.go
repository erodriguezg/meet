package security

import (
	"github.com/erodriguezg/meet/pkg/application/http/rest/fiber/fiberidentity"
	"github.com/gofiber/fiber/v2"
)

type HttpSecurityService interface {
	GetOpenIdLoginUrl(state string) string
	GetToken(code string, state string) (string, error)
	GetIdentity(c *fiber.Ctx) (*fiberidentity.FiberIdentity, error)
	MustHavePermission(permissionCode int, c *fiber.Ctx) error
	MustHaveProfile(profileCode int, c *fiber.Ctx) error
}
