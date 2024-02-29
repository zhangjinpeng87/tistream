package utils

import (
	"bytes"
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
