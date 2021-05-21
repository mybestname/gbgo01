package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"golang.org/x/time/rate"
)

var _ rate.Limit


// https://www.reddit.com/r/golang/comments/fpbvxm/help_using_golangorgxtimerate/
// Q:
// I'm trying to use `golang.org/x/time/rate` to enforce a rate limit for making client side requests to an API.
// The API states they enforce a rate limit on their side of 10 requests for every 5 seconds.
// I can't seem to wrap my head around what I am doing wrong though because the values I'm using don't seem to
// benchmark to the rates I'm expecting.
// Code:
//  - https://github.com/ewohltman/go-discordbotsgg/blob/develop/pkg/discordbotsgg/discordbotsgg.go#L61
// Benchmark:
//  - https://github.com/ewohltman/go-discordbotsgg/blob/develop/pkg/discordbotsgg/discordbotsgg_test.go#L48
// Benchmark output:
// ```
//    goos: linux
//    goarch: amd64
//    pkg: github.com/ewohltman/go-discordbotsgg/pkg/discordbotsgg
//
//    BenchmarkClient_QueryBot
//
//    BenchmarkClient_QueryBot: discordbotsgg_test.go:66: Requests: 10, Seconds: 4.500560, RPS: 2.221946, Max RPS: 2.000000
//
//    BenchmarkClient_QueryBot: discordbotsgg_test.go:75: Failed to enforce rate limit
//
//    --- FAIL: BenchmarkClient_QueryBot
// ```
// A:
// From documentation (https://godoc.org/golang.org/x/time/rate#Limiter)
// > It implements a "token bucket" of size b, initially full and refilled at rate r tokens per second.
// > Informally, in any large enough time interval, the Limiter limits the rate to r tokens per second,
// > with a maximum burst size of b events.
//
// So you're limiting to 2 queries per second but since the bucket starts full you can get a single
// extra "bucket load" (1 for you) query in over entire run time. That, is you'll do 2/sec but if
// you start out going as fast as you can for n seconds you'll get 2*n+1 (since you use a bucket
// size of 1). For your 10 requests that means you do it quicker than 5 seconds or, in other words,
// if you ran for 5 seconds you'd do 10+1.
//
// a clearer example. It's more or less what you're doing. Note the possible "solution" for your test,
// by starting with an empty bucket; another would be to subtract off the bucket size from your measurement.
// Making your non-test code initialise the limit with an empty bucket isn't really a "solution" in
// the sense that if you were to do so, after any delay sufficient to refill the bucket you could once
// again have a brief period where you get an "extra" bucket-size worth of queries for a specific duration.
// Depending on the details of how the actual source API limit is enforced, this shouldn't be a real
// issue/problem.
//
// Edit: another way to look at your unexpected test results is that you're measuring the time between
// the 1st and the 10th query and failing since it's less than 5 seconds. But that's not a failure unless
// you did an 11th query in that 5 second window. E.g.:
//
// ```
//  1  2  3  4  5  6  7  8  9  10  11
//  ^--------------------------^     Less than five seconds, doesn't violate the limit
//  ^------------------------------^ if this was < five seconds it would be beyond the limit; equal to five "should" be okay
// ```

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
