// original code from with modification
// https://www.ardanlabs.com/blog/2017/05/language-mechanics-on-escape-analysis.html
// Language Mechanics On Escape Analysis
package main

type user struct {
    name  string
    email string
}

func main() {
    u1 := createUserV1()
    u2 := createUserV2()
	u3 := createUserV3()

	println("u1", &u1)
	println("u2", &u2, "->",u2)
	println("u3", &u3, "->",u3)
}
// output:
// V1 0xc00003e6f8
// V2 0xc00008e000
// V3 0xc00008e020
// u1 0xc00003e758
// u2 0xc00003e750 -> 0xc00008e000
// u3 0xc00003e748 -> 0xc00008e020
//
// 分析：
//                         0xc00008e020  heap user    <---------------+----+
//                         0xc00008e000  heap user    <-----+----+    |    |
//                         ....                             |    |    |    |
//  main frame u1 -->      0xc00003e758  stack user         |    |    |    |  u1是复制到栈上的。
//  main frame u2 -->      0xc00003e750  c00008e000   ------^    |    |    |  u2引用堆上的
//  main frame u3 -->      0xc00003e748  c00008e020   -----------|----^    |  u3引用堆上的
//                         ....                                  |         |
//    V1 frame    -->      0xc00003e6f8  stack user              |         |   v1这个在栈上，不要逃逸，运行时调用完毕后即释放
//    v2 frame    -----------------------------------------------^         |   v2的这个在编译器时候逃逸分析，得出结论需要放到堆上
//    v3 frame    ---------------------------------------------------------^   v3的这个在编译器时候逃逸分析，得出结论需要放到堆上
//
// 注意在stack分配和在heap上分配的区别。


//go:noinline
func createUserV1() user {
	u := user{
		name:  "Bill",
		email: "bill@ardanlabs.com",
	}

	println("V1", &u)
	return u                         // the user value created by this function is being copied and passed
	                                 // up the call stack. This means the calling function is receiving a
	                                 // copy of the value itself.
	                                 // 因为是copy方式共享，所以被copy到main的栈上。而V1帧直接释放
}

//go:noinline
func createUserV2() *user {
	u := user{                       // moved to heap: u
		name:  "Bill",
		email: "bill@ardanlabs.com",
	}

	println("V2", &u)
	return &u                       // 明确表示返回指针，可读性更好。
	                                // 编译器会把该对象分配到堆上，而不是栈上。否则方法结束之后，局部变量就被回收，所以分配到堆上是理所当然的
	                                // 这种在C上就是野指针（悬垂指针），表示危险
	                                // 而对于golang来说，通过逃逸分析，给你分配到了堆上来保证使用。
}

//go:noinline
func createUserV3() *user {          // V3和V2的语法不同，1. V2显然可读性更好 2.和V2有差别，V3会多一个栈分配再copy（逃逸）的过程。
	u := &user{                      // escapes to heap 和 move to heap 是有区别的。
		name:  "Bill",
		email: "bill@ardanlabs.com",
	}

	println("V3", u)
	return u
}


// Anytime you share a value up the call stack, it is going to escape.
// 逃逸分析是编译器用于决定变量分配到堆上还是栈上的一种行为
// 在编译阶段确立逃逸，注意并不是在运行时。
// Go语言里没有一个关键字或者函数可以直接让变量被编译器分配到堆上，相反，编译器通过分析代码来决定将变量分配到何处。
// 对一个变量取地址，可能会被分配到堆上。但是编译器进行逃逸分析后，如果考察到在函数返回后，此变量不会被引用，
// 那么还是会被分配到栈上。
// 简单来说，编译器会根据变量是否被外部引用来决定是否逃逸：
// - 如果函数外部没有引用，则优先放到栈中；
// - 如果函数外部存在引用，则必定放到堆中；
// Go 语言的逃逸分析遵循以下两个不变性：
// - 指向栈对象的指针不能存在于堆中；
// - 指向栈对象的指针不能在栈对象回收后存活；
//

//go:generate go build -gcflags "-m -m" escape2.go
//go:generate rm escape2
//
//
//./escape2.go:44:6: cannot inline createUserV1: marked go:noinline
//./escape2.go:58:6: cannot inline createUserV2: marked go:noinline
//./escape2.go:72:6: cannot inline createUserV3: marked go:noinline
//./escape2.go:11:6: cannot inline main: function too complex: cost 205 exceeds budget 80
//./escape2.go:59:2: u escapes to heap:
//./escape2.go:59:2:   flow: ~r0 = &u:
//./escape2.go:59:2:     from &u (address-of) at ./escape2.go:65:9
//./escape2.go:59:2:     from return &u (return) at ./escape2.go:65:2
//./escape2.go:59:2: moved to heap: u
//./escape2.go:73:7: &user{...} escapes to heap:
//./escape2.go:73:7:   flow: u = &{storage for &user{...}}:
//./escape2.go:73:7:     from &user{...} (spill) at ./escape2.go:73:7
//./escape2.go:73:7:     from u := &user{...} (assign) at ./escape2.go:73:4
//./escape2.go:73:7:   flow: ~r0 = u:
//./escape2.go:73:7:     from return u (return) at ./escape2.go:79:2
//./escape2.go:73:7: &user{...} escapes to heap
//
// 逃逸分析输出可以看出:
// - 变量u和&user{...}都逃逸到堆上
// - literal语法的不同导致输出不同
// - move to heap 和 escape to hea是有不同的。应该使用V2的语义，而非V3的写法。


//
// 所以 一共有关两个user逃逸到heap上。（需要GC才能回收）
// 而值copy的是栈上拷贝，不需要GC回收。但是问题是需要copy副本。
// 而如果要共享数据，但是不要副本，那么就得用heap，只能堆上。但是就需要GC才能回收了。
// 到处都用指针传递并不一定是最好的。 那么到底那种更好，更有效率？需要具体问题具体分析。
// 不要盲目使用变量的指针作为函数参数，虽然它会减少复制操作。但其实当参数为变量自身的时候，复制是在栈上完成的操作，
// 开销远比变量逃逸后动态地在堆上分配内存少的多。
//
// 值的构建并不能决定它的所在，只有值的共享方式才能确定编译器将如何安排它。
// 并不是说所有指针对象，都应该在堆上，而是
// 每当你在调用栈中共享一个值时，它都会逃逸。
// 需要在值语义与指针语义之间权衡得失，
// - 值语义将值保留在栈(stack)上从而减轻了GC的压力，但必须存储、追踪和维护任何给定的值的副本
// - 而指针语义则将值放置在堆(heap)上，对GC造成了压力，但只有一个值需要被存储、追踪、维护，所以效率很高。
// 关键在于正确、一致、平衡地使用每种语义。
// 值语义和指针语义
// function is using value semantics on the return
//  - the user value created by this function is being copied and passed up the call stack.
//  - This means the calling function is receiving a copy of the value itself.
// function is using pointer semantics on the return
//  - the user value created by this function is being shared up the call stack.
//  - This means the calling function is receiving a copy of the address for the value.
//