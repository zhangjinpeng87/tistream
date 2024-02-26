package storage

import (
	"errors"

	"github.com/zhangjinpeng87/tistream/pkg/utils"
)

type ExternalStorage interface {
	// GetFile returns the file content of the given path.
	GetFile(filePath string) ([]byte, error)

	// PutFile puts the file content to the given path.
	PutFile(filePath string, content []byte) error

	// ListFiles returns the file names under the given path.
	ListFiles(dirPath string) ([]string, error)

	// ListSubDir returns the sub-directory of the given path.
	ListSubDir(rootDir string) ([]string, error)

	// DeleteFile deletes the file of the given path.
	DeleteFile(filePath string) error
}

var (
	// ErrUnknownStorage is the error of unknown storage.
	ErrUnknownStorage = errors.New("unknown storage")
)

const (
	// S3StorageURL is the url of the S3 storage.
	S3 = "s3"
)

// NewExternalStorage creates a new external storage.
func NewExternalStorage(storageType string, cfg *utils.StorageConfig) (ExternalStorage, error) {
	switch storageType {
	case S3:
		return NewS3Storage(cfg), nil
	default:
		return nil, ErrUnknownStorage
	}
}
