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
			oldMap[num] = cfg  // fatal error: concurrent map writes
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

// $ GORACE="halt_on_error=1" go run -race config.go
// ==================
// WARNING: DATA RACE
// Write at 0x00c000124180 by goroutine 8:
//   runtime.mapassign_fast64()
//       /usr/local/go/src/runtime/map_fast64.go:92 +0x0
//   main.main.func1()
//       /gbgo01/week03/u3_sync/06cow/config.go:28 +0x271
//
// Previous write at 0x00c000124180 by goroutine 7:
//   runtime.mapassign_fast64()
//       /usr/local/go/src/runtime/map_fast64.go:92 +0x0
//   main.main.func1()
//       /gbgo01/week03/u3_sync/06cow/config.go:28 +0x271
//
// Goroutine 8 (running) created at:
//   main.main()
//       /gbgo01/week03/u3_sync/06cow/config.go:33 +0x3e4
//
// Goroutine 7 (running) created at:
//   main.main()
//       /gbgo01/week03/u3_sync/06cow/config.go:32 +0x3c4
// ==================
// exit status 66
