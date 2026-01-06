package aws

import (
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/parxyws/nego-gin/config"
)

func NewAWSClient(config *config.Config) (*minio.Client, error) {
	client, err := minio.New(config.AWS.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AWS.MinioAccessKey, config.AWS.MinioSecretKey, ""),
		Secure: config.AWS.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("create minio client: %w", err)
	}

	return client, err
}
