package remotesorter

import (
	"context"

	"github.com/zhangjinpeng87/tistream/pkg/remotesorter/sorter"
	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
)

type sorterRpcServer struct {
	pb.UnimplementedSorterServiceServer

	// The task management.
	sorterMgr *sorter.SorterManager
}

func NewSorterRpcServer(m *sorter.SorterManager) *sorterRpcServer {
	return &sorterRpcServer{
		sorterMgr: m,
	}
}

func (m *sorterRpcServer) NewDataChanges(ctx context.Context, in *pb.RangeChangesReq) (*pb.RangeChangesResp, error) {
	// Received the data changes from the dispatcher, update the sorter.
	// TODO: Implement this.

	return nil, nil
}
