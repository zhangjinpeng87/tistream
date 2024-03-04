package datapool

import (
	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
)

// DataStream is a wrapper of data stream for a specified range, the range can be a table or a whole tenant.
type DataStream struct {
	// The tenant id.
	tenantID uint64
	// The range of the data stream.
	Range *pb.Task_Range
	// The low watermark of the data stream.
	lowWatermark uint64

	files []*RangeFiles
}

// NewDataStream creates a new data stream.
func NewDataStream(tenantID uint64, range_ *pb.Task_Range, lowWatermark uint64) *DataStream {
	return &DataStream{
		tenantID:     tenantID,
		Range:        range_,
		lowWatermark: lowWatermark,
	}
}

func (d *DataStream) Next() *pb.EventRow {
	return nil
}
