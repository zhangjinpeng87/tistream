package datapool

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/huandu/skiplist"
	"github.com/zhangjinpeng87/tistream/pkg/storage"
	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
)

type TenantDataMap struct {
	// The tenant id.
	tenantID uint64

	// Root Path of the tenant data.
	rootPath string

	// Backend Storage
	backendStorage storage.ExternalStorage

	// The skiplist to store the data files.
	// Range -> Ordered file list
	rangeFiles skiplist.SkipList
}

// NewTenantDataMap creates a new TenantDataMap.
func NewTenantDataMap(tenantID uint64, rootPath string, backendStorage storage.ExternalStorage) *TenantDataMap {
	return &TenantDataMap{
		tenantID:       tenantID,
		rootPath:       rootPath,
		backendStorage: backendStorage,
		rangeFiles:     *skiplist.New(skiplist.ByteAsc),
	}
}

func (m *TenantDataMap) GetRangeFiles(range_ *pb.Task_Range, lowWatermark uint64) []*RangeFiles {
	lowWatermarkStr := fmt.Sprintf("%020d", lowWatermark)

	ele := m.rangeFiles.Get(range_.Start)
	if ele == nil {
		ele = m.rangeFiles.Back()
		if ele != nil {
			rfiles := ele.Value.(*RangeFiles)
			if bytes.Compare(rfiles.Range_.End, range_.Start) > 0 && len(rfiles.Files) > 0 {
				idx := sort.Search(len(rfiles.Files), func(i int) bool {
					return rfiles.Files[i] > fmt.Sprintf("%020d", lowWatermark)
				})
				rf := &RangeFiles{
					Range_: range_,
					Files:  rfiles.Files[idx:],
				}
				return []*RangeFiles{rf}
			}
		}
		return nil
	} else {
		res := make([]*RangeFiles, 0)
		prev := ele.Prev()
		if prev != nil {
			rfiles := prev.Value.(*RangeFiles)
			if bytes.Compare(rfiles.Range_.End, range_.Start) > 0 && len(rfiles.Files) > 0 {
				idx := sort.Search(len(rfiles.Files), func(i int) bool {
					return rfiles.Files[i] > lowWatermarkStr
				})
				res = append(res, &RangeFiles{Range_: rfiles.Range_, Files: rfiles.Files[idx:]})
			}
		}

		for ele != nil {
			rfiles := ele.Value.(*RangeFiles)
			if bytes.Compare(rfiles.Range_.End, range_.Start) > 0 && len(rfiles.Files) > 0 {
				idx := sort.Search(len(rfiles.Files), func(i int) bool {
					return rfiles.Files[i] > lowWatermarkStr
				})
				res = append(res, &RangeFiles{Range_: rfiles.Range_, Files: rfiles.Files[idx:]})
			}
			ele = ele.Next()
		}

		return res
	}
}

func (m *TenantDataMap) AddRangeFiles(range_ *pb.Task_Range, files ...string) {
	// Todo: handle task split & merge cases.
	ele := m.rangeFiles.Get(range_.Start)
	if ele == nil {
		rfiles := NewRangeFiles(range_)
		rfiles.AddFiles(files...)
		m.rangeFiles.Set(range_.Start, rfiles)
	} else {
		rfiles := ele.Value.(*RangeFiles)
		rfiles.AddFiles(files...)
	}
}

type RangeFiles struct {
	// The range of the data files.
	Range_ *pb.Task_Range

	// Ordered data files, the file name is the max commit ts of the file.
	Files []string
}

// NewRangeFiles creates a new RangeFiles.
func NewRangeFiles(range_ *pb.Task_Range) *RangeFiles {
	return &RangeFiles{
		Range_: range_,
		Files:  make([]string, 0),
	}
}

func (r *RangeFiles) AddFiles(file ...string) {
	r.Files = append(r.Files, file...)
}
