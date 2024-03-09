package utils

import (
	"bytes"
	"encoding/binary"
	"sort"
)

func SortDedupSlice(src [][]byte) [][]byte {
	sort.Slice(src, func(i, j int) bool {
		return bytes.Compare(src[i], src[j]) < 0
	})
	dst := make([][]byte, 0, len(src))
	for i := 0; i < len(src); i++ {
		if i == 0 || bytes.Compare(src[i], src[i-1]) > 0 {
			dst = append(dst, src[i])
		}
	}
	return dst
}

func TsToReverseBytes(ts uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, ^ts)
	return b
}

func ReverseBytesToTs(b []byte) uint64 {
	return ^binary.BigEndian.Uint64(b)
}
