package main

import "fmt"

// Go 1.1 includes a race detector
// The race detector is integrated with the go tool chain. When the -race command-line flag is set,
// the compiler instruments all memory accesses with code that records when and how the memory was
// accessed, while the runtime library watches for unsynchronized accesses to shared variables.
// 算法： https://github.com/google/sanitizers/wiki/ThreadSanitizerAlgorithm
// form https://github.com/golang/blog/blob/master/support/racy/racy.go
func main() {
	done := make(chan bool)
	m := make(map[string]string)
	m["name"] = "world"
	go func() {
		m["name"] = "data race"
		done <- true
	}()
	fmt.Println("Hello,", m["name"])
	<-done
}
// go build -race race_detector.go
// ./race_detector
// ==================
// WARNING: DATA RACE
// Write at 0x00c00007a180 by goroutine 7:
//   runtime.mapassign_faststr()
//       /usr/local/go/src/runtime/map_faststr.go:202 +0x0
//   main.main.func1()
//       gcgbgo01/week03/u3_sync/02race/race_detector.go:16 +0x5d
//
// Previous read at 0x00c00007a180 by main goroutine:
//   runtime.mapaccess1_faststr()
//       /usr/local/go/src/runtime/map_faststr.go:12 +0x0
//   main.main()
//       gcgbgo01/week03/u3_sync/02race/race_detector.go:19 +0x138
//
// Goroutine 7 (running) created at:
//   main.main()
//       gcgbgo01/week03/u3_sync/02race/race_detector.go:15 +0x109
// ==================
// ==================
// WARNING: DATA RACE
// Write at 0x00c00010a088 by goroutine 7:
//   main.main.func1()
//       gcgbgo01/week03/u3_sync/02race/race_detector.go:16 +0x72
//
// Previous read at 0x00c00010a088 by main goroutine:
//   main.main()
//       gcgbgo01/week03/u3_sync/02race/race_detector.go:19 +0x14b
//
// Goroutine 7 (running) created at:
//   main.main()
//       gcgbgo01/week03/u3_sync/02race/race_detector.go:15 +0x109
// ==================
// Found 2 data race(s)