package main
import (
	"fmt"
	"golang.org/x/xerrors"
)

func main() {
	err := &MyError2{Message: "oops", frame: xerrors.Caller(1)}
	fmt.Printf("%v\n", err)
	fmt.Println()
	fmt.Printf("%+v\n", err)

}

type MyError2 struct {
	Message string
	frame   xerrors.Frame
}

func (m *MyError2) Error() string {
	return m.Message
}

func (m *MyError2) Format(f fmt.State, c rune) { // implements fmt.Formatter
	xerrors.FormatError(m, f, c)
}

func (m *MyError2) FormatError(p xerrors.Printer) error { // implements xerrors.Formatter
	p.Print(m.Message)
	if p.Detail() {
		m.frame.Format(p)
	}
	return nil
}