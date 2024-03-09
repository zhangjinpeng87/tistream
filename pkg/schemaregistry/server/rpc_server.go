package schemaregistry

import (
	"context"

	"github.com/zhangjinpeng87/tistream/pkg/schemaregistry/schemastorage"
	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
)

type schemaRpcServer struct {
	pb.UnimplementedSchemaServiceServer

	schemaMgr *schemastorage.SchemaManager
}

func NewSorterRpcServer(m *schemastorage.SchemaManager) *schemaRpcServer {
	return &schemaRpcServer{
		schemaMgr: m,
	}
}

func (m *schemaRpcServer) RegisterSchemaSnap(ctx context.Context, in *pb.RegisterSchemaSnapReq) (*pb.RegisterSchemaSnapResp, error) {
	// Received the schema snap initialization from the dispatcher, setup schema snap for this tenant.
	// TODO: Implement this.

	return nil, nil
}

func (m *schemaRpcServer) RegisterDDLChange(ctx context.Context, in *pb.RegisterDDLReq) (*pb.RegisterDDLResp, error) {
	// Received the ddl changes from the dispatcher.
	// TODO: Implement this.

	return nil, nil
}
