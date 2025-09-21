package openstack

import (
	"context"
	"fmt"

	"github.com/thekrauss/beto-shared/pkg/errors"
)

type NeutronClient struct {
	Endpoint string
	Token    string
}

func NewNeutronClient(endpoint, token string) *NeutronClient {
	return &NeutronClient{Endpoint: endpoint, Token: token}
}

type Network struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateNetworkRequest struct {
	Network struct {
		Name string `json:"name"`
	} `json:"network"`
}

// crée un réseau
func (c *NeutronClient) CreateNetwork(ctx context.Context, req CreateNetworkRequest) (*Network, error) {
	url := fmt.Sprintf("%s/v2.0/networks", c.Endpoint)
	var resp struct {
		Network Network `json:"network"`
	}
	if err := doRequest(ctx, "POST", url, c.Token, req, &resp); err != nil {
		return nil, errors.Wrap(err, errors.CodeNeutronError, "failed to create network")
	}
	return &resp.Network, nil
}

// liste les réseaux
func (c *NeutronClient) ListNetworks(ctx context.Context) ([]Network, error) {
	url := fmt.Sprintf("%s/v2.0/networks", c.Endpoint)
	var resp struct {
		Networks []Network `json:"networks"`
	}
	if err := doRequest(ctx, "GET", url, c.Token, nil, &resp); err != nil {
		return nil, errors.Wrap(err, errors.CodeNeutronError, "failed to list networks")
	}
	return resp.Networks, nil
}

// réserve une IP flottante
func (c *NeutronClient) AllocateFloatingIP(ctx context.Context, networkID string) (string, error) {
	url := fmt.Sprintf("%s/v2.0/floatingips", c.Endpoint)
	body := map[string]any{
		"floatingip": map[string]string{"floating_network_id": networkID},
	}
	var resp struct {
		FloatingIP struct {
			ID        string `json:"id"`
			IPAddress string `json:"floating_ip_address"`
		} `json:"floatingip"`
	}
	if err := doRequest(ctx, "POST", url, c.Token, body, &resp); err != nil {
		return "", errors.Wrap(err, errors.CodeNeutronError, "failed to allocate floating IP")
	}
	return resp.FloatingIP.IPAddress, nil
}

// associe une IP flottante à une VM (port)
func (c *NeutronClient) AttachFloatingIP(ctx context.Context, floatingIPID, portID string) error {
	url := fmt.Sprintf("%s/v2.0/floatingips/%s", c.Endpoint, floatingIPID)
	body := map[string]any{
		"floatingip": map[string]string{"port_id": portID},
	}
	if err := doRequest(ctx, "PUT", url, c.Token, body, nil); err != nil {
		return errors.Wrap(err, errors.CodeNeutronError, "failed to attach floating IP")
	}
	return nil
}
