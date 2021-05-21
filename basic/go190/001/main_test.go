package main_test

import (
	"fmt"
	"time"
)

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
//
//  Suppose a function G defers a function D that calls recover and a panic occurs
//  in a function on the same goroutine in which G is executing. When the running
//  of deferred functions reaches D, the return value of D's call to recover will
//  be the value passed to the call of panic. If D returns normally, without starting
//  a new panic, the panicking sequence stops.
func ExamplePanic(){
	F := func() {
		defer func() { fmt.Println("2")}()   // 2. F的defer按顺序调用
		defer func() {                            // 1. F的defer按顺序调用，recover不为空，返回panic函数传入的值
			if e:= recover(); e!=nil {            //    因为该函数正常终止，没有新的panic发生，所以panicking过程
				fmt.Println("1",e)           //    被停止，runtime不会终止，程序继续执行2
			}
		}()
		panic("panic")                         // 0. F函数显示调用panic会触发runtime的终止。
	}
	Caller := func() {
		defer func() { fmt.Println("3") }()   // 3. F的调用者的defer被调用
		F()
	}
	Caller()
	// Output:
	// 1 panic
	// 2
	// 3
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
// E:
//  A call to recover stops the unwinding and returns the argument passed to panic.
//  Because the only code that runs while unwinding is inside deferred functions,
//  recover is only useful inside deferred functions.
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
		if n:= recover(); n!=nil {               // 对于正常的执行流程，recover()没有意义，只是nil
			fmt.Println("n not nil")         // 所以E说：recover is only useful inside deferred function
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


// S:
//  a function G defers a function D that calls recover a panic occurs in a function
//  on **the same goroutine** in which G is executing
// E:
//  When panic is called, including implicitly for run-time errors such as indexing a
//  slice out of bounds or failing a type assertion, it immediately stops execution
//  of the current function and begins unwinding the stack of the goroutine, running
//  any deferred functions along the way. If that unwinding reaches the top of the
//  goroutine's stack, **the program dies**.
//  However, it is possible to use the built-in function recover to regain control of
//  **the goroutine** and resume normal execution.
//
// 注意 panic/recover 和 goroutine的对应关系！
func ExamplePanic_third(){
	F := func() {
		defer func() { fmt.Println("3")}()
		panic("panic")                          // goroutine发生panic时，只会调用自己栈上的defer
		                                           // 如果本身没有recover逻辑，程序就会终止，
		                                           // 和其它goroutine的调用栈上defer的recovery()是无关的
	}
	Caller := func() {
		defer func() { fmt.Println("2") }()
		defer func() {
			if n:= recover(); n!=nil {             // 这里的recover只能恢复调用caller的goroutine本身的panic！
				fmt.Println(n)                     // 无法恢复其它goroutine的panic。
			}
		}();
		go F()                                     // 这里的F是被一个不同的goroutine调用的。
		fmt.Println("1")
	}
	Caller()
	//output:
	//3
	//panic: panic
	//.. . ...
}

func ExamplePanic_fixed(){
	F := func() {
		defer func() { fmt.Println("3")}()
		panic("panic")                          // goroutine发生panic时，只会调用自己栈上的defer
		// 如果goroutine的调用栈上没有recover逻辑，程序就会终止，
		// 和其它goroutine的调用栈上defer的recovery()是无关的
		// F本身是没有defer的recovery的，那么它本身不能安全的被goroutine调用。即go F()是不安全的
	}
	Caller := func() {
		defer func() { fmt.Println("2") }()
		defer func() {
			if n:= recover(); n!=nil {             // 这里的recover只能恢复调用caller的goroutine本身的panic！
				fmt.Println(n)                     // 无法恢复其它goroutine的panic。
			}
		}();
		F()                                        // 这里的F不再是被一个不同的goroutine调用的。
		                                           // 而是和Caller一个goroutine
		                                           // 所以里面的panic可以被caller本身的defer recovery处理
		fmt.Println("1")
	}
	go Caller()                                    // 这里启的goroutine调用的是一个安全的函数
	time.Sleep(100*time.Millisecond)
	//output:
	//3
	//panic
	//2
}

// E上举的一个例子：安全调用：
// E：
// One application of recover is to shut down a failing goroutine inside a server without
// killing the other executing goroutines.
// ```go
//  func server(workChan <-chan *Work) {
//    for work := range workChan {
//        go safelyDo(work)                      // 使用独立goroutine调用时候要考虑panic的处理，来保证安全
//                                               // 如果safelyDo中出现panic，只会终止goroutine，不会影响整个程序。
//    }
//  }
//
//  func safelyDo(work *Work) {                  // 对do建立一个含recovery的保护层
//    defer func() {
//        if err := recover(); err != nil {
//            log.Println("work failed:", err)
//        }
//    }()
//    do(work)                                   // do的调用是安全的。panic被recover
// }
// ```
// E的第二个例子：适合内部封装错误，使用panic/recover连用，对外暴露错误，对内则使用panic来快速终止调用栈（simplify error handling）
// ```golang
//    // Error is the type of a parse error; it satisfies the error interface.
//    type Error string
//    func (e Error) Error() string {
//        return string(e)
//    }
//
//    // error is a method of *Regexp that reports parsing errors by
//    // panicking with an Error.
//    func (regexp *Regexp) error(err string) {
//        panic(Error(err))                                                   //对内用panic简化错误处理
//    }
//
//    // Compile returns a parsed representation of the regular expression.
//    func Compile(str string) (regexp *Regexp, err error) {
//        regexp = new(Regexp)
//        // doParse will panic if there is a parse error.
//        defer func() {
//            if e := recover(); e != nil {
//                regexp = nil    // Clear return value.
//                err = e.(Error) // Will re-panic if not a parse error.       //对外重新封装为错误
//            }
//        }()
//        return regexp.doParse(str), nil
//    }
// ```

//
// 总结：关于panic/recovery的注意点
// - recovery只有在defer中才有意义
// - 使用go关键字调用一定要特别注意，panic/recovery的安全性是和goroutine绑定的，不能只关注 defer()/recovery()的写法。
// - 所谓的特别注意指：
//   + panic/recovery是和goroutine绑定的。
//   + 这个原则使得对于每一个go关键字调用都要特别注意被调用函数的panic/recovery安全性。
//   + 而且一个上层函数的安全不能保证内层函数的安全，因为可能含有其它go关键字调用。
//   + 所以要仔细检查每一个go关键字的调用（即每一个goroutine的调用栈），才能确保panic/recovery安全性。
//
