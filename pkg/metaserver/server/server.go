package server

import (
	"context"
	"sync/atomic"

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

	// The task management.
	dataManagement *metadata.DataManagement

	// Master campaign.
	campaign *Campaign
	eg       *errgroup.Group
	ctx      context.Context
	cancel   context.CancelFunc

	// Indicate if this meta-server is master role.
	// Only master meta-server can serve the requests from the client.
	// Other-wise return the correct current master address to the client.
	isMaster atomic.Bool
}

func NewMetaServer(globalConfig *utils.GlobalConfig) (*MetaServer, error) {
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)

	dataManagement, err := metadata.NewDataManagement(&globalConfig.MetaServer)
	if err != nil {
		return nil, err
	}

	return &MetaServer{
		globalConfig:   globalConfig,
		dataManagement: dataManagement,
		eg:             eg,
		ctx:            ctx,
		cancel:         cancel,
	}, nil
}

func (s *MetaServer) Prepare() error {
	// Initialize the schema if not exists.
	if err := s.dataManagement.Backend().Bootstrap(); err != nil {
		return err
	}

	// Start the master campaign.
	s.campaign = NewCampaign(s.dbPool, &s.globalConfig.MetaServer)
	s.campaign.Start(s.taskManagement, &s.isMaster, s.eg, s.ctx)

	// Initialize the gRPC server.
	s.grpcServer = utils.NewGrpcServer(s.globalConfig.MetaServer.Addr, s.globalConfig.MetaServer.Port)
	s.grpcServer.RegisterService(&pb.MetaService_serviceDesc)

	return nil
}

func (s *MetaServer) ReadyToServe() bool {
	return s.isMaster.Load() && s.taskManagement.DataIsReady()
}

func (s *MetaServer) Start() error {
	return s.grpcServer.Start(s.eg, s.ctx)
}

func (s *MetaServer) Wait() error {
	return s.eg.Wait()
}

func (s *MetaServer) Stop() error {
	s.cancel()
	s.dbPool.Close()
	s.eg.Wait()
	return s.grpcServer.Stop()
}
