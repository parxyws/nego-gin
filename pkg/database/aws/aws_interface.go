package aws

import (
	"context"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
)

const (
	BucketAvatar  = ""
	BucketProduct = ""
)

type UploadInput struct {
	Object      io.Reader
	ObjectName  string
	ObjectSize  int64
	BucketName  string
	ContentType string
}

type AWSService interface {
	PutObject(ctx context.Context, entity *UploadInput) (*minio.UploadInfo, error)
	GetObject(ctx context.Context, bucketName, objectName string) (*minio.Object, error)
	RemoveObject(ctx context.Context, bucketName, objectName string) error
	PresignedGetObject(ctx context.Context, bucketName, objectName string, expiry time.Duration) (*url.URL, error)
}
