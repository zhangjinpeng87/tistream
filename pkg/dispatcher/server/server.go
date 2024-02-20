package dispatcher

import (
	"github.com/zhangjinpeng87/tistream/pkg/utils"
	pb "github.com/zhangjinpeng87/tistream/proto/go/dispatcher_grpc.pb.go"
)

// DispatchServer receives 2 types of messages from the meta server:
//  1) The notification of which tenant has new data change files, so that the dispatcher
//  can notify the worker start to dispatch the data change files for the corresponding
//  tenant.
//  2) Tenant task re-balancing messages, the dispatcher attaches or detaches the tenant's
//  tasks.
type DispatchServer struct {
	grpcServer *utils.GrpcServer
}

// NewDispatchServer creates a new DispatchServer.
func NewDispatchServer(addr string, port int) *DispatchServer {
	grpcServer := utils.NewGrpcServer(addr, port)

	grpcServer.RegisterService(&pb.DispatcherService_serviceDesc)
}

