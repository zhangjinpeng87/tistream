package utils

import (
	"bytes"

	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
)

func RangeOverlapped(r1 *pb.Task_Range, r2 *pb.Task_Range) bool {
	if bytes.Compare(r1.End, r2.Start) <= 0 || bytes.Compare(r2.End, r1.Start) <= 0 {
		return false
	}
	return true
}
