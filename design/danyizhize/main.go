package main

import "fmt"

type clothes struct {
}

func (w *clothes) onWork() {
	fmt.Println("上班时的装扮")
}
func (w *clothes) onShop() {
	fmt.Println("逛街的装扮")
}
func main() {
	c := clothes{}
	fmt.Println("在上班")
	c.onWork()
	fmt.Println("在购物")
	c.onShop()
}

type clothesWork struct {
}
type clothesShop struct {
}
