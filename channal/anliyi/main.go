package main

import (
	"fmt"
	"sync"
)

func main() {
	var chp = chPool{}
	chp.initChPool(11, 3)

	wg := &sync.WaitGroup{}
	wg.Add(4)
	chp.sendTask(1, func() {
		defer wg.Done()
		fmt.Println("我是任务1")
	})
	chp.sendTask(2, func() {
		defer wg.Done()
		fmt.Println("我是任务2")
	})
	chp.sendTask(3, func() {
		defer wg.Done()
		fmt.Println("我是任务3")
	})
	chp.sendTask(4, func() {
		defer wg.Done()
		fmt.Println("我是任务4")
	})
	wg.Wait()
}

type chPool struct {
	jobId   int
	workId  int
	taskNum int
	pool    []chan func()
}

func (c *chPool) initChPool(chSize int, taskNum int) {
	c.taskNum = taskNum
	for i := 1; i <= c.taskNum; i++ {
		ch := make(chan func(), chSize)
		c.pool = append(c.pool, ch)
		go c.task(ch, i)
	}
}
func (c *chPool) task(ch chan func(), index int) {
	for {
		f := <-ch
		f()
		fmt.Printf("%d号协程被用来执行任务了！\n", index)
	}

}
func (c *chPool) sendTask(taskId int, task func()) {
	num := taskId % c.taskNum
	tas := c.pool[num]
	tas <- task
}
