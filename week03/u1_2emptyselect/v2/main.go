package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "hello");
	})
	go func() {
		if err := http.ListenAndServe(":8080",nil); err != nil {
			log.Fatal(err)
		}
		fmt.Println("goroutine exit")
	}()
	// 空的 select 语句将永远阻塞。
	select {}
}
