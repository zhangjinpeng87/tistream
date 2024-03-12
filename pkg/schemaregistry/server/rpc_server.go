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

func NewSchemaRpcServer(m *schemastorage.SchemaManager) *schemaRpcServer {
	return &schemaRpcServer{
		schemaMgr: m,
	}
}

func (m *schemaRpcServer) RegisterSchemaSnap(ctx context.Context, in *pb.RegisterSchemaSnapReq) (*pb.RegisterSchemaSnapResp, error) {
	// Create schema storage from the schema snap.
	s := schemastorage.NewSchemaStorageFromSnapshot(in.TenantId, in)
	m.schemaMgr.AddSchemaStorage(s)

	// Create the response.
	resp := &pb.RegisterSchemaSnapResp{
		TenantId: in.TenantId,
		/*todo: other response fields*/
	}

	return resp, nil
}

func (m *schemaRpcServer) RegisterDDLChange(ctx context.Context, in *pb.RegisterDDLReq) (*pb.RegisterDDLResp, error) {
	resp := &pb.RegisterDDLResp{}
	resp.TenantId = in.TenantId

	switch in.Ddl.(type) {
	case *pb.RegisterDDLReq_TableDdl:
		t := in.Ddl.(*pb.RegisterDDLReq_TableDdl).TableDdl
		err := m.schemaMgr.AddTableDdl(in.TenantId, t)
		resp.ErrMsg = err.Error()
	case *pb.RegisterDDLReq_SchemaDdl:
		s := in.Ddl.(*pb.RegisterDDLReq_SchemaDdl).SchemaDdl
		err := m.schemaMgr.AddSchemaDdl(in.TenantId, s)
		resp.ErrMsg = err.Error()
	}

	return resp, nil
}
