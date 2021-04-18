package main

import (
	"errors"
	"fmt"
	"reflect"
)
var (
	ErrMyError = errors.New("pkg_name: my error")
)

func main() {
	err := errors.New("test")
	if err.Error() == "test" {
		fmt.Printf("Go error is just a pointer to struct %s \n",reflect.TypeOf(err))
	}
	//Each call to New returns a distinct error value even if the text is identical.
	err = errors.New("pkg_name: my error")
	if err != ErrMyError {
		fmt.Printf("Even the error text is identical, they are diffrent, since %p != %p\n", err, ErrMyError);
	}
}
