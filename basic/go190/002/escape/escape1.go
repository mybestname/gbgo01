// original code from with modification
// https://www.ardanlabs.com/blog/2017/05/language-mechanics-on-stacks-and-pointers.html
// Language Mechanics On Stacks And Pointers
//
package main

// The passing of data between two frames is done “by value” in Go.
func main() {

   // Declare variable of type int with a value of 10.
   count := 10

   // Display the "value of" and "address of" count.
   println("count: Value [", count, "] Addr [", &count, "]")

   // Pass the "value of" the count.
   increment(count)

   println("count: Value [", count, "] Addr [", &count, "]")

   // Pass the "addr of" the count.
   increment2(&count)

   println("count: Value [", count, "] Addr [", &count, "]")
}
// count: Value [ 10 ] Addr [ 0xc00003e770 ]
// inc:   Value [ 11 ] Addr [ 0xc00003e760 ]                          // 因为pass-by-value，所以copy了一份指
// count: Value [ 10 ] Addr [ 0xc00003e770 ]
// inc2:  Value [ 0xc00003e770 ] Addr [ 0xc00003e760 ] PointsTo[ 11 ] // 因为pass-by-ref， 所以copy了一份地址
// count: Value [ 11 ] Addr [ 0xc00003e770 ]

// 使用`go:noinline`标签来阻止编译器进行内联编译的原因
// using the `go:noinline` directive to prevent the compiler from inlining the code for these functions
// directly in main. Inlining would erase the function calls and complicate this example.

//go:noinline
func increment(inc int) {
   // Increment the "value of" inc.
   inc++
   println("inc:   Value [", inc, "] Addr [", &inc, "]")
} // when the frame is taken, that the stack memory for that frame is wiped clean.


// * PointerTypes,
// 1.) PointerType = "*" BaseType .
// - https://golang.org/ref/spec#PointerType
// 2.) all have the same memory size and representation
//     On 32bit 4 bytes, 64bit 8 bytes.
// 3.) Pointer are variables like any other variable. They have a memory allocation and they hold a value.
//     It just so happens that all pointer variables, regardless of the type of value they can point to,
//     are always the same size and representation.
//go:noinline
func increment2(inc *int) {
    // Increment the "value of" count that the "pointer points to". (dereferencing)
    *inc++
    println("inc2:  Value [", inc, "] Addr [", &inc, "] PointsTo[", *inc, "]")
} // when the frame is taken, that the stack memory for that frame is wiped clean.


// Summary
//  - Functions execute within the scope of frame boundaries that provide an individual memory
//    space for each respective function.
//  - When a function is called, there is a transition that takes place between two frames.
//  - The benefit of passing data “by value” is readability.
//  - The stack is important because it provides the physical memory space for the frame
//    boundaries that are given to each individual function.
//  - All stack memory below the active frame is invalid but memory from the active frame and
//    above is valid.
//  - Making a function call means the goroutine needs to frame a new section of memory on the
//    stack.
//  - It’s during each function call, when the frame is taken, that the stack memory for that
//    frame is wiped clean.
//  - Pointers serve one purpose, to share a value with a function so the function can read and
//    write to that value even though the value does not exist directly inside its own frame.
//  - For every type that is declared, either by you or the language itself, you get for free a
//    compliment pointer type you can use for sharing.
//  - The pointer variable allows indirect memory access outside of the function’s frame that is
//    using it.
//  - Pointer variables are not special because they are variables like any other variable. They
//    have a memory allocation and they hold a value.