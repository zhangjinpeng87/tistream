package datachangebuffer

import (
	"context"
	"fmt"

	"github.com/zhangjinpeng87/tistream/pkg/client"
	"github.com/zhangjinpeng87/tistream/pkg/codec"
	"golang.org/x/sync/errgroup"
)

type SenderManager struct {
	eg  *errgroup.Group
	ctx context.Context

	tenantChannels map[uint64]chan *codec.FileEvent
	tenantWorkers  map[uint64]*TenantDataSender

	// The meta server client. Used to:
	// 1) get the sorter address and schema-registry address.
	// 2) update the sorter address and schema-registry address when the tenant scheduled to another sorter or schema-registry.
	// 3) send the heartbeat to the meta server.
	metaServerCli *client.MetaServerClient
}

func NewSenderManager(ctx context.Context, eg *errgroup.Group) *SenderManager {
	return &SenderManager{
		eg:             eg,
		ctx:            ctx,
		tenantChannels: make(map[uint64]chan *codec.FileEvent),
	}
}

func (dr *SenderManager) AttachTenant(tenantID uint64) (chan<- *codec.FileEvent, error) {
	_, ok := dr.tenantChannels[tenantID]
	if ok {
		return nil, fmt.Errorf("tenant %d already exists", tenantID)
	}

	ch := make(chan *codec.FileEvent)
	dr.tenantChannels[tenantID] = ch

	tenantDataSender := NewTenantDataSender(tenantID, ch, dr.metaServerCli)
	dr.tenantWorkers[tenantID] = tenantDataSender
	tenantDataSender.Run(dr.ctx, dr.eg)

	return ch, nil
}

func (dr *SenderManager) DetachTenant(tenantID uint64) {
	ch, ok := dr.tenantChannels[tenantID]
	if ok {
		close(ch)
		delete(dr.tenantChannels, tenantID)
		delete(dr.tenantWorkers, tenantID)
	}
}

func (dr *SenderManager) Run(ctx context.Context) {
	for tenantID, ch := range dr.tenantChannels {
		tenantID := tenantID
		ch := ch

		tenantDataSender := NewTenantDataSender(tenantID, ch, dr.metaServerCli)
		tenantDataSender.Run(ctx, dr.eg)
		dr.tenantWorkers[tenantID] = tenantDataSender
	}
}

func (dr *SenderManager) Close() {
	for _, ch := range dr.tenantChannels {
		close(ch)
	}
}
