package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)


func serve(addr string, handler http.Handler, stop <-chan struct{}) error {
	s := http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		<-stop // wait for stop signal
		fmt.Printf("stop signal received, shutdown serv(%v) gracefully.\n", addr)
		s.Shutdown(context.Background())  //shutdown会导致ListenAndServe()返回，从而结束serve()
	}()

	return s.ListenAndServe()
}

func serveApp(stop <-chan struct{}) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(resp, "serve v3!")
	})
	return serve("0.0.0.0:8080", mux, stop)
}

func serveDebug(stop <-chan struct{}) error {
	return serve("127.0.0.1:8001", http.DefaultServeMux, stop)
}

func main() {
	done := make(chan error, 2)
	stop := make(chan struct{})
	go func() {
		done <- serveDebug(stop)
	}()
	go func() {
		done <- serveApp(stop)
	}()

	go func() {
		interruptChannel := make(chan os.Signal, 1)
		interruptSignals := []os.Signal{os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM}
		signal.Notify(interruptChannel, interruptSignals...)
		select {
		case sig := <-interruptChannel:
			fmt.Printf("signal recevied: %v\n", sig)
			stop<-struct{}{}
		}
	}()

	var stopped bool
	for i := 0; i < cap(done); i++ {
		if err := <-done; err != nil {
			fmt.Printf("error: %v\n", err)
		}
		if !stopped {
			stopped = true
			close(stop)
		}
	}
}