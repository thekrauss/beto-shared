package storage

import (
	"bytes"
	"context"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/containers"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/objects"
	"github.com/thekrauss/beto-shared/pkg/errors"
)

// implémente ObjectStorageBackend pour OpenStack Swift
type SwiftClient struct {
	client *gophercloud.ServiceClient
}

// initialise un client Swift (Object Storage V1)
func NewSwiftClient(provider *gophercloud.ProviderClient, region string) (*SwiftClient, error) {
	client, err := openstack.NewObjectStorageV1(provider, gophercloud.EndpointOpts{
		Region: region,
	})
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeInternal, "failed to init Swift client")
	}
	return &SwiftClient{client: client}, nil
}

// crée un "container" Swift (équivalent bucket)
func (c *SwiftClient) CreateBucket(ctx context.Context, name string) error {
	createOpts := containers.CreateOpts{}
	res := containers.Create(c.client, name, createOpts)
	if res.Err != nil {
		return errors.Wrap(res.Err, errors.CodeInternal, "failed to create bucket")
	}
	return nil
}

// retourne la liste des objets dans un bucket
func (c *SwiftClient) ListObjects(ctx context.Context, bucket string) ([]string, error) {
	allPages, err := objects.List(c.client, bucket, objects.ListOpts{}).AllPages()
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeInternal, "failed to list objects")
	}

	objList, err := objects.ExtractNames(allPages)
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeInternal, "failed to parse object list")
	}

	return objList, nil
}

// charge un fichier dans un bucket
func (c *SwiftClient) UploadObject(ctx context.Context, bucket, objectName string, content []byte) error {
	createOpts := objects.CreateOpts{
		Content: bytes.NewReader(content),
	}

	res := objects.Create(c.client, bucket, objectName, createOpts)
	if res.Err != nil {
		return errors.Wrap(res.Err, errors.CodeInternal, "failed to upload object")
	}
	return nil
}

// télécharge un objet depuis un bucket
func (c *SwiftClient) DownloadObject(ctx context.Context, bucket, objectName string) ([]byte, error) {
	res := objects.Download(c.client, bucket, objectName, nil)
	data, err := res.ExtractContent()
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeInternal, "failed to download object")
	}
	return data, nil
}

// supprime un objet dans un bucket
func (c *SwiftClient) DeleteObject(ctx context.Context, bucket, objectName string) error {
	res := objects.Delete(c.client, bucket, objectName, nil)
	if res.Err != nil {
		return errors.Wrap(res.Err, errors.CodeInternal, "failed to delete object")
	}
	return nil
}
