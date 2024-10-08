package openid

import "fmt"

type OpenIdUser struct {
	Email      string
	FirstName  string
	LastName   string
	PictureUrl string
	Locale     string
}

func (model *OpenIdUser) GetFullName() string {
	return fmt.Sprintf("%s %s", model.FirstName, model.LastName)
}

type OpenIdConfig struct {
	ResponseType string
	Scope        string
	ClientId     string
	ClientSecret string
	RedirectUri  string
}

type OpenIdService interface {
	GetLoginUrl(state string) string
	ProcessCallback(code string, state string, stateValidator func(string) bool) (OpenIdUser, error)
}
