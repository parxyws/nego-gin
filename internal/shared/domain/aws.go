package domain

import "io"

type UserUploadInput struct {
	Object      io.Reader
	ObjectName  string
	ObjectSize  int64
	BucketName  string
	ContentType string
}
