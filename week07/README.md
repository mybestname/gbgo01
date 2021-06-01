# 第7周 播放历史记录架构设计

## 目录
- 功能模块
- 架构设计
- 存储设计
- 可用性设计 
- References

## 功能模块

- 实例：https://www.bilibili.com/account/history
- 用户需求：
  - 播放历史查看
  - 播放进度同步。
  - 离线型用户，app本地保留历史记录数据。
  - 考虑平台化，视频、文章、漫画 等业务扩展接入。
- 功能划分    
  - 变更功能
    - 添加记录
    - 删除记录
    - 清空历史。
  - 读取功能
    - 按照timeline返回topN，点查获取进度信息。
  - 其他功能
    - 暂停/恢复记录
    - 首次观察增加经验
- 分析
  - 一个极高 tps 写入，高 qps 读取的业务服务。
  - **原因** 
    - 播放进度要不断上报云端进行同步（极高tps）
    - 每点击一个视频都要获取进度（高qps），无法预判点哪个，如何设计缓存？
  - 注意
    - 分析清楚系统的hot path再优化，不哪哪都优化。
    - 收益/产出

## 架构设计 

###  概览 

#### history BFF
- 接受来自外部用户的**读请求**，
  - 依赖其他(例如archive、comic)服务来组装完整的history⻚面需要的数据组合。
- 同时接受来自内部其他业务线的**写请求**，
  - 业务方自己进行业务ID的判定
  - 投递到history BFF写接口。
- history BFF 打包在 app-interface BFF (大杂烩BFF) 中
  - 考虑到隔离性，读写流量很大，独立成 history BFF。

#### Service: history-service 

- serivce层的特点
  - 去平台业务的逻辑，专注在数据的持久化
  - BFF专注平台业务数据组织，service负责数据的读、写、删、清理等操作。
  
- history service的特点
  - **播放进度**非常高频同步，需考虑性能优化。
  
