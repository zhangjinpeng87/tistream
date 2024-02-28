package server

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

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

	// last heartbeat meta-server
	lastMetaServerAddr       string
	lastSuccessHeartbeatTime time.Time
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

func (s *DispatchServer) Heartbeat() {
	s.eg.Go(func() error {
		timer := time.NewTicker(5 * time.Second)
		defer timer.Stop()
		for {
			select {
			case <-s.ctx.Done():
				return nil
			case <-timer.C:
				// Report the tenant stats.
				throughput, stats := s.dcbMgr.FetchStatsAndReset()
				heartbeat := &pb.DispatcherHeartbeatReq{
					Addr:        fmt.Sprintf("%s:%d", s.cfg.Addr, s.cfg.Port),
					Throughput:  throughput,
					TenantStats: stats,
				}
				// Send the heartbeat to the meta server.
				s.sendHeartbeat(heartbeat)
			}
		}
	})
}

func (s *DispatchServer) sendHeartbeat(heartbeat *pb.DispatcherHeartbeatReq) {
	// Send the heartbeat to the stored meta server, this is the most case.
	if len(s.lastMetaServerAddr) > 0 {
		// Try last meta server address first.
		if err := s.sendHeartbeatToMetaServer(s.lastMetaServerAddr, heartbeat); err == nil {
			s.lastSuccessHeartbeatTime = time.Now()
			return
		}
	}

	// For each meta server address, try to send the heartbeat.
	for _, addr := range s.cfg.MetaServerEndpoints {
		if err := s.sendHeartbeatToMetaServer(addr, heartbeat); err == nil {
			s.lastMetaServerAddr = addr
			return
		} else {
			// Todo: extract current meta server master from the heartbeat response and update the endpoints.
			// So in the case of meta-sever scale out/in, the dispatcher can send heartbeat to the correct master.
		}
	}

	// Failed to send heartbeat to all meta servers, if we cannot send heartbeat to any meta server for a long time,
	// we should stop the dispatcher. This may happen when there is a network partition between the dispatcher and
	// the meta server. In this case, the meta server would schedule these tenants' tasks to other dispatchers.
	if time.Since(s.lastSuccessHeartbeatTime) > time.Duration(s.cfg.SuicideDur)*time.Second {
		// Todo: Pause the dispatcher instead of stop until the received clear command from meta server.
		s.Stop()
	}
}

func (s *DispatchServer) sendHeartbeatToMetaServer(addr string, heartbeat *pb.DispatcherHeartbeatReq) error {
	conn, err := grpc.Dial(addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewMetaServiceClient(conn)
	_, err = client.DispatcherHeartbeat(context.Background(), heartbeat)
	return err
}
