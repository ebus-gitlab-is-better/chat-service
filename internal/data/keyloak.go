package data

import (
	"chat-service/internal/conf"
	"context"

	"github.com/Nerzal/gocloak/v13"
	"github.com/go-kratos/kratos/v2/log"
)

type KeycloakAPI struct {
	conf   *conf.Data
	client *gocloak.GoCloak
	logger *log.Helper
}

func NewKeyCloakAPI(conf *conf.Data, client *gocloak.GoCloak, logger log.Logger) *KeycloakAPI {
	return &KeycloakAPI{
		conf:   conf,
		client: client,
		logger: log.NewHelper(logger),
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
