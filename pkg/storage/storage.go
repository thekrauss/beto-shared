package storage

import (
	"context"

	"github.com/gophercloud/gophercloud"
	"github.com/thekrauss/beto-shared/pkg/errors"
)

// définit les opérations communes à tout backend de stockage objet
type ObjectStorageBackend interface {
	CreateBucket(ctx context.Context, name string) error
	ListObjects(ctx context.Context, bucket string) ([]string, error)
	UploadObject(ctx context.Context, bucket, objectName string, content []byte) error
	DownloadObject(ctx context.Context, bucket, objectName string) ([]byte, error)
	DeleteObject(ctx context.Context, bucket, objectName string) error
}

// Config générique pour sélectionner le backend
type Config struct {
	Backend string // "swift" ou "minio"

	// Swift
	Region   string
	Provider *gophercloud.ProviderClient

	// MinIO
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
}

// retourne une implémentation en fonction de la config
func NewObjectStorage(cfg Config) (ObjectStorageBackend, error) {
	switch cfg.Backend {
	case "swift":
		if cfg.Provider == nil {
			return nil, errors.New(errors.CodeInvalidInput, "provider client required for Swift backend")
		}
		client, err := NewSwiftClient(cfg.Provider, cfg.Region)
		if err != nil {
			return nil, err
		}
		return client, nil

	case "minio":
		client, err := NewMinIOBackend(cfg.Endpoint, cfg.AccessKey, cfg.SecretKey, cfg.UseSSL)
		if err != nil {
			return nil, err
		}
		return client, nil

	default:
		return nil, errors.Newf(errors.CodeInvalidInput, "unsupported storage backend: %s", cfg.Backend)
	}
}
