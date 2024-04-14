package main

import (
	"log"
	"net"
	"net/rpc"
	server2 "notes/rpc/demo2/server"
	"notes/rpc/demo3/server"
)

func main() {
	svc := server2.HelloService{}
	err := server.RegisterHelloService(&svc)
	if err != nil {
		log.Fatal("err:", err)
	}
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		return
	}
	conn, err := l.Accept()
	if err != nil {
		return
	}
	rpc.ServeConn(conn)
}
