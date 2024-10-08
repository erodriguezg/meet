package googleapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	oauth2BaseUrl  string = "https://oauth2.googleapis.com"
	oauth2TokenUrl string = oauth2BaseUrl + "/token"
)

type netHttpOauth2Api struct {
	httpClient *http.Client
}

func NewNetHttpOauth2Api(httpClient *http.Client) OAuth2Api {
	return &netHttpOauth2Api{httpClient}
}

func (port *netHttpOauth2Api) ValidateAndGetToken(input OAuth2TokenRequest) (OAuth2TokenResponse, error) {

	var output OAuth2TokenResponse

	data := url.Values{}
	data.Set("code", input.Code)
	data.Set("client_id", input.ClientId)
	data.Set("client_secret", input.ClientSecret)
	data.Set("redirect_uri", input.RedirectUri)
	data.Set("grant_type", input.GrantType)

	request, err := http.NewRequest(http.MethodPost, oauth2TokenUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return output, err
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := port.httpClient.Do(request)
	if err != nil {
		return output, err
	}

	defer response.Body.Close()

	responseRaw, err := io.ReadAll(response.Body)
	if err != nil {
		return output, err
	}

	if response.StatusCode >= 200 && response.StatusCode < 300 {
		err := json.Unmarshal(responseRaw, &output)
		return output, err
	}

	return output, fmt.Errorf("network error: %s", responseRaw)
}
