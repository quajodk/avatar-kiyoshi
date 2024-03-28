package gcstorage

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
)

const (
	productID  = "eng-pulsar-400816"
	bucketName = "wp-core-media"
)

type GCSUploader struct {
	client     *storage.Client
	projectID  string
	bucketName string
}

var Uploader *GCSUploader

func IntGCS() {
	ctx := context.Background()
	credPath := filepath.Join("gcs-credentials.json")
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", credPath)
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	Uploader = &GCSUploader{
		client: client, projectID: productID, bucketName: bucketName,
	}
}

func (c *GCSUploader) UploadFile(fileName string, contentType string, r *io.PipeReader) (string, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	fileName = path.Base(fileName)

	w := c.client.Bucket(c.bucketName).Object(fileName).NewWriter(ctx)
	w.ContentType = contentType

	if _, err := io.Copy(w, r); err != nil {
		return "", fmt.Errorf("io.Copy: %v", err)
	}
	if err := w.Close(); err != nil {
		return "", fmt.Errorf("Writer.Close: %v", err)
	}

	// get signed url for the file
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV2,
		Method:  "GET",
		Expires: time.Now().Add(43830 * time.Hour),
	}
	objUrl, err := c.client.Bucket(c.bucketName).SignedURL(fileName, opts)
	if err != nil {
		return "", fmt.Errorf("Bucket(%q).SignedURL: %w", c.bucketName, err)
	}
	return objUrl, nil

}
