package main

import (
	"fmt"
	"log"
	"net/http"
)

func serveApp() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(resp, "Serv v1!")
	})
	if err := http.ListenAndServe("0.0.0.0:8080", mux); err != nil {
		log.Fatal(err) // !!! BUG: should not use log.Fatal
		// log.Fatal calls os.Exit which will unconditionally exit the program;
		// defers won’t be called, other goroutines won’t be notified to shut down,
		// the program will just stop. This makes it difficult to write tests for those functions.
	}
	// !!! BUG: serve shut down without stopping the application.
}

func serveDebug() {
	if err := http.ListenAndServe("127.0.0.1:8001", http.DefaultServeMux); err != nil {
		log.Fatal(err) // !!! BUG: should not use log.Fatal (is os.Exit() inside)
	}
	// !!! BUG: debug shut down without stopping the application.
}

func main() {
	go serveDebug()
	go serveApp()
	select {}
}

// !!!   Only use log.Fatal from main.main or init functions.  !!!
