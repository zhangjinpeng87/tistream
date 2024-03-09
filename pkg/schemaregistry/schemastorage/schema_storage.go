package schemastorage

import (
	"bytes"
	"fmt"

	"github.com/huandu/skiplist"
	"github.com/zhangjinpeng87/tistream/pkg/storage"
	"github.com/zhangjinpeng87/tistream/pkg/utils"

	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
)

type SchemaSnapshot interface {
	// ListDatabases lists the databases in the snapshot.
	ListDatabases() []string

	// GetTable gets the schema of the specified table.
	GetTable(db string, tableID, ts uint64) (*pb.Table, error)

	// ListTables lists the tables in the specified database.
	ListTables(db string, ts uint64) ([]*pb.Table, error)
}

type schemaSnapshot struct {
	// The tenant id.
	tenantID uint64

	// Timestamp of the snapshot.
	ts uint64

	// The inner storage.
	innerStorage *SchemaStorage
}

func (s *schemaSnapshot) ListDatabases() []string {
	return s.innerStorage.ListDatabases(s.ts)
}

func (s *schemaSnapshot) GetTable(db string, tableID, ts uint64) (*pb.Table, error) {
	return s.innerStorage.GetTable(db, tableID, ts)
}

func (s *schemaSnapshot) ListTables(db string, ts uint64) ([]*pb.Table, error) {
	return nil, nil
}

func NewSchemaSnapshot(tenantID, ts uint64, innerStorage *SchemaStorage) SchemaSnapshot {
	return &schemaSnapshot{
		tenantID:     tenantID,
		ts:           ts,
		innerStorage: innerStorage,
	}
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
		tables:         skiplist.New(skiplist.ByteAsc),
	}
}

// AddTable puts the schema of the specified table.
func (s *SchemaStorage) AddTableRecord(db string, table *pb.Table, ts uint64) {
	// Put the table into the skiplist.
	// We use the reverse ts as the score, so we can get the largest version of  a specifed ts can see.
	k := TableKey(db, table.Id, ts)
	s.tables.Set(k, table)
}

// GetTable gets the schema of the specified table.
func (s *SchemaStorage) GetTable(db string, tableID, ts uint64) (*pb.Table, error) {
	// Get the table from the skiplist.
	k := TableKey(db, tableID, ts)
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
	k := DatabaseKey(db.Name, ts)
	s.tables.Set(k, db)
}

func (s *SchemaStorage) ListDatabases(ts uint64) []string {
	k := []byte("d.")
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

func (s *SchemaStorage) CompactTo(ts uint64) {
	ele := s.tables.Front()
	gcKeys := make([][]byte, 0)
	var curUserKey []byte
	var skippedLastVersion bool
	for ele != nil {
		key := ele.Key().([]byte)
		userKey := key[:len(key)-8]
		keyTs := utils.ReverseBytesToTs(key[len(key)-8:])

		// This is a new user key.
		if !bytes.Equal(userKey, curUserKey) {
			skippedLastVersion = false
			curUserKey = userKey
			continue
		}

		if keyTs <= ts {
			if skippedLastVersion {
				gcKeys = append(gcKeys, key)
			} else {
				skippedLastVersion = true
			}
		}

		ele = ele.Next()
	}

	for _, k := range gcKeys {
		s.tables.Remove(k)
	}
}

// GetSnapshot gets the snapshot of the schema storage.
func (s *SchemaStorage) GetSnapshot() SchemaSnapshot {
	return nil
}

func TableKey(db string, tableID, ts uint64) []byte {
	kPrefix := fmt.Sprintf("t.%s.%d", db, tableID)
	return append([]byte(kPrefix), utils.TsToReverseBytes(ts)...)
}

func DatabaseKey(db string, ts uint64) []byte {
	kPrefix := fmt.Sprintf("d.%s", db)
	return append([]byte(kPrefix), utils.TsToReverseBytes(ts)...)
}
