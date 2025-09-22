package storage

import (
	"bytes"
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/thekrauss/beto-shared/pkg/errors"
)

type MinIOBackend struct {
	client *minio.Client
}

// initialise un client MinIO
func NewMinIOBackend(endpoint, accessKey, secretKey string, useSSL bool) (*MinIOBackend, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeInternal, "failed to init MinIO client")
	}
	return &MinIOBackend{client: client}, nil
}

// crée un "bucket"
func (m *MinIOBackend) CreateBucket(ctx context.Context, name string) error {
	err := m.client.MakeBucket(ctx, name, minio.MakeBucketOptions{})
	if err != nil {
		// bucket existe déjà ?
		exists, errBucketExists := m.client.BucketExists(ctx, name)
		if errBucketExists == nil && exists {
			return nil
		}
		return errors.Wrap(err, errors.CodeInternal, "failed to create bucket")
	}
	return nil
}

// liste les objets dans un bucket
func (m *MinIOBackend) ListObjects(ctx context.Context, bucket string) ([]string, error) {
	var objs []string
	for obj := range m.client.ListObjects(ctx, bucket, minio.ListObjectsOptions{Recursive: true}) {
		if obj.Err != nil {
			return nil, errors.Wrap(obj.Err, errors.CodeInternal, "failed to list objects")
		}
		objs = append(objs, obj.Key)
	}
	return objs, nil
}

// envoie un objet
func (m *MinIOBackend) UploadObject(ctx context.Context, bucket, objectName string, content []byte) error {
	_, err := m.client.PutObject(ctx, bucket, objectName, bytes.NewReader(content), int64(len(content)), minio.PutObjectOptions{})
	if err != nil {
		return errors.Wrap(err, errors.CodeInternal, "failed to upload object")
	}
	return nil
}

// télécharge un objet
func (m *MinIOBackend) DownloadObject(ctx context.Context, bucket, objectName string) ([]byte, error) {
	obj, err := m.client.GetObject(ctx, bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeInternal, "failed to get object")
	}
	defer obj.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(obj); err != nil {
		return nil, errors.Wrap(err, errors.CodeInternal, "failed to read object content")
	}
	return buf.Bytes(), nil
}

// supprime un objet
func (m *MinIOBackend) DeleteObject(ctx context.Context, bucket, objectName string) error {
	err := m.client.RemoveObject(ctx, bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return errors.Wrap(err, errors.CodeInternal, "failed to delete object")
	}
	return nil
}
