# 第8周 分布式缓存 & 分布式事务

## 目录

分布式缓存
 - 缓存选型
 - 缓存模式 
 - 缓存技巧

分布式事务
 - 

## 缓存选型

### Memcache Memcache 

- 提供**简单的kv** cache 存储
  
- value 大小**不超过1MB**
  
- B站初始使用 Memcache 作为大文本（评论内容，弹幕文件集合）或者简单的kv结构
  - 初期选型：取大文本时候吞吐比redis吞吐好，（redis相对的操作更大，再早期版本表现明显）
  - 参考了微博的关系链（我关注你/你关注我）的设计。使用memcache在redis前做了一层遮挡。
    
- Memcache 使用了**slab方式**做内存管理
  - slab
    - Sun的Jeff Bonwick在Solaris 2.4中设计并实现
    - 后来被Linux所借鉴，用于实现内核中更小粒度(比page更小)的内存分配
    -  每个 slab 包含若干大小为1M的内存⻚，
       - 这些内存又被分割成多个 chunk，每个 chunk 存储一个 item;
       - chunk 的增⻓因子由 -f 指定，默认1.25，起始大小为48字节。
    - Memcache 初始化时，每个 slab 都预分配一个 1M 的内存⻚，
       - 由 slabs_preallocate 完成
       - 也可将相应代码注释掉关闭预分配功能
  - 存在一定的浪费
  - 如果大量接近的item，建议调整 Memcache 参数来优化每一个 slab 增⻓ 的 ratio
    - 可以通过设置 slab_automove & slab_reassign 开启 Memcache 的动态/手动 move slab
    - 防止某些 slab 热点
       - 导致内存足够的情况下引发 LRU。 
    
- 大部分情况下，简单KV推荐使用 Memcache，吞吐和响应都足够好。

#### 内存池设计参考
- nginx
  - ngx_pool_t
    - create
      - https://github.com/nginx/nginx/blob/5eadaf69e394c030056e4190d86dae0262f8617c/src/core/ngx_palloc.c#L19 
      - https://github.com/nginx/nginx/blob/0ab91b901299ac41e3867ebec7e04e5082a4c8b4/src/core/ngx_palloc.c#L6  (Igor Sysoev Jun 2004)
    - align
      - https://github.com/nginx/nginx/blob/5eadaf69e394c030056e4190d86dae0262f8617c/src/core/ngx_palloc.c#L253
    - 小块内存分配
      - https://github.com/nginx/nginx/blob/master/src/core/ngx_palloc.c#L149
    - 大块内存分配
      - https://github.com/nginx/nginx/blob/master/src/core/ngx_palloc.c#L214 
    - 一开始就想过 
     - https://github.com/nginx/nginx/blob/release-0.1.0/src/core/ngx_slab.c (2004) 
