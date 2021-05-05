package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Config struct {
	a []int
}

func main() {
	type Map map[int]*Config
	var v atomic.Value
	m := make(Map)
	m[1] = &Config{a: []int{0,0,0,0,0}}
	m[2] = &Config{a: []int{0,0,0,0,0}}
	v.Store(m)

	writer := func(num int ){
		i := 0
		for {
			i++
			oldMap := v.Load().(Map)
			cfg := &Config{}
			cfg.a = []int{i, i+num, i+2*num, i+3*num, i+4*num, i+5*num}
			oldMap[num] = cfg
			v.Store(oldMap)
		}
	}
	go writer(1)
	go writer(2)

	var wg sync.WaitGroup
	for n :=0 ; n < 4 ; n++ {

		wg.Add(1)
		go func(num int) {
			for i := 0; i <100; i++ {
				m := v.Load().(Map)
				fmt.Printf("r%d c1 %v %d\n",num, m[1].a, i)
				fmt.Printf("r%d c2 %v %d\n",num, m[2].a, i)
			}
			wg.Done()
		}(n+1)
	}
	wg.Wait()
}

// GORACE="halt_on_error=1" go run -race config.go
// WARNING: DATA RACE
// Write at 0x00c000142048 by goroutine 7:
// main.main.func1()
// /gbgo01/week03/u3_sync/06cow/config.go:30 +0x286
//
// Previous read at 0x00c000142048 by goroutine 9:
// main.main.func2()
// /gbgo01/week03/u3_sync/06cow/config.go:45 +0x7d
//
// Goroutine 7 (running) created at:
// main.main()
// /gbgo01/week03/u3_sync/06cow/config.go:35 +0x3c4
//
// r1 c1 [425 426 427 428 429 430] 11
// Goroutine 9 (running) created at:
// main.main()
// /gbgo01/week03/u3_sync/06cow/config.go:42 +0x485
// ==================
// exit status 66


