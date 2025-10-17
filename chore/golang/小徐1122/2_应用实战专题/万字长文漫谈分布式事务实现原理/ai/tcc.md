好的，这是对《从零到一搭建 TCC 分布式事务框架》这篇文章的深入详细讲解。

这篇文章是上一篇理论文章的实践篇，详细介绍了如何使用 Go 语言从零开始构建一个名为 `gotcc` 的 TCC 分布式事务框架。

### 1. 架构设计

文章首先回顾了 TCC（Try-Confirm/Cancel）作为两阶段提交（2PC）协议的实现，并明确了框架的设计思路：**分离通用逻辑与用户自定义实现**。

- **TCC SDK (`gotcc`) 负责的通用逻辑**：

  - **`TXManager` (事务协调器)**：实现核心的 2PC 流程（Try -> Confirm/Cancel）、TCC 组件的注册管理、以及用于保证最终一致性的异步轮询任务。
  - **定义接口规范**：
    - `TCCComponent` 接口：定义了 TCC 组件必须实现的 `Try`, `Confirm`, `Cancel` 方法。
    - `TXStore` 接口：定义了事务日志存储模块必须实现的 CRUD 和分布式锁等方法。

- **用户需要自行实现的部分**：
  - **具体的 `TCCComponent`**：根据业务需求，实现每个参与方（如订单服务、库存服务）的 `Try`, `Confirm`, `Cancel` 逻辑。
  - **具体的 `TXStore`**：实现事务日志的持久化存储，可以选择 MySQL、PostgreSQL 等数据库，并实现分布式锁（如基于 Redis）。

#### 核心角色

1.  **`TXManager`**：框架的核心，作为统一入口，协调整个事务流程。
2.  **`TCCComponent`**：用户实现的业务逻辑单元，在 `TXManager` 启动时被注册。
3.  **`TXStore`**：用户实现的事务日志存储模块，用于记录事务状态，是实现异步轮询和故障恢复的基础。
4.  **`RegistryCenter`**：`TXManager` 内置的组件，用于管理已注册的 `TCCComponent`。

### 2. `TXManager` 核心源码讲解

这部分深入剖析了 `gotcc` 框架的核心代码。

#### 核心类定义

- **`TXManager`**: 事务协调器，包含 `TXStore`、`registryCenter` 和配置项。
- **`RegistryCenter`**: 一个带读写锁 (`sync.RWMutex`) 的 `map`，用于并发安全地存储和查询 TCC 组件。
- **`TXStore` (interface)**: 定义了事务存储的契约，关键方法包括：
  - `CreateTX`: 创建事务记录并返回全局唯一的 `txID`。
  - `TXUpdate`: 更新某个组件的 `Try` 状态。
  - `TXSubmit`: 将事务状态提交为最终的成功或失败。
  - `GetHangingTXs`: 获取所有处于中间状态（未完成）的事务，供轮询任务使用。
  - `Lock`/`Unlock`: 分布式锁，防止多个 `TXManager` 节点同时执行轮询任务。

#### 事务主流程 (`Transaction` 方法)

这是用户启动一个分布式事务的入口。

1.  **获取组件**：根据请求参数，从 `RegistryCenter` 查找对应的 TCC 组件实例。
2.  **创建事务日志**：调用 `txStore.CreateTX`，在数据库中创建一条状态为 `hanging` 的事务记录，并获得 `txID`。
3.  **执行两阶段提交 (`twoPhaseCommit`)**：这是框架的精髓所在。
    - **第一阶段 (Try)**：
      - 使用 `goroutine` 和 `sync.WaitGroup` **并发**执行所有组件的 `Try` 方法。
      - 使用一个 `chan error` 来接收执行结果。任何一个 `Try` 失败，都会立即向 channel 发送错误。
      - 主流程阻塞等待 channel。一旦收到错误，立即调用 `context.CancelFunc` 来终止其他正在执行的 `Try` 操作，实现快速失败。
    - **第二阶段 (Confirm/Cancel)**：
      - `Try` 阶段结束后，事务的最终成败已定。
      - 框架会**异步地** (`go t.advanceProgressByTXID(txID)`) 启动第二阶段。
      - 这么做是因为第二阶段有异步轮询任务作为“安全网”来兜底，即使本次执行失败，轮询任务也会确保它最终被执行。

#### 异步轮询流程 (`run` 方法)

这是保证框架鲁棒性和数据最终一致性的关键。

1.  **启动时机**：`TXManager` 初始化时，一个后台 `goroutine` 就会启动该轮询任务。
2.  **核心逻辑**：
    - 使用 `time.After` 定时触发。
    - **获取分布式锁**：在处理前，调用 `txStore.Lock()`，确保在分布式部署时，只有一个 `TXManager` 节点在执行轮询，避免重复处理。
    - **获取待处理事务**：调用 `txStore.GetHangingTXs()` 从数据库获取所有未完成的事务。
    - **批量推进进度**：并发地为每个 `hanging` 状态的事务执行 `advanceProgress` 方法，即根据其 `Try` 阶段的结果来调用 `Confirm` 或 `Cancel`，并更新事务日志为最终状态。
    - **释放锁**：处理完毕后，调用 `txStore.Unlock()`。
    - **退避策略**：如果轮询过程中发生错误，会动态增加下一次轮询的间隔时间（退避），避免在系统故障时造成过大压力。

### 3. `gotcc` 使用案例讲解

文章最后提供了一个完整的示例，演示了如何使用 `gotcc` 框架。

1.  **实现 `TCCComponent`**：

    - 创建了一个 `MockComponent`，其 `Try`, `Confirm`, `Cancel` 方法都依赖 Redis。
    - **`Try`**：使用 Redis 分布式锁保证操作原子性，通过检查 Redis 中的键来实现幂等性，并将业务数据状态置为“冻结”。
    - **`Confirm`**：同样加锁和检查状态，将“冻结”状态改为“成功”。
    - **`Cancel`**：同样加锁和检查状态，删除“冻结”记录或将状态置为“取消”。支持**空回滚**。

2.  **实现 `TXStore`**：

    - 创建了一个 `MockTXStore`，使用 **MySQL** 存储事务日志，使用 **Redis** 实现分布式锁。
    - 定义了 `TXRecordPO` (持久化对象) 来映射数据库表结构。
    - 实现了 `TXStore` 接口的所有方法，例如 `CreateTX` 会在 `tx_record` 表中插入一条记录。在更新操作中，使用了 `SELECT ... FOR UPDATE` 来实现行级锁，保证并发安全。

3.  **使用示例**：
    - 初始化 MySQL 和 Redis 客户端。
    - 创建 `TCCComponent` 和 `TXStore` 的实例。
    - 创建 `TXManager` 实例，并将 `TXStore` 注入。
    - 向 `TXManager` 注册所有 `TCCComponent`。
    - 调用 `txManager.Transaction()` 方法，传入要执行的组件列表和参数，启动分布式事务。

通过这个实践，文章清晰地展示了框架设计与实际业务逻辑如何解耦并协同工作。
