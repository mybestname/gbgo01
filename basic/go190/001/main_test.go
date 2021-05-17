package main_test

import "fmt"

func ExampleDefer(){
	defer func() { fmt.Println("1")}()
	defer func() { fmt.Println("2")}()
	defer func() { fmt.Println("3")}()
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(4)
		}
	}()
	panic("panic")
	// Output:
	// 4
	// 3
	// 2
	// 1
}

// S:
//  While executing a function F, an explicit call to panic or a run-time panic
//	terminates the execution of F.
//	   - Any functions deferred by F are then executed as usual.
//	   - Next, any deferred functions run by F's caller are run, and so on up to
//	     any deferred by the top-level function in the executing goroutine.
//	   - At that point, the program is terminated and the error condition is
//	     reported, including the value of the argument to panic. This termination
//	     sequence is called panicking.
func ExamplePanic(){
	F := func() {
		defer func() { fmt.Println("2")}()
		defer func() {
			if e:= recover(); e!=nil {
				fmt.Println("3")
			}
		}()
		panic("panic")
	}
	Caller := func() {
		defer func() { fmt.Println("1") }()
		F()
	}
	Caller()
	// Output:
	// 3
	// 2
	// 1
}

// S:
//  The return value of recover is nil if any of the following conditions holds:
//    - panic's argument was nil;
//    - the goroutine is not panicking;
//    - recover was not called directly by a deferred function.
//
// recover() :
//     If recover is called outside the deferred function it will
//   not stop a panicking sequence. In this case, or when the goroutine is not
//   panicking, or if the argument supplied to panic was nil, recover returns
//   nil. Thus the return value from recover reports whether the goroutine is
//   panicking.
//
// [Go Blog: Defer, Panic, and Recover](https://blog.golang.org/defer-panic-and-recover):
//    During normal execution, a call to recover will return nil and have
//    no other effect.
//
func ExamplePanic_second(){
	F := func() {
		defer func() { fmt.Println("2")}()
		defer func() {
			if m:=recover(); m!=nil {
				fmt.Println(m)
			}
		}()
		if n:= recover(); n!=nil {
			fmt.Println("n not nil")
		}else{
			fmt.Println("n is nil")
		}
		panic("panic")
	}
	Caller := func() {
		defer func() { fmt.Println("1") }()
		F()
	}
	Caller()
	// Output:
	// n is nil
	// panic
	// 2
	// 1
}
