package tokensvc

import (
	"context"
	"errors"

	"github.com/datpham/user-service-ms/config"
	"github.com/datpham/user-service-ms/internal/client/oauth"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	GoogleOAuthEmailScope   = "https://www.googleapis.com/auth/userinfo.email"
	GoogleOAuthProfileScope = "https://www.googleapis.com/auth/userinfo.profile"
)

type OAuthService struct {
	config      *oauth2.Config
	oauthClient *oauth.OauthClient
}

func NewOAuthService(appConfig *config.Config, oauthClient *oauth.OauthClient) *OAuthService {
	config := &oauth2.Config{
		ClientID:     appConfig.OAuth.ClientID,
		ClientSecret: appConfig.OAuth.ClientSecret,
		RedirectURL:  appConfig.OAuth.RedirectURL,
		Scopes:       []string{GoogleOAuthEmailScope, GoogleOAuthProfileScope},
		Endpoint:     google.Endpoint,
	}

	return &OAuthService{config, oauthClient}
}

func (s *OAuthService) GetGoogleAuthUrl() string {
	return s.config.AuthCodeURL("state")
}

func (s *OAuthService) GetGoogleAccessToken(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := s.config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *OAuthService) GetGoogleUserInfo(ctx context.Context, accessToken string) (map[string]interface{}, error) {
	userInfo, err := s.oauthClient.GetGoogleUserInfo(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}

func (s *OAuthService) VerifyGoogleState(state string) error {
	if state != "state" {
		return errors.New("invalid state")
	}

	return nil
}
