package utils

import (
	"hash/crc32"
)

// Use crc32 to calculate the checksum of the data.
func IsChecksumMatch(data []byte, checksum uint32) bool {
	return checksum == crc32.ChecksumIEEE(data)
}
