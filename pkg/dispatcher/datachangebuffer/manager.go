package datachangebuffer

import (
	"fmt"
	"strconv"
	"sync"
	"context"
	"golang.org/x/sync/errgroup"

	"github.com/zhangjinpeng87/tistream/pkg/storage"
	"github.com/zhangjinpeng87/tistream/pkg/utils"
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
	// Dispatcher Config
	cfg *utils.DispatcherConfig

	// The backendStorage to store the data change files.
	backendStorage storage.BackendStorage

	// RWMutex to protect the tenantChanges.
	mu sync.RWMutex

	// TenantID -> TenantDataChanges
	tenantChanges map[uint64]*TenantDataChanges
}

func NewDataChangeBufferManager(cfg *utils.DispatcherConfig, backendStorage storage.BackendStorage) *DataChangeBufferManager {
	return &DataChangeBufferManager{
		backendStorage: backendStorage,
		cfg:         cfg,
		tenantChanges: make(map[uint64]*TenantDataChanges),
	}
}

func (m *DataChangeBufferManager) AttachTenant(tenantID uint64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.tenantChanges[tenantID]; ok {
		return fmt.Errorf("tenant %d already attached", tenantID)
	}

	tenantChanges := NewTenantDataChanges(tenantID, strconv.FormatUint(tenantID, 10), m.backendStorage)
	m.tenantChanges[tenantID] = tenantChanges

	return nil
}

func (m *DataChangeBufferManager) DetachTenant(tenantID uint64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.tenantChanges[tenantID]; !ok {
		return fmt.Errorf("tenant %d not attached", tenantID)
	}

	delete(m.tenantChanges, tenantID)

	return nil
}

func (m *DataChangeBufferManager) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, tenantChanges := range m.tenantChanges {
		eg.Go(func() { 
			tenantChanges.Run(ctx, m.cfg.CheckStoreInterval, m.cfg.CheckStoreTimeout)
		})
	}

	return nil

}
