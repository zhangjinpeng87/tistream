package storage

type BackendStorage interface {
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
