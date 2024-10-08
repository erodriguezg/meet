package security

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/erodriguezg/meet/pkg/application/http/rest/fiber/fiberidentity"
	"github.com/erodriguezg/meet/pkg/core/domain"
	"github.com/erodriguezg/meet/pkg/core/exception"
	"github.com/erodriguezg/meet/pkg/core/service"
	"github.com/erodriguezg/meet/pkg/util/openid"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type DefaultHttpSecurityService struct {
	openIdService     openid.OpenIdService
	personService     service.PersonService
	fiberIdentityUtil *fiberidentity.FiberIdentityUtil
	rsaPrivateKey     *rsa.PrivateKey
	statePassPhrase   string
}

func NewDefaultHttpSecurityService(
	openIdService openid.OpenIdService,
	personService service.PersonService,
	fiberIdentityUtil *fiberidentity.FiberIdentityUtil,
	rsaPrivateKeyBytes []byte,
	statePassPhrase string) HttpSecurityService {

	rsaPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM(rsaPrivateKeyBytes)
	if err != nil {
		panic(fmt.Sprintf("error getting jwt private key from pem: %s", err.Error()))
	}

	return &DefaultHttpSecurityService{
		openIdService,
		personService,
		fiberIdentityUtil,
		rsaPrivateKey,
		statePassPhrase,
	}
}

func (port *DefaultHttpSecurityService) GetOpenIdLoginUrl(state string) string {
	return port.openIdService.GetLoginUrl(state)
}

func (port *DefaultHttpSecurityService) GetToken(code string, state string) (string, error) {

	stateFunc := func(s string) bool {
		return true
	}

	openIdUser, err := port.openIdService.ProcessCallback(code, state, stateFunc)
	if err != nil {
		return "", fmt.Errorf("error at processing openid: %w", err)
	}

	person, err := port.personService.FindByEmail(openIdUser.Email)
	if err != nil {
		return "", fmt.Errorf("error at finding person by email: %w", err)
	}

	if person == nil {
		newPerson := domain.Person{
			Email:       openIdUser.Email,
			FirstName:   openIdUser.FirstName,
			LastName:    openIdUser.LastName,
			ProfileCode: domain.ProfileCodeUser,
		}
		person, err = port.personService.Save(newPerson)
		if err != nil {
			return "", fmt.Errorf("error at saving first time person: %w", err)
		}
	}

	if !person.Active {
		return "", exception.NewPersonIsNotActiveException(person)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"email":       person.Email,
		"firstName":   person.FirstName,
		"lastName":    person.LastName,
		"profileCode": person.ProfileCode,
		"exp":         time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(port.rsaPrivateKey)
}

func (port *DefaultHttpSecurityService) GetIdentity(c *fiber.Ctx) (*fiberidentity.FiberIdentity, error) {
	identity, err := port.fiberIdentityUtil.GetIdentity(c)
	if err != nil {
		return nil, err
	}
	return identity, nil
}

func (port *DefaultHttpSecurityService) MustHavePermission(permissionCode int, c *fiber.Ctx) error {
	identity, err := port.fiberIdentityUtil.GetIdentity(c)
	if err != nil {
		return err
	}
	return identity.MustHavePermission(permissionCode)
}

func (port *DefaultHttpSecurityService) MustHaveProfile(profileCode int, c *fiber.Ctx) error {
	identity, err := port.fiberIdentityUtil.GetIdentity(c)
	if err != nil {
		return err
	}
	return identity.MustHaveProfile(profileCode)
}
