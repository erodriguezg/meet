package fiberidentity

import (
	"crypto/rsa"
	"fmt"
	"strings"

	"github.com/erodriguezg/meet/pkg/core/service"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type FiberIdentity struct {
	PersonId         string `json:"personId"`
	Email            string `json:"email"`
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	ProfileCode      int    `json:"profileCode"`
	ProfileName      string `json:"profileName"`
	PermissionsCodes []int  `json:"permissionsCodes"`
}

type FiberAccessDeniedError struct {
	ErrorDetail error
}

// Error implements error
func (port *FiberAccessDeniedError) Error() string {
	return fmt.Sprintf("Access Denied: %v", port.ErrorDetail)
}

func NewAccessDeniedError(errorDetail error) error {
	return &FiberAccessDeniedError{errorDetail}
}

type FiberIdentityUtil struct {
	personService  service.PersonService
	profileService service.ProfileService
	modelService   service.ModelService
	rsaPublicKey   *rsa.PublicKey
}

func NewFiberIdentityUtil(
	personService service.PersonService,
	profileService service.ProfileService,
	modelService service.ModelService,
	rsaPublicKeyBytes []byte,
) *FiberIdentityUtil {

	rsaPublicKey, err := jwt.ParseRSAPublicKeyFromPEM(rsaPublicKeyBytes)
	if err != nil {
		panic(fmt.Sprintf("error getting jwt public key from pem: %s", err.Error()))
	}

	return &FiberIdentityUtil{
		personService,
		profileService,
		modelService,
		rsaPublicKey,
	}
}

func (port *FiberIdentityUtil) GetIdentity(c *fiber.Ctx) (*FiberIdentity, error) {

	authHeader := c.GetReqHeaders()["Authorization"]

	tokenHeader := ""

	if len(authHeader) > 0 {
		tokenHeader = strings.ReplaceAll(authHeader[0], "Bearer ", "")
	}

	token, err := jwt.Parse(tokenHeader, func(token *jwt.Token) (interface{}, error) {
		// verify if is a rsa sign method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("jwt parse error: %v", token.Header["alg"])
		}
		return port.rsaPublicKey, nil
	})

	if err != nil {
		return nil, NewAccessDeniedError(fmt.Errorf("error processing authorization header jwt: %w", err))
	}

	if token == nil || !token.Valid {
		return nil, NewAccessDeniedError(fmt.Errorf("no valid authorization token found"))
	}

	claims := token.Claims.(jwt.MapClaims)
	email := claims["email"].(string)

	person, err := port.personService.FindByEmail(email)
	if err != nil {
		return nil, NewAccessDeniedError(fmt.Errorf("error at finding by email on identity creation: %w", err))
	}

	if person == nil {
		return nil, NewAccessDeniedError(fmt.Errorf("no person was found with email: %s on identity creation", email))
	}

	profile, err := port.profileService.FindByCode(person.ProfileCode)
	if err != nil {
		return nil, NewAccessDeniedError(fmt.Errorf("error at finding profile on identity creation: %w", err))
	}

	if profile == nil {
		return nil, NewAccessDeniedError(fmt.Errorf("no profile was found with code: %d on identity creation", person.ProfileCode))
	}

	identity := FiberIdentity{
		PersonId:         person.Id.Hex(),
		Email:            person.Email,
		FirstName:        person.FirstName,
		LastName:         person.LastName,
		ProfileCode:      profile.Code,
		ProfileName:      profile.Name,
		PermissionsCodes: profile.PermissionsCodes,
	}

	return &identity, nil
}

func (port *FiberIdentity) MustHaveProfile(profileCode int) error {
	if profileCode == port.ProfileCode {
		return nil
	}
	return NewAccessDeniedError(fmt.Errorf("don't have profile: %d", profileCode))
}

func (port *FiberIdentity) MustHavePermission(permissionCode int) error {
	for _, permission := range port.PermissionsCodes {
		if permission == permissionCode {
			return nil
		}
	}
	return NewAccessDeniedError(fmt.Errorf("don't have permission: %d", permissionCode))
}

func (port *FiberIdentity) HasProfile(profileCode int) bool {
	return profileCode == port.ProfileCode
}

func (port *FiberIdentity) HasPermission(permissionCode int) bool {
	for _, permission := range port.PermissionsCodes {
		if permission == permissionCode {
			return true
		}
	}
	return false
}
