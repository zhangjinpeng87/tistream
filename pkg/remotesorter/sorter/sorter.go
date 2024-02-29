package sorter

import (
	"github.com/zhangjinpeng87/tistream/pkg/remotesorter/storage"
	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
)

// Sorter is the sorter for a specified range of a specified tenant.
type Sorter struct {
	// The tenant id.
	TenantID uint64

	Range *pb.Task_Range

	// The prewrite buffer.
	PrewriteBuffer *PrewriteBuffer

	// The sorter buffer.
	SorterBuffer *SorterBuffer

	// External Storage
	ExternalStorage storage.ExternalStorage
}

// NewSorter creates a new Sorter.
func NewSorter(tenantID uint64, range_ *pb.Task_Range, es storage.ExternalStorage) *Sorter {
	return &Sorter{
		TenantID:       tenantID,
		Range:          range_,
		ExternalStorage: es,
		PrewriteBuffer: NewPrewriteBuffer(tenantID, range_, es),
		SorterBuffer:   NewSorterBuffer(tenantID, range_, es),
	}
}

// AddEvent adds an event to the sorter.
func (s *Sorter) AddEventBatch(eventBatch *pb.EventBatch) {
	// Handle all events rows in the batch.
	for _, eventRow := range eventBatch.Rows {
		switch eventRow.OpType {
		case pb.OpType_PREWRITE:
			s.PrewriteBuffer.AddEvent(eventRow)
		case pb.OpType_COMMIT:
			commitRow := s.PrewriteBuffer.AddEvent(eventRow)
			if commitRow != nil {
				s.SorterBuffer.AddEvent(commitRow)
			}
		case pb.OpType_ROLLBACK:
			s.PrewriteBuffer.RemoveEvent(eventRow)
		default:
			panic("unreachable")
		}
	}

	// Handle all watermarks in the batch.
	

}
