package config

import (
	"net/http"

	"github.com/erodriguezg/meet/pkg/util/googleapi"
	"github.com/erodriguezg/meet/pkg/util/jwtutil"
)

var (
	httpClient      *http.Client
	jwtUtil         jwtutil.JwtUtil
	googleOAuth2Api googleapi.OAuth2Api
)

func configUtils() {
	httpClient = configHttpClient()
	jwtUtil = configGolangJwtUtil()
	googleOAuth2Api = configGoogleOAuth2Api()
}

func configHttpClient() *http.Client {
	return &http.Client{}
}

func configGolangJwtUtil() jwtutil.JwtUtil {
	return jwtutil.NewGolangJwtUtil()
}

func configGoogleOAuth2Api() googleapi.OAuth2Api {
	panicIfAnyNil(httpClient)
	return googleapi.NewNetHttpOauth2Api(httpClient)
}
