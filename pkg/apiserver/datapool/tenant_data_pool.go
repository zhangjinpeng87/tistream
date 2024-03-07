package datapool

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"sort"
	"sync/atomic"

	"github.com/huandu/skiplist"
	"github.com/zhangjinpeng87/tistream/pkg/codec"
	"github.com/zhangjinpeng87/tistream/pkg/storage"
	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
)

// The Committed Data Pool of a tenant is organized as follows:
// |
// |____committed_data_pool/Tenant-{1}
// |  |____manifest (this file only can be modified by meta-server during range split/merge, it contains range snapshots of different time slices)
// |  |____{range1-uuid}
// |  |  |____meta (this file is maintained by the sorter who own this task, contains all sub directories like watermark1, watermark2...)
// |  |  |____watermark1/
// |  |  |  |____file-{min-ts}-{max-ts}
// |  |  |  |____file-{min-ts}-{max-ts}
// |  |  |____watermark2/
// |  |  |  |____file-{min-ts}-{max-ts}
// |  |  |  |____file-{min-ts}-{max-ts}
// |  |  |____watermark3/
// |  |  |  |____file-{min-ts}-{max-ts}
// |  |  |  |____...
// |  |____{range2-uuid}
// |  |  |____meta (this file is maintained by the sorter who own this task)
// |  |  |____watermark1/
// |  |  |  |____file-{min-ts}-{max-ts}
// |  |  |  |____file-{min-ts}-{max-ts}

// TenantDataPool is the memory index for the tenant. It is a 2 dimensions map:
// 1. The first dimension is the watermark. (watermark -> RangesFilesSnap)
// 2. The second dimension is the range. ()
type TenantDataPool struct {
	// The tenant id.
	tenantID uint64

	// Root Path of the tenant data.
	rootPath string

	// Backend Storage
	backendStorage storage.ExternalStorage

	// The skiplist to store the data files. It is a 2 dimensions map:
	// Time Slices (highWatermark) -> RangesFilesSnap
	rangeFiles skiplist.SkipList

	inited atomic.Bool
}

// NewTenantDataPool creates a new TenantDataPool.
func NewTenantDataPool(tenantID uint64, rootPath string, backendStorage storage.ExternalStorage) *TenantDataPool {
	return &TenantDataPool{
		tenantID:       tenantID,
		rootPath:       rootPath,
		backendStorage: backendStorage,
		rangeFiles:     *skiplist.New(skiplist.ByteAsc),
		inited:         atomic.Bool{},
	}
}

