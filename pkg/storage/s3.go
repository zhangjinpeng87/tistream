package storage

import (
	"bytes"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/zhangjinpeng87/tistream/pkg/utils"
)

// S3Storage is aws S3 service.
type S3Storage struct {
	// The S3 client.
	client *s3.S3

	// Configuration.
	region    string
	bucket    string
	accessKey string
	secureKey string
}

// NewS3Storage creates a new S3 client.
func NewS3Storage(cfg *utils.StorageConfig) ExternalStorage {
	session := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(cfg.Region),
		Credentials: credentials.NewStaticCredentials(cfg.AccessKey, cfg.SecurityKey, ""),
	}))

	return &S3Storage{
		client:    s3.New(session),
		region:    cfg.Region,
		bucket:    cfg.Bucket,
		accessKey: cfg.AccessKey,
		secureKey: cfg.SecurityKey,
	}
}

// GetFile gets the file from the S3.
func (s *S3Storage) GetFile(key string) ([]byte, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	output, err := s.client.GetObject(input)
	if err != nil {
		return nil, err
	}
	defer output.Body.Close()

	return io.ReadAll(output.Body)
}

// PutFile puts the file to the S3.
func (s *S3Storage) PutFile(key string, content []byte) error {
	input := &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(content),
	}

	_, err := s.client.PutObject(input)
	if err != nil {
		return err
	}

	return nil
}

// ListFiles lists the files in the S3.
func (s *S3Storage) ListFiles(dirPath string) ([]string, error) {
	input := &s3.ListObjectsInput{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(dirPath),
	}

	output, err := s.client.ListObjects(input)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, obj := range output.Contents {
		files = append(files, *obj.Key)
	}

	return files, nil
}

// ListSubDir lists the sub-directory in the S3.
func (s *S3Storage) ListSubDir(rootDir string) ([]string, error) {
	input := &s3.ListObjectsInput{
		Bucket:    aws.String(s.bucket),
		Prefix:    aws.String(rootDir),
		Delimiter: aws.String("/"),
	}

	output, err := s.client.ListObjects(input)
	if err != nil {
		return nil, err
	}

	var dirs []string
	for _, commonPrefix := range output.CommonPrefixes {
		dirs = append(dirs, *commonPrefix.Prefix)
	}

	return dirs, nil
}

// DeleteFile deletes the file in the S3.
func (s *S3Storage) DeleteFile(filePath string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filePath),
	}

	_, err := s.client.DeleteObject(input)
	if err != nil {
		return err
	}

	return nil
}
