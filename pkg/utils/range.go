package utils

import (
	"bytes"

	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
)

func RangeOverlapped(r1 *pb.Task_Range, r2 *pb.Task_Range) bool {
	if bytes.Compare(r1.Start, r2.End) < 0 && bytes.Compare(r2.Start, r1.End) < 0 {
		return true
	}
	return false
}
