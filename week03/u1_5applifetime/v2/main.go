package main

import (
	"fmt"
	"net/http"
)

func serveApp() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(resp, "Serve v2!")
	})
	return http.ListenAndServe("0.0.0.0:8080", mux)
}

func serveDebug() error {
	return http.ListenAndServe("127.0.0.1:8001", http.DefaultServeMux)
}

func main() {
	done := make(chan error, 2)
	go func() {
		done <- serveDebug()
	}()
	go func() {
		done <- serveApp()
	}()

	for i := 0; i < cap(done); i++ {
		if err := <-done; err != nil {
			fmt.Println("error: %v", err)
		}
	}
}

// 问题，不能优雅的stop，无法优雅的调用shutdown，只是处理了错误的情况。
