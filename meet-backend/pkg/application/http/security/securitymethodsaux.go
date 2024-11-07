package security

import (
	"github.com/erodriguezg/meet/pkg/application/http/rest/fiber/fiberidentity"
	"github.com/gofiber/fiber/v2"
)

func MustHavePermission(securityService HttpSecurityService, c *fiber.Ctx, permissionCode int) (*fiberidentity.FiberIdentity, error) {
	identity, err := securityService.GetIdentity(c)
	if err != nil {
		return nil, err
	}
	err = identity.MustHavePermission(permissionCode)
	if err != nil {
		return nil, err
	}
	return identity, nil
}
