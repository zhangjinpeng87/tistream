package utils

import (
	"fmt"
)

var (
	ErrInvalidSchemaSnapFile = fmt.Errorf("invalid schema snap file")
	ErrInvalidDataChangeFile = fmt.Errorf("invalid data change file")
)
