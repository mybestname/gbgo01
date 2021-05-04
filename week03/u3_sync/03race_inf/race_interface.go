package main

import "fmt"

type IceCreamMaker interface {
	// Hello greets a customer
	Hello()
}

type Ben struct {
	name string
}

func (b *Ben) Hello() {
	fmt.Printf("Ben says, Hello my name is %s\n", b.name)
}

type Jerry struct {
	name string
}

func (j *Jerry) Hello() {
	fmt.Printf("Jerry says, Hello my name is %s\n", j.name)
}

// from https://dave.cheney.net/2014/06/27/ice-cream-makers-and-data-races
// 执行如下命令可以发现文问题：
// go run race_interface.go |grep "Ben says, Hello my name is Jerry"
// go run race_interface.go |grep "Jerry says, Hello my name is Ben"
// 问题原因：Russ Cox的解释
// https://research.swtch.com/gorace
// > The current Go representation of slices and interface values admits a data race: because they are multiword values,
// > if one goroutine reads the value while another goroutine writes it, the reader might see half of the old value and
// > half of the new value.
// > In Go, an interface value is represented as two words, a type and a value of that type.
// > race happens in all of Go's mutable multiword structures: slices, interfaces, and strings.
func main() {
	var ben = &Ben{"Ben"}
	var jerry = &Jerry{"Jerry"}
	var maker IceCreamMaker = ben

	var loop0, loop1 func()

	loop0 = func() {
		maker = ben       // 赋值操作并不原子
		go loop1()
	}

	loop1 = func() {
		maker = jerry     // 赋值操作并不原子
		go loop0()
	}

	go loop0()

	for {
		maker.Hello()
	}
}

