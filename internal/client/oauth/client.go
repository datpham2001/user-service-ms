package oauth

import (
	"context"
	"fmt"

	"github.com/datpham/user-service-ms/internal/pkg/httpclient"
)

const (
	GoogleUserInfoURL = "https://www.googleapis.com/oauth2/v3/userinfo"
)

type OauthClient struct {
	httpClient *httpclient.Client
}

func NewOauthClient(
	httpClient *httpclient.Client,
) *OauthClient {
	return &OauthClient{
		httpClient: httpClient,
	}
}

func (c *OauthClient) GetGoogleUserInfo(ctx context.Context, token string) (map[string]any, error) {
	optsHeaders := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}

	response, err := c.httpClient.Get(
		ctx,
		GoogleUserInfoURL,
		&httpclient.RequestOptions{
			Headers: optsHeaders,
		},
	)
	if err != nil {
		return nil, err
	}

	var userInfo map[string]any
	if err := response.DecodeJSON(&userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}
