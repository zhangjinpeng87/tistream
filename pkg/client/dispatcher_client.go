package client

import (
	"context"

	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DispatcherClient struct {
	addr                 string
	conn                 *grpc.ClientConn
	dispatcherServiceCli pb.DispatcherServiceClient
}

func NewDispatcherClient(addr string) (*DispatcherClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	cli := pb.NewDispatcherServiceClient(conn)

	return &DispatcherClient{
		addr:                 addr,
		conn:                 conn,
		dispatcherServiceCli: cli,
	}, nil
}

func (c *DispatcherClient) Close() {
	c.conn.Close()
}

func (c *DispatcherClient) NotifyTenantHasUpdate(ctx context.Context, req *pb.HasNewChangeReq) error {
	// todo: handle response from dispatcher
	_, err := c.dispatcherServiceCli.NotifiyTenantHasUpdate(ctx, req)
	return err
}

func (c *DispatcherClient) ScheduleNewTenant(ctx context.Context, req *pb.TenantTasksReq) error {
	// todo: handle response from dispatcher
	_, err := c.dispatcherServiceCli.ScheduleNewTenant(ctx, req)
	return err
}
