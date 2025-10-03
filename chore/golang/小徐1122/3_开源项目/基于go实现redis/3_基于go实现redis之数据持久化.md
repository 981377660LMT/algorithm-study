好的，这是 `基于go实现redis之数据持久化` 的详细内容。

---

### 基于 go 实现 redis 之数据持久化

#### 0 前言

欢迎回来，由我们一起继续推进技术分享专题—— **【基于 go 实现 redis 之数据持久化】**。

此前我已于 github 开源项目——`goredis`，由于我个人水平有限，如有实现不到位之处，欢迎批评指正：https://github.com/xiaoxuxiansheng/goredis

![图片](https://img-blog.csdnimg.cn/img_convert/06758652618037380128795085351235.png)

本系列正是围绕着开源项目`goredis`展开，共分为四篇内容，本文为该系列的完结篇——数据持久化篇：

- **基于 go 实现 redis 之主干框架（已完成）**： 在宏观视角下纵览 goredis 整体架构，梳理各模块间的关联性
- **基于 go 实现 redis 之指令分发（已完成）**： 聚焦介绍 goredis 服务端如何启动和运行，并在接收客户端请求后实现指令协议的解析和分发
- **基于 go 实现 redis 之存储引擎（已完成）**： 聚焦介绍数据存储层中单协程无锁化执行框架，各类基本数据类型的底层实现细节，以及过期数据的惰性和定期回收机制
- **基于 go 实现 redis 之数据持久化（本篇）**： 介绍 goredis 关于 aof 持久化机制的实现以及有关于 aof 重写策略的执行细节

（此外，这里需要特别提及一下，在学习过程中，很大程度上需借助了 `hdt3213` 系列博客和项目的帮助，在此致敬一下：https://www.cnblogs.com/Finley/category/1598973.html）

#### 1 redis 持久化机制

Redis 中的数据均存储在内存，而内存属于一种易失性存储介质，一旦进程或者机器崩溃，都会导致数据的丢失。为了提高数据存储的稳定性，Redis 建立了两种将内存数据溢写到磁盘的持久化机制，包括 RDB 和 AOF 两种。

##### 1.1 RDB

![图片](https://img-blog.csdnimg.cn/img_convert/06758652618037380128795085351235.png)

RDB，全称 Redis Database，指的是一种将内存数据序列化生成快照文件的持久化机制。

RDB 持久化方式注重的是结果。由于每轮 RDB 都涉及到全量数据的操作，因此从单次持久化过程来看显得比较耗费性能，但是从结果来看，RDB 是一种比较节省和高效的持久化方式：一方面，RDB 文件可以对文件内容高度压缩，从而实现存储空间的节省；另一方面，快照形式的持久化内容能保证不会存在冗余的内容，在加载还原内存数据时也会更加的高效。

##### 1.2 AOF

![图片](https://img-blog.csdnimg.cn/img_convert/06758652618037380128795085351235.png)

AOF，全称 Append Only File，是一种通过记录增量变化信息，从而回溯出全量数据的持久化策略。

AOF 的优势是，每次以指令的粒度进行持久化，因此从单次行为来看比较轻便和高效，但是这种方式也不可避免地存在两个缺陷：

- **持久化数据存在冗余**：AOF 持久化方式注重的是过程。针对同一笔 KV 数据，不论执行多少次变更操作都会被 AOF 事无巨细地记录下来。
- **数据还原流程低效**： 在还原内存数据时，需要通过遍历执行指令的方式，完整回溯一遍内存数据库的变更时间线，是一种很低效的数据加载方式。

受限于笔者水平，对 Redis RDB 持久化流程的认知水平还有所不足，本次 `goredis` 项目实现中，仅涉及到对 AOF 持久化机制的实现。

#### 2 goredis 持久化模块定位

##### 2.1 持久化模块接口定义

![图片](https://img-blog.csdnimg.cn/img_convert/06758652618037380128795085351235.png)

在 `goredis` 项目中，定义出一个持久化模块 `persister`，该模块包含两个明确的核心作用：

- **数据持久化**： 内存数据库发生变化时，持久化记录增量变更内容。
- **数据重加载**： 在服务重启时，能够读取持久化内容还原出内存中的数据。

`goredis` 针对于持久化模块定义了一个抽象 `interface`，对应声明了 `Reloader` 和 `PersistCmd` 两个方法，代码位于 `handler/persister.go`：

```go
// 持久化模块
type Persister interface {
    // 获取一个 reader，用于读取之前持久化的内容
    Reloader() (io.ReadCloser, error)
    // 写入增量持久化指令的入口
    PersistCmd(ctx context.Context, cmd [][]byte)
    // 关闭持久化模块
    Close()
}
```

##### 2.2 增量指令持久化

![图片](https://img-blog.csdnimg.cn/img_convert/06758652618037380128795085351235.png)

在 DB 的数据执行层中，针对数据的操作分为读操作和写操作两类：

- **读操作**： 如 `get`，`sismember` 等指令，并不会引起底层数据状态的变化，因此无需持久化。
- **写操作**： 如 `hset`、`sadd` 等，则会通过 `persister` 统一对操作指令进行持久化。

以 `hset`、`sadd` 两个写操作指令的代码实例加以说明：

```go
func (k *KVStore) HSet(cmd *database.Command) handler.Reply {
    // 1. ...执行数据写入 hashmap 操作
    // 2. 对 hset 指令进行持久化
    k.persister.PersistCmd(cmd.Ctx(), cmd.Cmd()) // 持久化
    // 3. ...返回响应结果
}

func (k *KVStore) SAdd(cmd *database.Command) handler.Reply {
    // 1. ...执行数据写入 set 操作
    // 2. 持久化 sadd 指令
    k.persister.PersistCmd(cmd.Ctx(), cmd.Cmd())
    // 3. ...返回响应结果
}
```

值得一提的是，涉及到过期时间的写操作在持久化时需要一些特殊的处理技巧。持久化过期时间的要点在于，不能以一个相对的过期时间 TTL（time to live）作为持久化目标，否则可能导致过期信息的失真。

基于以上，`goredis` 在持久化过期时间时，统一采用的是一个过期时间点绝对值，无论是 `set ex` 还是 `expire` 指令，在持久化内容上都需要转为 `expireat` 的形式。

下面以 `Expire` 指令持久化流程的代码示例加以说明：

```go
// 执行 expire 指令.
func (k *KVStore) Expire(cmd *database.Command) handler.Reply {
    // ...
    // 根据 ttl 推算出过期时间点
    expireAt := time.Now().Add(time.Duration(ttl) * time.Second)
    // 构造出要持久化的 expireAt 指令
    _cmd := [][]byte{[]byte(database.CmdTypeExpireAt), []byte(key), []byte(lib.TimeSecondFormat(expireAt))}
    // 在 expireAt 方法中，完成过期时间设置操作，并且完成指令持久化
    return k.expireAt(cmd.Ctx(), _cmd, key, expireAt)
}

func (k *KVStore) expireAt(ctx context.Context, cmd [][]byte, key string, expireAt time.Time) handler.Reply {
    // 1. 执行过期时间设置操作
    k.expire(key, expireAt)
    // 2. 持久化过期指令
    k.persister.PersistCmd(ctx, cmd) // 持久化
    // ...
}
```

##### 2.3 持久化数据恢复

![图片](https://img-blog.csdnimg.cn/img_convert/06758652618037380128795085351235.png)

在 `goredis` 服务重启时，会通过持久化模块 `persister` 获取指令加载器 `reloader`，从中依次读取此前完成持久化的指令。在具体技术实现上，`goredis` 会将 `reloader` 包装成一个虚拟连接，模拟服务端接收到客户端请求指令的流程，依次将指令分发到存储引擎层，最终实现内存数据的恢复：

```go
// 启动指令分发器
func (h *Handler) Start() error {
    // 1. 加载持久化指令，还原出内存中的数据
    reloader, err := h.persister.Reloader()
    if err != nil {
        return err
    }
    if reloader == nil {
        return nil
    }
    defer reloader.Close()
    // 2. 把持久化指令 reloader 包装成一个 fake 的连接，执行 handler 分发指令主流程
    h.handle(SetLoadingPattern(context.Background()), newFakeReaderWriter(reloader))
    return nil
}

func (h *Handler) handle(ctx context.Context, conn io.ReadWriter) {
    // 借助 parser 将连接转为 chan 形式，并持续接收解析到的指令
    stream := h.parser.ParseStream(conn)
    for {
        select {
        // ...
        // 接收到来自 chan 的指令，分发到存储引擎层执行
        case droplet := <-stream:
            if droplet == nil || droplet.Terminated() {
                return // 加载完成
            }
            if err := h.handleDroplet(ctx, conn, droplet); err != nil {
                // ...
            }
        }
    }
}
```

#### 3 goredis AOF 持久化实现

介绍完上游模块是如何对 `persister` 展开使用后，下面就打开 `persister` 模块的黑盒子，解释 `goredis` 中针对 AOF 持久化机制的实现细节。

##### 3.1 类定义

![图片](https://img-blog.csdnimg.cn/img_convert/06758652618037380128795085351235.png)

`aofPersister` 是对 `persister` 的具体实现，其实例在启动时，会伴生启动一个异步守护协程，持续接收来自数据执行层投递的增量指令并为之完成持久化操作。

有关 AOF 持久化模块的类定义如下：

- 通过 `ctx` 实现守护协程的生命周期控制。
- 通过 `buffer` channel 接收来自数据执行层 `executor` 投递的写指令。
- 通过 `aofFile` 文件持久化记录写指令内容。
- 此外基于一系列配置参数，定义 AOF 落盘和指令重写的策略。

```go
// aof 持久化模块实现类
type aofPersister struct {
    ctx    context.Context
    cancel context.CancelFunc

    buffer                 chan [][]byte
    aofFile                *os.File
    aofFileName            string
    appendFsync            appendSyncStrategy
    autoAofRewriteAfterCmd int64
    aofCounter             atomic.Int64
    mu                     sync.Mutex
    once                   sync.Once
}
```

##### 3.2 持久化策略

与 Redis AOF 机制相对应，在 `goredis` 的实现中，将 AOF 持久化策略同样划分为三种等级：

```go
// aof 持久化等级 always | everysec | no
type appendSyncStrategy string

const (
    alwaysAppendSyncStrategy   appendSyncStrategy = "always"   // 每条指令都进行持久化落盘
    everysecAppendSyncStrategy appendSyncStrategy = "everysec" // 每秒批量执行一次持久化落盘
    noAppendSyncStrategy       appendSyncStrategy = "no"       // 不主动进行指令的持久化落盘，由设备自行决定落盘节奏
)
```

在构造 `persister` 实例时，会读取 `redis.conf` 中的配置信息，决定是否启用 AOF 持久化策略，以及对应的持久化策略级别：

```go
// thinker 为收拢了持久化配置参数的 interface
func NewPersister(thinker Thinker) (handler.Persister, error) {
    // 不启用 aof 持久化
    if !thinker.AppendOnly() {
        return newFakePersister(nil), nil
    }
     // 启用 aof 持久化
    return newAofPersister(thinker)
}

// 构造 aof 持久化模块
func newAofPersister(thinker Thinker) (handler.Persister, error) {
    // ...
    // 构造持久化模块实例
    a := aofPersister{
        // ...
    }

    // 设置 aof 持久化策略级别
    switch thinker.AppendFsync() {
    case string(alwaysAppendSyncStrategy):
        a.appendFsync = alwaysAppendSyncStrategy
    case string(everysecAppendSyncStrategy):
        a.appendFsync = everysecAppendSyncStrategy
    default:
        a.appendFsync = noAppendSyncStrategy // 默认策略
    }

    // 启动持久化模块常驻运行协程
    pool.Submit(a.run)
    return &a, nil
}
```

##### 3.3 核心方法

![图片](https://img-blog.csdnimg.cn/img_convert/06758652618037380128795085351235.png)

`aofPersister` 的伴生协程会持续运行 `run` 方法，通过 `for-select` 的运行框架，持续监听 `buffer` channel，接收来自上游数据执行层 `executor` 投递的增量写指令，并对其完成持久化操作：

```go
func (a *aofPersister) run() {
    if a.appendFsync == everysecAppendSyncStrategy {
        pool.Submit(a.fsyncEverySecond)
    }

    for {
        select {
        case <-a.ctx.Done():
            return
        case cmd := <-a.buffer:
            a.writeAof(cmd)
            a.aofTick()
        }
    }
}
```

当 `aofPersister` 从 channel 中获取到拟持久化的写指令后，则会将其格式化成 Multi-Bulk Reply 的形式，并调用 `file.Write` 方法。如果持久化级别为 `always`，则会立即执行一次 `fsync` 操作，确保指令当即被落盘。

```go
func (a *aofPersister) writeAof(cmd [][]byte) {
    a.mu.Lock()
    defer a.mu.Unlock()

    persistCmd := handler.NewMultiBulkReply(cmd)
    if _, err := a.aofFile.Write(persistCmd.ToBytes()); err != nil {
        // log
        return
    }

    if a.appendFsync != alwaysAppendSyncStrategy {
        return
    }

    if err := a.fsyncLocked(); err != nil {
        // log
    }
}
```

而在 `goredis` 重启时，`handler` 会调用 `persister` 的 `Reloader` 方法，将 AOF 文件包装成一个 `io.ReadCloser`，用于后续的数据加载：

```go
func (a *aofPersister) Reloader() (io.ReadCloser, error) {
    file, err := os.Open(a.aofFileName)
    if err != nil {
        if os.IsNotExist(err) {
            return nil, nil
        }
        return nil, err
    }
    _, _ = file.Seek(0, io.SeekStart)
    return file, nil
}
```

#### 4 goredis AOF 指令重写流程

AOF 的一大弊端是存在数据冗余。`goredis` 建立了一套 AOF 指令“瘦身”机制，其思路是，针对内存数据“拷贝”出一份副本，然后从结果出发将其映射成一条简单的 AOF 指令，从而实现冗余指令的去重。

##### 4.1 重写流程启动

![图片](https://img-blog.csdnimg.cn/img_convert/06758652618037380128795085351235.png)

在 `aofPersister` 每持久化一条指令时，都会通过计数器进行记录，当执行的指令数量达到一定阈值后，就会异步启动一次 AOF 重写流程：

```go
func (a *aofPersister) aofTick() {
    if a.autoAofRewriteAfterCmd <= 1 {
        return
    }

    if ticked := a.aofCounter.Add(1); ticked < a.autoAofRewriteAfterCmd {
        return
    }

    _ = a.aofCounter.Add(-a.autoAofRewriteAfterCmd)
    pool.Submit(func() {
        if err := a.rewriteAOF(); err != nil {
            // log
        }
    })
}
```

##### 4.2 重写详细步骤

![图片](https://img-blog.csdnimg.cn/img_convert/06758652618037380128795085351235.png)

AOF 指令重写流程会分为三个核心步骤：

- **重写前准备**：记录此刻 AOF 文件的大小，并拷贝生成一份 AOF 临时文件副本（此阶段需要加锁，短暂停止 AOF 持久化流程）。
- **重写进行时**：拷贝生成一份内存数据副本，并映射成最直观的 AOF 指令落在 AOF 临时文件副本中（此流程可以与 AOF 持久化流程并发进行）。
- **重写收尾**： 将原 AOF 文件后继部分（在重写进行时，会有新的指令在并发地执行持久化操作）追加到 AOF 临时文件副本中，然后使用临时文件副本覆盖原 AOF 文件（此阶段需要加锁，暂停 AOF 持久化流程）。

```go
func (a *aofPersister) rewriteAOF() error {
    // 1. 重写前处理. 需要短暂加锁
    tmpFile, fileSize, err := a.startRewrite()
    if err != nil {
        return err
    }
    defer tmpFile.Close()

    // 2. aof 指令重写. 与主流程并发执行
    if err = a.doRewrite(tmpFile, fileSize); err != nil {
        return err
    }

    // 3. 完成重写. 需要短暂加锁
    return a.endRewrite(tmpFile, fileSize)
}
```

**重写进行时**：读取重写开始前的所有 AOF 指令，构造一份内存数据副本，然后遍历该副本，将其中的 KV 数据转换成最直接的 `SET`、`HSET` 等指令写入临时文件。

```go
func (a *aofPersister) doRewrite(tmpFile *os.File, fileSize int64) error {
    forkedDB, err := a.forkDB(fileSize)
    if err != nil {
        return err
    }

    // 将 db 数据转为 aof cmd
    forkedDB.ForEach(func(key string, adapter database.CmdAdapter, expireAt *time.Time) {
        _, _ = tmpFile.Write(handler.NewMultiBulkReply(adapter.ToCmd()).ToBytes())
        if expireAt != nil {
            expireCmd := [][]byte{[]byte(database.CmdTypeExpireAt), []byte(key), []byte(lib.TimeSecondFormat(*expireAt))}
            _, _ = tmpFile.Write(handler.NewMultiBulkReply(expireCmd).ToBytes())
        }
    })
    return nil
}
```

**重写收尾**：加锁，将重写期间新产生的 AOF 指令追加到临时文件，最后用临时文件替换旧的 AOF 文件，完成重写。

```go
func (a *aofPersister) endRewrite(tmpFile *os.File, fileSize int64) error {
    a.mu.Lock()
    defer a.mu.Unlock()

    // copy commands executed during rewriting to tmpFile
    src, err := os.Open(a.aofFileName)
    if err != nil {
        return err
    }
    defer src.Close()

    if _, err = src.Seek(fileSize, 0); err != nil {
        return err
    }

    // 把老的 aof 文件中后续内容 copy 到 tmp 中
    if _, err = io.Copy(tmpFile, src); err != nil {
        return err
    }

    // 关闭老的 aof 文件，准备废弃
    _ = a.aofFile.Close()
    // 重命名 tmp 文件，作为新的 aof 文件
    if err := os.Rename(tmpFile.Name(), a.aofFileName); err != nil {
        // log
    }

    // 重新开启
    aofFile, err := os.OpenFile(a.aofFileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
    if err != nil {
        panic(err)
    }
    a.aofFile = aofFile
    return nil
}
```

#### 5 展望

至此，整个【基于 go 实现 redis】系列内容全部完成，在此对本系列内容做个总结：

- **基于 go 实现 redis 之主干框架（已完成）**： 在宏观视角下纵览 goredis 整体架构，梳理各模块间的关联性。
- **基于 go 实现 redis 之指令分发（已完成）**： 聚焦介绍 goredis 服务端如何启动和运行，并在接收客户端请求后实现指令协议的解析和分发。
- **基于 go 实现 redis 之存储引擎（已完成）**： 聚焦介绍数据存储层中单协程无锁化执行框架，各类基本数据类型的底层实现细节，以及过期数据的惰性和定期回收机制。
- **基于 go 实现 redis 之数据持久化（已完成）**： 介绍 goredis 关于 AOF 持久化机制的实现以及有关于 AOF 重写策略的执行细节。
