package server

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/zhangjinpeng87/tistream/pkg/dispatcher/datachangebuffer"
	"github.com/zhangjinpeng87/tistream/pkg/storage"
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

	cfg *utils.DispatcherConfig

	// The data change buffer.
	dcbMgr *datachangebuffer.DataChangeBufferManager

	// The context to cancel the running goroutines.
	ctx    context.Context
	eg     *errgroup.Group
	cancel context.CancelFunc
}

// NewDispatchServer creates a new DispatchServer.
func NewDispatchServer(cfg *utils.DispatcherConfig) (*DispatchServer, error) {
	grpcServer := utils.NewGrpcServer(cfg.Addr, cfg.Port)
	ctx := context.Background()

	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)

	// create external storage
	backend, err := storage.NewExternalStorage("s3", &cfg.Storage)
	if err != nil {
		return nil, err
	}

	return &DispatchServer{
		grpcServer: grpcServer,
		dcbMgr:     datachangebuffer.NewDataChangeBufferManager(ctx, eg, cfg, backend),
		cfg:        cfg,
		ctx:        ctx,
		eg:         eg,
		cancel:     cancel,
	}, nil
}

// Prepare prepares the server.
func (s *DispatchServer) Prepare() error {
	s.grpcServer = utils.NewGrpcServer(s.cfg.Addr, s.cfg.Port)

	pb.RegisterDispatcherServiceServer(s.grpcServer.InternalServer, NewDispatcherRpcServer(s.dcbMgr))

	return nil
}

// Start starts the server.
func (s *DispatchServer) Start() error {
	// Start the data change buffer manager.
	s.dcbMgr.Run(s.ctx)

	// Start the gRPC server.
	err := s.grpcServer.Start(s.eg, s.ctx)
	if err != nil {
		return err
	}

	// Wait for the server to stop.
	return s.eg.Wait()
}

// Stop stops the server.
func (s *DispatchServer) Stop() error {
	s.cancel()
	return s.grpcServer.Stop()
}
