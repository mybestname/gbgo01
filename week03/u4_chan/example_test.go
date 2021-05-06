package main

import (
	"fmt"
	"sync"
	"time"
)

func ExampleUnbufferedChannel() {
	channelSyntaxDemo1(0)
	//Output:
	//sent foo
	//Message: foo
	//sent bar
	//Message: bar
}

func ExampleBufferedChannel() {
	channelSyntaxDemo1(2)
	//Output:
	//sent foo
	//sent bar
	//Message: foo
	//Message: bar
}

func channelSyntaxDemo1(size int) {
	c := make(chan string,size)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		defer close(c)
		c <- `foo`
		fmt.Println("sent foo")
		c <- `bar`
		fmt.Println("sent bar")
	}()

	go func() {
		defer wg.Done()
		for {
			v, ok := <-c     // 如何判断通道是否关闭，false表示当前接收的通道读完且关闭。但需要注意的时候通道必须本身被关闭，<-c 语法会导致死锁。
			if !ok {
				break
			}
			time.Sleep(100*time.Millisecond)
			fmt.Println(`Message: ` + v)
		}
	}()
	wg.Wait()
}

// receiver是总是被block的，和chan是否buffered无关。
// 而sender是否被block，按是否buffer
func ExampleDeadlock() {
	//c := make(chan int64, 2)
	c := make(chan int64)
	go func() {
		c <- 1
		c <- 2
		// c <- 3
	}()
	fmt.Println(<-c)
	fmt.Println(<-c)
	fmt.Println(<-c)  //dead lock (block forever)

	//output:
	//1
	//2
	//3
}

func ExampleTimeOut() {
	// 3是总是被block的（无法打出），和chan是否buffered无关。
	blockTimeOut(0)
	blockTimeOut(2)
	// Output:
	// 1
	// 2
	// timeout 100ms
	// 1
	// 2
	// timeout 100ms
}

func blockTimeOut(bufSize int) {
	c := make(chan int64,bufSize)
	go func() {
		c <- 1
		c <- 2
	}()
	fmt.Println(<-c)
	fmt.Println(<-c)
	select {
	case v :=<-c:  //block forever
		fmt.Println(v)
	case <-time.After(100 * time.Millisecond):
		fmt.Println("timeout 100ms")
	}
}

func ExampleRange() {
	c := make(chan int)
	go func() {
		defer close(c)  //这句一定不能忘，否则就是死锁
		c <- 1
		c <- 2
	}()

	for v := range c {  //因为这种语法的前提是有人在某个时间点上需要关闭channel
		fmt.Println(v)
	}
	// Output:
	// 1
	// 2
}