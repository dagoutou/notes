package main

import "fmt"

// 司机张三 李四 汽车 宝马 奔驰
// 1 张三 开奔驰
// 2 李四 开宝马
func main() {
	z := zhangsan{}
	z.drive(&benchi{})
	l := lisi{}
	l.drive(&bwm{})
}

type benchi struct {
}

func (b *benchi) run() {
	fmt.Println("benchi is running")
}

type bwm struct {
}

func (b *bwm) run() {
	fmt.Println("bwm is running")
}

type zhangsan struct {
	name string
}

type lisi struct {
	name string
}

func (z *zhangsan) drive(b *benchi) {
	b.run()
	fmt.Println(z.name + ":" + "开奔驰")
}
func (z *lisi) drive(b *bwm) {
	b.run()
	fmt.Println(z.name + ":" + "开宝马")
}
