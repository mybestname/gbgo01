package main

func main() {

}

// ListDirectoryV1 将目录读取到一个 slice 中，然后返回整个切片，或者如果出现错误，则返回错误。
// 这是同步调用，ListDirectory 的调用方会阻塞，直到读取所有目录条目。
// 根据目录的大小，这可能需要很长时间，并且可能会分配大量内存来构建目录条目名称的 slice。
func ListDirectoryV1(dir string) ([]string, error) {
	return []string{"foo"}, nil
}

// ListDirectoryV2 返回一个 chan string，将通过该 chan 传递目录。当通道关闭时，这表示不再有目录。
// 由于在返回后发生通道的填充，可在内部启动 goroutine 来填充通道。
// 不要这样做！！！
func ListDirectoryV2(dir string) chan string {
	return nil
}

// 问题，
// 1. 通过使用关闭通道作为信号存在二义性，例如无法区分是空目录与还是目录读取中遇到了错误这之间的区别。
//    这两种情况都会导致返回的通道立即关闭。
// 2. 调用者必须持续从通道读取，直到它关闭，因为这是调用者知道填充chan的goroutine已经停止的唯一方法。
//    这对调用者来说是一个严重的限制，调用者必须继续从通道读取数据直到chan关闭，即使已经收到了想要的答
//    案。调用者不能中通退出， 调用者对chan的处理如果不当（例如中途退出），那么也会使得API内部的
//    goroutine堵塞，这也很危险。

// 应该把是否使用并发的决定权留给调用者

// ListDirectoryV3 使用callback函数
// 更好的方式！ 将异步执行函数的决定权交给该函数的调用方通常更容易。
func ListDirectoryV3(dir string, fn func(string)) error{
	return nil
}

// 对比Go的API对类似功能的设计：
// https://golang.org/src/path/filepath/path.go?s=13729:13779#L450
// func WalkDir(root string, fn fs.WalkDirFunc) error
