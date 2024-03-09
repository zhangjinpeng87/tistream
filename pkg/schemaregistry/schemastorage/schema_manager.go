package schemastorage


type SchemaManager struct {
	// Tenant -> SchemaStorage
	storage map[uint64]SchemaStorage
}

func NewSchemaManager() *SchemaManager {
	return &SchemaManager{
		storage: make(map[uint64]SchemaStorage),
	}
}

func (m *SchemaManager) GetSchemaStorage(tenantID uint64) (SchemaStorage, bool) {
	s, ok := m.storage[tenantID]
	return s, ok
}

func (m *SchemaManager) SetSchemaStorage(tenantID uint64, storage SchemaStorage) {
	m.storage[tenantID] = storage
}