package datapool

import (
	"fmt"

	"github.com/zhangjinpeng87/tistream/pkg/storage"
	"github.com/zhangjinpeng87/tistream/pkg/utils"
)

// CommittedDataPool is the committed data pool for multiple tenants.
type CommittedDataPool struct {
	cfg *utils.StorageConfig

	// The data pool.
	dataPool map[uint64]*TenantDataPool

	// backend storage.
	backendStorage storage.ExternalStorage
}

// NewCommittedDataPool creates a new data management.
func NewCommittedDataPool(backendStorage storage.ExternalStorage) *CommittedDataPool {
	return &CommittedDataPool{
		dataPool:       make(map[uint64]*TenantDataPool),
		backendStorage: backendStorage,
	}
}

func (dm *CommittedDataPool) AddTenant(tenantID uint64) error {
	if _, ok := dm.dataPool[tenantID]; ok {
		return fmt.Errorf("tenant %d already exists", tenantID)
	}

	tenantRootDir := fmt.Sprintf("%s/Tenant-%d", dm.cfg.Prefix, tenantID)
	tenantDataPool := NewTenantDataPool(tenantID, tenantRootDir, dm.backendStorage)
	tenantDataPool.Init()

	return nil
}
