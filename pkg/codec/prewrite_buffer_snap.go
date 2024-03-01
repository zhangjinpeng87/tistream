package codec

import (
	"encoding/binary"
	"io"

	"github.com/zhangjinpeng87/tistream/pkg/utils"
	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
	"google.golang.org/protobuf/proto"
)

const (
	// The magic number of the prewrite buffer file.
	PrewriteBufferSnapFileMagicNumber = uint32(0x21647911)
	PrewriteBufferSnapFileVersion     = uint32(1)
)

type PrewriteBufferSnapEncoder struct {
	// The tenant id.
	TenantID uint64

	// range of the buffer.
	Range *pb.Task_Range
}

func NewPrewriteBufferSnapEncoder(tenantID uint64, range_ *pb.Task_Range) *PrewriteBufferSnapEncoder {
	return &PrewriteBufferSnapEncoder{
		TenantID: tenantID,
		Range:    range_,
	}
}

func (p *PrewriteBufferSnapEncoder) Encode(w io.Writer, prewrites map[string]*pb.EventRow, commits map[string]*pb.EventRow) error {
	// Write header
	if err := binary.Write(w, binary.LittleEndian, PrewriteBufferSnapFileMagicNumber); err != nil {
		return err
	}

	// Write version.
	if err := binary.Write(w, binary.LittleEndian, PrewriteBufferSnapFileVersion); err != nil {
		return err
	}

	// Write tenant id.
	if err := binary.Write(w, binary.LittleEndian, p.TenantID); err != nil {
		return err
	}

	// Write the range.
	if err := binary.Write(w, binary.LittleEndian, uint32(len(p.Range.Start))); err != nil {
		return err
	}
	if _, err := w.Write(p.Range.Start); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, uint32(len(p.Range.End))); err != nil {
		return err
	}
	if _, err := w.Write(p.Range.End); err != nil {
		return err
	}

	// Write prewrites.
	if err := binary.Write(w, binary.LittleEndian, uint32(len(prewrites))); err != nil {
		return err
	}
	for _, prewrite := range prewrites {
		data, err := proto.Marshal(prewrite)
		if err != nil {
			return err
		}
		if err := binary.Write(w, binary.LittleEndian, uint32(len(data))); err != nil {
			return err
		}
		if _, err := w.Write(data); err != nil {
			return err
		}
	}

	// Write commits.
	if err := binary.Write(w, binary.LittleEndian, uint32(len(commits))); err != nil {
		return err
	}
	for _, commit := range commits {
		data, err := proto.Marshal(commit)
		if err != nil {
			return err
		}
		if err := binary.Write(w, binary.LittleEndian, uint32(len(data))); err != nil {
			return err
		}
		if _, err := w.Write(data); err != nil {
			return err
		}
	}

	return nil
}

type PrewriteBufferSnapDecoder struct {
	// The tenant id.
	TenantID uint64

	// range of the buffer.
	Range *pb.Task_Range
}

func NewPrewriteBufferSnapDecoder(tenantID uint64, range_ *pb.Task_Range) *PrewriteBufferSnapDecoder {
	return &PrewriteBufferSnapDecoder{
		TenantID: tenantID,
		Range:    range_,
	}
}

func (p *PrewriteBufferSnapDecoder) Decode(r io.Reader) ([]*pb.EventRow, []*pb.EventRow, error) {
	// Read header
	version, err := p.decodeHeader(r)
	if err != nil {
		return nil, nil, err
	}

	// Todo: decode the file according to the version.
	switch version {
	case 1:
		return p.decodeBodyV1(r)
	default:
		return nil, nil, utils.ErrInvalidPrewriteBufferSnapFileVersion
	}
}

func (p *PrewriteBufferSnapDecoder) decodeHeader(r io.Reader) (uint32, error) {
	// Read header
	var magicNumber uint32
	if err := binary.Read(r, binary.LittleEndian, &magicNumber); err != nil {
		return 0, err
	}
	if magicNumber != PrewriteBufferSnapFileMagicNumber {
		return 0, utils.ErrInvalidPrewriteBufferSnapFile
	}

	// Read version.
	var version uint32
	if err := binary.Read(r, binary.LittleEndian, &version); err != nil {
		return 0, err
	}

	// Read tenant id.
	var tenantID uint64
	if err := binary.Read(r, binary.LittleEndian, &tenantID); err != nil {
		return 0, err
	}
	if tenantID != p.TenantID {
		return 0, utils.ErrUnmatchedTenantID
	}

	// Read the range.
	var startLen uint32
	if err := binary.Read(r, binary.LittleEndian, &startLen); err != nil {
		return 0, err
	}
	var start []byte
	if startLen > 0 {
		start := make([]byte, startLen)
		if _, err := io.ReadFull(r, start); err != nil {
			return 0, err
		}
	}
	var endLen uint32
	if err := binary.Read(r, binary.LittleEndian, &endLen); err != nil {
		return 0, err
	}
	var end []byte
	if endLen > 0 {
		end := make([]byte, endLen)
		if _, err := io.ReadFull(r, end); err != nil {
			return 0, err
		}
	}
	if string(start) != string(p.Range.Start) || string(end) != string(p.Range.End) {
		return 0, utils.ErrUnmatchedRange
	}

	return version, nil
}

func (p *PrewriteBufferSnapDecoder) decodeBodyV1(r io.Reader) (prewrites []*pb.EventRow, commits []*pb.EventRow, err error) {
	var prewriteCount uint32
	if err := binary.Read(r, binary.LittleEndian, &prewriteCount); err != nil {
		return nil, nil, err
	}
	prewrites = make([]*pb.EventRow, 0, prewriteCount)
	for i := uint32(0); i < prewriteCount; i++ {
		var eventSize uint32
		if err := binary.Read(r, binary.LittleEndian, &eventSize); err != nil {
			return nil, nil, err
		}
		data := make([]byte, eventSize)
		if _, err := io.ReadFull(r, data); err != nil {
			return nil, nil, err
		}
		event := &pb.EventRow{}
		if err := proto.Unmarshal(data, event); err != nil {
			return nil, nil, err
		}
		prewrites = append(prewrites, event)
	}

	var commitCount uint32
	if err := binary.Read(r, binary.LittleEndian, &commitCount); err != nil {
		return nil, nil, err
	}
	commits = make([]*pb.EventRow, 0, commitCount)
	for i := uint32(0); i < commitCount; i++ {
		var eventSize uint32
		if err := binary.Read(r, binary.LittleEndian, &eventSize); err != nil {
			return nil, nil, err
		}
		data := make([]byte, eventSize)
		if _, err := io.ReadFull(r, data); err != nil {
			return nil, nil, err
		}
		event := &pb.EventRow{}
		if err := proto.Unmarshal(data, event); err != nil {
			return nil, nil, err
		}
		commits = append(commits, event)
	}

	return prewrites, commits, nil
}
