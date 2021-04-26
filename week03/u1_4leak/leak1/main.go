package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	// Capture starting number of goroutines.
	startingGs := runtime.NumGoroutine()

	// leak
	leak();                           //goroutine inside never return

	// not leak
	fixV1();                          // print 0 (the zero value of chan int)
	time.Sleep(100*time.Millisecond)  // wait for v1 done, prettify output

	// not leak
	fixV2();                          // print nothing

	// not leak
	fixV3();                          // print 1, 2, 3

	// Hold the program from terminating for 1 second to see
	// if any goroutines created by process terminate.
	time.Sleep(time.Second)
	// Capture ending number of goroutines.
	endingGs := runtime.NumGoroutine()
	// Report the results.
	fmt.Println("========================================")
	//fmt.Println("Number of goroutines before:", startingGs)
	//fmt.Println("Number of goroutines after :", endingGs)
	fmt.Println("Number of goroutines leaked:", endingGs-startingGs)
}

// 问题函数 leak
func leak() {
	ch := make(chan int)
	go func() {
		val := <-ch                   // always blocked, waiting indefinitely <----+
		fmt.Println("received:", val) // will never happen.               //  |
	}()                                                                       //  |
	// the leak() returns. While that Goroutine is waiting -----------------------+
}

func fixV1() {
	ch := make(chan int)
	go func() {
		val := <-ch
	    fmt.Println("fix v1 received:", val)
	}()
	close(ch) //通道被关闭，任何接受操作立刻完成，同时获得与通道类型对应的零值。
}

func fixV2() {
	ch := make(chan int)
	go func() {
		for val := range ch {               // <--------------------------+
			fmt.Println("fix v2 received:", val)                  // |
		}                                                             // |
	}()                                                               // |
	close(ch) //通道被关闭，通道没有数据, 代码第60行不会执行，for循环直接结束，goroutine退出 --+
}

func fixV3() {
	ch := make(chan int)
	go func() {
		for val := range ch {
			fmt.Println("fix v3 received:", val)
		}
	}()
	ch<-1       //传输1从通道ch到goroutine，代码第70行被执行，随后继续阻塞在for循环
	ch<-2       //传输2
	ch<-3       //传输3
	close(ch)   //通道被关闭，goroutine中的for循环结束，goroutine退出。
}
