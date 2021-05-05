package main

import (
	"bytes"
	"io"
	"os"
	"sync"
	"time"
)

var bufPool = sync.Pool{
	New: func() interface{} {
		// The Pool's New function should generally only return pointer
		// types, since a pointer can be put into the return interface
		// value without an allocation:
		return new(bytes.Buffer)
	},
}

// timeNow is a fake version of time.Now for tests.
func timeNow() time.Time {
	return time.Unix(1136214245, 0)
}

func Log(w io.Writer, key, val string) {
	b := bufPool.Get().(*bytes.Buffer)
	b.Reset()
	// Replace this with time.Now() in a real logger.
	b.WriteString(timeNow().UTC().Format(time.RFC3339))
	b.WriteByte(' ')
	b.WriteString(key)
	b.WriteByte('=')
	b.WriteString(val)
	w.Write(b.Bytes())
	bufPool.Put(b)
}


// from https://golang.org/pkg/sync/#Pool
// A Pool is a set of temporary objects that may be individually saved and retrieved.
// Pool的使用场景是保存并复用临时对象，来减少内存分配，从而降低GC的压力。
func main() {
	// 这里以打日志这种场景来说明pool的用例。很高频的操作，（ngix在栈上申请局部的buffer，进行字符串拼接再打印，用栈上内存地址直接刷到磁盘，
	// 性能非常好，不申请任何堆上对象，同时复用的还是栈上的空间）
	// 但是go不行，因为如果go在栈上拿buffer，如果用fmt库，会被判定为逃逸，造成堆上分配。
	// 目前只能使用sync.Pool，但是只能放没有状态的，随时可以被回收的对象。因为在Pool中的对象是不确定在什么时候回收的。
	// sync.Pool是一个全局对象，类似threadlocal
	// 内部实现使用ringbuff加双向链表
	Log(os.Stdout, "path", "/search?q=flowers")
}