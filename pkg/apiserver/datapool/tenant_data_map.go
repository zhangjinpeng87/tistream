package datapool

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/huandu/skiplist"
	"github.com/zhangjinpeng87/tistream/pkg/storage"
	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
)

// A tenant may have multiple data ranges, each range has its own data files ordered by the max commit ts of the file.
// These ranges may split or merge as time goes by, so we need to maintain the data files for each range.
// The committed data files of a tenant are organized as follows:
// |      Range1       |       Range2       |       Range3       |
// |-------------------|--------------------|--------------------|
// | file1, file2, ... | file1, file2, ...  | file1, file2, ...  | watermark: 0 ~ 100
// | file3, file4, ... | file3, file4, ...  | file3, file4, ...  | watermark: 100 ~ 200
// |                 Range1                 |       Range3       | // Range2 is merged into Range1
// | file5, file6, ...  file5, file6, ...   | file5, file6, ...  | watermark: 200 ~ 300
// | file7, file8, ...  file7, file8, ...   | file7, file8, ...  | watermark: 300 ~ 400
// |      Range1       |       Range2       |       Range3       | // Range2 is split from Range1
// | file9, file10, ...| file9, file10, ... | file9, file10, ... | watermark: 400 ~ 500
// ...
// TenantDataMap is the data map for a tenant. It is a 2 dimensions map:
// 1. The first dimension is the watermark.
// 2. The second dimension is the range.
type TenantDataMap struct {
	// The tenant id.
	tenantID uint64

	// Root Path of the tenant data.
	rootPath string

	// Backend Storage
	backendStorage storage.ExternalStorage

	// The skiplist to store the data files. It is a 2 dimensions map:
	// Time Slices (highWatermark) -> TimeSliceRangeFiles
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

// GetRangeFiles returns the data files for a specified logical range.
// Returns
func (m *TenantDataMap) GetRangeFiles(range_ *pb.Task_Range, lowWatermark uint64, fileCntLimit int) map[uint64][]*RangeFiles {
	// Iterate the TimeSliceRangeFiles from lowWatermark to get the data files for the specified logical range.
	ele := m.rangeFiles.Find(lowWatermark)
	if ele == nil {
		return nil
	}

	res := make(map[uint64][]*RangeFiles)
	totalFileCnt := 0
	for ele != nil {
		tsrfiles := ele.Value.(*TimeSliceRangeFiles)
		maxTs := ele.Key().(uint64)

		r := tsrfiles.GetRangeFiles(range_, lowWatermark)
		if r != nil {
			for _, rf := range r {
				totalFileCnt += len(rf.Files)
			}
			res[maxTs] = r
		}

		if totalFileCnt >= fileCntLimit {
			// We have got enough files. Next time set the maxTs as the lowWatermark to fetch more files.
			break
		}

		ele = ele.Next()
	}

	return res
}

// TimeSliceRangeFiles is RangeFiles for a specific time slice.
type TimeSliceRangeFiles struct {
	// Time slice of the data files.
	HighWatermark uint64
	LowWatermark  uint64

	// Different Ranges -> Ordered Files
	rangeFiles skiplist.SkipList
}

// NewTimeSliceRangeFiles creates a new TimeSliceRangeFiles.
func NewTimeSliceRangeFiles(lowWatermark, highWatermark uint64) *TimeSliceRangeFiles {
	return &TimeSliceRangeFiles{
		LowWatermark:  lowWatermark,
		HighWatermark: highWatermark,
		rangeFiles:    *skiplist.New(skiplist.ByteAsc),
	}
}

func (t *TimeSliceRangeFiles) AddRangeFiles(range_ *pb.Task_Range, files ...string) {
	ele := t.rangeFiles.Get(range_.Start)
	if ele == nil {
		rfiles := NewRangeFiles(range_)
		rfiles.AddFiles(files...)
		t.rangeFiles.Set(range_.Start, rfiles)
	} else {
		rfiles := ele.Value.(*RangeFiles)
		if bytes.Compare(rfiles.Range_.End, range_.End) != 0 {
			// If there is range split or merge, should use a new TimeSliceRangeFiles.
			panic(fmt.Sprintf("range end not match, range: %v, file: %v", rfiles.Range_, files))
		}
		rfiles.AddFiles(files...)
	}
}

// GetRangeFiles returns the data files for a specified logical range.
// The logical range can be a table or a whole tenant.
func (t *TimeSliceRangeFiles) GetRangeFiles(range_ *pb.Task_Range, lowWatermark uint64) []*RangeFiles {
	ele := t.rangeFiles.Find(range_.Start)
	if ele == nil {
		// Because we organized the range by the start of the range, so there is a case that the requested
		// logical range's start is larger than the start of the last range. We need to check the last range.
		// | 	Range1       |       Range2       |       Range3            |
		//                                          |requested logical range|
		ele = t.rangeFiles.Back()
		if ele != nil {
			rfiles := ele.Value.(*RangeFiles)
			if bytes.Compare(rfiles.Range_.End, range_.Start) > 0 {
				idx := sort.Search(len(rfiles.Files), func(i int) bool {
					return rfiles.Files[i] > fmt.Sprintf("%020d", lowWatermark)
				})
				return []*RangeFiles{&RangeFiles{Range_: rfiles.Range_, Files: rfiles.Files[idx:]}}
			}
		}

		return nil
	}

	res := make([]*RangeFiles, 0)

	// Check the previous range.
	prev := ele.Prev()
	if prev != nil {
		rfiles := prev.Value.(*RangeFiles)
		if bytes.Compare(rfiles.Range_.End, range_.Start) > 0 {
			// The requested logical range has overlap with the previous range.
			idx := sort.Search(len(rfiles.Files), func(i int) bool {
				return rfiles.Files[i] > fmt.Sprintf("%020d", lowWatermark)
			})
			res = append(res, &RangeFiles{Range_: rfiles.Range_, Files: rfiles.Files[idx:]})
		}
	}

	// Check the current range and the following ranges which has overlap with the
	// requested logical range.
	for ele != nil {
		rfiles := ele.Value.(*RangeFiles)
		if bytes.Compare(rfiles.Range_.Start, range_.End) < 0 {
			break
		}

		idx := sort.Search(len(rfiles.Files), func(i int) bool {
			return rfiles.Files[i] > fmt.Sprintf("%020d", lowWatermark)
		})
		res = append(res, &RangeFiles{Range_: rfiles.Range_, Files: rfiles.Files[idx:]})

		ele = ele.Next()
	}

	return res
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
