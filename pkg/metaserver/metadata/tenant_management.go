package metadata

import "sync"

// Tenant is the tenant information.
type Tenant struct {
	// ID is the tenant ID.
	ID uint32 `json:"id"`
	// ClusterID is the cluster ID.
	ClusterID uint32 `json:"cluster_id"`
	// SnapshotAddr is the snapshot address.
	SnapshotAddr string `json:"snapshot_addr"`
	// DataChangeAddr is the data change address.
	DataChangeAddr string `json:"data_change_addr"`
	// KMS is the key management service.
	KMS string `json:"kms"`
}

type TenantManagement struct {
	sync.RWMutex

	// TenantId -> Tenant
	tenants map[uint32]*Tenant
}

// NewTenantManagement creates a new TenantManagement.
func NewTenantManagement() *TenantManagement {
	return &TenantManagement{
		tenants: make(map[uint32]*Tenant),
	}
}

// AddTenant adds a tenant.
func (tm *TenantManagement) AddTenant(tenant *Tenant) {
	tm.Lock()
	defer tm.Unlock()

	tm.tenants[tenant.ID] = tenant
}

// GetTenant gets a tenant.
func (tm *TenantManagement) GetTenant(tenantID uint32) *Tenant {
	tm.RLock()
	defer tm.RUnlock()

	return tm.tenants[tenantID]
}

// DeleteTenant deletes a tenant.
func (tm *TenantManagement) DeleteTenant(tenantID uint32) {
	tm.Lock()
	defer tm.Unlock()

	delete(tm.tenants, tenantID)
}