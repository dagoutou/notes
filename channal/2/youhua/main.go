package main

import (
	"fmt"
	"sync"
)

type chPools struct {
	chSize int
	jobNum int
	pools  []chan func()
	ch     chan func()
}

func (c *chPools) initPools(chSize, jobNum int) {
	c.chSize = chSize
	c.jobNum = jobNum
	for i := 1; i <= c.jobNum; i++ {
		ch := make(chan func(), c.chSize)
		c.ch = ch
		//c.pools = append(c.pools,ch)
		go c.exec(i, ch)
	}
}

func (c *chPools) exec(jobId int, ch chan func()) {
	for {
		fmt.Printf("%d 号job正在执行任务\n", jobId)
		c := <-ch
		c()
	}
}
func (c *chPools) addJob(jobId int, f func()) {
	//num := jobId % c.chSize
	//ch := c.pools[num]
	//ch <- f

	c.ch <- f
}
func main() {
	var chp chPools
	chp.initPools(100, 3)
	wg := &sync.WaitGroup{}
	wg.Add(4)
	chp.addJob(1, func() {
		defer wg.Done()
		fmt.Println("我是任务1")
	})
	chp.addJob(2, func() {
		defer wg.Done()
		fmt.Println("我是任务2")
	})
	chp.addJob(3, func() {
		defer wg.Done()
		fmt.Println("我是任务3")
	})
	chp.addJob(4, func() {
		defer wg.Done()
		fmt.Println("我是任务4")
	})
	wg.Wait()
}
