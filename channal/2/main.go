package main

import (
	"fmt"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	chWork := make(chan int, 100)
	chResult := make(chan int, 100)
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go work(i, &wg, chWork, chResult)
	}
	for i := 1; i <= 100; i++ {
		chWork <- i
	}
	close(chWork)
	go func() {
		wg.Wait()
		close(chResult)
	}()

	for i := range chResult {
		fmt.Println("输出结果:", i)
	}

}
func work(jobId int, wg *sync.WaitGroup, chWork <-chan int, chResult chan<- int) {
	defer wg.Done()
	for i := range chWork {
		fmt.Printf("使用 %d 号job, 执行第 %d 个任务,执行结果为 %d \n", jobId, i, i*2)
		chResult <- i * 2
	}
}
