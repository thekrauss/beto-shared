package openstack

import (
	"context"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/volumeattach"
	"github.com/thekrauss/beto-shared/pkg/errors"
)

type CinderClient struct {
	client *gophercloud.ServiceClient
}

// client Cinder (Block Storage)
func NewCinderClient(provider *gophercloud.ProviderClient, region string) (*CinderClient, error) {
	client, err := openstack.NewBlockStorageV3(provider, gophercloud.EndpointOpts{
		Region: region,
	})
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeInternal, "failed to init Cinder client")
	}
	return &CinderClient{client: client}, nil
}

type Volume struct {
	ID          string
	Name        string
	Size        int
	Description string
	Status      string
}

func (c *CinderClient) CreateVolume(ctx context.Context, name string, sizeGB int, desc string) (*Volume, error) {
	createOpts := volumes.CreateOpts{
		Name:        name,
		Size:        sizeGB,
		Description: desc,
	}

	v, err := volumes.Create(c.client, createOpts).Extract()
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeInternal, "failed to create volume")
	}

	return &Volume{
		ID:          v.ID,
		Name:        v.Name,
		Size:        v.Size,
		Description: v.Description,
		Status:      v.Status,
	}, nil
}

func (c *CinderClient) ListVolumes(ctx context.Context) ([]Volume, error) {
	allPages, err := volumes.List(c.client, volumes.ListOpts{}).AllPages()
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeInternal, "failed to list volumes")
	}

	vList, err := volumes.ExtractVolumes(allPages)
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeInternal, "failed to parse volumes")
	}

	var result []Volume
	for _, v := range vList {
		result = append(result, Volume{
			ID:          v.ID,
			Name:        v.Name,
			Size:        v.Size,
			Description: v.Description,
			Status:      v.Status,
		})
	}
	return result, nil
}

func (c *CinderClient) DeleteVolume(ctx context.Context, id string) error {
	res := volumes.Delete(c.client, id, volumes.DeleteOpts{})
	if res.Err != nil {
		return errors.Wrap(res.Err, errors.CodeInternal, "failed to delete volume")
	}
	return nil
}

func AttachVolume(ctx context.Context, computeClient *gophercloud.ServiceClient, serverID, volumeID, device string) error {
	attachOpts := volumeattach.CreateOpts{
		VolumeID: volumeID,
		Device:   device, // ex: "/dev/vdb"
	}

	_, err := volumeattach.Create(computeClient, serverID, attachOpts).Extract()
	if err != nil {
		return errors.Wrap(err, errors.CodeInternal, "failed to attach volume")
	}
	return nil
}

func DetachVolume(ctx context.Context, computeClient *gophercloud.ServiceClient, serverID, attachmentID string) error {
	err := volumeattach.Delete(computeClient, serverID, attachmentID).ExtractErr()
	if err != nil {
		return errors.Wrap(err, errors.CodeInternal, "failed to detach volume")
	}
	return nil
}
