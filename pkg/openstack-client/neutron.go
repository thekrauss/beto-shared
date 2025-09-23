package openstack

import (
	"context"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/thekrauss/beto-shared/pkg/errors"
)

type NeutronClient struct {
	client *gophercloud.ServiceClient
}

// client Neutron
func NewNeutronClient(provider *gophercloud.ProviderClient, region string) (*NeutronClient, error) {
	client, err := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{
		Region: region,
	})
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeNeutronError, "failed to init Neutron client")
	}
	return &NeutronClient{client: client}, nil
}

// CRUD Réseaux
type Network struct {
	ID   string
	Name string
}

func (c *NeutronClient) CreateNetwork(ctx context.Context, name string) (*Network, error) {
	createOpts := networks.CreateOpts{
		Name:         name,
		AdminStateUp: gophercloud.Enabled,
	}

	n, err := networks.Create(c.client, createOpts).Extract()
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeNeutronError, "failed to create network")
	}

	return &Network{ID: n.ID, Name: n.Name}, nil
}

func (c *NeutronClient) ListNetworks(ctx context.Context) ([]Network, error) {
	allPages, err := networks.List(c.client, networks.ListOpts{}).AllPages()
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeNeutronError, "failed to list networks")
	}

	nList, err := networks.ExtractNetworks(allPages)
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeNeutronError, "failed to parse networks")
	}

	var result []Network
	for _, n := range nList {
		result = append(result, Network{ID: n.ID, Name: n.Name})
	}
	return result, nil
}

// Floating IPs

// reserves a floating IP
func (c *NeutronClient) AllocateFloatingIP(ctx context.Context, networkID string) (string, error) {
	createOpts := floatingips.CreateOpts{
		FloatingNetworkID: networkID,
	}

	fip, err := floatingips.Create(c.client, createOpts).Extract()
	if err != nil {
		return "", errors.Wrap(err, errors.CodeNeutronError, "failed to allocate floating IP")
	}

	return fip.FloatingIP, nil
}

// associates a floating IP with a port (VM)
func (c *NeutronClient) AttachFloatingIP(ctx context.Context, floatingIPID, portID string) error {
	updateOpts := floatingips.UpdateOpts{
		PortID: &portID,
	}

	_, err := floatingips.Update(c.client, floatingIPID, updateOpts).Extract()
	if err != nil {
		return errors.Wrap(err, errors.CodeNeutronError, "failed to attach floating IP")
	}
	return nil
}
