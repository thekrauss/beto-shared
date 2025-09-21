package openstack

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/thekrauss/beto-shared/pkg/errors"
)

type SwiftClient struct {
	Endpoint string
	Token    string
}

// init Swift
func NewSwiftClient(endpoint, token string) *SwiftClient {
	return &SwiftClient{Endpoint: endpoint, Token: token}
}

// CreateContainer crée un "bucket"
func (c *SwiftClient) CreateContainer(ctx context.Context, name string) error {
	url := fmt.Sprintf("%s/%s", c.Endpoint, name)
	req, err := http.NewRequestWithContext(ctx, "PUT", url, nil)
	if err != nil {
		return errors.Wrap(err, errors.CodeInternal, "failed to build request")
	}
	req.Header.Set("X-Auth-Token", c.Token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, errors.CodeInternal, "swift create container failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return errors.Newf(errors.CodeInternal, "swift error %d while creating container", resp.StatusCode)
	}
	return nil
}

// upload un fichier dans un container
func (c *SwiftClient) UploadObject(ctx context.Context, container, objectName string, content io.Reader) error {
	url := fmt.Sprintf("%s/%s/%s", c.Endpoint, container, objectName)

	req, err := http.NewRequestWithContext(ctx, "PUT", url, content)
	if err != nil {
		return errors.Wrap(err, errors.CodeInternal, "failed to build request")
	}
	req.Header.Set("X-Auth-Token", c.Token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, errors.CodeInternal, "swift upload failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return errors.Newf(errors.CodeInternal, "swift error %d while uploading object", resp.StatusCode)
	}
	return nil
}

// récupère un fichier
func (c *SwiftClient) DownloadObject(ctx context.Context, container, objectName string) (io.ReadCloser, error) {
	url := fmt.Sprintf("%s/%s/%s", c.Endpoint, container, objectName)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeInternal, "failed to build request")
	}
	req.Header.Set("X-Auth-Token", c.Token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeInternal, "swift download failed")
	}

	if resp.StatusCode >= 400 {
		return nil, errors.Newf(errors.CodeInternal, "swift error %d while downloading object", resp.StatusCode)
	}
	return resp.Body, nil
}

// supprime un fichier
func (c *SwiftClient) DeleteObject(ctx context.Context, container, objectName string) error {
	url := fmt.Sprintf("%s/%s/%s", c.Endpoint, container, objectName)

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return errors.Wrap(err, errors.CodeInternal, "failed to build request")
	}
	req.Header.Set("X-Auth-Token", c.Token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, errors.CodeInternal, "swift delete failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return errors.Newf(errors.CodeInternal, "swift error %d while deleting object", resp.StatusCode)
	}
	return nil
}

// liste les objets d’un container
func (c *SwiftClient) ListObjects(ctx context.Context, container string) ([]string, error) {
	url := fmt.Sprintf("%s/%s", c.Endpoint, container)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeInternal, "failed to build request")
	}
	req.Header.Set("X-Auth-Token", c.Token)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeInternal, "swift list objects failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, errors.Newf(errors.CodeInternal, "swift error %d while listing objects", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	objects := strings.Split(strings.TrimSpace(string(body)), "\n")
	return objects, nil
}
