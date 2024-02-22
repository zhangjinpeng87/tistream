package server

import (
	"context"

	pb "github.com/zhangjinpeng87/tistream/proto/go/tistreampb"
)

type metaRpcServer struct {
	pb.UnimplementedMetaServiceServer

	
}

func NewMetaRpcServer() *metaRpcServer {
	return &metaHandler{}
}

func (h *metaHandler) 

