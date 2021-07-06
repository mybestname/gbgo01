# 第12课 消息队列 - Kafka

## 目录
- Kafka 基础概念
- Topic & Partition 
- Producer & Consumer 
- Leader & Follower
- 数据可靠性
- 性能优化

## Kafka
- Kafka 是由 LinkedIn 开发并开源的分布式消息系统，因其分布式及高吞吐率而被广泛使用，
  现已与 Cloudera Hadoop，Apache Storm，Apache Spark，Flink 集成。
- 常做为多种类型的数据管道和消息系统使用。
   - 行为流数据
     - 几乎所有站点在对其网站使用情况做报表时都要用到的数据。
     - ⻚面访问量 PV、
     - ⻚面曝光 Expose、
     - ⻚面点击 Click 等行为事件;
   - 实时计算中经常使用 Kafka 做为 Data Source 和 Dataflow Pipeline
   - 业务的消息系统，通过发布订阅消息解耦多组微服务，消除峰值;
     - 业务之间解耦，只在数据上依赖。
     - 海量流量，解决流量差速。进行消峰。

### Kafka 简介
- Kafka 是一种分布式的，基于发布/订阅的消息系统。
- 主要设计目标:
  - 以时间复杂度为 O(1) 的方式提供消息持久化能力，即使对 TB 级以上数据也能保证常数时间复杂度的访问性能;
  - 高吞吐率。即使在非常廉价的商用机器上也 能做到单机支持每秒 100K 条以上消息的传输; 
  - 支持 Kafka Server 间的消息分区，及分布式消费，
  - 同时保证每个 Partition 内的消息顺序传输;
  - 支持离线数据处理和实时数据处理; 
  - Scale out:支持在线水平扩展;


### 为何使用消息系统
- 解耦
  - 消息系统在处理过程中间插入了一个隐含的、基于数据的接口层，两边的处理过程都要实现这 一接口。
    这允许你独立的扩展或修改两边的处理过程，只要确保它们遵守同样的接口约束。
  - 而基于消息发布订阅的机制，可以联动多个业务下游子系统，能够不侵入的情况下分步编排和开发，来保证数据一致性。
- 冗余
  - 有些情况下，处理数据的过程会失败。除非数据被持久化，否则将造成丢失。
  - 消息队列把数据 进行持久化直到它们已经被完全处理，通过这一方式规避了数据丢失⻛险。
    许多消息队列所采用的”插入-获取-删除”范式中，在把一个消息从队列中删除之前，
    需要你的处理系统明确的指出该消息已经被处理完毕，从而确保你的数据被安全的保存
    直到你使用完毕。
- 扩展性
  - 因为消息队列解耦了你的处理过程，所以增大消息入队和处理的频率是很容易的，只要另外增加处理过程即可。
    不需要改变代码、不需要调节参数。扩展就像调大电力按钮一样简单。
- 灵活性 & 峰值处理能力
   - 在访问量剧增的情况下，应用仍然需要继续发挥作用，但是这样的突发流量并不常⻅;
     如果为以能处理这类峰值访问为标准来投入资源随时待命无疑是巨大的浪费。使用
     消息队列能够使关键 组件顶住突发的访问压力，而不会因为突发的超负荷的请求而完全崩溃。
- 可恢复性
   - 系统的一部分组件失效时，不会影响到整个系统。消息队列降低了进程间的耦合度，
     所以即使一个处理消息的进程挂掉，加入队列中的消息仍然可以在系统恢复后被处理。
- 顺序保证
   - 在大多使用场景下，数据处理的顺序都很重要。大部分消息队列本来就是排序的，并且能保证
     数据会按照特定的顺序来处理。Kafka 保证一个 Partition 内的消息的有序性。
- 缓冲
   - 在任何重要的系统中，都会有需要不同的处理时间的元素。
     消息队列通过一个缓冲层来帮助任务最高效率的执行———写入队列的处理会尽可能的快速。
     该缓冲有助于控制和优化数据流经过系统的速度。
- 异步通讯
   - 很多时候，用户不想也不需要立即处理消息。消息队列提供了异步处理机制，允许用户把一个
     消息放入队列，但并不立即处理它。想向队列中放入多少消息就放多少，然后在需要的
     时候再去处理它们。

## Topic & Partition
- Topic 在逻辑上可以被认为是一个 queue， 每条消费都必须指定它的 Topic，可以简单 理解为必须指明把这条消息放进哪个 queue 里。
  我们把一类消息按照主题来分类，有点类似于数据库中的表。
- 为了使得 Kafka 的吞吐率可以线性提高，物理上把 Topic 分成一个或多个 Partition。
  对应到系统上就是一个或若干个目录。

### Broker
- Broker
  - Kafka 集群包含一个或多个服务 器，每个服务器节点称为一个 Broker。
- Broker 存储 Topic 的数据。
  - 如果某 Topic 有 N 个 Partition，集群有 N 个 Broker，那么每个 Broker 存储该 Topic 的一个 Partition。
- 从 scale out 的性能⻆度思考，通过 Broker Kafka server 的更多节点，带更多的存储，建 立更多的 Partition 把 IO 负载到更多的物理节 点，提高总吞吐 IOPS。
- 从 scale up 的⻆度思考，一个 Node 拥有越多 的 Physical Disk，也可以负载更多的 Partition，提升总吞吐 IOPS。

#### Broker 例子
- 如果某 Topic 有 N 个 Partition，集群有 (N+M)个 Broker，那么其中有 N 个 Broker 存储该 Topic 的一个 Partition， 剩下的 M 个 Broker 不存储该 Topic 的 Partition 数据。
- 如果某 Topic 有 N 个 Partition，集群中 Broker 数目少于 N 个，那么一个 Broker 存储该 Topic 的一个或多个 Partition。
- Topic 只是一个逻辑概念，真正在 Broker 间 分布式的 Partition。
- 每一条消息被发送到 Broker 中，会根据 Partition 规则选择被存储到哪一个 Partition。 如果 Partition 规则设置的合理，所有消息可 以均匀分布到不同的 Partition中。

### Broker & Partition
- 实验条件:
  - 3个 Broker，1个 Topic，无 Replication，异步模式，3个 Producer，消息 Payload 为100字节:
