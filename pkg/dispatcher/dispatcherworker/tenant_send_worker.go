package dispatcherworker

import (
	"context"
	"fmt"

	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
	"golang.org/x/sync/errgroup"
)

type TenantSendingWorker struct {
	tenantId uint64
	ch       <-chan *FileEvent

	// Meta server client. Used to get the sorter address and schema-registry address.
	// And update the sorter address and schema-registry address when the tenant scheduled to another sorter or schema-registry.
	metaSrvCli *MetaServerClient

	// The remote sorter client.
	sorterCli *SorterClient

	// The schema registry client.
	schemeCli *SchemaRegistryClient
}

func NewTenantSendingWorker(tenantId uint64, ch <-chan *FileEvent, metaSrvCli *MetaServerClient) *TenantSendingWorker {
	return &TenantSendingWorker{
		tenantId:   tenantId,
		ch:         ch,
		metaSrvCli: metaSrvCli,
	}
}

func (tsw *TenantSendingWorker) Run(ctx context.Context, eg *errgroup.Group) {

	// Get the sorter address and schema-registry address.
	sorterAddr, err := tsw.metaSrvCli.FetchSorterAddr(ctx, &pb.FetchRangeSorterAddrReq{TenantId: tsw.tenantId})
	if err != nil {
		panic(fmt.Sprintf("fetch sorter addr failed, tenantID: %d, err: %v", tsw.tenantId, err))
	}
	sorterCli, err := NewSorterClient(sorterAddr)
	if err != nil {
		panic(fmt.Sprintf("create sorter client failed, tenantID: %d, err: %v", tsw.tenantId, err))
	}

	schemaRegistryAddr, err := tsw.metaSrvCli.FetchSchemaRegistryAddr(ctx, &pb.FetchSchemaRegistryAddrReq{TenantId: tsw.tenantId})
	if err != nil {
		panic(fmt.Sprintf("fetch schema registry addr failed, tenantID: %d, err: %v", tsw.tenantId, err))
	}
	schemaCli, err := NewSchemaRegistryClient(schemaRegistryAddr)
	if err != nil {
		panic(fmt.Sprintf("create schema registry client failed, tenantID: %d, err: %v", tsw.tenantId, err))
	}

	tsw.sorterCli = sorterCli
	tsw.schemeCli = schemaCli

	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return nil
			case fileEvent := <-tsw.ch:
				if fileEvent.TenantID != tsw.tenantId {
					panic(fmt.Sprintf("tenant id not match, expect %d, got %d", tsw.tenantId, fileEvent.TenantID))
				}

				// Send row event and watermark to the remote-sorter.
				req := &pb.RangeChangesReq{
					TenantId:   tsw.tenantId,
					Rows:       fileEvent.EventRows,
					Watermarks: fileEvent.EventWatermarks,
				}
				_, err := tsw.sorterCli.SendDataChanges(ctx, req)
				if err != nil {
					// Todo:
					// 1) handle the error, and retry.
					// 2) update the sorter address when the tenant scheduled to another sorter.
				}

				// Send the ddl event to the schema-registry.
				for _, ddl := range fileEvent.DdlChanges {
					req := &pb.RegisterDDLReq{
						TenantId: tsw.tenantId,
					}
					if ddl.Schema != nil {
						req.Ddl = &pb.RegisterDDLReq_SchemaDdl{SchemaDdl: ddl.Schema}
					} else if ddl.Table != nil {
						req.Ddl = &pb.RegisterDDLReq_TableDdl{TableDdl: ddl.Table}
					}

					_, err := tsw.schemeCli.SendDDLChange(ctx, req)
					if err != nil {
						// todo:
						// 1) handle the error, and retry.
						// 2) update the schema-registry address when the tenant scheduled to another schema-registry.
					}
				}
			}
		}
	})
}

func (tsw *TenantSendingWorker) Close() {
	tsw.sorterCli.Close()
	tsw.schemeCli.Close()
}
