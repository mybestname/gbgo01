package main

import (
	"fmt"
	"sync"
)

type Config struct {
	a []int
}

// from https://medium.com/a-journey-with-go/go-how-to-reduce-lock-contention-with-the-atomic-package-ba3b2664b549
func main() {
	cfg := &Config{}
	go func(){
		i := 0
		for {
			i++
			cfg.a = []int{i, i+1, i+2, i+3, i+4, i+5}
		}
	}()

	var wg sync.WaitGroup
	for n :=0 ; n < 4 ; n++ {
		wg.Add(1)
		go func() {
			for i := 0; i <100; i++ {
				fmt.Printf("%v\n",cfg.a)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
// $  GORACE="halt_on_error=1" go run -race .
//  ==================
//  WARNING: DATA RACE
//  Read at 0x00c0000b4018 by goroutine 8:
//    main.main.func2()
//        /gbgo01/week03/u3_sync/05atomic/config.go:29 +0x4a
//
//  Previous write at 0x00c0000b4018 by goroutine 7:
//    main.main.func1()
//        /gbgo01/week03/u3_sync/05atomic/config.go:20 +0x126
//
//  Goroutine 8 (running) created at:
//    main.main()
//        /gbgo01/week03/u3_sync/05atomic/config.go:27 +0x115
//
//  Goroutine 7 (running) created at:
//    main.main()
//        /gbgo01/week03/u3_sync/05atomic/config.go:16 +0x8e
//  ==================
//  exit status 66
