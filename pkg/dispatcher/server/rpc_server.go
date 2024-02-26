package server

import (
	"context"

	"github.com/zhangjinpeng87/tistream/pkg/dispatcher/datachangebuffer"
	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
)

type dispatcherRpcServer struct {
	pb.UnimplementedDispatcherServiceServer

	// The data change buffer.
	dcbMgr *datachangebuffer.DataChangeBufferManager
}

func NewDispatcherRpcServer(dcbMgr *datachangebuffer.DataChangeBufferManager) *dispatcherRpcServer {
	return &dispatcherRpcServer{
		dcbMgr: dcbMgr,
	}
}

func (d *dispatcherRpcServer) TenantHasNewChanges(ctx context.Context, in *pb.HasNewChangeReq) (*pb.HasNewChangeResp, error) {
	// Tenant has new chagnes, the dispatcher should notify the worker to dispatch the data change files.
	for _, tenantId := range in.TenantId {
		d.dcbMgr.TenantHasNewChanges(tenantId)
	}

	return &pb.HasNewChangeResp{}, nil
}

func (d *dispatcherRpcServer) HandleTenantTask(ctx context.Context, in *pb.TenantTasksReq) (*pb.TenantTasksResp, error) {
	// Handle the tenant's tasks.
	err := d.dcbMgr.HandleTenantTasks(in)

	return &pb.TenantTasksResp{}, err
}
