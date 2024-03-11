package codec

import (
	"encoding/binary"
	"io"

	"github.com/zhangjinpeng87/tistream/pkg/utils"
	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
	"google.golang.org/protobuf/proto"
)

const (
	SchemaFileMagicHeader = string("TiStreamSchemaSnap")
	SchemaFileVersion     = 1
)

type SchemaSnapFile struct {
	// All the schema and table info.
	snap *pb.SchemaSnapshot
}

// NewSchemaSnap creates a new SchemaSnap.
func NewEmptySchemaSnap() *SchemaSnapFile {
	return &SchemaSnapFile{}
}

func (s *SchemaSnapFile) EncodeTo(writer io.Writer) error {
	// Encode the header.
	_, err := writer.Write([]byte(SchemaFileMagicHeader))
	if err != nil {
		return err
	}

	// Encode the version.
	err = binary.Write(writer, binary.LittleEndian, uint32(SchemaFileVersion))
	if err != nil {
		return err
	}

	data, err := proto.Marshal(s.snap)
	if err != nil {
		return err
	}
	err = binary.Write(writer, binary.LittleEndian, uint32(len(data)))
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (s *SchemaSnapFile) DecodeFrom(reader io.Reader) error {
	// Read the header.
	header := make([]byte, len(SchemaFileMagicHeader))
	n, err := reader.Read(header)
	if err != nil {
		return err
	}
	if n != len(SchemaFileMagicHeader) {
		return io.ErrUnexpectedEOF
	}
	if string(header) != SchemaFileMagicHeader {
		return utils.ErrInvalidSchemaSnapFile
	}

	// Read the version.
	// Todo: check the version, and handle the different version.
	var version uint32
	err = binary.Read(reader, binary.LittleEndian, &version)
	if err != nil {
		return err
	}

	// Read the data.
	var dataLen uint32
	err = binary.Read(reader, binary.LittleEndian, &dataLen)
	if err != nil {
		return err
	}
	data := make([]byte, dataLen)
	_, err = io.ReadFull(reader, data)
	if err != nil {
		return err
	}

	// Decode the data.
	s.snap = &pb.SchemaSnapshot{}
	err = proto.Unmarshal(data, s.snap)
	if err != nil {
		return err
	}

	return nil
}
