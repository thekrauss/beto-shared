package authz

import (
	"context"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
	"github.com/thekrauss/beto-shared/pkg/errors"
)

type Claims struct {
	UserID    string
	UserName  string
	ProjectID string
	Project   string
	Roles     []string
}

type KeystoneValidator struct {
	Provider *gophercloud.ProviderClient
}

func NewKeystoneValidator(provider *gophercloud.ProviderClient) *KeystoneValidator {
	return &KeystoneValidator{Provider: provider}
}

func (k *KeystoneValidator) ValidateToken(ctx context.Context, token string) (*Claims, error) {
	client, err := openstack.NewIdentityV3(k.Provider, gophercloud.EndpointOpts{})
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeKeystoneAuthFailed, "failed to init identity client")
	}

	result := tokens.Get(client, token)
	_, err = result.ExtractToken()
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeKeystoneTokenInvalid, "keystone token invalid")
	}

	user, _ := result.ExtractUser()
	project, _ := result.ExtractProject()
	roles, _ := result.ExtractRoles()

	claims := &Claims{
		UserID:    user.ID,
		UserName:  user.Name,
		ProjectID: project.ID,
		Project:   project.Name,
	}
	for _, r := range roles {
		claims.Roles = append(claims.Roles, r.Name)
	}

	return claims, nil
}
