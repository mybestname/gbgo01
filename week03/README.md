# 第三周：并行编程

## Goroutine
- go的两个关键字在语言级别上对于并行的支持（go和chan）
    - 编程体验上：顺序的编码方式，而不是回调（callback）的方式 ，更符合人类直觉。
### goroutine（go关键字，go协程）的底层是什么？
  - 操作系统：进程和线程的区别。
    - 进程：应用程序的资源容器（地址空间，文件handle，device，线程）
    - 线程：OS调度的执行路径。应用程序如何执行（一定从主线程开始，主线程结束，进程结束）
    - 无论线程属于哪个进程，OS安排线程在CPU上执行。OS的调度算法。（哪个CPU运行哪个线程）
  - goroutine被Go runtime来调度，goroutine被映射到某个线程上。
    - OS不认识goroutine，只认识OS的线程
    - go runtime把goroutine调度到某个go runtime的逻辑处理器（即P上）
    - 在P上，goroutine像一个队列上的任务等待被处理。
    - 多个线程（OS级别）从P的队列中（goroutine）获得任务，在CPU上执行。
    - 这样goroutine就可以10w级别存在。
  - 并发不是并行（Concurrency is not Parallelism - Rob Pike）
    - 并行指两个或多个线程**同时**在**不同**CPU上执行代码。
  - goroutine泄漏
    - 对任何一个goroutine关键是搞懂如何结束
    - 什么时候结束和谁不让他结束

## 内存模型

> ...reads of a variable in one goroutine can be guaranteed to observe values produced by writes to the same variable in a different goroutine.
> 
> To serialize access, protect the data with channel operations or other synchronization primitives such as those in the sync and sync/atomic packages.
> 
> If you must read the rest of this document to understand the behavior of your program, you are being too clever. 
> 
> Don't be clever.
> 
> *-- from go语言官方文档《go语言内存模型》*

如果要保证串行访问，必须使用下面3选1
- chan
- sync
- sync/atomic

### 参考
- [Go官方文档 - Go内存模型](https://golang.org/ref/mem)
- [中文翻译](https://www.jianshu.com/p/5e44168f47a3)

## sync包

## chan
- 关闭channel必须先于receive发生
- 什么时候用，unbuffer vs. buffer
  - 选择unbuffer除非你非常清楚为何要buffer
  - 如果要用buffer，必须明白满的时候如何处理。
  - buffer的大小并不代表性能的提升，不是越大越好（你需要真正理解size的大小）

## context包

## 设计模式

- [Go Concurrency Patterns: Timing out, moving on (Andrew Gerrand)23/Sep/2010](https://blog.golang.org/concurrency-timeouts)
- [Go Concurrency Patterns (Rob Pike) Google I/O 2012](https://talks.golang.org/2012/concurrency.slide#1)
- [Advanced Go Concurrency Patterns (Andrew Gerrand) Google I/O 2013](https://blog.golang.org/io2013-talk-concurrency)
- [Go Concurrency Patterns: Pipelines and cancellation (SameerAjmani)13/Mar/2014](https://blog.golang.org/pipelines)
- [Go Concurrency Patterns: Context (SameerAjmani)29/Jul/2014](https://blog.golang.org/context)

### 基于Channel
  - 超时控制 （context的超时控制基于channel）
  - 管道
  - 扇入/扇出
  - 取消

## 参考

- https://www.ardanlabs.com/blog/2018/11/goroutine-leaks-the-forgotten-sender.html
- https://www.ardanlabs.com/blog/2019/04/concurrency-trap-2-incomplete-work.html
- https://www.ardanlabs.com/blog/2014/01/concurrency-goroutines-and-gomaxprocs.html
- https://dave.cheney.net/practical-go/presentations/qcon-china.html#_concurrency
- https://golang.org/ref/mem
- https://blog.csdn.net/caoshangpa/article/details/78853919
- https://blog.csdn.net/qcrao/article/details/92759907
- https://cch123.github.io/ooo/
- https://blog.golang.org/codelab-share
- https://dave.cheney.net/2018/01/06/if-aligned-memory-writes-are-atomic-why-do-we-need-the-sync-atomic-package
- http://blog.golang.org/race-detector
- https://dave.cheney.net/2014/06/27/ice-cream-makers-and-data-races
- https://www.ardanlabs.com/blog/2014/06/ice-cream-makers-and-data-races-part-ii.html
- https://medium.com/a-journey-with-go/go-how-to-reduce-lock-contention-with-the-atomic-package-ba3b2664b549
- https://medium.com/a-journey-with-go/go-discovery-of-the-trace-package-e5a821743c3c
- https://medium.com/a-journey-with-go/go-mutex-and-starvation-3f4f4e75ad50
- https://www.ardanlabs.com/blog/2017/10/the-behavior-of-channels.html
- https://medium.com/a-journey-with-go/go-buffered-and-unbuffered-channels-29a107c00268
- https://medium.com/a-journey-with-go/go-ordering-in-select-statements-fd0ff80fd8d6
- https://www.ardanlabs.com/blog/2017/10/the-behavior-of-channels.html
- https://www.ardanlabs.com/blog/2014/02/the-nature-of-channels-in-go.html
- https://www.ardanlabs.com/blog/2013/10/my-channel-select-bug.html
- https://blog.golang.org/io2013-talk-concurrency
- https://blog.golang.org/waza-talk
- https://blog.golang.org/io2012-videos
- https://blog.golang.org/concurrency-timeouts
- https://blog.golang.org/pipelines
- https://www.ardanlabs.com/blog/2014/02/running-queries-concurrently-against.html
- https://blogtitle.github.io/go-advanced-concurrency-patterns-part-3-channels/
- https://www.ardanlabs.com/blog/2013/05/thread-pooling-in-go-programming.html
- https://www.ardanlabs.com/blog/2013/09/pool-go-routines-to-process-task.html
- https://blogtitle.github.io/categories/concurrency/
- https://medium.com/a-journey-with-go/go-context-and-cancellation-by-propagation-7a808bbc889c
- https://blog.golang.org/context
- https://www.ardanlabs.com/blog/2019/09/context-package-semantics-in-go.html
- https://golang.org/ref/spec#Channel_types
- https://drive.google.com/file/d/1nPdvhB0PutEJzdCq5ms6UI58dp50fcAN/view
- https://medium.com/a-journey-with-go/go-context-and-cancellation-by-propagation-7a808bbc889c
- https://blog.golang.org/context
- https://www.ardanlabs.com/blog/2019/09/context-package-semantics-in-go.html
- https://golang.org/doc/effective_go.html#concurrency
- https://zhuanlan.zhihu.com/p/34417106?hmsr=toutiao.io
- https://talks.golang.org/2014/gotham-context.slide#1
- https://medium.com/@cep21/how-to-correctly-use-context-context-in-go-1-7-8f2c0fafdf39

