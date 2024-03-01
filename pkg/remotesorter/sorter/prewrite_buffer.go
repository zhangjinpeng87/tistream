package sorter

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/zhangjinpeng87/tistream/pkg/codec"
	"github.com/zhangjinpeng87/tistream/pkg/storage"
	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
)

type PrewriteBuffer struct {
	// The tenant id.
	tenantID uint64
	Range    *pb.Task_Range

	// All prewrites.
	// key: key-startTs
	// After a commit event arrives, the event will be moved to the SorterBuffer.
	prewrites map[string]*pb.EventRow
	// Commits arrived before prewrites.
	commitsBeforePrewrites map[string]*pb.EventRow

	// backend storage to store committed data files.
	externalStorage storage.ExternalStorage
}

func NewPrewriteBuffer(tenantID uint64, range_ *pb.Task_Range, es storage.ExternalStorage) *PrewriteBuffer {
	return &PrewriteBuffer{
		tenantID:               tenantID,
		Range:                  range_,
		prewrites:              make(map[string]*pb.EventRow),
		commitsBeforePrewrites: make(map[string]*pb.EventRow),
		externalStorage:        es,
	}
}

func (p *PrewriteBuffer) AddEvent(eventRow *pb.EventRow) *pb.EventRow {
	k := fmt.Sprintf("%s-%d", eventRow.Key, eventRow.StartTs)

	switch eventRow.OpType {
	case pb.OpType_PREWRITE:
		p.prewrites[k] = eventRow
	case pb.OpType_COMMIT:
		prewrite, ok := p.prewrites[k]
		if !ok {
			p.commitsBeforePrewrites[k] = eventRow
		} else {
			prewrite.CommitTs = eventRow.CommitTs
			prewrite.OpType = pb.OpType_COMMIT
			delete(p.prewrites, k)

			return prewrite
		}
	case pb.OpType_ROLLBACK:
		delete(p.prewrites, k)
	default:
		panic("unreachable")
	}

	return nil
}

func (p *PrewriteBuffer) SaveToSnap(rootPath string) error {
	snapName := fmt.Sprintf("%s/-%d-%s-%s.snap", rootPath, p.tenantID, p.Range.Start, p.Range.End)
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	encoder := codec.NewPrewriteBufferSnapEncoder(p.tenantID, p.Range)
	if err := encoder.Encode(w, p.prewrites, p.commitsBeforePrewrites); err != nil {
		return err
	}
	if err := w.Flush(); err != nil {
		return err
	}

	// Write the checksum.
	checksum := codec.CalcChecksum(buf.Bytes())
	if err := binary.Write(w, binary.LittleEndian, checksum); err != nil {
		return err
	}
	if err := w.Flush(); err != nil {
		return err
	}

	// Save the snap file to the external storage.
	if err := p.externalStorage.PutFile(snapName, buf.Bytes()); err != nil {
		return err
	}

	return nil
}

func (p *PrewriteBuffer) LoadFromSnap(rootPath string) error {
	// Fetch the snap file from the external storage.
	snapName := fmt.Sprintf("%s/-%d-%s-%s.snap", rootPath, p.tenantID, p.Range.Start, p.Range.End)
	snapContent, err := p.externalStorage.GetFile(snapName)
	if err != nil {
		return err
	}

	// Verify the checksum
	checksum := binary.LittleEndian.Uint32(snapContent[len(snapContent)-4:])
	if checksum != codec.CalcChecksum(snapContent[:len(snapContent)-4]) {
		return fmt.Errorf("checksum mismatch")
	}

	// Decode the snap file.
	decoder := codec.NewPrewriteBufferSnapDecoder(p.tenantID, p.Range)
	r := bytes.NewReader(snapContent[:len(snapContent)-4])
	prewrites, commits, err := decoder.Decode(r)
	if err != nil {
		return err
	}

	// Load the prewrites and commits.
	for _, prewrite := range prewrites {
		k := fmt.Sprintf("%s-%d", prewrite.Key, prewrite.StartTs)
		p.prewrites[k] = prewrite
	}
	for _, commit := range commits {
		k := fmt.Sprintf("%s-%d", commit.Key, commit.StartTs)
		p.commitsBeforePrewrites[k] = commit
	}

	return nil
}
