package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	pb "notes/grpc/hello/proto"
)

type server struct {
	pb.UnimplementedSayHelloServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (resp *pb.HelloResponse, err error) {
	return &pb.HelloResponse{ResponseMsg: "hello" + req.RequestName}, nil
}
func main() {
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("listen error", err)
	}
	rpcServer := grpc.NewServer()
	pb.RegisterSayHelloServer(rpcServer, &server{})
	if err = rpcServer.Serve(listen); err != nil {
		log.Fatal("rpcServer.Serve error", err)
	}
}
