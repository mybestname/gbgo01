package main

import "fmt"

func main() {
	slice := []int{0, 1, 2, 3}
	m := make(map[int]*int)
	for i:=0; i < len(slice); i++ { // 1. don't use for .. range to avoid value copy.
		m[i] = &slice[i]  // get addr from the source data.
	}
	for _, v := range slice {       // 2. range form `slice` for the fixed order.
		fmt.Println(v, "->", m[v], "->",*m[v])
	}
}

// range slice 可以保证有序
// S： https://golang.org/ref/spec#For_statements
// 1. For an array, pointer to array, or slice value a, the index iteration values
//    are produced in increasing order, starting at element index 0.
//    - If at most one iteration variable is present, the range loop produces iteration
//      values from 0 up to len(a)-1 and does not index into the array or slice itself.
//    - For a nil slice, the number of iterations is 0.