package metadata

import (
	"bytes"

	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
)

type RangeStart []byte

type TenantTask struct {
	// RangeStart is the start of the range.
	RangeStart

	// Task contains rangeStart, rangeEnd, tenantId, and other information.
	Task *pb.Task

	// assignedSorter indicates which sorter is responsible for this task.
	AssignedSorter uint32
}

func taskLess(a, b TenantTask) bool {
	return bytes.Compare(a.RangeStart, b.RangeStart) < 0
}
