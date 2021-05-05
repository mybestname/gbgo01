package main

import (
	"sync"
	"sync/atomic"
	"testing"
)

// sync.RWMutex is optimized when your program
// deals with multiples readers and very few writers.
// 大量读，少量写的情节
func TestRWMutex(t *testing.T) {
	cfg := &Config{}
	var lock sync.RWMutex

	// 一个goroutine写
	go func(){
		i := 0
		for {
			i++
			lock.Lock()
			cfg.a = []int{i, i+1, i+2, i+3, i+4, i+5}
			lock.Unlock()
		}
	}()

	var wg sync.WaitGroup

	for n :=0 ; n < 4 ; n++ {
		wg.Add(1)
		// 5个goroutine读
		go func() {
			for i := 0; i <100; i++ {
				lock.RLock()
				t.Logf("%v\n",cfg.a)
				lock.RUnlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkMutex(b *testing.B) {
	var lastvalue uint64
	var lock sync.Mutex
	cfg := Config{
		a:[]int{0,0,0,0,0},
	}

	go func(){
		i := 0
		for {
			i++
			lock.Lock()
			cfg.a = []int{i, i+1, i+2, i+3, i+4, i+5}
			lock.Unlock()
		}
	}()

	var wg sync.WaitGroup
	for n :=0 ; n < 4 ; n++ {
		wg.Add(1)
		go func() {
			for i := 0; i <b.N; i++ {
				lock.Lock()
				atomic.SwapUint64(&lastvalue, uint64(cfg.a[0]))
				lock.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkMutex0Write(b *testing.B) {
	var lastvalue uint64
	var lock sync.Mutex
	cfg := Config{
		a:[]int{0,0,0,0,0},
	}

	var wg sync.WaitGroup
	for n :=0 ; n < 4 ; n++ {
		wg.Add(1)
		go func() {
			for i := 0; i <b.N; i++ {
				lock.Lock()
				atomic.SwapUint64(&lastvalue, uint64(cfg.a[0]))
				lock.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}


func BenchmarkRWMutex(b *testing.B) {
	var lastValue uint64
	var lock sync.RWMutex
	cfg := Config{
		a:[]int{0,0,0,0,0},
	}

	go func(){
		i := 0
		for {
			i++
			lock.Lock()
			cfg.a = []int{i, i+1, i+2, i+3, i+4, i+5}
			lock.Unlock()
		}
	}()

	var wg sync.WaitGroup
	for n :=0 ; n < 4 ; n++ {
		wg.Add(1)
		go func() {
			for i := 0; i <b.N; i++ {
				lock.RLock()
				atomic.SwapUint64(&lastValue, uint64(cfg.a[0]))
				lock.RUnlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}


func BenchmarkRWMutex0Write(b *testing.B) {
	var lastValue uint64
	var lock sync.RWMutex
	cfg := Config{
		a:[]int{0,0,0,0,0},
	}

	var wg sync.WaitGroup
	for n :=0 ; n < 4 ; n++ {
		wg.Add(1)
		go func() {
			for i := 0; i <b.N; i++ {
				lock.RLock()
				atomic.SwapUint64(&lastValue, uint64(cfg.a[0]))
				lock.RUnlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkRWMutexLarge(b *testing.B) {
	length := 50000000
	var lastValue uint64
	var lock sync.RWMutex
	vals := make([]int,length)
	for i := 0; i < length; i++ {
		vals[i] = i
	}
	cfg := Config{a:vals}

	go func(){
		i := 0
		for {
			i++
			vals := make([]int, length)
			for n := 0; n < length; n++ {
				vals[n] = n + i
			}
			lock.Lock()
			cfg.a = vals
			lock.Unlock()
		}
	}()

	var wg sync.WaitGroup
	for n :=0 ; n < 4 ; n++ {
		wg.Add(1)
		go func() {
			for i := 0; i <b.N; i++ {
				index := i
				if i >= length {
					index = index % length
				}
				lock.RLock()
				atomic.SwapUint64(&lastValue, uint64(cfg.a[index]))
				lock.RUnlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

// $ go test -bench .
// goos: darwin
// goarch: amd64
// pkg: atomic
// cpu: Intel(R) Core(TM) i7-7920HQ CPU @ 3.10GHz
// BenchmarkAtomic-8               11429712                91.60 ns/op
// BenchmarkMutex-8                 1395115              1080 ns/op
// BenchmarkMutex0Write-8           3761952               268.1 ns/op
// BenchmarkRWMutex-8               1000000              2150 ns/op
// BenchmarkRWMutex0Write-8         5852595               205.3 ns/op
// PASS
// ok      atomic  9.283s
// 可以看出:
// 1. RWMutex对于多数读少数写的情况下有优势。但是如果写也很多的情况下，反而不如Mutex
// 2. 而Atomic相对效率上明显。但是需要注意的是，如果你的数据结构非常大，那么copy复制的成本增加。