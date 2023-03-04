package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"cloud.google.com/go/storage"
)

type CloudStorage struct {
	client *storage.Client
	bucket *storage.BucketHandle
}

func NewCloudStorage() *CloudStorage {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to create cloud storage client: %v", err)
	}

	bucket := client.Bucket(os.Getenv("CLOUD_STORAGE_BUCKET"))
	if _, err := bucket.Attrs(ctx); err != nil {
		log.Fatalf("failed to get bucket attributes: %v", err)
	}

	return &CloudStorage{
		client: client,
		bucket: bucket,
	}
}

func (s *CloudStorage) Upload(filename string, page string) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	source, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer source.Close()

	obj := s.bucket.Object(page)
	w := obj.NewWriter(ctx)
	if _, err := io.Copy(w, source); err != nil {
		return fmt.Errorf("unable to copy file: %v", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("unable to close file: %v", err)
	}
	if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return fmt.Errorf("unable to set object ACL: %v", err)
	}

	return nil
}
