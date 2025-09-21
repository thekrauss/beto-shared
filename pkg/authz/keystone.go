package authz

import (
	"context"
	"fmt"
	"net/http"

	"github.com/thekrauss/beto-shared/pkg/errors"
	"github.com/thekrauss/beto-shared/pkg/openstack-client"
)

type Claims struct {
	UserID    string
	UserName  string
	ProjectID string
	Project   string
	Roles     []string
}

type KeystoneValidator struct {
	Client *openstack.KeystoneClient
}

// init avec client OpenStack
func NewKeystoneValidator(client *openstack.KeystoneClient) *KeystoneValidator {
	return &KeystoneValidator{Client: client}
}

// appelle Keystone pour vérifier un token
func (k *KeystoneValidator) ValidateToken(ctx context.Context, token string) (*Claims, error) {
	url := fmt.Sprintf("%s/v3/auth/tokens", k.Client.AuthURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeInternal, "failed to build keystone validation request")
	}
	req.Header.Set("X-Auth-Token", token)
	req.Header.Set("X-Subject-Token", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeKeystoneAuthFailed, "keystone validation failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(errors.CodeKeystoneTokenInvalid, "invalid Keystone token")
	}

	claims := &Claims{
		UserID:    "user-123",
		UserName:  "demo",
		ProjectID: "project-abc",
		Project:   "demo-project",
		Roles:     []string{"admin"},
	}
	return claims, nil
}
