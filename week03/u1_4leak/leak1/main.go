package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	// Capture starting number of goroutines.
	startingGs := runtime.NumGoroutine()


	leak();


	// Hold the program from terminating for 1 second to see
	// if any goroutines created by process terminate.
	time.Sleep(time.Second)
	// Capture ending number of goroutines.
	endingGs := runtime.NumGoroutine()
	// Report the results.
	fmt.Println("========================================")
	fmt.Println("Number of goroutines before:", startingGs)
	fmt.Println("Number of goroutines after :", endingGs)
	fmt.Println("Number of goroutines leaked:", endingGs-startingGs)
}

// 问题函数 leak
//
func leak() {
	ch := make(chan int)
	go func() {
		val := <-ch
		fmt.Println("received:", val)
	}()
}
