package datachangebuffer

import (
	"os"
	"testing"

	"github.com/pingcap/tidb/pkg/parser/model"
)

func TestSchemaSnap_EncodeDecode(t *testing.T) {
	dbInfo := &model.DBInfo{
		ID:      1,
		Name:    model.NewCIStr("test"),
		Charset: "utf8",
		Collate: "utf8_bin",
	}

	column1 := &model.ColumnInfo{
		ID:   1,
		Name: model.NewCIStr("c1"),
	}

	tableInfo := &model.TableInfo{
		ID:      1,
		Name:    model.NewCIStr("t1"),
		Charset: "utf8",
		Collate: "utf8_bin",
		Columns: []*model.ColumnInfo{column1},
	}

	schemaSnap := NewEmptySchemaSnap()
	schemaSnap.dbs = append(schemaSnap.dbs, dbInfo)
	schemaSnap.tables = append(schemaSnap.tables, tableInfo)
	schemaSnap.currentTs = 123

	// Open a temp file.
	file, err := os.CreateTemp("", "test_schema_snap")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	// Encode the schema snap to the file.
	if err := schemaSnap.EncodeTo(file); err != nil {
		t.Fatal(err)
	}

	// Read the file and decode the schema snap.
	schemaSnap2 := NewEmptySchemaSnap()
	file.Seek(0, 0)
	err = schemaSnap2.DecodeFrom(file)
	if err != nil {
		t.Fatal(err)
	}

	// Compare the schema snap.
	if schemaSnap2.currentTs != schemaSnap.currentTs {
		t.Fatalf("currentTs: %d != %d", schemaSnap2.currentTs, schemaSnap.currentTs)
	}

	if !schemaSnap2.equalTo(schemaSnap) {
		t.Fatalf("dbs length: %d != %d", len(schemaSnap2.dbs), len(schemaSnap.dbs))
	}
}
