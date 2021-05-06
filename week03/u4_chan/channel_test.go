package main

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestNoBuffer(t *testing.T) {
	runWithBuffer(0)
}
func TestBufferSize1(t *testing.T) {
	runWithBuffer(1)
}
func TestBufferSize5(t *testing.T) {
	runWithBuffer(50)
}
func TestBufferSize25(t *testing.T) {
	runWithBuffer(25)
}
func TestBufferSize50(t *testing.T) {
	runWithBuffer(50)
}
//
// $ go test . -trace trace.out -v
// === RUN   TestNoBuffer
// --- PASS: TestNoBuffer (0.00s)
// === RUN   TestBufferSize1
// --- PASS: TestBufferSize1 (0.00s)
// === RUN   TestBufferSize5
// --- PASS: TestBufferSize5 (0.00s)
// === RUN   TestBufferSize25
// --- PASS: TestBufferSize25 (0.00s)
// === RUN   TestBufferSize50
// --- PASS: TestBufferSize50 (0.00s)
// PASS
// ok      chan    0.061s
//
// $ go tool trace trace.out
// nobuffer 2768us vs. 628us 50buff


func runWithBuffer(size int) {
	c := make(chan uint32, size)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		for i := uint32(0); i < 1000; i++ {
			c <- i%2
		}
		close(c)
	}()

	var total uint32
	for w := 0; w < 5; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				v, ok := <-c
				if !ok {
					break
				}
				atomic.AddUint32(&total, v)
			}
		}()
	}
	wg.Wait()
}

func BenchmarkWithNoBuffer(b *testing.B) {
	benchmarkWithBuffer(b, 0)
}

func BenchmarkWithBufferSizeOf1(b *testing.B) {
	benchmarkWithBuffer(b, 1)
}

func BenchmarkWithBufferSizeEqualsToNumberOfWorker(b *testing.B) {
	benchmarkWithBuffer(b, 5)
}

func BenchmarkWithBufferSizeExceedsNumberOfWorker(b *testing.B) {
	benchmarkWithBuffer(b, 25)
}

func benchmarkWithBuffer(b *testing.B, size int) {
	for i := 0; i < b.N; i++ {
		c := make(chan uint32, size)

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()

			for i := uint32(0); i < 1000; i++ {
				c <- i%2
			}
			close(c)
		}()

		var total uint32
		for w := 0; w < 5; w++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				for {
					v, ok := <-c
					if !ok {
						break
					}
					atomic.AddUint32(&total, v)
				}
			}()
		}

		wg.Wait()
	}
}
// $ go get golang.org/x/perf/cmd/benchstat
// $ for i in {1..20}; do go test -bench . >> run.txt; done
// $ benchstat run.txt
// name                                    time/op
// WithNoBuffer-8                          361µs ± 7%
// WithBufferSizeOf1-8                     301µs ± 8%
// WithBufferSizeEqualsToNumberOfWorker-8  212µs ± 3%
// WithBufferSizeExceedsNumberOfWorker-8   142µs ± 4%
//