- 当 Partition 数量小于 Broker个数时，
  - Partition 数量越大，吞吐率越高，且呈线性提 升。 
  - Kafka 会将所有 Partition 均匀分布到所有Broker 上，所以当只有2个 Partition 时，会有2个 Broker 为该 Topic 服务。3个 Partition 时同理会有3个 Broker 为该 Topic 服务。
- 当 Partition 数量多于 Broker 个数时，总吞吐量并 未有所提升，甚至还有所下降。
  - 可能的原因是，当 Partition 数量为4和5时，不同 Broker 上的 Partition 数量不同，而 Producer 会将数据均匀发送到各 Partition 上，
    这就造成各Broker的负载不同，不能最大化集群吞吐量。

### 存储原理：顺序IO和OS Page Cache
- Kafka 的消息是存在于文件系统之上的。 
- Kafka 高度依赖文件系统来存储和缓存消息。
- 操作系统将主内存剩余的所有空闲内存空间都用作磁盘缓存，
  所有的磁盘读写操作都会经过统一的磁盘缓存(除了直接 I/O 会绕过磁盘缓存)。
- 任何发布到 Partition 的消息都会被追加到 Partition 数据文件的尾部，
  这样的顺序写磁盘操作让 Kafka 的效率非常高。
- Kafka 利用**顺序IO**，以及**linux Page Cache** 达成超高吞吐。

### 存储原理：数据删除策略 和 Offset
- Kafka 集群保留所有发布的 message，不管这个 message 有没有被消费过，
  Kafka 提供可配置的保留策略去删除旧数据(还有 一种策略根据分区大小删除数据)。 
- 例如，如果将保留策略设置为两天，在 message 写入后两天内，它可用于消费，
  之后 它将被丢弃以腾出空间。
- Kafka 的性能跟存储的 数据量的大小无关， 所以将数据存储很⻓一段 时间是没有问题的。
- Offset:偏移量。
  - 每条消息都有一个当前 Partition 下唯一的 64 字节的 Offset，
    它是相当于当前分区第一条消息的偏移量，即第几条消息。
  - 消费者可以指定消费的位置信息，当消费者挂掉再重新恢复的时候，可以从消费位置继续消费。

### 存储原理 例子
- 假设我们现在 Kafka 集群只有一个 Broker，
   - 创建 2 个 Topic 名称分别为: 「Topic1」和「Topic2」，
   - Partition 数量分别为 1、2。
- 那么我们的根目录下就会创建如下三个文件夹:
```
 -- topic1-0
 -- topic2-0
 -- topic3-0
```
- 在 Kafka 的文件存储中，同一个 Topic 下有多个不同的 Partition，
  每个 Partition 都为一个目录。
- 而每一个目录又被平均分配成多个大小相等的 Segment File 中，
- Segment File 又由 index file 和 data file 组成，他们总是成对出现，
  - 后缀 ".index" 和 ".log" 分表表示 Segment 索引文件和数据文件。
```
 -- topic1-0
    -- 00000000000000000000.index 
    -- 00000000000000000000.log 
    -- 00000000000000368769.index 
    -- 00000000000000368769.log
    ...
 -- topic2-0
 -- topic3-0
```    

### 存储原理 Index文件
- 以索引文件中元数据 <3, 497> 为例，依次在数据文件中表示
  第 3 个 Message(在全局 Partition 表示第 368769 + 3 = 368772 个 message)
  以及该消息的物理偏移地址为 497。 
- 注意该 Index 文件并不是从0开始，也不是每次递增 1 的，
  这是因为 Kafka 采取稀疏索引存储的方式，每隔一定字节的数据建立一条索引。
- kafka 减少了索引文件大小，使得能够把 Index 映射到内存，降低了查询时的磁盘 IO 开销，
  同时也并没有给查询带来太多的时间消耗。
- 因为其文件名为上一个 Segment 最后一条消息的 Offset ，
  所以当需要查找一个指定 Offset 的 Message 时，
  通过在所有 Segment 的文件名中进行 二分查找就能找到它归属的 Segment。
- 再在其 Index 文件中找到其对应到文件上的物理位置，就能拿出该 Message。

### 存储原理 Message(log文件)
- Kafka 是如何准确的知道 Message 的偏移的呢? 
  这是因为在 Kafka 定义了标准的数据存储结构，
- Partition 中的每一条 Message 都包含了以下三个属性:
  - Offset:表示 Message 在当前 Partition 中的偏移量，是一个逻辑上的值，
    唯一确定了 Partition 中的一条 Message，可以 简单的认为是一个 ID。
  - MessageSize:表示 Message 内容 Data 的大小。
  - Data:Message 的具体内容。

### 存储原理:2个步骤查找
- 例如读取 offset=368776的 message，需要通过下面2个步骤查找。
  - 第一步查找 segment file 上述图为例，其中 00000000000000000000.index 
    表示最开始的文件，起始 偏移量(offset)为0。
    第二个文件 00000000000000368769.index 的消息量起始偏移量为 
    368770 = 368769 + 1，其他后续文件依次类推，
    以起始偏移量命名并排序这些文件，只要根据 offset 二分查找文件列表，
    就可以快速定位到具体文件。 
    当 offset=368776时 定位到00000000000000368769.index | log。
  - 第二步通过 segment file 查找 message 通过第一步定位到 segment file，
    当 offset=368776时，依次定位到 00000000000000368769.index 的
    元数据物理位置和 00000000000000368769.log 的物理偏移地址，
    然后再通 过00000000000000368769.log 顺序查找直到 offset=368776 为止。

### 存储原理: .timeindex 索引文件
- Kafka 从0.10.0.0版本起，为分片日志文件 中新增了一个 .timeindex 的索引文件，
  可以根据时间戳定位消息。
- 同样我们可以通过脚本 kafka-dump-log.sh 查看时间索引的文件内容。
- 首先定位分片，将 1570793423501 与每个分片的最大时间戳进行对比
  (最大时间戳取 时间索引文件的最后一条记录时间，如果时间为 0 则取该日志分段的
  最近修改时间)， 直到找到大于或等于 1570793423501 的日志分段，因此会定位
  到时间索引文件 00000000000003257573.timeindex，其最大时间戳为1570793423505。
- 重复 offset 找到 log 文件的步骤。

