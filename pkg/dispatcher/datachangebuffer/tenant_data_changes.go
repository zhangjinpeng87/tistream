package datachangebuffer

import (
	"bytes"

	"github.com/zhangjinpeng87/tistream/pkg/storage"
)

const (
	schemaSnap = "schema_snap"
)

// The data change buffer files for a specified tenant are organized as follows:
// |____data_change_buffer/cluster-{1}/Tenant-{1}
// |  |____schema_snap
// |  |____store-{1}
// |  |  |____file-{ts}
// |  |  |____file-{ts}
// |  |____store-{2}
// |  |  |____file-{ts}
// |  |  |____file-{ts}
// |  |____store-{3}
// |     |____file-{ts}
//
// There are 2 types of files:
//  1. schema_snap: the schema snapshot file.
//  2. file-{ts}: the data change file. It contains the data change events for a specified
//     time range for a specified store.

// Tenant is the data change buffer for a tenant.
type TenantDataChanges struct {
	// The tenant id.
	tenantID int

	// Root directory of the data change buffer for the tenant.
	// In above file org it should be "{prefix}data_change_buffer/cluster-{1}/Tenant-{1}"
	rootDir string

	// Backend Storage
	backendStorage storage.BackendStorage
}

// NewTenantDataChanges creates a new TenantDataChanges.
func NewTenantDataChanges(tenantID int, rootDir string, backendStorage storage.BackendStorage) *TenantDataChanges {
	return &TenantDataChanges{
		tenantID:       tenantID,
		rootDir:        rootDir,
		backendStorage: backendStorage,
	}
}

// GetSchemaSnap returns the schema snapshot of this tenant.
func (t *TenantDataChanges) GetSchemaSnap() (*SchemaSnap, error) {
	filePath := t.rootDir + schemaSnap

	content, err := t.backendStorage.GetFile(filePath)
	if err != nil {
		return nil, err
	}

	schemaSnap := NewEmptySchemaSnap()
	reader := bytes.NewReader(content)
	err = schemaSnap.DecodeFrom(reader)
	if err != nil {
		return nil, err
	}

	return schemaSnap, nil
}

// ListStoreDirs returns the list of store directories.
func (t *TenantDataChanges) ListStoreDir() ([]string, error) {
	return t.backendStorage.ListSubDir(t.rootDir)
}
