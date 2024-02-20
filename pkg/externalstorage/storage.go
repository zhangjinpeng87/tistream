package externalstorage

type ExternalStorage interface {
	// GetFile returns the file content of the given path.
	GetFile(path string) ([]byte, error)

	// PutFile puts the file content to the given path.
	PutFile(path string, content []byte) error

	// ListFiles returns the file names under the given path.
	ListFiles(path string) ([]string, error)

	// DeleteFile deletes the file of the given path.
	DeleteFile(path string) error
}