## Producer 如何选择分区
- Producer 发送消息到 Broker 时，会根据 Partition 机制选择将其存储到哪一个 Partition。
  - 如果 Partition 机制设置合理，所有消息可以均匀分布到不同的 Partition 里，这样就实现
    了负载均衡。
  - 指明 Partition 的情况下，
  直接将给定的 Value 作为 Partition 的值。
  - 没有指明 Partition 但有 Key 的情况下，
  将 Key 的 Hash 值与分区数取余得到 Partition 值。
  - 既没有 Partition 有没有 Key 的情况下，
  第一次调用时随机生成一个整数(后面每次调用 都在这个整数上自增)，
  将这个值与可用的分区数取余，得到 Partition 值，也就是常说的 Round-Robin 轮询算法。

### Producer 工作原理
- 为保证 Producer 发送的数据，能可靠地发送到指定的 Topic，
  Topic 的每个 Partition 收到 Producer 发送的数据后，
  都需要向 Producer 发送 ACK。
  如果 Producer 收到 ACK，就会进行下一轮的发送，否则重新发送数据。
- 选择完分区后，生产者知道了消息所属的主题和分区，
  它将这条记录添加到相同主题和 分区的批量消息中，
  另一个线程负责发送这 些批量消息到对应的 Kafka Broker。
- 当 Broker 接收到消息后，如果成功写入则返回一个包含消息的主题、分区及位移的 RecordMetadata 对象，否则返回异常。
- 生产者接收到结果后，对于异常可能会进行重试。

### Producer Exactly Once (Kafka对幂等性的支持)
- 0.11 版本的 Kafka，引入了幂等性:Producer 
  - 不论向 Server 发送多少重复数据， Server 端都只会持久化一条。
- 要启用幂等性，只需要将 Producer 的参数中 `enable.idompotence` 设置为 true 即可。 
- 开启幂等性的 Producer 在初始化时会被分配一个PID，
  - 发往同一 Partition 的消息会附带Sequence Number。
- Borker 端会对`<PID,Partition,SeqNumber>` 做缓存，
  - 当具有相同主键的消息提交时， Broker只会持久化一条。
#### 注意  
- **PID 重启后就会变化**
- **不同的 Partition 具有不同主键**
- 幂等性无法保证跨分区会话的 Exactly Once
- 建议不要使用Kafka的幂等，而是在消费端去做基于业务的幂等。

## Consumer

- 场景
  - 从 Kafka 中读取消息，并且进行检查，最后产生结果数据。
  - 创建一个消费者实例去做这件事情，如果生产者写入消息的速度比消费者读取的速度快怎么办呢?
  - 这样随着时间增⻓，消息堆积越来越严重。
  - 对于这种场景，我们需要增加多个消费者来进行水平扩展。
- Kafka 消费者
  - 是消费组的一部分，当多个消费者 形成一个消费组来消费主题时，每个消费者会收到不同分区的消息。
- 假设有一个 T1 主题，该主题有 4 个分区;
  - 同时 我们有一个消费组 G1，这个消费组只有一个消 费者 C1。
  - 那么消费者 C1 将会收到这 4 个分区的消息。

- 如果我们增加新的消费者 C2 到消费组 G1，那么每个消费者将会分别收到两个分 区的消息。
  - 相当于 T1 Topic 内的 Partition 均分给了 G1 消 费的所有消费者，在这里 C1 消费 P0 和 P2， C2 消费 P1 和 P3。

- 如果增加到 4 个消费者，那么每个消费者 将会分别收到一个分区的消息。
  - 这时候每个消费者都处理其中一个分区，满负 载运行。

- 但如果我们继续增加消费者到这个消费组，剩 余的消费者将会空闲，不会收到任何消息。
- 总而言之，我们可以通过增加消费组的消费者来进 行水平扩展提升消费能力。

- 这也是为什么建议创建主题时使用比较多的分区数， 
  - 这样可以在消费负载高的情况下增加消费者来提升性能。
- 另外，消费者的数量不应该比分区数多，因为多出来的消费者是空闲的，没有任何帮助。

- 如果我们的 C1 处理消息仍然还有瓶颈，我们如何 优化和处理?
  - 把 C1 内部的消息进行二次 sharding，
  - 开启多个 goroutine worker 进行消费，
  - 为了保障 offset 提交 的正确性，需要使用 watermark 机制，
  - 保障最小的 offset 保存，才能往 Broker 提交。

### Consumer Group
- Kafka 一个很重要的特性就是，只需写入一 次消息，可以支持任意多的应用读取这个消息。
- 换句话说，每个应用都可以读到全量的消息。 
- 为了使得每个应用都能读到全量消息，应用 需要有不同的消费组。
- 对于上面的例子，
  - 假如我们新增了一个新的消费 组 G2，而这个消费组有两个消费者如图。
  - 在这个场景中，消费组 G1 和消费组 G2 都能收 到 T1 主题的全量消息，
    在逻辑意义上来说它们 属于不同的应用。
  - 最后，总结起来就是:
    - 如果应用需要读取全量消息，那么请为该应用设置一个消费组;
    - 如果该应用消费能力不足，那么可以考虑在这个消费组里增加消费者。

### Consumer 重平衡(Rebalance) 
- 可以看到，当新的消费者加入消费组，它会消费一个或多个分区，而这些分区之前是由其他消费者负责的。
- 另外，当消费者离开消费组(比如重启、宕机等)时，它所消费的分区会分配给其他分区。
- 这种现象称为重平衡(Rebalance)。
- 重平衡是 Kafka 一个很重要的性质， 这个性质保证了高可用和水平扩展。
- 不过也需要注意到，在 重平衡期间，所有消费者都不能消费消息，因此会造成整个消费组短暂的不可用。
- 而且，将分区进行重平衡也会导致原来的消费者状态过期，从而导致消费者需要重新更新状态，这段期间也会降低消费性能。 
- 消费者通过定期发送心跳(Hearbeat)到一个作为组协调者(Group Coordinator)的 Broker 来保持在消费组内存活。
  - 这个 Broker 不是固定的，每个消费组都可能不同。
  - 当消费者拉取消息或者提交时，便会发送心跳。
  - 如果消费者超过一定时间没有发送心跳，那么它的会话(Session)就会过期，
  - 组协调者会认为该消费者已经宕机，然后触发重平衡。

- 可以看到，从消费者宕机到会话过期是有一定时间的，这段时间内该消费者的分区都不能进行消息消费。
  
- 通常情况下，我们可以进行优雅关闭，这样消费者会发送离开的消息到组协调者，
  这样组协调者可以立即进行重平衡而不需要等待会话过期。
  
