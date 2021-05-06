package main

import (
	"sync"
	"time"
)

// # 无缓冲channel，通信成功需要两个goroutine同时准备好。
// 含义：
//   - 1。无buff表示，当发送时，如果接收方没有准备好，那么锁住发送方，等待接收方。
//     2。接收方总是被锁住，等待是否有消息。
// 无缓冲的关键是: 如果接收没有准备好，那么发送这里会被锁住，直到接收方准备好。而缓冲则如果缓冲不满，发送这一句是不会被锁住的。
// 所以发送会不会成功（block住）是关键点！那么就可以保证只有消息被收到了，才会继续运行。（注意这里的保证，这里是一种同步保证），
// 但是这里这种保证是有代价的，代价就是发送会延迟，甚至一直被block。
//
// Rob Pike (Effective Go) https://golang.org/doc/effective_go#channels
// > - Unbuffered channels combine communication—the exchange of a value—with synchronization
// >   guaranteeing that two calculations (goroutines) are in a known state.
// > - Receivers always block until there is data to receive.
// > - If the channel is unbuffered, the sender blocks until the receiver has received the value.
// > - If the channel has a buffer, the sender blocks only until the value has been copied to the buffer;
// > - if the buffer is full, this means waiting until some receiver has retrieved a value.
//
// Go Programming Language Specification https://golang.org/ref/spec#Channel_types
// > - The capacity, in number of elements, sets the size of the buffer in the channel.
// > - If the capacity is zero or absent, the channel is unbuffered and communication succeeds only
// >   when both a sender and receiver are ready.
// > - Otherwise, the channel is buffered and communication succeeds without blocking if the buffer
// >   is not full (sends) or not empty (receives).
// > - A nil channel is never ready for communication.

// code from
// https://medium.com/a-journey-with-go/go-buffered-and-unbuffered-channels-29a107c00268
//
func main() {
	c := make(chan string, 2)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		c <- `foo`
		c <- `bar`
	}()

	go func() {
		defer wg.Done()

		time.Sleep(time.Second * 1)
		println(`Message: `+ <-c)
		println(`Message: `+ <-c)
	}()

	wg.Wait()
}
