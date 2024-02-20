package dispatcher

import (
	"github.com/zhangjinpeng87/tistream/pkg/utils"
	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
)

// DispatchServer receives 2 types of messages from the meta server:
//  1. The notification of which tenant has new data change files, so that the dispatcher
//     can notify the worker start to dispatch the data change files for the corresponding
//     tenant.
//  2. Tenant task re-balancing messages, the dispatcher attaches or detaches the tenant's
//     tasks.
type DispatchServer struct {
	grpcServer *utils.GrpcServer

	// The worker to handle the data change files.
}

// NewDispatchServer creates a new DispatchServer.
func NewDispatchServer(addr string, port int) *DispatchServer {
	grpcServer := utils.NewGrpcServer(addr, port)

	grpcServer.RegisterService(&pb.DispatcherService_serviceDesc)
}

// Start starts the server.
func (s *DispatchServer) Start() error {
	return s.grpcServer.Start()
}

// Stop stops the server.
func (s *DispatchServer) Stop() error {
	return s.grpcServer.Stop()
}
