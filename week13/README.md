# 第13课 Go语言实践 Runtime
## 目录
- Goroutine原理
- 内存分配原理
- GC原理
- Channel原理

## Goroutine原理
### Goroutine
- 定义
  - 在同一用户地址空间里并行独立执行的functions。
  - channel用于goroutine间的通信和同步访问控制。
  - “Goroutine 是一个与其他 goroutines 并行运 行在同一地址空间的 Go 函数或方法。一个运 行的程序由一个或更多个 goroutine 组成。它
     与线程、协程、进程等不同。它是一个 goroutine” —— Rob Pike
- 和线程的区别
  - 内存开销低（栈空间默认2kb，linux amd64 gov1.4)，栈空间自动扩缩容。
    - vs。线程（1-8MB 栈空间，POSIX），guard page区隔离。栈空间初始化完成后不能再变化。
  - 创建/销毁的销毁低（goroutine是用户态的，go runtime管理）
    - vs. 线程（os进程模型的扩展kernel API，内核级别的调用，需要进入内核调度）
  - 调度切换成本小（goroutine 用户态，3个寄存器，上下文简单，成本小，100～200纳秒）
    - vs. 线程 上下文复杂，保存成本高，较多寄存器，很多公平算法，复杂的时间统计（例如Linux CFS调度器之虚拟时钟vruntime与调度延迟）等。1000-1500纳秒）
  - 复杂性低（开发维护成本低很多）
    - vs. 线程 创建/退出很复杂，多个线程间通讯（sharememory）模型很复杂。
### GMP模型
- M：N模型
  - M线程（内核的task_struct)，N (goroutine)
    > struct task_struct
    >  - https://github.com/torvalds/linux/blob/master/include/linux/sched.h
    >  - 对于linux实现，进程/线程都是task_struct
    >  - 只是线程把字段（mm_struct）设置为共享同一地址空间，即为线程。
- GMP
  - G 代表 Goroutine，每一次 `go func()` 都代表一个G，无限制。
      - `struct runtime.g`
        - 当前goroutine的状态，堆栈，上下文
  - M 代表 工作线程（OS thread），也被称为Machine，
      - `struct runtime.m`
      - M有操作系统的线程栈
      - M.stack -> G.stack
      - M的PC寄存器执行G的函数，然后去执行
  - Go 1.2 只有 G和M，即GM模型。的问题
      - 限制了并发的伸缩性。吞吐不好
      - 单一的全局互斥锁和集中化的状态控制。（所有goroutine的相关操作：创建，结束，重调度都要上锁）
      - G的传递问题：M需要经常在M之间传递"可运行"的goroutine。刚创建的G放到全局队列，而不是本地M运行。
      - 每一个M持有一个内存缓存M.mcache和stacklloc，而大部分阻塞在syscall（1：100，Go代码：syscall），
        -而syscall并不需要内存缓存，只有Go代码才会用内存缓存，造成大量的内存浪费。
        -另外这个缓存的数据局部型（内存亲缘性）差，因为G被调用到同一个M的概率不高。
          - G刚在M1预热完，又切换到M2去了。
      - 线程阻塞和死锁：
        - M找不到G，M频繁进入阻塞/唤醒来进行检查，以便及时发现新的G来执行。
  - P 是为了解决Go1.2问题出现的。代表Processor
    - 代表M所需的上下文环境
    - 处理用户级别代码逻辑的处理器
    - 衔接M和G的调度，将等待执行的G和M对接。
    - P有任务时候创建/或唤醒 一个M来执行队列里面的任务。
    - P决定了并行任务的数量。
    - GOMAXPROCS （go1。5被默认设置为可用的CPU核数）
    - mcache和stacklloc转移到了P，
    - G队列分为两类：1个全局G队列，每个P有一个本地G队列。
      - P优先去本地队列，让M执行G
      - 看全局队列还有没有G  
      - 同时P还能去别人的本地队列里面"偷"G
        - 本地队列是一个LockFree的队列，使用CAS原子操作来保证无锁的原子性。

