package main


type H interface {
	Hello()
}

type A struct {}

func (A) Hello()  {
}

type B struct {}

func (B) Hello() {
}

func sayHelloA(a A) {
	a.Hello()
}

func SayHello(h H) {
	h.Hello()
}

func main() {
	var h H;
	var a A;
	a = A{}   // not-escape
	sayHelloA(a)
	h = A{}     //A{} 逃逸
	SayHello(h) // 方法本身调用的 call parameter
	h = &A{}   // &A{} 逃逸  因为这里有一次 interface-converted
	SayHello(h) // 另外方法本身调用的 call parameter
	b := B{};   // 逃逸被取消，因为激活了devirtualizing,
	SayHello(b) // 如果加上-l，则b逃逸 (call parameter) 因为方法调用是interface
	            // 这个如果去掉-l，则不会逃逸，因为激活了devirtualizing
	            // devirtualizing h.Hello to B
}

//go:generate go build -gcflags "-m=2 -l" main.go
//go:generate go build -gcflags "-m=2" main.go
//go:generate go build -gcflags "-m" main.go
//go:generate rm main


//
// cmd/compile: de-virtualize interface calls
// https://go-review.googlesource.com/c/go/+/38139
// https://github.com/golang/go/commit/295307ae78f8dd463a2ab8d85a1592ca76619d36
//
//  After this change, code like
//
//      h := sha1.New()
//      h.Write(buf)
//      sum := h.Sum()
//
//  gets compiled into static calls rather than
//  interface calls, because the concrete type of
//  'h' is known statically.

//  ./main.go:10:6: can inline A.Hello with cost 0 as: method(A) func() {  }
//  ./main.go:15:6: can inline B.Hello with cost 0 as: method(B) func() {  }
//  ./main.go:18:6: can inline sayHelloA with cost 3 as: func(A) { a.Hello() }
//  ./main.go:19:9: inlining call to A.Hello method(A) func() {  }
//  ./main.go:22:6: can inline SayHello with cost 60 as: func(H) { h.Hello() }
//  ./main.go:26:6: cannot inline main: function too complex: cost 221 exceeds budget 80
//  ./main.go:30:11: inlining call to sayHelloA func(A) { a.Hello() }
//  ./main.go:30:11: inlining call to A.Hello method(A) func() {  }
//  ./main.go:32:10: inlining call to SayHello func(H) { h.Hello() }
//  ./main.go:34:10: inlining call to SayHello func(H) { h.Hello() }
//  ./main.go:36:10: inlining call to SayHello func(H) { h.Hello() }
//  ./main.go:36:10: devirtualizing h.Hello to B
//  ./main.go:22:15: parameter h leaks to {heap} with derefs=0:
//  ./main.go:22:15:   flow: {heap} = h:
//  ./main.go:22:15:     from h.Hello() (call parameter) at ./main.go:23:9
//  ./main.go:22:15: leaking param: h
//  ./main.go:33:6: &A{} escapes to heap:
//  ./main.go:33:6:   flow: h = &{storage for &A{}}:
//  ./main.go:33:6:     from &A{} (spill) at ./main.go:33:6
//  ./main.go:33:6:     from &A{} (interface-converted) at ./main.go:33:4
//  ./main.go:33:6:     from h = &A{} (assign) at ./main.go:33:4
//  ./main.go:33:6:   flow: h = h:
//  ./main.go:33:6:     from h := h (assign-pair) at ./main.go:34:10
//  ./main.go:33:6:   flow: {heap} = h:
//  ./main.go:33:6:     from h.Hello() (call parameter) at ./main.go:34:10
//  ./main.go:31:4: A{} escapes to heap:
//  ./main.go:31:4:   flow: h = &{storage for A{}}:
//  ./main.go:31:4:     from A{} (spill) at ./main.go:31:4
//  ./main.go:31:4:     from h = A{} (assign) at ./main.go:31:4
//  ./main.go:31:4:   flow: h = h:
//  ./main.go:31:4:     from h := h (assign-pair) at ./main.go:34:10
//  ./main.go:31:4:   flow: {heap} = h:
//  ./main.go:31:4:     from h.Hello() (call parameter) at ./main.go:34:10
//  ./main.go:31:4: A{} escapes to heap
//  ./main.go:33:6: &A{} escapes to heap
//  ./main.go:36:10: b does not escape
//  <autogenerated>:1: parameter .this leaks to {heap} with derefs=0:
//  <autogenerated>:1:   flow: {heap} = .this:
//  <autogenerated>:1:     from .this.Hello() (call parameter) at <autogenerated>:1
//  <autogenerated>:1: leaking param: .this
//  <autogenerated>:1: inlining call to A.Hello method(A) func() {  }
//  <autogenerated>:1: .this does not escape
//  <autogenerated>:1: inlining call to B.Hello method(B) func() {  }
//  <autogenerated>:1: .this does not escape