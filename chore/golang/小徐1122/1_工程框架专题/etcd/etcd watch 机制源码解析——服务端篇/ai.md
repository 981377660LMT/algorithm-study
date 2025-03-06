# etcd watch 机制源码解析——服务端篇

https://mp.weixin.qq.com/s/-Vxu7jQZ-7ID-4oUF_0Agg

### **etcd Watch 服务端机制核心解析**

---

#### **一、服务端架构**

**1. 整体分层结构**

- **serverWatchStream**：处理与客户端的 gRPC 长连接，负责接收请求和发送响应。
  - 包含读协程 `recvLoop`（处理客户端请求）和写协程 `sendLoop`（推送事件和响应）。
- **watchableStore**：底层存储 Watcher，基于内存管理，分为 `synced`、`unsynced` 和 `victims` 三部分。
  - **synced**：无需回溯历史事件的 Watcher。
  - **unsynced**：需回溯历史事件的 Watcher。
  - **victims**：因通道容量不足暂存的事件。

**2. 事件触发流程**

- **数据变更切面**：当状态机写入数据时（Raft 日志提交后），触发 `notify` 操作。
- **事件生成**：将数据变更与 `synced` 组的 Watcher 匹配，生成事件批次。
- **事件传递**：通过 `watchStream` 通道将事件推送至 `sendLoop`，最终通过 gRPC 发送给客户端。

---

#### **二、核心数据结构**

**1. `serverWatchStream`**

- 管理客户端长连接，通过 `watchStream` 与底层 `watchableStore` 交互。
- 关键字段：
  - `gRPCStream`：与客户端的通信通道。
  - `ctrlStream`：内部协程间传递控制响应的通道。

**2. `watchableStore`**

- 存储 Watcher 的核心模块：
  - `synced` 和 `unsynced` 组：基于 `watcherGroup` 实现。
  - `victims`：暂存发送失败的事件批次（`watcherBatch`）。
- 关键方法：
  - `notify()`：生成事件并尝试发送。
  - `syncWatchersLoop()` 和 `syncVictimsLoop()`：异步处理未同步 Watcher 和暂存事件。

**3. `watcherGroup`**

- 通过红黑树（`IntervalTree`）和哈希表（`watcherSetByKey`）高效管理 Watcher：
  - **单 Key 监听**：哈希表快速查找。
  - **Range 监听**：红黑树支持区间查询（`Stab` 操作）。

**4. `watcher`**

- 监听器的核心属性：
  - `minRev`：记录事件处理进度，避免重复通知。
  - `ch`：事件推送通道（绑定至 `watchStream`）。
  - `victim` 标记：标识事件是否进入暂存。

---

#### **三、关键流程解析**

**1. 创建 Watcher**

- **客户端请求**：通过 gRPC 发送 `WatchRequest`。
- **服务端处理**：
  - `recvLoop` 接收请求，调用 `watchStream.Watch()`。
  - `watchableStore` 根据 `startRev` 决定将 Watcher 放入 `synced` 或 `unsynced` 组。
  - 响应通过 `ctrlStream` 发送至客户端。

**2. 事件回调**

- **触发时机**：状态机写入数据后调用 `notify()`。
- **事件匹配**：通过 `newWatcherBatch` 将变更事件与 `synced` 组的 Watcher 匹配。
- **发送逻辑**：
  - 成功：直接通过 `watcher.ch` 推送。
  - 失败：事件暂存至 `victims`，后续由 `syncVictimsLoop` 重试。

**3. 数据同步与补偿**

- **unsynced 组处理**：`syncWatchersLoop` 定期拉取历史数据，填充未同步事件。
- **victims 处理**：`syncVictimsLoop` 每隔 10ms 重试发送暂存事件，失败则保留。

---

#### **四、设计亮点**

**1. 内存高效管理**

- Watcher 完全基于内存存储，通过红黑树和哈希表实现高效查询。
- `victims` 机制避免高并发下通道阻塞，牺牲部分内存换取系统稳定性。

**2. 事件去重与进度控制**

- `minRev` 字段确保仅推送新于当前进度的事件。
- 过滤器（`FilterFunc`）支持自定义事件过滤，减少无效传输。

**3. 异步化处理**

- 读写协程分离（`recvLoop` 和 `sendLoop`），避免阻塞。
- 后台协程（`syncWatchersLoop` 和 `syncVictimsLoop`）异步补偿异常场景。

---

#### **五、注意事项**

**1. Watcher 的易失性**

- Watcher 存储在内存中，节点重启或连接断开会导致 Watcher 丢失。
- **解决方案**：客户端需实现重连和 Watcher 重建逻辑。

**2. 网络分区问题**

- 若客户端连接的节点与集群多数派分区，Watcher 可能停滞。
- **规避方案**：使用 `WithRequireLeader` 选项，确保连接至 Leader 节点。

**3. 性能影响**

- 大量 Watcher 或大范围监听（如全 Key 监听）会显著增加内存和 CPU 开销。
- **优化建议**：合理设计监听范围，避免全量监听。

---

#### **六、总结**

etcd 的 Watch 服务端机制通过多层协作实现高效事件推送：

1. **连接管理**：基于 gRPC 长连接，通过分离的读写协程处理高并发。
2. **事件生成**：在数据提交切面触发事件匹配，减少延迟。
3. **异常处理**：通过 `victims` 和后台协程实现故障自愈。
4. **内存优化**：数据结构设计兼顾查询效率与内存占用。

源码设计体现了分布式系统**可靠性**与**性能**的平衡，是理解 etcd 高可用特性的关键切入点。
