package server

import (
	"context"

	"github.com/zhangjinpeng87/tistream/pkg/dispatcher/datachangebuffer"
	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
)

type dispatcherRpcServer struct {
	pb.UnimplementedDispatcherServiceServer

	// The reader and sender manager.
	readerMgr *datachangebuffer.ReaderManager
	senderMgr *datachangebuffer.SenderManager
}

func NewDispatcherRpcServer(readerMgr *datachangebuffer.ReaderManager, senderMgr *datachangebuffer.SenderManager) *dispatcherRpcServer {
	return &dispatcherRpcServer{
		readerMgr: readerMgr,
		senderMgr: senderMgr,
	}
}

func (d *dispatcherRpcServer) NotifiyTenantHasUpdate(ctx context.Context, in *pb.HasNewChangeReq) (*pb.HasNewChangeResp, error) {
	// Tenant has new chagnes, the dispatcher should notify the worker to dispatch the data change files.
	for _, tenantId := range in.TenantId {
		d.readerMgr.TenantHasNewChanges(tenantId)
	}

	return &pb.HasNewChangeResp{}, nil
}

func (d *dispatcherRpcServer) ScheduleNewTenant(ctx context.Context, in *pb.TenantTasksReq) (*pb.TenantTasksResp, error) {
	switch in.Op {
	case pb.TaskOp_Attach:
		sender, err := d.senderMgr.AttachTenant(in.TenantId)
		if err != nil {
			return nil, err
		}

		err = d.readerMgr.AttachTenant(in.TenantId, sender)
		if err != nil {
			return nil, err
		}
	case pb.TaskOp_Detach:
		err := d.readerMgr.DetachTenant(in.TenantId)
		if err != nil {
			return nil, err
		}

		d.senderMgr.DetachTenant(in.TenantId)
	}

	return &pb.TenantTasksResp{}, nil
}
