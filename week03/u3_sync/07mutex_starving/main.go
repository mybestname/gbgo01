package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// about starving
// https://github.com/golang/go/issues/13086
// https://github.com/rsc/tmp/blob/b6656a82b129/lockskew/lockskew.go
// 这种所谓的饥饿状态的情况（goroutine 2 拿不到锁，都被goroutine 1 抢走的情况），在go1.9中会自动判断这种情况
// 即如果 goroutine 超过 1ms 都没有获取到锁就会进饥饿模式 (https://github.com/golang/go/blob/master/src/sync/mutex.go#L139)

var (
	// A duration string such as "300ms", "-1.5h" or "2h45m".
	// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
	D1  = flag.Duration("d1", 100*time.Microsecond, "goroutine 1 time to sleep")
	D2  = flag.Duration("d2", 100*time.Microsecond, "goroutine 2 time to sleep")
	N   = flag.Int("n", 10, "goroutine2 count number")
)


func main() {
	done := make(chan struct{})
	stop := make(chan struct{})
	var mu sync.Mutex
	var c1 uint64
	var c2 uint64

	flag.Parse()

	fmt.Printf("g1.sleep=%v, g2.sleep=%v, g2.iter=%v \n", *D1,*D2, *N)

	start := time.Now()

	go func() {
		interruptChannel := make(chan os.Signal, 1)
		interruptSignals := []os.Signal{os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM}
		signal.Notify(interruptChannel, interruptSignals...)
		select {
		case <-interruptChannel:
			println()
			stop<- struct{}{}
		}
	}()

	// goroutine 1
	go func() {
		for {
			select {
			case <-done:
				//fmt.Printf("goroutine 1 stopped!\n")
				stop<- struct{}{}
				return
			default:
				mu.Lock()
				time.Sleep(*D1)  // 100 mics by default
				c1++
				mu.Unlock()
			}
		}
	}()

	// goroutine 2
	go func() {
		for i:=0; i<*N; i++ {
			time.Sleep(*D2)          // 100 s (default)
			mu.Lock()
			c2++
			mu.Unlock()
		}

		done<- struct{}{}
	}()

	<-stop

	fmt.Printf("goroutine 1 lock acquired : %d \n", c1)
	fmt.Printf("goroutine 2 lock acquired : %d \n", c2)

	end := time.Now()
	elapsed := end.Sub(start)
	fmt.Printf("-- the total time elapsed : %v!\n", elapsed)

}
//
// $ go run main.go
// g1.sleep=100µs, g2.sleep=100µs, g2.iter=10
// goroutine 1 lock acquired : 100
// goroutine 2 lock acquired : 10
// -- the total time elapsed : 13.408418ms!
//
// $ go run main.go  -d1 1ns -d2 1ns -n 10000
// g1.sleep=1ns, g2.sleep=1ns, g2.iter=10000
// goroutine 1 lock acquired : 318253
// goroutine 2 lock acquired : 10000
// -- the total time elapsed : 148.534405ms!
//
// $ go run main.go  -d1 1000ns -d2 1ns -n 10000
// g1.sleep=1µs, g2.sleep=1ns, g2.iter=10000
// goroutine 1 lock acquired : 269032
// goroutine 2 lock acquired : 10000
// -- the total time elapsed : 1.284663781s!
//
// $ go run main.go  -d1 1ns -d2 1000ns -n 10000
// g1.sleep=1ns, g2.sleep=1µs, g2.iter=10000
// goroutine 1 lock acquired : 1355597
// goroutine 2 lock acquired : 10000
// -- the total time elapsed : 599.220628ms!
//
// $ go run main.go  -d1 1ns -d2 1000ns -n 100000
// g1.sleep=1ns, g2.sleep=1µs, g2.iter=100000
// goroutine 1 lock acquired : 12223047
// goroutine 2 lock acquired : 100000
// -- the total time elapsed : 5.522939886s!
//