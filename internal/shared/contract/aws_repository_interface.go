package contract

import (
	"context"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/parxyws/nego-gin/internal/shared/domain"
)

type AWSRepository interface {
	PutObject(ctx context.Context, entity *domain.UserUploadInput) (*minio.UploadInfo, error)

	GetObject(ctx context.Context, bucketName, objectName string) (*minio.Object, error)

	RemoveObject(ctx context.Context, bucketName, objectName string) error

	PresignedGetObject(ctx context.Context, bucketName, objectName string, expiry time.Duration) (*url.URL, error)
}
