package main

import (
	"fmt"
	"sync"
	"testing"
)

func TestDataRace(t *testing.T) {
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
// $  GORACE="halt_on_error=1" go test -race .
//   ==================
//   WARNING: DATA RACE
//   Read at 0x00c000140030 by goroutine 9:
//     atomic.TestDataRace.func2()
//         /gbgo01/week03/u3_sync/05atomic/config_test.go:24 +0x4a
//
//   Previous write at 0x00c000140030 by goroutine 8:
//     atomic.TestDataRace.func1()
//         /gbgo01/week03/u3_sync/05atomic/config_test.go:15 +0x126
//
//   Goroutine 9 (running) created at:
//     atomic.TestDataRace()
//         /gbgo01/week03/u3_sync/05atomic/config_test.go:22 +0x115
//     testing.tRunner()
//         /usr/local/go/src/testing/testing.go:1193 +0x202
//
//   Goroutine 8 (running) created at:
//     atomic.TestDataRace()
//         /gbgo01/week03/u3_sync/05atomic/config_test.go:11 +0x8e
//     testing.tRunner()
//         /usr/local/go/src/testing/testing.go:1193 +0x202
//   ==================
//   FAIL    atomic  0.065s
//   FAIL