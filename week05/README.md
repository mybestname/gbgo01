# 第5周：微服务可用性设计

## 1 隔离
- 服务隔离
  +  动静分离、读写分离
- 轻重隔离
  + 核心与非核心、快慢、热点
- 物理隔离
  + 线程、进程、集群、机房
### 服务隔离 

#### 动静隔离
- CPU cache
  + cache line 的 false sharing
    - CPU是以cache line 为单位存储的，当多线程修改互相独立的变量时，如果这些变量共享同一个cache line，就会无意中影响彼此的性能，这就是伪共享 false sharing。
    - cache line上的写竞争在特定情况下是并行线程性能的一个非常重要限制因素。
    - 从代码层面中很难看出是否会出现false sharing，这个涉及语言实现的底层和硬件。
    - 解决的办法是padding和对齐（align），让数据对象（语言级别）处于不同的cache line（CPU级别）。
  > 这个问题的本质是语言级别的数据结构的设计和硬件之间由于抽象不同而造成矛盾。
  > - vs. 游戏开发上面所谓的Data Oriented Design，其实说的也是这个问题。
  > - 不是说建立抽象层不好（vs. 经典的SICP书中的思想)，计算科学解决任何问题本身就是基于分层和抽象。
  > - 任何数据结构本身都基于一层或多层的抽象，而抽象必然带来性能损失。关键还是要看要解决什么问题。
- MySQL的buffer pool（`innodb_buffer_pool_size`）
  + 本质上是Mysql向OS申请的一段连续内存空间, 链表的数据结构, 优化的LRU算法（Least recently used）
  + flush 链表，存储脏页的链表，把修改过的缓存页数据（脏页）加入到一个链表中，在未来的某个时间点进行同步，而不是频繁写磁盘。  
  + 缓冲池污染：某一个SQL语句，要批量扫描大量数据时，可能导致把缓冲池的所有页都替换出去，导致大量热数据被替换出，造成MySQL性能急剧下降。
  + 在表设计上要避免buffer pool的频繁过期。
     - 经常更新的表（动表），不更新的表（静表）。两者避免放在一个buffer pool里面，避免更新造成从buffer pool中不必要的剔除过程。
    > - 这里所谓的隔离的buffer pool那么也就是分库的意思了。 
    > - 因为并不存在可以把一个库的不同表存在不同的buffer pool里面。
    > - buffer pool是存储引擎的缓冲池，针对的是数据读写的底层操作。不可能基于表进行控制。
    > - 但是分库本身的cache效果是否就更好，这里存疑。
    > - 这里的意思应该是指：本来有关系（业务层面）的一张表，分开为两张表（没有关系，sql层面），这样更新操作在缓存层面就被隔离了。
- CDN的场景
  + 静态资源（图片，css等）请求和动态API请求进行分离。
  + 降低应用服务器负载，静态文件负载都走CDN

