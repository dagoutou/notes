package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"notes/rpc/demo1/server"
)

func main() {
	arith := new(server.Arith)
	rpc.Register(arith)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("listen error:", err)
	}
	go http.Serve(l, nil)
	select {}
}
