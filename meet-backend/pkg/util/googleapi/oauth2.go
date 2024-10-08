package googleapi

type OAuth2TokenRequest struct {
	Code         string
	ClientId     string
	ClientSecret string
	RedirectUri  string
	GrantType    string
}

type OAuth2TokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int32  `json:"expires_in"`
	IdToken      string `json:"id_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
}

type OAuth2Api interface {
	ValidateAndGetToken(request OAuth2TokenRequest) (OAuth2TokenResponse, error)
}
