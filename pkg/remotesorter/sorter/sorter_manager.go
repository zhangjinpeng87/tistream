package sorter

import (
	"context"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/zhangjinpeng87/tistream/pkg/storage"
	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
)

// SorterManager is the manager of the all tenants' sorters for this remote-sorter instance.
type SorterManager struct {
	// The backendStorage to store the sorter files.
	backendStorage storage.ExternalStorage

	// The tenant sorters.
	// TenantID -> TenantSorter
	sorterMutex   sync.RWMutex
	tenantSorters map[uint64]*TenantSorter

	// The context to cancel the running goroutines.
	ctx context.Context
	eg  *errgroup.Group
}

func NewSorterManager(ctx context.Context, eg *errgroup.Group, es storage.ExternalStorage) *SorterManager {
	return &SorterManager{
		backendStorage: es,
		tenantSorters:  make(map[uint64]*TenantSorter),
		ctx:            ctx,
		eg:             eg,
	}
}

func (m *SorterManager) AddTaskRange(tenantID uint64, taskRange *pb.Task_Range) {
	m.sorterMutex.Lock()
	defer m.sorterMutex.Unlock()
	tenantSorter, ok := m.tenantSorters[tenantID]
	if !ok {
		tenantSorter = NewTenantSorter(tenantID, m.backendStorage)
		m.tenantSorters[tenantID] = tenantSorter
	}

	tenantSorter.AddNewRange(taskRange)
}

func (m *SorterManager) Run() error {
	// Todo: add more dedicated go routines to do different works.
	// Create a background routine to flush all tenant sorters periodically if needed.
	// Also, flush the committed data to the downstream periodically.
	m.eg.Go(func() error {
		ticker1 := time.NewTicker(5 * time.Second)
		defer ticker1.Stop()
		ticker2 := time.NewTicker(2 * time.Second)
		defer ticker2.Stop()
		for {
			select {
			case <-m.ctx.Done():
				return nil
			case <-ticker1.C:
				if err := m.SaveSnapshot(); err != nil {
					return err
				}
			case <-ticker2.C:
				if err := m.FlushCommittedData(); err != nil {
					return err
				}
			}
		}
	})

	return nil
}

func (m *SorterManager) AddEventBatch(tenantID uint64, taskRange *pb.Task_Range, eventBatch *pb.EventBatch) error {
	m.sorterMutex.RLock()
	ts, ok := m.tenantSorters[tenantID]
	m.sorterMutex.RUnlock()
	if !ok {
		return nil
	}
	return ts.AddEventBatch(tenantID, taskRange, eventBatch)
}

func (m *SorterManager) SaveSnapshot() error {
	// Todo: lock just get the sorters, release before call SaveSnapshot since it may take a long time.
	m.sorterMutex.RLock()
	defer m.sorterMutex.RUnlock()
	for _, ts := range m.tenantSorters {
		if err := ts.SaveSnapshot(); err != nil {
			return err
		}
	}
	return nil
}

func (m *SorterManager) FlushCommittedData() error {
	// Todo: lock just get the sorters, release before call FlushCommittedData since it may take a long time.
	m.sorterMutex.RLock()
	defer m.sorterMutex.RUnlock()
	for _, ts := range m.tenantSorters {
		if err := ts.FlushCommittedData(); err != nil {
			return err
		}
	}
	return nil
}
