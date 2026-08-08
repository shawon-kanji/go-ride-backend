package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go-ride-backend/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const defaultPresignExpiry = 15 * time.Minute

// S3Storage is a single S3-API adapter that targets either real AWS S3
// (cfg.Endpoint empty, credentials from the pod's IAM role via IRSA) or a
// local/self-hosted S3-compatible store such as AIStor (cfg.Endpoint set,
// cfg.UsePathStyle true, static cfg.AccessKeyID/SecretAccessKey) — both
// speak the identical S3 API, so one client type covers both environments.
type S3Storage struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	bucket        string
	presignExpiry time.Duration
}

func NewS3Storage(ctx context.Context, cfg config.StorageConfig) (*S3Storage, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(cfg.Region))
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}
	if cfg.AccessKeyID != "" {
		awsCfg.Credentials = credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, "")
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
		}
		o.UsePathStyle = cfg.UsePathStyle
	})

	if cfg.Endpoint != "" {
		if err := ensureBucket(ctx, client, cfg.Bucket); err != nil {
			return nil, fmt.Errorf("ensure bucket %q: %w", cfg.Bucket, err)
		}
	}

	presignExpiry := cfg.PresignExpiry
	if presignExpiry == 0 {
		presignExpiry = defaultPresignExpiry
	}

	return &S3Storage{
		client:        client,
		presignClient: s3.NewPresignClient(client),
		bucket:        cfg.Bucket,
		presignExpiry: presignExpiry,
	}, nil
}

// ensureBucket idempotently creates the bucket when it doesn't already
// exist. Only ever called for a local/dev endpoint (see caller) so local
// setup doesn't need a manual `mc mb` step; bucket lifecycle against real
// AWS belongs to Terraform (go-ride-infra), never to app code.
func ensureBucket(ctx context.Context, client *s3.Client, bucket string) error {
	_, err := client.HeadBucket(ctx, &s3.HeadBucketInput{Bucket: aws.String(bucket)})
	if err == nil {
		return nil
	}

	var notFound *s3types.NotFound
	if !errors.As(err, &notFound) {
		return fmt.Errorf("head bucket: %w", err)
	}

	if _, err := client.CreateBucket(ctx, &s3.CreateBucketInput{Bucket: aws.String(bucket)}); err != nil {
		return fmt.Errorf("create bucket: %w", err)
	}
	return nil
}

// PresignUploadURL returns a time-limited URL the caller can PUT the raw
// document bytes to directly — the file never transits the API server.
// contentType is signed into the request, so the PUT must send an identical
// Content-Type header or S3 rejects it as a signature mismatch.
func (s *S3Storage) PresignUploadURL(ctx context.Context, key string, contentType string) (string, error) {
	out, err := s.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(s.presignExpiry))
	if err != nil {
		return "", fmt.Errorf("presign upload url for %q: %w", key, err)
	}
	return out.URL, nil
}

// PresignDownloadURL returns a time-limited URL a backoffice reviewer's
// browser can GET the document from directly.
func (s *S3Storage) PresignDownloadURL(ctx context.Context, key string) (string, error) {
	out, err := s.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(s.presignExpiry))
	if err != nil {
		return "", fmt.Errorf("presign download url for %q: %w", key, err)
	}
	return out.URL, nil
}

// Exists reports whether an object was actually written to key.
func (s *S3Storage) Exists(ctx context.Context, key string) (bool, error) {
	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		var notFound *s3types.NotFound
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, fmt.Errorf("head object %q: %w", key, err)
	}
	return true, nil
}

// Delete removes the object at key (e.g. cleanup of an abandoned,
// never-confirmed upload).
func (s *S3Storage) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("delete object %q: %w", key, err)
	}
	return nil
}
