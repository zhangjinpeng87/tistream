package utils

import (
	"fmt"
)

var (
	ErrContentTooShort                       = fmt.Errorf("content too short")
	ErrChecksumNotMatch                      = fmt.Errorf("checksum not match")
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
	ErrInvalidRange                          = fmt.Errorf("invalid range")
	ErrInvalidRangeWatermarksSnapFile        = fmt.Errorf("invalid range watermarks snap file")
	ErrInvalidRangeWatermarksSnapFileVersion = fmt.Errorf("invalid range watermarks snap file version")
	ErrInvalidRangeWatermarksSnapTenantID    = fmt.Errorf("invalid range watermarks snap tenant id")
	ErrInvalidRangeWatermarksSnapRange       = fmt.Errorf("invalid range watermarks snap range")
	ErrInvalidCommittedDataFile              = fmt.Errorf("invalid committed data file")
	ErrInvalidCommittedDataFileVersion       = fmt.Errorf("invalid committed data file version")
	ErrInvalidCommittedManifestFile          = fmt.Errorf("invalid committed manifest file")
)
