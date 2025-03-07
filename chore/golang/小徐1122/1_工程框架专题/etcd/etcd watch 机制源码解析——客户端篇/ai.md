# etcd watch 机制源码解析——客户端篇

https://mp.weixin.qq.com/s/2TEgbOoX36PwSWzbKq0Qsg

etcd 的 Watch 机制是其核心功能之一，允许客户端监听特定键或键范围的变更。本文将从客户端角度深入解析 Watch 机制的实现细节，结合源码详细分析其架构、核心数据结构和关键流程。

---

### 一、背景与核心概念

#### 1.1 Watch 机制的作用

- **分布式锁**：通过监听锁键的删除事件，避免无效轮询。
- **配置中心**：实时获取配置变更，如服务发现中的节点变化。
- **数据同步**：监听数据变更，实现跨系统的状态同步。

#### 1.2 客户端核心架构

客户端通过 gRPC 长连接与 etcd 服务端通信，关键组件包括：

- **Client**：封装与 etcd 交互的接口。
- **Watcher**：管理 Watch 请求和事件分发。
- **watchGrpcStream**：处理 gRPC 流通信。
- **watcherStream**：单个 Watch 的事件处理单元。

---

### 二、核心数据结构解析

#### 2.1 Client

```go
type Client struct {
    Watcher
    conn *grpc.ClientConn // gRPC 连接
}
```

- **Watcher**：负责 Watch 请求的创建和管理。
- **conn**：维护与 etcd 服务端的 gRPC 长连接。

#### 2.2 Watcher

```go
type watcher struct {
    remote pb.WatchClient          // gRPC 客户端
    mu     sync.Mutex              // 保护 streams
    streams map[string]*watchGrpcStream // 按上下文分组的流
}
```

- **streams**：管理多个 `watchGrpcStream`，按上下文隔离（如不同租户）。

#### 2.3 watchGrpcStream

```go
type watchGrpcStream struct {
    owner      *watcher
    remote     pb.WatchClient
    substreams map[int64]*watcherStream // 按 Watch ID 管理的子流
    reqc       chan watchStreamRequest  // 接收 Watch 请求
    respc      chan *pb.WatchResponse   // 接收服务端响应
}
```

- **substreams**：每个活跃的 Watch 对应一个 `watcherStream`。
- **reqc/respc**：异步处理请求和响应的通道。

#### 2.4 watcherStream

```go
type watcherStream struct {
    initReq watchRequest          // 初始请求参数
    outc    chan WatchResponse    // 推送事件给上层
    recvc   chan *WatchResponse   // 接收下层事件
    buf     []*WatchResponse      // 事件缓冲区
}
```

- **outc**：应用层通过此通道接收事件。
- **buf**：缓存未处理事件，确保顺序性。

---

### 三、创建 Watch 的完整流程

#### 3.1 应用层入口：endpointManager.NewWatchChannel

```go
func (m *endpointManager) NewWatchChannel(ctx context.Context) (WatchChannel, error) {
    // 1. 获取当前数据快照
    resp, _ := m.client.Get(ctx, key, clientv3.WithPrefix())
    initUpdates := parseInitUpdates(resp.Kvs) // 解析历史数据

    // 2. 创建事件通道
    upch := make(chan []*Update, 1)
    if len(initUpdates) > 0 {
        upch <- initUpdates // 推送初始数据
    }

    // 3. 启动后台监听协程
    go m.watch(ctx, resp.Header.Revision+1, upch)
    return upch, nil
}
```

- **初始数据获取**：确保应用获取 Watch 开始前的已有数据。
- **Revision 控制**：从当前版本+1 开始监听，避免遗漏变更。

#### 3.2 监听协程：endpointManager.watch

```go
func (m *endpointManager) watch(ctx context.Context, rev int64, upch chan []*Update) {
    wch := m.client.Watch(ctx, key, clientv3.WithRev(rev), clientv3.WithPrefix())
    for {
        select {
        case wresp := <-wch:
            deltaUps := parseEvents(wresp.Events) // 解析事件
            upch <- deltaUps // 推送变更
        case <-ctx.Done():
            return
        }
    }
}
```

- **事件解析**：将 etcd 事件转换为应用层理解的 `Update` 结构。
- **非阻塞推送**：通过缓冲通道 `upch` 避免阻塞。

#### 3.3 Watcher 核心逻辑：watcher.Watch

