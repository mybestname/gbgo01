// original code from with modification
// http://www.agardner.me/golang/garbage/collection/gc/escape/analysis/2015/10/18/go-escape-analysis.html
// Golang escape analysis
//
package main

// The basic rule is
//  if a reference to a variable is returned from the function where it is declared, it “escapes”
//  it can be referenced after the function returns, so it must be heap-allocated.
//
// - functions calling other functions
// - references being assigned to struct members
// - slices and maps
// - cgo taking pointers to variables
//
// How to
//   Go builds a graph of function calls at compile time, and traces the flow of input arguments and
//   return values.
//    - A function may take a reference to one of it’s arguments, but if that reference is not returned,
//   the variable does not escape.
//    - A function may also return a reference, but that reference may be dereferenced or not returned
//   by another function in the stack before the function which declared the variable returns.

type S struct {}

func v1() {
	var x S
	_ = identityV1(x) // not escape
}

// 全部都是pass-by-value，不会逃逸
func identityV1(x S) S {
	return x
}

func v2() {
	var x S
	y := &x
	_ = *identityV2(y)
}

//z只是流过，没有任何reference,所以还是不会逃逸
func identityV2(s *S) *S {
	return s
}

func v3() {
	var x S
	_ = *ref(x)
	_ = ref2()
	_ = ref3()
}

// ref发生逃逸
// 入参 pass-by-value 语义
func ref(s S) *S { //入参s被逃逸
	return &s //返回参数s的地址，不可能在`ref`的栈上，只能把s分配在堆上，直接返回堆地址。-> moved to heap
}

// ref2发生逃逸
func ref2() *S {
	s := S{}  //试图声明本地变量z（试图栈上分配）// s被逃逸
	return &s //返回本地变量z的地址，不可能在`ref2`的栈上，只能把z分配在堆上，直接返回堆地址。-> move to heap
}

// ref3发生逃逸
func ref3() *S {
	s := &S{} //试图用z存储栈上内存的地址  // &S{}被逃逸
	return s // 返回本地变量z，不可能在`ref3`的栈上，只能把栈上内容复制到堆上，然后返回堆上地址。 -> escape to heap
}

type Z struct {
	M *int
}

func v4() {
	var i int
	refStruct(i)
}

// refStruct发生逃逸
// 入参 pass-by-copy，返回值语义，表面上应该没有问题。但入参y会逃逸
func refStruct(y int) (z Z) { //入参y会逃逸
	z.M = &y   // Z的成员是一个引用，而这个引用在refStruct中被指定。
	return z   // z被返回，那么对y的引用必须不能在栈中，必须在堆上分配，所以y必须逃逸。
}

func v5() {
	var i int  // 这个i不会逃逸，
	refStructV2(&i)
}
// refStructV2不会发生逃逸
func refStructV2(y *int) (z Z) { //y不会逃逸 （leaking param: y to result z level=0，调用v5者需要考虑）
	z.M = y    //引用被传递进来，那么只要y本身的生存周期ok，不需要给y在堆上分配。
	return z
}

// v6/v7/v8 涉及的都是逃逸分析的粒度
//
func v6() {
	var z Z
	var i int    // 这个i会逃逸，其实是不用逃逸的。但是go的逃逸分析器能力有限。
	refStructV3(&i, &z)
}

//不发生逃逸
func refStructV3(y *int, z *Z) { // z不逃逸， 入参y泄漏：leaking param: y（调用者v6需要重点考虑）
	z.M = y
}

func v7() {
	i1 := 1                   // i1 和 i2 都会被逃逸到堆上
	i2 := 2                   // 但其实是不需要的，但是go无法分析
	z1 := Z{&i1}
	z2 := Z{&i2}
	refStructV4(&z1, &z2)
}

func refStructV4(z1 *Z, z2 *Z) { // z1和z2都泄漏，调用者v7需要考虑？
	temp := z1.M
	z1.M = z2.M
	z2.M = temp
}

func v8() {   //对比v7和v8
	i := 1
	i2 := 2
	z1 := Z{&i}
	z2 := Z{&i2}
	temp := z1.M
	z1.M = z2.M
	z2.M = temp
}

