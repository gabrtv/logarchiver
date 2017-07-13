package main

import (
	"log"

	"io"

	"github.com/minio/minio-go"
)

func newMinioClient(cfg *minioConfig) (*minio.Client, error) {

	minioClient, err := minio.New(cfg.Endpoint, cfg.AccessKey, cfg.SecretKey, cfg.SSL)
	if err != nil {
		return nil, err
	}
	return minioClient, nil
}

func createBucket(cfg *minioConfig, c *minio.Client) error {

	err := c.MakeBucket(cfg.BucketName, cfg.Location)
	if err != nil {
		exists, err := c.BucketExists(cfg.BucketName)
		if err == nil && exists {
			log.Printf("We already own %s\n", cfg.BucketName)
		} else {
			return err
		}
	}
	log.Printf("Successfully created %s\n", cfg.BucketName)
	return nil
}

func uploadStream(cfg *minioConfig, c *minio.Client, r io.Reader) error {

	objectName := "logs.out"

	n, err := c.PutObjectStreaming(cfg.BucketName, objectName, r)
	if err != nil {
		return err
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, n)
	return nil
}
