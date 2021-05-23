package main

import (
	"fmt"
)

// 关于 `for ... range` 循环
// S:
// A "for" statement with a "range" clause iterates through all entries of an array, slice, string or map,
// or values received on a channel. For each entry it assigns iteration values to corresponding iteration
// variables if present and then executes the block.
// The iteration variables may be declared by the "range" clause using a form of short variable declaration (:=).
// In this case their types are set to the types of the respective iteration values and their scope is the block
// of the "for" statement; they are re-used in each iteration.
// If the iteration variables are declared outside the "for" statement, after execution their values will be
// those of the last iteration.
// E:
//

func main() {
	slice := []int{0, 1, 2, 3}              // slice 使用literal初始化，len=4， cap=4， [0,1,2,3]
	m := make(map[int]*int)                 // 使用 make 初始化一个 int -> *int 的 map，没有指定任何其它参数，那么默认是空
	for key, val := range slice {           // for range 遍历 slice，把key为 0， val值为 0
		m[key] = &val                       // m[0] -> addr of val, &val的地址是val的地址，但是不是slice的地址
	}
	for k, v := range m {                   // 遍历 m
		fmt.Println(k, "->", *v)            // 0 -> value from addr of val, 因为都是一个相同的地址，所以一定都是一样的值。
		                                    // 问题：那么这个地址的内容一定会是3吗？
	}
}
// output :
// 1 -> 3
// 2 -> 3
// 3 -> 3
// 4 -> 3
//
// !!! 注意 !!! 这只是一种可能性，因为对于map的range遍历的顺序是不确定的。
// 也可能是：
// 2 -> 3
// 3 -> 3
// 0 -> 3
// 1 -> 3
//
// 参考 https://blog.golang.org/maps#TOC_7.
// Go maps in action
//  When iterating over a map with a range loop, the iteration order is not specified and is not
//  guaranteed to be the same from one iteration to the next. If you require a stable iteration
//  order you must maintain a separate data structure that specifies that order.
//
// ```
//    import "sort"
//
//    var m map[int]string
//    var keys []int
//    for k := range m {
//        keys = append(keys, k)
//    }
//    sort.Ints(keys)
//    for _, k := range keys {
//        fmt.Println("Key:", k, "Value:", m[k])
//    }
// ```
// 另见规范：S: https://golang.org/ref/spec#For_statements
//   3. The iteration order over maps is not specified and is not guaranteed to be the same from
//      one iteration to the next.
//        - If a map entry that has not yet been reached is removed during iteration, the
//          corresponding iteration value will not be produced.
//        - If a map entry is created during iteration, that entry may be produced during
//          the iteration or may be skipped.
//        - The choice may vary for each entry created and from one iteration to the next.
//        - If the map is nil, the number of iterations is 0.
//
// 附注：对于go1.12以上版本，为了方便测试，对于使用fmt.Println(map)这种模式，将结果排序后的结果
// https://tip.golang.org/doc/go1.12#fmt
// ```
//   func main() {
//       m := map[int]int{3: 5, 2: 4, 1: 3}
//       fmt.Println(m)
//
//       // In Go 1.12+
//       // Output: map[1:3 2:4 3:5]
//
//       // Before Go 1.12 (the order was undefined)
//       // map[3:5 2:4 1:3]
//   }
// ```
// 注意：这个只是对于fmt.Println如何打印一个map的情况，而不是for..range循环的遍历的情况，所以原代码的输出
// 还是不确定的，不管是否使用了go.12以上版本。
//

//go:generate go build -gcflags "-m" main.go
//go:generate go tool objdump -s main.main -S main
//go:generate rm main

// 问题：那么这个地址的内容一定会是3吗？
// ?