//go:generate go tool compile -m=2 -l escape5.go
////go:generate go tool objdump escape5.o
//go:generate rm escape5.o
//
// escape5.go:43:17: parameter s leaks to ~r1 with derefs=0:
// escape5.go:43:17:   flow: ~r1 = s:
// escape5.go:43:17:     from return s (return) at escape5.go:44:2
// escape5.go:43:17: leaking param: s to result ~r1 level=0
// escape5.go:56:10: parameter s leaks to ~r1 with derefs=0:
// escape5.go:56:10:   flow: ~r1 = &s:
// escape5.go:56:10:     from &s (address-of) at escape5.go:57:9
// escape5.go:56:10:     from return &s (return) at escape5.go:57:2
// escape5.go:56:10: s escapes to heap:
// escape5.go:56:10:   flow: ~r1 = &s:
// escape5.go:56:10:     from &s (address-of) at escape5.go:57:9
// escape5.go:56:10:     from return &s (return) at escape5.go:57:2
// escape5.go:56:10: moved to heap: s
// escape5.go:62:2: s escapes to heap:
// escape5.go:62:2:   flow: ~r0 = &s:
// escape5.go:62:2:     from &s (address-of) at escape5.go:63:9
// escape5.go:62:2:     from return &s (return) at escape5.go:63:2
// escape5.go:62:2: moved to heap: s
// escape5.go:68:7: &S{} escapes to heap:
// escape5.go:68:7:   flow: s = &{storage for &S{}}:
// escape5.go:68:7:     from &S{} (spill) at escape5.go:68:7
// escape5.go:68:7:     from s := &S{} (assign) at escape5.go:68:4
// escape5.go:68:7:   flow: ~r0 = s:
// escape5.go:68:7:     from return s (return) at escape5.go:69:2
// escape5.go:68:7: &S{} escapes to heap
// escape5.go:83:16: parameter y leaks to z with derefs=0:
// escape5.go:83:16:   flow: z = &y:
// escape5.go:83:16:     from &y (address-of) at escape5.go:84:8
// escape5.go:83:16:     from z.M = &y (assign) at escape5.go:84:6
// escape5.go:83:16: y escapes to heap:
// escape5.go:83:16:   flow: z = &y:
// escape5.go:83:16:     from &y (address-of) at escape5.go:84:8
// escape5.go:83:16:     from z.M = &y (assign) at escape5.go:84:6
// escape5.go:83:16: moved to heap: y
// escape5.go:93:18: parameter y leaks to z with derefs=0:
// escape5.go:93:18:   flow: z = y:
// escape5.go:93:18:     from z.M = y (assign) at escape5.go:94:6
// escape5.go:93:18: leaking param: y to result z level=0
// escape5.go:107:18: parameter y leaks to {heap} with derefs=0:
// escape5.go:107:18:   flow: {heap} = y:
// escape5.go:107:18:     from z.M = y (assign) at escape5.go:108:6
// escape5.go:107:18: leaking param: y
// escape5.go:107:26: z does not escape
// escape5.go:102:6: i escapes to heap:
// escape5.go:102:6:   flow: {heap} = &i:
// escape5.go:102:6:     from &i (address-of) at escape5.go:103:14
// escape5.go:102:6:     from refStructV3(&i, &z) (call parameter) at escape5.go:103:13
// escape5.go:102:6: moved to heap: i
// escape5.go:119:18: parameter z1 leaks to {heap} with derefs=1:
// escape5.go:119:18:   flow: temp = *z1:
// escape5.go:119:18:     from z1.M (dot of pointer) at escape5.go:120:12
// escape5.go:119:18:     from temp := z1.M (assign) at escape5.go:120:7
// escape5.go:119:18:   flow: {heap} = temp:
// escape5.go:119:18:     from z2.M = temp (assign) at escape5.go:122:7
// escape5.go:119:25: parameter z2 leaks to {heap} with derefs=1:
// escape5.go:119:25:   flow: {heap} = *z2:
// escape5.go:119:25:     from z2.M (dot of pointer) at escape5.go:121:11
// escape5.go:119:25:     from z1.M = z2.M (assign) at escape5.go:121:7
// escape5.go:119:18: leaking param content: z1
// escape5.go:119:25: leaking param content: z2
// escape5.go:113:2: i2 escapes to heap:
// escape5.go:113:2:   flow: z2 = &i2:
// escape5.go:113:2:     from &i2 (address-of) at escape5.go:115:10
// escape5.go:113:2:     from Z{...} (struct literal element) at escape5.go:115:9
// escape5.go:113:2:     from z2 := Z{...} (assign) at escape5.go:115:5
// escape5.go:113:2:   flow: {heap} = z2:
// escape5.go:113:2:     from &z2 (address-of) at escape5.go:116:19
// escape5.go:113:2:     from refStructV4(&z1, &z2) (call parameter) at escape5.go:116:13
// escape5.go:112:2: i1 escapes to heap:
// escape5.go:112:2:   flow: z1 = &i1:
// escape5.go:112:2:     from &i1 (address-of) at escape5.go:114:10
// escape5.go:112:2:     from Z{...} (struct literal element) at escape5.go:114:9
// escape5.go:112:2:     from z1 := Z{...} (assign) at escape5.go:114:5
// escape5.go:112:2:   flow: {heap} = z1:
// escape5.go:112:2:     from &z1 (address-of) at escape5.go:116:14
// escape5.go:112:2:     from refStructV4(&z1, &z2) (call parameter) at escape5.go:116:13
// escape5.go:112:2: moved to heap: i1
// escape5.go:113:2: moved to heap: i2
//
// 注意
// - `moved to heap` 表示一个本地变量从栈上被移动到堆上。
//   - moved to heap: z
// - 和"escapes to heap"是有差异的。后者表示本地变量被复制到堆上。有一个复制本地的值，重新在堆上alloc的动作。
//   - &S{} escapes to heap
//
// 参考：
// =========================================
// 逃逸分析输出的解释 by Ian Lance Taylor
// =========================================
//  from [golang-dev mail-list](https://groups.google.com/g/golang-dev/c/Cf4tpaWP6rc/m/iUgOcpZvAQAJ)
//
//  "moved to heap" means that a local variable was allocated on the heap
//  rather than the stack.
//
//  "leaking param" means that the memory associated with some parameter
//  (e.g., if the parameter is a pointer, the memory to which it points)
//  will escape. This typically means that the caller must allocate that
//  memory on the heap.
//
//  "escapes to heap" means that some value was copied into the heap.
//  This differs from "moved to heap" in that with "moved to heap" the
//  variable was allocated in the heap. With "escapes to heap" the value
//  of some variable was copied, for example when assigning to a variable
//  of interface type, and that copy forced the value to be copied into a
//  newly allocated heap slot.
//
// 参考：
// Go Escape Analysis Flaws by Dmitry Vyukov (dvyukov) (Feb 10, 2015)
//  - https://docs.google.com/document/d/1CxgUBPlx9iJzkz9JWkb6tIpTe5q32QDmz8l0BouG0Cw/preview
//
// 注意粒度：不能区分struct的field，和slice的field，是把struct本身和slice本身当作最小粒度看的。
//
// from https://github.com/golang/go/blob/go1.16.4/src/cmd/compile/internal/gc/escape.go#L66-L82
// Every Go language construct is lowered into this representation,
// generally without sensitivity to flow, path, or context; and
// without distinguishing elements within a compound variable. For
// example:
//
//     var x struct { f, g *int }
//     var u []*int
//
//     x.f = u[0]
//
// is modeled simply as
//
//     x = *u
//
// That is, we don't distinguish x.f from x.g, or u[0] from u[1],
// u[2], etc. However, we do record the implicit dereference involved
// in indexing a slice.

//// ~/work/golang/binary/go1.13.0/go/bin/go tool compile -m=2 -l escape5.go
//// ~/work/golang/goroot/bin/go tool compile -m=2 -l escape5.go