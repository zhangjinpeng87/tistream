package datachangebuffer

import (
	"context"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/zhangjinpeng87/tistream/pkg/codec"
	"github.com/zhangjinpeng87/tistream/pkg/storage"
	"github.com/zhangjinpeng87/tistream/pkg/utils"
	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
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

// ReaderManager is the manager of the data change buffer.
type ReaderManager struct {
	// Dispatcher Config
	cfg *utils.DispatcherConfig

	// The backendStorage to store the data change files.
	backendStorage storage.ExternalStorage

	// RWMutex to protect the tenantReaders.
	mu sync.RWMutex

	// TenantID -> TenantReader
	tenantReaders map[uint64]*TenantReader
	// TenantID -> channel to notify the tenant has new data change files.
	tenantChannels map[uint64]chan struct{}

	// Channel that receives current status from all tenants.
	statusCh chan *pb.TenantSubStats

	statsMu sync.RWMutex
	// TenantId-Range -> TenantSubStats
	dispatcherStats map[string]*pb.TenantSubStats
	totalThroughput uint64

	// The context to cancel the running goroutines.
	ctx context.Context
	eg  *errgroup.Group
}

func NewReaderManager(ctx context.Context, eg *errgroup.Group, cfg *utils.DispatcherConfig, backendStorage storage.ExternalStorage) *ReaderManager {
	return &ReaderManager{
		backendStorage: backendStorage,
		cfg:            cfg,
		tenantReaders:  make(map[uint64]*TenantReader),
		statusCh:       make(chan *pb.TenantSubStats, 128),
		eg:             eg,
		ctx:            ctx,
	}
}

func (m *ReaderManager) AttachTenant(tenantID uint64, fileSender chan<- *codec.FileEvent) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.tenantReaders[tenantID]; ok {
		return fmt.Errorf("tenant %d already attached", tenantID)
	}

	tenantRootDir := fmt.Sprintf("%s/Tenant-%d", m.cfg.Storage.Prefix, tenantID)

	tenantReader := NewTenantReader(tenantID, tenantRootDir, fileSender, m.backendStorage)
	m.tenantReaders[tenantID] = tenantReader

	// Create channel to notify the tenant has new data change files.
	c := make(chan struct{}, 8)
	m.tenantChannels[tenantID] = c

	return nil
}

func (m *ReaderManager) DetachTenant(tenantID uint64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.tenantReaders[tenantID]; !ok {
		return fmt.Errorf("tenant %d not attached", tenantID)
	}

	delete(m.tenantReaders, tenantID)
	close(m.tenantChannels[tenantID])
	delete(m.tenantChannels, tenantID)

	return nil
}

func (m *ReaderManager) Run(ctx context.Context) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Start the tenant data changes.
	for _, tenantReader := range m.tenantReaders {
		tc := tenantReader
		recv, ok := m.tenantChannels[tenantReader.tenantID]
		if !ok {
			errMsg := fmt.Sprintf("tenant %d no notification channel", tenantReader.tenantID)
			panic(errMsg)
		}
		m.eg.Go(func() error {
			return tc.Run(ctx, m.cfg.CheckStoreInterval, m.cfg.CheckFileInterval, recv, m.statusCh)
		})
	}

	// Start the status collector.
	m.eg.Go(func() error {
		return m.collectStats(ctx)
	})
}

func (m *ReaderManager) TenantHasNewChanges(tenantID uint64) {
	// Find the notify channel according to the tenantID.
	c, ok := m.tenantChannels[tenantID]
	if !ok {
		return
	}

	// Notify the tenant has new changes without block.
	select {
	case c <- struct{}{}:
	default:
	}
}

func (m *ReaderManager) collectStats(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case status := <-m.statusCh:
			m.statsMu.Lock()

			k := fmt.Sprintf("%d-%v-%v", status.TenantId, status.Range.Start, status.Range.End)
			_, ok := m.dispatcherStats[k]
			if !ok {
				m.dispatcherStats[k] = &pb.TenantSubStats{
					TenantId: status.TenantId,
					Range:    status.Range,
				}
			}

			m.dispatcherStats[k].Throughput += status.Throughput
			m.totalThroughput += status.Throughput

			m.statsMu.Unlock()
		}
	}
}

func (m *ReaderManager) FetchStatsAndReset() (uint64, []*pb.TenantSubStats) {
	m.statsMu.Lock()
	defer m.statsMu.Unlock()

	// Copy the stats.
	res := make([]*pb.TenantSubStats, 0, len(m.dispatcherStats))
	for _, v := range m.dispatcherStats {
		res = append(res, v)
	}
	totalThroughput := m.totalThroughput

	// Reset the stats.
	m.totalThroughput = 0
	m.dispatcherStats = make(map[string]*pb.TenantSubStats)

	return totalThroughput, res
}
