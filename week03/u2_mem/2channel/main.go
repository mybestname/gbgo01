package main

func main() {
	go f()
	<-c  // receive 发生在 send 后
	println(a)  // 所以a 一定是hello
}

// Channel communication

var c = make(chan int, 10)
var a string

func f() {
	a = "hello"
	c <- 0  // A send on a channel happens before the corresponding receive from that channel completes.
}


// closing of a channel


