package main
//go:generate go build -gcflags "-m -m" -o main2 main2.go
//go:generate rm main2
//
import (
	"fmt"
)

func main() {
	slice := []int{0, 1, 2, 3}
	println("slice",&slice)
	println("slice 0",&slice[0])
	println("slice 1",&slice[1])
	println("slice 2",&slice[2])
	println("slice 2",&slice[3])
	m := make(map[int]*int)
	println("m",&m)
	for key, val := range slice {
		println("key",&key,"val",&val)
		m[key] = &val
	}
	for k, v := range m {
		println("k",&k,"v",&v,"->",v)
		fmt.Println(k, "->", *v)
	}
}

// 问题：那么这个地址的内容一定会是3吗？
//
// 1. 观察逃逸分析的结果：
//./main2.go:12:11: val escapes to heap:
//./main2.go:12:11:   flow: {heap} = &val:
//./main2.go:12:11:     from &val (address-of) at ./main2.go:13:12
//./main2.go:12:11:     from m[key] = &val (assign) at ./main2.go:13:10 
// 说明：
// 因为13行处对于&val的引用，使得val逃逸到堆上。
// 即使key,val的作用域在for..range范围内，但是因为发生了逃逸。所以m中的只有
//
//
// 2. 输出分析：
//
// slice 0xc000068e10                                      slice在栈上
// slice 0 0xc000068dc8
// slice 1 0xc000068dd0
// slice 2 0xc000068dd8
// slice 2 0xc000068de0
//
// m 0xc000068df0                                          map在栈上
//
// key 0xc000068da8 val 0xc00001c088                       key在栈上，但是val因为逃逸，分配在堆上。
// key 0xc000068da8 val 0xc00001c088                       可以看出，key的地址总是相同，val的地址总是相同。
// key 0xc000068da8 val 0xc00001c088                       相同的key/value地址说明了for..range作用域的值copy
// key 0xc000068da8 val 0xc00001c088                       但是因为val有超出作用域的引用，所以被逃逸到堆
//
// k 0xc000068db0 v 0xc000068de8 -> 0xc00001c088           k在栈上，v也在栈上，但v中存着指向堆的地址（即val）
// 2 -> 3                                                  因为该堆内存一直被引用，所以不会GC，所以输出一定为3
// k 0xc000068db0 v 0xc000068de8 -> 0xc00001c088           相同的k/v地址再次说明了for..range作用域的值copy
// 3 -> 3
// k 0xc000068db0 v 0xc000068de8 -> 0xc00001c088
// 0 -> 3
// k 0xc000068db0 v 0xc000068de8 -> 0xc00001c088
// 1 -> 3