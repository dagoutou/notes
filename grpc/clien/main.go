package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	pb "notes/grpc/hello/proto"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("grpc.Dial error", err)
	}
	defer conn.Close()
	clien := pb.NewSayHelloClient(conn)
	resp, err := clien.SayHello(context.Background(), &pb.HelloRequest{RequestName: "狗头"})
	if err != nil {
		log.Fatal("error", err)
	}
	fmt.Println(resp.GetResponseMsg())

}