- 在 0.10.1 版本，Kafka 对心跳机制进行了修改，将发送心跳与拉取消息进行分离，
  这样使得发送心跳的频率不受拉取的频率影响。
  
- 另外更高版本的 Kafka 支持配置一个消费者多⻓时间不拉取消息但仍然保持存活，
  - 这个配置可以避免活锁(livelock)。
  - 活锁，是指应用没有故障但是由于某些原因不能进一步消费。
  - 但是活锁也很容易导致连锁故障，当消费端下游的组件性能退化，那么消息消费会变的很慢，
    会很容易出发 livelock 的重新均衡机制，反而影响力吞吐。

### Consumer Group consumer_offsets
- Partition 会为每个 Consumer Group 保存 一个偏移量，记录 Group 消费到的位置。
- Kafka 0.9开始将消费端的位移信息保存在集群的内部主题(__consumer_offsets)中，
  - 该主 题默认为50个分区，每条日志项的格式都是:`<TopicPartition, OffsetAndMetadata>`，
  - key 为主题分区主要存放主题、分区以及消费组信息，
  - value 为 OffsetAndMetadata 对象 主要包括 位移、位移提交时间、自定义元数据等信息。

### Consumer Commit Offset
- 消费端可以通过设置参数 enable.auto.commit 来控制是自动提交还是手动，
  - 如果值为 true 则 表示自动提交，在消费端的后台会定时的提交消费位移信息，时间间隔由 auto.commit.interval.ms(默认为5秒)。
  - 可能存在重复的位移数据提交到消费位移主题中，因为每隔5秒会往主题中写入一条消息，
    不管是否有新的消费记录，这样就会产生大量的同 key 消息，其实只需要一条，
    因此需要依赖前面提到日志压 缩策略来清理数据。
  - 重复消费，假设位移提交的时间间隔为5秒，那么在5秒内如果发生了 rebalance，
    则所有的消费者会 从上一次提交的位移处开始消费，那么期间消费的数据则会再次被消费。
    
- 集中 Delivery Guarantee:
  - 读完消息先 commit 再处理消息。这种模式下，
    如果 Consumer 在 commit 后还没来得及处理消息就 crash 了，
    下次重新开始工作后就无法读到刚刚已提交而未处理的消息，
     - 这就对应于 At most once。 
   - 读完消息先处理再 commit。这种模式下，
     如果在处理完消息之后 commit 之前 Consumer crash 了，
     下次重新开始工作时还会处理刚刚未 commit 的消息，实际上该消息已经被处理过了。
     - 这就对应于At least once。

### Consumer Exactly Once
- Flink 提供的 checkpoint 机制，
  结合 Source/Sink 端配合支持 Exactly Once 语义，以 Hive 为例:
  - 1.从 Kafka 消费数据，写入到临时目录
  - 2.ck snapshot 阶段，
      - 将 Offset 存储到 State 中，Sink 端关闭写入的文件句柄，
      - 以及保存 ckid 到 State 中
  - 3.ck complete 阶段，
      - commit kafka offset，将 临时目录中的数据移到正式目录
  - 4.ck recover 阶段，
      - 恢复 state 信息，reset kafka offset
      - 恢复 last ckid，将临时目录的 数据移动到正式目录

### Push vs Pull
- 作为一个消息系统，Kafka遵循了传统的方式，
  - 选择由 Producer 向 Broker push 消息
  - 并由 Consumer 从 Broker pull 消息。
- 一些 logging-centric system，
  - 比如 Facebook 的 Scribe 和 Cloudera 的 Flume，采用 push 模式。
- 事实上，push 模式 和 pull 模式各有优劣。
  - push 模式很难适应消费速率不同的消费者，因为消息发送速率是由 Broker 决定的。
  - push 模式的目标是尽可能以最快速度传递消息，但是这样很容易造成 Consumer 来不及处理消息，
    - 典型的表现就是拒绝服务以及网络拥塞。
  - 而pull 模式则可以**根据 Consumer 的消费能力以适当的速率消费消息**。
- 对于 Kafka 而言，pull 模式更合适。
  - pull 模式可简化 Broker 的设计，Consumer 可**自主控制消费消息的速率**，
    同时 Consumer 可以自己控制消费方式——即可批量消费也可逐条消费，
    同时还能选择 不同的提交方式从而实现不同的传输语义。
  - 而 Pull 模式则可以根据 Consumer 的消费能力以适当的速率消费消息。
- Pull 模式不足之处是，
  - 如果 Kafka 没有数据，消费者可能会陷入循环中，一直返回空数据。
    - 因为消费者从 Broker 主动拉取数据，需要维护一个⻓轮询，
    - 针对这一点，Kafka 的消费者在消费 数据时会传入一个时⻓参数timeout。
      如果当前没有数据可供消费，Consumer 会等待一段时间之后 再返回，这段时⻓即为 timeout。
      
## Leader & Follower

### 0.8版本前 Replication机制的缺失
- Kafka 在0.8以前的版本中，并不提供 HA 机制，
  - 一旦一个或多个 Broker 宕机， 则宕机期间 其上所有 Partition 都无法继续提供服务。
  - 若该 Broker 永远不能再恢复，亦或磁盘故障，则 其上数据将丢失。
- 在 Kafka 在0.8以前的版本中，是没有 Replication 的，
  - 一旦某一个 Broker 宕机，则其上所有的 Partition 数据都不可被消费，
  - 这与 Kafka 数据持久性及 Delivery Guarantee 的设计目标相悖。
  - 同时 Producer 都不能再将数据存于这些 Partition 中。 
  - 如果 Producer 使用同步模式
    - 则 Producer 会在尝试重新发送 `message.send.max.retries`(默认值为 3)次后
      抛出 Exception，用户可以选择停止发送后续数据也可选择继续选择发送。
      - 而前者会造成数 据的阻塞，后者会造成本应发往该 Broker 的数据的丢失。
    - 如果 Producer 使用异步模式，则 Producer 会尝试重新发送 `message.send.max.retries`(默认值为 3)次
      后记录该异常并继续发送后续数据，
      - 这会造成数据丢失并且用户只能通过日志发现该问题。
  - 由此可⻅，在没有 Replication 的情况下，
    - 一旦某机器宕机或者某个 Broker 停止工作则会造成整个系统的可用性降低。
    - 随着集群规模的增加，整个集群中出现该类异常的几率大大增加，
    - 因此对于生产系统而言 Replication 机制的引入非常重要。

