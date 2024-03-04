package remotesorter

import (
	"context"

	"github.com/zhangjinpeng87/tistream/pkg/remotesorter/sorter"
	"github.com/zhangjinpeng87/tistream/pkg/storage"
	"github.com/zhangjinpeng87/tistream/pkg/utils"
	"golang.org/x/sync/errgroup"

	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
)

type SorterServer struct {
	// Global configuration.
	globalConfig *utils.GlobalConfig

	// The gRPC server.
	grpcServer *utils.GrpcServer

	// Sorter backend storage, used to store the snapshot and flush committed data.
	backend storage.ExternalStorage

	// Sorter Manager.
	sorterMgr *sorter.SorterManager

	eg     *errgroup.Group
	ctx    context.Context
	cancel context.CancelFunc
}

func NewSorterServer(globalConfig *utils.GlobalConfig) (*SorterServer, error) {
	// Initialize the backend.
	backend, err := storage.NewExternalStorage("S3", &globalConfig.Sorter.Storage)
	if err != nil {
		return nil, err
	}

	// Initialize the task management.
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)

	mgr := sorter.NewSorterManager(ctx, eg, backend)

	return &SorterServer{
		globalConfig: globalConfig,
		sorterMgr:    mgr,
		backend:      backend,
		eg:           eg,
		ctx:          ctx,
		cancel:       cancel,
	}, nil
}

func (s *SorterServer) Prepare() error {
	// Initialize the gRPC server.
	s.grpcServer = utils.NewGrpcServer(s.globalConfig.Sorter.Addr, s.globalConfig.Sorter.Port)
	pb.RegisterSorterServiceServer(s.grpcServer.InternalServer, NewSorterRpcServer(s.sorterMgr))

	return nil
}

func (s *SorterServer) Start() error {
	return s.grpcServer.Start(s.eg, s.ctx)
}

func (s *SorterServer) Wait() error {
	return s.eg.Wait()
}

func (s *SorterServer) Stop() error {
	s.cancel()
	s.eg.Wait()
	return s.grpcServer.Stop()
}