#### 读写分离
  + 主从（master/salve之间的读写分离）
  + Replicaset （radis）
    > 这个名词用的不好，和k8s混淆了，而k8s的概念（指pod的副本）和这里完全不同。
  + CQRS（Command Query Responsibility Segregation）(https://martinfowler.com/bliki/CQRS.html) (Martin fowler 2011)
  
### 轻重隔离

- 核心业务和非核心业务的隔离
- 快慢：如果把数据流想象为水流，那么流量吞吐能力不同，就会有快慢之分
  - 例子：kafka的同一topic，不同sink端的速度不同带来的问题
    - Source负责导入数据到Kafka，Sink负责从Kafka导出数据，针对相同Topic的多个Sink，如果速度（吞）不一，必然影响上游吞吐。
    - 解决：
      - sink端隔离：
        - 建consumer Group
      - topic拆分：业务线
        - topic变多的问题是顺序IO的好处没有了。
- 热点： 经常访问的数据点，或者业务突然密集的时间点（秒杀，热门直播）
  + 这种场景下，对redis cache可能单点打爆
  + 解决办法：将redis cache 提升到local cache
  + 高频访问，但是内容不怎么变的场景
      + 服务启动时伴随一个goroutine一直从数据库 polling数据
      + 使用`atomic.Value` ，不用`sync.Map`原因是不能同时读写，不如`atomic.Value`进行原子替换。可以CoW无锁访问。
  + 大量在线客户端大量刷新的场景
    + 主动预热（即通过监控来主动防御）
    + 通过监控nginx的live-streaming情况，来通知服务，要求服务将cache提升为local cache，来防止cache穿透。
    + 更好的方式：不用监控，通过进程级别自主根据服务访问频次主动提升为local cache，通过将功能集成到基础库中，来提升整体可用性。
   
### 物理隔离

- 线程隔离
  + 线程隔离指线程池分业务，不同业务不同线程池管理，当某个某业务出问题，故障只限于本线程池。
  + Java的线程池中的线程耗尽 (https://mp.weixin.qq.com/s/PmU14UsJOb4IiH_81RlJMA)
  + go和java的区别，go不需要担心线程池，只需要考虑控制goroutine总量。
  + Java的解决办法 
    - [Netflix/Hystrix](https://github.com/Netflix/Hystrix/blob/master/hystrix-core/src/main/java/com/netflix/hystrix/)
    - [resilience4j](https://github.com/resilience4j/resilience4j)  
    - 熔断器：配置值进行控制是否进行熔断。
    - 基于信号量(semaphore)：获得信号的可以访问，完成后返回信号。通过信号量总数有限控制访问。
    - 问题：不论那种方式，都需要手工配置值。一麻烦，二值是难以手工确定的，多少合适呢？
    - 更好的方式：自适应，不需要手工设定。
- 进程隔离
  + docker
  + k8s
  + kvm
  + yarn

- 集群隔离
  - 逻辑1，物理多。
  - 物理机房隔离（多活）

### 基于隔离的案例
 - 早期转码集群被超大视频攻击，导致转码大量延迟。
   + 解决：按视频规格，重要性等指标走不同集群，使得影响被隔离。
 - 缩略图服务，被大图（GIF）实时缩略吃完所有 CPU，导致正常的小图缩略被丢弃，大量503。
   + 解决：把图片按照种类和规格分隔为不同集群，大图走特殊集群。全局故障变局部故障。
 - 数据库实例 cgroup 未隔离，导致大 SQL 引起的集体故障。（虚拟机和物理机的隔离）
   + 解决：加cgroup，通过CGroup进行CPU、内存等资源控制
 - INFO 日志量过大，导致异常 ERROR 日志采集延迟。

## 2 超时控制

- 内网服务要求100ms （最多不能超300ms）
- 公网服务不能超1s
- 注意超时叠加：调用是互相叠加的。
- 注意网络传递的不确定性。
- 注意c/s两端由于**超时策略不一致造成的资源浪费** ，例如客户端设定为100ms超时，服务端设定为500ms，那么服务请求对于客户端已经失败，但这个调用还被服务端继续执行。
- 默认值的问题：一般基础库的默认值都很保守，不适合实际情况。
- 高延迟的服务使用**超时传递** ：把超时策略传递进来。
- 超时控制是微服务可用性的第一道关，良好的超时策略，可以尽可能让服务**不堆积请求**，**尽快清空高延迟的请求**，**释放 Goroutine**。

- Service level objectives (SLOs)  a target level for the reliability of your service.
- 可以把Laency SLO 描述在 gRPC Proto 定义中

https://github.com/googleapis/googleapis/blob/master/google/monitoring/v3/service.proto#L170-L184
```
// A Service-Level Objective (SLO) describes a level of desired good service. It
// consists of a service-level indicator (SLI), a performance goal, and a period
// over which the objective is to be evaluated against that goal. The SLO can
// use SLIs defined in a number of different manners. Typical SLOs might include
// "99% of requests in each rolling week have latency below 200 milliseconds" or
// "99.5% of requests in each calendar month return successfully."
message ServiceLevelObjective {
  option (google.api.resource) = {
    type: "monitoring.googleapis.com/ServiceLevelObjective"
    pattern: "projects/{project}/services/{service}/serviceLevelObjectives/{service_level_objective}"
    pattern: "organizations/{organization}/services/{service}/serviceLevelObjectives/{service_level_objective}"
    pattern: "folders/{folder}/services/{service}/serviceLevelObjectives/{service_level_objective}"
    pattern: "*"
    history: ORIGINALLY_SINGLE_PATTERN
  };

```

- 基础库兜底：基础库配置100ms，进行防御保护，避免超大。
- 默认值（公共配置）兜底：对于未配置的服务使用公共配置。

### 超时传递 
  - 上游服务已经超时，但下游服务仍然在执行，会导致浪费资源做无用功。
  - 把当前服务的剩余超时量传递到下游服务中，继承超时策略，控制请求级别的全局超时控制。
  - 实现
    + go的`context.WithTimeout` 
    + 所有的服务调用首参数都是context, 构建带timeout的context即可。
    + 跨进程的传递：依赖gRPC Metadata，HTTP2-Header 传递 grpc-timeout 字段，自动传递到下游，
      
### 双峰分布 

- 95%的请求耗时在100ms内，5%的请求可能永远不会完成(长超时)。
- 监控不能只看平均，关注长耗时的分布统计，比如 95th，99th。
- 关注 5% 的请求 -> dead cases.
- 超时分布一般不是正态分布（unimodal)，而是双峰分布(bimodal)。（最快和最慢的双峰）  
- 设置合理的超时，拒绝超长请求，或者当Server 不可用要主动失败。
- 超时决定着服务线程耗尽。

### 基于超时的案例

- SLB (Server Load Balancer) 入口 Nginx 没配置超时导致连锁故障。
- 服务依赖的 DB 连接池漏配超时，导致请求阻塞，最终服务集体 OOM。
- 下游服务发版耗时增加，而上游服务配置超时过短，导致上游请求失败。

## 3 过载保护和限流

超时保护和过载保护的目的是都是让节点（服务）能够最存活：
 - 超时是让流量能尽快的消耗。
   > 服务消费者的放弃，不再出新请求，已经发的请求被迅速消耗掉。
 - 过载是当流量过多时候，服务的主动拒绝。
   > 服务提供者的自我保护，不接收请求。

### 令牌桶算法 (Token bucket)

#### 原理
  + 按照设定速率向一个固定容量的令牌桶中添加令牌。
  + 桶中最多存放固定容量的令牌，桶满时候，新添加令牌被丢弃（或拒绝）。
  + 服务请求需要获取令牌（即从桶中删除令牌）
  + 如果桶中当前令牌数量不满足最小设定，则服务无法获取令牌（无法再从桶中删除令牌），此时服务无法执行（过载保护）
  + 通过令牌桶的设定容量和添加令牌的速度可以让令牌桶设定一定的峰值。
    - ex: size:20 rate: 10/s 那么该令牌桶在某个峰值上，可以响应20/s个服务，
      然后会退化为10/s个服务。
    - 添加的速率代表令牌桶一般的服务吞吐能力。而令牌桶的最大容量代表瞬时的最大峰值。
  - https://en.wikipedia.org/wiki/Token_bucket
#### 实现    
- go实现 [`golang.org/x/time`:`rate.go`](https://github.com/golang/time/blob/f8bda1e9f3badef837c98cbaf4f7c335de90f266/rate/rate.go#L32-L64)  
- nginx实现 - [`ngx_http_limit_req_module.c`](https://github.com/nginx/nginx/blob/130a3ec5010227ca93498a1eb3a182062daeb349/src/http/modules/ngx_http_limit_req_module.c#L40-L47)

#### 问题
- 这种基于阈值的算法，主要问题是如何能知道正确的配置值，什么是合适的配置值？
- 这种值可能跟具体物理硬件有关，也可能和业务代码有关。
- 物理硬件是变化的 业务代码也是经常更新的，阈值很难设定！
- 实现简单，使用简单，核心问题是你不知道阈值如何设定。

### 漏桶算法 (Leaky Bucket)
#### 原理
- 一个固定容量的令牌桶，按照设定速率（常量固定）流出令牌。
- 流入令牌的速率为任意。
- 如果流入速度过快，超过桶的容量，直接丢弃
#### 实现
- https://github.com/uber-go/ratelimit

### 这两种算法的本质缺陷
- 防护思路都是设定一个指标, 当超过该指标后就阻止或减少流量的继续进入，当系统负载降低到某一水平后则恢复流量的进入。
- 都是被动的，其实际效果取决于限流阈值设置是否合理，但往往设置合理不是一件容易的事情。
  + 集群增加机器或者减少机器限流阈值是否要重新设置?
  + 设置限流阈值的依据是什么?
  + 人力运维成本是否过高?
  + 当调用方反馈429(429 Too Many Requests)时, 这个时候重新设置限流, 其实流量高峰已经过了重新评估限流是否有意义?
- 被动, 不能快速适应流量变化。
- 需要自适应的限流算法，根据系统当前的负载自动丢弃流量。

### 利特尔法制

```
L(门店最大顾客容量）= 速度（单位时间进门的顾客数） * 时间（顾客从进门到出门耗时）
```
这说明可以用流入速度QPS和响应时间latency

- 服务器临近过载时，主动抛弃一定量的负载，目标是自保。
- 在系统稳定的前提下，保持系统的吞吐量。
- CPU、内存作为信号量进行节流。
- 队列管理: 队列长度、LIFO。

### CoDel
- CoDel(Controlled Delay)队列管理算法
  - CoDel算法是诸多AQM策略算法中的一个 一个合理的队列管理策略
  - CoDel pseudocode https://queue.acm.org/appendices/codel.html
- BBR (Google的TCP BBR拥塞控制算法)
  - https://en.wikipedia.org/wiki/TCP_congestion_control#TCP_BBR
  - https://blog.csdn.net/dog250/article/details/52830576
  - https://blog.csdn.net/dog250/article/details/72042516
  - https://blog.csdn.net/dog250/article/details/72849890
  - https://blog.csdn.net/dog250/article/details/72849893
- Kratos的控流算法
  - https://github.com/go-kratos/kratos/blob/v1.0.x/pkg/ratelimit/bbr/bbr.go
- [阿里Sentinel的系统自适应限流](https://github.com/alibaba/Sentinel/wiki/%E7%B3%BB%E7%BB%9F%E8%87%AA%E9%80%82%E5%BA%94%E9%99%90%E6%B5%81)

## 限流
过载保护是自适应的（如果使用前述的算法），而限流是限定的，通过指标来控制。
- 令牌、漏桶针对单个节点，无法分布式限流。
  - 一个节点一个节点的配置难以达到合适效果，无法精确达到要求。
### 分布式限流
- 单个大流量的接口，使用 redis 容易产生热点。
- pre-request 模式对性能有一定影响，高频的网络往返。
- 从获取单个 quota 升级成批量 quota。获取后使用令牌桶算法来限制。
  - 每次心跳后，异步批量获取 quota，可以大大减少请求 redis 的频次，获取完以后本地消费，基于令牌桶拦截。
  - 初次使用默认值，一旦有过去历史窗口的数据，可以基于历史窗口数据进行 quota 请求。
  - 使用“最大最小公平分享”（Max-Min Fairness）。首先均分，然后按需求量从最小开始分配，先分配最小需求者，未得到满足的用户再均分，直到分配完备的原则。
  - quota针对服务级别（粗粒度），更细粒度的接口级别，设定重要性
    + 注意调用链的情况下，重要性需要传递（A->B,A->C），那么B，C重要性等于A
  - 全局quota不足时候，优先拒绝重要性低的。
  
### 熔断（Circuit Breakers）
- 前述的限流保护的都是服务本身。
- 熔断是指client side限流，保护的是下游的服务。
- 实现
  - https://martinfowler.com/bliki/CircuitBreaker.html
    + 熔断器`关闭`状态
    + 默认状态，持续统计一些指标（rpc错误率）。
    + 达到指标指，进入熔断器打开状态。
    + 熔断器`打开`状态
      + client不再对外处理服务请求。
      + client自己内部等一个timeout，然后进入熔断器`半开`状态。
    + 熔断器`半开`状态
      + 该状态下，client向后台发起一个try 请求。
      + 如果try失败，回到熔断器`打开`状态。
      + 如果try成功，进入熔断器`关闭`状态，这时对外处理服务请求。
  - 参考java的[Hystrix](https://github.com/Netflix/Hystrix)
- 改进
  - google sre  (Client request rejection probability)
    - https://sre.google/sre-book/handling-overload/#eq2101
  - https://github.com/tal-tech/go-zero/blob/master/core/breaker/googlebreaker.go#L41
  - https://github.com/go-kratos/kratos/blob/v1.0.x/pkg/net/netutil/breaker/sre_breaker.go#L74
    
### Kafka Gutter
- 主kafka 配一个副kafka（非等比例，只有主的10%），用于接管限流溢出的负载
- 副kafka只使用10%的资源，平常不接受流量时，避免资源浪费。
- 核心利用熔断的思路，是把抛弃的流量转移到 gutter 集群，如果 gutter 也接受不住的流量，重新回抛到主集群，最大力度来接受。

### 客户端限流

- 客户端休眠的策略（参考google的grpc的backoff）

```golang
// Backoff returns the amount of time to wait before the next retry given the
// number of retries.
func (bc Exponential) Backoff(retries int) time.Duration {
	if retries == 0 {
		return bc.Config.BaseDelay
	}
	backoff, max := float64(bc.Config.BaseDelay), float64(bc.Config.MaxDelay)
	for backoff < max && retries > 0 {
		backoff *= bc.Config.Multiplier
		retries--
	}
	if backoff > max {
		backoff = max
	}
	// Randomize backoff delays so that if a cluster of requests start at
	// the same time, they won't operate in lockstep.
	backoff *= 1 + bc.Config.Jitter*(grpcrand.Float64()*2-1)
	if backoff < 0 {
		return 0
	}
	return time.Duration(backoff)
}
```
- [GRPC Connection Backoff Protocol](https://github.com/grpc/grpc/blob/master/doc/connection-backoff.md)
  + https://github.com/grpc/grpc-go/blob/master/backoff/backoff.go
    + https://github.com/grpc/grpc-go/blob/master/internal/backoff/backoff.go
  + https://github.com/grpc/grpc/blob/master/src/core/lib/backoff/backoff.cc   
- https://github.com/go-kratos/kratos/blob/v1.0.x/pkg/net/netutil/backoff.go

## 4 降级

- 降级本质为: 提供有损服务。
- UI 模块化，非核心模块降级。
- BFF 层聚合 API，模块降级。
- 页面上一次缓存副本。
- 默认值、热门推荐等。
- 流量拦截 + 定期数据缓存(过期副本策略)。
- 处理策略
  - 页面降级、延迟服务、写/读降级、缓存降级
  - 抛异常、返回约定协议、Mock 数据、Fallback 处理

## 5 重试 & 负载均衡

### 重试
- 留意重试带来的流量放大:
  + 限制重试次数和基于重试分布的策略（例如重试比率: 不能超正常流量的10%）。
  + 随机化、指数型递增的重试周期: exponential ackoff + jitter。
    - backoff限流算法
  + client 测记录重试次数直方图，传递到 server，进行分布判定，交由 server 判定拒绝。
  + 只应该在失败的这层进行重试，当重试仍然失败，全局约定错误码“过载，无须重试”，避免级联重试。
    - 重要：级联重试会造成重试流量的迅速放大，必须避免。
- 业务不幂等，导致重试结果不一。  
- “幂等”原为数学概念，指函数被执行多次，结果始终不变。
  + 数字乘1：幂等。执行多少次运算，结果都是相同的。
  + 数字取整：幂等：n次 floor(3.4)，结果都是3。
  + 数字加1：不幂等，因为执行一次和执行多次的结果不同。

### 负载均衡

- P2C算法 The power of two choices in randomized load balancing
  - 随机选取的两个节点进行打分，选择更优的节点:
     + 选择 backend：CPU，client：health、inflight、latency 作为指标，使用一个简单的线性方程进行打分。
     + 对新启动的节点使用常量惩罚值（penalty），以及使用探针方式最小化放量，进行预热。
     + 打分比较低的节点，避免进入“永久黑名单”而无法恢复，使用统计衰减的方式，让节点指标逐渐恢复到初始状态(即默认值)。
     + 当前发出去的请求超过了 predict lagtency，就会加惩罚。
     + 指标计算结合 moving average，使用时间衰减，计算vt = v(t-1) * β + at * (1-β) ，β 为若干次幂的倒数即: Math.Exp((-span) / 600ms)
  - 参考：
    + https://www.eecs.harvard.edu/~michaelm/postscripts/tpds2001.pdf
    + https://www.nginx.com/blog/nginx-power-of-two-choices-load-balancing-algorithm/ 
  - 实现
    + nginx C 
      - https://github.com/nginx/nginx/blob/master/src/http/modules/ngx_http_upstream_least_conn_module.c#L100
    + envoy C++ 
      - https://github.com/envoyproxy/envoy/blob/main/source/common/upstream/load_balancer_impl.h#L471-L509
    + tower Rust 
      - https://github.com/tower-rs/tower/blob/master/tower/src/balance/p2c/service.rs#L181-L210
  - 对于go来说，一般选择扩展 [gRPC balancer](https://github.com/grpc/grpc-go/tree/master/balancer/base)
    + douyu 
      - https://github.com/douyu/jupiter/blob/master/pkg/util/xp2c/leastloaded/least_loaded.go#L49-L80
    + kratos 
      - https://github.com/go-kratos/kratos/blob/v1.0.x/pkg/net/rpc/warden/balancer/p2c/p2c.go#L190-L270
  
## 最佳实践

- 变更管理:
  + 70％的问题是由变更引起的，恢复可用代码并不总是坏事。
- 避免过载:
  + 过载保护、流量调度等。
- 依赖管理:
  + 任何依赖都可能故障，做 chaos monkey testing，注入故障测试。
- 优雅降级:
  + 有损服务，避免核心链路依赖故障。
- 重试退避:
   + 退让算法，冻结时间，API retry detail 控制策略。
- 超时控制:
   + 进程内 + 服务间 超时控制。
- 极限压测 + 故障演练。
- 扩容 + 重启 + 消除有害流量。

## References

- hhtp://www.360doc.com/content/16/1124/21/31263000_609259745.shtml
- http://www.infoq.com/cn/articles/basis-frameworkto-implement-micro-service/
- http://www.infoq.com/cn/news/2017/04/linkerd-celebrates-one-year
- https://medium.com/netflix-techblog/netflix-edge-load-balancing-695308b5548c
- https://mp.weixin.qq.com/s?__biz=MzAwNjQwNzU2NQ==&mid=402841629&idx=1&sn=f598fec9b370b8a6f2062233b31122e0&mpshare=1&scene=23&srcid=0404qP0fH8zRiIiFzQBiuzuU#rd
- https://mp.weixin.qq.com/s?__biz=MzIzMzk2NDQyMw==&mid=2247486641&idx=1&sn=1660fb41b0c5b8d8d6eacdfc1b26b6a6&source=41#wechat_redirect
- https://blog.acolyer.org/2018/11/16/overload-control-for-scaling-wechat-microservices/
- https://www.cs.columbia.edu/~ruigu/papers/socc18-final100.pdf
- https://github.com/alibaba/Sentinel/wiki/系统负载保护
- https://blog.csdn.net/okiwilldoit/article/details/81738782
- http://alex-ii.github.io/notes/2019/02/13/predictive_load_balancing.html
- https://blog.csdn.net/m0_38106113/article/details/81542863