### Leader
- 引入 Replication 之后
  - 同一个 Partition 可能会有多个 Replica，
  - 而这时需要在这些 Replication 之间选出一个 Leader，
    - Producer 和 Consumer 只与这个 Leader 交互，
  - 其它 Replica 作为 Follower 从 Leader 中复制数据。
- 因为需要保证同一个 Partition 的多个 Replica 之间的数据一致性
  - 其中一个宕机后其它 Replica 必须要能继续服务并且即不能造成数据重复也不能造成数据丢失。
- 如果没有一个 Leader，所有 Replica 都可同时读/写数据，
  那就需要保证多个 Replica 之间互相 (N×N条通路)同步数据，
  数据的一致性和有序性非常难保证，大大增加了 Replication 实现的复杂性，
  同时也增加了出现异常的几率。
- 而引入 Leader 后，只有 Leader 负责数据读写，
  Follower 只向 Leader 顺序 Fetch 数据(N条通路)，系统更加简单且高效。

### Leader
- 由于 Kafka 集群依赖 zookeeper 集群，所以最简单最直观的方案是，
  - 所有 Follower 都 在 ZooKeeper 上设置一个 Watch，一旦 Leader 宕机，
  - 其对应的 ephemeral znode 会 自动删除，此时所有 Follower 都尝试创建该节点，
  - 而创建成功者(ZooKeeper 保证只 有一个能创建成功)即是新的 Leader，
  - 其它 Replica 即为Follower。
- 前面的方案有以下缺点:
  - split-brain (脑裂): 
    - 这是由 ZooKeeper 的特性引起的，虽然 ZooKeeper 能保证所有 Watch按顺序触发，
      但并不能保证同一时刻所有 Replica “看”到的状态是一样的，
      这就可能造成不同 Replica 的响应不一致 。
  - herd effect (羊群效应): 
    - 如果宕机的那个 Broker 上的 Partition 比较多，会造成多个 Watch被触发，
      造成集群内大量的调整。
  - ZooKeeper 负载过重 : 
    - 每个 Replica 都要为此在 ZooKeeper 上注册一个 Watch，当集群规模增加到
      几千个 Partition 时 ZooKeeper 负载会过重。

### Controller
- Kafka 的 Leader Election 方案解决了上述问题，
  - 它在所有 Broker 中选出一个 controller， 所有 Partition 的 Leader 选举都由 Controller 决定。
    - Controller 会将Leader 的改变直接通过 RPC 的方式
      (比 ZooKeeper Queue 的方式更高效)
      通知需为此作为响 应的 Broker。
- Kafka 集群 controller 的选举过程如下 :
  - 每个 Broker 都会在 Controller Path (/controller)上注册一个 Watch。
  - 当前 Controller 失败时，对应的 Controller Path 会自动消失
    (因为它是 ephemeral Node)，
  - 此时该 Watch 被 fire，所有“活”着的 Broker 都会去竞选成为新的 Controller
    (创建新的Controller Path)，
  - 但是只会有一个竞选成功(这点由 Zookeeper 保证)。
  - 竞选成功者即为新的 Leader，竞选失败者则重新在新的 Controller Path 上注册 Watch。
    因为 Zookeeper 的 Watch 是一次性的，被 fire 一次之后即失效，所以需要重新注册。
- Kafka partition Leader 的选举过程如下 (由 Controller 执行):
  - 从 Zookeeper 中读取当前分区的所有 ISR(in-sync replicas) 集合。
  - 调用配置的分区选择算法选择分区的 Leader。

- Kafka 集群 Partition Replication 默认自动分配。
  - 在 Kafka 集群中，每个 Broker 都有均等分配Partition 的 Leader 机会。
  - 创建1个 Topic 包含4个 Partition，2 Replication
    - 上述图 Broker Partition 中，箭头指向为副本，以 Partition-0 为例:Broker1 中 parition-0 为 Leader，Broker2 中 Partition-0 为副本。
  - 当集群中新增2节点，Partition 增加到6个
    - 上述图种每个 Broker (按照 BrokerId 有序)依次分配 主 Partition,下一个 Broker 为副本，如此循环迭代 分配，多副本都遵循此规则。
- 副本分配算法如下:
  - 将所有 N Broker 和待分配的 i 个 Partition 排序。 
  - 将第 i 个 Partition 分配到第(i mod n)个 Broker 上。
  - 将第 i 个 Partition 的第 j 个副本分配到第((i + j) mod n)个 Broker 上。

### Leader
- 和大部分分布式系统一样，Kafka 处理失败需要明确定义一个 Broker 是否“活着”。
- 对于 Kafka 而言，Kafka 存活包含两个条件:
  - 副本所在节点需要与 ZooKeeper 维持 session 
    (这个通过 ZK 的 Heartbeat 机制来实现)。
  - 从副本的最后一条消息的 offset 需要与主副本的最后一条消息
    offset 差值不超过设定阈值 (replica.lag.max.messages)
    或者副本的 LEO 落后于主副本的 LEO 时⻓
    不大于设定阈值 (replica.lag.time.max.ms)，
    - 官方推荐使用后者判断，并在新版本 kafka0.10.0 
      移除了 replica.lag.max.messages 参数。
#### ISR      
- Leader 会跟踪与其保持同步的 Replica 列表，
  该列表称为 ISR(即in-sync Replica)。
- 如果一个 Follower 宕机，或者落后太多，Leader 将把它从 ISR 中移除。
  当其再次满足以上条件之后又会 被重新加入集合中。
- ISR 的引入主要是解决同步副本与异步复制两种方案各自的缺陷:
  - 同步副本中如果有个副本宕机或者超时就会拖慢该副本组的整体性能。
  - 如果仅仅使用异步副本，当所有的副本消息均远落后于主副本时，一旦主副本宕机重新选举， 
    那么就会存在消息丢失情况。

### Replicated log
- Replicated log 是分布式日志系统，主要保证:
  - commit log 不会丢失
  - commit log 在不同机器上是一致的 
- 几个常⻅的基于主从复制的 replicated log 实现:
  - raft:基于多数节点的 ack，节点一般称为 leader/follower， kafka 将要使用
  - pacificA:基于所有节点的 ack，节点一般称为 primary/secondary，kafka 正在使用 
  - bookkeeper:基于法定个数节点的 ack，节点一般称为 writer/bookie，pulsar 正在使用
