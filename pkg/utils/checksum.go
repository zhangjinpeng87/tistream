package utils

import (
	"hash/crc32"
)

// Use crc32 to calculate the checksum of the data.
func IsChecksumMatch(checksum uint32, data []byte) bool {
	return checksum == crc32.ChecksumIEEE(data)
}
