# 第6周 评论系统架构设计 

## 概述
- 架构设计最重要的就是理解整个产品体系在系统中的定位。
- 不要做需求的翻译机，先理解业务背后的本质，事情的初衷。
- 搞清楚系统背后的背景，才能做出最佳的设计和抽象。
- 架构设计等同于数据设计，梳理清楚数据的走向和逻辑。尽量避免环形依赖、数据双向请求等。
- 大型系统最关键的是：缓存c，消息队列m，存储d（按重要性排序），一定要把缓存和消息队列用好。
  - 如果按8/2原则则是：c/m 80% d 20% 的去解决问题的原则
  - 因为d是难以被scale的（太复杂），最容易成为系统瓶颈。
  - 所以把上层（c/m）做好解决大部分问题。
  - 朴素架构的核心
    - 不依赖过于复杂的存储层扩展模式
    - 而是通过减少存储层的负载来解决问题。

### 性质分析
- 评论系统是读/写均衡，典型的读和写都密集型系统。

## 功能模块

### 评论系统功能模块
- 发布评论: 支持回复楼层、楼中楼。
- 读取评论: 按照时间、热度排序。
- 删除评论: 用户删除、作者删除。
- 管理评论: 作者置顶、后台运营管理(搜索、删除、审核等)。

## 架构设计
- BFF: comment
  + gPRC
  + 复杂评论业务的服务编排，
     + 比如访问账号服务进行等级判定，
  + 在 BFF 面向移动端/WEB场景来设计 API 
  + BFF抽象把评论的本身的内容列表处理(加载、分页、排序等)进行了隔离，关注在业务平台化逻辑上。