- Kafka 在 Zookeeper 中动态维护了一个 ISR(in-sync replicas)，
  这个 ISR 里的所有 Replica都 跟上了 leader，只有 ISR 里的成员才有被选为 Leader 的可能。
  - 在这种模式下，对于 f+1 个 Replica，一个 Partition 能在保证
    不丢失已经 commit的消息的前提下容忍 f 个 Replica 的失败。 
  - 在大多数使用场景中，这种模式是非常有利的。
  - 事实上，为了容忍 f 个 Replica 的失败，
    Majority Vote 和 ISR 在 commit 前需要等待的 Replica 数量是一样的，
    但是 ISR 需要的总的Replica 的个 数几乎是 Majority Vote 的一半。

### High Watermark & Log End Offset
- 每个Kafka副本对象都有两个重要属性
  - LEO
  - HW
  - 注意是所有副本，不只是leader
- LEO
  - 日志末端位移，下一条消息的位移值，注意是下一条，
    如果LEO=10，表示该副本保存了10个消息，位移值范围[0,9]
- HW
  - 水位值，对于同一个副本，HW不会大于LEO，小于等于HW的消息被认为已复制（replicated）
  
- 初始时 Leader 和 Follower 的 HW 和 LEO 都是0。
  - Leader 中的 remote LEO 指的就是 leader 端保存的 follower LEO，也被初始化成0。
  - 此时，Producer 没有发送任何消息给 Leader，
    而 Follower 已经开始不断地给 Leader 发送 FETCH 请求了，
    但因为没有数据因此什么都不会发生。
  - 值得一提的是，Follower 发送过来的 FETCH 请求因为无数据而暂时会被寄存到
    Leader 端的 purgatory 中，待 500ms(replica.fetch.wait.max.ms参数)
    超时后会强制完成。
    - 倘若在寄存期间 Producer 端发送过 来数据，那么会 Kafka 会自动唤醒该 FETCH 请求，
    让 Leader 继续处理之。

- Follower 发送 FETCH 请求在 Leader 处理完 PRODUCE 请求之后。 
  - Producer 给该 topic 分区发送了一条消息。

- Follower 发送 FETCH 请求在 Leader 处理完 PRODUCE 请求之后。
  - Producer 给该 topic 分区发送了一条消息。
  - 把消息写入写底层 log(同时也就自动地更新了 Leader 的 LEO)。
  - 尝试更新 Leader HW 值(前面 Leader 副本何时更新 HW 值一节中的第三个条件触发)。我们已经假设此时 Follower 尚未发送 FETCH 请求，那么 Leader 端保存的 remote LEO 依然是0，因此 Leader 会比较它自己的 LEO 值和 remote LEO 值，发现最小 值是0，与当前 HW值相同，故不会更新分区 HW 值。

- 所以，PRODUCE 请求处理完成后 Leader 端的 HW 值依然是0，而 LEO 是1，remote LEO 是1。假设此时 Follower 发送了 FETCH 请求。

- 本例中当 Follower 发送 FETCH 请求时，Leader 端的处理依次是:
  - 读取底层 log 数据。
  - 更新 remote LEO = 0(为什么是0? 因为此时 Follower 还没有写入这条消息。Leader 如何确 认 Follower 还未写入呢?这是通过 Follower 发来的 FETCH 请求中的 fetch offset 来确定 的)。
  - 尝试更新分区 HW —— 此时 Leader LEO = 1，remote LEO = 0，故分区 HW 值= min(leader LEO, follower remote LEO) = 0。
  - 把数据和当前分区 HW 值(依然是0)发送给 Follower 副本。 
- 而 Follower 副本接收到 FETCH response 后依次执行下列操作:
  - 写入本地 log(同时更新 Follower LEO)。
  - 更新 Follower HW —— 比较本地 LEO 和当前 Leader HW 取小者，故 Follower HW = 0。

- 此时，第一轮 FETCH RPC 结束，我们会发现虽然 Leader 和 Follower 都已经在 log中 保存了这条消息，但分区 HW 值尚未被更新。实际上，它是在第二轮 FETCH RPC中 被更新的。

- Follower 发来了第二轮 FETCH 请求，Leader 端接收到后仍然会依次执行下列操作:
  - 读取底层 log 数据。
  - 更新 remote LEO = 1(这次为什么是1了? 因为这轮 FETCH RPC 携带的 fetch offset 是1，
    那么为什么这轮携带的就是1了呢，因为上一轮结束后 follower LEO 被更新为1了)。
  - 尝试更新分区 HW —— 此时 Leader LEO = 1，remote LEO = 1，故分区 HW 值= min(leader
LEO, follower remote LEO) = 1。
  - 把数据(实际上没有数据)和当前分区 HW 值(已更新为1)发送给 Follower 副本。
    
- 同样地，Follower 副本接收到 FETCH response 后依次执行下列操作:
  - 写入本地 log，当然没东⻄可写，故 Follower LEO 也不会变化，依然是1。
  - 更新 Follower HW —— 比较本地 LEO 和当前 Leader HW 取小者。由于此时两者都是1，故更新 Follower HW = 1 。
    
- Producer 端发送消息后 Broker 端完整的处理流程就讲完了。
  - 此时消息已经成功地被复制到 Leader 和 Follower 的 log 中且分区 HW 是1，
    表明 consumer 能够消费 offset = 0 的这条消息。 
    
- 下面我们来分析下 PRODUCE 和 FETCH 请求交互的第二种情况。
  - 第二种情况:FETCH 请求保存在 purgatory 中 PRODUCE 请求到来。
  - 这种情况实际上和第一种情况差不多。
    前面说过了，当 Leader 无法立即满足 FECTH 返回要求的时候(比如没有数据)，
    那么该 FETCH 请求会被暂存到 Leader 端的 purgatory 中，
    待时机成熟时会尝试再次处理它。
  - 不过 Kafka 不会无限期地将其缓存着，默认有个超时时间(500ms)，一旦超时时间已过，
    则这个请求会被强制完成。
    
- 不过我们要讨论的场景是在寄存期间，producer 发送 PRODUCE 请求从而使之满足了条件从而被唤醒。
  此时，Leader端处理流程如下:
  - Leader 写入本地 log(同时自动更新 leader LEO)。 
  - 尝试唤醒在 purgatory 中寄存的 FETCH 请求。 
  - 尝试更新分区 HW。

