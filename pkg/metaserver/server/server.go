package server

import (
	"context"

	"github.com/zhangjinpeng87/tistream/pkg/metaserver/metadata"
	"github.com/zhangjinpeng87/tistream/pkg/utils"
	"golang.org/x/sync/errgroup"

	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
)

type MetaServer struct {
	// Global configuration.
	globalConfig *utils.GlobalConfig

	// The gRPC server.
	grpcServer *utils.GrpcServer

	// Meta Backend. Used to store the metadata like tasks, tenants, lease, etc.
	backend *metadata.Backend

	// The task management.
	taskManagement *metadata.TaskManagement

	// Master campaign.
	campaign *Campaign
	eg       *errgroup.Group
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewMetaServer(globalConfig *utils.GlobalConfig) (*MetaServer, error) {
	// Initialize the backend.
	backend, err := metadata.NewBackend(&globalConfig.MetaServer)
	if err != nil {
		return nil, err
	}

	// Initialize the task management.
	taskManagement := metadata.NewTaskManagement(&globalConfig.MetaServer, backend)

	// Initialize the errgroup and context.
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)

	return &MetaServer{
		globalConfig:   globalConfig,
		taskManagement: taskManagement,
		backend:        backend,
		eg:             eg,
		ctx:            ctx,
		cancel:         cancel,
	}, nil
}

func (s *MetaServer) Prepare() error {
	// Initialize the schema if not exists.
	if err := s.taskManagement.Prepare(); err != nil {
		return err
	}

	// Start the master campaign.
	s.campaign = NewCampaign(&s.globalConfig.MetaServer)
	s.campaign.Start(s.taskManagement, s.eg, s.ctx)

	// Initialize the gRPC server.
	s.grpcServer = utils.NewGrpcServer(s.globalConfig.MetaServer.Addr, s.globalConfig.MetaServer.Port)
	pb.RegisterMetaServiceServer(s.grpcServer.InternalServer, NewMetaRpcServer(s.taskManagement))

	return nil
}

func (s *MetaServer) ReadyToServe() bool {
	// Only master meta-server can serve the requests from the client.
	return s.campaign.IsMaster() && s.taskManagement.DataIsReady()
}

func (s *MetaServer) Start() error {
	return s.grpcServer.Start(s.eg, s.ctx)
}

func (s *MetaServer) Wait() error {
	return s.eg.Wait()
}

func (s *MetaServer) Stop() error {
	s.cancel()
	s.eg.Wait()
	return s.grpcServer.Stop()
}