```go
// G - goroutine.
// M - worker thread, or machine.
// P - processor, a resource that is required to execute Go code.
//     M must have an associated P to execute Go code, however it can be
//     blocked or in a syscall w/o an associated P.  
```
### work-stealing 调度算法
- 算法描述
  - 当一个P执行完本地所有的G，
    - 全局队列不为空
      - 从全局队列获取（当前个数/GOMAXPROCS）个
        > - 这里说法有误，查网上的说法也有问题
        >   - [Golang调度器GMP原理与调度全分析 by 刘丹冰Aceld](https://mp.weixin.qq.com/s/SEPP56sr16bep4C_S0TLgA)
        >   - `n = min(len(GQ)/GOMAXPROCS + 1, len(GQ/2))`
        >   - 但这里也有问题
        > - https://github.com/golang/go/blob/master/src/runtime/proc.go#L5777-L5806
        > ```go
        > 	n := sched.runqsize/gomaxprocs + 1
        > 	if n > sched.runqsize {
        > 		n = sched.runqsize
        > 	}
        > 	if max > 0 && n > max {
        > 		n = max
        > 	}
        > 	if n > int32(len(_p_.runq))/2 {
        > 		n = int32(len(_p_.runq)) / 2
        > 	}
        > ```   
        > - 应该是： `min(len(GQ)/GOMAXPROCS + 1, len(GQ), len(LQ)/2)`
    - 全局队列为空
      - 尝试挑选一个P（受害者），从它的本地队列中窃取一半的G
        > - https://github.com/golang/go/blob/master/src/runtime/proc.go#L6182-L6199
        >   ```go
        >   // Steal half of elements from local runnable queue of p2
        >   // and put onto local runnable queue of p.
        >   // Returns one of the stolen elements (or nil if failed).
        >   func runqsteal(_p_, p2 *p, stealRunNextG bool) *g {
        >   ``` 
      - 为了保证公平性，随机选一个P，而且遍历的顺序也随机化
        - 选择一个小于GOMAXPROCS且和它互质数的步长，来保证遍历顺序是随机的。
        > - https://github.com/golang/go/blob/master/src/runtime/proc.go#L3024
        > ```go
        > for enum := stealOrder.start(fastrand()); !enum.done(); enum.next() {
        > ```
        > - https://github.com/golang/go/blob/master/src/runtime/proc.go#L6379-L6395
        > ```go
        > var stealOrder randomOrder
        > 
        > // randomOrder/randomEnum are helper types for randomized work stealing.
        > // They allow to enumerate all Ps in different pseudo-random orders without repetitions.
        > // The algorithm is based on the fact that if we have X such that X and GOMAXPROCS
        > // are coprime, then a sequences of (i + X) % GOMAXPROCS gives the required enumeration.
        > type randomOrder struct {
        > 	count    uint32
        > 	coprimes []uint32
        > }
        > 
        > type randomEnum struct {
        > 	i     uint32
        > 	count uint32
        > 	pos   uint32
        > 	inc   uint32
        > }
        > ``` 
  - 如果某P的本地队列已经放满（256个时候），会放一半的G到全局队列。
  - 如果阻塞的syscall返回时候，找不到空闲的P，那么把G放到全局队列。
  - P每N个调度后，会去全局拿一个G。（1/61的时间）
    - 也就是当P老不饥饿时候，也不能老不查全局队列。会有一个查的机会。
     > - `schedtick%61 == 0`
     > - https://github.com/golang/go/blob/master/src/runtime/proc.go#L3351-L3360
     > ```go
     > 	if gp == nil {
     > 		// Check the global runnable queue once in a while to ensure fairness.
     > 		// Otherwise two goroutines can completely occupy the local runqueue
     > 		// by constantly respawning each other.
     > 		if _g_.m.p.ptr().schedtick%61 == 0 && sched.runqsize > 0 {
     > 			lock(&sched.lock)
     > 			gp = globrunqget(_g_.m.p.ptr(), 1)
     > 			unlock(&sched.lock)
     > 		}
     > 	}
     > ```    

#### Syscall
- 当调用syscall会解绑P，然后，M和G进入阻塞
- 而此时P进入特殊状态（syscall状态），表明这个P的G正在syscall中，此时P不能调度给其它的M
  - 因为大概率的M会短时间唤醒，那么M会优先和这个P绑定，这样有利于数据的局部性
- 但也有彻底解绑P的情况
  - 在执行syscall时候，如果某个P执行G超过10ms（一个sysmon tick 10ms），这时强制解绑P，状态为idle，放入idle list
  - 系统监视器 (system monitor)，称为 sysmon，会定 时扫描。在执行 syscall 时, 如果某个 P 的 G 执行超 过一个 sysmon tick(10ms)，就会把他设为 idle，重 新调度给需要的 M，强制解绑。
- syscall结束后，
  - 对 M，
    - 尝试获取同一个P，恢复执行G
    - 尝试获取idle list中的其它空闲P，恢复执行G
    - 找不到空闲的P，把G放回全局队列。M放入idle list
#### OS thread数量
- 当使用了Syscall，Go无法限制被阻塞的系统线程的数量。
- 这是一个重要的坑，对于syscall是不同的。
- 不要想当然的认为GOMAXPROCS会限制系统OS线程的数量。
> The GOMAXPROCS variable limits the number of operating system threads 
> that can execute user-level Go code simultaneously. There is no limit
> to the number of threads that can be blocked in system calls on behalf 
> of Go code; those do not count against the GOMAXPROCS limit. This 
> package's GOMAXPROCS function queries and changes the limit.
- GOMAXPROCS只是限制了P的数量，即让你的用户态的Go code可以执行的系统线程。
- 但是Block在syscall上的系统线程数量是不受限制的。
- 当涉及了syscall调用一定要小心有没有线程耗尽的问题。
  - 认真考虑 pthread exhaust 问题。

### Spining thread 自旋线程
- 线程自旋是相对于线程阻塞而言的，表象就是循环执行一个指定逻辑 
  - 调度逻辑，目的是不停 地寻找 G。
- 这样做的问题显而易⻅，如果 G 迟迟不来，CPU 会白白浪费在这无意义的计算上。
- 但好处也很明显，降低了 M 的上下文切换成本，提高了性能。
- 在两个地方引入自旋:
  - 类型1:M 不带 P 的找 P 挂载(一有 P 释放就结合) 
  - 类型2:M 带 P 的找 G 运行(一有 runable 的 G 就执行)
- 为了避免过多浪费 CPU 资源，**自旋的 M 最多只允许 GOMAXPROCS** (Busy P)。
- 同时当有类型1的自旋 M 存在时，类型2的自旋 M 就不阻塞，
   - 阻塞会释放 P，一释放 P 就⻢上被类型1的自旋 M 抢走了，没必要。
- 在新 G 被创建、M 进入系统调用、M 从空闲被激活这三种状态变化前，
  调度器会确保至少有一个自旋 M 存在(唤醒或者创建一个 M)，除非没有空闲的 P。
  - 当新 G 创建，如果有可用 P，就意味着新 G 可以被立即执行，
    即便不在同一个 P 也无妨，所以我们 保留一个自旋的 M
    - 这时应该不存在类型1的自旋只有类型2的自旋 
    - 就可以保证新 G 很快被运行。
  - 当 M 进入系统调用，意味着 M 不知道何时可以醒来，
    - 那么 M 对应的 P 中剩下的 G 就得有新的 M 来执行，
    - 所以我们保留一个自旋的 M 来执行剩下的 G
    - 这时应该不存在类型2的自旋只有类型1的自旋。
  - 如果 M 从空闲变成活跃，意味着可能一个处于自旋状态的 M 进入工作状态了，
    - 这时要检查并确保还 有一个自旋 M 存在，
    - 以防还有 G 或者还有 P 空着的。

### GMP 问题总结
- 单一全局互斥锁(Sched.Lock)和集中状态存储
  - G 被分成全局队列和 P 的本地队列，
  - 全局队列依旧是全局锁，但是使用场景明显很少，
  - P 本地队列使用无锁队列，使用原子操作来面对可能的并发场景。
- Goroutine 传递问题
  - G 创建时就在 P 的本地队列，可以避免在 G 之间传递(窃取除外)，
  - G 对 P 的数据局部性 好; 
  - 当 G 开始执行了，系统调用返回后 M 会尝试获取可用 P，获取到了的话可以避免在 M 之间 传递。而且优先获取调用阻塞前的 P，所以 G 对 M 数据局部性好，G 对 P 的数据局部性也好。
- Per-M 持有内存缓存 (M.mcache)
  - 内存 mcache 只存在 P 结构中，P 最多只有 GOMAXPROCS 个，远小于 M 的个数，
  - 所以内 存没有过多的消耗。
- 严重的线程阻塞/解锁
  - 通过引入自旋，保证任何时候都有处于等待状态的自旋 M， 避免在等待可用的 P 和 G 时频繁 的阻塞和唤醒。

#### sysmon
- sysmon 也叫监控线程，它无需 P 也可以运行， 
- 他是一个死循环，每20us~10ms循环一次，循环完一次就 sleep 一会，
- 为什么会是一个变动 的周期呢，主要是避免空转，如果每次循环都没什么需要做的事，
  那么 sleep 的时间就会加大。
- 功能  
  - 释放闲置超过5分钟的 span 物理内存;
  - 如果超过2分钟没有垃圾回收，强制执行;
  - 将⻓时间未处理的 netpoll 添加到全局队列; 
  - 向⻓时间运行的 G 任务发出抢占调度;
  - 收回因 syscall ⻓时间阻塞的 P;
- 抢占调度
  - 当 P 在 M 上执行时间超过10ms，sysmon 调用 preemptone 将 G 标记为 stackPreempt 。
  -  因此需要 在某个地方触发检测逻辑，Go 当前是在检查栈是否溢出的地方判定(morestack())，
     M 会保存当前 G 的 上下文，重新进入调度逻辑。
  - 死循环:issues/11462 
  - 信号抢占: [go1.14基于信号的抢占式调度实现原理](http://xiaorui.cc/archives/6535)
    - 异步抢占，注册 sigurg 信号，通过 sysmon 检测，
    - 对 M 对应的线程发送信号，触发注册的 handler，
    - 它 往当前 G 的 PC 中插入一条指令(调用某个方法)，
    - 在 处理完 handler，G 恢复后，自己把自己推到了 global queue 中。

#### Network poller
- Go 所有的 I/O 都是阻塞的。
  - 然后通过 goroutine + channel 来处理并发。
  - 因此所有的 IO 逻辑都是直 来直去的，
  - 你不再需要回调，不再需要 future，要 的仅仅是 step by step。
  - 这对于代码的可读性是很 有帮助的。
- G 发起网络 I/O 操作也不会导致 M 被阻塞(仅阻塞 G)，
  从而不会导致大量 M 被创建出来。
- 将异步 I/O 转换为阻塞 I/O 的部分称为 netpoller。
- 打开或接受 连接都被设置为非阻塞模式。
- 如果你试图对其进行 I/O 操作，并且文件描述符数据还没有准备好，
  - G 会进入 gopark 函数，将当前正在执行的 G 状态保存起来，然后切换到新的堆栈上执行新的 G。

- 那什么时候 G 被调度回来呢?
  - sysmon
  - schedule():M 找 G 的调度函数
  - GC:start the world
  
- 调用 netpoll() 在某一次调度 G 的过程中， 处于就绪状态的 fd 对应的 G 就会被调度 回来。
  
- G 的 gopark 状态
  - G 置为 waiting 状态，等待显示 goready 唤醒
  - 在 poller 中用得 较多，
  - 还有锁、chan 等场景会进入gopark。
  
#### Scheduler Affinity
- GM 调度器时代的，chan 操作导致的切换代价。
  - Goroutine#7 正在等待消息，阻塞在 chan。
    一旦收到消息，这个 goroutine 就被推到全局队列。 
    然后，chan 推送消息，goroutine#X 将在可用线程上运行，
    而 goroutine#8 将阻塞在 chan。 goroutine#7 现在在可用线程上运行。
- 在 chan 来回通信的 goroutine 会导致频繁 的 blocks，即频繁地在本地队列中重新排队。
- 然而，由于本地队列是 FIFO 实现， 如果另一个 goroutine 占用线程，
  unblock goroutine 不能保证尽快运行。
- 同时 Go 亲 缘性调度的一些限制:Work-stealing、系统调用。
- goroutine #9 在 chan 被阻塞后恢复。但是，它必须等待#2、#5和#4之后才能运行。
  goroutine #5将阻塞其线程，从而延迟goroutine #9，并使其面临被另一个 P 窃取的⻛险。
- 针对 communicate-and-wait 模式，进行了 **亲缘性调度的优化**。
- Go 1.5 在 P 中引入了一个 **runnext** 特殊字段，可以高优先级执行 unblock G。
- goroutine #9现在被标记为下一个可运行的。这种新的优先级排序允许 goroutine 在再次被阻塞 之前快速运行。
  这一变化对运行中的标准库产生了总体上的积极影响，提高了一些包的性能。
  
### Goroutine的生命周期

#### Go 程序启动
- 整个程序始于一段汇编， 而在随后的 runtime·rt0_go(也是汇编程序)中，会执行很多初 始化工作。
  - 绑定 m0 和 g0，
     - m0就是程序的主线程，程序启动必然会拥有一个主线程，这个就是 m0。
     - g0 负责调度，即 shedule() 函数。
  - 创建 P，绑定 m0 和 p0，
     - 首先会创建 GOMAXPROCS 个 P ，
     - 存储在 sched 的 空闲链表(pidle)。 
  - 新建任务 g 到 p0 本地队列，
     - m0 的 g0 会创建一个 指向 runtime.main() 的 g ，并放到 p0 的本地队列。
  - runtime.main(): 
     - 启动 sysmon 线程;
     - 启动 GC 协程;
     - 执行 init，即代码中的各种 init 函数;
     - 执行 main.main 函数。

#### OS thread 创建
- 准备运行的新 goroutine 将唤醒 P 以更好地分发工作。
- 这个 P 将创建一个与之关联的 M 绑定到一个 OS thread。 
- go func() 中 触发 Wakeup 唤醒机制:
  - 有空闲的 P 而没有在 spinning 状态的 M 时候, 
    需要去唤醒一个 空闲(睡眠)的 M 或者新建一个。
  - 当线程首次创建时，会执行一个 特殊的 G，即 g0，
    它负责管理和调度 G。
- 特殊的 g0
  - Go 基于两种断点将 G 调度到线程上:
     - 当 G 阻塞时:系统调用、互斥锁或 chan。
       阻塞的 G 进入睡眠模式/ 进入队列，
       并允许 Go 安排和运行等待其他的 G。
     - 在函数调用期间，如果 G 必须扩展其堆栈。
       这个断点允许 Go 调度 另一个 G 
       并避免运行 G 占用 CPU。
     - 在这两种情况下，运行调度程序的 g0 
       - 将当前 G 替换为另一个 G，即 ready to run。
       - 然后，选择的 G 替换 g0 并在线程上运行。
       - 与常规 G 相 反，g0 有一个固定和更大的栈。
  - Defer 函数的分配
  - GC 收集，比如 STW、扫描 G 的堆栈和标记、清除操作 
  - 栈扩容，当需要的时候，由 g0 进行扩栈操作

#### Schedule
- 在 Go 中，G 的切换相当轻便，其中需要保存的状 态仅仅涉及以下两个:
  - Goroutine 在停止运行前执行的指令，程序当前要运 行的指令是记录在程序计数器(PC)中的， 
    G 稍后 将在同一指令处恢复运行;
  - G 的堆栈，以便在再次运行时还原局部变量;
    - 在切换之前，堆栈将被保存，以便在 G 再次运行时进行 恢复:

- 从 g 到 g0 或从 g0 到 g 的切换是相当迅速的( 9-10ns)，它们只包含少量固定的指令。
  相反，对于调度阶段，调度程序需要检查许多资源以便确定下一个要运行的 G。
  - 当前 g 阻塞在 chan 上并切换到 g0:
    - 1、PC 和堆栈指针一起保存在内部结构中;
    - 2、将 g0 设置 为正在运行的 goroutine;
    - 3、g0 的堆栈替换当前堆栈;
  - g0 寻找新的 Goroutine 来运行
  - g0 使用所选的 Goroutine 进行切换: 
    - 1、PC 和堆栈指针是从其内部结构中获取的;
    - 2、程序跳 转到对应的 PC 地址;

#### Goroutine Recycle
- G 很容易创建，栈很小以及快速的上下文 切换。基于这些原因，开发人员非常喜欢 并使用它们。
- 然而，一个产生许多 shortlive 的 G 的程序将花费相当⻓的时间 来创建和销毁它们。
- 每个 P 维护一个 freelist G，保持这个列表是本地的，
  这样做的好处是不使用任何锁来 push/ get 一个空闲的 G。
- 当 G 退出当前工作时，它 将被 push 到这个空闲列表中。
- 为了更好地分发空闲的 G ，调度器也有自己的列表。
  - 它实际上有**两个列表**:
    - 一个包含已分配栈的 G，
    - 另一个包含释放过堆栈的 G (无栈)。 
- 锁保护 central list，因为任何 M 都可以访问它。 
  - 当**本地列表⻓度超过64**时，调度程序持有的列表 从 P 获取 G。
  - 然后一半的 G 将移动到中心列表。 
  - 需求回收 G 是一种节省分配成本的好方法。
    - 但是，由于堆栈是动态增⻓的，现有的G 最终可能会有一个大栈。
    - 因此，当堆栈增⻓(即超过2K) 时，Go 不会保留这些栈。

## 内存分配原理
- 堆栈 & 逃逸分析 
- 连续栈
- 内存结构
- 优化实践

### 堆和栈的定义
- Go 有两个地方可以分配内存
- 一个全局堆空间用来 动态分配内存，
- 另一个是每个 goroutine 都有的自身 栈空间。
#### 栈
- 栈区的内存一般由编译器自动进行分配和释放，其中存储着函数的入参以及局部变量，
  这些参数会随着函数的创建 而创建，函数的返回而销毁。(通过 CPU push & release)。
- A function has direct access to the memory inside its frame,
  through the frame pointer, but access to memory outside its
  frame requires indirect access.

#### 堆
- 堆区的内存一般由编译器和工程师自己共同进行管理分配，交给 Runtime GC 来释放。
- 堆上分配必须找到一块足够大的内存来存放新的变量数据。
- 后续释放时，垃圾回收器扫描堆空间寻找不再被使用的对象。
  - Anytime a value is shared outside the scope of a function’s 
    stack frame, it will be placed (or allocated) on the heap.
- 栈分配廉价，堆分配昂贵。
  - stack allocation is cheap and heap allocation is expensive.

#### 变量是在堆还是栈上?
- 其他语言，比如C，有明确的栈和堆的相关概念。
- 而 Go 声明语法并没有提到栈和 堆，而是交给 Go 编译器决定在哪分配内存， 保证程序的正确性，
- 在 Go FAQ 里面提到这么一段解释:
  - 从正确的⻆度来看，你不需要知道。Go 中的每个变量只要有引用就会一直存在。
    变量的存储位置(堆 还是栈)和语言的语义无关。
- 存储位置对于写出高性能的程序确实有影响。
  - 如果可能，Go 编译器将为该函数的堆栈侦(stack frame)中的函数分配本地变量。
  - 但是如果编译器在函数返回后无法证明变量未被引用，则编译器必须在会被垃圾回收的堆上分配
    变量以避免悬空指针错误。
  - 此外，如果局部变量非常大，将它存储在堆而不是 栈上可能更有意义。
- 在当前编译器中，如果变量存在取址，则该变量是堆上分配的候选变量。
  - 但是基础的逃逸分析可以 将那些生存不超过函数返回值的变量识别出来，并且因此可以分配在栈上。

#### 逃逸分析
- “通过检查变量的作用域是否超出了它所在的栈来决定 是否将它分配在堆上”的技术，
  其中“变量的作用域超出了它所在的栈”这种行为即被称为逃逸。
- 逃逸分析在大多数语言里属于静态分析:
  - 在编译期由静态代码分析来决定一个值是否能被分配在栈帧上，还是需要“逃逸”到堆上。
- 减少 GC 压力，栈上的变量，随着函数退出后系统 直接回收，不需要标记后再清除
- 减少内存碎片的产生
- 减轻分配堆内存的开销，提高程序的运行速度

`go build -gcflags '-m'`

#### 超过栈帧(stack frame)
- 当一个函数被调用时，会在两个相关的帧边界间进行上下文切换。
  - 从调用函数切换到被调用函数，如果函数调用时需要传递参数，
    那么这些参数值也要传递到被调用函数的帧边界中。
  - Go语言中帧边界间的数据传递是按值传递的。
  - 任何在函数 getRandom 中的变量在函数返回时，都将不能访问。
  - Go 查找所有变量超过当前函数栈侦的，把它们分配到堆上，避免 outlive 变量。

- 超过栈帧(stack frame) 上述情况中，num 变量不能指向之前的栈。
  - Go 查找所有变量超过当前函数栈侦的，把它们分配到堆上，避免 outlive 变量。
  
- 变量 tmp 在栈上分配，但是它包含了指向堆内存的地址，所以可以安全的
  从一个函数的栈帧复制到另外一个函数的栈帧。

#### 逃逸案例
 - 还存在大量其他的 case 会出现逃逸，
   比较典型的就是 “多级间接赋值容易导致逃逸”，这里的多级间接指的是，
   对某个引用类对象中的引用类成员进行赋值
 - 记住公式 Data.Field = Value，如果 Data, Field 都是引用类的数据类型，
   则会导致 Value 逃逸。这里的等号 = 不单单只赋值，也表示参数传递。
 - Go 语言中的引用类数据类型有 func, interface, slice, map, chan, *Type 
 - 一个值被分享到函数栈帧范围之外
 - 在 for 循环外申明，在 for 循环内分配，同理闭包
 - 发送指针或者带有指针的值到 channel 中
 - 在一个切片上存储指针或带指针的值
 - slice 的背后数组被重新分配了
 - 在 interface 类型上调用方法 
   ....
   
`go build -gcflags '-m'`

### 分段栈(Segmented stacks)
- Go 应用程序运行时，每个 goroutine 都维护着一 个自己的栈区，
  这个栈区只能自己使用不能被其他 goroutine 使用。
- 栈区的初始大小是2KB(比 x86_64 架构下线程的默认栈2M 要小很多)，
  在 goroutine 运行的时候栈区会按照需要增⻓和收缩， 
  占用的内存最大限制的默认值在64位系统上是 1GB。 
  - v1.0 ~ v1.1 — 最小栈内存空间为 4KB 
  - v1.2 — 将最小栈内存提升到了 8KB
  - v1.3 — 使用连续栈替换之前版本的分段栈 
  - v1.4 — 将最小栈内存降低到了 2KB

#### Hot split 问题
- 分段栈的实现方式存在 “hot split” 问题，
  - 如果栈快满了，那么下一次的函数调用会强制触发栈扩容。
- 当函数返回时，新分配的 “stack chunk” 会被清理掉。
- 如果这个函数调用产生的范围是在一个循环中，会导致严重的性能问题，
  频繁的 alloc/free。
- Go 不得不在1.2版本把栈默认大小改为8KB， 降低触发热分裂的问题，
  但是每个 goroutine 内存开销就比较大了。
- 直到实现了连续栈 (contiguous stack)，栈大小才改为2KB。

### 连续栈(Contiguous stacks)
- 采用复制栈的实现方式，在热分裂场景中不会频发释放内存，
  即不像分配一个新的内存块并链接到老的栈内存块，
- 而是会**分配一个两倍大的内存块并把老的内存块内容复制到新的内存块**里，
  当栈缩减回之前大小时，我们不需要做任何事情。
- runtime.newstack 分配更大的栈内存空间 
- runtime.copystack 将旧栈中的内容复制到新栈中
- 将指向旧栈对应变量的指针重新指向新栈 
- runtime.stackfree 销毁并回收旧栈的内存空间
- 如果栈区的空间使用率不超过1/4，那么在垃圾回收的时候使用 runtime.shrinkstack 进行栈缩容，
  同样使用 copystack

### 栈扩容
- Go 运行时判断栈空间是否足够，所以在 call function 中会插 入 runtime.morestack，
- 每个函数调用都判定的话，成本比较高。
- 在编译期间通过计算 sp、func stack framesize 确定需要哪个函数调用中插入
  runtime.morestack。
  
- 当函数是叶子节点，且栈帧小于等于 112 
  - 不插入指令
- 插入
   - 当叶子函数栈帧大小为 120 -128 或者 非叶子函数栈帧大小为 0-128， SP < stackguard0 
   - 当函数栈帧大小为 128 - 4096, `SP - framesize < stackguard0 - StackSmall`
   - 大于 StackBig, `SP-stackguard+StackGuard <= framesize + (StackGuard- StackSmall)`


## 内存管理
- TCMalloc 是 Thread Cache Malloc 的简称，
- 是Go 内 存管理的起源，Go的内存管理是借鉴了TCMalloc:

### 内存碎片
- 随着内存不断的申请和释放，内存上会存在大量的碎片， 降低内存的使用率。
- 为了解决内存碎片，可以将2个连续的 未使用的内存块合并，减少碎片。

### 大锁 (全局锁）
- 同一进程下的所有线程共享相同的内存空间，它们申请内存时需要加锁，
- 如果不加锁就存在同一块内存被2个线程同时访问的问题。

### 内存布局
几个重要的概念:
- page: 内存⻚
  - 一块 8K 大小的内存空间。
  - Go 与操作系统 之间的内存申请和释放，都是以 page 为单位的。
- span: 内存块，
  - 一个或多个连续的 page 组成一个 span。 
- sizeclass: 空间规格，
  - 每个 span 都带有一个 sizeclass，标记着该 span 中的 page 应该如何使用。
- object: 对象，
  - 用来存储一个变量数据内存空间，一个 span 在初始化时，
    会被切割成一堆等大的 object。
  - 假设 object 的大小是 16B，span 大小是 8K，
    那么就会把 span 中的 page 就会被初始化 8K / 16B = 512 个 object。

### 小于 32kb 内存分配

#### mcache
- 当程序里发生了 32kb 以下的小块内存申请时，
- Go 会从一个叫做的 mcache 的本地缓存给程序分配内存。
- 这样的一个内存块里叫做 mspan，它是要给程序分配内存时的分配单元。
- 在 Go 的调度器模型里，每个线程 M 会绑定给一个处理 器 P，
  - 在单一粒度的时间里只能做多处理运行一个 goroutine，
  - 每个 P 都会绑定一个上面说的本地缓存 mcache。
- 当需要进行内存分配时，当前运行的 goroutine 会从 mcache 中查找可用的 mspan。
- 从本地 mcache 里 分配内存时不需要加锁，这种分配策略效率更高。

#### mspan
- 申请内存时都分给他们一个 mspan 这样的单元会 不会产生浪费。
  - 其实 mcache 持有的这一系列的 mspan 并不都是统一大小的，
  - 而是按照大小，从 8kb 到 32kb 分了大概 67*2 类的 mspan。
- 每个内存⻚分为多级固定大小的“空闲列表”，这有助于 减少碎片。
  - 类似的思路在 Linux Kernel、Memcache 都以⻅到 Slab-Allactor。
    
#### mcentral 
- 如果分配内存时 mcachce 里没有空闲的对口 sizeclass 的 mspan 了
  - Go 里还为每种类别的 mspan 维护着一个 mcentral。
- mcentral 的作用是为所有 mcache 提供切分好的 mspan 资源。
  - 每个 central 会持有一种特定大小的全局 mspan 列表， 包括已分配出去的 和未分配出去的。 
  - 每个 mcentral 对应一种 mspan，
    - 当工作线程的 mcache 中没有合适(也就是特定大小的)的mspan 
      时就会从 mcentral 去获取。 
- mcentral 被所有的工作线程共同享有，存在多个 goroutine 竞争的情 况，
  - 因此从 mcentral 获取资源时需要加锁。
- mcentral 里维护着两个双向链表，
  - nonempty list 表示链表里还有空闲的 mspan 待分配。
  - empty list 表示这条链表里的 mspan 都被分配了object 或缓存 mcache 中。

- 程序申请内存的时候，mcache 里已经没有合适的空闲 mspan了，
  那么工作线程就会像下图这样去 mcentral 里 去申请。
- mcache 从 mcentral 获取和归还 mspan 的流 程:
  - 获取 加锁;
    - 从 nonempty 链表找到一个可用的 mspan;
    - 并将其从 nonempty 链表删除;
    - 将取出的 mspan 加入到 empty 链表;
    - 将 mspan 返回给工作线程;
    - 解锁。
  - 归还 加锁;
    - 将 mspan 从 empty 链表删除;
    - 将mspan 加入到 nonempty 链表;
    - 解锁。
- mcentral 是 sizeclass 相同的 span 会以链表的形式组织 在一起, 就是指该 span 用来存储哪种大小的对象。

#### mheap
- 当 mcentral 没有空闲的 mspan 时，会向 mheap 申请。
- 而 mheap 没有资源时，会向操作系统申 请新内存。
- mheap 主要用于大对象的内存分配， 以及管理未切割的 mspan，
  用于给 mcentral 切 割成小对象。
- mheap 中含有所有规格的 mcentral，
  - 所以当一 个 mcache 从 mcentral 申请 mspan 时，
    只需要 在独立的 mcentral 中使用锁，
    并不会影响申请其 他规格的 mspan。

#### arena 区域
- 所有 mcentral 的集合则是存放于 mheap 中的。
- mheap 里的 arena 区域是真正的堆区，
  运行时会将 8KB 看做一⻚，这些内存⻚中存储了所有在堆上初始化的对象。
- 运行时使用二维的 runtime.heapArena 数组管理所有的内存，
  - 每个 runtime.heapArena 都会管理 64MB 的内存。
- 如果 arena 区域没有足够的空间，会调用 runtime.mheap.sysAlloc 
  从操作系统中申请更多的内存。(如下图:Go 1.11 前的内存布局)
- 1.11 最大是512G
- 1.11+，对arena结构修改，使得可以支持到256T


#### 小于 16b 内存分配
- 对于小于16字节的对象(且无指针)，Go 语言将其划分 为了tiny 对象。
- 划分 tiny 对象的主要目的是为了处理 极小的字符串和独立的转义变量。
- 对 json 的基准测试 表明，使用 tiny 对象减少了12%的分配次数和20%的堆大小。
- tiny 对象会被放入class 为2的 span 中。

分配方法  
- 首先查看之前分配的元素中是否有空余的空间
- 如果当前要分配的大小不够，例如要分配16字节的大 小，这时就需要找到下一个空闲的元素
- tiny 分配的第一步是尝试利用分配过的前一个元素的空间， 达到节约内存的目的。

#### 大于 32kb 内存分配
- Go 没法使用工作线程的本地缓存 mcache 和全局中心缓存 mcentral 
  上管理超过32KB 的内存分配，所以对于那些超过32KB的内 存申请，
  会直接从堆上(mheap)上分配对应的数量的内存⻚(每⻚大小是8KB)给程序。
- 实现数据结构（每版不一）  
 - freelist
 - treap
 - radix tree + pagecache

### 内存分配总结
- 一般小对象通过 mspan 分配内存
- 大对象则直接由 mheap 分配内存。
- Go 在程序启动时，会向操作系统申请一大块内存，由 mheap 结构全局管理
  - 现在 Go 版本不需要连续地址了，所以不会申请一大堆地址
- Go 内存管理的基本单元是 mspan，每种 mspan 可以分配特定大小的 object
- mcache, mcentral, mheap 是 Go 内存管理的三大组件，
  - mcache 管理线程在本地缓存的 mspan
  - mcentral 管理全局的 mspan 供所有线程

## GC 原理
- Mark & Sweep
- Tri-color Mark & Sweep 
- Write Barrier
- Stop The World

### Garbage Collection
- 现代高级编程语言管理内存的方式分为两种: 自动和手动，
  - 像 C、C++ 等编程语言使用手动管理内存的方式，
    工程师编写代码过程中需要主动申请或者释放内存;
  - 而 PHP、Java 和 Go 等语言使用自动的内存管理系统，
    有内存分配器和垃圾收集器来代为分配和回收内存，
- 其中垃圾收集器就是我们常说的 GC。
  
- 主流的垃圾回收算法:
  - 引用计数
  - 追踪式垃圾回收
- Go 现在用的三色标记法就属于追踪式垃圾回收算法的一种。

### Mark & Sweep
#### Mark-Sweep 概念
- Mark Sweep 两个阶段
  - 标记(Mark)和 
  - 清除(Sweep)两个 阶段，
  - 所以也叫 Mark-Sweep 垃圾回收算法。
  
- STW stop the world
  - GC 的一些阶段需要停止所有的 mutator 以确定当前的引用关系。
  - 这便是很多人对GC担心的来源，这也是 GC 算法优化的重点。
    
- Root
  - 根对象是 mutator 不需要通过其他对象就可以直接访问到 的对象。
     - 比如全局对象，栈对象中的数据等。
  - 通过Root 对象，可以追踪到其他存活的对象。

- Mark & Sweep 这个算法就是严格按照追踪式算法的思路来实现的:
  - Stop the World
  - Mark:通过 Root 和 Root 直接间接访问到的对象， 来
  - 寻找所有可达的对象，并进行标记。 Sweep:对堆对象迭代，已标记的对象置位标记。所有
  - 未标记的对象加入freelist， 可用于再分配。 
  - Start the Wrold
  
- 这个算法最大的问题是 GC 执行期间需要把整个程序完全 暂停，
  - 朴素的 Mark Sweep 是整体 STW，并且分配速度慢，内存碎片率高。
  - Go v1.1，STW 可能秒级

#### gov1.3 并发 Sweep
- 标记过程需的要 STW，因为对象引用关系如果在 标记阶段做了修改，会影响标记结果的正确性。
- 并发 GC 分为两层含义:
  - 每个 mark 或 sweep 本身是多个线程(协程)执行的 (concurrent)
  -  mutator 和 collector 同时运行(background)
- concurrent 这一层是比较好实现的, GC 时整体进行 STW，那么对象引用关系不会再改变，
  对 mark 或者 sweep 任务进行分块，就能多个线程(协程) conncurrent 
  执行任务 mark 或 sweep。
  - Go v1.3，标记 STW，并发 Sweep
- 而对于 backgroud 这一层, 也就是说 mutator 和 mark，sweep 同时运行，则相对复杂。
  - 1.3以前的版本使用标记-清扫的方式，整个过程都需要 STW。 
  - 1.3版本分离了标记和清扫的操作，标记过程STW，清扫过程并发执行。
- backgroup sweep 是比较容易实现的，因为 mark 后，哪些对象是存活， 
  哪些是要被 sweep 是已知的，sweep 的是不再引用的对象。sweep 结束前，
  这些对象不会再被分配到，所以 sweep 和 mutator 运行共存。
  无论 全局还是栈不可能能访问的到这些对象，可以安全清理。
#### gov1.5 三色标记法   
- 1.5版本在标记过程中使用三色标记法。标记和清扫都并发执行的，
- 但标 记阶段的前后需要 STW 一定时间来做 GC 的准备工作和栈的re-scan。

### Tri-color Mark & Sweep

#### 三色标记法
- 三色标记是对标记清除法的改进，标记清除法在整个执行时要求⻓时间 STW，
  Go 从1.5版本开始改为三色标记法，
- 初始将所有内存标记为白色，然后将 roots 加入待扫描队列(进入队列即被视为变成灰色)，
  然后使用并发 goroutine 扫描队列中的指针，如果指针还引用 了其他指针，
  那么被引用的也进入队列，被扫描的对象视为黑色。
- 三色标记 
  - 白色对象:潜在的垃圾，其内存可能会被垃圾收集器回收。
  - 黑色对象:活跃的对象，
    包括不存在任何引用外部指针的对象以及从根对象可达的对象，
    垃圾回收器不会扫描这些对象的子对象。
  - 灰色对象 :活跃的对象，因为存在指向白色对象的外部指针，
    垃圾收集器会扫描这些对象的子对象。
-  Go v1.5，并行 Mark & Sweep

#### Tri-color Marking
- 垃圾收集器从 root 开始然后跟随指针递归整个内存空间。
  分配于 noscan 的 span 的对象, 不会进行扫描。
-  然而，此过程不是由同一个 goroutine 完成的，每个指针 都排队在工作池中 
-  然后，先看到的被标记为工作协程 的后台协程从该池中出队，扫描对象，
-  然后将在其中找 到的指针排入队列。

#### Tri-color Coloring
- 染色流程:
  -  一开始所有对象被认为是白色
  -  根节点(stacks，heap，global variables)被染 色为灰色 
  - 一旦主流程走完，gc 会:
   - 选一个灰色对象，标记为黑色
   - 遍历这个对象的所有指针，标记所有其引用 的对象为灰色
   - 最终直到所有对象需要被染色。

- 标记结束后，黑色对象是内存中正在使用的对象，而白色对象是要收集的对象。
- 由于 struct2的实例是在匿名函数中创建的，并且 无法从堆栈访问，因此它保持为白色，可以清除。
- 颜色在内部实现原理:
  - 每个 span 中有一个名为 gcmarkBits 的位图属性，该属性跟踪扫描，并将相应的位设置为1。

### Write Barrier
- 1.5版本在标记过程中使用三色标记法。
- 回收过程主要 有四个阶段，其中，标记和清扫都并发执行的，
  但标记阶段的前后需要 STW 一定时间来做GC 的准备工作和 栈的 re-scan。
- 使用并发的垃圾回收，也就是多个 Mutator 与 Mark 并发执 行，
  想要在并发或者增量的标记算法中保证正确性，
  我们需要 达成以下两种三色不变性(Tri-color invariant)中的任意一种:
  - 强三色不变性:
    - 黑色对象不会指向白色对象，只会指向灰色对象或者黑色对象。
  - 弱三色不变性:
    - 黑色对象指向的白色对象必须包含一条从灰色对象经由多个白色对象的可达路径。

#### 对象丢失问题
- 可以看出，一个白色对象被黑色对象引用，是注定无法通过
  这个黑色对象来保证自身存活的，与此同时，如果所有能到 
  达它的灰色对象与它之间的可达关系全部遭到破坏，那么这 
  个白色对象必然会被视为垃圾清除掉。 
- 故当上述两个条件同时满足时，就会出现对象丢失的问题。
  如果这个白色对象 下游还引用了其他对象，并且这条路径是指向下游对象的唯 一路径，
  那么他们也是必死无疑的。
  
- 为了防止这种现象的发生，最简单的方式就是 STW，
  直接禁止掉其他用户程序对对象引用关系的干扰，
  但是 STW 的过程有明显的资源浪费，对所有的用户程序都有很大影响，
  如何能在保证对象不丢失的情况下合理的尽可能的提高 GC 效率，减少 STW 时间呢?

#### Write Barrier - Dijkstra 写屏障

- 插入屏障拦截将白色指针插入黑色对象的操作，标记其对应对象为灰色状态，
  这样就不存在黑色 对象引用白色对象的情况了，满足强三色不变式， 
  在插入指针 f 时将 C 对象标记为灰色。
- 如果对栈上的写做拦截，那么流程代码会非常复杂，并且性能下降会非常大，得不偿失。
  根据局部性的原理来说，其实我们程序跑起来，大部分的其实都是操作在栈上，
  函数参数啊、函数调用导致的压栈出栈、局部变量啊，协程栈，这些如果也弄起写屏障，
  那么可想而知了，根本就不现实，复杂度和性能就是越不过去的坎。
- Go1.5版本使用的 Dijkstra 写屏障就是这个原理

- 所以 Go 选择仅对堆上的指针插入增加写屏障，这样就会出现在扫描结束后，
  栈上仍存在引用白色对象的情况，这时的栈是灰色的，不满足三色不变式，
  所以需要对栈进行重新扫描使其变黑，完成剩余对象的标记，这个过程 需要 STW。

- 初始化 GC 任务，包括
  - 开启写屏障(write barrier)和开启辅助 GC(mutator assist)，
    统计 root 对象的任务数量等，这个过程需要 STW。
- 扫描所有 root 对象，包括全局指针和 goroutine(G) 
  栈上的指针(扫描对应 G 栈时需停止该 G)，将其加入标记队列(灰色队列)，
  并循环处理灰色队列的对 象，直到灰色队列为空，该过程后台并行执行。
- 完成标记工作，重新扫描(re-scan)全局指针和栈。
- 因为 Mark 和 mutator 是并 行的，所以在 Mark 过程中可能会有新的对象分配和指针赋值，
  这个时候就需 要通过写屏障(write barrier)记录下来，
  re-scan 再检查一下，这个过程也是会 STW 的。
- 按照标记结果回收所有的白色对象，该过程后台并行执行。

#### Write Barrier - Yuasa 删屏障
- 删除屏障也是拦截写操作的，但是是通过保护灰色对象到白色对象的路径不会断来实现的。
- 如上图例中，在删除指针 e 时将对象 C 标记为灰色，
  这样 C 下游的所有白色对象，即使会被黑色对象 引用，
  最终也还是会被扫描标记的，满足了弱三 色不变式。
- 这种方式的回收精度低，一个对象即 使被删除了最后一个指向它的指针也依旧可以活过这一轮，
  在下一轮 GC 中被清理掉。

#### Write Barrier - 混合屏障
- 插入屏障和删除屏障各有优缺点，Dijkstra 的插入写屏障在标记开始时无需 STW，
  可直接开 始，并发进行，但结束时需要 STW 来重新扫描 栈，
  标记栈上引用的白色对象的存活;
- Yuasa 的 删除写屏障则需要在 GC 开始时 STW 扫描堆栈 来记录初始快照，
  这个过程会保护开始时刻的所 有存活对象，但结束时无需 STW。
- Golang 中的混合写屏障满足的是变形的弱三色不变式， 同样允许黑色对象引用白色对象，
  白色对象处于灰色保 护状态，但是只由堆上的灰色对象保护。
- Go1.8 混合写屏障结合了Yuasa的删除写屏障和Dijkstra的插入写屏障

#### Write Barrier - 混合屏障
- 由于结合了 Yuasa 的删除写屏障和 Dijkstra 的插入 写屏障的优点，
  只需要在开始时并发扫描各个 goroutine 的栈，使其变黑并一直保持，
  这个过程不 需要 STW，而标记结束后，因为栈在扫描后始终是 黑色的，
  也无需再进行 re-scan 操作了，减少了 STW 的时间。
- 为了移除栈的重扫描过程，除了引入混合写屏障之外，在 垃圾收集的标记阶段，
  我们还需要将创建的所有堆上新对 象都标记成黑色，防止新分配的栈内存和堆内存中的对象
  被错误地回收，因为栈内存在标记阶段最终都会变为黑色，
  所以不再需要重新扫描栈空间。

### Sweep
- Sweep 让 Go 知道哪些内存可以重新分配使用， 
  然而，Sweep 过程并不会处理释放的对象内存置 为0(zeroing the memory)。
- 而是在分配重新使用 的时候，重新 reset bit。
- 每个 span 内有一个 bitmap allocBits，他表示上 一次 GC 之后每一个 object 的分配情况，
  1:表 示已分配，0:表示未使用或释放。
- 内部还使用了 uint64 allocCache(deBruijn)，加速寻找 freeobject。

- GC 将会启动去释放不再被使用的内存。
  在标记期 间，GC 会用一个位图 gcmarkBits 来跟踪在使用 中的内存。
- 正在被使用的内存被标记为黑色，然而当前执行并不能 够到达的那些内存会保持为白色。
- 现在，我们可以使用 gcmarkBits 精确查看可用于分配的 内存。
- Go 使用 gcmarkBits 赋值了 allocBits，这个操作 就是内存清理。
- 然而必须每个 span 都来一次类似的处理，需要耗费大 量时间。
- Go 的目标是在清理内存时不阻碍执行，并为 此提供了两种策略。

- Go 提供两种方式来清理内存:
  - 在后台启动一个 worker 等待清理内存，一个一个 mspan 处理
    - 当开始运行程序时，Go 将设置一个后台运行的 Worker(唯一的任务就是去清理内存)，
    - 它将进入睡眠状态并等待内存段扫描。
  - 当申请分配内存时候 lazy 触发
    - 当应用程序 goroutine 尝试在堆内存中分配新内存时，会触发该操作。
    - 清理导致的延迟和吞吐量降低被分散到 每次内存分配时。

- 清理内存段的第二种方式是即时执行。
  - 但是，由于这些 内存段已经被分发到每一个处理器 P 的本地缓存 mcache 中，
    因此很难追踪首先清理哪些内存。
  - 这就是 为什么 Go 首先将所有内存段移动到 mcentral 的原因。 
  - 然后，它将会让本地缓存 mcache 再次请求它们，去即 时清理。
-  即时扫描确保所有内存段都会得到清理(节省资源)，同时不 会阻塞程序执行。
- 由于后台只有一个 worker 在清理内存块，清理过程可能会花费 一些时间。
  - 但是，我们可能想知道如果另一个 GC 周期在一次 清理过程中启动会发生什么。
  - 在这种情况下，这个运行 GC 的 Goroutine 就会在开始标记阶段前去协助完成剩余的清理工作。



### STW
- 在垃圾回收机制 (GC) 中，"Stop the World" (STW) 是一个重要阶段。
- 顾名思义， 在 "Stop the World" 阶段， 当前运行的所有程序将被暂 停， 
  - 扫描内存的 root 节点和添加写屏障 (write barrier) 。
- 这个阶段的第一步， 是抢占所有正在运行的 goroutine，
- 被抢占之后， 这些 goroutine 会被悬停在 一个相对安全的状态。

- 处理器 P (无论是正在运行代码的处理器还是已在 idle 列表中的处 理器)， 
  都会被被标记成停止状态 (stopped)， 不再运行任何代码。 
- 调度器把每个处理器的 M 从各自对应的处理器 P 分离出 来， 放到 idle 列表中去。
- 对于 Goroutine 本身， 他们会被放到一个全局队列中等待。

### Pacing
- 运行时中有 GC Percentage 的配置选项，默认情况下为 100。
- 此值表示在下一次垃圾收集必须启动之前可以分配多少新内存的比率。
  - 将 GC 百分比设置为100意味着，基于在垃圾收集完成后标记为 活动的堆内存量，
    下次垃圾收集前，堆内存使用可以增加100%。
- 如果超过2分钟没有触发，会强制触发 GC。
- 使用环境变量 GODEBUG 和 gctrace = 1选项生成GC trace
  - `GODEBUG=gctrace=1 ./app`

## References
- https://medium.com/a-journey-with-go/go-goroutine-os-thread-and-cpu-management-2f5a5eaf518a
- http://www.sizeofvoid.net/goroutine-under-the-hood/ 
- https://zhuanlan.zhihu.com/p/84591715 
- https://rakyll.org/scheduler/ 
- https://zhuanlan.zhihu.com/p/248697371 
- https://zhuanlan.zhihu.com/p/68299348 
- https://blog.csdn.net/qq_25504271/article/details/81000217 
- https://blog.csdn.net/ABo_Zhang/article/details/90106910 
- https://zhuanlan.zhihu.com/p/66090420 
- https://zhuanlan.zhihu.com/p/27056944 
- https://www.cnblogs.com/sunsky303/p/11058728.html 
- https://www.cnblogs.com/zkweb/p/7815600.html
- https://yizhi.ren/2019/06/03/goscheduler/
- https://morsmachine.dk/netpoller
- https://segmentfault.com/a/1190000022030353?utm_source=sf-related 
- https://www.jianshu.com/p/0083a90a8f7e 
- https://www.jianshu.com/p/1ffde2de153f 
- https://www.jianshu.com/p/63404461e520 
- https://www.jianshu.com/p/7405b4e11ee2 
- https://www.jianshu.com/p/518466b4ee96 
- https://zhuanlan.zhihu.com/p/59125443 
- https://www.codercto.com/a/116486.html 
- https://www.jianshu.com/p/db0aea4d60ed 
- https://www.jianshu.com/p/ef654413f2c1 
- https://zhuanlan.zhihu.com/p/248697371
- https://medium.com/a-journey-with-go/go-how-does-a-goroutine-start-and-exit-2b3303890452
- https://medium.com/a-journey-with-go/go-g0-special-goroutine-8c778c6704d8
- https://medium.com/a-journey-with-go/go-how-does-go-recycle-goroutines-f047a79ab352
- https://medium.com/a-journey-with-go/go-what-does-a-goroutine-switch-actually-involve-394c202dddb7
- http://xiaorui.cc/archives/6535 http://xiaorui.cc/archives/category/golang
- https://docs.google.com/document/d/1lyPIbmsYbXnpNj57a261hgOYVpNRcgydurVQIyZOz_o/pub
- https://medium.com/a-journey-with-go/go-asynchronous-preemption-b5194227371c
- https://medium.com/a-journey-with-go/go-goroutine-and-preemption-d6bc2aa2f4b7 
- http://xiaorui.cc/archives/6535 
- https://medium.com/a-journey-with-go/go-gsignal-master-of-signals-329f7ff39391 
- https://www.jianshu.com/p/1ffde2de153f
- https://kirk91.github.io/posts/2d571d09/ 
- http://yangxikun.github.io/golang/2019/11/12/go-goroutine-stack.html
- https://www.ardanlabs.com/blog/2017/05/language-mechanics-on-stacks-and-pointers.html
- https://www.ardanlabs.com/blog/2017/05/language-mechanics-on-escape-analysis.html
- https://zhuanlan.zhihu.com/p/237870981
- https://www.ardanlabs.com/blog/2017/05/language-mechanics-on-stacks-and-pointers.html
- https://blog.csdn.net/qq_35587463/article/details/104221280 
- https://www.jianshu.com/p/63404461e520
- https://www.do1618.com/archives/1328/go-%E5%86%85%E5%AD%98%E9%80%83%E9%80%B8%E8%AF%A6%E7%BB%86%E5%88%86%E6%9E%90/
- https://www.jianshu.com/p/518466b4ee96 
- https://zhuanlan.zhihu.com/p/28484133 
- http://yangxikun.github.io/golang/2019/11/12/go-goroutine-stack.html 
- https://kirk91.github.io/posts/2d571d09/ 
- https://zhuanlan.zhihu.com/p/237870981 
- https://agis.io/post/contiguous-stacks-golang/
- https://docs.google.com/document/d/13v_u3UrN2pgUtPnH4y-qfmlXwEEryikFu0SQiwk35SA/pub
- https://docs.google.com/document/d/1lyPIbmsYbXnpNj57a261hgOYVpNRcgydurVQIyZOz_o/pub
- https://zhuanlan.zhihu.com/p/266496735 
- http://dmitrysoshnikov.com/compilers/writing-a-memory-allocator/ 
- https://studygolang.com/articles/22652?fr=sidebar 
- https://studygolang.com/articles/22500?fr=sidebar 
- https://www.cnblogs.com/unqiang/p/12052308.html
- https://blog.csdn.net/weixin_33869377/article/details/89801587
- https://www.cnblogs.com/smallJunJun/p/11913750.html
- https://zhuanlan.zhihu.com/p/53581298
- https://zhuanlan.zhihu.com/p/141908054
- https://zhuanlan.zhihu.com/p/143573649
- https://zhuanlan.zhihu.com/p/145205154
- https://www.jianshu.com/p/47735dfb0b81
- https://zhuanlan.zhihu.com/p/266496735
- https://dave.cheney.net/high-performance-go-workshop/dotgo-paris.html#memory-and-gc
- https://juejin.im/post/6844903917650722829
- https://spin.atomicobject.com/2014/09/03/visualizing-garbage-collection-algorithms/
- https://zhuanlan.zhihu.com/p/245214547
- https://www.jianshu.com/p/2f94e9364ec4
- https://www.jianshu.com/p/ebd8b012572e
- https://www.ardanlabs.com/blog/2018/12/garbage-collection-in-go-part1-semantics.html
- https://segmentfault.com/a/1190000012597428 
- https://www.jianshu.com/p/bfc3c65c05d1
- https://golang.design/under-the-hood/zh-cn/part2runtime/ch08gc/sweep/ 
- https://zhuanlan.zhihu.com/p/74853110 
- https://www.jianshu.com/p/2f94e9364ec4 
- https://juejin.im/post/6844903917650722829 
- https://zhuanlan.zhihu.com/p/74853110
- https://www.jianshu.com/p/ebd8b012572e
- https://www.jianshu.com/p/2f94e9364ec4
- https://www.jianshu.com/p/bfc3c65c05d1
- https://zhuanlan.zhihu.com/p/92210761
- https://blog.csdn.net/u010853261/article/details/102945046
- https://blog.csdn.net/hello_bravo_/article/details/103840054
- https://segmentfault.com/a/1190000020086769
- https://blog.csdn.net/cyq6239075/article/details/106412038
- https://zhuanlan.zhihu.com/p/77943973
- https://www.ardanlabs.com/blog/2018/12/garbage-collection-in-go-part1-semantics.html
- https://www.ardanlabs.com/blog/2019/05/garbage-collection-in-go-part2-gctraces.html
- https://www.ardanlabs.com/blog/2019/07/garbage-collection-in-go-part3-gcpacing.html
- https://github.com/dgraph-io/badger/tree/master/skl
- https://dgraph.io/blog/post/manual-memory-management-golang-jemalloc/
