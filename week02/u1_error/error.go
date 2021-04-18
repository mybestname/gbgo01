package main

import (
	"errors"
	"fmt"
	"reflect"
)
// Error Type vs. Sentinel Error

// ErrMyError is Sentinel Error
var (
	ErrMyError = errors.New("pkg_name: my error")
)

// MyError is Error Type
type MyError struct {
	Msg string
	File string
	Line int
}
func (e *MyError) Error() string {
	return fmt.Sprintf("%s:%d: %s", e.File, e.Line, e.Msg)
}

func doMyStuff() error {
	return &MyError{"some bad happen", "server.go", 42};
}

func main() {
	err := errors.New("test")
	if err.Error() == "test" {
		fmt.Printf("Go error is just a pointer to struct %s \n",reflect.TypeOf(err))
	}
	//Each call to New returns a distinct error value even if the text is identical.
	err = errors.New("pkg_name: my error")
	// using sentinel error
	if err != ErrMyError {
		fmt.Printf("Even the error text is identical, they are diffrent, since %p != %p\n", err, ErrMyError);
	}
	// using error type
	err = doMyStuff();
	if err!= nil {
		fmt.Printf("%v\n", err);
	}

	err = doMyStuff();
	switch err := err.(type) {
	case nil:
		// call succeeded, do nothing
	case *MyError:
		// 因为 MyError 是一个 type，调用者可以使用断言转换成这个类型，来获取更多的上下文信息。
		// 永远不要操作err.Error()的结果，那个是给人用的，不是给程序用的。
		fmt.Printf("error occurred at line:%d\n", err.Line)
	default:
		// unknown error
	}
}