// GetRangeFiles returns the data files for a specified logical range.
// Returns
func (m *TenantDataPool) GetRangeFiles(range_ *pb.Task_Range, lowWatermark uint64, fileCntLimit int) map[uint64][]*RangeFiles {
	// Iterate the RangesFilesSnap from lowWatermark to get the data files for the specified logical range.
	ele := m.rangeFiles.Find(lowWatermark)
	if ele == nil {
		return nil
	}

	res := make(map[uint64][]*RangeFiles)
	totalFileCnt := 0
	for ele != nil {
		tsrfiles := ele.Value.(*RangesFilesSnap)
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

func (m *TenantDataPool) Init() error {
	// Read the manifest file
	fName := fmt.Sprintf("%s/%s", m.rootPath, codec.CommittedManifestName)
	data, err := m.backendStorage.GetFile(fName)
	if err != nil {
		return err
	}

	// check the checksum
	checksum := codec.CalcChecksum(data[:len(data)-4])
	expectedChecksum := binary.LittleEndian.Uint32(data[len(data)-4:])
	if checksum != expectedChecksum {
		return fmt.Errorf("checksum not match, expected: %d, actual: %d", expectedChecksum, checksum)
	}

	// Decode the manifest
	r := bytes.NewReader(data[:len(data)-4])
	manifestDecoder := codec.NewCommittedManifestDecoder(m.tenantID)
	snapshots, err := manifestDecoder.Decode(r)
	for i, s := range snapshots {
		low := s.Ts
		high := uint64(math.MaxUint64)
		if i < len(snapshots)-1 {
			high = snapshots[i+1].Ts
		}
		rfs, err := NewRangesFilesSnapFromStorage(s.Ranges, low, high, m.tenantID, m.rootPath, m.backendStorage)
		if err != nil {
			return err
		}
		m.rangeFiles.Set(high, rfs)
	}

	m.inited.Store(true)

	return nil
}

// RangesFilesSnap is RangeFiles for a specific time slice like ts100~200.
type RangesFilesSnap struct {
	// Time slice of the data files.
	HighWatermark uint64
	LowWatermark  uint64

	// Data Range -> Ts Range -> Files
	// | 	Range1       |       Range2       |       Range3       |
	// |-----------------|--------------------|--------------------|
	// | file1, file2    | file3, file4       | file5, file6       | ts100~200
	// | file7, file8    | file9, file10      | file11, file12     | ts200~300
	rangeFiles skiplist.SkipList
}

// NewRangesFilesSnap creates a new RangesFilesSnap.
func NewRangesFilesSnap(lowWatermark, highWatermark uint64) *RangesFilesSnap {
	return &RangesFilesSnap{
		LowWatermark:  lowWatermark,
		HighWatermark: highWatermark,
		rangeFiles:    *skiplist.New(skiplist.ByteAsc),
	}
}

func NewRangesFilesSnapFromStorage(ranges []*pb.Task_Range, lowWatermark, highWatermark uint64, tenantId uint64,
	rootPath string, s storage.ExternalStorage) (*RangesFilesSnap, error) {

	t := NewRangesFilesSnap(lowWatermark, highWatermark)

	for _, r := range ranges {
		rangeIndexFileName := fmt.Sprintf("%s/%s/%s", rootPath, r.Uuid, codec.CommittedRangeIndexName)
		data, err := s.GetFile(rangeIndexFileName)
		if err != nil {
			return nil, err
		}

		// The range index file contains sub directories of the range.
		subDirs, err := codec.GetCommittedRangeIndex(data, tenantId, r)
		if err != nil {
			return nil, err
		}

		rtf := NewTimeRangeFiles(r)
		// Iterate the sub directories to get the data files.
		for i, subDir := range subDirs {
			low := subDir
			high := uint64(math.MaxUint64)
			if i < len(subDirs)-1 {
				high = subDirs[i+1]
			}
			if high < lowWatermark {
				// The range is before the lowWatermark, skip it.
				continue
			}
			if low > highWatermark {
				// The range is after the highWatermark, skip it.
				break
			}

			subDirName := fmt.Sprintf("%s/%s/%s", rootPath, r.Uuid, subDir)
			list, err := s.ListSubDir(subDirName)
			if err != nil {
				return nil, err
			}
			if len(list) > 0 {
				// Sort the files
				// Todo: filter files by lowWatermark and highWatermark
				sort.Strings(list)

				rf := NewRangeFiles(r, low, high)
				rf.AddFiles(list...)
				rtf.AddRangeFiles(rf)
			}
		}
		if rtf.Len() > 0 {
			t.rangeFiles.Set(r.End, rtf)
		}
	}

	return t, nil
}

func (t *RangesFilesSnap) Len() int {
	return t.rangeFiles.Len()
}

func (t *RangesFilesSnap) AddRangeFiles(rf *TimeRangeFiles) {
	t.rangeFiles.Set(rf.Range.End, rf)
}

// GetRangeFiles returns the data files for a specified logical range.
// The logical range can be a table or a whole tenant.
func (t *RangesFilesSnap) GetRangeFiles(range_ *pb.Task_Range, lowWatermark uint64) []*RangeFiles {
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

type TimeRangeFiles struct {
	Range *pb.Task_Range

	// ts -> RangeFiles
	Files skiplist.SkipList
}

func NewTimeRangeFiles(range_ *pb.Task_Range) *TimeRangeFiles {
	return &TimeRangeFiles{
		Range: range_,
		Files: *skiplist.New(skiplist.Uint64Asc),
	}
}

func (r *TimeRangeFiles) AddRangeFiles(rf *RangeFiles) {
	r.Files.Set(rf.HighWatermark, rf)
}

func (r *TimeRangeFiles) Len() int {
	return r.Files.Len()
}

type RangeFiles struct {
	// The range of the data files.
	Range_ *pb.Task_Range

	// Time Range of the data files
	LowWatermark  uint64
	HighWatermark uint64

	// Ordered data files, the file name is the max commit ts of the file.
	Files []string
}

// NewRangeFiles creates a new RangeFiles.
func NewRangeFiles(range_ *pb.Task_Range, low, high uint64) *RangeFiles {
	return &RangeFiles{
		Range_:        range_,
		LowWatermark:  low,
		HighWatermark: high,
		Files:         make([]string, 0),
	}
}

func (r *RangeFiles) AddFiles(file ...string) {
	r.Files = append(r.Files, file...)
}
