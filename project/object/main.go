package main

import "fmt"

func main() {
	c := canary{
		base: base{"金丝雀"},
		name: "金丝雀A",
	}
	cr := crow{base: base{"乌鸦"},
		name: "乌鸦A",
	}
	info(c)
	info(cr)
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
