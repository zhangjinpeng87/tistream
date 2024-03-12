package client

import (
	"context"
	"fmt"

	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
	"google.golang.org/grpc"
)

// Todo: handle meta-server scale-out and scale-in cases, dynamically add or remove meta-server addr.
type MetaServerClient struct {
	addrs                []string
	conns                []*grpc.ClientConn
	schemaServiceClients []pb.MetaServiceClient
}

func NewMetaServerClient(addrs []string) (*MetaServerClient, error) {
	conns := make([]*grpc.ClientConn, 0, len(addrs))
	schemaServiceClients := make([]pb.MetaServiceClient, 0, len(addrs))
	for _, addr := range addrs {
		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		conns = append(conns, conn)
		schemaServiceClients = append(schemaServiceClients, pb.NewMetaServiceClient(conn))
	}

	return &MetaServerClient{
		addrs:                addrs,
		conns:                conns,
		schemaServiceClients: schemaServiceClients,
	}, nil
}

func (c *MetaServerClient) Close() {
	for _, conn := range c.conns {
		conn.Close()
	}
}

func (c *MetaServerClient) FetchSorterAddr(ctx context.Context, req *pb.FetchRangeSorterAddrReq) (string, error) {
	// Todo: handle "meta-server not leader" error, and retry to fetch sorter addr from the new leader.
	for _, client := range c.schemaServiceClients {
		resp, err := client.FetchRangeSorterAddr(ctx, req)
		if err == nil {
			return resp.Addr, nil
		}
	}

	return "", fmt.Errorf("cann't connect to any meta server to fetch sorter addr")
}

func (c *MetaServerClient) FetchSchemaRegistryAddr(ctx context.Context, req *pb.FetchSchemaRegistryAddrReq) (string, error) {
	// Todo: handle "meta-server not leader" error, and retry to fetch sorter addr from the new leader.
	for _, client := range c.schemaServiceClients {
		resp, err := client.FetchSchemaRegistryAddr(ctx, req)
		if err == nil {
			return resp.Addr, nil
		}
	}

	return "", fmt.Errorf("cann't connect to any meta server to fetch schema registry addr")
}
