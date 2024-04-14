package main

import "fmt"

type catA interface {
	eat()
}
type catB struct {
	cat
}
type catC struct {
	catB
}
type cat struct {
}

func (c *cat) eat() {
	fmt.Println("小猫吃饭")
}

func (c *catB) sleep() {
	fmt.Println("小猫睡觉")
}
func (c *catC) run() {
	fmt.Println("小猫跑步")
}
func main() {
	a := catB{}
	a.eat()
	a.sleep()

}
