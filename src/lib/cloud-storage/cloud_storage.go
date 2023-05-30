package clod_storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"time"

	"cloud.google.com/go/storage"
)

type Interface interface {
}

type cloudStorage struct {
	cl         *storage.Client
	projectID  string
	bucketName string
}

func Init(path string) Interface {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", path)
	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("failed to create client %v", err)
	}

	c := &cloudStorage{
		cl:         client,
		bucketName: "officebuddy-images",
		projectID:  "officebuddy-388208",
	}

	return c
}

func (c *cloudStorage) UploadFile(ctx context.Context, file multipart.File, object string, uploadPath string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	wc := c.cl.Bucket(c.bucketName).Object(uploadPath + object).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return fmt.Errorf("io.Copy, %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writter.Close : %v", err)
	}

	return nil
}
