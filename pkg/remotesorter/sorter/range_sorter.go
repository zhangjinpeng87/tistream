package sorter

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/zhangjinpeng87/tistream/pkg/storage"
	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
)

// RangeSorter is the sorter for a specified range of a specified tenant.
type RangeSorter struct {
	// Mutex to protect the meta of sorter.
	sync.Mutex

	// The tenant id.
	TenantID uint64
	// Range of the sorter
	Range *pb.Task_Range

	// The prewrite buffer, sorter buffer and range watermarks.
	prewriteBuffer  *PrewriteBuffer
	sorterBuffer    *SorterBuffer
	rangeWatermarks *RangeWatermarks

	// External Storage
	ExternalStorage storage.ExternalStorage
	// rootPath is the root path of the external storage.
	rootPath string

	// Last checkpoint of the sorter.
	lastCheckpoint uint64
	// Last snapshot time of the sorter.
	lastSnapshotTime time.Time
	// Last update time of the sorter.
	// Monotonic time
	lastUpdateTime time.Time
}

// NewRangeSorter creates a new RangeSorter.
func NewRangeSorter(tenantID uint64, range_ *pb.Task_Range, es storage.ExternalStorage) *RangeSorter {
	return &RangeSorter{
		TenantID:         tenantID,
		Range:            range_,
		ExternalStorage:  es,
		prewriteBuffer:   NewPrewriteBuffer(tenantID, range_, es),
		sorterBuffer:     NewSorterBuffer(tenantID, range_, es, &SkiplistFactory{}),
		rangeWatermarks:  NewRangeWatermarks(tenantID, range_, es),
		lastSnapshotTime: time.Now(),
		lastUpdateTime:   time.Now(),
	}
}

// AddEvent adds an event to the sorter.
func (s *RangeSorter) AddEventBatch(eventBatch *pb.EventBatch) error {
	// Handle all events rows in the batch.
	for _, eventRow := range eventBatch.Rows {
		switch eventRow.OpType {
		case pb.OpType_PREWRITE, pb.OpType_ROLLBACK:
			s.prewriteBuffer.AddEvent(eventRow)
		case pb.OpType_COMMIT:
			commitRow := s.prewriteBuffer.AddEvent(eventRow)
			if commitRow != nil {
				s.sorterBuffer.AddEvent(commitRow)
			}
		default:
			panic("unreachable")
		}
	}

	// Handle all watermarks in the batch.
	for _, watermark := range eventBatch.Watermarks {
		if err := s.rangeWatermarks.UpdateWatermark(watermark); err != nil {
			return err
		}
	}

	// Update last update time.
	s.Lock()
	defer s.Unlock()
	s.lastUpdateTime = time.Now()

	return nil
}

func (s *RangeSorter) skipSnapshot() bool {
	s.Lock()
	defer s.Unlock()
	return s.lastSnapshotTime.After(s.lastUpdateTime)
}

func (s *RangeSorter) SaveSnapshot() error {
	if s.skipSnapshot() {
		return nil
	}

	if err := s.prewriteBuffer.SaveSnapTo(fmt.Sprintf("%s-prewritebuffer", s.rootPath)); err != nil {
		return err
	}
	if err := s.sorterBuffer.SaveSnapTo(fmt.Sprintf("%s-sorterbuffer", s.rootPath)); err != nil {
		return err
	}
	if err := s.rangeWatermarks.SaveSnapTo(fmt.Sprintf("%s-watermarks", s.rootPath)); err != nil {
		return err
	}

	s.Lock()
	defer s.Unlock()
	s.lastSnapshotTime = time.Now()

	return nil
}

func (s *RangeSorter) LoadSnapshot() error {
	if err := s.prewriteBuffer.LoadSnapFrom(fmt.Sprintf("%s-prewritebuffer", s.rootPath)); err != nil {
		return err
	}
	if err := s.sorterBuffer.LoadSnapFrom(fmt.Sprintf("%s-sorterbuffer", s.rootPath)); err != nil {
		return err
	}
	if err := s.rangeWatermarks.LoadSnapFrom(fmt.Sprintf("%s-watermarks", s.rootPath)); err != nil {
		return err
	}
	return nil
}

func (s *RangeSorter) FlushCommittedData() error {
	// Flush the committed data in the sorter buffer.
	latestWatermark := s.rangeWatermarks.LatestWatermark()
	if err := s.sorterBuffer.FlushCommittedData(s.lastCheckpoint, latestWatermark, s.rootPath); err != nil {
		return err
	}

	// Update the last checkpoint.
	s.Lock()
	s.lastCheckpoint = latestWatermark
	s.Unlock()

	// Clean the sorter buffer.
	s.sorterBuffer.CleanData(s.lastCheckpoint, latestWatermark)

	return nil
}

func (s *RangeSorter) Split(splitPoint []byte) (*RangeSorter, *RangeSorter) {
	s.Lock()
	defer s.Unlock()

	// Split the prewrite buffer.
	leftPrewriteBuffer, rightPrewriteBuffer := s.prewriteBuffer.Split(splitPoint)

	// Split the sorter buffer.
	leftSorterBuffer, rightSorterBuffer := s.sorterBuffer.Split(splitPoint)

	// Split the range watermarks.
	leftRangeWatermarks, rightRangeWatermarks := s.rangeWatermarks.Split(splitPoint)

	// Create the left and right range sorter.
	leftSorter := &RangeSorter{
		TenantID:        s.TenantID,
		Range:           &pb.Task_Range{Start: s.Range.Start, End: splitPoint},
		ExternalStorage: s.ExternalStorage,
		rootPath:        s.rootPath,
		prewriteBuffer:  leftPrewriteBuffer,
		sorterBuffer:    leftSorterBuffer,
		rangeWatermarks: leftRangeWatermarks,
	}
	rightSorter := &RangeSorter{
		TenantID:        s.TenantID,
		Range:           &pb.Task_Range{Start: splitPoint, End: s.Range.End},
		ExternalStorage: s.ExternalStorage,
		rootPath:        s.rootPath,
		prewriteBuffer:  rightPrewriteBuffer,
		sorterBuffer:    rightSorterBuffer,
		rangeWatermarks: rightRangeWatermarks,
	}
	return leftSorter, rightSorter
}

func (left *RangeSorter) MergeWith(right *RangeSorter) {
	left.Lock()
	defer left.Unlock()

	if left.TenantID != right.TenantID {
		panic("the ranges are not belongs to the same tenant")
	}
	if bytes.Compare(left.Range.End, right.Range.Start) != 0 {
		panic("the ranges are not adjacent")
	}

	// Merge the prewrite buffer.
	left.prewriteBuffer.MergeWith(right.prewriteBuffer)

	// Merge the sorter buffer.
	left.sorterBuffer.MergeWith(right.sorterBuffer)

	// Merge the range watermarks.
	left.rangeWatermarks.MergeWith(right.rangeWatermarks)
}
