package openid

import (
	"fmt"
	"net/url"

	"github.com/erodriguezg/meet/pkg/util/googleapi"
	"github.com/erodriguezg/meet/pkg/util/jwtutil"
	"go.uber.org/zap"
)

const (
	googleAccountsHostname string = "https://accounts.google.com/o/oauth2/v2/auth"
)

type googleOpenIdService struct {
	config    OpenIdConfig
	oauth2Api googleapi.OAuth2Api
	jwtUtil   jwtutil.JwtUtil
	log       *zap.Logger
}

func NewGoogleOpenIdService(config OpenIdConfig, oauth2Api googleapi.OAuth2Api, jwtUtil jwtutil.JwtUtil, log *zap.Logger) OpenIdService {
	return &googleOpenIdService{config, oauth2Api, jwtUtil, log}
}

// GetLoginUrl implements OpenIdService
func (port *googleOpenIdService) GetLoginUrl(state string) string {
	urlTarget := googleAccountsHostname + "?"
	urlTarget += "response_type=" + port.config.ResponseType
	urlTarget += "&client_id=" + port.config.ClientId
	urlTarget += "&scope=" + url.QueryEscape(port.config.Scope)
	urlTarget += "&redirect_uri=" + url.QueryEscape(port.config.RedirectUri)
	urlTarget += "&state=" + state
	urlTarget += "&prompt=consent"
	return urlTarget
}

// ProcessCallback implements OpenIdService
func (port *googleOpenIdService) ProcessCallback(code string, state string, stateValidator func(string) bool) (OpenIdUser, error) {
	var openIdUser OpenIdUser
	stateValidatorSuccess := stateValidator(state)

	if !stateValidatorSuccess {
		return openIdUser, fmt.Errorf("invalid google openid state: '%s'", state)
	}

	tokenRequest := googleapi.OAuth2TokenRequest{
		Code:         code,
		ClientId:     port.config.ClientId,
		ClientSecret: port.config.ClientSecret,
		RedirectUri:  port.config.RedirectUri,
		GrantType:    "authorization_code",
	}

	port.log.Debug("before call api", zap.Any("token request", tokenRequest))

	tokenResponse, err := port.oauth2Api.ValidateAndGetToken(tokenRequest)
	if err != nil {
		return openIdUser, err
	}

	port.log.Debug("api response", zap.Any("token", tokenResponse))

	jwt, err := port.jwtUtil.ParseWithoutKey(tokenResponse.IdToken)
	if err != nil {
		return openIdUser, err
	}

	if jwt["email"] != nil {
		openIdUser.Email = jwt["email"].(string)
	}
	if jwt["given_name"] != nil {
		openIdUser.FirstName = jwt["given_name"].(string)
	}
	if jwt["family_name"] != nil {
		openIdUser.LastName = jwt["family_name"].(string)
	}
	if jwt["picture"] != nil {
		openIdUser.PictureUrl = jwt["picture"].(string)
	}
	if jwt["locale"] != nil {
		openIdUser.Locale = jwt["locale"].(string)
	}

	return openIdUser, nil
}
