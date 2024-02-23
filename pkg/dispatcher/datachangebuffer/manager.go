package datachangebuffer


import (
	"github.com/zhangjinpeng87/tistream/pkg/storage"
)

// The data change buffer files are organized as follows:
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
// |____data_change_buffer/cluster-{2}/Tenant-{1}
// |  |____store-{id1}
// |  |____store-{id2}
// |  |____store-{id3}

// DataChangeBufferManager is the manager of the data change buffer.
type DataChangeBufferManager struct {
	// The backendStorage to store the data change files.
	backendStorage *storage.BackendStorage
}

func NewDataChangeBufferManager(backendStorage *storage.BackendStorage) *DataChangeBufferManager {
	return &DataChangeBufferManager{
		backendStorage: backendStorage,
	}
}
