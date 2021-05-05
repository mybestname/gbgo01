package main

import (
	"sync"
	"sync/atomic"
	"testing"
)

// 这里的解决办法是使用 atomic.Value
// 相对mutex来说
// two kinds of mutex with the sync package: sync.Mutex and sync.RWMutex; the latter is optimized when your program
// deals with multiples readers and very few writers.
func TestAtomic(t *testing.T) {
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
				t.Logf("%v\n",cfg.a)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
// GORACE="halt_on_error=1" go test -race .

func BenchmarkAtomic(b *testing.B) {
	var v atomic.Value
	var lastValue uint64
	v.Store(&Config{a:[]int{0,0,0,0,0}})

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
			for i := 0; i < b.N; i++ {
				cfg := v.Load().(*Config)
				atomic.SwapUint64(&lastValue, uint64(cfg.a[0]))
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkAtomic0Write(b *testing.B) {
	var v atomic.Value
	var lastValue uint64
	v.Store(&Config{a:[]int{0,0,0,0,0}})

	var wg sync.WaitGroup
	for n :=0 ; n < 4 ; n++ {
		wg.Add(1)
		go func() {
			for i := 0; i < b.N; i++ {
				cfg := v.Load().(*Config)
				atomic.SwapUint64(&lastValue, uint64(cfg.a[0]))
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkAtomicLarge(b *testing.B) {
	length := 50000000
	var v atomic.Value
	var lastValue uint64
	vals := make([]int,length)
	for i := 0; i < length; i++ {
		vals[i] = i
	}
	v.Store(&Config{a:vals})

	go func(){
		i := 0
		for {
			i++
			cfg := &Config{}
			vals := make([]int, length)
			for n := 0; n < length; i++ {
				vals[n] = n + i
			}
			cfg.a = vals
			v.Store(cfg)
		}
	}()

	var wg sync.WaitGroup
	for n :=0 ; n < 4 ; n++ {
		wg.Add(1)
		go func() {
			for i := 0; i < b.N; i++ {
				cfg := v.Load().(*Config)
				index := i
				if i >= length {
					index = index % length
				}
				atomic.SwapUint64(&lastValue, uint64(cfg.a[index]))
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

// $ go test -bench Atomic
// goos: darwin
// goarch: amd64
// pkg: atomic
// cpu: Intel(R) Core(TM) i7-7920HQ CPU @ 3.10GHz
// BenchmarkAtomic-8               11809771                95.15 ns/op
// BenchmarkAtomic0Write-8         13933311                90.81 ns/op
// BenchmarkAtomicLarge-8             14973             82367 ns/op
// PASS
// ok      atomic  11.803s

// $ go test -trace=trace.out
// $ go tool trace trace.out