- **核心解决思路** 
  - 使用 **write-back思路** (缓存回写）
    - **状态数据先入分布式缓存，再回写数据库**。
    - 写极高时，先别急写状态存储，先写缓存是一个常见思路
      - 存储的raid的思路
      - linux内核对I/O的page cache的思路
      - 游戏时候掉装备（回档）的原因（大规模PK时候常见），因为写入没有持久而是在缓存。
  - **播放进度**写入kafka时候，**先聚合数据**
    - 不要来一条记录写入一次kafka，而是把零散小消息，变成一个大消息，一次IO
      - 一个kafka message可能有100个用户消息。
    - kafka里面记录的是事件，相应于事件发生序，而redis里面是事件的内容。
    - job消费kafka消息，并读redis内容构造为完整的数据记录，再写入hbase。
    - 这样就等于有三级的写入差速：redis，kafka，hbase，来平滑掉IO
  
- Job: history-job
  - job 消费上游kafka 的数据，
  - 利用消息队列的堆积能力，对于存储层的差速(消费能力跟不上生产速度时)，可以进行一定的数据反压。
  - 配合上游 service 批量打包过来的数据持久化。

- Upstream: some-app，some-api 
  - 整个历史服务还会被一些外部 gRPC 服务所依赖，
  - 所以 history 还充当了内网的 gRPC Provider， 这些上游服务，使用历史服务的写接口，把自己业 务的数据进行持久化。
  
- **History Service最重要的设计**
  - 批量打包(pipeline)聚合数据。
  - 将高频、密集的写请求先入缓存(write-back)
  - 批量消费 减少对存储的直接压力。

### 架构设计 - history-service

history-service 专注在历史数据处理。 

#### 写的核心逻辑: 

在 in-process 内存中，定时定量聚合数据，使用 kafka 消息队列来消除写入峰值。

  - 用户观看的稿件、漫画等，带有进度信息的数据，同一个id 最后一次的数据即可
  - last-write win，高频的用户端同步逻辑，只需要**最后一次数据持久化**即可。
  - 聚合不同用户的 **“同一个对象的最后一次进度”** ，

我们需要保证用户数据可以实时被观察到

  - 不能出现上报进度后，需要一阵子才能体现进度变化。
  - 在内存中打包数据，同时实时写入到 redis中，这样即保证了实时，又避免海量写入冲击存储。
    
kafka 是为高吞吐设计，超高频的写入并不是最优

  - **内存聚合和分片算法比较重要**，
  - 按照 uid 来sharding 数据，写放大仍然很大，
  - 使用 **region sharding**，打包一组数据当作一个 kafka message
     - 比如 uid % 100数据打包
    
写逻辑的数据流向
 - 实时写 redis -> 内存维护用户数据 -> 定时/定量写入到 kafka。
       
#### 读的核心逻辑:
- 历史数据，实时写入 redis 后，不会无限制的存储， 会按量截断，
  - 所以分布式缓存中数据不是完整数据，
- 历史数据从 redis sortedset 中读取后，如果发现尾部数据不足，会触发 **cache-aside 模式** 
  - 从存储中回捞数据，但是**不会重新回填缓存**
    - 因为拉取过去更久远的数据，属于用户纬度的低频度行为。
  > cache-aside是一种缓存更新模式 : 指不更新缓存，只是让缓存失效的模式，
  >  - 写操作马上让缓存失效，后续的读操作再把数据从数据库中拉出来。
  >  - 参考：https://coolshell.cn/articles/17416.html
  >  - 但作者似乎并不会更新缓存，而只是意思是从存储上读数据，而ingnore缓存这里和cache-aside的原意是不同的。
- 历史数据通常是按照 timeline 来组织，游标的 key 可以使用时间戳进行翻⻚或者下拉。

### 架构设计 - history-job 
- history-job，获取打包好的用户数据，进行批量持久化。
- 上游 history-service 
  - 按照 uid region sharding 聚合好的数据， 在 job 中消费取出，
  - 节约传输过程，以及 history-service 的 in-process cache 的内存使用
     - **只维护用户的 uid 以及 id 列表** 
        - kafka最小化存储和传输。 
  - 因为数据是不完整的，需要从 redis 中获得id对应的数据内容，再持久化。
     - 从原来的N条记录变为一个用户一条记录。
- HBase 非常合适高密度写入。
   - **注意job是batch write**，这样最大化利用了db的I/O

### 架构设计 - history BFF
- history 作为 BFF，
  - 对 用户端 提供统一的用户记录入口的接口，
  - 对 内 提供 gRPC 写入历史接口。
- 如果业务场景中不存在统一的用户入口
  - 可以去掉BFF层，直接使用history-service提供读接口，
  - 这样需要每个业务方自己实 现自己的数据组装。 
- BFF层为何要依赖redis => 解决B站内部功能 
  - 功能：根据用户是否首次播放、观看等加经验 奖励积分
       - 读redis，判定用户当天是否是首次访问
        - 容易想到使用 bitmap 或者 bloom filter 来进行判断
       - 然后往下游 kafka 投递消息，而不直接依赖后续业务。(加经验/积分业务解耦)
  - 功能：关闭历史记录的功能 （暂停历史记录）
       - 每次写入操作读一次redis（flag标签），是否打开了开关，
  - 功能：首次发送奖励
       - 每次写入操作读一次redis（flag标签），是否打开了开关，
  - 是否有更好的方法？

## 存储设计 

### 数据库设计 
- 存储选型: HBase。

#### 数据写入
- PUT mid, values，只需要写入到 column_family 的 info 列簇
- rowkey 使用用户 id md5 以后的头两位 + 用户，避免 rowkey 热点密集到一个 region 中，导致写/读热点。
- 对于 column_family: info，存储一个列 obj_id + obj_type，
  例如 稿件业务:1、稿件ID: 100， 100_1 作为列名，
- 对于 value 使用 protobuf 序列 化一个结构体接入。
  - 所以只需要单次更新 kv store。
- 使用 HBase TTL 的能力，只需要 保存90天的用户数据即可。
  - (删除同理)
  
#### 数据读取

- 列表获取为 GET mid，直接获取1000 条，在内存中排序和翻⻚。
- 点查 GET mid columns，在茫茫多视频查看当前视频的阅读进度，
   - cache miss 会非常严重，
   - 虽然支持点查，但是对于上层 cache miss 后，不再回源请求 HBase。

### 存储设计 - 缓存设计

#### 数据写入

每次产生的历史数据，需要立⻢更新 
   - redis，使用 **sorted set** 基于时间排序的列表
   - member 为业务 ID。同时存储一份数据到 redis string 中
      - `(obj_id + obj_type)` 业务唯一ID
   - 使用 protobuf 序列化完整的数据内容。
   - 为了避免 redis 中单个用户数据无 限增⻓，需要超过一定量后对数据进行截断。

#### 数据读取 

分为两个场景，
  - 一个是历史⻚面， 这时候使用 sorted set，排序查找即可
     - 拿到 列表后，**mget 批量获取** history_content 内容。 
  - 一个是点查进度，比如我们点击进入一个视频详情⻚
     - 这时候直接查找 history_content 进行点查
     - 不再回源 HBase，因为命中率太低。
#### 缓存设计

我们在缓存设计中，经常使用 bitmap(roaring bitmap)、 bloom filter 缓存加速访问，

但是在使用缓存时， **需要注意规避热点问题**，
 - 某个key sharding 命 中 node 是固定的，
 - 因此我们可以利用**构建多组** bitmap 或 bloom filter，来**进行打散**。 
    - 例如 `prefix_key = hash(mid) % 1000`
    - 根据 `prefix_key` 找到对应的 cache 再进行操作，
    - 这样 `1000` 个 key 尽可能均匀的分布到更小集合的node，而不会产生数据热点。

但是仍然每次触发行为，都为前置判定，有更 好的优化方案吗?
  - bloom filter 的误判率，可以前置计算预估下

## 可用性设计

### 可用性设计- Write-Back

1. history-service 中实时写入 redis 数据，
  - 因此只需要重点优化缓存架构中，扛住峰值的流量写入。
  - 之后在服务内存中，
    - 使用 `map[int]map[int]struct{}` 聚合数据，
      - int_1 : 用户id
      - int_2 ：obj_id
    - 利用 `chan` 在内部发送每个小消息，再聚合成一个大`map`
    - 在 `sendproc` 中，使用 timer 和 定量判定逻辑，发送到下游 kafka 中。
  
2. history-job 中
  - 获取消息后，重新去 redis 中回捞数据即: history-content
  - 然后构建完整的数据批量写入到 HBase 中。
    
存在两个⻛险:
  - 1、history-service 重启过程中，预聚合的消息丢失; (map丢了)
  - 2、history-job 读取 redis 构建数据，但 redis 丢失;

风险说明：    
- 我们在这里进行了 trade-off，高收敛比的设计，意味着**存在数据丢失的⻛险**，
- 对于历史场景，非 L0 的业务服务/ 数据，我们认为极端情况下**可接受**。

### 可用性设计 - 聚合

经过 BFF history 的流量 per-request 都会发送给 history-service，
  - 最容易想到的优化就是**聚合上移**来减少发送给下游的 rpc。
  
但是按照 mid region sharding 的思路非常具有业务的耦合性

  - 所以**不应该把逻辑上移，而只是数据上移**，
  - 所以可以考虑简单 batch put 请求，做一个**无逻辑的数据聚合**
  - 再发送给 history-service，这样可以大大的减少内网的流量，节约资源。
    
经过 API Gateway 的流量都会触发高频的 per-rpc auth，给内网的 identify-service 带来了不少压力。
 - 我们认为大部分历史行为通过心跳的方式同步进度 
 - 连接一个⻓连接，握手后先进行用户级的身份验证，之后维持身份信息，
 - 不是每次发送 request 都进行验证，
 - 这样可以大大减少内网的 identify-service 的流量。
 - 我们内网使用 boardcast(goim) 服务维护⻓连接，⻓连接一次验证，不断使用。

### 可用性设计 - 广播

用户首次触发的行为，需要发送消息给下游系统进行 触发其他奖励等。

如何减少这类一天只用一次的标记位缓存请求? (前述的redis中为了该功能设的flag位)

思考🤔 用 in-process localcache，
  - 使用LRU维护这个集合，
  - 只有高频的用户访问收益，
  - 用户分布很广，很难覆盖，命中率很低。 
    
越源头解决架构问题，通常越简单，效率越高。 

我们在写操作(高频请求)中，
  - 把当前的 flag 返回到 API 协议中，作为一个日期值，
  - 客户端保存到本地，下次请求的时候带上，
  - 如果发现该值在，获取以后直接使用不再请求缓存，
  - 例如:
    - 2021-1-1，发现当前时间还是 2021-1-1，
      - 直接不再请求 redis，
    - 如果发现当前时间是 2021-1-2，
      - 需要触发一次 redis 访问，
      - 返回新的 flag 到客户端，
  - 这样把状态广播同步到任何其他设备，可以大大减少判定缓存。

实现成本在于，你认为的代价高低。

## References
- https://en.wikipedia.org/wiki/Cache#Writing_Policies 
- https://blog.csdn.net/jiaomeng/article/details/1495500 
- https://blog.csdn.net/yizishou/article/details/78342499 
- https://blog.csdn.net/caoshangpa/article/details/78783749

