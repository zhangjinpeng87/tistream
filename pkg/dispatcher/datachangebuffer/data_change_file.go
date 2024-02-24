package datachangebuffer

import (
	"encoding/binary"
	"io"

	"github.com/zhangjinpeng87/tistream/pkg/utils"
	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
	"google.golang.org/protobuf/proto"
)

const (
	// Magic header of the data change file.
	fileHeader = "TISTREAM_DATA_CHANGE_FILE"
	version    = 1
)

// Data change file is the file to store the data change.
// The file payload is as follows:
// Header (24 bytes) | File version (4 bytes) | Event size (4 bytes) | Event batch (event size bytes) | ... | Checksum (4 bytes)
type DataChangesFileDecoder struct {
	// All the event rows in the file.
	eventRows []*pb.EventRow

	// The event watermark in the file.
	eventWatermarks []*pb.EventWatermark
}

// NewDataChangesFileDecoder creates a new DataChanges.
func NewDataChangesFileDecoder() *DataChangesFileDecoder {
	return &DataChangesFileDecoder{
		eventRows:      make([]*pb.EventRow, 0),
		eventWatermarks: make([]*pb.EventWatermark, 0),
	}
}

// DecodeFrom decodes the data change from the reader.
func (f *DataChangesFileDecoder) DecodeFrom(reader io.Reader) error {
	// Read the header.
	header := make([]byte, len(fileHeader))
	if _, err := reader.Read(header); err != nil {
		return err
	}

	// Check the header.
	if string(header) != fileHeader {
		return utils.ErrInvalidDataChangeFile
	}

	// Read the file version.
	var fileVersion uint32
	if err := binary.Read(reader, binary.LittleEndian, &fileVersion); err != nil {
		return err
	}

	// TODO: decode the file according to the version.
	// Newer versions may support different compression algorithm, etc.

	// Read the event rows.
	for {
		eventBatch := &pb.EventBatch{}

		// Read the size of the event batch.
		var eventBatchSize uint32
		if err := binary.Read(reader, binary.LittleEndian, &eventBatchSize); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		// Read and decode the event batch.
		eventBatchBytes := make([]byte, eventBatchSize)
		if _, err := reader.Read(eventBatchBytes); err != nil {
			return err
		}
		if err := proto.Unmarshal(eventBatchBytes, eventBatch); err != nil {
			return err
		}

		// Append the rows and watermarks.
		if eventBatch.Rows != nil {
			f.eventRows = append(f.eventRows, eventBatch.Rows...)
		}
		if eventBatch.Watermark != nil {
			f.eventWatermarks = append(f.eventWatermarks, eventBatch.Watermark...)
		}
	}

	return nil
}
