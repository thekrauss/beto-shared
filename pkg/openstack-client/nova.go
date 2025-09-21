package openstack

import (
	"context"
	"fmt"

	"github.com/thekrauss/beto-shared/pkg/errors"
)

type NovaClient struct {
	Endpoint string
	Token    string
}

func NewNovaClient(endpoint, token string) *NovaClient {
	return &NovaClient{Endpoint: endpoint, Token: token}
}

type Server struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateServerRequest struct {
	Server struct {
		Name      string                 `json:"name"`
		ImageRef  string                 `json:"imageRef"`
		FlavorRef string                 `json:"flavorRef"`
		Networks  []map[string]string    `json:"networks"`
		Metadata  map[string]interface{} `json:"metadata,omitempty"`
	} `json:"server"`
}

// crée une VM
func (c *NovaClient) CreateVM(ctx context.Context, req CreateServerRequest) (*Server, error) {
	url := fmt.Sprintf("%s/servers", c.Endpoint)
	var resp struct {
		Server Server `json:"server"`
	}
	if err := doRequest(ctx, "POST", url, c.Token, req, &resp); err != nil {
		return nil, errors.Wrap(err, errors.CodeNovaError, "failed to create VM")
	}
	return &resp.Server, nil
}

// liste les VMs
func (c *NovaClient) ListVMs(ctx context.Context) ([]Server, error) {
	url := fmt.Sprintf("%s/servers/detail", c.Endpoint)
	var resp struct {
		Servers []Server `json:"servers"`
	}
	if err := doRequest(ctx, "GET", url, c.Token, nil, &resp); err != nil {
		return nil, errors.Wrap(err, errors.CodeNovaError, "failed to list VMs")
	}
	return resp.Servers, nil
}

// supprime une VM
func (c *NovaClient) DeleteVM(ctx context.Context, id string) error {
	url := fmt.Sprintf("%s/servers/%s", c.Endpoint, id)
	if err := doRequest(ctx, "DELETE", url, c.Token, nil, nil); err != nil {
		return errors.Wrap(err, errors.CodeNovaError, "failed to delete VM")
	}
	return nil
}

// démarre une VM
func (c *NovaClient) StartVM(ctx context.Context, id string) error {
	url := fmt.Sprintf("%s/servers/%s/action", c.Endpoint, id)
	body := map[string]any{"os-start": nil}
	if err := doRequest(ctx, "POST", url, c.Token, body, nil); err != nil {
		return errors.Wrap(err, errors.CodeNovaError, "failed to start VM")
	}
	return nil
}

// arrête une VM
func (c *NovaClient) StopVM(ctx context.Context, id string) error {
	url := fmt.Sprintf("%s/servers/%s/action", c.Endpoint, id)
	body := map[string]any{"os-stop": nil}
	if err := doRequest(ctx, "POST", url, c.Token, body, nil); err != nil {
		return errors.Wrap(err, errors.CodeNovaError, "failed to stop VM")
	}
	return nil
}
