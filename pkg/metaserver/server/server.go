package server

import (
	"github.com/zhangjinpeng87/tistream/pkg/metaserver/metadata"
	"github.com/zhangjinpeng87/tistream/pkg/utils"

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

}

func NewMetaServer(globalConfig *utils.GlobalConfig) *MetaServer {
	return &MetaServer{
		globalConfig:   globalConfig,
		taskManagement: metadata.NewTaskManagement(),
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
	// Load tasks from the backend DB.
	bootstrapper := metadata.NewBootstrapper(s.dbPool)
	if err := bootstrapper.Bootstrap(s.taskManagement); err != nil {
		return err
	}

	// Initialize the gRPC server.
	s.grpcServer = utils.NewGrpcServer(s.globalConfig.MetaServer.Addr, s.globalConfig.MetaServer.Port)
	s.grpcServer.RegisterService(&pb.MetaService_serviceDesc)

	return nil
}

func (s *MetaServer) Start() error {
	return s.grpcServer.Start()
}

