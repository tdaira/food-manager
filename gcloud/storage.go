package gcloud

import (
	"bufio"
	"cloud.google.com/go/storage"
	"context"
	"os"
	"path/filepath"
)

type Storage struct {
	ctx    context.Context
	client *storage.Client
	bucket *storage.BucketHandle
}

func NewStorage(bucketName string) (*Storage, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	bucket := client.Bucket(bucketName)
	return &Storage{ctx, client, bucket}, nil
}

func (s *Storage) Upload(path string) error {
	fileName := filepath.Base(path)
	if fileName == "." {
		fileName = "default"
	}

	fp, err := os.OpenFile(path, os.O_RDONLY, 0644)
	defer fp.Close()
	if err != nil {
		return err
	}

	reader := bufio.NewReader(fp)
	writer := s.bucket.Object(fileName).NewWriter(s.ctx)
	defer writer.Close()
	_, err = reader.WriteTo(writer)
	if err != nil {
		return err
	}
	return nil
}
