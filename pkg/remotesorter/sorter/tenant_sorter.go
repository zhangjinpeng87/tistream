package sorter

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/huandu/skiplist"
	"github.com/zhangjinpeng87/tistream/pkg/storage"
	"github.com/zhangjinpeng87/tistream/pkg/utils"
	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
)

// TenantSorter is the sorter of the tenant.
type TenantSorter struct {
	sync.RWMutex

	// The tenant id.
	TenantID uint64

	// RangeSorters of the tenant.
	rangeSorters *skiplist.SkipList

	// external storage
	externalStorage storage.ExternalStorage
}

// NewTenantSorter creates a new TenantSorter.
func NewTenantSorter(tenantID uint64, es storage.ExternalStorage) *TenantSorter {
	return &TenantSorter{
		TenantID:        tenantID,
		rangeSorters:    skiplist.New(skiplist.BytesAsc),
		externalStorage: es,
	}
}

func (s *TenantSorter) AddEventBatch(tenantId uint64, task_range *pb.Task_Range, eventBatch *pb.EventBatch) error {
	if s.TenantID != tenantId {
		return fmt.Errorf("tenant id not match")
	}

	// Get the range sorter.
	s.Lock()
	ele := s.rangeSorters.Find(task_range.Start)
	if ele == nil {
		return fmt.Errorf("range not found")
	}
	rangeSorter := ele.Value.(*RangeSorter)
	if bytes.Compare(rangeSorter.Range.End, task_range.End) != 0 {
		return fmt.Errorf("range not match")
	}
	s.Unlock()

	return rangeSorter.AddEventBatch(eventBatch)
}

func (s *TenantSorter) AddNewRange(range_ *pb.Task_Range) error {
	// Check range conflict.
	if err := s.CheckRangeOverlap(range_); err != nil {
		return err
	}

	// Add the range sorter.
	s.Lock()
	defer s.Unlock()
	s.rangeSorters.Set(range_.Start, NewRangeSorter(s.TenantID, range_, s.externalStorage))

	return nil
}

func (s *TenantSorter) AddRangeFromSnap(range_ *pb.Task_Range) error {
	// Check range conflict.
	if err := s.CheckRangeOverlap(range_); err != nil {
		return err
	}

	// Init the range sorter.
	rangeSorter := NewRangeSorter(s.TenantID, range_, s.externalStorage)
	if err := rangeSorter.LoadSnapshot(); err != nil {
		return err
	}

	// Add the range sorter.
	s.Lock()
	defer s.Unlock()
	s.rangeSorters.Set(range_.Start, rangeSorter)

	return nil
}

func (s *TenantSorter) CheckRangeOverlap(range_ *pb.Task_Range) error {
	s.RLock()
	defer s.RUnlock()

	// Check range conflict.
	ele := s.rangeSorters.Find(range_.Start)
	if ele != nil {
		if utils.RangeOverlapped(ele.Value.(*RangeSorter).Range, range_) {
			return fmt.Errorf("range conflict")
		}

		prev := ele.Prev()
		if prev != nil {
			if utils.RangeOverlapped(prev.Value.(*RangeSorter).Range, range_) {
				return fmt.Errorf("range conflict")
			}
		}
	} else {
		last := s.rangeSorters.Back()
		if last != nil {
			if utils.RangeOverlapped(last.Value.(*RangeSorter).Range, range_) {
				return fmt.Errorf("range conflict")
			}
		}
	}

	return nil
}

func (s *TenantSorter) SaveSnapshot() error {
	// Todo: just lock when get these sorters and release lock before savesnapshot
	// since savesnapshot may take several seconds.
	s.RLock()
	defer s.RUnlock()

	ele := s.rangeSorters.Front()
	for ele != nil {
		rangeSorter := ele.Value.(*RangeSorter)
		if err := rangeSorter.SaveSnapshot(); err != nil {
			return err
		}
		ele = ele.Next()
	}
	return nil
}

func (s *TenantSorter) FlushCommittedData() error {
	// Todo: just lock when get these sorters and release lock before flush,
	// since flush may take several seconds.
	s.RLock()
	defer s.RUnlock()

	ele := s.rangeSorters.Front()
	for ele != nil {
		rangeSorter := ele.Value.(*RangeSorter)
		if err := rangeSorter.FlushCommittedData(); err != nil {
			return err
		}
		ele = ele.Next()
	}
	return nil
}
