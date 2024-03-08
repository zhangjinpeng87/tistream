package schemastorage

import (
	"fmt"

	"github.com/huandu/skiplist"
	"github.com/zhangjinpeng87/tistream/pkg/storage"

	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
)

const (
	DropDatabase   uint32 = 0
	CreateDatabase uint32 = 1
)

type SchemaSnapshot interface {
	// ListDatabases lists the databases in the snapshot.
	ListDatabases() []string

	// GetTable gets the schema of the specified table.
	GetTable(db string, tableID, ts uint64) (*pb.Table, error)

	// ListTables lists the tables in the specified database.
	ListTables(db string, ts uint64) ([]*pb.Table, error)
}

// SchemaStorage is a mvcc schema storage for a specified tenant.
// It is constructed from the SchemaSnap file and SchemaChange files from the external storage.
type SchemaStorage struct {
	// The tenant id.
	tenantID uint64

	// The backendStorage to store the schema files.
	backendStorage storage.ExternalStorage

	// Use skiplist to store table mvcc for last N days.
	// db.table.ts -> pb.Table
	tables *skiplist.SkipList
}

// NewSchemaStorage creates a new SchemaStorage.
func NewSchemaStorage(tenantID uint64, backendStorage storage.ExternalStorage) *SchemaStorage {
	return &SchemaStorage{
		tenantID:       tenantID,
		backendStorage: backendStorage,
		tables:         skiplist.New(skiplist.StringDesc), // We use desc order to get the latest table schema.
	}
}

// AddTable puts the schema of the specified table.
func (s *SchemaStorage) AddTableRecord(db string, table *pb.Table, ts uint64) {
	// Put the table into the skiplist.
	// We use the reverse ts as the score, so we can get the largest version of  a specifed ts can see.
	k := fmt.Sprintf("%s.%d.%020d", db, table.Id, ts)
	s.tables.Set(k, table)
}

// GetTable gets the schema of the specified table.
func (s *SchemaStorage) GetTable(db string, tableID, ts uint64) (*pb.Table, error) {
	// Get the table from the skiplist.
	k := fmt.Sprintf("t.%s.%d.%020d", db, tableID, ts)
	ele := s.tables.Find(k)
	if ele != nil {
		t := ele.Value.(*pb.Table)
		if t.Id == tableID {
			return t, nil
		}
	}

	return nil, nil
}

func (s *SchemaStorage) AddDBRecord(db *pb.Schema, ts uint64) {
	k := fmt.Sprintf("d.%s.%020d", db, ts)
	s.tables.Set(k, db)
}

func (s *SchemaStorage) ListDatabases(ts uint64) []string {
	k := "d."
	ele := s.tables.Find(k)
	var currentDb string
	res := make([]string, 0)
	for ele != nil {
		db := ele.Value.(*pb.Schema)
		// Check first visible version of this db with the specified ts.
		if currentDb != db.Name && db.Ts <= ts {
			if db.Op == pb.Schema_CREATE {
				res = append(res, db.Name)
			}
			currentDb = db.Name
		}

		ele = ele.Next()
	}

	return res
}

// GetSnapshot gets the snapshot of the schema storage.
func (s *SchemaStorage) GetSnapshot() SchemaSnapshot {
	return nil
}
