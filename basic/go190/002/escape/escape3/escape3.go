// original code from with modification
// https://www.ardanlabs.com/blog/2017/06/language-mechanics-on-memory-profiling.html
// https://www.youtube.com/watch?v=2557w0qsDV0
// https://github.com/ardanlabs/gotraining/blob/master/topics/go/profiling/memcpu/stream.go
// Language Mechanics On Memory Profiling

package main

import (
	"bytes"
	"fmt"
	"io"
)

// Given a stream of bytes, write a function that can find the string elvis
// and replace it with the capitalized version of the string Elvis.

// data represents a table of input and expected output.
var data = []struct {
	input  []byte
	output []byte
}{
	{[]byte("abc"), []byte("abc")},
	{[]byte("elvis"), []byte("Elvis")},
	{[]byte("aElvis"), []byte("aElvis")},
	{[]byte("abcelvis"), []byte("abcElvis")},
	{[]byte("eelvis"), []byte("eElvis")},
	{[]byte("aelvis"), []byte("aElvis")},
	{[]byte("aabeeeelvis"), []byte("aabeeeElvis")},
	{[]byte("e l v i s"), []byte("e l v i s")},
	{[]byte("aa bb e l v i saa"), []byte("aa bb e l v i saa")},
	{[]byte(" elvi s"), []byte(" elvi s")},
	{[]byte("elvielvis"), []byte("elviElvis")},
	{[]byte("elvielvielviselvi1"), []byte("elvielviElviselvi1")},
	{[]byte("elvielviselvis"), []byte("elviElvisElvis")},
}

// assembleInputStream combines all the input into a
// single stream for processing.
func assembleInputStream() []byte {
	var in []byte
	for _, d := range data {
		in = append(in, d.input...)
	}
	return in
}

// assembleOutputStream combines all the output into a
// single stream for comparing.
func assembleOutputStream() []byte {
	var out []byte
	for _, d := range data {
		out = append(out, d.output...)
	}
	return out
}

func algOne(data []byte, find []byte, repl []byte, output *bytes.Buffer) {

	// Use a bytes Buffer to provide a stream to process.
	input := bytes.NewBuffer(data)
	// The number of bytes we are looking for.
	size := len(find)

	// Declare the buffers we need to process the stream.
	buf := make([]byte, size)
	end := size - 1

	// Read in an initial number of bytes we need to get started.
	if n, err := io.ReadFull(input, buf[:end]); err != nil {
		output.Write(buf[:n])
		return
	}

	for {
		// Read in one byte from the input stream.
		if _, err := io.ReadFull(input, buf[end:]); err != nil {
			// Flush the reset of the bytes we have.
			output.Write(buf[:end])
			return
		}
		// If we have a match, replace the bytes.
		if bytes.Compare(buf, find) == 0 {
			output.Write(repl)
			// Read a new initial number of bytes.
			if n, err := io.ReadFull(input, buf[:end]); err != nil {
				output.Write(buf[:n])
				return
			}
			continue
		}

		// Write the front byte since it has been compared.
		output.WriteByte(buf[0])
		// Slice that front byte out.
		copy(buf, buf[1:])
	}
}

// algTwo is a second way to solve the problem.
// Provided by Tyler Stillwater https://twitter.com/TylerStillwater
func algTwo(data []byte, find []byte, repl []byte, output *bytes.Buffer) {

	// Use the bytes Reader to provide a stream to process.
	input := bytes.NewReader(data)
	// The number of bytes we are looking for.
	size := len(find)

	// Create an index variable to match bytes.
	idx := 0
	for {
		// Read a single byte from our input.
		b, err := input.ReadByte()
		if err != nil {
			break
		}
		// Does this byte match the byte at this offset?
		if b == find[idx] {
			// It matches so increment the index position.
			idx++
			// If every byte has been matched, write
			// out the replacement.
			if idx == size {
				output.Write(repl)
				idx = 0
			}
			continue
		}
		// Did we have any sort of match on any given byte?
		if idx != 0 {
			// Write what we've matched up to this point.
			output.Write(find[:idx])
			// Unread the unmatched byte so it can be processed again.
			input.UnreadByte()
			// Reset the offset to start matching from the beginning.
			idx = 0
			continue
		}

		// There was no previous match. Write byte and reset.
		output.WriteByte(b)
		idx = 0
	}
}

