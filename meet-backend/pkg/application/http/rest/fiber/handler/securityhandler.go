package handler

import (
	"github.com/erodriguezg/meet/pkg/application/http/rest"
	"github.com/erodriguezg/meet/pkg/application/http/security"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type securityHandler struct {
	securityService security.HttpSecurityService
	log             *zap.Logger
}

func NewSecurityHandler(
	securityService security.HttpSecurityService,
	log *zap.Logger) FiberHandler {
	return &securityHandler{
		securityService,
		log,
	}
}

// RegisterRoutes implements FiberHandler
func (port *securityHandler) RegisterRoutes(fiberRouter *fiber.Router) {
	router := *fiberRouter
	group := router.Group("/security")
	group.Get("/redirect-to-openid-login-url", port.redirectToOpenIdLoginUrl)
	group.Get("/login-url", port.getLoginUrl)
	group.Post("/token", port.getToken)
	group.Get("/identity", port.getIdentity)
}

// ShowAccount godoc
// @Summary      Redirect To Login
// @Description  Redirect To OpenId Login Url
// @Tags         Security
// @Accept       json
// @Param        state  query     string  false  "state from frontend"
// @Success      301  {object}  string
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/security/redirect-to-openid-login-url [get]
func (port *securityHandler) redirectToOpenIdLoginUrl(c *fiber.Ctx) error {
	port.log.Debug("-> getOpenIdLoginUrl")
	state := c.Query("state")
	loginUrl := port.securityService.GetOpenIdLoginUrl(state)
	return c.Redirect(loginUrl)
}

// ShowAccount godoc
// @Summary      Get Login Url
// @Description  Get the login url of the sso
// @Tags         Security
// @Accept       json
// @Param        state  query     string  false  "state from frontend"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/security/login-url [get]
func (port *securityHandler) getLoginUrl(c *fiber.Ctx) error {
	port.log.Debug("-> getLoginUrl")
	state := c.Query("state")
	loginUrl := port.securityService.GetOpenIdLoginUrl(state)
	return c.JSON(map[string]string{
		"loginUrl": loginUrl,
	})
}

// ShowAccount godoc
// @Summary      Get Token
// @Description  Get the token from OpenId Response Data
// @Tags         Security
// @Accept       json
// @Param        code  query     string  false  "code from openid"
// @Param        state  query     string  false  "state from openid"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/security/token [post]
func (port *securityHandler) getToken(c *fiber.Ctx) error {
	port.log.Debug("-> getToken")
	code := c.Query("code")
	state := c.Query("state")
	token, err := port.securityService.GetToken(code, state)
	if err != nil {
		return err
	}
	return c.JSON(map[string]string{
		"jwt": token,
	})
}

// ShowAccount godoc
// @Summary      Get Identity
// @Description  Get the authenticated identity
// @Tags         Security
// @Accept       json
// @Success      200  {object}  rest.ApiResponse[fiberidentity.FiberIdentity]
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /v1/security/identity [get]
func (port *securityHandler) getIdentity(c *fiber.Ctx) error {
	port.log.Debug("-> getIdentity")
	identity, err := port.securityService.GetIdentity(c)
	if err != nil {
		return err
	}
	return c.JSON(rest.ApiOk(identity))
}
