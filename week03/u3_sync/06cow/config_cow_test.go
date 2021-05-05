package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

type Map map[int]*Config

func TestCow(t *testing.T) {
	var v atomic.Value
	var mu sync.Mutex

	// init
	m := make(Map)
	m[1] = &Config{a: []int{0,0,0,0,0}}
	m[2] = &Config{a: []int{0,0,0,0,0}}
	v.Store(m)

	// cow
	writer := func(num int ){
		for i := 0; i< 1000; i++{

			mu.Lock()
			// first copy old one
			oldMap := v.Load().(Map)
			newMap := make(Map)
			for k,value := range oldMap{
				newMap[k] = value
			}
			cfg := &Config{}
			cfg.a = []int{i, i+num, i+2*num, i+3*num, i+4*num, i+5*num}
			newMap[num] = cfg
			v.Store(newMap)
			mu.Unlock()
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
				fmt.Printf("r%d c1 %v %d\n", num, m[1].a, i)
				fmt.Printf("r%d c2 %v %d\n", num, m[2].a, i)
			}
			m := v.Load().(Map)
			fmt.Printf("done! r%d c1 %v c2 %v \n", num, m[1].a, m[2].a)
			wg.Done()
		}(n+1)
	}
	wg.Wait()
}
// GORACE="halt_on_error=1" go test -race -v