package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"golang.org/x/time/rate"
)

var _ rate.Limit

func main() {
	const (
		n  = 10
		tf = 5 * time.Second
		b  = 1
	)
	lim := rate.NewLimiter(rate.Every(tf/n), b)
	// _ = lim.ReserveN(time.Now(), b)
	//<<<--- see what happens if you start with an empty bucket
	fmt.Println("expected limit:", lim.Limit())
	fmt.Println("completed elapsed actual rate")
	ctx := context.Background()
	start := time.Now()
	for i := 0; i < 2*n+2; i++ {
		err := lim.Wait(ctx)
		if err != nil {
			log.Fatal(err)
		}
		elapsed := time.Since(start)
		completed := i + 1
		// <<<--- see what happens if you did
		// completed := i
		actual := float64(completed) / elapsed.Seconds()
		fmt.Printf("%8d %8v %v\n", completed, elapsed, actual)
	}
}