func main() {
	var output bytes.Buffer
	in := assembleInputStream()
	out := assembleOutputStream()

	find := []byte("elvis")
	repl := []byte("Elvis")

	fmt.Println("=======================================\nRunning Algorithm One")
	output.Reset()
	algOne(in, find, repl, &output)
	matched := bytes.Compare(out, output.Bytes())
	fmt.Printf("Matched: %v\nInp: [%s]\nExp: [%s]\nGot: [%s]\n", matched == 0, in, out, output.Bytes())

	fmt.Println("=======================================\nRunning Algorithm Two")
	output.Reset()
	algTwo(in, find, repl, &output)
	matched = bytes.Compare(out, output.Bytes())
	fmt.Printf("Matched: %v\nInp: [%s]\nExp: [%s]\nGot: [%s]\n", matched == 0, in, out, output.Bytes())

}

//go:generate go test -run none -bench . -benchtime 3s -benchmem
// output :
// goos: darwin
// goarch: amd64
// pkg: escape3
// cpu: Intel(R) Core(TM) i7-7920HQ CPU @ 3.10GHz
// BenchmarkAlgorithmOne-8          1886578              1922 ns/op              53 B/op          2 allocs/op
// BenchmarkAlgorithmTwo-8          7578488               453.8 ns/op             0 B/op          0 allocs/op
// PASS
// ok      escape3 9.531s
//
// 注意bench输出中的allocs：
// 算法1中的 2 个 allocation ，以及 53 Byte 的 allocs容量

// 寻找算法1中哪里导致了内存allocs：
//go:generate  go test -run none -bench AlgorithmOne -benchtime 3s -benchmem -memprofile mem.out

// go tool pprof -alloc_space escape3.test mem.out
//
//      mem.out - the file contains the profile data
// escape3.test - a test binary contains symbols when looking at the profile data
//
// (pprof) list algOne
//      16MB   139.01MB (flat, cum) 99.64% of Total
//         .          .     56:}
//         .          .     57:
//         .          .     58:func algOne(data []byte, find []byte, repl []byte, output *bytes.Buffer) {
//         .          .     59:
//         .          .     60:   // Use a bytes Buffer to provide a stream to process.
//         .   123.01MB     61:   input := bytes.NewBuffer(data)
//         .          .     62:   // The number of bytes we are looking for.
//         .          .     63:   size := len(find)
//         .          .     64:
//         .          .     65:   // Declare the buffers we need to process the stream.
//      16MB       16MB     66:   buf := make([]byte, size)
//         .          .     67:   end := size - 1
//         .          .     68:
//         .          .     69:   // Read in an initial number of bytes we need to get started.
//         .          .     70:   if n, err := io.ReadFull(input, buf[:end]); err != nil {
//         .          .     71:           output.Write(buf[:n])
// 通过pprof的输出，可以发现：两处allocs
//
// 注意这句： input := bytes.NewBuffer(data)
//
// 返回指向Buffer的指针：*Buffer，这里是一个指针语义。input变成一个指向某个内存位置的指针。
// 但问题是为何input这里变成一个堆逃逸呢？
//
//go:generate go build -gcflags "-m=2"
//
// ./escape3.go:61:26: &bytes.Buffer{...} escapes to heap:
// ./escape3.go:61:26:   flow: ~R0 = &{storage for &bytes.Buffer{...}}:
// ./escape3.go:61:26:     from &bytes.Buffer{...} (spill) at ./escape3.go:61:26
// ./escape3.go:61:26:     from ~R0 = <N> (assign-pair) at ./escape3.go:61:26
// ./escape3.go:61:26:   flow: input = ~R0:
// ./escape3.go:61:26:     from input := (*bytes.Buffer)(~R0) (assign) at ./escape3.go:61:8
// ./escape3.go:61:26:   flow: io.r = input:
// ./escape3.go:61:26:     from input (interface-converted) at ./escape3.go:86:28
// ./escape3.go:61:26:     from io.r, io.buf := input, buf[:end] (assign-pair) at ./escape3.go:86:28
// ./escape3.go:61:26:   flow: {heap} = io.r:
// ./escape3.go:61:26:     from io.ReadAtLeast(io.r, io.buf, len(io.buf)) (call parameter) at ./escape3.go:86:28
//
// 观察逃逸分析：
// 果然input逃逸到heap
// 原因在86行：interface-converted 以及 call parameter 什么意思呢？
// - interface-converted 因为从 bytes.Buffer 转变为 io.Reader 接口
// - call parameter 因为 io.ReadAtLeast调用的方法签名是接口。
// 所以逃逸分析判断&bytes.Buffer{...} 需要逃逸。
// 思考🤔：
//  - 算法1是否真的需要使用io.ReadFull()呢？
//  - 有什么必要先转成一个接口呢？
//  - bytes.Buffer自己不能read吗？
//
// 第二处逃逸分析：
// ./escape3.go:66:13: make([]byte, size) escapes to heap:
// ./escape3.go:66:13:   flow: {heap} = &{storage for make([]byte, size)}:
// ./escape3.go:66:13:     from make([]byte, size) (non-constant size) at ./escape3.go:66:13
//
// 这个比较简单，因为使用make构造了一个大小不固定的slice，当然没法在栈上，肯定逃逸。
// 思考🤔：
//  - 这个slice的size必须是不固定的吗？不能是静态的吗？
//


