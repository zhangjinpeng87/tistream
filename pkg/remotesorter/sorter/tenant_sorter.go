package sorter

import (
	"fmt"

	"github.com/huandu/skiplist"
	"github.com/zhangjinpeng87/tistream/pkg/storage"
	"github.com/zhangjinpeng87/tistream/pkg/utils"
	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
)

// TenantSorter is the sorter of the tenant.
type TenantSorter struct {
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

func (s *TenantSorter) AddRange(range_ *pb.Task_Range) error {
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

	// Add the range sorter.
	s.rangeSorters.Set(range_.Start, NewRangeSorter(s.TenantID, range_, s.externalStorage))
}
