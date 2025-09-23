package openstack

import (
	"context"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/startstop"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/thekrauss/beto-shared/pkg/errors"
)

type NovaClient struct {
	client *gophercloud.ServiceClient
}

func NewNovaClient(provider *gophercloud.ProviderClient, region string) (*NovaClient, error) {
	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: region,
	})
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeNovaError, "failed to init Nova client")
	}
	return &NovaClient{client: client}, nil
}

type Server struct {
	ID   string
	Name string
}

func (c *NovaClient) CreateVM(ctx context.Context, name, imageRef, flavorRef string, networkID string) (*Server, error) {
	createOpts := servers.CreateOpts{
		Name:      name,
		ImageRef:  imageRef,
		FlavorRef: flavorRef,
		Networks: []servers.Network{
			{UUID: networkID},
		},
	}

	s, err := servers.Create(c.client, createOpts).Extract()
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeNovaError, "failed to create VM")
	}

	return &Server{
		ID:   s.ID,
		Name: s.Name,
	}, nil
}

func (c *NovaClient) ListVMs(ctx context.Context) ([]Server, error) {
	pager := servers.List(c.client, servers.ListOpts{})
	var result []Server

	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		sList, err := servers.ExtractServers(page)
		if err != nil {
			return false, err
		}
		for _, s := range sList {
			result = append(result, Server{ID: s.ID, Name: s.Name})
		}
		return true, nil
	})
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeNovaError, "failed to list VMs")
	}

	return result, nil
}

func (c *NovaClient) DeleteVM(ctx context.Context, id string) error {
	err := servers.Delete(c.client, id).ExtractErr()
	if err != nil {
		return errors.Wrap(err, errors.CodeNovaError, "failed to delete VM")
	}
	return nil
}

func (c *NovaClient) StartVM(ctx context.Context, id string) error {
	err := startstop.Start(c.client, id).ExtractErr()
	if err != nil {
		return errors.Wrap(err, errors.CodeNovaError, "failed to start VM")
	}
	return nil
}

func (c *NovaClient) StopVM(ctx context.Context, id string) error {
	err := startstop.Stop(c.client, id).ExtractErr()
	if err != nil {
		return errors.Wrap(err, errors.CodeNovaError, "failed to stop VM")
	}
	return nil
}
