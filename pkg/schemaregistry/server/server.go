package schemaregistry

import (
	"context"

	"github.com/zhangjinpeng87/tistream/pkg/schemaregistry/schemastorage"
	"github.com/zhangjinpeng87/tistream/pkg/utils"
	"golang.org/x/sync/errgroup"

	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
)

type SchemaServer struct {
	// Configuration.
	cfg *utils.SchemaRegistryConfig

	// The gRPC server.
	grpcServer *utils.GrpcServer

	// Schema Manager.
	schemaMgr *schemastorage.SchemaManager

	eg     *errgroup.Group
	ctx    context.Context
	cancel context.CancelFunc
}

func NewSchemaServer(cfg *utils.SchemaRegistryConfig) (*SchemaServer, error) {
	// Initialize the task management.
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)

	mgr := schemastorage.NewSchemaManager()

	return &SchemaServer{
		cfg:       cfg,
		schemaMgr: mgr,
		eg:        eg,
		ctx:       ctx,
		cancel:    cancel,
	}, nil
}

// Prepare prepares the server.
func (s *SchemaServer) Prepare() error {
	s.grpcServer = utils.NewGrpcServer(s.cfg.Addr, s.cfg.Port)

	pb.RegisterSchemaServiceServer(s.grpcServer.InternalServer, NewSchemaRpcServer(s.schemaMgr))

	return nil
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
