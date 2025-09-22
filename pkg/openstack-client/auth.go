package openstack

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
)

// définit les infos nécessaires pour s’authentifier
type AuthOptions struct {
	IdentityEndpoint string
	Username         string
	Password         string
	DomainName       string
	ProjectName      string
}

// crée un provider Keystone
func NewIdentityClient(opts AuthOptions) (*gophercloud.ProviderClient, error) {
	authOpts := gophercloud.AuthOptions{
		IdentityEndpoint: opts.IdentityEndpoint,
		Username:         opts.Username,
		Password:         opts.Password,
		DomainName:       opts.DomainName,
		Scope: &gophercloud.AuthScope{
			ProjectName: opts.ProjectName,
			DomainName:  opts.DomainName,
		},
	}

	provider, err := openstack.AuthenticatedClient(authOpts)
	if err != nil {
		return nil, err
	}

	return provider, nil
}

func ValidateToken(provider *gophercloud.ProviderClient, token string) (*tokens.Token, error) {
	client, err := openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{})
	if err != nil {
		return nil, err
	}

	result := tokens.Get(client, token)
	return result.Extract()
}
