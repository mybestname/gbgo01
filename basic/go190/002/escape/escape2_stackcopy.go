// original code from with modification
// https://www.ardanlabs.com/blog/2017/05/language-mechanics-on-escape-analysis.html
// Language Mechanics On Escape Analysis
// Sample program to show how stacks grow/change.
package main

// Number of elements to grow each stack frame.
// Run with 10 and then with 1024
const size = 1024

func main() {
	s := "HELLO"
	stackCopy(&s, 0, [size]int{})
	//output:
	//0 0xc0000eff60 HELLO
	//1 0xc0000eff60 HELLO
	//2 0xc0000fff60 HELLO
	//3 0xc0000fff60 HELLO
	//4 0xc0000fff60 HELLO
	//5 0xc0000fff60 HELLO
	//6 0xc00011ff60 HELLO                HELLO的地址改变了！因为递归调用时候，stack帧自增长来容纳更多
	//7 0xc00011ff60 HELLO
	//8 0xc00011ff60 HELLO
	//9 0xc00011ff60 HELLO
}
// stackCopy recursively runs increasing the size
// of the stack.
func stackCopy(s *string, c int, a [size]int) {
	println(c, s, *s)
	c++
	if c == 10 { return }
	stackCopy(s, c, a)
}

