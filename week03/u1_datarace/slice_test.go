package main

import (
	"sync"
	"testing"
)

// go test ./.. -race
func TestAppend(t *testing.T) {
	// x := []string{"start"}  //OK
	// Data races don’t happen when multiple threads read memory, x, that doesn’t change.
	x := make([]string, 0, 6)  //FAIL TO race
	// The race happens because both goroutines are trying to write to
	// the same spare memory and it’s not clear who wins.

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		y := append(x, "hello", "world")
		t.Log(cap(y), len(y))
	}()
	go func() {
		defer wg.Done()
		z := append(x, "goodbye", "bob")
		t.Log(cap(z), len(z))
	}()
	wg.Wait()
}
