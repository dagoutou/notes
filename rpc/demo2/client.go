package main

import (
	"fmt"
	"log"
	"net/rpc"
)

func main() {
	clent, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dial error", err)
	}
	var reply string
	if err = clent.Call("HelloService.Hello", "hello", &reply); err != nil {
		log.Fatal("Call error", err)
	}
	fmt.Println(reply)
}
