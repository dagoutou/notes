package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go request("127.0.0.1:8080", &wg)
	}
	wg.Wait()
}
func request(url string, wg *sync.WaitGroup) {
	i := 0
	for {
		if _, err := http.Get(url); err == nil {
			break
		}
		i++
		if i >= 3 {
			fmt.Println("aaa")
			break
		}
		time.Sleep(time.Second)
	}
	wg.Done()
}

type base struct {
	Type string
}

func (b base) Class() string {
	return b.Type
}

type bird interface {
	Name() string
	Class() string
}
type canary struct {
	base
	name string
}

func (c canary) Name() string {
	return c.name
}

type crow struct {
	base
	name string
}

func (c crow) Name() string {
	return c.name
}
func info(b bird) {
	fmt.Printf("im %s,i belong to %s bird class!\n", b.Name(), b.Class())
}
