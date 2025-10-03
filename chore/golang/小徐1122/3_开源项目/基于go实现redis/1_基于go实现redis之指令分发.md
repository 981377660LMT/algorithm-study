好的，这是 `基于go实现redis之指令分发` 的详细内容。

---

### 基于 go 实现 redis 之指令分发

#### 0 前言

欢迎回来，由我们一起继续推进技术分享专题—— **【基于 go 实现 redis】**。

此前我已于 github 开源项目——`goredis`，由于我个人水平有限，如有实现不到位之处，欢迎批评指正：https://github.com/xiaoxuxiansheng/goredis

![图片](https://img-blog.csdnimg.cn/img_convert/06758652618037380128795085351235.png)

本系列正是围绕着开源项目 `goredis` 展开，共分为四篇内容，本篇为其中的第二篇——指令分发篇：

- **基于 go 实现 redis 之主干框架（已完成）**： 在宏观视角下纵览 goredis 整体架构，梳理各模块间的关联性
- **基于 go 实现 redis 之指令分发（本篇）**： 聚焦介绍 goredis 服务端如何启动和运行，并在接收客户端请求后实现指令协议的解析和分发
- **基于 go 实现 redis 之存储引擎（待填坑）**： 聚焦介绍数据存储层中单协程无锁化执行框架，各类基本数据类型的底层实现细节，以及过期数据的惰性和定期回收机制
- **基于 go 实现 redis 之数据持久化（待填坑）**： 介绍 goredis 关于 aof 持久化机制的实现以及有关于 aof 重写策略的执行细节

#### 1 架构梳理

![图片](https://img-blog.csdnimg.cn/img_convert/06758652618037380128795085351235.png)

正如前言中所介绍，本篇我们聚焦探讨两个问题：

- **goredis 服务端如何运行**： 即如何快速搭建一个 TCP server，能够正常接收处理来自客户端的 TCP 连接和请求。
- **goredis 请求内容解析**： 即遵循怎样的协议，针对 TCP 请求内容进行解析转换，并给出遵循协议的响应格式。

为了解决上述问题，本篇涉及到的模块包括：

- **应用程序 `application`**： 对应为全局 goredis 应用程序的抽象载体（第 2 章介绍）。
- **服务端 `server`**： 对应为 TCP server 的具体实现模块（第 2 章介绍）。
- **指令分发器 `handler`**： 负责为到来的 TCP 连接服务，将请求内容转为操作指令，并分发给存储引擎层（第 3 章介绍）。
- **协议解析器 `parser`**： 遵循 Redis 文本解析协议，将请求内容解析为 Redis 操作指令（第 3 章介绍）。

#### 2 服务运行

本章中，我们将展开介绍，如何支撑 `goredis` 应用程序的启动，并作为一个 TCP server 面向客户端提供服务。

![图片](https://img-blog.csdnimg.cn/img_convert/06758652618037380128795085351235.png)

`goredis` 启动与运行流程如上图所示，分为几个核心步骤：

- **配置加载**： 感知并启用用户设定的运行策略。
- **server 启动**： server 启动后持续监听 TCP 端口，接收到来的 TCP 连接。
- **handler 启动**：handler 用于服务到来的 TCP 连接，将请求内容转为 Redis 指令，并分发至存储引擎侧处理。

##### 2.1 启动入口

`goredis` 程序启动入口位于 `main.go`，完成对 `conf` 的加载，`server` 和 `application` 实例的构造，然后通过 `application` 一键启动程序：

```go
func main() {
    server, err := app.ConstructServer()
    if err != nil {
        panic(err)
    }

    app := app.NewApplication(server, app.SetUpConfig())
    defer app.Stop()

    if err := app.Run(); err != nil {
        panic(err)
    }
}
```

##### 2.2 配置加载

`redis` 运行策略的配置文件位于 `./redis.conf`。目前支持配置的策略包括：

- `bind`： 服务端支持接收的客户端源 IP。
- `port`： 服务端启动时监听的端口。
- `appendonly`： (`true`/`false`) 是否启用 AOF 持久化机制。
- `appendfilename`： 持久化时使用的 AOF 文件路径+名称。
- `appendfsync`：AOF 持久化策略。(1) `always`： 针对每个指令进行 `fsync` 落盘； (2) `everysec`： 指令先集中溢写到设备 buffer， 每秒集中 `fsync` 落盘一次； (3) `no`： 指令溢写到设备 buffer 后即结束流程，由设备自行决定何时落盘。
- `auto-aof-rewrite-after-cmds`： 每持久化多少条指令后进行一次 AOF 文件重写。

```ini
# 接收的源 ip 地址. 设置为 0.0.0.0 代表接受任意源 ip.
bind 0.0.0.0
# 程序绑定的端口号
port 6379

# 是否启用 aof 持久化策略
appendonly yes
# aof 文件名
appendfilename appendonly.aof
# aof 级别，分为 always | everysec | no
appendfsync everysec
# 每执行多少次 aof 指令后，进行一次指令重写
auto-aof-rewrite-after-cmds 1000
```

`goredis` 启动时，首先会通过 `SetUpConfig` 方法对配置文件进行加载，生成 `Config` 实例并应用到服务各个模块中。此处代码位于 `app/config.go`：

```go
// 与 redis.conf 一一对应的各项配置参数
type Config struct {
    Bind                    string `cfg:"bind"`                        // 接收的源 ip 地址
    Port                    int    `cfg:"port"`                        // 程序绑定的端口号
    AppendOnly_             bool   `cfg:"appendonly"`                  // 是否启用 aof 持久化
    AppendFileName_         string `cfg:"appendfilename"`              // aof 文件路径+名称
    AppendFsync_            string `cfg:"appendfsync"`                 // aof 持久化策略
    AutoAofRewriteAfterCmd_ int    `cfg:"auto-aof-rewrite-after-cmds"` // 每执行多少次 aof 操作后，进行一次重写
}

// 加载生成配置项实例
func SetUpConfig() *Config {
    // 懒加载机制，保证全局只加载一次配置
    confOnce.Do(func() {
        defer func() {
            // 倘若配置加载失败，则使用默认的兜底策略
            if globalConf == nil {
                globalConf = defaultConf()
            }
        }()
        // 打开根目录下的 redis.conf 文件
        file, err := os.Open("./redis.conf")
        if err != nil {
            return
        }
        // 读取文件内容，反序列化为 conf 实例
        globalConf = setUpConfig(file)
    })

    return globalConf
}

// 默认的兜底策略
func defaultConf() *Config {
    return &Config{
        Bind:        "0.0.0.0", // 默认面向所有源 ip 地址开放
        Port:        6379,      // 默认端口号为 6379
        AppendOnly_: false,     // 默认不启用 aof
    }
}
```

##### 2.3 服务运行

支撑 `goredis` 运行 TCP server 的核心代码位于 `server/server.go` 文件。在实现过程中，体现到的核心技术细节主要包含三点：

- **IO 多路复用技术**：server 运行过程中，在监听 TCP 端口和处理 TCP 连接时，通过对 Go 语言 `net` 包的封装，隐藏了底层 `epoll` (Linux) 或 `kqueue` (macOS) 的复杂性。`listener.Accept()` 和 `conn.Read()` 等操作都构建于此之上。
- **一比一异步协程服务连接**： 对于到来的每个 TCP 连接，会为其启动一个独立的 goroutine 进行处理，实现了高并发的 "Goroutine-per-Connection" 模型。
- **优雅关闭策略**： 通过 `context` 和 `sync.WaitGroup` 等工具控制和守护其生命周期，确保在服务关闭时能妥善处理所有活跃连接。

![图片](https://img-blog.csdnimg.cn/img_convert/06758652618037380128795085351235.png)

有关服务端 `server` 的类定义如下，其中包含如下核心成员属性：

- **单例工具 `sync.Once`**： 通过 `runOnce`、`stopOnce` 两个单例工具，规避单 `server` 实例下启动和停止动作的重复执行。
- **指令分发器实例 `handler`**： `server` `Accept` 到 TCP 连接后，会分配给 `handler` 进行处理。
- **关闭控制器 `stopc`**： 通过 `stopc` 感知到服务退出信号，进行资源回收。

```go
type Server struct {
    runOnce  sync.Once
    stopOnce sync.Once
    handler  Handler
    // ...
    stopc    chan struct{}
}
```

服务端启动的入口方法为 `server.Serve()`：

- **退出信号感知**： 启动一个异步协程监听操作系统信号（如 `SIGINT`, `SIGTERM`），以便完成退出前的资源回收、优雅关闭。
- **端口监听器初始化**： 根据配置的 IP 和端口，创建 `net.Listener` 实例。
- **服务端主运行流程**： 基于 `for` 循环 + `listener.Accept()` 模型，持续探测并接收到来的 TCP 连接。为每个到来的连接分配 goroutine，并将其托付给指令分发器 `handler` 进行处理。

```go
func (s *Server) Serve(address string) error {
    if err := s.handler.Start(); err != nil {
        return err
    }
    var _err error
    s.runOnce.Do(func() {
        // 监听进程信号，完成程序退出前的资源回收
        exitWords := []os.Signal{syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT}
        sigc := make(chan os.Signal, 1)
        signal.Notify(sigc, exitWords...)
        closec := make(chan struct{}, 1)
        pool.Submit(func() {
            select {
            case <-sigc:
                closec <- struct{}{}
            case <-s.stopc:
                closec <- struct{}{}
            }
        })

        // 为指定端口申请监听器
        listener, err := net.Listen("tcp", address)
        if err != nil {
            _err = err
            return
        }
        // 启动经典的 tcp server 运行框架
        s.listenAndServe(listener, closec)
    })
    return _err
}

func (s *Server) listenAndServe(listener net.Listener, closec chan struct{}) {
    errc := make(chan error, 1)
    defer close(errc)

    ctx, cancel := context.WithCancel(context.Background())
    pool.Submit(func() {
        select {
        case <-closec:
        case <-errc:
        }
        cancel()
        s.handler.Close()
        _ = listener.Close()
    })

    var wg sync.WaitGroup
    for {
        conn, err := listener.Accept()
        if err != nil {
            if ne, ok := err.(net.Error); ok && ne.Timeout() {
                time.Sleep(5 * time.Millisecond)
                continue
            }
            errc <- err
            break
        }

        // 为每个到来的 conn 分配一个 goroutine 处理
        wg.Add(1)
        pool.Submit(func() {
            defer wg.Done()
            s.handler.Handle(ctx, conn)
        })
    }
    // 通过 waitGroup 保证优雅退出
    wg.Wait()
}
```

#### 3 指令分发

##### 3.1 指令分发器

每当有 TCP 连接到达后，`server` 层会将其统一分配到下层，由 `handler` 为其提供服务。

入口方法为 `handler.Handle()`。为了更好地支持优雅关闭策略，`handler` 会缓存活跃的连接，以便在退出前对连接进行关闭：

```go
func (h *Handler) Handle(ctx context.Context, conn net.Conn) {
    h.mu.Lock()
    if h.closed.Load() {
        h.mu.Unlock()
        _ = conn.Close()
        return
    }
    // 当前 conn 缓存起来
    h.conns[conn] = struct{}{}
    h.mu.Unlock()

    // 核心逻辑所在
    h.handle(ctx, conn)
}

// 关闭 handler
func (h *Handler) Close() {
    h.Once.Do(func() {
        h.closed.Store(true)
        h.mu.RLock()
        defer h.mu.RUnlock()
        // 依次关闭连接
        for conn := range h.conns {
            _ = conn.Close()
        }
        h.conns = nil
        // 关闭存储引擎
        h.db.Close()
        // 关闭持久化模块
        h.persister.Close()
    })
}
```

![图片](https://img-blog.csdnimg.cn/img_convert/06758652618037380128795085351235.png)

`handler` 层会紧密依赖协议解析器 `parser`：

1.  将到来的 TCP 连接 (`io.ReadWriter`) 转换成 Go 语言中易于处理的通道 (`channel`) 形式。
2.  借由 `parser` 的能力将到来的请求内容一一转换成符合 Redis 协议的指令。
3.  进一步将其分发给下层的存储引擎 `db` 进行处理。

在此处的实现上，我们将用于传递指令的通道类比为一股源源不断的水流——`stream`，将到来的每一个指令类比为水流中的一个水滴——`droplet`，并采用 Go 语言中经典的 `for-select` 模式组织 `handler` 流程的运行框架。

```go
func (h *Handler) handle(ctx context.Context, conn io.ReadWriter) {
    // 将连接转换成通道 channel 实例，抽象为一股水流
    stream := h.parser.ParseStream(conn)
    for {
        select {
        case <-ctx.Done():
            return
        // 承载水流中的每一个水滴
        case droplet := <-stream:
            if err := h.handleDroplet(ctx, conn, droplet); err != nil {
                // 如果是EOF或连接关闭等错误，则正常退出
                return
            }
        }
    }
}

// 处理到来的水滴(请求指令)
func (h *Handler) handleDroplet(ctx context.Context, conn io.ReadWriter, droplet *Droplet) error {
    // 请求是否终止. 如包含了 EOF 类错误
    if droplet.Terminated() {
        return droplet.Err
    }

    // 请求指令必须为 multiBulkReply 类型
    multiReply, ok := droplet.Reply.(MultiReply)
    if !ok {
        return errors.New("invalid request")
    }

    // 将请求指令发往存储引擎 db
    if reply := h.db.Do(ctx, multiReply.Args()); reply != nil {
       // 将 db 给予的响应通过 tcp 连接返回给客户端
        _, _ = conn.Write(reply.ToBytes())
    }
    return nil
}
```

##### 3.2 RESP

在正式介绍协议解析模块之前，我们有必要先对 Redis 文本解析协议（RESP）的理论知识进行补充。

![图片](https://img-blog.csdnimg.cn/img_convert/06758652618037380128795085351235.png)

RESP 全称 Redis Serialization Protocol，是 Redis 用于解析 TCP 文本内容采用的协议：

- **内容分类**： RESP 将文本内容分为简单字符串、错误信息、整数、定长字符串以及数组五类。
- **换行分隔符**： RESP 以行为解析单元，文本内容一律以换行符 `\r\n` 进行分隔。
- **二进制安全**： RESP 中定长字符串和数组类型的文本内容能够保证二进制安全，即使正文内容混淆了类似 `\r\n` 这样的特殊字符，也不会发生混淆歧义。

接下来我们针对上文提到的五类内容逐一展开介绍：

- **简单字符串 (Simple String)**： 以 `+` 开头，随后紧跟文本内容直到换行符出现。非二进制安全，主要用于 `OK` 等简单响应。
  `+OK\r\n`

- **错误 (Error)**： 以 `-` 开头，随后紧跟错误信息直到换行符出现。非二进制安全，用于服务端报错。
  `-ERR syntax error\r\n`

- **整数 (Integer)**： 以 `:` 开头，随后紧跟一个 64 位有符号整数。用于返回整数结果，如 `INCR` 命令。
  `:1\r\n`

- **定长字符串 (Bulk String)**： 由两行内容组成：第一行以 `$` 开头，后跟字符串字节长度；第二行是实际的字符串内容。二进制安全，用于返回单个值，如 `GET` 命令的结果。
  表示字符串 "hello"：
  `$5\r\nhello\r\n`

- **数组 (Array / Multi Bulk String)**： 由 `2n + 1` 行内容组成（n 为数组长度）：第一行以 `*` 开头，后跟数组元素个数 n；接下来每组元素都遵循其自身的 RESP 格式（通常是 Bulk String）。二进制安全，客户端发送的命令和一些命令（如 `LRANGE`）的返回值都采用此格式。
  表示 `SET key value` 命令：
  ```
  *3\r\n
  $3\r\n
  SET\r\n
  $3\r\n
  key\r\n
  $5\r\n
  value\r\n
  ```

##### 3.3 协议解析器

协议解析器 `parser` 模块遵循 RESP，针对各类文本内容，定义了相应的解析方法。对应代码位于 `protocol/parser.go` 文件中：

```go
// 协议解析器 parser
type Parser struct {
    // 以文本首个字符为标志，映射到不同文本类型的解析方法
    lineParsers map[byte]lineParser
    // ...
}

// 协议解析器构造器函数
func NewParser() handler.Parser {
    p := Parser{}
    // 以文本首个字符为标志，映射到不同文本类型的解析方法
    p.lineParsers = map[byte]lineParser{
        '+': p.parseSimpleString, // 解析简单字符串
        '-': p.parseError,        // 解析错误
        ':': p.parseInt,          // 解析整数
        '$': p.parseBulk,         // 解析定长字符串
        '*': p.parseMultiBulk,    // 解析数组
    }
    return &p
}
```

指令分发器 `handler` 接收到 TCP 连接时，会调用 `parser.ParseStream` 方法。这个方法的设计非常巧妙：

1.  构造一个 `channel` 实例返回给 `handler`。
2.  异步启动一个 goroutine，该协程持续从 TCP 连接中读取数据，遵循 RESP 将其解析成 Redis 指令，然后通过 `channel` 发往 `handler`。

```go
// 由 parser 提供能力，将 tcp 连接中的请求内容解析成指令并通过 channel 发往 handler
func (p *Parser) ParseStream(reader io.Reader) <-chan *handler.Droplet {
    // 构造 channel 实例
    ch := make(chan *handler.Droplet)
    pool.Submit(func() {
        // 异步启动 goroutine，负责完成 tcp 请求内容解析，并通过 channel 传输
        p.parse(reader, ch)
    })
    return ch
}

// 负责完成 tcp 请求内容解析，并通过 channel 传输
func (p *Parser) parse(rawReader io.Reader, ch chan<- *handler.Droplet) {
    defer close(ch)
    reader := bufio.NewReader(rawReader)
    for {
        // 以换行符为分割，读取首行内容
        firstLine, err := reader.ReadBytes('\n')
        if err != nil {
            ch <- &handler.Droplet{Err: err}
            return
        }

        length := len(firstLine)
        if length <= 2 || firstLine[length-2] != '\r' {
            continue // or handle error
        }

        firstLine = bytes.TrimSuffix(firstLine, []byte{'\r', '\n'})
        // 以文本首个字符为标志，映射到不同文本类型的解析方法
        lineParseFunc, ok := p.lineParsers[firstLine[0]]
        if !ok {
            // handle error
            continue
        }
        // 解析成指令内容后，通过 channel 发送往 handler
        ch <- lineParseFunc(firstLine, reader)
    }
}
```

#### 4 使用示例

介绍完 `goredis` 的服务运行框架以及文本解析协议后，接下来于本章中给出 `goredis` 的使用示例。

![图片](https://img-blog.csdnimg.cn/img_convert/06758652618037380128795085351235.png)

1.  **首先，通过 `goredis` 下的启动脚本一键启动服务：**

    ```sh
    ./start.sh
    ```

2.  **接下来，使用 `redis-cli` 或 `telnet` 访问服务对应的 TCP 端口，建立 TCP 连接：**

    ```sh
    telnet localhost 6379
    ```

3.  **成功建立 TCP 连接后，遵循 RESP 协议，执行一笔 `SET` 指令——`SET a b`**
    (输入以下内容，注意换行)

    ```
    *3
    $3
    set
    $1
    a
    $1
    b
    ```

4.  **设置成功后，获得简单字符串响应 `+OK`：**

    ```
    +OK
    ```

5.  **接下来遵循 RESP 协议，执行一笔 `GET` 指令——`GET a`**

    ```
    *2
    $3
    get
    $1
    a
    ```

6.  **最终成功获得定长字符串响应，内容为 `b`：**
    ```
    $1
    b
    ```

#### 5 展望

至此为指令分发篇的全部内容，在此对本系列内容做个小结和展望：

- **基于 go 实现 redis 之主干框架（已完成）**： 在宏观视角下纵览 goredis 整体架构，梳理各模块间的关联性。
- **基于 go 实现 redis 之指令分发（已完成）**： 聚焦介绍 goredis 服务端如何启动和运行，并在接收客户端请求后实现指令协议的解析和分发。
- **基于 go 实现 redis 之存储引擎（待填坑）**： 聚焦介绍数据存储层中单协程无锁化执行框架，各类基本数据类型的底层实现细节，以及过期数据的惰性和定期回收机制。
- **基于 go 实现 redis 之数据持久化（待填坑）**： 介绍 goredis 关于 AOF 持久化机制的实现以及有关于 AOF 重写策略的执行细节。
