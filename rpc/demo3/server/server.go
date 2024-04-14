package server

import (
	"net/rpc"
)

const HelloServiceName = "notes/rpc/demo2/server.HelloService"

type HelloServiceInterface interface {
	Hello(request string, reply *string) error
}

func RegisterHelloService(svc HelloServiceInterface) error {
	return rpc.RegisterName(HelloServiceName, svc)
}
