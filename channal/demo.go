package main

import (
	"fmt"
	"sync"
)

// 协程池
type ChMgr struct {
	taskNum int
	ChList  []chan func()
	stat    bool
}

func (c *ChMgr) InitChMgr(num int) {
	c.taskNum = num
	for i := 0; i < c.taskNum; i++ {
		tmpCh := make(chan func(), 1000)
		c.ChList = append(c.ChList, tmpCh)
		go c.task(tmpCh, i)
	}
}

func (c *ChMgr) task(ch chan func(), index int) {

	for {
		t := <-ch
		t()
		fmt.Println("index", index)
	}
}

func (c *ChMgr) stop() {
	c.stat = false
}

func (c *ChMgr) sendTask(taskId int, task func()) {
	index := taskId % c.taskNum
	ch := c.ChList[index]
	ch <- task
}

var chPool ChMgr
var chPool2 *ChMgr

func main() {

	chPool.InitChMgr(1021)

	//var tmpList []uint64
	//mt := &sync.mut{}
	wg := &sync.WaitGroup{}
	wg.Add(2)
	chPool.sendTask(1, func() {
		// 查第一个
		defer wg.Done()
		fmt.Println("第一个")
		//mt.lock()
		//tmpList = append(tmpList , 1343)
		//mt.ulock()

	})

	chPool.sendTask(2, func() {
		// 查第二个
		defer wg.Done()
		fmt.Println("第二个")

	})
	wg.Wait()
	//chPool.sendTask(44434, func(){
	//	// 查第3个
	//	def wg.done()
	//	fmt.println("123456")
	//	mt.lock()
	//	tmpList = append(tmpList , 1343)
	//	mt.ulock()
	//})
	//
	//wg.wirt()
	//
	//
	//
	//2 个进程
	//
	//
	//50 个用户
	//30 用户是同一个组 1
	//20 用户是同一个组 2
	//
	//
	//
	//
	//按用户id分片
	//
	//30 用户请求到同一个进程
	//
	//
	//1，2，3 ，4 ，5，6 同一个  并发
	//同时请求到1号进程
	//
	//var groutMap1 map[int32]uint32
	//
	//userId % 2
	//chPool.sendTask(userInfo.groupId, func(){
	//
	//	// 都要修改 组信息的访问次数 map同时访问，会崩溃
	//	// 这样不用加锁
	//	tmp , ok := groutMap1[1]
	//	groutMap1[2] = 1
	//
	//})
}

// 对客户端的接口 tcp
//func sendCli(userId, groupInfo) {
//	chPool.sendTask(groupInfo.groupId, func(){
//
//		// 都要修改 组信息的访问次数 map同时访问，会崩溃
//		// 这样不用加锁
//		// tmp , ok := groutMap1[1]
//		groutMap1[userId] = userId
//
//
//	})
//}
//
//分流  1 % 1021   1
//
//进程内部 1 %1021 1
