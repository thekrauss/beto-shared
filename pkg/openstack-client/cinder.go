package openstack

import (
	"context"
	"fmt"

	"github.com/thekrauss/beto-shared/pkg/errors"
)

type CinderClient struct {
	Endpoint string
	Token    string
}

func NewCinderClient(endpoint, token string) *CinderClient {
	return &CinderClient{Endpoint: endpoint, Token: token}
}

type Volume struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Size        int    `json:"size"` // en Go
	Description string `json:"description"`
	Status      string `json:"status"`
}

// crée un volume
func (c *CinderClient) CreateVolume(ctx context.Context, name string, sizeGB int, desc string) (*Volume, error) {
	url := fmt.Sprintf("%s/volumes", c.Endpoint)
	body := map[string]any{
		"volume": map[string]any{
			"name":        name,
			"size":        sizeGB,
			"description": desc,
		},
	}
	var resp struct {
		Volume Volume `json:"volume"`
	}
	if err := doRequest(ctx, "POST", url, c.Token, body, &resp); err != nil {
		return nil, errors.Wrap(err, errors.CodeInternal, "failed to create volume")
	}
	return &resp.Volume, nil
}

// liste les volumes
func (c *CinderClient) ListVolumes(ctx context.Context) ([]Volume, error) {
	url := fmt.Sprintf("%s/volumes/detail", c.Endpoint)
	var resp struct {
		Volumes []Volume `json:"volumes"`
	}
	if err := doRequest(ctx, "GET", url, c.Token, nil, &resp); err != nil {
		return nil, errors.Wrap(err, errors.CodeInternal, "failed to list volumes")
	}
	return resp.Volumes, nil
}

// supprime un volume
func (c *CinderClient) DeleteVolume(ctx context.Context, id string) error {
	url := fmt.Sprintf("%s/volumes/%s", c.Endpoint, id)
	if err := doRequest(ctx, "DELETE", url, c.Token, nil, nil); err != nil {
		return errors.Wrap(err, errors.CodeInternal, "failed to delete volume")
	}
	return nil
}

// attache un volume à une VM
func (c *CinderClient) AttachVolume(ctx context.Context, serverID, volumeID, device string) error {
	url := fmt.Sprintf("%s/servers/%s/os-volume_attachments", c.Endpoint, serverID)
	body := map[string]any{
		"volumeAttachment": map[string]any{
			"volumeId": volumeID,
			"device":   device, // ex: "/dev/vdb"
		},
	}
	if err := doRequest(ctx, "POST", url, c.Token, body, nil); err != nil {
		return errors.Wrap(err, errors.CodeInternal, "failed to attach volume")
	}
	return nil
}

// détache un volume d’une VM
func (c *CinderClient) DetachVolume(ctx context.Context, serverID, attachmentID string) error {
	url := fmt.Sprintf("%s/servers/%s/os-volume_attachments/%s", c.Endpoint, serverID, attachmentID)
	if err := doRequest(ctx, "DELETE", url, c.Token, nil, nil); err != nil {
		return errors.Wrap(err, errors.CodeInternal, "failed to detach volume")
	}
	return nil
}
