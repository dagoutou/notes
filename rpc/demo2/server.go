package main

import (
	"log"
	"net"
	"net/rpc"
	"notes/rpc/demo2/server"
)

func main() {
	rpc.RegisterName("HelloService", new(server.HelloService))
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("listen error", err)
	}
	conn, err := l.Accept()
	if err != nil {
		log.Fatal("Accept error", err)
	}
	rpc.ServeConn(conn)
}
