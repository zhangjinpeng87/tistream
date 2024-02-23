package datachangebuffer

import (
	"encoding/binary"
	"encoding/json"
	"io"

	"github.com/pingcap/tidb/pkg/parser/model"
	"github.com/zhangjinpeng87/tistream/pkg/utils"
)

const (
	SchemaFileMagicHeader = string("TiStreamSchemaSnap")
	SchemaFileVersion     = 1
)

type SchemaSnap struct {
	// All the schema and table info.
	dbs    []*model.DBInfo
	tables []*model.TableInfo

	// The timestamp of the schema snapshot.
	currentTs uint64
}

// NewSchemaSnap creates a new SchemaSnap.
func NewSchemaSnap() *SchemaSnap {
	return &SchemaSnap{
		dbs:       make([]*model.DBInfo, 0),
		tables:    make([]*model.TableInfo, 0),
		currentTs: 0,
	}
}

func (s *SchemaSnap) EncodeTo(writer io.Writer) error {
	// Encode the header.
	writer.Write([]byte(SchemaFileMagicHeader))
	// Encode the version.
	binary.Write(writer, binary.LittleEndian, uint32(SchemaFileVersion))
	// Encode the timestamp.
	binary.Write(writer, binary.LittleEndian, s.currentTs)
	// Encode the number of dbs.
	binary.Write(writer, binary.LittleEndian, uint64(len(s.dbs)))
	// Encode the number of tables.
	binary.Write(writer, binary.LittleEndian, uint64(len(s.tables)))

	// Encode the dbs.
	for _, db := range s.dbs {
		marshal, err := json.Marshal(db)
		if err != nil {
			return err
		}
		writer.Write(marshal)
	}

	// Encode the tables.
	for _, table := range s.tables {
		marshal, err := json.Marshal(table)
		if err != nil {
			return err
		}
		writer.Write(marshal)
	}

	return nil
}

func (s *SchemaSnap) DecodeFrom(reader io.Reader) error {
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

	// Read the timestamp.
	err = binary.Read(reader, binary.LittleEndian, &s.currentTs)
	if err != nil {
		return err
	}

	// Read the number of dbs.
	var dbNum uint64
	err = binary.Read(reader, binary.LittleEndian, &dbNum)
	if err != nil {
		return err
	}

	// Read the number of tables.
	var tableNum uint64
	err = binary.Read(reader, binary.LittleEndian, &tableNum)
	if err != nil {
		return err
	}

	// Read the dbs.
	for i := 0; i < int(dbNum); i++ {
		db := &model.DBInfo{}
		err = json.NewDecoder(reader).Decode(db)
		if err != nil {
			return err
		}
		s.dbs = append(s.dbs, db)
	}

	// Read the tables.
	for i := 0; i < int(tableNum); i++ {
		table := &model.TableInfo{}
		err = json.NewDecoder(reader).Decode(table)
		if err != nil {
			return err
		}
		s.tables = append(s.tables, table)
	}

	return nil
}

// Compare the schema snap. Used for test.
func (s *SchemaSnap) equalTo(s2 *SchemaSnap) bool {
	if s.currentTs != s2.currentTs {
		return false
	}

	if len(s.dbs) != len(s2.dbs) {
		return false
	}

	for i := 0; i < len(s.dbs); i++ {
		if s.dbs[i].ID != s2.dbs[i].ID {
			return false
		}
		if s.dbs[i].Name.L != s2.dbs[i].Name.L {
			return false
		}
		if s.dbs[i].Charset != s2.dbs[i].Charset {
			return false
		}
		if s.dbs[i].Collate != s2.dbs[i].Collate {
			return false
		}
	}

	if len(s.tables) != len(s2.tables) {
		panic("tables length not equal")
	}

	for i := 0; i < len(s.tables); i++ {
		if s.tables[i].ID != s2.tables[i].ID {
			return false
		}
		if s.tables[i].Name.L != s2.tables[i].Name.L {
			return false
		}
		if s.tables[i].Charset != s2.tables[i].Charset {
			return false
		}
		if s.tables[i].Collate != s2.tables[i].Collate {
			return false
		}
		if len(s.tables[i].Columns) != len(s2.tables[i].Columns) {
			return false
		}
		for j := 0; j < len(s.tables[i].Columns); j++ {
			if s.tables[i].Columns[j].ID != s2.tables[i].Columns[j].ID {
				return false
			}
			if s.tables[i].Columns[j].Name.L != s2.tables[i].Columns[j].Name.L {
				return false
			}
		}
	}

	return true
}