- LEO 总结
  - follower LEO 何时更新
    - 每新写入一条消息，LEO更新
  - leader remote LEO何时更新
    - 在leader处理follower FETCH请求时
  - leader LEO何时更新
    - leader写日志成功，自动更新自己的LEO
- HW总结
  - follower 副本何时更新HW
    - follower更新LEO之后。比较当前LEO和FETCH响应中leader的HW，取小者
  - leader副本何时更新HW
    - 选所有满足条件的副本，选择最小的LEO做为HW
      - producer向leader写入消息，leader更新LEO后，会检查HW是否也需要更新
      - leader处理follower的FETCH请求，会从log读数据，会尝试更新HW值。

### LEO下的 数据丢失场景
- 初始情况为主副本 A 已经写入了两条消息，对应 HW=1， LEO=2，LEOB=1，
  从副本 B 写入了一条消息，对应 HW=1，LEO=1。
- 此时从副本 B 向主副本 A 发起 fetchOffset=1 请求，
  主副本 收到请求之后更新 LEOB=1，表示副本 B 已经收到了消息 0，
  然后尝试更新 HW 值，min(LEO,LEOB)=1，即不需要 更新，
  然后将消息1以及当前分区 HW=1 返回给从副本 B， 
- 从副本 B 收到响应之后写入日志并更新 LEO=2，然后更新 其 HW=1，
  虽然已经写入了两条消息，但是 HW 值需要在 下一轮的请求才会更新为2。
- 此时从副本 B 重启，重启之后会根据 HW 值进行日志截断，即消息1会被删除。
- 从副本 B 向主副本 A 发送 fetchOffset=1 请求，如果此时主 副本 A 没有什么异常，
  则跟第二步骤一样没有什么问题， 假设此时主副本也宕机了，那么从副本 B 会变成主副本。
- 当副本 A 恢复之后会变成从副本并根据 HW 值进行日志截断，
  即把消息1丢失，此时消息1就永久丢失了。


### LEO下的 数据不一致场景
- 初始状态为主副本 A 已经写入了两条消息对应 HW=1，LEO=2，LEOB=1，
  从副本 B 也同步了两 条消息，对应HW=1，LEO=2。
- 此时从副本 B 向主副本发送 fetchOffset=2 请求， 
  主副本 A 在收到请求后更新分区 HW=2 并将该值 返回给从副本 B，
  如果此时从副本 B 宕机则会导致 HW 值写入失败。
- 我们假设此时主副本 A 也宕机了，从副本 B 先恢复 并成为主副本，
  此时会发生日志截断，只保留消息 0，然后对外提供服务，
  假设外部写入了一个消息1 (这个消息与之前的消息1不一样，用不同的颜色 标识不同消息)。
- 等副本 A 起来之后会变成从副本，不会发生日志截断，
  因为 HW=2，但是对应位移1的消息其实是不一致的。

## Leader epoch
- HW 值被用于衡量副本备份成功与否以及出现失败情况时候的日志截断依据
  可能会导致数据丢失与数据不一致情况，
  因此在新版的 Kafka(0.11.0.0)引入了 leader epoch 概 念。
- leader epoch 表示一个键值对<epoch, offset>，
  其中 epoch 表示 leader 主副本的版本号，从0开始编码，
  当 leader 每变更一次就会+1，offset 表示该 epoch 版本的主副本写入第一条消息的位置。
- 比如<0,0>表示第一个主副本从位移0开始写入消息，
  <1,100>表示第二个主副本版本号 为1并从位移100开始写入消息，
  主副本会将该信息保存在缓存中并定期写入到 checkpoint 文件中，
  每次发生主副本切换都会去从缓存中查询该信息。

### Leader epoch下的数据丢失场景
- 如图所示，当从副本 B 重启之后向主副本 A 发送 offsetsForLeaderEpochRequest，
  epoch 主从副本相等， 则 A 返回当前的 LEO=2，从副本 B 中没有任何大于2的位 移，
  因此不需要截断。
- 当从副本 B 向主副本 A 发送 fetchoffset=2 请求时，A宕机，所以从副本 B 成为主副本，
  并更新 epoch 值为 <epoch=1, offset=2>，HW 值更新为2。
- 当 A 恢复之后成为从副本，并向 B 发送fetcheOffset=2 请 求，
  B 返回 HW=2，则从副本 A 更新 HW=2。
- 主副本 B 接受外界的写请求，从副本 A 向主副本 A 不断发 起数据同步请求。 
- 从上可以看出引入 leader epoch 值之后避免了前面提到的数 据丢失情况，
  
- 但是这里需要注意的是如果在上面的第一步， 
  - 从副本 B 起来之后向主副本 A 发送 offsetsForLeaderEpochRequest 请求失败，
  - 即主副本 A同时 也宕机了，那么消息1就会丢失，具体可⻅下面数据不一致场 景中有提到。

### Leader epach下的 数据不一致场景
- 从副本 B 恢复之后向主副本 A 发送 offsetsForLeaderEpochRequest 请求，
  由于主副 本也宕机了，因此副本 B 将变成主副本并将消息1 截断，此时接受到新消息1的写入。
- 副本 A 恢复之后变成从副本并向主副本 A 发送 offsetsForLeaderEpochRequest 请求，
  请求的 epoch 值小于主副本 B，因此主副本 B 会返回 epoch=1 时的开始位移，
  即 lastoffset=1，因此从 副本 A 会截断消息1。
- 从副本 A 从主副本 B 拉取消息，并更新 epoch 值 <epoch=1, offset=1>。
- 可以看出 epoch 的引入避免的数据不一致，但是两个副本均宕机，则还是存在数据丢失的场景。

### Producer `required.acks`
- 对于某些不太重要的数据，对数据的可靠性要求不是很高，能够容忍数据的少量丢失，
  所以没必要等 ISR 中的 Follower 全部接受成功。
  - 只有被 ISR 中所有 Replica 同步 的消息才被 Commit，
  - 但Producer 发布数据时，Leader 并不需要 ISR 中的所有 Replica 
    同步该数据才确认收到数据。
- 0
  - Producer 不等待 Broker 的 ACK，这提供了最低延迟， 
    Broker 一收到数据还没有写入磁盘就已经返回，
  - 当 Broker 故障时有可能丢失数据。
