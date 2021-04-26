package main

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"time"
)

func main() {
	// Capture starting number of goroutines.
	startingGs := runtime.NumGoroutine()

	term:="test"
	if err := processV1(term); err!=nil {
		fmt.Println(err)
	}

	ctx, cancel :=context.WithTimeout(context.Background(),100*time.Millisecond)
	defer cancel()

	// leak
	if err := processV2(term, ctx); err!=nil {
		fmt.Println(err)
	}

	// not leak
	if err := processV3(term, ctx); err!=nil {
		fmt.Println(err)
	}

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

func search(term string) (string, error) {
	time.Sleep(200*time.Millisecond)
	return "find foo"+term, nil
}

func processV1(term string) error {
	if record, err := search(term); err!=nil {
		return err
	}else {
		fmt.Println("v1 received:",record)
	}
	return nil
}

type result struct{
	record string
	err error
}

// leak!!!
func processV2(term string, ctx context.Context) error {
	ch := make(chan result)

	go func() {  // <--------------------------------------------------+
		record, err := search(term)                                 // |
		ch <- result{record,err}                                  // |
		//无缓冲通道，所以这里会阻塞住，直到有人接受，所以这个goroutine无法结束，除非有人读（即代码85行必须执行）
	}()                                                             // |
	select {                                                        // |
		case <-ctx.Done():                                          // |
			return errors.New("ctx timeout, search canceled")  // |
			//如果done先执行，直接return，那上面👆--------------------------+
			//on line 72 it sends on the channel. Sending on this channel blocks execution
			//until another Goroutine is ready to receive the value. In the timeout case,
			//the receiver stops waiting to receive from the Goroutine and moves on. This
			//will cause the Goroutine to block forever waiting for a receiver to appear
			//which can never happen. This is when the Goroutine leaks.
			//这个goroutine就永远不会停止，造成goroutine泄露。
		case r := <-ch:
			if r.err!= nil { return r.err }
			fmt.Println("v2 received:", r.record)
			return nil
	}

}

func processV3(term string, ctx context.Context) error {
	// The simplest way to resolve this leak is to change the channel from an
	// unbuffered channel to a buffered channel with a capacity of 1.
	// 那么，为什么改为buffer为1的chan就可以了呢？
	ch := make(chan result,1)
	go func() {
		record, err := search(term)
		ch <- result{record,err}
		//因为这里是一个缓冲通道，所以只要还有缓冲，那么就不会阻塞，所以这个goroutine直接结束了，
		//并不会等待代码第109行的的读取。
	}()
	select {
	case <-ctx.Done():
		return errors.New("ctx timeout, search canceled")
	case r := <-ch:
		if r.err!= nil { return r.err }
		fmt.Println("v3 received:", r.record)
		return nil
	}

}