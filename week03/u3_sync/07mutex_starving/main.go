package main

import (
	"sync"
	"time"
)

// about starving
// https://github.com/golang/go/issues/13086
// https://github.com/rsc/tmp/blob/b6656a82b129/lockskew/lockskew.go
// 这种所谓的饥饿状态的情况（goroutine 2 拿不到锁，都被goroutine 1 抢走的情况），在go1.9中会自动判断这种情况
// 即如果 goroutine 超过 1ms 都没有获取到锁就会进饥饿模式 (https://github.com/golang/go/blob/master/src/sync/mutex.go#L139)

func main() {
	done := make(chan bool, 1)
	var mu sync.Mutex

	// goroutine 1
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				mu.Lock()
				time.Sleep(100*time.Microsecond)
				mu.Unlock()
			}
		}
	}()

	// goroutine 2
	for i:=0; i<10; i++ {
		time.Sleep(100*time.Microsecond)
		mu.Lock()
		mu.Unlock()
	}
	<-done

}
