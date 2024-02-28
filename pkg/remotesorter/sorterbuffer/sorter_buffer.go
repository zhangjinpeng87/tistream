package sorterbuffer

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"math"

	"github.com/zhangjinpeng87/tistream/pkg/codec"
	"github.com/zhangjinpeng87/tistream/pkg/storage"
	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
)

// SorterBuffer is the buffer of committed data.
type SorterBuffer struct {
	// The tenant id.
	tenantID uint64

	// range of the buffer.
	Range *pb.Task_Range

	// The ordered event map.
	orderedEventMap OrderedEventMap

	// backend storage to store committed data files.
	externalStorage storage.ExternalStorage
}

// NewSorterBuffer creates a new SorterBuffer.
func NewSorterBuffer(tenantID uint64, range_ *pb.Task_Range, s storage.ExternalStorage) *SorterBuffer {
	return &SorterBuffer{
		tenantID:        tenantID,
		Range:           range_,
		orderedEventMap: NewSkipListEventMap(),
		externalStorage: s,
	}
}

func (s *SorterBuffer) AddEvent(eventRow *pb.EventRow) {
	key := fmt.Sprintf("%020d-%s", eventRow.CommitTs, eventRow.Key)
	s.orderedEventMap.Put(key, eventRow)
}

func (s *SorterBuffer) FlushCommittedData(lowWatermark, highWatermark uint64, path string) error {
	events := s.orderedEventMap.GetRange(fmt.Sprintf("%020d-", lowWatermark), fmt.Sprintf("%020d-", highWatermark))
	if len(events) > 0 {

		encoder := codec.NewCommittedDataEncoder(s.tenantID, s.Range, lowWatermark, highWatermark)
		var buf bytes.Buffer
		w := bufio.NewWriter(&buf)
		if err := encoder.Encode(w, events); err != nil {
			return err
		}
		w.Flush()

		// append checksum
		checksum := codec.CalcChecksum(buf.Bytes())
		if err := binary.Write(w, binary.LittleEndian, checksum); err != nil {
			return err
		}
		w.Flush()

		// write to storage
		xPath := fmt.Sprintf("%s/%d-%s-%s/%020d-%020d", path, s.tenantID, s.Range.Start, s.Range.End, lowWatermark, highWatermark)

		if err := s.externalStorage.PutFile(xPath, buf.Bytes()); err != nil {
			// Todo: retry in case of external storage temporarily unavailable.
			return err
		}
	}

	return nil
}

func (s *SorterBuffer) CleanData(lowWatermark, highWatermark uint64) {
	s.orderedEventMap.DelRange(fmt.Sprintf("%020d-", lowWatermark), fmt.Sprintf("%020d-", highWatermark))
}

func (s *SorterBuffer) SaveSnapTo(path string) error {
	codec := codec.NewSorterBufferSnapEncoder(s.tenantID, s.Range)
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	minKey := fmt.Sprintf("%020d", uint64(0))
	maxKey := fmt.Sprintf("%020d", uint64(math.MaxUint64))
	if err := codec.Encode(w, s.orderedEventMap.GetRange(minKey, maxKey)); err != nil {
		return err
	}
	w.Flush()

	// append checksum
	checksum := codec.CalcChecksum(buf.Bytes())
	if err := binary.Write(w, binary.LittleEndian, checksum); err != nil {
		return err
	}
	w.Flush()

	// write to storage
	xPath := fmt.Sprintf("%s/%d-%s-%s", path, s.tenantID, s.Range.Start, s.Range.End)

	if err := s.externalStorage.PutFile(xPath, buf.Bytes()); err != nil {
		return err
	}

	return nil
}

func (s *SorterBuffer) LoadSnapFrom(path string) error {
	// Reset the ordered event map.
	s.orderedEventMap.Reset()

	// Load the snap from storage.
	xPath := fmt.Sprintf("%s/%d-%d", path, s.Range.Start, s.Range.End)
	data, err := s.externalStorage.GetFile(xPath)
	if err != nil {
		return err
	}

	expectedChecksum := binary.LittleEndian.Uint32(data[len(data)-4:])
	checksum := codec.CalcChecksum(data[:len(data)-4])
	if expectedChecksum != checksum {
		return fmt.Errorf("file %s checksum not match", xPath)
	}

	// Decode the snap.
	decoder := codec.NewSorterBufferSnapDecoder(s.tenantID, s.Range)
	reader := bytes.NewReader(data[:len(data)-4])
	events, err := decoder.Decode(reader)
	if err != nil {
		return err
	}

	for _, event := range events {
		s.orderedEventMap.Put(fmt.Sprintf("%020d-%s", event.CommitTs, event.Key), event)
	}

	return nil
}
