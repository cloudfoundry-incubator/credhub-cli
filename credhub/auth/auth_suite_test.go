package auth_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestAuth(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Auth Suite")
}

type DummyClient struct {
	Request  *http.Request
	Response *http.Response
	Error    error
}

func (d *DummyClient) Do(req *http.Request) (*http.Response, error) {
	d.Request = req
	return d.Response, d.Error
}

type dummyUaaClient struct {
	ClientId     string
	ClientSecret string
	Username     string
	Password     string
	AccessToken  string
	RefreshToken string

	NewAccessToken  string
	NewRefreshToken string
	Error           error
}

func (d *dummyUaaClient) ClientCredentialGrant(clientId, clientSecret string) (string, error) {
	d.ClientId = clientId
	d.ClientSecret = clientSecret

	return d.NewAccessToken, d.Error
}

func (d *dummyUaaClient) PasswordGrant(clientId, clientSecret, username, password string) (string, string, error) {
	d.ClientId = clientId
	d.ClientSecret = clientSecret
	d.Username = username
	d.Password = password

	return d.NewAccessToken, d.NewRefreshToken, d.Error
}

func (d *dummyUaaClient) RefreshTokenGrant(clientId, clientSecret, refreshToken string) (string, string, error) {
	d.ClientId = clientId
	d.ClientSecret = clientSecret
	d.RefreshToken = refreshToken

	return d.NewAccessToken, d.NewRefreshToken, d.Error
}

func (d *dummyUaaClient) RevokeToken(token string) error {
	d.AccessToken = token
	return d.Error
}