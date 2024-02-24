package main

import (
	"fmt"
	"log"
	"net/rpc"
	"notes/rpc/demo1/server"
)

func main() {
	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	args := &server.Args{A: 7, B: 8}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Arith:%d*%d=%d", args.A, args.B, reply)

	//quotient := new(server.Quotient)
	//divc := client.Go("Arith.Divide", args, quotient, nil)
	//replyCall := <-divc.Done
	//
	//if err != nil {
	//	log.Fatal("arith error:", err)
	//}
	//fmt.Printf("Arith:%d*%d=%d", args.A, args.B, reply)

}