- 1
  - Producer 等待 Broker 的 ACK，Partition 的 Leader 落盘 成功后返回 ACK，
  - 如果在 Follower 同步成功之前 Leader 故 障，那么将会丢失数据。
- -1(all)
  - Producer 等待 Broker 的 ACK，Partition 的 Leader 和 Follower
    全部落盘成功后才返回 ACK。
  - 但是在 Broker 发送 ACK 时，Leader 发生故障，则会造成数据重 复。

### `request.required.acks=-1` && `min.insync.replicas`
- 如果要提高数据的可靠性，在设置request.required.acks=-1的同时，
  也要 min.insync.replicas 这个参数(可以在 Broker 或者 Topic 层面进行设置)的配合，
  这样才能发挥最大的功效。
- min.insync.replicas这个参数设定 ISR 中的最小副本数是多少，默认值为1，
  - 当且仅当 request.required.acks 参数设置为-1时，此参数才生效。
- 如果 ISR 中的副本数少于 min.insync.replicas 配置的数量时，
  - 客户端会返回异常: `org.apache.kafka.common.errors.NotEnoughReplicasExceptoin: Messages are rejected since there are fewer in-sync replicas than required`

### `request.required.acks=1`的数据丢失
- Producer 发送数据到 Leader，Leader 写本地日志成功，返回客户端成功;
   - 此时 ISR 中的副本还没有来得及拉取该消息，Leader 就宕机了，那么此次发送的消息就会丢失。

### `request.required.acks=-1`
- 同步(Kafka 默认为同步，即 producer.type=sync)的发送模式，
  replication.factor>=2 且 min.insync.replicas>=2 的情况下，不会丢失数据。
- 有两种典型情况。
  - acks=-1 的情况下，数据发送到 Leader, ISR 的 Follower 全部完成数据同步后，
    Leader此时挂掉，那么会选举出新的 Leader，数据不会丢失。
  - acks=-1 的情况下，数据发送到 Leader 后 ，部分 ISR 的副本同步，Leader 此时挂 掉。 
    比如 follower1 和 follower2 都有可能变成新的 Leader, 
    Producer 端会得到返回异常，Producer 端会重新发送数据，数据可能会重复。

### 关于 HW 的进一步探讨
- 考虑上图(即 acks=-1，部分 ISR 副本同步)中的另一种情况，
  - 如果在 Leader 挂 掉的时候，follower1 同步了消息4,5， follower2 同步了消息4，
    与此同时 follower2 被选举为leader，
    那么此时 follower1 中的多出的消息5该做如何处理呢?
  - 这里就需要 HW 的协同配合。
    - 如前所述，一 个 partition 中的 ISR 列表中，
      Leader 的 HW 是所有 ISR 列表里副本中最小的那个的 LEO。 
    - 类似于木桶原理，水位取决于最低那块短板。
- 当 ISR 中的各副本的 LEO 不一致时，如果此时 leader 挂掉，
  选举新的 leader 时并不是按照 LEO 的 高低进行选举，
  而是按照 ISR 中的顺序选举。

## Kafka 高性能 
### 架构层面:
- Partition 级别并行
   - Broker
   - Disk
   - Consumer 端 
- ISR
### IO 层面:
- Batch 读写
- 磁盘顺序 IO
- Page Cache
- Zero Copy
- 压缩

## References
- https://mp.weixin.qq.com/s/fX26tCdYSMgwM54_2CpVrw 
- https://www.jianshu.com/p/bde902c57e80
- https://mp.weixin.qq.com/s?__biz=MzUxODkzNTQ3Nw==&mid=2247486202&idx=1&sn=23f249d3796eb53aff9cf41de6a41761  
- https://zhuanlan.zhihu.com/p/27551928
- https://zhuanlan.zhihu.com/p/27587872
- https://zhuanlan.zhihu.com/p/31322316
- https://zhuanlan.zhihu.com/p/31322697 
- https://zhuanlan.zhihu.com/p/31322840 
- https://zhuanlan.zhihu.com/p/31322994 
- https://mp.weixin.qq.com/s/X301soSDWRfOemQhk9AuPw
- https://www.cnblogs.com/wxd0108/p/6519973.html 
- https://tech.meituan.com/2015/01/13/kafka-fs-design-theory.html 
- https://mp.weixin.qq.com/s/fX26tCdYSMgwM54_2CpVrw 
- https://mp.weixin.qq.com/s/TUFNictt8XXLmmyWlfnj4g 
- https://mp.weixin.qq.com/s/EY6-rA5DJr28-dyTh5BP8w
- https://mp.weixin.qq.com/s/ByIqEgKIdQ2CRsq4_rTPmA
- https://zhuanlan.zhihu.com/p/77677075?utm_source=wechat_timeline&utm_medium=social&utm_oi=670706646783889408&from=timeline
- https://mp.weixin.qq.com/s/LRM8GWFQbxQnKoq6HgCcwQ
- https://www.slidestalk.com/FlinkChina/ApacheKafka_in_Meituan 
- https://tech.meituan.com/2021/01/14/kafka-ssd.html
- https://www.infoq.cn/article/eq3ecYUJSGgWVDGqg5oE?utm_source=related_read_bottom&utm_medium=article
- https://mp.weixin.qq.com/s/Zz35bvw7Sjdn3c8B12y8Mw 
- https://tool.lu/deck/pw/detail?slide=20 
- https://www.jiqizhixin.com/articles/2019-07-23-11
- https://www.jianshu.com/p/c987b5e055b0 
- https://blog.csdn.net/u013256816/article/details/71091774 
- https://zhuanlan.zhihu.com/p/107705346 
- https://www.cnblogs.com/huxi2b/p/7453543.html 
- https://blog.csdn.net/qq_27384769/article/details/80115392 
- https://blog.csdn.net/u013256816/article/details/80865540
- https://tech.meituan.com/2021/01/14/kafka-ssd.html
- https://www.infoq.cn/article/eq3ecYUJSGgWVDGqg5oE?utm_source=related_read_bottom&utm_medium=article
- https://mp.weixin.qq.com/s/Zz35bvw7Sjdn3c8B12y8Mw 
- https://tool.lu/deck/pw/detail?slide=20 
- https://www.jiqizhixin.com/articles/2019-07-23-11 
- https://mp.weixin.qq.com/s/LRM8GWFQbxQnKoq6HgCcwQ 
- https://mp.weixin.qq.com/s/EY6-rA5DJr28-dyTh5BP8w
