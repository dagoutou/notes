package main

import "fmt"

func main() {
	c := canary{
		Name: "金丝雀",
	}
	cr := crow{
		Name: "乌鸦",
	}
	letItFly(c)
	letItFly(cr)
}

type bird interface {
	Fly()
	Type() string
}
type canary struct {
	Name string
}

func (c canary) Fly() {
	fmt.Printf("我是 %s ，用黄色的翅膀飞\n", c.Name)
}
func (c canary) Type() string {
	return c.Name
}

type crow struct {
	Name string
}

func (c crow) Fly() {
	fmt.Printf("我是 %s ，用黑色的翅膀飞\n", c.Name)
}
func (c crow) Type() string {
	return c.Name
}
func letItFly(b bird) {
	b.Fly()
}
