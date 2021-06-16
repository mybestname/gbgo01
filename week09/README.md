# 第9周 Go 语言实践- 网络编程

## 目录
- 网络通信协议
- Go 实现网络编程 
- Goim ⻓连接网关 
- ID 分布式生成器 
- IM 私信系统

## 网络通信协议
互联网协议 (Internet Protocol Suite)
- 接口抽象层
  - Socket
- 面向连接(可靠) / 无连接(不可靠)
  - TCP / UDP
- 超文本传输协议
  - HTTP1.1 
  - HTTP2 
  - QUIC(HTTP3)

## Socket 抽象层 
  
- 应用程序通常通过“套接字”向网络发出请求或者应答网络请求。
- 一种通用的面向流的网络接口
- 主要操作:
  - 建立、接受连接 
  - 读写、关闭、超时
  - 获取地址、端口

## TCP 可靠连接，面向连接的协议
- TCP/IP 即传输控制协议/网间协议，是一种面向连接(连接导向)的、可靠的、基于字节流的传输层(Transport layer)通信协议。
  - 三次握手（建）和四次挥手（断）
- 服务端流程:
  - 监听端口 接收客户端请求建立连接 创建 goroutine 处理连接
- 客户端流程:
  - 建立与服务端的连接
  - 进行数据收发
  - 关闭连接


## UDP 不可靠连接，允许广播或多播
- UDP 协议(User Datagram Protocol) 用户数据报协议，一种无连接的传输层协议。
- 一个简单的传输层协议:
  - 不需要建立连接 
  - 不可靠的、没有时序的通信
  - 数据报是有⻓度(65535-20=65515) 
  - 支持多播和广播 低延迟，实时性比较好 
  - 应用于用于视频直播、游戏同步

## HTTP 超文本传输协议
- HTTP(HyperText Transfer Protocol) 详细规定了浏览器和web服务器之间互相通信的规则，通过因特网传送web文档的数据传送协议。
- 请求报文:
```
 Method: HEAD/GET/POST/PUT/DELETE
 Accept:text/html、application/json
 Content-Type:
   application/json
   application/x-www-form-urlencoded
 请求正文 
```
- 响应报文:
```
  状态行(200/400/500) 
  响应头(Response Header) 
  响应正文
```

- 例子
```bash
$ curl -v www.google.com
*   Trying xxx.xxx.xxx.xxx...
* TCP_NODELAY set
* Connected to www.google.com (xxx.xxx.xxx.xxx) port 80 (#0)
> GET / HTTP/1.1
> Host: www.google.com
> User-Agent: curl/version
> Accept: */*
>
< HTTP/1.1 302 Found
< Location: http://www.google.com/url?xxxxxxxxx
< Cache-Control: private
< Content-Type: text/html; charset=UTF-8
< P3P: CP="This is not a P3P policy! See g.co/p3phelp for more info."
< Date: mm-dd-yyyy hh:mm:ss GMT
< Server: gws
< Content-Length: 370
< X-XSS-Protection: 0
< X-Frame-Options: SAMEORIGIN
< Set-Cookie: 1P_JAR=YYYY-MM-DD-HH; expires= mm-dd-yyyy hh:mm:ss GMT; path=/; domain=.google.com; Secure
< Set-Cookie: xxxxxxxxxx; expires=mm-dd-yyyy hh:mm:ss GMT; path=/; domain=.google.com; HttpOnly
<
<HTML><HEAD><meta http-equiv="content-type" content="text/html;charset=utf-8">
<TITLE>...</TITLE></HEAD><BODY>
<H1>...</H1>
...
</BODY></HTML>
* Connection #0 to host www.google.com left intact
* Closing connection 0
```
- 工具
- tcpflow
- nload
- ss
- netstat 
- nmon
- top

## gRPC 基于 HTTP2 协议扩展
- Request
```
Headers
  method = POST
  scheme = https
  path = /api.echo.v1.Echo/SayHello 
  content-type = application/grpc+proto 
  grpc-encoding = gzip
Data
  <Length-Prefixed Message> (带消息长度的序列化的proto消息）
    1 byte of zero (not compressed). （0表示没有压缩，1压缩）
    network order 4 bytes of proto message length. 
    serialized proto message.
```
- Length-Prefixed-Message
   > The repeated sequence of **Length-Prefixed-Message** items is delivered in DATA frames
   >   - **Length-Prefixed-Message** → 
   >      - Compressed-Flag 
   >      - Message-Length 
   >      - Message
   >   - Compressed-Flag → 0 / 1 
   >      - encoded as **1 byte unsigned integer**
   >   - Message-Length → {length of Message} 
   >      - encoded as **4 byte unsigned integer (big endian)**
   >   - Message → *{binary octet}
   > 
   > from https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md
 
- Response
```
Headers
  status = 200
  grpc-encoding = gzip
  content-type = application/grpc+proto
Data
  <Length-Prefixed Message>
Trailers
  grpc-status = 0 
  grpc-message = OK 
  grpc-details-bin = base64(pb)
```
> trailer不是尾，英文单词trail易与tail混淆，trailer在定义上反而恰恰是头，虽然在位置上放在
> 数据的后面。
> trailer是一种response header。trailer原意是拖车，表示挂载，这里解释有误，容易造成误解。
> 因为trailer不是和header对应的，而是像拖车一样“拴在”发送数据的后面，如果非说“尾”的话，
> 也是data的尾，而不是“头和尾”。
>
> - https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Headers/Trailer
>   - Trailer是一个响应首部，允许发送方在分块发送的消息后面添加额外的元信息，这些元信息可能
>     是随着消息主体的发送动态生成的，比如消息的完整性校验，消息的数字签名，或者消息经过处理
>     之后的最终状态等。
>   - The Trailer response header allows the sender to include additional fields at the end of chunked messages in order to supply metadata that might be dynamically generated while the message body is sent, such as a message integrity check, digital signature, or post-processing status.
> - https://blog.cloudflare.com/road-to-grpc/
>   - gRPC uses HTTP trailers for two purposes. To begin with, it sends its 
>      final status (grpc-status) as a trailer header after the content has 
>      been sent. The second reason is to support streaming use cases. These
>      use cases last much longer than normal HTTP requests. The HTTP trailer
>      is used to give the post processing result of the request or the response. 
>      For example if there is an error during streaming data processing, 
>      you can send an error code using the trailer, which is not possible 
>      with the header before the message body.
> - https://datatracker.ietf.org/doc/html/rfc7230#section-4.4 
>   - When a message includes a message body encoded with the chunked
>     transfer coding and the sender desires to send metadata in the form
>     of trailer fields at the end of the message, the sender SHOULD
>     generate a Trailer header field before the message body to indicate
>     which fields will be present in the trailers.  This allows the
>     recipient to prepare for receipt of that metadata before it starts
>     processing the body, which is useful if the message is being streamed
>     and the recipient wishes to confirm an integrity check on the fly.
> - https://datatracker.ietf.org/doc/html/rfc7230#section-4.1.2
>   - A trailer allows the sender to include additional fields at the end
>     of a chunked message in order to supply metadata that might be
>     dynamically generated while the message body is sent, such as a
>     message integrity check, digital signature, or post-processing
>     status.  The trailer fields are identical to header fields, except
>     they are sent in a chunked trailer instead of the message's header
>     section.
> - https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md
>    - For responses end-of-stream is indicated by the presence of the END_STREAM flag on the last received HEADERS frame that carries Trailers.

