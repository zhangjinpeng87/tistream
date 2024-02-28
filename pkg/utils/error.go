package utils

import (
	"fmt"
)

var (
	ErrInvalidSchemaSnapFile         = fmt.Errorf("invalid schema snap file")
	ErrInvalidDataChangeFile         = fmt.Errorf("invalid data change file")
	ErrInvalidDataChangeFileChecksum = fmt.Errorf("invalid data change file checksum")
	ErrInvalidSchemaSnapFileChecksum = fmt.Errorf("invalid schema snap file checksum")
	ErrTenantAlreadyAttched          = fmt.Errorf("tenant already attached")
)
