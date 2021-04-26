package main

import (
	"fmt"
	"sync"
	"testing"
)

func TestMapWrite(t *testing.T) {
	var counter = struct{
		sync.RWMutex
		m map[string]int
	}{m: make(map[string]int)}

	counter.RLock()
	n := counter.m["some_key"]
	counter.RUnlock()
	fmt.Println("some_key:", n)
}

func TestConcurrentWritesAfterGrowth(t *testing.T) {
	t.Parallel()
	numLoop := 10
	numGrowStep := 250
	numReader := 16

	for i := 0; i < numLoop; i++ {
		m := make(map[int]int, 0)
		for gs := 0; gs < numGrowStep; gs++ {
			m[gs] = gs
			var wg sync.WaitGroup
			wg.Add(numReader * 2)
			for nr := 0; nr < numReader; nr++ {
				go func() {
					defer wg.Done()
					for range m {
					}
				}()
				go func() {
					defer wg.Done()
					for key := 0; key < gs; key++ {
						_ = m[key]
						m[key] = key
					}
				}()
			}
			wg.Wait()
		}
	}
}