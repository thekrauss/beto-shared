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

type SwiftClient struct {
	client *gophercloud.ServiceClient
}

func NewSwiftClient(provider *gophercloud.ProviderClient, region string) (*SwiftClient, error) {
	client, err := openstack.NewObjectStorageV1(provider, gophercloud.EndpointOpts{
		Region: region,
	})
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeInternal, "failed to init Swift client")
	}
	return &SwiftClient{client: client}, nil
}

func (c *SwiftClient) CreateBucket(ctx context.Context, name string) error {
	createOpts := containers.CreateOpts{}
	res := containers.Create(c.client, name, createOpts)
	if res.Err != nil {
		return errors.Wrap(res.Err, errors.CodeInternal, "failed to create bucket")
	}
	return nil
}

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

func (c *SwiftClient) DownloadObject(ctx context.Context, bucket, objectName string) ([]byte, error) {
	res := objects.Download(c.client, bucket, objectName, nil)
	data, err := res.ExtractContent()
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeInternal, "failed to download object")
	}
	return data, nil
}

func (c *SwiftClient) DeleteObject(ctx context.Context, bucket, objectName string) error {
	res := objects.Delete(c.client, bucket, objectName, nil)
	if res.Err != nil {
		return errors.Wrap(res.Err, errors.CodeInternal, "failed to delete object")
	}
	return nil
}
