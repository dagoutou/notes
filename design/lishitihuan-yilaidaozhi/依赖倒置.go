package main

import "fmt"

type car interface {
	running()
}
type Driver interface {
	driving(ca car)
}
type Bwm struct {
}

func (b *Bwm) running() {
	fmt.Println("bwm is running")
}

type Lisi struct {
}

func (l *Lisi) driving(ca car) {
	fmt.Println("李四开宝马")
	ca.running()
}

func main() {
	var bw car
	var d Driver
	bw = new(Bwm)
	d = new(Lisi)

	d.driving(bw)
}
