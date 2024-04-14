package main

import (
	"fmt"
	"log"
	"sync"
	"testing"
)

func TestGeneratorID(t *testing.T) {
	wg := sync.WaitGroup{}
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			orderID, er := generatorID("order")
			if er != nil {
				log.Fatal(er)
			}
			fmt.Println(orderID)
		}()
	}
	wg.Wait()
}
