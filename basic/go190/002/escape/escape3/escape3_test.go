// original code from with modification
// https://www.ardanlabs.com/blog/2017/06/language-mechanics-on-memory-profiling.html
// https://www.youtube.com/watch?v=2557w0qsDV0
// https://github.com/ardanlabs/gotraining/blob/master/topics/go/profiling/memcpu/stream.go
// Language Mechanics On Memory Profiling
package main

import (
	"bytes"
	"testing"
)

// Capture the time it takes to execute algorithm one.
func BenchmarkAlgorithmOne(b *testing.B) {
	var output bytes.Buffer
	in := assembleInputStream()
	find := []byte("elvis")
	repl := []byte("Elvis")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		output.Reset()
		algOne(in, find, repl, &output)
	}
}

// Capture the time it takes to execute algorithm two.
func BenchmarkAlgorithmTwo(b *testing.B) {
	var output bytes.Buffer
	in := assembleInputStream()
	find := []byte("elvis")
	repl := []byte("Elvis")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		output.Reset()
		algTwo(in, find, repl, &output)
	}
}

// Capture the time it takes to execute algorithm one.
func BenchmarkAlgorithmOneModified(b *testing.B) {
	var output bytes.Buffer
	in := assembleInputStream()
	find := []byte("elvis")
	repl := []byte("Elvis")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		output.Reset()
		algOne_modified(in, find, repl, &output)
	}
}
