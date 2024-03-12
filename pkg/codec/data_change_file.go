package codec

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
	TenantID uint64

	// All the event rows in the file.
	EventRows []*pb.EventRow

	// The event watermark in the file.
	EventWatermarks []*pb.EventWatermark

	// DDL changes
	DdlChanges []*pb.DDLChange
}

// NewDataChangesFileDecoder creates a new DataChanges.
func NewDataChangesFileDecoder(tenantId uint64) *DataChangesFileDecoder {
	return &DataChangesFileDecoder{
		TenantID:        tenantId,
		EventRows:       make([]*pb.EventRow, 0),
		EventWatermarks: make([]*pb.EventWatermark, 0),
		DdlChanges:      make([]*pb.DDLChange, 0),
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

		// Append the rows, watermarks and schema changes.
		if eventBatch.Rows != nil {
			f.EventRows = append(f.EventRows, eventBatch.Rows...)
		}
		if eventBatch.Watermarks != nil {
			f.EventWatermarks = append(f.EventWatermarks, eventBatch.Watermarks...)
		}
		if eventBatch.DdlChanges != nil {
			f.DdlChanges = append(f.DdlChanges, eventBatch.DdlChanges...)
		}
	}

	return nil
}
