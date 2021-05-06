package greeter

import (
	"context"
)

type Server struct {
	// listener net.Listener

	UnimplementedGreeterServer
}

func (*Server) SayHello(ctx context.Context, req *HelloRequest) (*HelloResponse, error) {
	// time.Sleep(200 * time.Millisecond)
	return &HelloResponse{Message: "Hello " + req.GetName()}, nil
}

// func (srv *Server) Run(addr string) error {
// 	l, err := net.Listen("tcp", addr)
// 	if err != nil {
// 		return err
// 	}
// 	srv.listener = l

// 	var opts []grpc.ServerOption
// 	grpcServer := grpc.NewServer(opts...)

// 	RegisterGreeterServer(grpcServer, srv)

// 	return grpcServer.Serve(l)
// }

// func (srv *Server) Addr() string {
// 	return srv.Addr()
// }
