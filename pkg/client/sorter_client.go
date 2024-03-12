package client

import (
	"context"

	"google.golang.org/grpc"

	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
)


// Current implementation is one tenant one sorter simple model.
// Todo:
// 1) support multiple sorters for one tenant.
// 2) update the sorter address when the tenant scheduled to another sorter.
type SorterClient struct {
	addr string
	conn *grpc.ClientConn

	client pb.SorterServiceClient
}

func NewSorterClient(addr string) (*SorterClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := pb.NewSorterServiceClient(conn)

	return &SorterClient{
		addr:   addr,
		conn:   conn,
		client: client,
	}, nil
}

func (s *SorterClient) Close() {
	s.conn.Close()
}

func (s *SorterClient) SendDataChanges(ctx context.Context, req *pb.RangeChangesReq) (*pb.RangeChangesResp, error) {
	return s.client.NewDataChanges(ctx, req)
}

func (s *SorterClient) RefreshAddr(newAddr string) error {
	if s.addr == newAddr {
		return nil
	}

	conn, err := grpc.Dial(newAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}

	s.conn.Close()
	s.conn = conn
	s.client = pb.NewSorterServiceClient(conn)
	s.addr = newAddr

	return nil
}
