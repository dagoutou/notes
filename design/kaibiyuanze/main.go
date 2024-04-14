package main

import "fmt"

// type banker struct {
//
// }
//
//	func (b *banker) save(){
//		fmt.Println("进行了存款业务")
//	}
//
//	func (b *banker) tansfer(){
//		fmt.Println("进行了转账业务")
//	}
//
//	func (b *banker) pay(){
//		fmt.Println("进行了支付业务")
//	}
//
//	func (b *banker) shares(){
//		fmt.Println("进行了股票业务")
//	}
type abstractBanker interface {
	doBusy()
}
type save struct {
}

func (s *save) doBusy() {
	fmt.Println("进行了存款业务")
}
func doBusy(a abstractBanker) {
	a.doBusy()
}
func main() {
	s := &save{}
	s.doBusy()
	doBusy(&save{})
}
