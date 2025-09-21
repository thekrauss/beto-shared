package openstack

import (
	"context"
	"fmt"

	"github.com/thekrauss/beto-shared/pkg/errors"
)

// client Keystone
type KeystoneClient struct {
	AuthURL string
}

// initialise un client Keystone
func NewKeystoneClient(authURL string) *KeystoneClient {
	return &KeystoneClient{AuthURL: authURL}
}

type TokenResponse struct {
	Token struct {
		ExpiresAt string `json:"expires_at"`
		User      struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"user"`
		Project struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"project"`
		Roles []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"roles"`
	} `json:"token"`
}

// Login authentifie via Keystone et retourne un token
func (c *KeystoneClient) Login(ctx context.Context, username, password, domain, project string) (*TokenResponse, error) {
	url := fmt.Sprintf("%s/v3/auth/tokens", c.AuthURL)

	reqBody := map[string]any{
		"auth": map[string]any{
			"identity": map[string]any{
				"methods": []string{"password"},
				"password": map[string]any{
					"user": map[string]any{
						"name":     username,
						"password": password,
						"domain":   map[string]string{"name": domain},
					},
				},
			},
			"scope": map[string]any{
				"project": map[string]string{
					"name":   project,
					"domain": domain,
				},
			},
		},
	}

	var resp TokenResponse
	if err := doRequest(ctx, "POST", url, "", reqBody, &resp); err != nil {
		return nil, errors.Wrap(err, errors.CodeKeystoneAuthFailed, "keystone login failed")
	}
	return &resp, nil
}
