package filestorage

import (
	"bytes"
	"context"
	"errors"

	"github.com/JuliusMoehring/court-judgment-finder-crawler/logger"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3FileStorage struct {
	logger logger.Logger

	client *s3.Client
	bucket string
}

func NewS3FileStorage(ctx context.Context, logger logger.Logger, bucket string) FileStorage {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		logger.Fatalf("file-storage", "Unable to load AWS SDK config: %s", err)
	}

	client := s3.NewFromConfig(cfg)

	return &S3FileStorage{
		logger: logger,
		client: client,
		bucket: bucket,
	}
}

func (s *S3FileStorage) Exists(ctx context.Context, path string) (bool, error) {
	_, err := s.client.GetObjectAttributes(ctx, &s3.GetObjectAttributesInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
		ObjectAttributes: []types.ObjectAttributes{
			types.ObjectAttributesObjectSize,
		},
	})
	if err != nil {
		var err *types.NoSuchKey
		if errors.As(err, &err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (s *S3FileStorage) Save(ctx context.Context, data []byte, path string) error {
	if _, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
		Body:   bytes.NewReader(data),
	}); err != nil {
		return err
	}

	return nil
}
