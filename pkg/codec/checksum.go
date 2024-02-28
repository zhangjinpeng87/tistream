package codec

import (
	"hash/crc32"
)

func CalcChecksum(data []byte) uint32 {
	return crc32.ChecksumIEEE(data)
}
