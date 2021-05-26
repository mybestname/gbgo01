package main

import (
	"fmt"
	"testing"
)

func ExampleMakeSlice() {
	var s []int
	fmt.Println(s, len(s), cap(s))
	s = make([]int,0)
	fmt.Println(s, len(s), cap(s))
	s = []int{}
	fmt.Println(s, len(s), cap(s))
	s = []int{0}
	fmt.Println(s, len(s), cap(s))
	s = make([]int, 3)               //如果省略第二个参数, 表示第二个等于前者，即cap=3
	fmt.Println(s, len(s), cap(s))
	s = append(s, 0)         //如果省略第二个参数，表示第二个等于前者，即cap=0
	fmt.Println(s, len(s), cap(s))
	s = make([]int, 0, 3)            //第一个参数是长度，第二个参数是capacity
	fmt.Println(s, len(s), cap(s))
	s = make([]int, 1, 3)
	fmt.Println(s, len(s), cap(s))
	s = make([]int, 1, 3)
	//output:
	// [] 0 0
	// [] 0 0
	// [] 0 0
	// [0] 1 1
	// [0 0 0] 3 3
	// [0 0 0 0] 4 6
	// [] 0 3
	// [0] 1 3
}

func TestUnInitializeSlice(t *testing.T) {
	defer func() {
		if e := recover(); e!= nil {
			fmt.Println(e)  //panic runtime error:
		}
	}()
	var s = make([]int,1)
	fmt.Println(s)
	s[0] = 2
	fmt.Println(s)
	var p = new([]int)
	(*p)[0]=2  //runtime painc
	t.Fail()   //should not go there
}