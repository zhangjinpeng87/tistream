package datachangebuffer

import (
	"context"
	"fmt"

	"github.com/zhangjinpeng87/tistream/pkg/client"
	"github.com/zhangjinpeng87/tistream/pkg/codec"
	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
	"golang.org/x/sync/errgroup"
)

type TenantDataSender struct {
	tenantId uint64
	ch       <-chan *codec.FileEvent

	// Meta server client. Used to get the sorter address and schema-registry address.
	// And update the sorter address and schema-registry address when the tenant scheduled to another sorter or schema-registry.
	metaSrvCli *client.MetaServerClient

	// The remote sorter client.
	sorterCli *client.SorterClient

	// The schema registry client.
	schemeCli *client.SchemaRegistryClient
}

func NewTenantDataSender(tenantId uint64, ch <-chan *codec.FileEvent, metaSrvCli *client.MetaServerClient) *TenantDataSender {
	return &TenantDataSender{
		tenantId:   tenantId,
		ch:         ch,
		metaSrvCli: metaSrvCli,
	}
}

func (tsw *TenantDataSender) Run(ctx context.Context, eg *errgroup.Group) {

	// Get the sorter address and schema-registry address.
	sorterAddr, err := tsw.metaSrvCli.FetchSorterAddr(ctx, &pb.FetchRangeSorterAddrReq{TenantId: tsw.tenantId})
	if err != nil {
		panic(fmt.Sprintf("fetch sorter addr failed, tenantID: %d, err: %v", tsw.tenantId, err))
	}
	sorterCli, err := client.NewSorterClient(sorterAddr)
	if err != nil {
		panic(fmt.Sprintf("create sorter client failed, tenantID: %d, err: %v", tsw.tenantId, err))
	}

	schemaRegistryAddr, err := tsw.metaSrvCli.FetchSchemaRegistryAddr(ctx, &pb.FetchSchemaRegistryAddrReq{TenantId: tsw.tenantId})
	if err != nil {
		panic(fmt.Sprintf("fetch schema registry addr failed, tenantID: %d, err: %v", tsw.tenantId, err))
	}
	schemaCli, err := client.NewSchemaRegistryClient(schemaRegistryAddr)
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

func (tsw *TenantDataSender) Close() {
	tsw.sorterCli.Close()
	tsw.schemeCli.Close()
}
