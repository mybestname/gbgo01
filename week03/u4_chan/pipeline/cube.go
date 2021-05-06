package main

import (
	"fmt"
	"math/big"
)
// 3个goroutine 被连接在一起。
// n ( 0, 1 ,2 ,3 ...) -> c (0, 1, 8, 27 ...) -> (forever println)
func main() {
	n := make(chan int64)
	c := make(chan *big.Int)

	go func() {
		for i:=int64(0);;i++{
			n <- i
		}
	}()
	go func() {
		for {
			v := <-n
			var cube big.Int
			s := cube.Mul(big.NewInt(v),big.NewInt(v))
			cu := cube.Mul(s,big.NewInt(v))
			c <- cu
		}
	}()
	for i:=0;; i++{
		fmt.Printf("%d^3=%v\n",i,<-c)
	}
}
