package sorter

import (
	"context"

	"github.com/zhangjinpeng87/tistream/pkg/remotesorter/prewritebuffer"
	"github.com/zhangjinpeng87/tistream/pkg/remotesorter/sorterbuffer"

	"github.com/huandu/skiplist"
)

// TenantSorter is the sorter of the tenant.
type TenantSorter struct {
	// The tenant id.
	TenantID uint64

	prewriteBuffers skiplist.SkipList
	sorterBuffers skiplist.SkipList
}

// NewTenantSorter creates a new TenantSorter.
func NewTenantSorter(tenantID uint64) *TenantSorter {
	return &TenantSorter{
		TenantID: tenantID,
		prewriteBuffers: skiplist.New(skiplist.StringAsc),
		sorterBuffers: skiplist.New(skiplist.StringAsc),
	}
}


