package main

import "fmt"

type IceCreamMaker interface {
	// Great a customer.
	Hello()
}

type Ben struct {
	name string
}

func (b *Ben) Hello() {
	if b.name != "Ben" {
		fmt.Printf("Ben says, Hello my name is %s\n", b.name)
	}
}

// the layout of the Ben and Jerry structs were identical in memory, so they were in some sense compatible.
type Jerry struct {
	field1 *[5]byte
	field2 int
}

// 如果我们修改Jerry的layout，运行一段时间后会panic
// 原因就在于这时候 Ben and Jerry2 的类型已经不再兼容了。
type Jerry2 struct {
	field2 int
	field1 *[5]byte
}
//
// Ben: 0xc000010240 Jerry: 0xc000010250
//   panic: runtime error: invalid memory address or nil pointer dereference
//   [signal SIGSEGV: segmentation violation code=0x1 addr=0x5 pc=0x1065253]
//
//   goroutine 1 [running]:
//   fmt.(*buffer).writeString(...)
//           /usr/local/go/src/fmt/print.go:82
//   fmt.(*fmt).padString(0xc00006ed40, 0x5, 0xc0000160b2)
//           /usr/local/go/src/fmt/format.go:110 +0x8e
//   fmt.(*fmt).fmtS(0xc00006ed40, 0x5, 0xc0000160b2)
//           /usr/local/go/src/fmt/format.go:359 +0x65
//   fmt.(*pp).fmtString(0xc00006ed00, 0x5, 0xc0000160b2, 0xc000000073)
//           /usr/local/go/src/fmt/print.go:446 +0x1ba
//   fmt.(*pp).printArg(0xc00006ed00, 0x10afdc0, 0xc000010270, 0xc000000073)
//           /usr/local/go/src/fmt/print.go:694 +0x875
//   fmt.(*pp).doPrintf(0xc00006ed00, 0x10cdd21, 0x1e, 0xc000068ec0, 0x1, 0x1)
//           /usr/local/go/src/fmt/print.go:1026 +0x168
//   fmt.Fprintf(0x10ea4a8, 0xc00000e018, 0x10cdd21, 0x1e, 0xc000068ec0, 0x1, 0x1, 0xc000068ee8, 0xc000068ef0, 0x0)
//           /usr/local/go/src/fmt/print.go:204 +0x72
//   fmt.Printf(...)
//           /usr/local/go/src/fmt/print.go:213

func (j *Jerry) Hello() {
	name := string((*j.field1)[:])
	if name != "Jerry" {
		fmt.Printf("Jerry says, Hello my name is %s \n", name)
	}
}

// from https://www.ardanlabs.com/blog/2014/06/ice-cream-makers-and-data-races-part-ii.html
// go run race_infv2.go | grep "Ben says, Hello my name is Jerry"
// go run race_infv2.go | grep "Jerry says, Hello my name is Ben"
func main() {
	var ben = &Ben{"Ben"}
	var jerry = &Jerry{field1:&[5]byte{'J', 'e', 'r', 'r', 'y'}, field2:5}
	var maker IceCreamMaker = ben

	var loop0, loop1 func()

	loop0 = func() {
		maker = ben
		go loop1()
	}

	loop1 = func() {
		maker = jerry
		go loop0()
	}
	fmt.Printf("Ben: %p Jerry: %p\n", ben, jerry)

	go loop0()

	for {
		maker.Hello()
	}
}
