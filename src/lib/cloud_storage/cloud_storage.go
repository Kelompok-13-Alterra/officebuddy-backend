package cloud_storage

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
	UploadFile(ctx context.Context, file multipart.File, object string, uploadPath string) error
	GetSignedURL(ctx context.Context, objectName string, uploadPath string) (string, error)
}

type Config struct {
	PrivateKey  string `mapstructure:"private_key"`
	ClientEmail string `mapstructure:"client_email"`
}

type cloudStorage struct {
	cl         *storage.Client
	projectID  string
	bucketName string
	conf       Config
}

func Init(path string, cfg Config) Interface {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", path)
	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("failed to create client %v", err)
	}

	c := &cloudStorage{
		cl:         client,
		bucketName: "officebuddy-images",
		projectID:  "officebuddy-388208",
		conf:       cfg,
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

func (c *cloudStorage) GetSignedURL(ctx context.Context, objectName string, uploadPath string) (string, error) {
	opts := &storage.SignedURLOptions{
		GoogleAccessID: c.conf.ClientEmail,
		PrivateKey:     []byte(c.conf.PrivateKey),
		Method:         "GET",
		Expires:        time.Now().Add(24 * time.Hour),
	}
	url, err := storage.SignedURL(c.bucketName, fmt.Sprintf("%s%s", uploadPath, objectName), opts)
	if err != nil {
		return "", err
	}

	return url, nil
}
