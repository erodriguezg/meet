package security

import (
	"github.com/erodriguezg/meet/pkg/application/http/rest/fiber/fiberidentity"
	"github.com/erodriguezg/meet/pkg/core/domain"
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

func MustHavePermissionToEditModel(securityService HttpSecurityService, c *fiber.Ctx, modelNickName string) (*fiberidentity.FiberIdentity, error) {
	identity, err := securityService.GetIdentity(c)
	if err != nil {
		return nil, err
	}
	if identity.ModelNickName == nil || *identity.ModelNickName != modelNickName {
		err = identity.MustHavePermission(domain.PermissionCodeEditAllModels)
		if err != nil {
			return nil, err
		}
	}
	return identity, nil
}
