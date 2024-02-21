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
	// DB pool.
	dbPool *utils.DBPool

	// The gRPC server.
	grpcServer *utils.GrpcServer

	// The task management.
	taskManagement *metadata.TaskManagement

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

func NewMetaServer(globalConfig *utils.GlobalConfig) *MetaServer {
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)

	return &MetaServer{
		globalConfig:   globalConfig,
		taskManagement: metadata.NewTaskManagement(),
		eg:             eg,
		ctx:            ctx,
		cancel:         cancel,
	}
}

func (s *MetaServer) Prepare() error {
	// Connect to backend DB.
	dbPool, err := utils.NewDBPool(s.globalConfig.MetaServer.MysqlHost,
		s.globalConfig.MetaServer.MysqlPort,
		s.globalConfig.MetaServer.MysqlUser, s.globalConfig.MetaServer.MysqlPassword, "")
	if err != nil {
		return err
	}
	s.dbPool = dbPool

	// Initialize the schema if not exists.
	bootstrapper := metadata.NewBootstrapper(s.dbPool)
	if err := bootstrapper.InitSchema(); err != nil {
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
