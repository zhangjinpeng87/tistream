package schemaregistry

import (
	"context"

	"github.com/zhangjinpeng87/tistream/pkg/schemaregistry/schemastorage"
	"github.com/zhangjinpeng87/tistream/pkg/utils"
	"golang.org/x/sync/errgroup"
)

type SchemaServer struct {
	// Global configuration.
	globalConfig *utils.GlobalConfig

	// The gRPC server.
	grpcServer *utils.GrpcServer

	// Schema Manager.
	schemaMgr *schemastorage.SchemaManager

	eg     *errgroup.Group
	ctx    context.Context
	cancel context.CancelFunc
}

func NewSchemaServer(globalConfig *utils.GlobalConfig) (*SchemaServer, error) {
	// Initialize the task management.
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)

	mgr := schemastorage.NewSchemaManager()

	return &SchemaServer{
		globalConfig: globalConfig,
		schemaMgr:    mgr,
		eg:           eg,
		ctx:          ctx,
		cancel:       cancel,
	}, nil
}

func (s *SchemaServer) Start() error {
	return s.grpcServer.Start(s.eg, s.ctx)
}

func (s *SchemaServer) Wait() error {
	return s.eg.Wait()
}

func (s *SchemaServer) Stop() error {
	s.cancel()
	s.eg.Wait()
	return s.grpcServer.Stop()
}
