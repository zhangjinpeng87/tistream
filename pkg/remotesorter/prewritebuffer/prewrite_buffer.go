package prewritebuffer

import (
	"fmt"

	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
)

type PrewriteBuffer struct {
	// The tenant id.
	tenantID uint64
	Range    *pb.Task_Range

	prewrites map[string]*pb.EventRow

	// Commits arrived before prewrites.
	commitsBeforePrewrites map[string]*pb.EventRow
}

func NewPrewriteBuffer(tenantID uint64, range_ *pb.Task_Range) *PrewriteBuffer {
	return &PrewriteBuffer{
		tenantID:  tenantID,
		Range:     range_,
		prewrites: make(map[string]*pb.EventRow),
	}
}

func (p *PrewriteBuffer) AddEvent(eventRow *pb.EventRow) *pb.EventRow {
	k := fmt.Sprintf("%s-%d", eventRow.Key, eventRow.StartTs)

	switch eventRow.OpType {
	case pb.OpType_PREWRITE:
		p.prewrites[k] = eventRow
	case pb.OpType_COMMIT:
		prewrite, ok := p.prewrites[k]
		if !ok {
			p.commitsBeforePrewrites[k] = eventRow
		} else {
			prewrite.CommitTs = eventRow.CommitTs
			prewrite.OpType = pb.OpType_COMMIT
			delete(p.prewrites, k)

			return prewrite
		}
	case pb.OpType_ROLLBACK:
		delete(p.prewrites, k)
	default:
		panic("unreachable")
	}

	return nil
}
