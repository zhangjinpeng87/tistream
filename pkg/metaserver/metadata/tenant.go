package metadata

// Tenant is the tenant basic information like tenant id.

type Tenant struct {
	TenantID string

	// ClusterID is the cluster id of the tenant.
	ClusterID string
}