- Service: comment-service
  + 服务层，去平台业务的逻辑，专注在评论功能的 API 实现上，
     + 比如发布、读取、删除等，
     + 上游去做平台化的组织能力
        - 关注分离：上层(BEF)关注平台化接入
        - 关注分离：下层(服务层）关注评论场景本身（数据的写/读/管理）
  + 关注在稳定性、可用性上，这样让上游可以灵活组织逻辑把基础能力和业务能力剥离。
  + kafka/redis/mysql   
- Job: comment-job
  + 订阅kafka消息。
  + 消息队列的最大用途是消峰处理（数据先入消息队列，再慢慢处理）
- Admin: comment-admin
  + 管理平台，按照安全等级划分服务，尤其划分运营平台。
    - 独立对于安全很重要，这样可以保证服务不会开给前端
  + 共享服务层的存储层(MySQL、Redis)。
  + 运营体系的数据大量都是检索：
     - 使用 canal 进行同步到 ElasticSearch 中，
     - 整个数据的展示都是通过 ES，再通过业务主键更新业务数据层，
     - 运营端的查询压力下放给了独立的 fulltext search 系统。
     - 也就是说运营体系不要读源数据，而是搜到到对象之后，再操作源数据。
- Dependency: account-service、filter-service
  + 整个评论服务还会依赖一些外部 gRPC 服务，
  + 统一的平台业务逻辑在 comment BFF 层收敛，
  + account-service 账号服务，filter-service 敏感词过滤服务。
### comment-service
 - 专注在评论的**数据处理**
   + 仔细考虑关注分离： Separation of Concerns。
      - 策略是多变的，和平台相关的扔到上层去。（先审后发/先发后审。。）
      - 读写（数据处理）是稳定的（存储的构建，评论的读/写/排序/删除/）
 - 读的核心逻辑:
   + Cache-Aside 模式
     - 先读取缓存，再读取存储。
     - 一般会使用 read ahead 的思路，即预读，
     - 用户访问了第一页，很有可能访问第二页，所以缓存会超前加载，避免频繁 cache miss。
   + 早期 cache rebuild 是做到服务里的，对于重建逻辑，
     - 当缓存抖动时候，特别容易引起集群 thundering herd 现象，
        - 大量的请求会触发 cache rebuild，
        - 因为使用了预加载，容易导致服务 OOM。
   + 解决：回源的逻辑, 不要把缓存的rebuild(写)放在服务层，使用comment-job处理cache rebuild
     - 我们使用了消息队列来进行逻辑异步化，
       - 读请求先读缓存，cache miss，读db，直接返回。同时不rebuild，而是向消息队列投递一个cache miss消息。
     - comment-job读消息队列  
     - 对于当前请求只返回 mysql 中部分数据即止。
- 写的核心逻辑:
     - 写和读相比较，写可以认为是透穿到存储层的，
     - 系统的瓶颈往往就来自于存储层，或者有状态层
     - 把写压力先放给kafka（消息队列吞吐大，因为顺序写）
       - 这样不会写堆积在服务层，而只是堆积在消息队列里面。
       - 同样压力也不会直接冲到存储层，因为kafka进行了消峰。
     - 刚发布的评论有极短的延迟(通常小于几 ms)对用户可见是可接受的
       - 把对存储的直接冲击下放到消息队列，
       - 按照消息反压的思路，
          - 即如果存储 latency 升高，消费能力就下降，自然消息容易堆积，
          - 系统始终以最大化方式消费。
     - Kafka 是存在 partition 概念的，处理回源消息也是类似的思路。
       - 可以认为是物理上的一个小队列， 一个 topic 是由一组 partition 组成的，
       - 所以 Kafka 的吞吐模型理解为: 全局并行，局部串行的生产消费方式。
       - 对于入队的消息，可以按照 hash(comment_subject) % N(partitions) 的方式进行分发。 那么某个 partition 中的 评论主题的数据一定都在一起，这样方便我们串行消费。
         - 问题：如何解决 “明星出轨”等热点事件的发生，
    
### comment-admin

- mysql binlog 中的数据被 canal 中间件流式消费，
- 获取到业务的原始 CRUD 操作，需要回放录入到 es 中
  - 但是 es 中的数据最终是面向运营体系提供服务能力，需要检索的数据维度比较多，
  - 在入 es 前需要做一个异构的 joiner，把单表变宽预处理好 join 逻辑，然后倒入到 es 中。
- 一般来说，运营后台的检索条件都是组合的，
  - 使用 es 的好处是避免依赖 mysql 来做多条件组合检索
    - 这样不用给mysql加大量索引（索引越多，写入越慢）
  - 同时 mysql 毕竟是 oltp 面向线上联机事务处理的。 通过冗余数据的方式，使用其他引擎来实现检索。
    - 目的是mysql只面对写操作，和简单读。复杂读不去mysql
- 这样对于运营来说，因为异步流程的delay是完全可以容忍的。  
- es 一般会存储检索、展示、primary key 等数据，
  - 操作编辑时，找到记录的 primary key，最后交由 comment-admin 进行运营测的 CRUD 操作。
- b站的内部运营体系基本都基于 es 来完成的。

### comment 
- comment 作为 BFF，是面向端，面向平台，面向业务组合的服务。
- 平台扩展的能力，我们都在 comment 服务来实现，方便统一和准入平台，
- 以统一的接口形式提供平台化的能力。
   - 依赖其他 gRPC 服务，整合统一平台测的逻辑(比如发布评论用户等级限定)。
   - 直接向端上提供接口，提供数据的读写接口，甚至可以整合端上，提供统一的端上 SDK。
   - 需要对非核心依赖的 gRPC 服务进行降级，当这些服务不稳定时。

## 存储设计

### 数据库设计

- 三张表
  - comment_subject  (主题，根评论总数/总评论数)
  - comment_index    (索引表，楼层/数量，评论的位置，数量，点赞)
  - comment_content  (内容表，内容/ip/设备)
- 数据写入: 事务更新    
  - content 属于非强制需要一致性考虑的。
    - 可以先写入 content，之后事务更新其他表。
    - 即便 content 先成功，后续失败仅仅存在一条 ghost 数据。
  - subject和index两张必须同步更新。
  - 只有subject/index这个事务完成，那么content才可见。

- 数据读取: 
  + 基于 obj_id + obj_type 在 comment_index 表找到评论列表
  + 父楼层 
    - WHERE root = 0 ORDER BY floor。
    - 之后根据 comment_index 的 id 字段捞出 comment_content 的评论内容。
  + 对于二级的子楼层，WHERE parent/root IN (id...)。
    - 因为产品形态上只存在二级列表，因此只需要迭代查询两次即可。
    - 对于嵌套层次多的，产品上，可以通过二次点击支持。
  - 未来可以 Graph 存储
    - DGraph、HugeGraph 类似的图存储思路。
    - 避免嵌套。

### 索引/内容分离
- comment_index: 
   - 评论楼层的索引组织表，实际并不包含内容。
- comment_content: 
   - 评论内容的表，包含评论的具体内容。
- 其中 comment_index 的 id 字段和 comment_content 是1对1的关系，这里面包含几种设计思想。
  - 表都有主键，即 cluster index，是物理组织形式存放的，
  - comment_content 没有 id，是为了减少一次二级索引查找，直接基于主键检索，
  - 同时 comment_id 在写入要尽可能的顺序自增。
- 索引、内容分离，方便 mysql datapage 缓存更多的 row，
  - 如果和 content 耦合，会导致更大的 IO。
- 长远来看 content 信息可以直接使用 KV storage 存储。

### 缓存设计

- comment_subject_cache:
  - value 使用 protobuf 序列化的方式存入。
  - 我们早期使用 memcache 来进行缓存，因为 redis 早期单线程模型，吞吐能力不高。 
- comment_index_cache: 
  - 使用 redis sortedset 进行索引的缓存，
  - 索引即数据的组织顺序，而非数据内容。
  - 参考过百度贴吧，百度使用自研的拉链存储来组织索引，不需要
    - mysql 作为主力存储，利用 redis 来做加速完全足够，
  - 因为 cache miss 的构建，前面讲过使用 kafka 的消费者中处理，
     - 预加载少量数据，通过增量加载的方式逐渐预热填充缓存，
     - 而 redis sortedset skiplist 的实现，可以做到 O(logN) + O(M) 的时间复杂度，效率很高。
     - sorted set 是要增量追加的，因此必须判定 key 存在，才能 zdd。
- comment_content_cache: 
  - 使用 protobuf 序列化的方式存入。类似的我们早期使用 memcache 进行缓存。
  
- 增量加载 + lazy 加载

## 可用性设计

### Singleflight

- 对于热门的主题，如果存在缓存穿透的情况，会导致大量的同进程、跨进程的数据回源到存储层，可能会引起存储过载的情况，
- 只交给同进程内一个人去做加载存储
   + 使用归并回源的思路:
      - https://pkg.go.dev/golang.org/x/sync/singleflight
      - 同进程只交给一个人去获取 mysql 数据，然后批量返回。
- 同时这个 lease owner 投递一个 kafka 消息，做 index cache 的 recovery 操作。
   - 这样可以大大减少 mysql 的压力，
   - 以及大量透穿导致的密集写 kafka 的问题。
- 更进一步的，后续连续的请求，仍然可能会短时 cache miss，
   + 我们可以在进程内设置一个 short-lived flag，
     - 标记最近有一个人投递了 cache rebuild 的消息
     - 发现标记请求则直接 drop。
 
总结
- 消息生产者进程内
   - 用归并（singleflight）读mysql
   - 用归并写kafka
- 消息消费者进程内
   - 用预判定缓存是否存在，是否丢弃cache rebuild指令。
   - 当rebuild没有完成时，在进程内使用短时flag去丢弃rebuild指令。
- 扩展
   - 使用一个短时的LRU缓存，去减少对mysql的读。
     - 例：5s内有一个人查了mysql，就缓存住，其它人查，走local cache
         - local cache 5s后失效 
         - 这时redis早已构建完成
     - LRU（Least recently used）指缓存淘汰机制。
        + 如果缓存已满，删除访问时间最早的数据。(最近最少使用的)
  
- 为什么我们不用分布式锁之类的思路？
  - 太复杂
  - 不想依赖一个集中式状态服务


### 热点

#### 问题
- 流量热点是因为突然热门的主题，被高频次的访问，
- 因为底层的 cache 设计，一般是按照主题 key 进行一致性 hash 来进行分片，
- 但是热点 key 一定命中某一个节点，
- 这时候 remote cache 可能会变为瓶颈，-> 单个redis被打爆。

#### 单进程自适应发现热点，在进程内吞掉大量的读请求
- cache 的升级 -> remote ->  local 
- 利用框架的能力自动发现热点。
- 使用单进程自适应发现热点的思路，附加一个短时的 ttl local cache，在进程内吞掉大量的读请求。
- 实现  
  - 在内存中使用 hashmap 统计每个 key 的访问频次，
  - 使用一个环形数组，过去数据清零。  
  - 使用滑动窗口统计，
     - 即每个窗口中，维护一个 hashmap，之后统计所有未过去的 bucket，汇总所有 key 的数据。
     - 之后使用小堆计算 TopK 的数据，自动进行热点识别。

## References