package dispatcherworker

import (
	"context"
	"fmt"

	"github.com/zhangjinpeng87/tistream/pkg/codec"
	"golang.org/x/sync/errgroup"
)

type FileEvent codec.DataChangesFileDecoder

type DataRouter struct {
	eg  *errgroup.Group
	ctx context.Context

	tenantChannels map[uint64]chan *FileEvent
	tenantWorkers  map[uint64]*TenantSendingWorker

	// The meta server client. Used to:
	// 1) get the sorter address and schema-registry address.
	// 2) update the sorter address and schema-registry address when the tenant scheduled to another sorter or schema-registry.
	// 3) send the heartbeat to the meta server.
	metaServerCli *MetaServerClient
}

func NewDataRouter(ctx context.Context, eg *errgroup.Group) *DataRouter {
	return &DataRouter{
		eg:             eg,
		ctx:            ctx,
		tenantChannels: make(map[uint64]chan *FileEvent),
	}
}

func (dr *DataRouter) AddTenant(tenantID uint64) (chan<- *FileEvent, error) {
	_, ok := dr.tenantChannels[tenantID]
	if ok {
		return nil, fmt.Errorf("tenant %d already exists", tenantID)
	}

	ch := make(chan *FileEvent)
	dr.tenantChannels[tenantID] = ch

	tenantSendingWorker := NewTenantSendingWorker(tenantID, ch, dr.metaServerCli)
	dr.tenantWorkers[tenantID] = tenantSendingWorker
	tenantSendingWorker.Run(dr.ctx, dr.eg)

	return ch, nil
}

func (dr *DataRouter) RemoveTenant(tenantID uint64) {
	ch, ok := dr.tenantChannels[tenantID]
	if ok {
		close(ch)
		delete(dr.tenantChannels, tenantID)
		delete(dr.tenantWorkers, tenantID)
	}
}

func (dr *DataRouter) Run(ctx context.Context) {
	for tenantID, ch := range dr.tenantChannels {
		tenantID := tenantID
		ch := ch

		tenantSendingWorker := NewTenantSendingWorker(tenantID, ch, dr.metaServerCli)
		tenantSendingWorker.Run(ctx, dr.eg)
		dr.tenantWorkers[tenantID] = tenantSendingWorker
	}
}

func (dr *DataRouter) Close() {
	for _, ch := range dr.tenantChannels {
		close(ch)
	}
}
