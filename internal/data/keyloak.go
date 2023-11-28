package data

import (
	"chat-service/internal/conf"
	"context"

	"github.com/Nerzal/gocloak/v13"
)

type KeycloakAPI struct {
	conf   *conf.Data
	client *gocloak.GoCloak
}

func NewKeyCloakAPI(conf *conf.Data, client *gocloak.GoCloak) *KeycloakAPI {
	return &KeycloakAPI{
		conf:   conf,
		client: client,
	}
}

func (api *KeycloakAPI) CheckToken(accessToken string) (*gocloak.IntroSpectTokenResult, error) {
	return api.client.RetrospectToken(
		context.TODO(),
		accessToken,
		api.conf.Keycloak.ClientId,
		api.conf.Keycloak.ClientSecret,
		api.conf.Keycloak.Realm)
}

func (api *KeycloakAPI) GetUserInfo(accessToken string) (*gocloak.UserInfo, error) {
	return api.client.GetUserInfo(
		context.TODO(),
		accessToken,
		api.conf.Keycloak.Realm)
}
