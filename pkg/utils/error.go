package utils

import (
	"fmt"
)

var (
	ErrInvalidSchemaSnapFile                 = fmt.Errorf("invalid schema snap file")
	ErrInvalidDataChangeFile                 = fmt.Errorf("invalid data change file")
	ErrInvalidDataChangeFileChecksum         = fmt.Errorf("invalid data change file checksum")
	ErrInvalidSchemaSnapFileChecksum         = fmt.Errorf("invalid schema snap file checksum")
	ErrTenantAlreadyAttched                  = fmt.Errorf("tenant already attached")
	ErrInvalidSorterBufferSnapFile           = fmt.Errorf("invalid sorter buffer snap file")
	ErrInvalidSorterBufferSnapFileChecksum   = fmt.Errorf("invalid sorter buffer snap file checksum")
	ErrUnmatchedTenantID                     = fmt.Errorf("unmatched tenant id")
	ErrUnmatchedRange                        = fmt.Errorf("unmatched range")
	ErrInvalidPrewriteBufferSnapFile         = fmt.Errorf("invalid prewrite buffer snap file")
	ErrInvalidPrewriteBufferSnapFileChecksum = fmt.Errorf("invalid prewrite buffer snap file checksum")
	ErrInvalidSorterBufferSnapFileVersion    = fmt.Errorf("invalid sorter buffer snap file version")
	ErrInvalidPrewriteBufferSnapFileVersion  = fmt.Errorf("invalid prewrite buffer snap file version")
)
