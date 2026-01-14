package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"github.com/yourusername/gin-collection-saas/pkg/logger"
)

// S3Client handles S3 storage operations
type S3Client struct {
	client   *s3.S3
	uploader *s3manager.Uploader
	bucket   string
	region   string
}

// S3Config holds S3 configuration
type S3Config struct {
	Bucket          string
	Region          string
	Endpoint        string // Optional for S3-compatible services
	AccessKeyID     string
	SecretAccessKey string
}

// UploadResult represents the result of an upload
type UploadResult struct {
	Key       string
	URL       string
	SizeBytes int64
}

// NewS3Client creates a new S3 client
func NewS3Client(cfg *S3Config) (*S3Client, error) {
	// Create AWS session
	awsConfig := &aws.Config{
		Region:      aws.String(cfg.Region),
		Credentials: credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
	}

	// Add custom endpoint if provided (for MinIO, etc.)
	if cfg.Endpoint != "" {
		awsConfig.Endpoint = aws.String(cfg.Endpoint)
		awsConfig.S3ForcePathStyle = aws.Bool(true)
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	client := s3.New(sess)
	uploader := s3manager.NewUploader(sess)

	logger.Info("S3 client initialized", "bucket", cfg.Bucket, "region", cfg.Region)

	return &S3Client{
		client:   client,
		uploader: uploader,
		bucket:   cfg.Bucket,
		region:   cfg.Region,
	}, nil
}

// UploadPhoto uploads a photo to S3
func (c *S3Client) UploadPhoto(ctx context.Context, tenantID int64, ginID int64, filename string, data []byte, contentType string) (*UploadResult, error) {
	// Generate unique key
	ext := filepath.Ext(filename)
	key := fmt.Sprintf("tenants/%d/gins/%d/%s%s", tenantID, ginID, uuid.New().String(), ext)

	logger.Info("Uploading photo to S3", "key", key, "size", len(data), "content_type", contentType)

	// Upload to S3
	result, err := c.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket:      aws.String(c.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
		ACL:         aws.String("private"), // Private by default
	})

	if err != nil {
		return nil, fmt.Errorf("failed to upload to S3: %w", err)
	}

	logger.Info("Photo uploaded successfully", "key", key, "location", result.Location)

	return &UploadResult{
		Key:       key,
		URL:       result.Location,
		SizeBytes: int64(len(data)),
	}, nil
}

// DeletePhoto deletes a photo from S3
func (c *S3Client) DeletePhoto(ctx context.Context, key string) error {
	logger.Info("Deleting photo from S3", "key", key)

	_, err := c.client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return fmt.Errorf("failed to delete from S3: %w", err)
	}

	logger.Info("Photo deleted successfully", "key", key)

	return nil
}

// GetPresignedURL generates a presigned URL for temporary access
func (c *S3Client) GetPresignedURL(key string, expiration time.Duration) (string, error) {
	req, _ := c.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})

	url, err := req.Presign(expiration)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return url, nil
}

// DownloadPhoto downloads a photo from S3
func (c *S3Client) DownloadPhoto(ctx context.Context, key string) ([]byte, error) {
	logger.Debug("Downloading photo from S3", "key", key)

	result, err := c.client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to download from S3: %w", err)
	}
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read S3 object: %w", err)
	}

	return data, nil
}

// CheckExists checks if a file exists in S3
func (c *S3Client) CheckExists(ctx context.Context, key string) (bool, error) {
	_, err := c.client.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		// Check if error is "not found"
		return false, nil
	}

	return true, nil
}

// GetObjectSize gets the size of an object in bytes
func (c *S3Client) GetObjectSize(ctx context.Context, key string) (int64, error) {
	result, err := c.client.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return 0, fmt.Errorf("failed to get object metadata: %w", err)
	}

	if result.ContentLength == nil {
		return 0, fmt.Errorf("content length is nil")
	}

	return *result.ContentLength, nil
}