- tcmalloc (go的内存管理基于tcmalloc)
   - Thread-Cache Malloc
   - https://github.com/google/tcmalloc 
   - https://github.com/golang/go/blob/master/src/runtime/malloc.go#L7-L8
   - [译文：Go 内存分配器可视化指南](https://www.linuxzen.com/go-memory-allocator-visual-guide.html)
     - https://medium.com/@ankur_anand/a-visual-guide-to-golang-memory-allocator-from-ground-up-e132258453ed

### 缓存选型 - Redis
- Redis 有**丰富的数据类型**，支持增量方式的修改部分数据，比如排行榜，集合，数组等。
- 比较常用的方式是使用 Redis 作为数据索引
  - 比如评论的列表 ID，播放历史的列表ID集合
  - B站的关系链列表 ID
- Redis 自己没有使用内存池，所以是存在一定的内存碎片
  - 一般会使用 jemalloc 来优化内存分配 (有些命令是implemented only when using jemalloc
    - By default, Redis uses jemalloc memory allocator on Linux.
  - 需要编译时候使用 jemalloc 库 代替 glib 的 malloc 使用。
    - jemalloc是freebsd的libc allocator（https://github.com/jemalloc/jemalloc）
  - Redis支持(glibc malloc/jemalloc11/tcmalloc) 
    - https://github.com/redis/redis/blob/2.4/src/zmalloc.h (2.4开始支持jemalloc)
    - https://github.com/redis/redis/blob/3e39ea0b83f5588a5460a366072a5c7b3bd42635/src/zmalloc.h#L38-L73
  - 有意思的比较
    - http://ithare.com/testing-memory-allocators-ptmalloc2-tcmalloc-hoard-jemalloc-while-trying-to-simulate-real-world-loads/  
    
### 缓存选型 - Redis vs Memcache
- Redis 和 Memcache 最大的区别其实是 
  - redis 单线程
     - 新版本双线程
        - 一个线程单独负责io
  - memcache 多线程
- 所以 
  - QPS 可能两者差异不大，但是吞吐会有很大的差别
   - 比如大数据 value 返回的时候，redis qps 会抖动下降的的很厉害，
   - 因为单线程工作，其他查询进不来(新 版本有不少的改善)。
- 所以
  - 建议纯 kv 都走 memcache， 
- B站的关系链服务中用了 hashes 存储双向关系
  - hashes 
    - https://redis.io/topics/data-types#hashes
    - redis数据类型，maps between string fields and string values
      - https://redis.io/commands#hash
      - Hashes, which are maps composed of fields associated with values. 
        - Both the field and the value are strings. 
        - very similar to Ruby or Python hashes.
        - 注意key是string，value是fields，也就是用一个key，和一个map
        - key -> { field1: value1, field2 : value2} 这样的数据结构。
  - 但是我们也会使用 memcache 档一层来避免 hgetall 导致的吞吐下降问题。
- B站系统中多次使用 memcache + redis 双缓存设计。

### 缓存选型 - Proxy
- 早期使用 twemproxy 作为缓存代理
  - https://github.com/twitter/twemproxy 
    - A fast, light-weight proxy for memcached and redis Resources
    - 0.4.1 2015年
  - 目的是把redis shading出去，twemproxy内置一致性hash算法，封装了redis cluster，对外redis api。
  - 使用上有如下一些痛点:
    - 单进程单线程模型和 Redis 类似，在处理一 些大 key 的时候可能出现 io 瓶颈;
    - 二次开发成本难度高，难以于公司运维平台 进行深度集成;
    - 不支持自动伸缩，不支持 autorebalance 增 删节点需要重启才能生效;
    - 运维不友好，没有控制面板;
- 也考察过业界开源的的其他代理工具:
  - codis: 只支持 Redis 协议，且需要使用 patch 版本的 redis;
     - https://github.com/CodisLabs/codis (最新版3.2.2 2018年)
  - mcrouter: 只支持 memcache 协议，C 开发， 与运维集成开发难度高;
     - https://github.com/facebook/mcrouter (v41 2019)

- 自研（b站）了 overload（go语言实现）做为redis cluster代理。  

- 从集中式访问缓存到Sidecar模式访问缓存:
  - 微服务强调去中心化;
  - LVS 运维困难，容易流量热点，随下游扩容而扩容，连接不均衡等问题; 
  - Sidecar伴生容器随App容器启动而启动，配置简化
    -k8s伴生容器，把`overload`伴生进去。
    -不算service mesh，只是缓存代理使用sidebar模式，轻量。
    
###  一致性 Hash

#### 概念
- 一致性hash 
  - 将数据按照特征值映射到一个首尾相接的hash环上
  - 同时也将节点 (按照IP地址或者机器名hash)映射到这个环上。
  - 对于数据，从数据在环上的位置开始，顺时针找到的第一个节点即为数据的存储节点。
  - 余数分布式算法由于保存key的服务器会发生巨大变化而影响缓存的命中率
  - 但一致性hash中，只有在环上增加服务器的地方的逆时针方向的第一台服务器上的键会受到影响。
    - 这样加入新节点时候，对缓存的影响不大。 
  
#### 目标
- 平衡性(Balance)
  - 尽可能平衡分布到所有的缓存
- 单调性(Monotonicity):
  - 如果已经有一些内容通过哈希分派到了相应的缓存中，又有新的缓存加入到系统中，
     - 那么哈希的结果应能够保证原有已分配的内容可以被映射到新的缓存中去，
     - 而不会被映射到旧的缓冲集合中的其他缓冲区。
- 分散性(Spread):
  - 相同内容被存储到不同缓冲中去，降低了系统存储的效率，需要尽量降低分散性。
- 负载(Load)
  - 哈希算法应能够尽量降低缓冲的负荷。
- 平滑性(Smoothness)
  - 缓存服务器的数目平滑改变和缓存对象的平滑改变是一致的。
  
#### 虚拟节点机制
- 一致性哈希算法在服务节点太少时，容易因为节点分部不均匀而造成数据倾斜问题。
- 例如2台节点，hash不平衡，如果偏A，必然造成大量数据集中到 NodeA 上，只有极少量会定位到NodeB。
- 为了解决数据倾斜问题，引入虚拟节点机制
   - 即对每一个服务节点计算多个哈希，每个计算结果位置都放置一个此服务节点，称为虚拟节点。
- 具体做法
  - 可以在服务器 IP 或主机名的后面增加编号来实现。
  - 例如
    - 可以为NodeA和NodeB，每台服务器计算三个虚拟节点，
      - 于是可以分别计算“Node A#1”、“Node A#2”、“Node A#3”、 “Node B#1”、“Node B#2”、“Node B#3”的哈希值，
      - 于是形成六个虚拟节点。 
      - 同时数据定位算法不变，只是多了一步虚拟节点到实际节点的映射
      - 例如定位到 “Node A#1”、“Node A#2”、“Node A#3”三个虚拟节点的数据均定位到 Node A 上。
      - 这样就解决了服务节点少时数据倾斜的问题。

#### 一致性Hash和微信红包的写合并优化:
- https://www.cnblogs.com/chinanetwind/articles/9460820.html
- 在网关层，使用一致性 hash，对红包 ID 进行分片，命中到某一个逻辑服务器处理，
- 在进程内做写操作的合并，减少存储层的单行锁争用。
- 更好的做法
  - 有界负载一致性hash。

#### 一致性Hash和数据分片

- 按照数据的某一特征(key)来计算哈希值，并将哈希值与系统中的节点建立映射关系,从而将哈希值不同的数据分布到不同的节点上。
- 按照 hash 方式做数据分片，映射关系非常简单;
- 需要管理的元数据也非常之少，只需要记录节点的数目以及 hash 方式。
  
问题:
 - 当加入或者删除一个节点的时候，大量的数据需要移动。
 - 原始数据的特征值分布不均匀，导致大量的数据集中到一个物理节点上
 - 对于可修改的记录数据，单条记录的数据变大。

解决：
- 高级玩法是抽象 slot，基于 Hash 的 Slot Sharding，例如 Redis-Cluster。

#### redis-cluster的Slot sharding

- redis-cluster 一开始分配了16384个槽（0～16383）
- 把16384 槽按照节点数量进行平均分配，由节点进行管理。（例如5个节点，那么node1管理（0～3276））
- 对每个 key 按照 CRC16 规则进行 hash 运算，
- 把hash结果对16383进行取余，根据余数找到槽位，再找到Redis节点。
- 还是hash求余，但是hash求余先命中槽位，然后看这个槽谁负责，再到节点
  - 好处是什么？
    - 当加入node6时候，原来5个节点，每人贡献一些槽位给node6，那么数据迁移就可以均匀。

- 需要注意
  - Redis Cluster 的节点之间会 共享消息
  - 每个节点都会知道是哪个节点负 责哪个范围内的数据槽

- slot sharding 是很常见的数据分片做法。

## 缓存模式

### 缓存模式 - 数据一致性 

- **Storage 和 Cache 同步更新容易出现数据不一致**。

- 模拟 MySQL Slave 做数据复制，再把消息投递到 Kafka，保证至少一次消费: 
  - 1.同步操作 DB;
  - 2.同步操作 Cache;
  - 3.利用 Job 消费消息，重新补偿一次缓存操作
  - 保证时效性和一致性。利用Job回放解决**最终一致性**

- Cache Aside模型
  - 读缓存 Miss 的回填操作，和修改数据同步更新缓存，
  - 包括消息队列的异步补偿缓存，都无法满足 “Happens Before”，会存在相互覆盖的情况。
  > Cache Aside的脏数据问题：
  >  - 首先一个读操作，没有命中缓存，到数据库中取数据，此时获得v1，
  >  - 同时有一个写操作，写数据库形成v2，让缓存失效，
  >  - 之前的那个读操作再把老的数据(v1)放进去缓存，这样造成脏数据。
  >  - 实际上出现的概率可能非常低
  >    - 因为这个条件需要发生在读缓存时缓存失效，而且并发着有一个写操作。
  >    - 而实际上数据库的写操作会比读操作慢得多，而且还要锁表，
  >    - 而读操作必需在写操作前进入数据库操作，而又要晚于写操作更新缓存，
  >    - 所有的这些条件都具备的概率基本并不大。
  >  - 为缓存设置上过期时间 -> 降低并发时脏数据的概率 
  >  - 不需要强一致性，最终一致性就可以了。

### 数据一致性的一个解决例子

- 读操作和写操作同时进行:
  - 1 读操作，读缓存，缓存 MISS
  - 2 读操作，读 DB，读取到数据
  - 3 写操作，更新 DB 数据
  - 4 写操作，SET(or DELETE) Cache
  - 5 读操作，SET Cache 操作数据回写缓存
  - 这种交互下，由于4和5操作步骤都是设置缓存，导致写入的值互相覆盖
    - 并且操作的顺序性不确定，从而导致 cache 存在脏缓存的情况。
- 解决：读操作和写操作同时进行：
  - 1 读操作，读缓存，缓存 MISS
  - 2 读操作，读 DB，读取到数据
  - 3 写操作，更新 DB 数据
  - 4 写操作，SET Cache （写操作保证SET）
  - 5 **读操作，ADD Cache** 操作数据回写缓存 （读操作不一定回填缓存）
  - 第5步，读操作，改使用 ADD操作 回写 MISS 数据， 从而保证写操作的最新数据不会被读操作的回写数据覆盖。
    - ADD操作 
      - 可用 Job 异步操作
      - Redis 可以使用 SETNX (SET if Not eXists)，即写操作优先级更高，我自己主动放弃。

### 缓存模式 - 多级缓存

- 整合服务(聚合服务) 用于提供粗粒度的接口，以及二级缓存加速，减少扇出的RPC网络请求，减少延迟。
- 问题 ：**最重要是保证多级缓存的一致性**:
  - **清理的优先级**是有要求的，先优先清理下游再上游;
  - **下游的缓存 expire 要大于上游**，里面穿透回源;

### 缓存模式 - 热点缓存 

- 对于热点缓存 Key解决思路:
  - 小表广播：把Remote Cache（redis） 提升为 Local Cache（进程）
      - App 定时更新 
      - 甚至可以让运营平台支持广播刷新Local Cache 
  - 主动监控防御预热
      - 比如直播房间⻚高在线情况下直接外挂服务主动防御 
  - 基础库框架支持热点发现
    - 自动短时的 short-live cache
    - krotos v2 middle-ware
  - 多 Cluster 支持
    - 多 Key 设计: 使用多副本，减小节点热点
    - 使用多副本 ms_1,ms_2,ms_3 每个节点保存一份数据，使得请求分散到多个节点，避免单点热点问题。
  - 建立多个 Cluster ，和微服务、存储等一起组成一个Region。 
      - 这样相当于是用空间换时间:
      - 同一个 key 在每一个 frontend cluster 都可能有一个 copy，这样会带来一致性问题，
        - 但是这样能够降低 latency 和提高 availability。
      - 利用 MySQL Binlog 消息 anycast 到不同集群的某个节点清理或者更新缓存
  - 当业务频繁更新时候，cache频繁过期，会导致命中率低
    - 例如一个热点的key被删除，那么大量的请求透传到DB，就可能打爆DB
    - 使用 stale sets (短时间把脏数据返回给用户，而不要打爆DB)
      - 当一个 key 被删除，被放倒一 个临时的数据结构里，会再续上比较短的一段时间。
      - 当有请求进来的时候会返回这个数据
      - 并标记为“Stale”。
      - 对于大部分应用场景而言，Stale Value 是可以忍受的。
      - 需要改 memcache、Redis 源码，或者基础库支持

### 缓存模式 - 穿透缓存
- singleflight
  - 对关键字进行一致性 hash，使其某一个维度的 key 一定命中某个节点，
  - 然后在节点内使用互斥锁，保证归并回源 (见week06)
  - 但是对于批量查询无解
  
- 分布式锁（最好不要用，强一致性是很难保证的，往往得不偿失，而且bug密集型思路）
  - 设置一个 lock key，有且只有一个人成功，并且返回，交由这个人来执行回源操作，
  - 其他候选者轮询 cache 这个 lock key，
     - 如果不存在去读数据缓存，
     - hit 就返回
     - miss 继续抢锁

- 队列（建议的做法，和singlefight配合一起）
  - 如果 cache miss，交由队列聚合一个 key，来 load 数据回写缓存，
  - 对于 miss 当前请求可以使用 singlefly 保证回源，
     - 如评论架构实现。适合回源加载数据重的任务，
     - 比如评论 miss 只返回 第一⻚，
        - 但是需要构建完成评论数据索引。
  
- lease (facebook的租约机制)
  - lease 是 64-bit 的 token，与客户端请求的key是绑定的，
  - 有token的才能访问DB， 在写入缓存时需要验证token，（有token才能更新缓存）
  - 每个key 10s 重写分配一次token（token过期时间）
  - 当 client 在没有获取到 token 时，就等待cache被构建好。
  - 基础库支持 & 修改 cache 源码

- CRDT(Conflict-Free Replicated Data Type)
  - https://hal.inria.fr/inria-00609399/document
  - 数据结构  
  - 各种基础数据结构最终一致算法的理论总结，能根据一定的规则自动合并，解决冲突，达到强最终一致的效果。
  - Eric Brewer的回首CAP20年文章提到：C和A并不是完全互斥，建议大家使用CRDT来保障一致性。
    - https://www.infoq.com/articles/cap-twelve-years-later-how-the-rules-have-changed/ (Eric Brewer, May30 2012)
      - commutative replicated data types (CRDTs) 
        - a class of data structures that provably converge after a partition
        - ensure that all operations during a partition are commutative, or 
        - represent values on a lattice and ensure that all operations during a partition are monotonically increasing with respect to that lattice.


## 缓存技巧

### 缓存技巧 - Incast Congestion
- 如果在网路中的包太多，就会发生 Incast Congestion 的问题
  - network很多 switch/router，一次性发一堆包，这些包同时到达switch，switch忙不过来
- 解决：不要让大量包在同一时间发送出去
  - 客户端限制每次发出去的包的数量
  - 具体实现：客户端队列。
- 每次发送的包的数量称为“Window size”。
  - 值太小，发送太慢，延迟会变高
  - 值太大，发送包太多 -> switch 崩溃，可能发生丢包，可能被当作cache miss，延迟也会变高。
  - 这个值需要调，一般在proxy层面实现

### 缓存技巧 - 通用技巧

- 易读性前提下，key设置尽可能小，减少资源的占用
  - redis value 可以用 int 就不要用 string
  - 对于小于 N 的 value，redis 内部有 shared_object 缓存。
  
- 拆分 key。
   - 主要是用在 redis 使用 hashes 情况下。
    - 同一个 hashes key 会落到同一个 redis 节点，
    - hashes 过大的情况下会导致内存及请求分布的不均匀。
    - 考虑对 hash 进行拆分为小的 hash，使得节点内存均匀及避免单节点请求热点。
  
- 空缓存设置
    - 对于部分数据，可能数据库始终为空，
    - 攻击场景：故意构造请求，来透传到DB
      - key不应该容易被枚举，对外暴露的API需要考虑这个场景。
    如果key还是能透传
    - 这时应该设置空缓存，避免每次请求都缓存 miss 直接打到 DB。
  
- 空缓存保护策略 
  
- 读失败后的写缓存策略(降级后一般读失败不触发回写缓存)。
  
- 序列化使用 protobuf，尽可能减少 size。
  
- 工具化胶水代码
  - Java有annotation
  - go可以基于protobuf写插件，code生成。

### 缓存技巧 - memcache 小技巧
- flag 使用:标识 compress、encoding、large value 等;
- memcache 支持 gets，尽量读取，尽可能的 pipeline，减少网络往返; 使用二进制协议，支持 pipeline delete，UDP 读取、TCP 更新;

### 缓存技巧 - Redis 小技巧
- 增量更新一致性
   - EXPIRE、ZADD/HSET 等，保证索引结构体务必存在的情况下去操作新 增数据;
- BITSET
  - 存储每日登陆用户，单个标记位置(boolean)，
  - 为了避免单个 BITSET 过大或者 热点，需要使用 region sharding，
    - 比如按照 mid求余 %和/ 10000，商为 KEY、余数作为 offset;
- List
   - 抽奖的奖池、顶弹幕，用于类似 Stack PUSH/POP操作;
- Sortedset
   - 翻⻚、排序、有序的集合，
  - 杜绝 zrange 或者 zrevrange 返回的集合过大
- Hashs 
  - 过小的时候会使用压缩列表、
  - 过大的情况容易导致 rehash 内存浪费，也杜绝返回 hgetall，
  - 对于小结构体，建议直接使用 memcache KV;
- String
  - SET 的 EX/NX 等 KV 扩展指令，
  - SETNX 可以用于分布式锁、
  - SETEX 聚合了SET + EXPIRE;
- Sets
  - 类似 Hashs，无 Value，去重等;
  
- 尽可能打包指令，批量操作，但是避免集合过大 

- 避免超大Value

## References
- [微博应对日访问量百亿级的缓存架构设计](https://mp.weixin.qq.com/s?__biz=MzkwOTIxNDQ3OA==&mid=2247533269&idx=1&sn=b0b0146d8afa51ece102f9a20edc7417&source=41)
- [Redis 集群中的纪元(epoch)](https://blog.csdn.net/chen_kkw/article/details/82724330)
- [一万字详解 Redis Cluster Gossip 协议](https://zhuanlan.zhihu.com/p/328728595)
- [微信红包系统架构的设计和优化分享](https://www.cnblogs.com/chinanetwind/articles/9460820.html)
- [Improving load balancing with a new consistent-hashing algorithm](https://medium.com/vimeo-engineering-blog/improving-load-balancing-with-a-new-consistent-hashing-algorithm-9f1bd75709ed)
- [浅谈分布式存储系统数据分布方法](https://www.jianshu.com/p/5fa447c60327)
- [一致性哈希算法（一）- 问题的提出](https://writings.sh/post/consistent-hashing-algorithms-part-1-the-problem-and-the-concept)
- [高可用Redis(十二)：Redis Cluster](https://www.cnblogs.com/renpingsheng/p/9862485.html)  

## 分布式事务

### 经典转账问题
- 支付宝账户表:A (id, user_id, amount)
- 余额宝账户表:B (id, user_id, amount)
- 用户的 user_id = 1，从支付宝转帐1万到余额宝分为两个步骤:
  - 1. 支付宝表扣除1万:
   UPDATE A SET amount = amount - 10000
   WHERE user_id = 1;
  - 2. 余额宝表增加1万:
   UPDATE B SET amount = amount + 10000
   WHERE user_id = 1;

- 如何保证数据一致性呢?
- 单个数据库，我们保证 ACID 使用 数据库事务。

### 转账问题的微服务场景

- 微服务架构改造，每个微服务独占了一个数据库实例
- 从 user_id = 1 发起的转帐动作，跨越了两个 微服务:pay 和 balance 服务。
- 我们需要保证，跨多个服务的步骤数据一致性:
  - 1. 微服务 pay 的支付宝表扣除1万
  - 2. 微服务 balance 的余额宝表增加1万
- 每个系统都对应一个独立的数据源，且可能位于不同机房，同时调用多个系统的服务很难保证同时成功
- 跨服务分布式事务
   - 保证每个服务自身的ACID，
  - 事务消息解决分布式事务问题。

### 分布式事务 - 事务消息
- 小吃店场景，点单付费，给顾客小票，顾客拿着小票到出货区排队去取。
- 为什么要将付钱和取货两个动作分开呢?
  - 使他们接待能力增强 -> 并发量更高
  - 只要这张小票在最终是能拿到货。
  - 同理转账服务也是如此。

- 当账户扣除1万后
  - 生成一个凭证 (消息)
    - 这个凭证(消息)上写着“让余额宝账户增加 1万”
    - 只要这个凭证(消息)能可靠保存，最终可以拿着这个凭证(消息)让余额宝账户增加1万的
    - 即我们能依靠这个凭证(消息)完成最终一致性。

### 分布式事务 - 如何可靠的保存消息凭证?

解决消息可靠存储

 - 解决本地的 sql 存储和 msg 存储 的一致性问题。
 - 步骤  
   - 1. Transactional outbox 
   - 2. Polling publisher
   - 3. Transaction log tailing
   - 4. 2PC Message Queue
   
- 事务消息一旦被可靠的持久化
  - 我们整个分布式事务，变为了最终一致性
  - **消息的消费才能保障最终业务数据的完整性**
  - 所以我们要尽最大努力，把消息送达到下游的业务消费方
    - 称为:**Best Effort**。只有消息被消费，整个交易才能算是完整完结。
    > https://en.wikipedia.org/wiki/Best-effort_delivery
    > 这个术语是个网络通信术语
  - 支付宝交易接口
    - 一般会在支付宝的回调⻚面和接口里，解密参数，然后调用系统中更新交易状态相关的服务，将订单更新为付款成功。
    - 同时，只有当我们回调⻚面中输出了 success 字样或者标识业务处理成功相应状态码时，支付宝才会停止回调请求。
    - 否则，支付宝会每间隔一段时间后，再向客户方发起回调请求，直到输出成功标识为止。
      > - 这种不停callback直到成功就是Best Effort模式
      > - 注意：
      >   - 这种重试模型结果的正确性是由回调服务提供方保证的，所以必须注意幂等性，否则会造成攻击风险（例如：一次付费，多个商品（不幂等））
      >   - **任何存在多次重试（消息重复投递/消息重复消费）的情况，都要考虑服务的幂等性**

### 分布式事务 - 事务消息 1. Transactional outbox
- 支付宝在完成扣款的同时，同时记录消息数据
  - 假设消息数据与业务据保存在同一数据库实例 
  - 假设**消息**表名为msg
  ```sql
  BEGIN TRANSACTION
  UPDATE A SET amount = amount - 10000 WHERE user_id = 1;
  INSERT INTO msg(user_id, amount, status) VALUES(1, 10000, 1);
  END TRANSACTION COMMIT;
  ```
- 上述事务能保证只要支付宝账户里被扣了钱，**消息**一定能保存下来
- 当上述事务提交成功后，将此消息通知余额宝
  - 余额宝处理成功后发送回复成功消息
  - 支付宝收到回复后删除该条消息数据

### 分布式事务 - 事务消息 2. Polling publisher 
- 一个独立的 pay_task 服务，
  - 定时轮询 msg 表，把 status = 1 的消息统统拿出来消费，
   - 可以按照自增 id 排序，保证顺序消费。
  - 把拖出来的消息 publish 给我们消息队列，
- balance 服务自己来消费队列
  - 或者直接 rpc 发送给 balance 服 务。
- pull的问题  
  - B站第一个版本的 archive-service 在实现 CQRS 时
    > - CQRS(stands for Command Query Responsibility Segregation) 参考week01
    > - 注：CQRS模式本身对更新的描述是订阅模式，而不是pull模式，其实说明实现的不够好，不能称为CQRS。
    - 就使用这个 Pull 的模型
    - 延迟不够好，
      - Pull 太猛对 Database 有一定压力
      - Pull 频次低了，延迟比较高。

### 分布式事务 - 事务消息 2. Transaction log tailing

- 上述保存消息的方式 使得消息数据和业务数据紧耦合在一起，从架构上看不够优雅，而且容易诱发其他问题。
  - 改为订阅模式
  
- 有一些业务场景，
   - 可以直接使用主表被 canal 订阅使用
- 有一些业务场景自带这类 message 表
  - 比如订单或者交易流水，可以直接使用这类流水表作为 message 表使用。
- 使用 canal 订阅以后，是实时流式消费数据，
  - 在消费者 balance 或者 balance-job 必须努力送达到。

- 所有努力送达的模型，必须是先预扣(预占资源)的模型。

### 分布式事务 - 幂等
- 很严重问题就是消息重复投递并被重复消费
- 例如：如果相同的消息被重复投递两次，那么我们余额宝账户将会增加2万而不是1万了。
  - 比如余额宝处理完消息 msg 后，发送了处理成功的消息给支付宝
  - 正常情况下支付宝应该要删除消息msg，
  - 但如果支付宝这时候悲剧的挂了，重启后一看消息 msg 还在，就会继续发送消息 msg。
  - 解决->不要重复消费消息   
- 全局唯一ID + 去重表
  - 在余额宝这边增加消息应用状态表 msg_apply（该表为去重表），用于记录消息的消费情况，
  - 每次来一个消息，同一事务下
      - 先去消msg_apply中查询，（全局唯一ID保证消息唯一性）
      - 如果找到说明是重复消息，丢弃即可，
      - 如果没找到，进行业务操作，更新balance，同时插入到msg_apply

### 分布式事务 - 2PC

- 两阶段提交协议(Two Phase Commitment Protocol)中，涉及到两种⻆色
  - 一个事务协调者(coordinator):负责协调多个参与者进行事务投票及提交(回滚) 
  - 多个事务参与者(participants):即本地事务执行者
- 总共处理步骤有两个
  - (1)投票阶段(voting phase):
    - 协调者将通知事务参与者准备提交或取消事务，然后进入表决过程。
    - 参与者将告知协调者自己的决策
       - 同意(事务参与者本地事务执行成功，但未提交)
       - 取消(本地事务执行故障)
  - (2)提交阶段(commit phase)
    - 收到参与者的通知后，协调者再向参与者发出通知，
       - 根据反馈情况决定各参与者是否要提交还是回滚
> https://en.wikipedia.org/wiki/Two-phase_commit_protocol
> - 关键是怎么定义分布式事务，因为按原生定义，分布式事务是局限在数据库事务这个定义下的
> - 2PC其实只是一种思路，一种协议。但因为年代久远，所以也和基于数据库事务的分布式事务绑定在一起而造成歧义。
>    - 应该和 X/Open Distributed Transaction Processing (DTP) Model (X/Open XA)
>    - EJB/MTS
>    - 这样已经淘汰的技术划清界限，因为已经没人在生产上再使用基于数据库事务的分布式事务。
> - 分布式事务越来越是一个基于服务级别的概念，事务本身这个词已经脱离了数据库事务这个定义。
> - 理解2PC应该回归本源，所谓的提交，理解应该更加抽象，而非和数据库事务绑定。
>
### 分布式事务 - 2PC Message Queue

- 思路按2PC，但是使用消息队列取代协调者角色。
- 通过消息中间件对分布式事务进行解耦

- 1. 生产者发送`prepare`消息到消息队列
- 2. 生产者执行本地事务。
     - 如果成功，生产者发送`commit-msg`
     - 如果失败，生产者发送`rollback-msg`
  3. 当`commit-msg`后，消息队列通知下游
     - 或者通知
     - 或者某订阅可以被消费
  4. 消费者在第三步后，消费到消息
     - 消费者执行本地事务（消费者需要保证业务幂等性）即对消费消息本身去重
     - 本地事务如果成功，消费者发送`ack-msg`
  5. 消费队列获知`ack`，则这个消费本正式可以从队列中剔除（被真正消费了）
  

- 注意：对于消费者执行本地事务一直失败的情况，采取**人工介入**的方式，
   - 而不是整个回滚，因为整个回滚代价太大，更容易出bug
   - 实践证明，这时候人工介入，而非系统回滚的方式是有效的。（支付宝实践的解决办法）

### 分布式事务 - Seata 2PC
> Seata
>  - https://github.com/seata/seata
>  - 阿里开源，java实现（也有go版本，非官方开发，还不成熟）
>  - 传统的2PC是以数据库为中心，而Seata是以服务为中心实现2PC，即基于微服务架构去做2PC

- Seata 2PC 与传统 2PC 的差别
 - 架构层次方面
   - 传统 2PC 方案的 RM 实际上是在数据库层
     - RM 本质上就是数据库自身，通过 XA 协议实现
   - 而 Seata的RM 是以 jar 包的形式作为中间件层部署在应用程序这一侧
  
 - 两阶段提交方面
   - 传统 2PC无论第二阶段的决议是 commit 还是 rollback
      - 事务性资源的锁都要保持到 Phase2 完成才释放。
   - Seata 的做法是在 Phase1 就将本地事务提交
      - 这样就可以省去 Phase2 持锁的时间，整体提高效率。

### 分布式事务 - TCC

TCC 是 **Try、Confirm、Cancel 三个词语的缩写**

TCC 要求每个分支事务实现三个操作: 预处理 Try、确认 Confirm、撤销 Cancel。 
   - Try 操作做业务检查及资源预留
   - Confirm 做业务确认操作
   - Cancel 实现一个与 Try 相反的操作即回滚操作。

TM首先发起所有的分支事务的Try 操作
  - 任何一个分支事务的 Try 操作执行失败
    - TM 将会发起所有分支事务的 Cancel 操作
  - 若 Try 操作全部成功
    - TM 将会发起所有分支事务的Confirm 操作
      - 其中 Confirm/Cancel 操作若执行失败，TM会进行重试。

需要注意:
- 幂等
- 空回滚
  - try的update操作没有执行成功，那么回滚时候需要考虑。
  - 没有调用Try的情况下，调用了 Cancel，Cancel需识别出这是一个空回滚，然后直接返回成功。  
- 防悬挂
  - cancle比try还先到达的情况。交易应该彻底结束，try不能成功。
  - 出现原因是在 RPC 调用分支事务 Try 时，先注册分支事务，再执行 RPC 调用，
     - 如果此时 RPC 调用的网络发生拥堵，通常 RPC 调用是有超时时间的，
       - RPC 超时以后，TM 就会通知 RM 回滚该分布式事务，
       - 可能回滚完成后，RPC 请求才到达参与者真正执行，
       - 而一个 Try 方法预留的业务资源，只有该分布式事务才能使用，
       - 该分布式事务第一阶段预留的业务资源就再也没有人能够处理了，
       - 对于这种情况，我们就称为悬挂，即业务资源预留后没法继续处理。  

### Saga模式
- 强调同步的模型，不使用传统2PC来解决问题。
  > - 注意这里(microservices.io/patterns) 说的分布式事务和本文讨论的分布式事务存在歧义，更多指旧的传统数据库的2PC。
  > - 本文说的更接近event-sourcing这种模式。
- 托管给编排服务。编排服务决定如何进行服务调用和顺序。


## References

- [Seata实战-分布式事务简介及demo上手](https://blog.csdn.net/hosaos/article/details/89136666)
- [面试必问：分布式事务六种解决方案](https://zhuanlan.zhihu.com/p/183753774)
- [分布式事务有这一篇就够了](https://www.cnblogs.com/dyzcs/p/13780668.html)
- [漫画：什么是分布式事务？](https://blog.csdn.net/bjweimengshu/article/details/79607522)
- https://microservices.io/patterns/data/event-sourcing.html 
- https://microservices.io/patterns/data/saga.html 
- http://chrisrichardson.net/post/microservices/2019/07/09/developing-sagas-part-1.html
- https://microservices.io/patterns/data/polling-publisher.html 
- https://microservices.io/patterns/data/polling-publisher.html 
- https://microservices.io/patterns/data/transaction-log-tailing.html
