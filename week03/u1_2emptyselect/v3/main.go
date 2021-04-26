package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "hello v3");
	})
	if err := http.ListenAndServe(":8080",nil); err != nil {
			log.Fatal(err)
	}
	//如果你的 goroutine 在从另一个 goroutine 获得结果之前无法取得进展，
	//那么通常情况下，你自己去做这项工作比委托它( go func() )更简单。
	//这通常消除了将结果从 goroutine 返回到其启动器所需的大量状态跟踪和 chan 操作。

	//问题！
	//如果要启动两个server怎么办？
}
