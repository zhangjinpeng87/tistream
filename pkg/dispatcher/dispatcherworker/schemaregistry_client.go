package dispatcherworker

import (
	"context"

	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
	"google.golang.org/grpc"
)

type SchemaRegistryClient struct {
	addr   string
	conn   *grpc.ClientConn
	client pb.SchemaServiceClient
}

func NewSchemaRegistryClient(addr string) (*SchemaRegistryClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &SchemaRegistryClient{
		addr:   addr,
		conn:   conn,
		client: pb.NewSchemaServiceClient(conn),
	}, nil
}

func (c *SchemaRegistryClient) Close() {
	c.conn.Close()
}

func (c *SchemaRegistryClient) SendSchemaSnap(ctx context.Context, req *pb.RegisterSchemaSnapReq) (*pb.RegisterSchemaSnapResp, error) {
	return c.client.RegisterSchemaSnap(ctx, req)
}

func (c *SchemaRegistryClient) SendDDLChange(ctx context.Context, req *pb.RegisterDDLReq) (*pb.RegisterDDLResp, error) {
	return c.client.RegisterDDLChange(ctx, req)
}

func (c *SchemaRegistryClient) RefreshAddr(newAddr string) error {
	if c.addr == newAddr {
		return nil
	}

	conn, err := grpc.Dial(newAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}

	c.conn.Close()
	c.conn = conn
	c.client = pb.NewSchemaServiceClient(conn)
	c.addr = newAddr

	return nil
}