func algOne_modified(data []byte, find []byte, repl []byte, output *bytes.Buffer) {

	// Use a bytes Buffer to provide a stream to process.
	input := bytes.NewBuffer(data)
	// The number of bytes we are looking for.
	// size := len(find)
	const size int = 10  //预设输入buff大小为10

	// Declare the buffers we need to process the stream.
	buf := make([]byte, size)   // 因为大小固定，所以不会逃逸。
	end := len(find) - 1
	if end > size {
		s := 2 * size
		for i := s ; s < end + 1 ; s = 2*i {}
		buf = make([]byte, s);  // 只有真正需要，才到heap上。
	}
	// Read in an initial number of bytes we need to get started.
	if n, err := input.Read(buf[:end]); err != nil {  // 使用byte.buffer自身，不会逃逸
		output.Write(buf[:n])
		return
	}
	for {
		// Read in one byte from the input stream.
		if _, err := input.Read(buf[end:]); err != nil {  // 使用byte.buffer自身，不会逃逸
			// Flush the reset of the bytes we have.
			output.Write(buf[:end])
			return
		}
		// If we have a match, replace the bytes.
		if bytes.Compare(buf, find) == 0 {
			output.Write(repl)
			// Read a new initial number of bytes.
			if n, err := input.Read(buf[:end]); err != nil { // 使用byte.buffer自身，不会逃逸
				output.Write(buf[:n])
				return
			}
			continue
		}

		// Write the front byte since it has been compared.
		output.WriteByte(buf[0])
		// Slice that front byte out.
		copy(buf, buf[1:])
	}
}