## HTTP2 如何提升网络速度
### HTTP/1.1 为网络效率做了几点优化
  - keep-alive 增加了持久连接，每个复用连接进行串行请求。 
    > https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Keep-Alive                 
    >  - Connection-specific header fields such as Connection and Keep-Alive
    >    are prohibited in HTTP/2. Chrome and Firefox ignore them in HTTP/2 responses
    >  - Disadvantages 
    >    - https://en.wikipedia.org/wiki/HTTP_persistent_connection#Disadvantages 
    >    - If the client does not close the connection, the resources needed to keep it open and will be unavailable for other clients
    >    - a race condition, Clients must be prepared to retry requests if the connection closes before they receive the entire response （这样就要求服务必须是幂等服务）
  - 浏览器为每个域名最多同时维护6个TCP持久连接。
    > 这只是一个实现参考值， rfc2616(1999)的原文是2个，而chrome使用6个。
    > - https://stackoverflow.com/questions/985431/max-parallel-http-connections-in-a-browser
    > ```shell
    >     Firefox 2:  2
    >     Firefox 3+: 6
    >     Opera 9.26: 4
    >     Opera 12:   6
    >     Safari 3:   4
    >     Safari 5:   6
    >     IE 7:       2
    >     IE 8:       6
    >     IE 10:      8
    >     Edge:       6
    >     Chrome:     6
    > ```
    > - https://chromium.googlesource.com/chromium/chromium/+/trunk/content/browser/loader/resource_scheduler.cc#26  
    > ```
    > static const size_t kMaxNumDelayableRequestsPerHost = 6;
    > ```
    > - https://datatracker.ietf.org/doc/html/rfc2616#page-44
    >    - Clients that use persistent connections SHOULD limit the number of simultaneous connections that they maintain to a given server. A single-user client SHOULD NOT maintain more than 2 connections with any server or proxy. A proxy SHOULD use up to 2*N connections to another server or proxy, where N is the number of simultaneously active users. These guidelines are intended to improve HTTP response times and avoid congestion. 
  - 使用CDN实现域名分片机制。
    > https://stackoverflow.com/questions/56022404/cdn-server-with-http-1-1-vs-webserver-with-http-2 
    > - Under HTTP/1.1 you are limited to 6 connections to a domain. So hosting content on a separate domain (e.g. static.example.com) or loading from a CDN was a way to increase that limit beyond 6. These separate domains are also often cookie-less as they are on separate domains which is good for performance and security. And finally if loading jQuery from code.jquery.com then you might benefit from the user already having downloaded it for another site so save that download completely (though with the number of versions of libraries and CDNs the chance of having a commonly used library already downloaded and in the browser cache is questionable in my opinion).
    > - However separate domains requires setting up a separate connection. Which means a DNS lookup, a TCP connection and usually an HTTPS handshake too. This all takes time and especially if downloading just one asset (e.g. jQuery) then those can often eat up any benefits from having the assets hosted on a separate site! This is in fact why browsers limit the connections to 6 - there was a diminishing rate of return in increasing it beyond that. I've questioned the value of sharded domains for a while because of this and people shouldn't just assume that they will be faster.
    > - HTTP/2 aims to solve the need for separate domains (aka sharded domains) by removing the need for separate connections by allowing multiplexing, thereby effectively removing the limit of 6 "connections", but without the downsides of separate connections. They also allow HTTP header compression, reducing the performance downside to sending large cookies back and forth.
    > - answered by Barry Pollard (http2-in-action author)
  
