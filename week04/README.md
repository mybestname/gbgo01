# 第4周：工程化实践

## 工程项目结构
### Google wire 工具
- 依赖注入工具
- 采用静态代码生成方式，编译期实现，而非runtime的reflection

## API 设计

### API repo
使用API repo来统一访问（对公的一个repo，使用IDL定义（语言无关），使用gRPC做sub
    - git 进行版本管理
    - git hook可以进行规范性检查（变成CI的一部分，例如gitlab/github-action之类）
    - 变更检查（危险性操作的检查，git diff，变成CI的一部分）
    - git的OWNERS文件，细粒度的目录级权限控制。
    - github/gitlab的机器人的协助（使用comments控制）

### API的命名
  - 参考谷歌（https://cloud.google.com/apis/design/naming_convention）
    
      | API Name         | Example                            |
      | --------------   | ---------------------------------- |
      | Product Name     | Google Calendar API                |
      | Service Name     | calendar.googleapis.com            |
      | Package Name     | google.calendar.v3                 |
      | Interface Name   | google.calendar.v3.CalendarService |
      | Source Directory | //google/calendar/v3               |
      | API Name         | calendar                           |

### API的错误处理
   - 参考谷歌设计文档（https://cloud.google.com/apis/design/errors) 
   - rpc错误码：[google.rpc.Code](https://github.com/googleapis/googleapis/blob/master/google/rpc/code.proto)
     
      | HTTP | gRPC | Description |
      | ---- | -----| ----------  |
      | 200  | OK                  | No error.
      | 400  | INVALID_ARGUMENT    | Client specified an invalid argument. Check error message and error details for more information.
      | 400  | FAILED_PRECONDITION | Request can not be executed in the current system state, such as deleting a non-empty directory.
      | 400  | OUT_OF_RANGE        | Client specified an invalid range.
      | 401  | UNAUTHENTICATED     | Request not authenticated due to missing, invalid, or expired OAuth token.
      | 403  | PERMISSION_DENIED   | Client does not have sufficient permission. This can happen because the OAuth token does not have the right scopes, the client doesn't have permission, or the API has not been enabled.
      | 404  | NOT_FOUND           | A specified resource is not found.
      | 409  | ABORTED             | Concurrency conflict, such as read-modify-write conflict.
      | 409  | ALREADY_EXISTS      | The resource that a client tried to create already exists.
      | 429  | RESOURCE_EXHAUSTED  | Either out of resource quota or reaching rate limiting. The client should look for google.rpc.QuotaFailure error detail for more information.
      | 499  | CANCELLED           | Request cancelled by the client.
      | 500  | DATA_LOSS           | Unrecoverable data loss or data corruption. The client should report the error to the user.
      | 500  | UNKNOWN             | Unknown server error. Typically a server bug.
      | 500  | INTERNAL            | Internal server error. Typically a server bug.
      | 501  | NOT_IMPLEMENTED     | API method not implemented by the server.
      | 502  | N/A                 | Network error occurred before reaching the server. Typically a network outage or misconfiguration.
      | 503  | UNAVAILABLE         | Service unavailable. Typically the server is down.
      | 504  | DEADLINE_EXCEEDED   | Request deadline exceeded. This will happen only if the caller sets a deadline that is shorter than the method's default deadline (i.e. requested deadline is not enough for the server to process the request) and the request did not finish within the deadline.

      注意：code为大类错误，例如400，表示找不到。具体可以在message中再带具体的小类错误。这样减少code定义的复杂性。通用错误吃到小类错误。
   - 参考kratos (https://github.com/go-kratos/kratos/blob/main/errors/errors.go)
      - 对404的封装 (https://github.com/go-kratos/kratos/blob/main/errors/types.go)
    
      ```golang
      import (
      	"google.golang.org/grpc/codes"
      	"google.golang.org/grpc/status"
      )
      // NotFound new NotFound error that is mapped to a 404 response.
      func NotFound(domain, reason, message string) *Error {
      	return Newf(codes.NotFound, domain, reason, message)
      }
      // Error is describes the cause of the error with structured details.
      // For more details see https://github.com/googleapis/googleapis/blob/master/google/rpc/error_details.proto.
      type Error struct {
      	s *status.Status
      	Domain   string            `json:"domain"`
      	Reason   string            `json:"reason"`
      	Metadata map[string]string `json:"metadata"`
      }
      // New returns an error object for the code, message.
      func New(code codes.Code, domain, reason, message string) *Error {
      	return &Error{
      		s:      status.New(code, message),
      		Domain: domain,
      		Reason: reason,
      	}
      }
      // Newf New(code fmt.Sprintf(format, a...))
      func Newf(code codes.Code, domain, reason, format string, a ...interface{}) *Error {
      	return New(code, domain, reason, fmt.Sprintf(format, a...))
      }
      ```
### pb的FieldMask的功能
- 对数据进行部分更新的功能。
- 只覆盖指定字段，而非覆盖全部字段。

## 配置管理

### 种类
- 环境变量：物理机/容器/Cluster/OS/Mem 。。。的信息，使用者一般是基础库
- 静态配置：资源进行初始化配置信息，如http，db，redis，需要的ip/port/pass。这种配置被认为是静态的（才是安全的，不应该动态修改，容易导致事故），
  变更流程应该等同于一次app重新部署。
- 动态配置：可以在线修改的配置。
- 全局配置：全局配置模版，再局部替换。以避免重复copy文件的问题。

### Options的设计模式
- [Rob Pike - 01/24/04 Self-referential functions and the design of options](https://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html)
- [Deve Cheney - 10/17/04 Functional options for friendly APIs](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis)

```golang

// Option is config option.
type Option func(*options)

// New new a config with options.
func New(opts ...Option) Config {
	options := options{
		...
	}
	for _, o := range opts {
	    o(&options)
	}
	return &config{
	    opts:   options,
	}
}

// Name with service name.
func Name(name string) Option {
    return func(o *options) { o.name = name }
}

// Version with service version.
func Version(version string) Option {
    return func(o *options) { o.version = version }
}
```

改之后还要还原的例子，可以让option的函数指针同时返回option本身。

```golang
type option func(f *Foo) option

// Option sets the options specified.
// It returns an option to restore the last arg's previous value.
func (f *Foo) Option(opts ...option) (previous option) {
    for _, opt := range opts {
        previous = opt(f)
    }
    return previous
}

// Verbosity sets Foo's verbosity level to v.
func Verbosity(v int) option {
    return func(f *Foo) option {
        previous := f.verbosity
        f.verbosity = v
        return Verbosity(previous)
    }
}

prevVerbosity := foo.Option(pkg.Verbosity(3))
foo.DoSomeDebugging()
foo.Option(prevVerbosity)

func DoSomethingVerbosely(foo *Foo, verbosity int) {
    // Could combine the next two lines,
    // with some loss of readability.
    prev := foo.Option(pkg.Verbosity(verbosity))
    defer foo.Option(prev)
    // ... do some stuff with foo under high verbosity.
}
```
gRPC的callOptions的例子 
(https://github.com/grpc/grpc-go/blob/master/examples/features/interceptor/client/main.go)

```golang

// CallOption configures a Call before it starts or extracts information from
// a Call after it completes.
type CallOption interface {
    // before is called before the call is sent to any server.  If before
    // returns a non-nil error, the RPC fails with that error.
    before(*callInfo) error
    
    // after is called after the call has completed.  after cannot return an
    // error, so any failures should be reported via output parameters.
    after(*callInfo, *csAttempt)
}

// EmptyCallOption does not alter the Call configuration.
// It can be embedded in another structure to carry satellite data for use
// by interceptors.
type EmptyCallOption struct{}

func (EmptyCallOption) before(*callInfo) error      { return nil }
func (EmptyCallOption) after(*callInfo, *csAttempt) {}


// unaryInterceptor is an example unary interceptor.
func unaryInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	var credsConfigured bool
	for _, o := range opts {
		_, ok := o.(grpc.PerRPCCredsCallOption)
		if ok {
			credsConfigured = true
			break
		}
	}
	if !credsConfigured {
		opts = append(opts, grpc.PerRPCCredentials(oauth.NewOauthAccess(&oauth2.Token{
			AccessToken: fallbackToken,
		})))
	}
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	end := time.Now()
	logger("RPC: %s, start time: %s, end time: %s, err: %v", method, start.Format("Basic"), end.Format(time.RFC3339), err)
	return err
}
```

kratos v1的例子
https://github.com/go-kratos/kratos/blob/v1.0.x/pkg/net/rpc/warden/client.go#L98

```golang
// Client is the framework's client side instance, it contains the ctx, opt and interceptors.
// Create an instance of Client, by using NewClient().
type Client struct {
	conf    *ClientConfig
	breaker *breaker.Group
	mutex   sync.RWMutex

	opts     []grpc.DialOption
	handlers []grpc.UnaryClientInterceptor
}

// TimeoutCallOption timeout option.
type TimeoutCallOption struct {
	*grpc.EmptyCallOption
	Timeout time.Duration
}

// WithTimeoutCallOption can override the timeout in ctx and the timeout in the configuration file
func WithTimeoutCallOption(timeout time.Duration) *TimeoutCallOption {
	return &TimeoutCallOption{&grpc.EmptyCallOption{}, timeout}
}

// handle returns a new unary client interceptor for OpenTracing\Logging\LinkTimeout.
func (c *Client) handle() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
```

kratos v2的例子 (https://github.com/go-kratos/kratos/blob/main/transport/grpc/client.go)
```golang
// WithTimeout with client timeout.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.timeout = timeout
	}
}
// clientOptions is gRPC Client
type clientOptions struct {
	endpoint   string
	timeout    time.Duration
	middleware middleware.Middleware
	discovery  registry.Discovery
	grpcOpts   []grpc.DialOption
}

func dial(ctx context.Context, insecure bool, opts ...ClientOption) (*grpc.ClientConn, error) {
    options := clientOptions{
        timeout: 500 * time.Millisecond,
        middleware: middleware.Chain(
            recovery.Recovery(),
        ),
    }
    for _, o := range opts {
        o(&options)
    }
    var grpcOpts = []grpc.DialOption{
        grpc.WithUnaryInterceptor(unaryClientInterceptor(options.middleware, options.timeout)),
    }
    return grpc.DialContext(ctx, options.endpoint, grpcOpts...)
}

func unaryClientInterceptor(m middleware.Middleware, timeout time.Duration) grpc.UnaryClientInterceptor {
    return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	...
}
```

## 包管理
### 历史
- GOPATH 环境变量  (大家都依赖一个目录，没有版本，多个项目的依赖互相影响)
- go1.6 vendor目录 (把单个工程的依赖包copy到项目目录)
- go1.11 `go mod`工具 (基于版本的包管理工具，依赖包的升级更新)
- go1.13 mod模式为默认

### GOPROXY和GOPRIVATE的设置
- GOPROXY可以设置代理
  + 访问内部git server
  + 防止公网仓管变更或消失，导致线上编译失败
  + 安全/审计要求
  + cache，速度+，-带宽消耗
  + import path泄漏（伪造的或不当的import包造成的安全隐患） 
  + goproxy.io是一个开源代理  
- GOPRIVATE可以跳过代理
  + 哪些不走代理的，例如公司内部的git repo
  + 跳过checksum  
- 最好GOPRXOY和PRIVATE结合在公司内部用。（自建goproxyio代理）

## 测试

### 单元测试
- docker的帮助
- subtest + gomock
  - subtest https://golang.org/pkg/testing/#T.Run
    - Run runs f as a subtest of t called name. It runs f in a separate goroutine and blocks until f returns or calls t.Parallel to become a parallel test. Run reports whether f succeeded (or at least did not fail before calling t.Parallel). 
    - Run may be called simultaneously from multiple goroutines, but all such calls must return before the outer test function for t returns.
  - https://github.com/golang/mock
    - https://github.com/golang/mock/tree/master/sample

## 一个建议的Go框架选型（第三周答疑）
- 服务发现：nacos
- 配置中心：nacos或 Applo
- 链路追踪：jadger
- 监控：prometheus
- 日志：ELK
- 容器：k8s
- 框架：kratos
- 缓存：redis-cluster
- 数据库：mysql/mariadb/tidb


## References 

 - https://www.ardanlabs.com/blog/2017/02/package-oriented-design.html
 - https://www.ardanlabs.com/blog/2017/02/design-philosophy-on-packaging.html
 - https://github.com/golang-standards/project-layout
 - https://github.com/golang-standards/project-layout/blob/master/README_zh.md
 - https://www.cnblogs.com/zxf330301/p/6534643.html
 - https://blog.csdn.net/k6T9Q8XKs6iIkZPPIFq/article/details/109192475
 - https://blog.csdn.net/chikuai9995/article/details/100723540
 - https://blog.csdn.net/Taobaojishu/article/details/101444324
 - https://blog.csdn.net/taobaojishu/article/details/106152641
 - https://cloud.google.com/apis/design/errors
 - https://kb.cnblogs.com/page/520743/
 - https://zhuanlan.zhihu.com/p/105466656
 - https://zhuanlan.zhihu.com/p/105648986
 - https://zhuanlan.zhihu.com/p/106634373
 - https://zhuanlan.zhihu.com/p/107347593
 - https://zhuanlan.zhihu.com/p/109048532
 - https://zhuanlan.zhihu.com/p/110252394
 - https://www.jianshu.com/p/dfa427762975
 - https://www.citerus.se/go-ddd/
 - https://www.citerus.se/part-2-domain-driven-design-in-go/
 - https://www.citerus.se/part-3-domain-driven-design-in-go/
 - https://www.jianshu.com/p/dfa427762975
 - https://www.jianshu.com/p/5732b69bd1a1
 - https://www.cnblogs.com/qixuejia/p/10789612.html
 - https://www.cnblogs.com/qixuejia/p/4390086.html
 - https://www.cnblogs.com/qixuejia/p/10789621.html
 - https://zhuanlan.zhihu.com/p/46603988
 - https://github.com/protocolbuffers/protobuf/blob/master/src/google/protobuf/wrappers.proto
 - https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
 - https://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html
 - https://blog.csdn.net/taobaojishu/article/details/106152641
 - https://apisyouwonthate.com/blog/creating-good-api-errors-in-rest-graphql-and-grpc
 - https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html
 - https://www.youtube.com/watch?v=oL6JBUk6tj0
 - https://github.com/zitryss/go-sample
 - https://github.com/danceyoung/paper-code/blob/master/package-oriented-design/packageorienteddesign.md
 - https://medium.com/@eminetto/clean-architecture-using-golang-b63587aa5e3f
 - https://hackernoon.com/golang-clean-archithecture-efd6d7c43047
 - https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1
 - https://medium.com/wtf-dial/wtf-dial-domain-model-9655cd523182
 - https://hackernoon.com/golang-clean-archithecture-efd6d7c43047
 - https://hackernoon.com/trying-clean-architecture-on-golang-2-44d615bf8fdf
 - https://manuel.kiessling.net/2012/09/28/applying-the-clean-architecture-to-go-applications/
 - https://github.com/katzien/go-structure-examples
 - https://www.youtube.com/watch?v=MzTcsI6tn-0
 - https://www.appsdeveloperblog.com/dto-to-entity-and-entity-to-dto-conversion/
 - https://travisjeffery.com/b/2019/11/i-ll-take-pkg-over-internal/
 - https://github.com/google/wire/blob/master/docs/best-practices.md
 - https://github.com/google/wire/blob/master/docs/guide.md
 - https://blog.golang.org/wire
 - https://github.com/google/wire
 - https://www.ardanlabs.com/blog/2019/03/integration-testing-in-go-executing-tests-with-docker.html
 - https://www.ardanlabs.com/blog/2019/10/integration-testing-in-go-set-up-and-writing-tests.html
 - https://blog.golang.org/examples
 - https://blog.golang.org/subtests
 - https://blog.golang.org/cover
 - https://blog.golang.org/module-compatibility
 - https://blog.golang.org/v2-go-modules
 - https://blog.golang.org/publishing-go-modules
 - https://blog.golang.org/module-mirror-launch
 - https://blog.golang.org/migrating-to-go-modules
 - https://blog.golang.org/using-go-modules
 - https://blog.golang.org/modules2019
 - https://blog.codecentric.de/en/2017/08/gomock-tutorial/
 - https://pkg.go.dev/github.com/golang/mock/gomock
 - https://medium.com/better-programming/a-gomock-quick-start-guide-71bee4b3a6f1