//go:generate  go test -run none -bench AlgorithmOneModified -benchtime 3s -benchmem -memprofile mem.out
//
// goos: darwin
// goarch: amd64
// pkg: escape3
// cpu: Intel(R) Core(TM) i7-7920HQ CPU @ 3.10GHz
// BenchmarkAlgorithmOneModified-8   	10012016	       350.9 ns/op	       0 B/op	       0 allocs/op
//
// 可见已经没有没有堆内存的alloc，同时性能大幅提高。
//
//
//
//
//
// ！！遗留问题！！
//
// ============================================ TL/DR ===========================================
// Q ：逃逸分析为何是86行输出？70和77行的调用为何没有输出呢？第70行不是完全一样的代码吗？难道原因不是第70行导致的吗？
//
// A : 这是因为go的逃逸分析内部逻辑的优化所致：
// ============================================ TL/DR ===========================================
//
// go1.13之前，执行的结果确实会定位到70行，但是从go1.13开始，添加了新版的逃逸分析 (可以添加"-newescape=false"改变1.13默认行为，继续在1.13上使用旧版逃逸分析，但只在1.13版本有效，1.14完全删除了旧版)
// 同时从1.14beta1开始，逃逸分析的行为有了比较大的改变：
//  - https://github.com/golang/go/issues/23109 (rewrite escape analysis)
//  - https://github.com/golang/go/commit/4cde749f633753cf59d0cfc912351e1b1def2b4f  (restore more missing -m=2 escape analysis details)
//  - https://github.com/golang/go/commit/991b0fd46c3e8160c9b5c622478caf7b5ebe139c  (remove  -newescape flag 1.14beta1)
//  - https://github.com/golang/go/commit/de454eef5f47212dc8a9d9c2c8b598fa343d2c2b  (rsc修改目录结构)
//    + 最新版trunk的代码在：go/src/cmd/compile/internal/escape/escape.go
// 版本比较：
// 1.13
//  - https://github.com/golang/go/blob/go1.13.15/src/cmd/compile/internal/gc/esc.go
//  - vs.
//  - https://github.com/golang/go/blob/go1.13.15/src/cmd/compile/internal/gc/escape.go
// 1.14beta1
//  - https://github.com/golang/go/blob/go1.14beta1/src/cmd/compile/internal/gc/esc.go
//  - https://github.com/golang/go/blob/go1.14beta1/src/cmd/compile/internal/gc/escape.go
// 1.14
//  - https://github.com/golang/go/blob/go1.14/src/cmd/compile/internal/gc/esc.go
//  - https://github.com/golang/go/blob/go1.14/src/cmd/compile/internal/gc/escape.go
// 1.16.4
//  - https://github.com/golang/go/blob/go1.16.4/src/cmd/compile/internal/gc/esc.go
//  - https://github.com/golang/go/blob/go1.16.4/src/cmd/compile/internal/gc/escape.go
// Trunk
//  - https://github.com/golang/go/blob/master/src/cmd/compile/internal/escape/escape.go
//
// 70和77不被显示的根源：在于：
// mdempsky ( Sep 4, 2019 , go1.14beta1)
//
// cmd/compile: silence esc diagnostics about directiface OCONVIFACEs
//
// In general, a conversion to interface type may require values to be
// boxed, which in turn necessitates escape analysis to determine whether
// the boxed representation can be stack allocated.
//
// However, esc.go used to unconditionally print escape analysis
// decisions about OCONVIFACE, even for conversions that don't require
// boxing (e.g., pointers, channels, maps, functions).
//
// For test compatibility with esc.go, escape.go similarly printed these
// useless diagnostics. This CL removes the diagnostics, and updates test
// expectations accordingly.
//
// https://github.com/golang/go/commit/9f89edcd9668bb3b011961fbcdd8fc2796acba5d#diff-9325f8b93fd0e344d1334adf9bacf53849fe22223330fd4900ea8e6bc9a242b5
//
// 在trunk源上rollback这个提交，得到的输出如下：
//
// escape3.go:61:26: &bytes.Buffer{...} escapes to heap:
// escape3.go:61:26:   flow: ~R0 = &{storage for &bytes.Buffer{...}}:
// escape3.go:61:26:     from &bytes.Buffer{...} (spill) at escape3.go:61:26
// escape3.go:61:26:     from ~R0 = &bytes.Buffer{...} (assign-pair) at escape3.go:61:26
// escape3.go:61:26:   flow: input = ~R0:
// escape3.go:61:26:     from input := (*bytes.Buffer)(~R0) (assign) at escape3.go:61:8
// escape3.go:61:26:   flow: io.r = input:
// escape3.go:61:26:     from input (interface-converted) at escape3.go:86:28
// escape3.go:61:26:     from io.r, io.buf := input, buf[:end] (assign-pair) at escape3.go:86:28
// escape3.go:61:26:   flow: {heap} = io.r:
// escape3.go:61:26:     from io.ReadAtLeast(io.r, io.buf, len(io.buf)) (call parameter) at escape3.go:86:28
// 上述内容和未改前输出完全一样，后续输出的内容如下：
// escape3.go:86:28: input escapes to heap:
// escape3.go:86:28:   flow: io.r = &{storage for input}:
// escape3.go:86:28:     from input ( spill) at escape3.go:86:28
// escape3.go:86:28:     from io.r, io.buf := input, buf[:end] (assign-pair) at escape3.go:86:28
// escape3.go:86:28:   flow: {heap} = io.r:
// escape3.go:86:28:     from io.ReadAtLeast(io.r, io.buf, len(io.buf)) (call parameter) at escape3.go:86:28
// escape3.go:77:27: input escapes to heap:
// escape3.go:77:27:   flow: io.r = &{storage for input}:
// escape3.go:77:27:     from input (spill) at escape3.go:77:27
// escape3.go:77:27:     from io.r, io.buf := input, buf[end:] (assign-pair) at escape3.go:77:27
// escape3.go:77:27:   flow: {heap} = io.r:
// escape3.go:77:27:     from io.ReadAtLeast(io.r, io.buf, len(io.buf)) (call parameter) at escape3.go:77:27
// escape3.go:70:26: input escapes to heap:
// escape3.go:70:26:   flow: io.r = &{storage for input}:
// escape3.go:70:26:     from input (spill) at escape3.go:70:26
// escape3.go:70:26:     from io.r, io.buf := input, buf[:end] (assign-pair) at escape3.go:70:26
// escape3.go:70:26:   flow: {heap} = io.r:
// escape3.go:70:26:     from io.ReadAtLeast(io.r, io.buf, len(io.buf)) (call parameter) at escape3.go:70:26
//
// 那么回到这个问题，即使让9f89edcd9668bb3b011961fbcdd8fc2796acba5d的修改，那么为何还是定位在86行，而不是第70行呢？
// 输出的顺序到底是怎么决定的呢？这是因为新版逃逸分析的大致思路是：
// 0. inlining的输出不在逃逸分析代码内：忽略。
// 1. 首先分析一定会逃逸的： HeapAllocReason()，并给出理由。
//    - 例如：make并且没有指定容量的slice，一定是heap。
// 2. 然后wall-all (把位置加入一个LIFO队列) 这个队列顺序是输出的关键。
//    - https://github.com/golang/go/commit/867ea9c17f1031e27f4a2d17d9f7dfb270f73fa1 (mdempsky  Sep 26, 2019)
//      - cmd/compile: use proper work queue for escape graph walking
//    - 这个提交使得wall-all的顺序有了大改变。
//      - https://github.com/golang/go/blob/go1.14/src/cmd/compile/internal/gc/escape.go#L1113-L1144
//    - 修改的目的是优化逃逸分析算法，使得算法更快（目的是提高编译速度）
//    - 但是优化的结果就导致从输出来说，目的并不是更好有利于人阅读来改进代码。
//
// ==============
// 总结：
// ==============
//   - 逃逸分析的行为随着go版本的升级和优化，输出的结果可能有比较大的不同
//   - 所谓的逃逸分析的-m选项是显示go在编译时候的优化过程。所以这个开关的第一目的并不是给程序员检查代码用的。
//     逃逸分析的第一服务对象肯定是编译器本身，而不是为了程序员检查代码，所以肯定优化为先。
//   - 从程序员分析代码的角度出发，可以考虑使用go1.12.x或者go1.13.x（使用-newescape=false)进行辅助
//     逃逸分析显示，旧版的输出从程序员角度可能更容易阅读，但是需要注意的是：旧版因为没有很多后续版本的优化，
//     所以结果可能相互并不一致（或者说必然不一致，因为毕竟是不同版本的编译器），例如本来已经被新版编译器
//     优化为非逃逸的，旧版还会编译为逃逸。