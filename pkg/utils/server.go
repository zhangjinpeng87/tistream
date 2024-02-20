package utils

import (
	"net"

	"google.golang.org/grpc"
)

// GrpcServer is the generic tcp server which provide protobuf rpc service.
type GrpcServer struct {
	// The server address and port.
	addr string
	port int

	// The listener to listen to the server.
	listener *net.Listener

	// The internal grpc server.
	internalServer *grpc.Server
}

// NewGrpcServer creates a new TcpServer.
func NewGrpcServer(addr string, port int) *GrpcServer {
	return &GrpcServer{
		addr: addr,
		port: port,
		internalServer: grpc.NewServer(),
	}
}

// RegisterService registers the service to the server.
func (s *GrpcServer) RegisterService(sd *grpc.ServiceDesc) {
	s.internalServer.RegisterService(sd, struct{}{})
}

// Start starts the server.
func (s *GrpcServer) Start() error {
	// Start to listen to the server.
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	s.listener = &lis

	// Start the server.
	go s.internalServer.Serve(*s.listener)

	return nil
}

// Stop stops the server.
func (s *GrpcServer) Stop() error {
	s.internalServer.Stop()
	(*s.listener).Close()
	return nil
}