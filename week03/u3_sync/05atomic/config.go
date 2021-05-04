package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Config struct {
	a []int
}

// from https://medium.com/a-journey-with-go/go-how-to-reduce-lock-contention-with-the-atomic-package-ba3b2664b549

func main() {
	var v atomic.Value
	v.Store(&Config{})

	go func(){
		i := 0
		for {
			i++
			cfg := &Config{}
			cfg.a = []int{i, i+1, i+2, i+3, i+4, i+5}
			v.Store(cfg)
		}
	}()

	var wg sync.WaitGroup
	for n :=0 ; n < 4 ; n++ {
		wg.Add(1)
		go func() {
			for i := 0; i <100; i++ {
				cfg := v.Load().(*Config)
				fmt.Printf("%v\n",cfg.a)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

//  GORACE="halt_on_error=1" go run -race .