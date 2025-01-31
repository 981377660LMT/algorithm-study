Go 语言的 `net/http` 标准库是一个功能完备且易用的 HTTP 客户端、服务端实现。它内部包含大量源代码文件（如 `server.go`、`transport.go`、`client.go`、`request.go` 等），为我们提供了从 **创建 HTTP 服务器**、**处理路由和请求**、**编写响应** 到 **发起 HTTP 客户端请求** 的一整套基础设施。要真正深入理解 `net/http`，需要先掌握其核心设计思路、主要结构和调用流程。本文将从 **HTTP Server 端** 与 **HTTP Client 端** 两大方面讲解相关源码及原理，并辅以关键类型、函数的源码剖析，以便读者能对 `net/http` 的工作机制形成系统认知。

---

# 目录

- [目录](#目录)
  - [概览：整体设计与分层](#概览整体设计与分层)
  - [Server 端核心结构与流程](#server-端核心结构与流程)
    - [2.1 Server 结构](#21-server-结构)
    - [2.2 ListenAndServe 入口](#22-listenandserve-入口)
    - [2.3 TCP 连接的处理：Serve](#23-tcp-连接的处理serve)
    - [2.4 HTTP 协议解析：`conn.serve` \& `readRequest`](#24-http-协议解析connserve--readrequest)
    - [2.5 请求的多路分发：ServeMux](#25-请求的多路分发servemux)
    - [2.6 Handler 处理请求](#26-handler-处理请求)
    - [2.7 ResponseWriter 的实现](#27-responsewriter-的实现)
  - [Client 端核心结构与请求过程](#client-端核心结构与请求过程)
    - [3.1 Transport 与 RoundTripper](#31-transport-与-roundtripper)
    - [3.2 请求的发起与执行](#32-请求的发起与执行)
    - [3.3 连接池与复用](#33-连接池与复用)
  - [请求与响应的结构：Request / Response](#请求与响应的结构request--response)
  - [HTTP/2 支持简介](#http2-支持简介)
  - [关键点与扩展](#关键点与扩展)
    - [6.1 并发处理模型](#61-并发处理模型)
    - [6.2 超时与上下文](#62-超时与上下文)
    - [6.3 中间件与拓展](#63-中间件与拓展)
    - [6.4 Net 包与 I/O 模型](#64-net-包与-io-模型)
  - [小结](#小结)

---

## 概览：整体设计与分层

在 Go 的 `net/http` 包中，大体可以分为如下几个层次：

1. **底层网络通信**：基于 `net` 包提供的 `net.Listener`, `net.Conn` 等抽象进行网络收发。
2. **HTTP 协议解析**：对请求头、响应头、Body 数据做解析和封装，包括 chunked 编码、Content-Length 解析、保持长连接等逻辑。
3. **Server 端**：`Server` 结构、`ListenAndServe` 等，用于启动 HTTP 服务器并管理连接；`ServeMux` 作为默认路由器，`Handler` 作为处理请求的抽象接口。
4. **Client 端**：`Client` 结构、`Transport` 结构，以及背后的连接池、复用与请求调度。
5. **公共数据结构**：`Request`/`Response`，以及辅助的 `Cookie`, `Header`, `URL` 等类型。

每一层都紧密协作，构成了 Go 原生 HTTP 库的使用框架。

---

## Server 端核心结构与流程

### 2.1 Server 结构

在 `server.go` 中，有一个核心的 `Server` 结构，它几乎包含了对 HTTP 服务器的所有可配置选项：

```go
type Server struct {
    Addr    string
    Handler Handler

    TLSConfig *tls.Config
    ReadTimeout  time.Duration
    WriteTimeout time.Duration
    IdleTimeout  time.Duration

    // 其他很多字段, 如最大 header 大小, 最大并发处理数量等
}
```

- `Addr`: 监听地址，常见格式如 `":8080"`。
- `Handler`: 实际处理请求的对象，一般用 `http.Handler` 接口或 `http.HandlerFunc` 来实现。
- `ReadTimeout`, `WriteTimeout`: 读写超时等。
- 运行时会根据 `Server` 的配置来进行监听、处理连接、管理超时与并发等操作。

> **默认行为**：如果你使用 `http.ListenAndServe(":8080", nil)`，它背后会创建一个默认的 `Server` 并调用其 `ListenAndServe` 方法。而 `nil` 表示使用默认的路由器 `DefaultServeMux`。

---

### 2.2 ListenAndServe 入口

最常见的启动服务器方法是：

```go
func ListenAndServe(addr string, handler Handler) error {
    server := &Server{Addr: addr, Handler: handler}
    return server.ListenAndServe()
}
```

它本质只是帮你创建一个 `Server` 实例，然后调用 `(*Server).ListenAndServe`：

```go
func (srv *Server) ListenAndServe() error {
    ln, err := net.Listen("tcp", srv.Addr)
    if err != nil {
        return err
    }
    return srv.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
}
```

- 首先使用 `net.Listen("tcp", srv.Addr)` 获取一个 TCP 监听器；
- 包装成 `tcpKeepAliveListener`（设置一些 TCP keep-alive 选项）；
- 调用 `srv.Serve(listener)` 真正进入**事件循环**，开始接收连接并处理请求。

---

### 2.3 TCP 连接的处理：Serve

`func (srv *Server) Serve(l net.Listener) error` 是服务端核心逻辑所在：

1. **主循环**：在一个 for 循环中，调用 `l.Accept()` 接收新的连接；若出错会决定是否重试或直接返回。
2. **启动 goroutine 处理**：每拿到一个新的 `net.Conn`，会包装成 `srv.newConn(conn)` 返回一个 `*conn` 结构，然后起一个新 goroutine 执行 `c.serve()`.

简化示意：

```go
for {
    rw, e := l.Accept()
    if e != nil {
        // 错误处理 ...
        continue
    }
    c := srv.newConn(rw)
    go c.serve()
}
```

这就意味着 **每个 TCP 连接** 对应**一个** Goroutine，负责该连接上的全部请求-响应交互（HTTP/1.1 可能在一个连接里多次请求/响应）。

---

### 2.4 HTTP 协议解析：`conn.serve` & `readRequest`

`conn.serve()` 是单个连接的处理流程，其核心代码在 `server.go`：

- 循环调用 `w, err := c.readRequest(context.Background())` 解析下一个 HTTP 请求头；
- 若有请求，就调用 `serverHandler{c.server}.ServeHTTP(w, w.req)` 分发给具体 `Handler`；
- 处理完后，根据 `Connection: keep-alive` 或请求头判断是否继续读取下一次请求；
- 如果遇到 EOF 或错误，退出循环并关闭连接。

**Request 解析**主要是通过一个带缓冲的 `*bufio.Reader` 结合 `textproto.Reader` 来读取 HTTP 头部和 body：

```go
func (c *conn) readRequest(ctx context.Context) (w *response, err error) {
    // 1. 解析请求行，如 "GET /index.html HTTP/1.1"
    // 2. 解析 header
    // 3. 根据 header 信息（Content-Length / Transfer-Encoding）构造 Request.Body
    // 4. 生成 http.Request 结构
    // 5. 最终返回 *response 结构（内部持有 *Request, 用于后续写回）
}
```

---

### 2.5 请求的多路分发：ServeMux

当服务端获取到一个完整 `Request` 后，会调用 `serverHandler{srv}.ServeHTTP(w, req)`，而其内部逻辑是：

```go
handler := srv.Handler
if handler == nil {
    handler = DefaultServeMux // 如果用户没指定，就用默认
}
handler.ServeHTTP(w, req)
```

- **默认 Handler**：如果不指定，则使用 `DefaultServeMux`；这是一种简单的路由器，通过 `HandleFunc(pattern, func)` 实现 URL 路径的匹配和分发。
- **自定义 Handler**：如果提供了 `srv.Handler`，则请求直接交给这个处理器。

**DefaultServeMux** 会在内部维护一个路由表（`map[string]muxEntry`），根据请求的 Path 找到最合适的 handler 进行调用。

---

### 2.6 Handler 处理请求

Go 定义了这样一个接口：

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

- 用户只要实现了 `ServeHTTP` 方法，就能成为一个 Handler；
- 最常见做法是：
  ```go
  http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
      w.Write([]byte("Hello, world!"))
  })
  ```
  这里的 `HandleFunc` 是一个帮助函数，它把回调转换成实现了 `ServeHTTP` 的 `HandlerFunc` 类型。

当调用 `handler.ServeHTTP(w, req)` 时，用户的业务逻辑就可以通过 `w` 写回响应内容、设置响应码、header 等。

---

### 2.7 ResponseWriter 的实现

`ResponseWriter` 是接口，实际在内部会有一个具体类型，如 `*response` 或 `*chunkWriter`.

- 当用户调用 `w.Write()` 时，底层会把数据写到该连接对应的 `bufio.Writer` 中，并自动生成 HTTP 响应头（若尚未写出）。
- 同时也会进行一些判断，比如是否要分块传输（chunked），是否有 Content-Length，是否连接复用等。
- 写完后，如果没有再次请求，就关闭连接或保持 keep-alive。

---

## Client 端核心结构与请求过程

### 3.1 Transport 与 RoundTripper

在 `net/http` 中，HTTP 客户端由 `Client` 结构配合 `Transport` 一起实现。主要接口是：

```go
type RoundTripper interface {
    RoundTrip(*Request) (*Response, error)
}
```

- `Transport` 是 Go 默认实现的 `RoundTripper`，负责连接管理、请求发送、响应读取、复用 keep-alive 等；
- `Client` 结构通常会含有一个 `Transport` 字段，执行请求时就是调用 `Transport.RoundTrip()`。

### 3.2 请求的发起与执行

用户端一般写法：

```go
resp, err := http.Get("https://example.com")
```

内部相当于：

```go
func Get(url string) (*Response, error) {
    return DefaultClient.Get(url)
}
```

而 `DefaultClient` 相当于：

```go
var DefaultClient = &Client{}
```

它在 `client.go` 中：当你调用 `Client.Get(url)` → `Client.do(req)` → `Transport.RoundTrip(req)` → 最终建立或复用 TCP 连接（若 https 则走 TLS），然后发送请求行、头部、body，等待响应；  
请求结束时返回 `*Response` 给上层。

### 3.3 连接池与复用

`Transport` 内部会维护一个连接池（或称 idleConn），当请求结束且 `Connection: keep-alive`，这个连接会被保留在空闲池里，下次遇到同一个目标（host:port + TLS 配置）时可直接复用，减少握手/三次握手开销。

- `Transport` 有一些可调参数，如 `MaxIdleConns`、`IdleConnTimeout` 等，用来控制空闲连接的最大数和生命周期。
- 如果池子里没有可用空闲连接，则创建新连接。

---

## 请求与响应的结构：Request / Response

在 `request.go` 与 `response.go` 中，定义了 HTTP 请求、响应的核心数据结构：

- `type Request struct {  
    Method string  
    URL *url.URL  
    Header Header  
    Body io.ReadCloser  
    ContentLength int64  
    ...  
}`
- `type Response struct {  
    Status     string  
    StatusCode int  
    Header     Header  
    Body       io.ReadCloser  
    ...  
}`

它们都将首部（header）封装为 `Header` 类型，使用一个 `map[string][]string` 来存储多值 header，且大小写不敏感。`Body` 均为一个流式接口（实现 `io.ReadCloser`），可边读边处理。

---

## HTTP/2 支持简介

Go 1.6 起，`net/http` 默认内置对 HTTP/2 的支持（带 TLS/ALPN）。

- 若是 HTTPS 且远端支持 HTTP/2，会自动升级到 http2；
- 底层切换到新的帧协议和多路复用；
- 对用户而言，API 不变：`http.Client` / `http.Server` 均可正常使用；
- 具体实现在子包 `golang.org/x/net/http2` 或已整合进标准库中，通过 `http2.ConfigureServer` 等 API 进行配置。

---

## 关键点与扩展

### 6.1 并发处理模型

- **Server 端**：一个连接一条 goroutine（HTTP/1.1 下），若请求行和 header 解析成功，就会再进一步调度 Handler 处理请求。
- **Handler**：可以并发处理，多个 goroutine 同时执行 `ServeHTTP`；要注意并发安全。
- **Client 端**：发送请求时，若有空闲连接，会直接用，若无则建新连接。每次请求通常在一个新的 goroutine 中执行（如果用户并发调用），`Transport` 负责同步处理复用等细节。

### 6.2 超时与上下文

- 在 `Server` 可以设置 `ReadTimeout`, `WriteTimeout`, `IdleTimeout` 等来防止阻塞或无效连接耗时。
- 在 `Client` 可以设置 `Timeout`, 或使用 `context.Context` 来控制请求全局超时、取消等。
- `http.Request` 里可以带 `Context()`，在服务器端 Handler 可以用 `<-r.Context().Done()` 来感知客户端中断。

### 6.3 中间件与拓展

- 通过自定义 Handler 链的方式可以实现**中间件**（类似 `func(w http.ResponseWriter, r *http.Request, next http.Handler)`）。
- 自定义路由器、实现 `http.RoundTripper`、在 `Transport` 里嵌入代理逻辑，都是常见的拓展方式。

### 6.4 Net 包与 I/O 模型

- `net/http` 是构建于 `net` 包之上，最终使用非阻塞 I/O + goroutine 并发处理。
- Go runtime 层面会把大量阻塞的 I/O 调用调度到多线程上，以保证并发性。

---

## 小结

`net/http` 标准库以 **简单易用** 作为主要目标，通过**分层抽象**把网络监听、HTTP 协议解析、请求多路分发、路由与处理器、连接池及 Keep-Alive 管理、并发模型**等细节封装起来，形成了一个**内置、开箱即用\*\*的 Web 开发基础设施：

1. **Server** 端：`Server` 结构控制监听、连接处理；`Handler` 以及 `ServeMux` 完成请求分发与处理；
2. **Client** 端：`Client`、`Transport`、`RoundTripper` 协作，提供复用连接、并发请求、以及 HTTP/2 等特性；
3. **并发 & 简洁**：通过 Goroutine + channel / context 模型处理高并发；用户只需专注业务 Handler；
4. **可扩展**：使用中间件、定制 Transport、钩子函数等方式进行功能扩展；或直接使用 `net/http` 构建自己的高级框架。

理解 `net/http` 源码可以帮助你在 **排查连接故障**、**调优高并发服务**、**编写中间件** 等方面更加得心应手，同时也能启发你在其他网络开发中模仿其优雅的抽象与并发模型。
