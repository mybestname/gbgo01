package main

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
)

// errgroup的源码非常简单（含注释才67行），就是对waitGroup和context的使用。
// 对比直接使用waitGroup
// - 不用手工defer wg.Done()
// - 加入cancel的支持


func main() {
	g := new(errgroup.Group)
	var urls = []string{
		"http://www.golang.org/",
		"http://www.google.com/",
		"http://www.somestupidname.com/",
	}
	for _, url := range urls {
		// Launch a goroutine to fetch the URL.
		url := url // https://golang.org/doc/faq#closures_and_goroutines
		g.Go(func() error {
			// Fetch the URL.
			resp, err := http.Get(url)
			if err == nil {
				resp.Body.Close()
			}
			return err
		})
	}
	// Wait for all HTTP fetches to complete.
	// Wait blocks until all function calls from the Go method have returned,
	// then returns the first non-nil error (if any) from them.
	if err := g.Wait(); err == nil {
		fmt.Println("Successfully fetched all URLs.")
	}
}


