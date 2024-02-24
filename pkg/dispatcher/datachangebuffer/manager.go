package datachangebuffer

import (
	"context"
	"fmt"
	"strconv"
	"sync"

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
	// TenantID -> channel to notify the tenant has new data change files.
	tenantChannels map[uint64]chan struct{}

	// The context to cancel the running goroutines.
	ctx context.Context
	eg  *errgroup.Group
}

func NewDataChangeBufferManager(ctx context.Context, cfg *utils.DispatcherConfig, backendStorage storage.BackendStorage) *DataChangeBufferManager {
	eg, ctx := errgroup.WithContext(ctx)

	return &DataChangeBufferManager{
		ctx:            ctx,
		eg:             eg,
		backendStorage: backendStorage,
		cfg:            cfg,
		tenantChanges:  make(map[uint64]*TenantDataChanges),
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

	// Create channel to notify the tenant has new data change files.
	c := make(chan struct{}, 8)
	m.tenantChannels[tenantID] = c

	return nil
}

func (m *DataChangeBufferManager) DetachTenant(tenantID uint64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.tenantChanges[tenantID]; !ok {
		return fmt.Errorf("tenant %d not attached", tenantID)
	}

	delete(m.tenantChanges, tenantID)
	close(m.tenantChannels[tenantID])
	delete(m.tenantChannels, tenantID)

	return nil
}

func (m *DataChangeBufferManager) Run(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, tenantChanges := range m.tenantChanges {
		tc := tenantChanges
		c, ok := m.tenantChannels[tenantChanges.tenantID]
		if !ok {
			return fmt.Errorf("tenant %d no notification channel", tenantChanges.tenantID)
		}
		m.eg.Go(func() error {
			return tc.Run(ctx, m.cfg.CheckStoreInterval, m.cfg.CheckFileInterval, c)
		})
	}

	return nil
}
