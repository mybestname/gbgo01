package main

// E:
// 1. Go's defer statement schedules a function call (the deferred function)
//    to be run immediately before the function executing the defer returns.
// 2. Deferred functions are executed in LIFO order.
// S:
// 1. A "return" statement in a function F terminates the execution of F, and
//   optionally provides one or more result values. Any functions deferred
//   by F are executed before F returns to its caller.
// 2. A "defer" statement invokes a function whose execution is deferred to
//   the moment the surrounding function returns, either because the surrounding
//   function executed a return statement, reached the end of its function body,
//   or because the corresponding goroutine is panicking.
// 3. While executing a function F, an explicit call to panic or a run-time panic
//	 terminates the execution of F.
//	   - Any functions deferred by F are then executed as usual.
//	   - Next, any deferred functions run by F's caller are run, and so on up to
//	     any deferred by the top-level function in the executing goroutine.
//	   - At that point, the program is terminated and the error condition is
//	     reported, including the value of the argument to panic. This termination
//	     sequence is called panicking.
func main() {
	defer func() { println("1")}()
	defer func() { println("2")}()
	defer func() { println("3")}()
	panic("panic")
	println("done")
}
// output
// 3
// 2
// 1
// panic: panic
//
// goroutine 1 [running]:
// main.main()
// ...... main.go:11 .....



