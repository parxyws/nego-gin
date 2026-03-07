package util

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/parxyws/nego-gin/pkg/database/aws"
)

func ReadUserImageRequest(c *gin.Context, fieldName string) (*aws.UploadInput, error) {
	image, err := c.FormFile(fieldName)
	if err != nil {
		return nil, err
	}

	if image.Size > 1<<20 {
		return nil, fmt.Errorf("file size exceeds 1 MB")
	}

	allowedContentTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/jpg":  true,
	}

	file, err := image.Open()
	if err != nil {
		return nil, fmt.Errorf("open image: %w", err)
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)

	buf := new(bytes.Buffer)
	if _, err = buf.ReadFrom(file); err != nil {
		return nil, fmt.Errorf("read image: %w", err)
	}

	contentType := http.DetectContentType(buf.Bytes())
	if !allowedContentTypes[contentType] {
		return nil, fmt.Errorf("content type not allowed: %s", contentType)
	}

	fileExtension := ".jpg"
	if contentType == "image/png" {
		fileExtension = ".png"
	}

	if contentType == "image/jpeg" {
		fileExtension = ".jpeg"
	}

	newFileName := fmt.Sprintf("%s%s", GenerateShortUUID(), fileExtension)

	model := &aws.UploadInput{
		Object:      bytes.NewReader(buf.Bytes()),
		ObjectName:  newFileName,
		ObjectSize:  image.Size,
		BucketName:  aws.BucketAvatar,
		ContentType: contentType,
	}

	return model, nil

}

func ReadProductImageRequest(c *gin.Context, fieldName string) (*aws.UploadInput, error) {
	image, err := c.FormFile(fieldName)
	if err != nil {
		return nil, err
	}

	if image.Size > 3<<20 {
		return nil, fmt.Errorf("file size exceeds 3 MB")
	}

	allowedContentTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/jpg":  true,
	}

	file, err := image.Open()
	if err != nil {
		return nil, fmt.Errorf("open image: %w", err)
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)

	buf := new(bytes.Buffer)
	if _, err = buf.ReadFrom(file); err != nil {
		return nil, fmt.Errorf("read image: %w", err)
	}

	contentType := http.DetectContentType(buf.Bytes())
	if !allowedContentTypes[contentType] {
		return nil, fmt.Errorf("content type not allowed: %s", contentType)
	}

	fileExtension := ".jpg"
	if contentType == "image/png" {
		fileExtension = ".png"
	}

	if contentType == "image/jpeg" {
		fileExtension = ".jpeg"
	}

	newFileName := fmt.Sprintf("%s%s", GenerateShortUUID(), fileExtension)

	model := &aws.UploadInput{
		Object:      bytes.NewReader(buf.Bytes()),
		ObjectName:  newFileName,
		ObjectSize:  image.Size,
		BucketName:  aws.BucketProduct,
		ContentType: contentType,
	}

	return model, nil
}

func GenerateShortUUID() string {
	u := uuid.New()
	return u.String()[:8]
}

func GenerateAWSMinioURL(bucket, key, endpoint string) string {
	return fmt.Sprintf("%s/%s/%s", endpoint, bucket, key)
}
