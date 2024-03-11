package schemastorage

import (
	"fmt"

	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
)

type SchemaManager struct {
	// Tenant -> SchemaStorage
	storage map[uint64]*SchemaStorage
}

func NewSchemaManager() *SchemaManager {
	return &SchemaManager{
		storage: make(map[uint64]*SchemaStorage),
	}
}

func (m *SchemaManager) GetSchemaStorage(tenantID uint64) (*SchemaStorage, bool) {
	s, ok := m.storage[tenantID]
	return s, ok
}

func (m *SchemaManager) SetSchemaStorage(tenantID uint64, storage *SchemaStorage) {
	m.storage[tenantID] = storage
}

// Add a schema storage for a tenant.
func (m *SchemaManager) AddSchemaStorage(storage *SchemaStorage) {
	m.storage[storage.tenantID] = storage
}

func (m *SchemaManager) AddTableDdl(tenantId uint64, t *pb.Table) error {
	storage, ok := m.GetSchemaStorage(tenantId)
	if !ok {
		return fmt.Errorf("tenant %d not found", tenantId)
	}

	return storage.AddTableRecord(t.Schema.Name, t)
}

func (m *SchemaManager) AddSchemaDdl(tenantId uint64, s *pb.Schema) error {
	storage, ok := m.GetSchemaStorage(tenantId)
	if !ok {
		return fmt.Errorf("tenant %d not found", tenantId) 
	}

	return storage.AddSchemaRecord(s)
}