```go
func (w *watcher) Watch(ctx context.Context, key string, opts ...OpOption) WatchChan {
    // 1. 创建 Watch 请求
    wr := &watchRequest{
        key: key,
        retc: make(chan chan WatchResponse, 1),
    }

    // 2. 获取或创建 gRPC 流
    w.mu.Lock()
    wgs := w.streams[ctxKey]
    if wgs == nil {
        wgs = w.newWatcherGrpcStream(ctx)
        w.streams[ctxKey] = wgs
    }
    w.mu.Unlock()

    // 3. 发送请求并等待响应
    wgs.reqc <- wr
    return <-wr.retc // 返回事件通道
}
```

- **流复用**：相同上下文的 Watch 复用同一 gRPC 流，减少连接开销。
- **异步处理**：通过 `retc` 通道返回事件通道，实现非阻塞。

#### 3.4 gRPC 流管理：watchGrpcStream.run

```go
func (w *watchGrpcStream) run() {
    // 1. 建立 gRPC 连接
    wc, _ := w.newWatchClient()
    go w.serveWatchClient(wc) // 启动响应接收协程

    // 2. 主循环处理请求和响应
    for {
        select {
        case req := <-w.reqc:
            // 处理创建 Watch 请求
            ws := newWatcherStream(req)
            w.resuming = append(w.resuming, ws)
            wc.Send(ws.initReq.toPB())

        case pbresp := <-w.respc:
            // 分发响应到对应子流
            w.dispatchEvent(pbresp)
        }
    }
}
```

- **请求队列**：通过 `reqc` 缓冲处理多个 Watch 创建请求。
- **事件分发**：根据 Watch ID 将事件路由到正确的 `watcherStream`。

#### 3.5 子流处理：watchGrpcStream.serveSubstream

```go
func (w *watchGrpcStream) serveSubstream(ws *watcherStream) {
    for {
        select {
        case wr := <-ws.recvc:
            if wr.Created {
                ws.initReq.retc <- ws.outc // 返回事件通道给上层
            }
            ws.buf = append(ws.buf, wr) // 缓冲事件
        case ws.outc <- *ws.buf[0]:
            ws.buf = ws.buf[1:] // 弹出已发送事件
        }
    }
}
```

- **双缓冲机制**：`recvc` 接收事件，`buf` 缓存确保顺序，`outc` 非阻塞推送。

---

### 四、Watch 回调事件的处理

#### 4.1 接收响应：serveWatchClient

```go
func (w *watchGrpcStream) serveWatchClient(wc pb.Watch_WatchClient) {
    for {
        resp, _ := wc.Recv() // 接收 gRPC 响应
        w.respc <- resp      // 转发到主处理循环
    }
}
```

- **持续监听**：通过 gRPC 流式接口持续接收服务端推送。

#### 4.2 事件路由：dispatchEvent

```go
func (w *watchGrpcStream) dispatchEvent(pbresp *pb.WatchResponse) bool {
    // 1. 查找对应的子流
    ws, ok := w.substreams[pbresp.WatchId]
    if !ok {
        return false
    }

    // 2. 转换事件格式并发送
    wr := convertResponse(pbresp)
    select {
    case ws.recvc <- wr: // 非阻塞发送
        return true
    default:
        return false
    }
}
```

- **精准路由**：通过 `WatchId` 关联到具体的 `watcherStream`。

#### 4.3 最终推送：watcherStream 到应用层

```go
// 在 serveSubstream 中处理事件缓冲
case ws.outc <- *ws.buf[0]:
    ws.buf = ws.buf[1:]
```

- **顺序保证**：通过先进先出的缓冲区确保事件顺序与 etcd 日志一致。

---

### 五、关键设计思想总结

1. **连接复用**：通过 `watchGrpcStream` 复用 gRPC 流，降低连接开销。
2. **异步非阻塞**：全链路采用 Channel 实现异步处理，避免阻塞主流程。
3. **精准路由**：通过 `WatchId` 实现事件到具体 Watch 的精准分发。
4. **顺序保证**：缓冲区确保事件顺序，与 etcd 的 Revision 机制一致。
5. **错误恢复**：隐含的重连机制（服务端篇详述）保证断线后自动恢复。

---

### 六、性能优化点

1. **批量事件处理**：合并多个事件到单个响应，减少网络开销。
2. **智能缓冲区**：动态调整缓冲区大小，平衡内存和吞吐量。
3. **流控机制**：通过 `WithFragment` 选项控制大值分片传输。
4. **心跳检测**：定期 ProgressNotify 确认连接健康状态。

通过深入源码分析，可以看出 etcd 客户端 Watch 机制通过精巧的并发模型和高效的事件路由机制，在保证强一致性的同时实现了高性能的事件推送能力。服务端篇将深入探讨 etcd 如何高效管理海量 Watch 连接及事件分发。
