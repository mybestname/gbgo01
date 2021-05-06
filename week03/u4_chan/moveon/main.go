package main

import (
	"fmt"
	"time"
)

//
func main() {
	cons := []Conn{ Conn1{}, Conn2{}, Conn3{}}
	result := Query(cons, "test")
	fmt.Println(result.name)  // result1 (因为最快）
}

type Result struct {
	name string
}
type Conn interface{
	DoQuery(query string) Result
}

type Conn1 struct{}
type Conn2 struct{}
type Conn3 struct{}

func (c Conn1) DoQuery(query string) Result {
	time.Sleep(100*time.Millisecond)
	return Result{name:"result in 100ms"}
}
func (c Conn2) DoQuery(query string) Result {
	time.Sleep(200*time.Millisecond)
	return Result{name:"result in 200ms"}
}
func (c Conn3) DoQuery(query string) Result {
	time.Sleep(300*time.Millisecond)
	return Result{name:"result in 300ms"}
}

// code from https://blog.golang.org/concurrency-timeouts
//
func Query(conns []Conn, query string) Result {
	ch := make(chan Result)      //注意！这里是阻塞channel
	for _, conn := range conns {
		go func(c Conn) {
			// the closure does a non-blocking send, which it achieves by using the send operation in
			// select statement with a default case. If the send cannot go through immediately
			// the default case will be selected.
			select {
			case ch <- c.DoQuery(query):  //这个send一定是执行最快的query，因为第一个是不存塞的，而后面的被阻塞
			default:                      //但是这时候可以通过default走掉，即使send被阻塞，即其它的慢查询走了default
			}
		}(conn)
	}
	return <-ch  //这里收到的，一定是最快的查询（即第一个send）
}
