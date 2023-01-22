package gcs

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"cloud.google.com/go/storage"
)

const (
	TIMEOUT_VALUE = 50 // sec
)

type GCSKey struct {
	PrivateKey  string `json:"private_key"`
	ClientEmail string `json:"client_email"`
}

// holds the necessary data to interact with GCS
type ClientUploader struct {
	Cl         *storage.Client
	ProjectID  string
	BucketName string
	UploadPath string
}

// creates a new instance of ClientUploader
func NewClientUploader(projectID, bucketName, uploadPath string) (*ClientUploader, error) {
	client, err := storage.NewClient(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Failed to create client: %v", err)
	}

	return &ClientUploader{
		Cl:         client,
		ProjectID:  projectID,
		BucketName: bucketName,
		UploadPath: uploadPath,
	}, nil
}

func (c *ClientUploader) UploadFile(file multipart.File, object string, ctx context.Context) error {
	bucket := c.Cl.Bucket(c.BucketName)
	objectPath := c.UploadPath + object
	wc := bucket.Object(objectPath).NewWriter(ctx)

	if _, err := io.Copy(wc, file); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}

	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	return nil
}
