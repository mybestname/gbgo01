package main
import (
	"fmt"
	"time"

	"go.uber.org/ratelimit"
)

func main() {

	prev := time.Now()
	rl := ratelimit.New(100) // per second

	for i := 0; i < 10; i++ {
		now := rl.Take()
		fmt.Println(i, now.Sub(prev))
		prev = now
	}

	// Output:
	// 0 0
	// 1 10ms
	// 2 10ms
	// 3 10ms
	// 4 10ms
	// 5 10ms
	// 6 10ms
	// 7 10ms
	// 8 10ms
	// 9 10ms
}
