package sorter

import (
	"bytes"
	"math"
	"sync"

	"github.com/huandu/skiplist"
	"github.com/zhangjinpeng87/tistream/pkg/utils"
	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
)

type RangeWatermarks struct {
	sync.RWMutex
	Range *pb.Task_Range

	// range-tree organized watermarks.
	// RangeStart -> EventWatermark
	watermarks *skiplist.SkipList
}

func NewRangeWatermarks(range_ *pb.Task_Range) *RangeWatermarks {
	return &RangeWatermarks{
		Range:      range_,
		watermarks: skiplist.New(skiplist.BytesAsc),
	}
}

// UpdateWatermark updates the watermark of the range.
// Todo: currently we care more about the corretness of overlapping ranges,
// we may need to optimize the performance later.
func (r *RangeWatermarks) UpdateWatermark(wm *pb.EventWatermark) error {
	r.Lock()
	defer r.Unlock()

	// Check the range of the watermark.
	if bytes.Compare(wm.RangeStart, r.Range.Start) < 0 || bytes.Compare(wm.RangeEnd, r.Range.End) > 0 {
		// Should not happen.
		return utils.ErrInvalidRange
	}

	// Get all existing watermarks that overlap with the new watermark.
	overlaps := r.getOrderedOverlaps(wm)
	if len(overlaps) == 0 {
		// No overlap, add the new watermark.
		r.watermarks.Set(wm.RangeStart, wm)
		return nil
	}

	var updateRanges []*pb.EventWatermark

	// Case 2: if there is a gap at the beginning of the range, add the gap.
	if bytes.Compare(wm.RangeStart, overlaps[0].RangeStart) < 0 {
		// Add the gap.
		updateRanges = append(updateRanges, &pb.EventWatermark{
			RangeId:      wm.RangeId,
			RangeVersion: wm.RangeVersion,
			RangeStart:   wm.RangeStart,
			RangeEnd:     overlaps[0].RangeStart,
			Ts:           wm.Ts,
		})
	}

	lastEnd := overlaps[0].RangeStart
	for _, overlap := range overlaps {
		updateRanges = append(updateRanges, r.calcWatermarkUpdates(lastEnd, overlap, wm)...)
		lastEnd = overlap.RangeEnd
	}

	// Case 11: if there is a gap at the end of the range, add the gap.
	if bytes.Compare(lastEnd, wm.RangeEnd) < 0 {
		// Add the gap.
		updateRanges = append(updateRanges, &pb.EventWatermark{
			RangeId:      wm.RangeId,
			RangeVersion: wm.RangeVersion,
			RangeStart:   lastEnd,
			RangeEnd:     wm.RangeEnd,
			Ts:           wm.Ts,
		})
	}

	// Merge updateRanges into the existing watermarks.
	finalUpdateRanges := mergeContinuousRanges(updateRanges)
	for _, olvelap := range overlaps {
		r.watermarks.Remove(olvelap.RangeStart)
	}
	for _, range_ := range finalUpdateRanges {
		r.watermarks.Set(range_.RangeStart, range_)
	}

	return nil
}

func mergeContinuousRanges(ranges []*pb.EventWatermark) []*pb.EventWatermark {
	var res []*pb.EventWatermark
	lastRange := ranges[0]
	for i, range_ := range ranges {
		if i == 0 {
			continue
		}

		if bytes.Compare(lastRange.RangeEnd, range_.RangeStart) != 0 {
			panic("invalid ranges")
		}

		if lastRange.RangeId == range_.RangeId &&
			lastRange.RangeVersion == range_.RangeVersion &&
			lastRange.Ts == range_.Ts {
			lastRange.RangeEnd = range_.RangeEnd
		} else {
			res = append(res, lastRange)
			lastRange = range_
		}
	}
	res = append(res, lastRange)

	return res
}

