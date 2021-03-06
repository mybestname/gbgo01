#  第一周：微服务（微服务概览与治理）

- [直播地址](http://gk.link/a/10pMe)
- [第一周答疑文档](https://shimo.im/docs/prWppTp6qKPgrHwY/)


## 1.微服务概览

- 小即是美：代码少，bug少，易测试，易维护
- 单一职责：一个服务一件事
- 尽早协议：尽早提供API建立契约。
- 重可移植：效率 vs.可移植，首考兼容和移植。

### 要点
- 类依赖 -> 消息交互
- 不同语言解决不同问题（服务）：网关接口（go） vs. 计算密集（c++）
- 消息协议通用（不需要分别打造协议栈）RPC协议
- 原子服务
- 独立进程
- 隔离部署
- 去中心化服务治理

## 2.微服务设计 

- 网关层 (API Gateway)
  - BFF层 (Backend For Frontend)
- 服务层 (Microservice)

### 知识点
- 扇出 fan-out（一个请求上出很多个请求，并行地访问底层的RPC服务）

### 网关
- 避免客户端直接和内部直接通信，暴露的结果在于无法升级。
- 避免客户端来聚合数据。客户端会太重，无法快速迭代和交付。
- 避免客户端进行协议兼容。
- 避免内部去兼容客户端（第3点的反向）。导致面向外部的兼容耦合入业务逻辑。
- 避免内部去兼容客户端版本。导致内部升级困难（这个第1点已经说过）
- 某些统一性要求永远无法收敛。例如安全、限流。（因为对外暴露无法收敛）

#### BFF层
- Backend for frontend， 面向前端应用的后端服务
- BFF 可以认为是一种适配服务，将后端的微服务进行适配(主要包括聚合裁剪和格式适配等逻辑)，向无线端设备暴露友好和统一的API，方便无线设备接入访问后端服务。
- app-interface (非面向资源接口，面向资源接口适用于内网，而不适用于外网）
- data-set join
- 从内外网的视角上去看，BFF层是内网和外网直接的桥梁。
- BFF把内网细粒度的微服务，聚合为面向前端业务留的粗粒度的接口层。

##### BFF的好处

- 轻量（协议可以精简）
- 差异（差异化功能
- 动态升级 （对外协议不变，内部服务升级）
- 组织效率提升 
  - 专门为前端服务的独立团队，而不是直接推到后端，沟通效率提升。或前端+网关组成一个团队提高前后端沟通效率。
  - 康威定律 (Conway's Law, 1967) 设计系统的架构受制于产生这些设计的组织的沟通结构。系统设计本质上反映了企业的组织机构。系统各个模块间的接口也反映了企业各个部门之间的信息流动和合作方式。
  
#### 网关层的问题  
- 单点失败！
- 业务复杂度上升-> BFF层分块儿以减少沟通成本-> 横切面逻辑（安全、日志、限流。。）复杂性！

#### 解决
- 将横切面逻辑（Cross-Cutting）从BFF中抽离（例如路由、安全、认证、限流）上升到API网关层
- API网关对业务解耦（通过路由，哪个接口应该路由到哪个BFF，认证，统一的鉴权）
  - 例如鉴权要升级，那么和业务无关。
  - API网关做为基础设施。
  - 限流功能在API网关层做就可以了（以前在BFF层）
  - 网关层用（https://github.com/envoyproxy/envoy）实现
- BFF层为业务功能


### 微服务层

#### 微服务的划分
- 业务职能（类似公司部门）
- Domain-driven design（DDD）/ Bounded Context （限界上下文）
  - 划分业务边界 
- CQRS (Command Query Responsibility Segregation) M.Fowler
  - 将应用分为命令端和查询端。
    - 命令端处理程序创建，更新和删除请求，并在数据更改时发出事件。
    - 查询端通过针对一个或多个物化视图执行查询来处理查询，这些物化视图通过订阅数据更改时发出的事件流而保持最新。
  - 稿件服务演进
    - 创作稿件、审核稿件、最终发布有大量的逻辑揉和，其中稿件本身的状态多种，最终前台用户只关注稿件能否查看
    - 两个数据库（数据修改+数据结果）+ kafka。数据消费者访问结果库。数据更新入修改库。
      中间用kafka传递更新消息。
    - Polling publisher (早期：查询Mysql，最近变更时间的数据，获取到某处，然后再投入消息队列) 
      - 需要频繁polling数据库，非流式，非实时。   
    - Transaction log tailing (Tx log指Mysql的binlog，包含最终结果库的CRUD)
      - https://microservices.io/patterns/data/transaction-log-tailing.html

#### 微服务安全
##### 外部认证安全 
- API网关统一认证拦截，客户端获得 AccessToken （cookie中）OAuth
- API请求中带有Token，API网关给Header中加入userid（身份信息），
- 挂载在gRPC的metaData中/或者放在HTTP的header中
- BFF的请求参数里面不会有userid，（requst参数里面没有userid），而是在header里面，通过API网关注入。
- API网关会清理外网的请求的Header（攻击），在认证成功后注入新的header
- 而BFF的下游的RPC服务则接口参数里面有userid，userid是服务的参数。
- 具体在Go中，BFF从context取得用户身份信息。

##### 内部服务安全 
- 内网的服务也有安全的要求 例如VIP服务（用户提权），提转币（价值）
- 对于内网，关键：知道谁调用。
- RBAC (Role-based access control) 基于角色的权限调用。
-（Authentication - 认证) vs. (Authorization - 授权)
- 内网传输加密，防嗅探。 gRPC支持证书。即知道是谁，也可以通信加密。
- 一般简单的来说基于APP Token 
- 复杂些基于证书。

## 3.gRPC & 服务发现 
### gPRC
- 语言中立
- ProtoBuff（高性能序列化）
- IDL（proto3）代码生成（生产、维护效率）
- HTTP2 性能 （双向，头压缩，单TCP多路复用，服务端推）
  - go的net/rpc就是基于单TCP的多路复用（通过在一个连接上通过标志req/rsp的ID，称为CallID或RPCID）
  - server push
- 传输层弱依赖，方便以后升级支持HTTP3
  - HTTP3 QUIC  user space congestion control over UDP. (谷歌主导的协议)
- 与负载无关，可以正式跑PB，而测试跑JSON（方便）
- 支持流式API，方便某些业务场景。
- 支持阻塞和非阻塞（同步/异步）调用。
- 支持元数据交互（RPC框架必备）。这点比go的net/rpc要强大。（通过HTTP2的header支持）对于元数据的支持，在认证、跟踪、限流等横切关注点上很重要。
- 错误使用标注状态码（HTTP2）
   - gPRC的标准化做的很好（先标准化再性能，不要过早关注性能）
- 支持Health Check
   - 利用gRPC的健康检查可以和k8s的服务发现/节点管理/ 结合。
   - 在平滑发布中很有用。
  
### 服务发现
- 客户端发现
  - 客服端直接找到服务实例，并访问服务。
  - 一个服务实例被启动时，它的网络地址会被写到注册表上；当服务实例终止时，再从注册表中删除； 服务实例的注册表通过心跳机制动态刷新； 客户端使用一个负载均衡算法，去选择一个可用的服务实例，来响应这个请求。
  - 直连，比服务端服务发现少一次网络跳转，Consumer 需要内置特定的服务发现客户端和发现逻辑。
- 服务端发现
  - 客户端不能感知服务实例，请求被代理转发到服务实例。（可能代理中心化失败）
  - 客户端通过负载均衡器向一个服务发送请求，这个负载均衡器会查询服务注册表， 并将请求路由到可用的服务实例上。服务实例在服务注册表上被注册和注销(Consul Template+Nginx，kubernetes+etcd)。
  - Consumer 无需关注服务发现具体细节，只需知道服务的 DNS 域名即可，支持异构语言开发，需要基础设施支撑，多了一次网络跳转，可能有性能损失。
- zookeep vs. eureka (CP vs. AP)
- eureka 介绍
  - 优点：自我保护机制（错也比没有可用的数据强）
  - 缺点：广播式的全量复制，当服务数量很大时候，写放大太大。
    - 2.0解决：数据shading（一致性hash）；读写server分离（大量读server，少量写server）。

## 4.多集群 & 多租户 

- 多集群的情况下一定要注意数据正交问题：
   - 不要切换集群（而是多集群同时在线，这样数据cache基本都会被hit到）
   - 当某集群出问题时候，只要下线该集群的全部节点就可以（而不是切换）。
   - 这样虽然不同集群的cache是各种独立的。但是cache全部都充分预热过了，所以不会因为集群切换造成数据正交问题。
- 多租户
   - [multi-tenancy](https://en.wikipedia.org/wiki/Multitenancy)
   - 通过"染色"的方式，在真实环境中把测试流量进行数据隔离的流量路由（比灰度测试更好）。
     - 实现上是go的context，使用metadata，注入到HTTP header中来区分流量。
   - 多租户架构本质: 跨服务传递请求携带上下文(context)，数据隔离的流量路由方案。利用服务发现注册租户信息，注册成特定的租户。 
  

## 参考资料

- [Chris Richardson【POJOs in Action/Microservices patterns】作者博客](https://microservices.io/index.html)
- [Microservice 微服务的理论模型和现实路径 【CSDN mindwind-_-】](https://blog.csdn.net/mindfloating/article/details/51221780)
- [微服务架构~BFF和网关是如何演化出来的](https://www.cnblogs.com/dadadechengzi/p/9373069.html)
- [微服务中的设计模式](https://www.cnblogs.com/viaiu/archive/2018/11/24/10011376.html)
- [微服务架构的故障隔离及容错处理](https://www.cnblogs.com/lfs2640666960/p/9543096.html)
- [为什么Uber微服务架构使用多租户？](https://mp.weixin.qq.com/s/L6OKJK1ev1FyVDu03CQ0OA)
- [面向资源的设计](https://www.bookstack.cn/read/API-design-guide/API-design-guide-02-面向资源的设计.md)
- [How To Design Great APIs With API-First Design](https://www.programmableweb.com/news/how-to-design-great-apis-api-first-design-and-raml/how-to/2015/07/10)
- [微服务实战（一）：微服务架构的优势与不足 【Nginx blog/Chris Richardson】](http://www.dockone.io/article/394)
- [微服务实战（二）：构建微服务：使用API Gateway](https://www.jianshu.com/p/3c7a0e81451a)
- [微服务实战（三）：深入微服务架构的进程间通信](https://www.jianshu.com/p/6e539caf662d)
- [微服务实战（四）：服务发现的可行方案以及实践案例](https://my.oschina.net/CraneHe/blog/703173)
- [微服务实践（五）：微服务的事件驱动数据管理](https://my.oschina.net/CraneHe/blog/703169)
- [微服务实践（七）：从单体式架构迁移到微服务架构](https://my.oschina.net/CraneHe/blog/703160)
