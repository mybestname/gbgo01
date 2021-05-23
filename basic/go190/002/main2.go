package main

import (
	"fmt"
)

func main() {
	var slice []int
	slice = append(slice,0, 1, 2, 3)
	m := make(map[int]*int)
	for key, val := range slice {
		m[key] = &val
	}
	for k, v := range m {
		fmt.Println(k, "->", *v)

	}
}


//go:generate go build -gcflags "-m" main2.go
//go:generate go tool objdump -s main.main -S main
//go:generate rm main

// 问题：那么这个地址的内容一定会是3吗？
// slice编译器知道大小，所以会被分配在栈上。
// 一定是3。这个指针指向的内存空间在栈上，直到main frame销毁前都存在
// 所以内容不会改变，一定是3。

// ./main.go:27:14: inlining call to fmt.Println
//./main.go:23:11: moved to heap: val
//./main.go:21:16: []int{...} does not escape
//./main.go:22:11: make(map[int]*int) does not escape
//./main.go:27:14: k escapes to heap
//./main.go:27:18: "->" escapes to heap
//./main.go:27:24: *v escapes to heap
//./main.go:27:14: []interface {}{...} does not escape
