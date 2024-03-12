package metaserver

import (
	"context"

	"github.com/zhangjinpeng87/tistream/pkg/metaserver/metadata"
	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
)

type metaRpcServer struct {
	pb.UnimplementedMetaServiceServer

	// The task management.
	taskManagement *metadata.TaskManagement
}

func NewMetaRpcServer(m *metadata.TaskManagement) *metaRpcServer {
	return &metaRpcServer{
		taskManagement: m,
	}
}

func (m *metaRpcServer) TenantHasNewChange(ctx context.Context, in *pb.HasNewChangeReq) (*pb.HasNewChangeResp, error) {
	// Tenant has new chagnes, the meta-server should notify the dispatcher to dispatch the data change files.
	// TODO: Implement this.

	return nil, nil
}

func (m *metaRpcServer) DispatcherHeartbeat(ctx context.Context, in *pb.DispatcherHeartbeatReq) (*pb.DispatcherHeartbeatResp, error) {
	// Meta-server receives the heartbeat from the dispatcher, update the

	return nil, nil
}

func (m *metaRpcServer) SorterHeartbeat(ctx context.Context, in *pb.SorterHeartbeatReq) (*pb.SorterHeartbeatResp, error) {
	// Meta-server receives the heartbeat from the sorter, update

	return nil, nil
}

func 