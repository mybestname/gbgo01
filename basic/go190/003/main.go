package main

import "fmt"

func main() {
	f1()
	f2()
}

func f1() {
	s := make([]int, 5)               // [0,0,0,0,0]  len=5, cap=5,
	s = append(s, 1, 2, 3)    // [0,0,0,0,0,1,2,3]  len=8,cap=10
	fmt.Println(s)
}

func f2() {
	s := make([]int, 0)                // [] cap=0, len=0
	s = append(s, 1, 2, 3, 4)  // [1,2,3,4] cap=4,len=4
	fmt.Println(s)
}

// E: https://golang.org/doc/effective_go#allocation_make
//   The built-in function make(T, args) serves a purpose different from new(T).
//     - It creates slices, maps, and channels only, and it returns an initialized (not zeroed) value of type T (not *T).
//     - The reason for the distinction is that these three types (slice, map, channel) represent, under the covers,
//       references to data structures that must be initialized before use.
//    A slice, for example, is a three-item descriptor containing a pointer to the data (inside an array), the length,
//    and the capacity, and until those items are initialized, the slice is nil.
//
// S: https://golang.org/ref/spec#Slice_types
//  - A new, initialized slice value for a given element type T is made using the built-in function make,
//    which takes a slice type and parameters specifying the length and optionally the capacity.
//  - A slice created with make always allocates a new, hidden array to which the returned slice value refers.
//  ```
//   make([]T, length, capacity)
//  ```
//  produces the same slice as allocating an array and slicing it, so these two expressions are equivalent:
//  ```
//    make([]int, 50, 100)
//    new([100]int)[0:50]
//  ```
// 关于slice的实现见main1.go
//
// 理解关键：
// cap才是涉及到alloc（内存/数组），而length/size只是slice如何看内存的问题。slice只是一个指向数据的胖指针，而不是数据本身。


