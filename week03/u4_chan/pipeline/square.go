package main

import "fmt"
// 3个goroutine 被连接在一起。
// n ( 0, 1 ,2 ,3 ...) -> s (0, 1, 4, 9 ...) -> (forever println)
func main() {
	n := make(chan int)
	s := make(chan int)

	go func() {
		for i:=0;;i++{
			n <- i
		}
	}()
	go func() {
		for {
			v := <-n
			s <- v*v
		}
	}()
	for i:=0;; i++{
		fmt.Println(i,"^2=",<-s)
	}
}
