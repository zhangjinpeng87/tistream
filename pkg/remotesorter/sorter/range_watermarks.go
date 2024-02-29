package sorter

import (
	"bytes"

	"github.com/huandu/skiplist"
	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
)

type RangeWatermarks struct {
	Range *pb.Task_Range

	// B-Tree organized watermarks.
	watermarks *skiplist.SkipList
}

func NewRangeWatermarks(range_ *pb.Task_Range) *RangeWatermarks {
	return &RangeWatermarks{
		Range:      range_,
		watermarks: skiplist.New(skiplist.BytesAsc),
	}
}

func (r *RangeWatermarks) AddWatermark(wm *pb.EventWatermark) {
	overlaps := r.getOrderedOverlaps(wm)
	if len(overlaps) == 0 {
		// No overlap, add the new watermark.
		r.watermarks.Set(wm.RangeStart, wm)
		return
	}

	for _, overlap := range overlaps {
		r.mergeWatermark(overlap, wm)
	}
}

func (r *RangeWatermarks) mergeWatermark(overlap, wm *pb.EventWatermark) {
	// Merge the two watermarks.
	if overlap.RangeVersion > wm.RangeVersion {

	}
}

func (r *RangeWatermarks) getOrderedOverlaps(wm *pb.EventWatermark) []*pb.EventWatermark {
	// Ranges can be mreged or split.
	// |   r1    |      inital status
	// | r1 | r2 |      after split
	// |    r1   |      after merge
	// |    r1       |  after merge

	// There are several different cases to handle:
	//        |  r1   | r2  | r3    |       : existing ranges
	// 1) | r |                             : totally new range
	// 2)     |  r    |                     : r1
	// 3)     |  r          |               : r1 and r2
	// 4)     | r |                         : part of r1
	// 5)         | r |                     : part of r1
	// 6)         |   r  |                  : part of r1 and part of r2
	// 7)         |   r     |               : part of r1 and full r2
	// 8)                   | r |           : part of r3
	// 9)                       | r |       : part of last range
	// 10)                      | r     |   : r3 and new range
	// 11)                          | r |   : totally new range

	// Find the first watermark that its start is greater than or equal to the given watermark's start.
	curEle := r.watermarks.Find(wm.RangeStart)
	if curEle != nil {
		// case 1, 2, 3, 4, 5, 6, 7, 8
		var overlaps []*pb.EventWatermark

		preEle := curEle.Prev()
		if preEle != nil {
			preWm := preEle.Value().(*pb.EventWatermark)
			if bytes.Compare(preWm.End, wm.RangeStart) > 0 {
				overlaps = append(overlaps, preWm)
			}
		}

		for e := curEle; e != nil; e = e.Next() {
			curWm := e.Value().(*pb.EventWatermark)
			if bytes.Compare(curWm.Start, wm.RangeEnd) >= 0 {
				break
			}
			overlaps = append(overlaps, curWm)
		}
		return overlaps
	} else {
		backEle := r.watermarks.Back()
		if backEle == nil {
			// existing ranges is empty
			return nil
		}
		backWm := backEle.Value().(*pb.EventWatermark)
		if bytes.Compare(backWm.End, wm.RangeStart) <= 0 {
			// case 11, no overlap
			return nil
		} else {
			// case 9, 10
			return []*pb.EventWatermark{oldWm.Value().(*pb.EventWatermark)}
		}
	}
}