### HTTP/2 的多路复用 (multiplexing)
> - 目标
>   - 原则：不改动HTTP的语义、方法、状态码、URI、Header字段这些核心概念
>   - 要求：突破性能限制，改进传输性能，实现低延迟和高吞吐量。
> - 解决
>   - 要保证HTTP1.x的各种动词，方法，首部都不受影响
>   - 思路：
>      - 多路复用：将每个文本资源的有效负载切割成更小的部分或“块”，我们可以在线路上混合或“交错”这些块
>      - 在应用层(HTTP2)和传输层(TCP/UDP)之间增加一个层: **二进制分帧层**
>   ![](https://developers.google.com/web/fundamentals/performance/http2/images/binary_framing_layer01.svg)
>      - HTTP1.x的Header被封装到Headers帧，request body则封装到Data帧里面。
>      - HTTP2通信在一个连接上完成，这个连接可承载任意数量的双向数据流。每个数据流以消息的形式发送，而消息由一或多个帧组成，这些帧可以乱序发送，然后再根据每个帧首部的流标识符重新组装。
>   ![](https://developers.google.com/web/fundamentals/performance/http2/images/streams_messages_frames01.svg)
> - 一句话总结：
>   - **将HTTP协议通信分解为二进制编码帧的交换，这些帧对应着特定数据流中的消息。所有这些都在一个TCP连接内复用**
> - 又一个FTSE的经典例子
>   - https://en.wikipedia.org/wiki/Fundamental_theorem_of_software_engineering
> - Reference
>   - https://hpbn.co/http2/#binary-framing-layer
>   - https://halfrost.com/http2-http-frames-definitions/
> 
 
- 通过引入**二进制分帧层**, 实现了HTTP的多路复用。
  - 请求数据二进制分帧层处理之后，会转换成请求 ID 编号的帧，通过协议栈将这些帧发送给服务器。
  - 服务器接收到所有帧之后，会将所有相同 ID 的帧合并为一条完整 的请求信息。
  - 然后服务器处理该条请求，并将处理的响应行、响应头和响应体分 别发送至二进制分帧层。
  - 同样，二进制分帧层会将这些响应数据转换为一个个带有请求 ID 编号的帧，经过协议栈发送给浏览器。
  - 浏览器接收到响应帧之后，会根据 ID 编号将帧的数据提交给对应的请求。

- 类比go的net/rpc包的client实现
  - 发 
    - https://github.com/golang/go/blob/go1.16.5/src/net/rpc/client.go#L83-L90
     ```
          func (client *Client) send(call *Call) {
          ...
          	seq := client.seq
          	client.seq++
          	client.pending[seq] = call  //map中维护编号和call的引用
          	client.mutex.Unlock()
           
          	// Encode and send the request.
          	client.request.Seq = seq
     ```
  - 收
    - https://github.com/golang/go/blob/go1.16.5/src/net/rpc/client.go#L113-L115
    ```go 
    func (client *Client) input() {
    	var err error
    	var response Response
    	for err == nil {
    		response = Response{}
    		err = client.codec.ReadResponseHeader(&response)
    		if err != nil {
    			break
    		}
    		seq := response.Seq        //从resp中拿到编号
    		client.mutex.Lock()
    		call := client.pending[seq]  //通过编号找到call引用。
    		delete(client.pending, seq)
    		client.mutex.Unlock()
    ```
  - 通过类比`net/rpc`包，可以看出，go可以用一个连接发多个rpc，内部通过一个自增
    序列号来引用rpc的请求。 思路上和http2是一样的。


## HTTP 超文本传输协议-演进
### HTTP 发展史
- 1991 HTTP/0.9 
- 1996 HTTP/1.0
- 1997 HTTP/1.1 
   - 最广泛版本
- 2015 HTTP/2.0 
  - 优化HTTP/1.1的性能和安全性 
- 2018 HTTP/3.0 
  - 使用 UDP 取代 TCP 协议

#### HTTP2

- 二进制分帧，按帧方式传输
- 多路复用，代替原来的序列和阻塞机制 • 头部压缩，通过HPACK压缩格式
- 服务器推送，服务端可以主动推送资源

#### HTTP3
- 连接建立延时低，一次往返可建立HTTPS连接
- 改进的拥塞控制，高效的重传确认机制
- 切换网络保持连接，从4G切换到WIFI不用重建连接

> ![](https://images.ctfassets.net/ee3ypdtck0rk/4du4aqnKuOLU4YbbHWYSfv/c88d774f278090e6ef3a5435b46bbfea/Screen_Shot_2021-01-28_at_6.54.47.png) 
> ![](https://camo.githubusercontent.com/21b23f947236de560cb7373659c6c8033c729e3d126303b14046c82ca5401140/68747470733a2f2f757365722d676f6c642d63646e2e786974752e696f2f323031392f31302f31352f313664636437326136353664393034643f773d35373926683d32343526663d706e6726733d3634313932) 
>
> HTTP2对TCP的**队头阻塞**无法彻底解决 (TCP-level head-of-the line blocking)
>  - HOL blocking
>  - TCP为了保证可靠传输，有个特别的“丢包重传”机制，丢失的包必须要等待重新传输确认，HTTP2出现丢包时，整个TCP都要开始等待重传，那么就会阻塞该TCP连接中的所有请求
>  - 这个问题只有替换传输层协议才能解决，所以HTTP3用了UDP(QUIC)
> 
> **Reference**
>  - https://github.com/ljianshu/Blog/issues/57
>  - https://ably.com/topic/http-2-vs-http-3
>  - https://calendar.perfplanet.com/2020/head-of-line-blocking-in-quic-and-http-3-the-details/
>  - https://http3-explained.haxx.se/zh/why-quic/why-tcphol
>

## HTTPS 超文本传输安全协议 

- 常称为HTTP over TLS、HTTP over SSL或HTTP Secure)
- 是一种通过计算机网络进行安全通信的传输协议。

SSL ( Secure Sockets Layer )
  - 网景公司(Netscape)开发
  - 1.0 (?)
  - 2.0 (1995) 2011作废  
  - 3.0 (1996) 2015作废 

TLS (Transport Layer Security)
  - 1.0 (1999) 2021作废  
     - IETF将SSL标准化，即RFC2246
     - 并将其称为TLS。 
  - 1.1 (2006) 2021作废
     - 添加对CBC攻击的保护、支持IANA登记的参数。
  - 1.2 (2008)
    - 增加SHA2
    - 增加AEAD加密算法，如GCM模式
    - 添加 TLS 扩展定义和 AES 密码组合。
  - 1.3 (2018)
    - 较TLS1.2速度更快，性能更好、更加安全。

### SSL/TLS 重要概念 

SSL/TLS 协议提供主要的作用有:

- 认证用户和服务器，确保数据发送到正确的客户端和服务器。
- 加密数据以防止数据中途被窃取。
- 维护数据的完整性，确保数据在传输过程中不被改变。

对称加密:
- 指的就是加、解密使用的同是一串密钥，所以被称做对称加密。对称加密只有一个密钥作为私钥。

非对称加密:
- 指的是加、解密使用不同的密钥，一把作为公开的公钥，另一把作为私钥。公钥加密的信息，只有私钥才能解密。

CA证书:
- CA是负责签发证书、认证证书、管理已颁发证书的机关 
- CA证书通常内置在操作系统，或者浏览器中
- CA用自己的私钥对指纹签名，浏览器通过内置CA跟证书公钥进行解密，如果解密成功就确定证书是CA颁发的。

### TLS 1.2 如何解决安全问题?

要解决的问题:
- 防窃听(eavesdropping)，对应加密(Confidentiality) 
- 防篡改(tampering)，对应完整性校验(Integrity)
- 防伪造(forgery)，对应认证过程(Authentication) 
  
如何保证公钥不被篡改?
- 解决方法:将公钥放在数字证书中。只要证书是可信的，公钥就是可信的。

公钥加密计算量太大，如何减少耗用的时间?
- 解决方法:每一次对话(session)，客户端和服务器端都生成 一个“对话密钥”(session key)，用它来加密信息。由于“对话密钥”是对称加密，所以运算速度非常快，而服务器公钥只用于 加密“对话密钥”本身，这样就减少了加密运算的消耗时间。

因此，SSL/TLS协议的基本过程:
- 客户端向服务器端索要证书，并通过签名验证公钥。
- 双方协商生成“对话密钥”，加密类型、随机串(非对称加密)。 
- 双方采用“对话密钥”进行加密通信(对称加密)。

### TLS 1.3 : Faster & More Secure

TLS 1.3 与之前的协议有较大差异，主要在于:
- RSA密钥交换被废弃，引入新的密钥协商机制PSK。
- 支持0-RTT数据传输，复用PSK无握手时间。
- 废弃若干加密组件，SHA1、MD5等hash算法。
- 不再允许压缩加密报文，不允许重协商，不发ChangeCipher了。

密钥协商机制:
- RSA是常用且简单的一个交换密钥的算法，即客户端决定密钥后， 用服务器的公钥加密传输给对方，这样通信双方就都有了后续通信的密钥。
- Diffie–Hellman(DH)是另一种交换密钥的算法，客户端和服务器都生成一对公私钥，然后将公钥发送给对方，双方得到对方的公钥后，用数字签名确保公钥没有被篡改，然后与自己的私钥结合，就可以计算得出相同的密钥。

为了保证前向安全，TLS 1.3 中移除了 RSA 算法，Diffie–Hellman 是唯一的密钥交换算法。

## Go 网络编程 - 基础概念

### 基础概念:
- Socket:数据传输
- Encoding:内容编码 
- Session:连接会话状态 
- C/S模式:通过客户端实现双端通信 
- B/S模式:通过浏览器即可完成数据的传输

### 简单例子
- TCP/UDP/HTTP
- 网络轮询器
  - 多路复用模型
  - 多路复用模块 
  - 文件描述符 
  - Goroutine 唤醒
    
#### Go 网络编程 - TCP 简单用例

```go
func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:10000")
	if err != nil {
		log.Fatalf("listen error: %v\n", err)
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Printf("accept error: %v\n", err)
			continue
		}
		// 开始goroutine监听连接（每一个conn一个goroutine）
		go handleConn(conn)
	}
}
func handleConn(conn net.Conn) {
	defer conn.Close() // 首先defer资源关闭，go的习惯做法1
	// 读写缓冲区，通过使用bufio的缓冲，减少os的实际syscall对网络设备的读写请求
	rd := bufio.NewReader(conn)
	wr := bufio.NewWriter(conn)
	for {
		line, _, err := rd.ReadLine()
		if err != nil {
			log.Printf("read error: %v\n", err)
			return
		}
		wr.WriteString("hello ")
		wr.Write(line)
		wr.Flush() // 一次性syscall ，这就是使用bufio的好处。
		// 注意：这里说一次是理想情况，具体根据实际情况和底层实现细节的约束来决定具体的syscall调用和次数。不能想当然。
	}
}
```

#### Go 网络编程 - UDP 简单用例
```go
func main() {
	listen, err := net.ListenUDP("udp", &net.UDPAddr{Port: 20000})
	if err != nil {
		log.Fatalf("listen error: %v\n", err)
	}
	defer listen.Close()
	for {
		var buf [1024]byte
		n, addr, err := listen.ReadFromUDP(buf[:])
		if err != nil {
			log.Printf("read udp error: %v\n", err)
			continue
		}
		data := append([]byte("hello "), buf[:n]...)
		listen.WriteToUDP(data, addr)
	}
}
```
#### Go 网络编程 - HTTP 简单用例
```go
// HTTPServer
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		fmt.Fprintf(w, "Welcome to the home page!")
	})
	s := &http.Server{
		Addr:           ":8080",
		Handler:        mux,
		ReadTimeout:    1 * time.Second,
		WriteTimeout:   1 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}

```
```go
// HTTPClient
func main() {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr, Timeout: 1 * time.Second}
	resp, err := client.Get("http://127.0.0.1:8080/")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	fmt.Println(ioutil.ReadAll(resp.Body))
}
```

### Go 网络编程 - I/O模型

Linux下主要的IO模型分为:
- Blocking IO - 阻塞IO
- Nonblocking IO - 非阻塞IO
- IO multiplexing - IO多路复用
- Signal-driven IO - 信号驱动式IO(异步阻塞) 
- Asynchronous IO - 异步IO

> 这种分法源于：UNP
> - https://en.wikipedia.org/wiki/UNIX_Network_Programming
> 

#### 同步
请求导致进程阻塞，直到I/O操作完成。
  - 阻塞: 服务端返回结果之前，调用端会被挂起，暂停运行。 
     - Blocking I/O (`read/write(O_SYNC)`) 
     - 系统I/O调用本身是阻塞的。 
  - 非阻塞: 
     - Nonblocking I/O (`read/write(O_NONBLOCK)`)
        - 不会阻塞调用端，而会立刻返回error，直到准备好。调用端需要不断的重试直到成功。
        - I/O系统调用本身是不阻塞的
        - 用户态的轮询
     - I/O multiplixing：非阻塞的改进模式：更好的轮询方法：系统级别的轮询，而不是用户态的轮询。
        - 阻塞在系统调用select/epoll上面，而不是阻塞在I/O系统调用上。
        - 监视的事情交给了内核，内核负责数据到达的处理。
        - `select` 
        - `epoll`
        - select/epoll通过后，就可以调I/O系统调用。
        - select/epoll的优势并不是对于单个连接能处理得更快，而是在于能处理更多的连接。
   - SIGIO
     - 使用信号的方式（SIGIO信号） 告诉内核要调某I/O系统调用
     - 通过接收SIGIO知道内核已经准备好。
     - 调系统I/O调用。
前四种都是同步模型，因为其中的I/O操作将阻塞进程。不论在第一阶段采取的什么样的策略，第二阶段都要阻塞。

#### 异步
不导致请求进程阻塞
  - Asynchronous I/O  (所谓的 POSIX asynchronous I/O )
  - I/O系统调用本身就是一个异步系统调用（需要操作系统本身支持这种API）
>  - Linux 内核从 2.5 版开始就有异步 I/O，到2.6以上版本才可用，但内核一直支持不够完善，被认为难以使用且效率低下
>  - 例如nginx只是读用，而写不能用。
>  - 在linux 5.1中加入的io_uring的架构，可能对未来linux AIO的使用有改观。
>    - https://lwn.net/Articles/776703/
>    - https://github.com/Linkerist/blog/issues/25  
>    - https://lore.kernel.org/io-uring/20210127212541.88944-1-axboe@kernel.dk/ 
>    - https://static.sched.com/hosted_files/kvmforum2020/9c/KVMForum_2020_io_uring_passthrough_Stefano_Garzarella.pdf
>    - https://cor3ntin.github.io/posts/iouring/
>    - https://github.com/axboe/liburing
> ![](https://cor3ntin.github.io/posts/iouring/uring.svg)
>  - 在go支持uring的可能性
>    - https://github.com/golang/go/issues/31908
>    - https://developers.mattermost.com/blog/hands-on-iouring-go/
>    - https://github.com/hodgesds/iouring-go 
>    - https://github.com/Iceber/iouring-go
>    - http://icebergu.com/archives/go-iouring
>    - https://mp.weixin.qq.com/s/YPiYNPa3xVD9Il1HeB5pTw 

### Go 网络编程 - I/O多路复用
Go 语言在采用 I/O 多路复用模型处理 I/O 操作，但是他没有选择最常⻅的系统调用 
select，例如在 Linux 上使用Epoll。 虽然select 也可以提供I/O 多路复用的能力，
但是使用它有比较多的限制:
- 监听能力有限 — 最多只能监听 1024 个文件描述符;
- 内存拷⻉开销大 — 需要维护一个较大的数据结构存储文件描述符，该结构需要拷⻉到内核中;
- 时间复杂度 𝑂(𝑛) — 返回准备就绪的事件个数后，需要遍历所有的文件描述符;

I/O 多路复用:
- 进程阻塞于 select，等待多个IO中的任一个变为可读，select调用返回，通知相应IO可以读。 
它可以支持单线程响应多个请求这种模式。


### Go 网络编程 - 多路复用模块 

为了提高 I/O 多路复用的性能 不同的操作系统也都实现了自己的 I/O 多路复用函数，
例如:
- epoll 
> epoll is a Linux kernel system call for a scalable I/O event notification mechanism, first introduced in version 2.5.44 of the Linux kernel.
- kqueue 
> Kqueue is a scalable event notification interface introduced in FreeBSD 4.1 on July 2000, also supported in NetBSD, OpenBSD, DragonFly BSD, and macOS.  
- evport 
> Solaris 10 event ports.

> - evport = Solaris 10
> - epoll = Linux
> - kqueue = OS X，FreeBSD
> - select =通常以fallback形式安装在所有平台上
> - Evport，Epoll和KQueue具有 O(1)描述符选择算法复杂度，并且它们都使用内部内核空间内存结构。他们还可以提供很多(成千上万个)文件描述符。
> - 除其他以外，select最多只能为提供最多1024个描述符，并且对描述符进行完全扫描(因此每次迭代所有描述符以选择一个可使用的描述符)，因此复杂度为 O(n)。
> - libevent
>   - https://github.com/libevent/libevent   
>   - Unix-like, Windows, OS X
>   - libevent supports /dev/poll, kqueue(2), POSIX select(2), Windows IOCP, poll(2), epoll(7) and Solaris event ports.
> - libuv
>   - libuv is a multi-platform support library with a focus on asynchronous I/O. It was primarily developed for use by Node.js, but it's also used by Luvit, Julia, pyuv, and others.
>   - https://en.wikipedia.org/wiki/Libuv
>   - https://github.com/libuv/libuv 
>   - https://github.com/libuv/libuv/pull/2322 (io_uring快搞定了)

Go 语言为了提高在不同操作系统上的 I/O 操作性能，使用平台特定的函数实现了多个版本的网络轮询模块

- src/runtime/netpoll_epoll.go
- src/runtime/netpoll_kqueue.go 
- src/runtime/netpoll_solaris.go
- src/runtime/netpoll_windows.go 
- src/runtime/netpoll_aix.go
- src/runtime/netpoll_fake.go

## Goim ⻓连接网关 

### Goim ⻓连接 TCP 编程 - 概览

#### Comet
- ⻓连接层，主要是监控外网 TCP/Websocket端口， 
- 通过设备 ID 进行绑定 Channel 实现，
- 实现了聊天室，适合直播时房间消息的广播。

#### Logic
逻辑层
- 监控连接 Connect、Disconnect 事件
- 可自定义鉴权，进行记录 Session 信息(设备 ID、 ServerID、用户 ID)
- 业务可通过设备 ID、用户 ID、RoomID、全局广播进行消息推送。

#### Job 
- 通过消息队列的进行推送消峰处理，
- 把消息推送到对应 Comet 节点

#### 模块间通信
- 各个模块之间通过 gRPC 进行通信。
  

### Goim ⻓连接 TCP 编程 - 协议设计
```
- Package Length 包⻓度         4 bytes
- Header Length，头⻓度         2 bytes
- Protocol Version，协议版本    2 bytes
- Operation，操作码             4 bytes
    - Auth  
    - Heartbeat
    - Message
- Sequence 请求序号 ID          4 bytes
    - 按请求、响应对应递增 ID
- Body，包内容                  (PackLen-HeaderLen) bytes
```

### Goim ⻓连接 TCP 编程 - 边缘节点
Comet ⻓连接连续节点，通常部署在距离用户比较近，通过 TCP 或者 Websocket 建立连接，并且通过应用层 Heartbeat 进行保活检测，保证连接可用性。

节点之间通过云 VPC 专线通信，按地区部署分布 
- 国内:
  - 华北(北京) 华中(上海、杭州) 华南(广州、深圳) 华⻄(四川)
- 国外:
  - 香港、日本、美国、欧洲

### Goim ⻓连接 TCP 编程 - 负载均衡

- **⻓连接负载均衡比较特殊**，需要按一定的负载算法进行分配节点
  - 可以通过 HTTPDNS 方式，请求获致到对应的节点 IP 列表，
  - 例如，返回固定数量 IP，按一定的权重或者最少连接数进行排序，客户端通过 IP 逐个重试连接
- 流程    
  - Comet 注册 IP 地址，以及节点权重，定时 Renew 当前节点连接数量;
  - Balancer 按地区经纬度计算，按最近地区(经纬度)提供 Comet 节点 IP 列表，以及权重计算排序;
  - BFF 返回对应的⻓连接节点 IP，客户端可以通过 IP直接连;
  - 客户端 按返回 IP 列表顺序，逐个连接尝试建立⻓ 连接

### Goim ⻓连接 TCP 编程 - 心跳保活机制
⻓连接断开的原因:
- ⻓连接所在进程OS被杀死
- NAT 超时
- 网络状态发生变化，如移动网络 & Wifi 切换、断开、重连
- 其他不可抗因素 (网络状态差、DHCP 的租期等等)

高效维持⻓连接方案
- 进程保活(防止进程被杀死) 
- 心跳保活(阻止 NAT 超时) 
- 断线重连(断网以后重新连接网络)

自适应心跳时间
- 心跳可选区间，[min=60s，max=300s] 
- 心跳增加步⻓，step=30s
- 心跳周期探测，success=current + step、fail=current - step

### Goim ⻓连接 TCP 编程 - 用户鉴权和 Session 信息 
用户鉴权，在⻓连接建立成功后，需要先进行连接鉴权，并且绑定对应的会话信息; 
1. Connect，建立连接进行鉴权，保存 Session 信息:
- DeviceID，设备唯一 ID
- Token，用户鉴权 Token，认证得到用户 ID 
- CometID，连接所在 comet 节点

2. Disconnect，断开连接，删除对应 Session 信息:
- DeviceID，设备唯一 ID 
- CometID，连接所在 Comet 节点 
- UserID，用户 ID

3. Session，会话信息通过 Redis 保存连接路由信息:
- 连接维度，通过 设备 ID 找到所在 Comet 节点
- 用户维度，通过 用户 ID 找到对应的连接和 Comet 所在节点

### Goim ⻓连接 TCP 编程 
#### Comet
Comet ⻓连接层，实现连接管理和消息推送:
- Protocol，TCP/Websocket 协议监听;
- Packet，⻓连接消息包，每个包都有固定⻓度;
- Channel，消息管道相当于每个连接抽象，最终TCP/ Websocket 中的封装，进行消息包的读写分发;
- Bucket，连接通过 DeviceID 进行管理，用于读写锁拆散，并且实现房间消息推送，类似 Nginx Worker;
  - bucket是为了拆锁，把数据进行分片，按bucket对整个数据分片，按bucket维度进行管理。
- Room，房间管理通过 RoomID 进行管理，通过链表进行 Channel 遍历推送消息;

每个 Bucket 都有独立的 Goroutine 和读写锁优化: 
```go
    Buckets {
        channels map[string]*Channel
        rooms map[string]*Room 
    }
```  
#### Comet Bucket
维护当前消息通道和房间的信息，有独立的 Goroutine 和 读写锁优化，
用户可以自定义配置对应的 buckets 数量， 
在大并发业务上尤其明显。

#### Comet Room
结构也比较简单，维护了的房间的通道 Channel, 推送消 息进行了合并写，
即 Batch Write, 如果不合并写，每来一 个小的消息都通过⻓连接写出去，
系统 Syscall 调用的开 销会非常大，Pprof 的时候会看到网络 Syscall 是大头。

#### Comet Channel
一个连接通道。Writer/Reader 就是对网络 Conn 的封装， 
cliProto 是一个 Ring Buffer，
保存 Room 广播或是直接发 送过来的消息体。

### Goim ⻓连接 TCP 编程 - 内存优化
内存优化主要几个方面 
一个消息一定只有一块内存:
 - 使用 Job 聚合消息，Comet 指针引用。 
一个用户的内存尽量放到栈上:
 - 内存创建在对应的用户 Goroutine中。 
内存由自己控制:
 - 主要是针对 Comet 模块所做的优化，可以查看模块中各个分配内存的地方，都使用了内存池。

### Goim ⻓连接 TCP 编程 - 模块优化
模块优化也分为以下几个方面 

- 消息分发一定是并行的并且互不干扰:
  - 要保证到每一个 Comet 的通讯通道必须是相互独立的，
  - 保证消息分发必须是完全并列的，并且彼此之间互不干扰。
- 并发数一定是可以进行控制的:
  - 每个需要异步处理开启的 Goroutine(Go 协程)都必须预先 
  - 创建好固定的个数，
     - 如果不提前进行控制，Goroutine 就随时存在爆发的可能。
- 全局锁一定是被打散的:
  - Socket 链接池管理、用户在线数据管理都是多把锁
  - 打散的个数通常取决于CPU，
    - 往往需要考虑 CPU 切换时造成的负担，并非是越多越好。

### Goim ⻓连接 TCP 编程 - Logic 

Logic 业务逻辑层，处理连接鉴权、消息路由，用户会话管理; 主要分为三层:
- sdk，通过 TCP/Websocket 建立⻓连 接，进行重连、心跳保活;
- goim，主要负责连接管理，提供消息⻓ 连能力;
- backend，处理业务逻辑，对推送消息过 虑，以及持久化相关等;

### Goim ⻓连接 TCP 编程 - Job

业务通过对应的推送方式，可以对连接设备、房间、 用户 ID 进行推送，通过 Session 信息定位到所在的 Comet 连接节点，并通过 Job 推送消息;

通过 Kafka 进行推送消峰，保证消息逐步推送成功; 支持的多种推送方式:
- Push(DeviceID, Message) 
- Push(UserID, Message)
- Push(RoomID, Message) 
- Push(Message)

### Goim ⻓连接 TCP 编程 - 唯一 ID 设计 
唯一 ID，需要保证全局唯一，绝对不会出现重复的ID，且 ID 整体趋势递增。

通常情况下，ID 的设计主要有以下几大类:
- UUID
- 基于 Snowflake 的 ID 生成方式 
- 基于申请 DB 步⻓的生成方式 
- 基于数据库多主集群模式
- 基于 Redis 或者 DB 的自增 ID生成方式 
- 特殊的规则生成唯一 ID

### Goim ⻓连接 TCP 编程 - 唯一 ID 设计 Snowflake

Snowflake，is a network service for generating
unique ID numbers at high scale with some simple guarantees.

id is composed of:
- time - 41 bits (millisecond precision w/ a custom epoch gives us 69 years)
- configured machine id - 10 bits - gives us up to 1024 machines
- sequence number - 12 bits - rolls over every 4096 per machine (with protection to avoid rollover in the same ms)

```
 64 bits Id
 
 0 - 00000000 00000000 00000000 00000000 00000000 0 - 00000000 00 - 00000000 0000
 
 - 1. 首位0 不用
 - 2. 41 bit 时间戳 （69年）
 - 3. 10 bit 工作机器id （表示1024个机器）
 - 4. 12 bit 序列号     (每台机器每毫秒可以发出4096个id）
``` 

可以保证粗略有序。

10bit - 数据中心ID + 机器ID，保证不同机器生成唯一


### Goim ⻓连接 TCP 编程 - 唯一 ID 设计 Sonyflake

Sonyflake is a distributed unique ID generator
inspired by Twitter's Snowflake. id is composed of:

```
 00000000 00000000 00000000 0000000 - 00000000 - 00000000 00000000 
- 1. 39 bits for time in units of 10 msec （单位变为10毫秒，可用174年）
- 2. 8 bits for a sequence number  (256个id/每10毫秒)
- 3. 16 bits for a machine id      (65536个机器）
```

As a result, Sonyflake has the following advantages
and disadvantages:
- The lifetime (174 years) is longer than that of Snowflake (69 years)
- It can work in more distributed machines (2^16) than Snowflake (2^10)
- It can generate 2^8 IDs per 10 msec at most in a single machine/thread (slower than Snowflake)

16bit - 机器ID，可用网段地址 255.255，保证不同机器生成唯一

### Goim ⻓连接 TCP 编程 - 唯一 ID 设计 基于步⻓递增
- 基于步⻓递增的分布式ID生成器，可以生成基于递增，并且比较小的唯一ID
- 强调Id连续性的场景

服务主要分为:
- 通过 gRPC 通信，提供 ID 生成接口，并且携带业 务标记，为不同业务分配 ID;
- 部署多个 id-server 服务，通过数据库进行申请 ID 步⻓，并且持久化最大的 ID，
   - 例如，每次批量取 1000到内存中，可减少对 DB 的压力;
- 数据库记录分配的业务 MAX_ID 和对应 Step ，供 Sequence 请求获取;

### Goim ⻓连接 TCP 编程 - 唯一 ID 设计 基于数据库集群模式

基于数据库集群模式，
在 MySQL 中的双主集群模式采用的是这个方案; 

服务主要分为:
- 两个 MySQL 实例都能单独的生产自增ID; 
- 设置 ID 起始值和自增步⻓;
  - MySQL_1 配置:
  ```mysql
  set @@auto_increment_offset = 1; -- 起始值 
  set @@auto_increment_increment = 2; -- 步⻓
  ```
  - MySQL_2 配置:
  ```mysql
  set @@auto_increment_offset = 2; -- 起始值 
  set @@auto_increment_increment = 2; -- 步⻓
  ```

## IM 私信系统

### IM 私信系统 - 基本概念 
- 在聊天系统中，我们几乎每个人都在使用聊天应用，
- 并且对消息及时性要求也非常高; 
- 对消息也需要有一致性保证; 
- 并且都有着丰富的多媒体传输功能:
  - 1 on 1 (1对1)
  - Group chat(群聊)
  - Online presence(在线状态) 
  - Multiple device support(多端同步) 
  - Push notifications(消息通知)
- 客户端可以是 Android、iOS、Web 应用;
- 通常客户端之间不会进行直接通信，而是客户端连接到服务端进行通信;
- 服务端需要支持:
  - 接收各个客户端消息 
  - 消息转发到对应的人
  - 用户不在线，存储新消息 
  - 用户上线，同步所有新消息

### IM 私信系统 - 实时通信协议 
- 在聊天系统中，最重要的是通信协议，如何有保证地及时送达消息;
- 一般来看，移动端基本都是通过⻓连方式实现， 而 Web 端可以使用 HTTP、Websocket 实现实时通信;
常用通信方式:
- TCP
- WebSocket
- HTTP ⻓轮询 
- HTTP 定时轮询

### IM 私信系统 - 服务类型

在聊天系统中，有着很多用户、消息功能，比如:
- 登录、注册、用户信息，可以通过 HTTP API 方式;
- 消息、群聊、用户状态，可以通过 实时通信 方式;
- 可能集群一些三方的服务，比如 小米、华为推送、 APNs等;

主要服务可为三大类:
- 无状态服务
- 有状态服务
- 第三方集成

### IM 私信系统 - 模块功能 
- 在聊天系统中，Goim 主要⻆色是 Real time service，
- 实现对 `连接` 和 `状态` 的管理: 
- 可以通过 API servers 进行系统之间的解耦;
- 各个服务的主要功能为:
  - 聊天服务，进行消息的 发送 和 接收 
  - 在线状态服务，管理用户 在线 和 离线
  - API 服务处理，用户登录、注册、修改信息 
  - 通知服务器，发送推送通知(Notification) 
  - 存储服务，通过 KV 进行 存储、查询 聊天信息

### IM 私信系统 - 消息发送流程
- 一对一聊天，主要的消息发送流程:
   - 用户 A 向 聊天服务 发送消息给 用户 B 
   - 聊天服务从生成器获取消息 ID 聊天服务将消息发到消息队列 
   - 消费保存在 KV 存储中 如果用户在线，则转发消息给用户
   - 如果用户不在线，则转发到通知服务 (Notification)

### IM 私信系统 - 发信箱 / 收信箱
- 两概念:
   - 收件箱(inbox): 该用户收到的消息。
   - 发件箱(outbox): 该用户发出的消息。 
- Timeline 模型:
   - 每个消息拥有一个唯一的顺序ID (SequenceID)，消息按 SequenceID 排序。
   - 新消息写入能自动分配递增的顺序 ID，保证永远 插入队尾。
   - 支持根据顺序 ID 的随机定位，可根据 SequenceID 随机定位到 Timeline 中的某个位置。

### IM 私信系统 - 存储类型选择
- 聊天系统中，消息存储是最主要，通常有海量消息需要存储
- `关系数据库` 还是 `NoSQL数据库`
- 关系数据库主要进行存储用户信息，好友列表，群组信息，通过主从、分片基本满足
- 消息存储比较单一，可以通过 KV 存储
- KV 存储消息的好处:
     - 水平扩展
     - 延迟低
     - 访问成本低

### IM 私信系统 - 消息存储设计(1 on 1)
- 在 1 对 1 聊天消息中，最重要的是数据格式，以及消息主键 message_id，
- 需要保证一定的顺序，并且可以按规则 Scan PrefixKey。
- 消息数据模型:
   - message_id，消息ID。 
   - message_from，消息发送者ID。
   - message_to，消息接收者ID。 
   - content，消息内容。 
   - created_at，消息发送时间。
- 收件箱(Inbox)，KV:
   - `<message_to>_<message_id> : <outbox_message_key>`
- 发件箱(Outbox)，KV:
   - `<message_from>_<message_id >: <message>`

### IM 私信系统 - 消息存储设计(Group chat)

- 在群聊中，存在读写放大问题 所以需要按具体场景考虑主键设计。
- 消息数据模型:
   - channel_id，频道ID。
   - message_id，消息ID。 
   - user_id，消息发送者ID。
   - content，消息内容。 
   - created_at，消息发送时间。
- 收件箱(Inbox)，KV:
  - `<channel_id>_<message_id> : <outbox_message_key>` (多读)
  - `<user_id>_<message_id> : <outbox_message_key>` (多写)
- 发件箱(Outbox)，KV:
  - `<user_id>_<message_id> : <message>`

### IM 私信系统 - 群聊 / 订阅号
- 群聊，较为复杂，通常有多写、多读两种方式; 
- 单件箱(多写)，每个用户都保存一份消息:
   - 消息同步流程比较简单，每个客户端仅需要读取自己的信箱，即可获取新消息
   - 当群组比较小时，成本也不是很高，
      - 例如微信群通 常为 500 用户上限
   - 对群组数量无上限
- 多件箱(多读)，每个群仅保存一份消息:
  - 用户需要同时查询多个信箱 
  - 如果信箱比较多，查询成本比较高 
  - 需要控制群组上限

### IM 私信系统 - 读写扩散
- 一般消息系统中，通常会比较关注消息存储; 
- 主要进行考虑“读”、“写”扩散，也就是性能问题; 在不同场景，可能选择不同的方式:
- 读扩散: 在 IM 系统里的读扩散通常是每两个相 关联的人就有一个信箱，或者每个群一个信箱。
  - 优点:写操作(发消息)很轻量，只用写自己信箱 
  - 缺点:读操作(读消息)很重，需要读所有人信箱
- 写扩散: 每个人都只从自己的信箱里读取消息， 但写(发消息)的时候需要所有人写一份
   - 优点:读操作很轻量
   - 缺点:写操作很重，尤其是对于群聊来说

### IM 私信系统 - 推拉结合模式
- 在⻓连接中，如果想把消息通知所有人，主要有两种模式:
  - 一种是自己拿广播通知所有人，这叫 “推”模式
  - 一种是有人主动来找你要，这叫“拉” 模式。;
- 在 IM 系统中，通常会有三种可能的做法:
  - 推模式
     - 有新消息时服务器主动推给所有端 (iOS、Android、PC 等)
     - 为了保证消息的实时性，一般采用推模式，
  - 拉模式
     - 由前端主动发起拉取消息的请求，
     - 拉模式一般用于获取历史消息
  - 推拉结合模式:
     - 有新消息时服务器会先推一个有新消息的通知给前端，前端接收到通知后就向服务器拉取消息


# References
 - https://draveness.me/golang/docs/part3-runtime/ch06-concurrency/golang-netpoller/ https://www.liwenzhou.com/posts/Go/15_socket/ https://hit-alibaba.github.io/interview/basic/network/HTTP.html https://www.cdn77.com/blog/improving-webperf-security-tls-1-3 https://cloud.google.com/dns/docs/dns-overview?hl=zh-cn https://cloud.tencent.com/developer/article/1030660 https://juejin.cn/post/6844903827536117774 https://xie.infoq.cn/article/19e95a78e2f5389588debfb1c https://tech.meituan.com/2019/03/07/open-source-project-leaf.html https://mp.weixin.qq.com/s/8WmASie_DjDDMQRdQi1FDg https://www.imooc.com/article/265871
 - https://www.infoq.cn/article/the-road-of-the-growth-weixin-backgroundhttps:// systeminterview.com/design-a-chat-system.php
 - https://blog.discord.com/how-discord-stores-billions-of-messages-7fa6ec7ee4c7
 - https://www.facebook.com/notes/facebook-engineering/the-underlying-technology-of- messages/454991608919/
 - https://www.infoq.cn/article/the-road-of-the-growth-weixin-background https://slack.engineering/flannel-an-application-level-edge-cache-to-make-slack-scale/ https://www.infoq.cn/article/emrual7ttkl8xtr-dve4
 - http://www.91im.net/im/1130.html https://xie.infoq.cn/article/19e95a78e2f5389588debfb1c
