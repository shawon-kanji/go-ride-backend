package storage

import (
	"context"
	"strings"
	"testing"
	"time"

	"go-ride-backend/internal/config"
)

// newTestS3Storage builds an S3Storage with static credentials and no
// endpoint override, so presigning is a pure local computation (SigV4
// signing) with no network calls — no live S3/AIStor endpoint required.
func newTestS3Storage(t *testing.T, presignExpiry time.Duration) *S3Storage {
	t.Helper()

	s, err := NewS3Storage(context.Background(), config.StorageConfig{
		Bucket:          "test-bucket",
		Region:          "us-east-1",
		AccessKeyID:     "test-access-key",
		SecretAccessKey: "test-secret-key",
		PresignExpiry:   presignExpiry,
	})
	if err != nil {
		t.Fatalf("NewS3Storage returned error: %v", err)
	}
	return s
}

func TestS3StoragePresignUploadURL(t *testing.T) {
	s := newTestS3Storage(t, 5*time.Minute)

	url, err := s.PresignUploadURL(context.Background(), "drivers/123/selfie.jpg", "image/jpeg")
	if err != nil {
		t.Fatalf("PresignUploadURL returned error: %v", err)
	}

	if !strings.Contains(url, "test-bucket") {
		t.Errorf("expected url to reference bucket, got %q", url)
	}
	if !strings.Contains(url, "drivers/123/selfie.jpg") {
		t.Errorf("expected url to reference key, got %q", url)
	}
	if !strings.Contains(url, "X-Amz-Expires=300") {
		t.Errorf("expected url to encode 5m expiry as 300s, got %q", url)
	}
}

func TestS3StoragePresignDownloadURL(t *testing.T) {
	s := newTestS3Storage(t, 0) // exercise the default-expiry fallback

	url, err := s.PresignDownloadURL(context.Background(), "drivers/123/selfie.jpg")
	if err != nil {
		t.Fatalf("PresignDownloadURL returned error: %v", err)
	}

	if !strings.Contains(url, "test-bucket") {
		t.Errorf("expected url to reference bucket, got %q", url)
	}
	if !strings.Contains(url, "drivers/123/selfie.jpg") {
		t.Errorf("expected url to reference key, got %q", url)
	}
	if !strings.Contains(url, "X-Amz-Expires=900") {
		t.Errorf("expected url to encode default 15m expiry as 900s, got %q", url)
	}
}
