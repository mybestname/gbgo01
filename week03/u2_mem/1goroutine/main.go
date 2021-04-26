package main

import "time"

func main() {
	println("before hello(), a=", a)
	hello()
	println("after hello(), a=", a)
	time.Sleep(100*time.Millisecond) //wait for fn()
	hello2()
}

// Goroutine creation

var a string

func f() {
	println(a)  // a一定是"hello"，但该语言也许要等hello()结束后才会执行
	//但是这里一定能保证a一定是hello，因为f()的创建一定在f的执行之前，所以a一定已经被设置为hello了。
}

func hello() {
	a = "hello"
	go f()  //The go statement that starts a new goroutine happens before the goroutine's execution begins.
}


// Goroutine destruction

var b string

func hello2() {
	go func() { b = "hello" }()   //The exit of a goroutine is not guaranteed to happen before any event in the program.
	println("b=",b)  //所以b很可能是空，因为没有并保证goroutine先结束。
}