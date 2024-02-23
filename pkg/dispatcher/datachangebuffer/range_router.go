package datachangebuffer

import (
	"bytes"
	"sort"

	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
)

type Range struct {
	Start []byte // base64 encoded
	End   []byte // base64 encoded
}

type RangeChanges struct {
	Range      *Range
	Rows       []*pb.EventRow
	Watermarks []*pb.EventWatermark
}

// RouteByRange routes the data changes by the range.
// SortedRanges is the sorted ranges to route the data changes, there is no overlap between the ranges.
func RouteByRange(sortedRanges []*Range, eventRows []*pb.EventRow, eventWatermarks []*pb.EventWatermark) []*RangeChanges {
	if len(eventRows) == 0 && len(eventWatermarks) == 0 {
		return nil
	}

	// Sort the event rows by its key and start ts
	sort.Slice(eventRows, func(i, j int) bool {
		if !bytes.Equal(eventRows[i].Key, eventRows[j].Key) {
			return bytes.Compare(eventRows[i].Key, eventRows[j].Key) < 0
		}
		if eventRows[i].CommitTs != eventRows[j].CommitTs {
			return eventRows[i].CommitTs < eventRows[j].CommitTs
		}
		return eventRows[i].StartTs < eventRows[j].StartTs
	})

	// Sort the event watermarks by its range and ts.
	sort.Slice(eventWatermarks, func(i, j int) bool {
		if !bytes.Equal(eventWatermarks[i].RangeStart, eventWatermarks[j].RangeStart) {
			return bytes.Compare(eventWatermarks[i].RangeStart, eventWatermarks[j].RangeStart) < 0
		}
		if eventWatermarks[i].RangeVersion != eventWatermarks[j].RangeVersion {
			return eventWatermarks[i].RangeVersion < eventWatermarks[j].RangeVersion
		}
		return eventWatermarks[i].Ts < eventWatermarks[j].Ts
	})

	res := make([]*RangeChanges, len(sortedRanges))
	for i := range res {
		res[i] = &RangeChanges{
			Range:      sortedRanges[i],
			Rows:       make([]*pb.EventRow, 0, len(eventRows)),
			Watermarks: make([]*pb.EventWatermark, 0, len(eventWatermarks)),
		}
	}

	// Dispatch the event rows to the ranges.
	idx := 0
	for _, r := range eventRows {
		// Advance the current range index to the range that the event row belongs to.
		for idx < len(sortedRanges) && bytes.Compare(r.Key, sortedRanges[idx].End) >= 0 {
			idx++
		}
		if idx >= len(sortedRanges) {
			break
		}

		// Skip the event row if it is before the start of the current range.
		if bytes.Compare(r.Key, sortedRanges[idx].Start) < 0 {
			continue
		}

		// Append the event row to the current range changes.
		res[idx].Rows = append(res[idx].Rows, r)
	}

	idx = 0
	for _, w := range eventWatermarks {
		// Advance the current range index to the range that the event row belongs to.
		for idx < len(sortedRanges) && bytes.Compare(w.RangeStart, sortedRanges[idx].End) >= 0 {
			idx++
		}

		if idx >= len(sortedRanges) {
			break
		}

		// Skip the event row if it is before the start of the current range.
		if bytes.Compare(w.RangeStart, sortedRanges[idx].Start) < 0 {
			continue
		}

		// Append the event row to the current range changes.
		res[idx].Watermarks = append(res[idx].Watermarks, w)
	}

	return res
}