func (r *RangeWatermarks) calcWatermarkUpdates(lastEnd []byte, overlap, wm *pb.EventWatermark) []*pb.EventWatermark {
	var updateRanges []*pb.EventWatermark
	// There is a gap between the last exiting watermark and the current existing watermark.
	if bytes.Compare(lastEnd, overlap.RangeStart) < 0 {
		// Add the gap.
		updateRanges = append(updateRanges, &pb.EventWatermark{
			RangeStart:   lastEnd,
			RangeEnd:     overlap.RangeStart,
			RangeId:      wm.RangeId,
			RangeVersion: wm.RangeVersion,
			Ts:           wm.Ts,
		})

		lastEnd = overlap.RangeStart
	}

	// Keep the overlap's left part which is not covered by the current watermark.
	if bytes.Compare(overlap.RangeStart, wm.RangeStart) < 0 {
		updateRanges = append(updateRanges, &pb.EventWatermark{
			RangeStart:   overlap.RangeStart,
			RangeEnd:     wm.RangeStart,
			RangeId:      overlap.RangeId,
			RangeVersion: overlap.RangeVersion,
			Ts:           overlap.Ts,
		})
		lastEnd = wm.RangeStart
	}

	// Compare range version and ts
	end := overlap.RangeEnd
	if bytes.Compare(wm.RangeEnd, overlap.RangeEnd) < 0 {
		end = wm.RangeEnd
	}
	if wm.RangeVersion > overlap.RangeVersion || (wm.RangeVersion == overlap.RangeVersion && wm.Ts > overlap.Ts) {
		updateRanges = append(updateRanges, &pb.EventWatermark{
			RangeStart:   lastEnd,
			RangeEnd:     end,
			RangeId:      wm.RangeId,
			RangeVersion: wm.RangeVersion,
			Ts:           wm.Ts,
		})
	} else {
		updateRanges = append(updateRanges, &pb.EventWatermark{
			RangeStart:   lastEnd,
			RangeEnd:     end,
			RangeId:      overlap.RangeId,
			RangeVersion: overlap.RangeVersion,
			Ts:           overlap.Ts,
		})
	}
	lastEnd = end

	if bytes.Compare(lastEnd, overlap.RangeEnd) < 0 {
		// Add the gap.
		updateRanges = append(updateRanges, &pb.EventWatermark{
			RangeStart:   lastEnd,
			RangeEnd:     overlap.RangeEnd,
			RangeId:      overlap.RangeId,
			RangeVersion: overlap.RangeVersion,
			Ts:           overlap.Ts,
		})
	}

	return updateRanges
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
	// 2) |      r    |                     : r1 + new range
	// 3)     |  r    |                     : r1
	// 4)     |  r          |               : r1 and r2
	// 5)     | r |                         : part of r1
	// 6)         | r |                     : part of r1
	// 7)         |   r  |                  : part of r1 and part of r2
	// 8)         |   r     |               : part of r1 and full r2
	// 9)                   | r |           : part of r3
	// 10)                      | r |       : part of last range
	// 11)                      | r     |   : r3 and new range
	// 12)                          | r |   : totally new range

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

func (r *RangeWatermarks) RangeWatermark() uint64 {
	r.RLock()
	defer r.RUnlock()

	// Use the minimal ts of all ranges as the range watermark.
	// If there is a gap between two ranges, it means the whole range is not ready,
	// and the watermark is 0.
	ele := r.watermarks.Front()
	if ele == nil {
		return 0
	}

	res := uint64(math.MaxUint64)
	lastEnd := r.Range.Start
	for ele != nil {
		wm := ele.Value().(*pb.EventWatermark)
		if bytes.Compare(wm.RangeStart, lastEnd) != 0 {
			return 0
		}
		lastEnd = wm.RangeEnd
		if wm.Ts < res {
			res = wm.Ts
		}

		ele = ele.Next()
	}

	if bytes.Compare(lastEnd, r.Range.End) < 0 {
		// There is a gap between the last watermark and the end of the range.
		return 0
	}

	return res
}

func (r *RangeWatermarks) SaveSnapTo() {
}

func (r *RangeWatermarks) LoadSnapFrom() {
}
