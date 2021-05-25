package escape4

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"testing"
	"time"
)

func TestStringsApi(t *testing.T) {
	// strings API 为何使用值语义？
	// 因为是Go的基础类型。go的基础类型应该使用值语义。
	s := strings.Replace("oink oink oink", "k", "ky", 2)
	if s != "oinky oinky oink" {
		t.Fail()
	}
	s = strings.Replace("oink oink oink", "oink", "moo", -1)
	if s!= "moo moo moo" {
		t.Fail()
	}
	//func Replace(s, old, new string, n int) string  值语义
}

func TestNetApi(t *testing.T) {
	// slice, map, inteface, function, channel -> reference type
	// 为何这些引用类型应该使用值语义？因为他们设计上就是放在栈上的。
	ip := net.IPv4(192,168,1,1)
	mask := net.IPv4Mask(255,255,255,0)
	s := ip.Mask(mask).To4().String()
	if s!= "192.168.1.0" {
		t.Fail()
	}
	// func (ip IP) Mask(mask IPMask) IP  值语义
	// func (ip IP) To4() IP              值语义
	// func (ip IP) String() string       值语义
}

func TestAppend(t *testing.T) {
	s := []byte{1,2,3}
	s = append(s, 4)

	for i,v := range []byte{1,2,3,4} {
		if s[i] != v {
			t.Fail()
		}
	}
	var data []string
	data = append(data, "string")

	if data[0] != "string" {
		t.Fail()
	}
	// func append(slice []Type, elems ...Type) []Type  值语义

}

func TestUnmarshalText(t *testing.T){
	// 只有特殊的情况才需要指针语义
	text  := []byte("192.168.0.1")
	ip := &net.IP{}
	err := ip.UnmarshalText(text)
	if err!=nil || ip.String() != "192.168.0.1" {
		t.Fail()
	}
	// func (ip *IP) UnmarshalText(text []byte) error  指针语义
	// 这个必须是指针语义。因为如果是值语义，那么UnmarshalText改变的是
	// 不是receiver的内部状态，那么是没有效果的。
	// 参考下面的myIp.UnmarshalText(text)
	var myIp *MyIp  // 注意这里是否声明为指针是一样的，go会自动进行判断
	// 参考： Methods and pointer indirection
	//  as a convenience, Go interprets
	//    v.method() as (&v).method() if method()_is a pointer receiver.
	//    p.method() as (*p).method() if method() is a value receiver.
	// - https://tour.golang.org/methods/7
	// - https://tour.golang.org/methods/6
	// 所以：
	// 关键是方法是怎么声明的，引用变量会自动转换，以适应方法的签名。
	// 以方法签名为准。来决定是值语义，还是指针语义
	myIp = &MyIp{}
	err = myIp.UnmarshalText(text)  // 这里go自动变为 (*myIp).UnmarshalText，所以关键还是需要在方法签名上使用指针语义。
	if err!=nil ||  myIp.String() != ""{ //没有效果，还是空
		println(myIp.String())
		t.Fail()
	}
}

type MyIp []byte

func (ip MyIp) String() string {
	return fmt.Sprintf("%v",string(ip));
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The IP address is expected in a form accepted by ParseIP.
func (ip MyIp) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		ip = nil   //这个也是一样，对外界没有影响。
		return nil
	}
	ip = text  //这个是不会传出这个方法的，所以没有效果。
	return nil
	// func (ip MyIp) UnmarshalText(text []byte)
	// 等效于
	// func UnmarshalText(ip MyIp, text []byte) error
	// 而
	// func (ip *IP) UnmarshalText(text []byte) error
	// func UnmarshalText(ip *IP, text []byte) error
	//
	// 所以为了改变receiver的内部状态，必须使用指针语义。
}

func (ip *MyIp) UnmarshalTextOk(text []byte) error {
	if len(text) == 0 {
		*ip = nil
		return nil
	}
	*ip = text
	return nil

}

func TestTimeApi(t *testing.T) {
	// time包也是这样，都是值语义。
	now := time.Now()                     // 值语义
	add := now.Add(100*time.Millisecond)  // 值语义
	d := add.Sub(now)                     // 值语义
	if d != 100*time.Millisecond {
		t.Fail();
	}

	// 只有Unmarshal和Decode是指针语义
	//
	// func (t *Time) GobDecode(data []byte) error
	// func (t *Time) UnmarshalBinary(data []byte) error
	// func (t *Time) UnmarshalJSON(data []byte) error
	// func (t *Time) UnmarshalText(data []byte) error
	s := time.Time{}
	err := s.UnmarshalText([]byte("2021-05-25T00:00:00Z"))
	if err!=nil || s.Format(time.RFC3339) != "2021-05-25T00:00:00Z"{
		t.Fail()
	}
}

func TestOsFileApi(t *testing.T) {
	// 全部是指针语义
	//    func (f *File) Chdir() error
	//    func (f *File) Chmod(mode FileMode) error
	//    func (f *File) Chown(uid, gid int) error
	//    func (f *File) Close() error
	//    func (f *File) Fd() uintptr
	//    func (f *File) Name() string
	//    func (f *File) Read(b []byte) (n int, err error)
	//    func (f *File) ReadAt(b []byte, off int64) (n int, err error)
	//    func (f *File) ReadDir(n int) ([]DirEntry, error)
	//    func (f *File) ReadFrom(r io.Reader) (n int64, err error)
	//    func (f *File) Readdir(n int) ([]FileInfo, error)
	//    func (f *File) Readdirnames(n int) (names []string, err error)
	//    func (f *File) Seek(offset int64, whence int) (ret int64, err error)
	//    func (f *File) SetDeadline(t time.Time) error
	//    func (f *File) SetReadDeadline(t time.Time) error
	//    func (f *File) SetWriteDeadline(t time.Time) error
	//    func (f *File) Stat() (FileInfo, error)
	//    func (f *File) Sync() error
	//    func (f *File) SyscallConn() (syscall.RawConn, error)
	//    func (f *File) Truncate(size int64) error
	//    func (f *File) Write(b []byte) (n int, err error)
	//    func (f *File) WriteAt(b []byte, off int64) (n int, err error)
	//    func (f *File) WriteString(s string) (n int, err error)
	// 显然，文件一定是指针语义，值语义显然是不适当的。
	// 即使不改变内部状态，依旧使用指针语义。

	f, e:= os.Open("test")
	if e!=nil {
		if !errors.Is(e, os.ErrNotExist)  || f!=nil {
			t.Fail()
		}
	}

}

//go:generate go tool compile -m=2 escape4_test.go

