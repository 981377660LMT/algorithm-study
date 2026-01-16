## Node

https://the-pocket.github.io/PocketFlow/core_abstraction/node.html

这是一个关于 **Pocket Flow 中 `Node`（节点）** 的深度解析。它是框架中最小的构建块，也是执行具体业务逻辑的核心单元。

### 1. 核心设计：三明治结构 (The Sandwich Pattern)

Pocket Flow 的 Node 强制采用了 `prep` -> `exec` -> `post` 的三段式执行生命周期。这种设计不仅仅是为了代码整洁，更是为了实现 **关注点分离 (Separation of Concerns)**。

```python
# 数据流向：
Shared Store --> [prep] --> (prep_res) --> [exec] --> (exec_res) --> [post] --> Shared Store
```

#### 🔍 详细拆解

1.  **`prep(shared)` - 准备阶段**

    - **核心职责**: **只读**。从全局 `shared` 字典中提取原料。
    - **最佳实践**: 包括查询数据库、读取文件、序列化对象。不要在这里做繁重的计算。
    - **输出**: 返回 `prep_res`。这是专门为了喂给下一步 `exec` 的干净数据。

2.  **`exec(prep_res)` - 执行阶段**

    - **核心职责**: **纯计算**。这是放置 LLM 调用、API 请求或核心算法的地方。
    - **关键隔离**: 请注意，`exec` **不接收** `shared` 字典。它只能看到 `prep` 传给它的 `prep_res`。
      - _为什么？_ 这使得 `exec` 函数变成了一个纯粹的、无副作用的函数（Stateless）。这极大地简化了单元测试，因为你不需要模拟整个 Context，只需传入特定的输入即可测试。
    - **幂等性 (Idempotency)**: 由于支持自动重试，这里的逻辑必须是幂等的（即执行一次和执行多次的效果应该相同，或者是安全的）。

3.  **`post(shared, prep_res, exec_res)` - 收尾阶段**
    - **核心职责**: **写回与路由**。
    - **双重任务**:
      1.  **更新状态**: 将 `exec_res` 处理后写入 `shared` 字典。
      2.  **决定方向**: 返回一个字符串（Action），告诉 Flow 下一步去哪个节点。如果未返回任何值，默认为 `"default"`。

---

### 2. 容错机制 (Fault Tolerance) & 自动重试

在 AI 应用开发中，网络波动和 LLM API 不稳定是常态。Pocket Flow 将错误处理内置到了 Node 基类中，开发者无需手动写 `try-catch` 循环。

#### ⚙️ 配置参数

在实例化节点时配置：

```python
my_node = SummarizeFile(max_retries=3, wait=10)
```

- **`max_retries`**: 最大尝试次数（默认 1）。
- **`wait`**: 每次失败后等待的秒数（默认 0）。这对处理 API Rate Limit（速率限制）非常有效。

#### 🔄 重试逻辑

当 `exec()` 抛出异常时：

1.  捕获异常。
2.  检查是否达到 `max_retries`。
3.  如果未达到：`time.sleep(wait)`，然后重试。
4.  如果已耗尽次数：进入 **Fallback **流程。

开发者可以通过 `self.cur_retry` 属性知道当前是第几次重试（从 0 开始计数）。

---

### 3. 优雅降级 (Graceful Fallback)

当所有重试都失败后，程序不应该直接崩溃（Crash）。你可以通过覆盖 `exec_fallback` 方法来定义“B 计划”。

```python
def exec_fallback(self, prep_res, exc):
    # 默认行为是重新抛出异常 (raise exc)
    # 这里的 exc 是最后一次失败捕获的异常对象

    # 自定义行为：返回一个兜底值
    logging.error(f"处理失败: {exc}")
    return "默认摘要：处理过程中发生错误，请稍后重试。"
```

一旦 `exec_fallback` 返回了值，该值会被视为正常的 `exec_res` 传递给 `post` 方法，流程得以继续运行，从而保证系统的稳定性。

---

### 4. 代码结构示例

这是一个标准的 Node 定义模板，请将其添加到你的笔记中：

```python
class MyWorkerNode(Node):
    def prep(self, shared):
        # 1. READ: 从共享存储获取数据
        print(f"Prop from shared: {shared.get('input_text')}")
        return shared.get("input_text", "")

    def exec(self, prep_res):
        # 2. COMPUTE: 纯计算/IO，可能会失败，自带重试
        # 注意：这里不能访问 shared
        if not prep_res:
             raise ValueError("No input data")
        # 模拟 API 调用
        result = call_external_api(prep_res)
        return result

    def exec_fallback(self, prep_res, exc):
        # 3. FALLBACK: 所有重试失败后的兜底
        print(f"Error handling data: {exc}")
        return {"status": "error", "data": None}

    def post(self, shared, prep_res, exec_res):
        # 4. WRITE & ROUTE: 更新共享存储并决定下一步
        shared["result"] = exec_res

        # 路由逻辑
        if exec_res.get("status") == "error":
            return "retry_later" # 返回特定 Action
        return "success"         # 返回成功 Action
```

这段代码清晰地展示了 Pocket Flow 的核心哲学：**数据获取、逻辑计算、状态更新严格分离，且内置健壮的错误处理。**

## Flow

这是关于 **Pocket Flow 中 `Flow`（流程）** 的深度解析。如果说 `Node` 是独立的工匠，那么 `Flow` 就是**指挥家**。

### 1. 核心定义：图的编排者

`Flow` 负责将多个 `Node` 连接成一个有向图（Graph），并管理它们之间的执行顺序。

- **机制**: 它完全依赖上一个节点的 `post()` 方法返回的 **Action 字符串** 来决定去往哪里。
- **终止**: 当当前节点执行完毕，且没有定义针对其返回 Action 的后续节点时，Flow 自动结束。

---

### 2. DSL 语法：像画图一样写代码

Pocket Flow 使用 Python 的运算符重载，创造了一套极其直观的领域特定语言（DSL）来定义拓扑结构。

#### A. 默认流转 (`>>`)

最常见的情况。当节点返回 `"default"`（或 `post` 没返回任何值）时触发。

```python
# 语义：做完 A，紧接着做 B
node_a >> node_b
```

#### B. 条件流转 (`-` + `>>`)

实现**分支（Branching）**逻辑的核心。根据 `post` 返回的具体字符串跳转。

```python
# 语义：如果 A 返回 "success" 去 B，返回 "fail" 去 C
node_a - "success" >> node_b
node_a - "fail"    >> node_c
```

#### C. 循环 (Looping)

Pocket Flow 没有专门的循环语法，**循环就是指向自己的流转**。

```python
# 语义：如果需要修改，这就形成了一个闭环，直到 approved
review - "needs_revision" >> revise
revise >> review
```

---

### 3. 分形架构：Flow 即 Node (Everything is a Node)

这是 Pocket Flow 最强大的抽象设计：**`Flow` 本身也是一个 `Node`**。

这意味着你可以像俄罗斯套娃一样，将一个小流程封装起来，作为一个单一节点嵌入到更大的流程中。

- **组合性**: 你可以构建 `PaymentFlow`、`InventoryFlow`、`ShippingFlow`，然后将它们串联成一个主 `OrderFlow`。
- **生命周期**:
  - `Flow` 作为节点运行时，也会执行 `prep()` 和 `post()`。
  - **区别**: 它**不会**执行 `exec()`。它的“执行”逻辑就是运行内部的子节点网络。

```python
# 1. 定义子流程
payment_chain = Flow(start=validate_payment)
inventory_chain = Flow(start=check_stock)

# 2. 将子流程作为节点串联
# 这让代码具有极高的可读性和模块化程度
master_flow = Flow(start=payment_chain)
payment_chain >> inventory_chain >> shipping_node
```

---

### 4. 运行机制：`node.run` vs `flow.run`

混淆这两者是初学者常见的错误：

| 方法                   | 作用范围                                                              | 适用场景                 |
| :--------------------- | :-------------------------------------------------------------------- | :----------------------- |
| **`node.run(shared)`** | **单点执行**。只跑这一个节点的 prep->exec->post，**不**触发后续节点。 | 单元测试、调试单个组件。 |
| **`flow.run(shared)`** | **链路执行**。从 `start` 节点开始，沿着 Action 路径一直跑到底。       | 生产环境、集成测试。     |

---

### 5. 快速示例

将这段代码加入笔记，它展示了分支和状态流转：

```python
# 1. 定义连接
# 审批流程：通过 -> 付款；驳回 -> 结束；需修改 -> 修改 -> 重新审批
review_node - "approved"       >> pay_node
review_node - "rejected"       >> end_node
review_node - "needs_revision" >> revise_node
revise_node >> review_node # 闭环

# 2. 实例化 Flow
# 指定入口节点
expense_flow = Flow(start=review_node)

# 3. 运行
# shared 是贯穿全程的上下文
expense_flow.run(shared_context)
```

**总结：Flow 是通过 Action 字符串编织 Node 的轻量级编排器，且支持无限嵌套。**

## Communication

这是关于 **Pocket Flow 中通信机制 (Communication)** 的深度解析。在 Pocket Flow 中，节点之间不直接传递对象，而是通过 **Shared Store（共享存储）** 和 **Params（参数）** 进行交互。

---

### 1. 核心理念：解耦与角色分工

Pocket Flow 严格区分了“全局持久状态”和“局部临时标识”，旨在实现**计算逻辑与数据架构的彻底分离**。

| 机制             | 作用域        | 生命周期             | 模拟对象       | 适用场景                       |
| :--------------- | :------------ | :------------------- | :------------- | :----------------------------- |
| **Shared Store** | 全局 (Global) | 贯穿整个 Flow 运行   | **堆 (Heap)**  | 大多数场景、业务数据结果       |
| **Params**       | 局部 (Local)  | 仅在当前节点运行周期 | **栈 (Stack)** | 批处理、任务标识符 (ID/文件名) |

---

### 2. Shared Store：核心通信中心

这是 Pocket Flow 最推荐的通信方式。它通常是一个简单的 Python 字典。

- **设计原则**: 应该先设计数据结构（Schema），再编写节点。
- **流转方式**:
  - `prep(shared)`: **读取**数据。
  - `post(shared, ...)`: **写入**结果。
- **优势**: 极大降低了节点间的耦合度。节点只需要知道去哪个 Key 读/写数据，而不必关心是谁产生的。

#### 示例

Shared Store 是一个全局的内存字典，它是节点间通信的“主战场”。

- **关注点分离**: 将数据架构（Data Schema）与计算逻辑分离。
- **透明调试**: 随时打印 `shared` 即可捕获整个应用的状态快照。

```python
class Summarize(Node):
    def prep(self, shared):
        return shared["raw_text"]  # 读取

    def post(self, shared, prep_res, exec_res):
        shared["summary"] = exec_res # 写入
        return "default"
```

---

### 3. Params：任务的“通行证”

Params 是局部、不可变的。它们主要服务于 **Batch（批处理）** 模式。

- **特性**:
  - **临时性**: 每次父级 Flow 调用节点时都会刷新。
  - **不可变性**: 在 `prep` -> `exec` -> `post` 周期内，Params 是只读的。
- **用途**: 作为执行任务的“索引”或“标识符”。
  - 例如：在处理一个文件夹的 100 份文件时，Params 存储当前正在处理的 `filename`。

#### 示例

Params 是由父级 Flow 传入节点的局部、临时配置。

- **适用场景**: 在批处理 (Batch) 或并行任务中区分不同的 Sub-task。
- **不可变性**: 在单个节点的生命周期内，Params 是只读的。

```python
class ProcessFile(Node):
    def prep(self, shared):
        # 使用 Params 作为 Key 去 Shared Store 查找特定数据
        target_file = self.params["filename"]
        return shared["files_content"][target_file]

    def post(self, shared, prep_res, exec_res):
        target_file = self.params["filename"]
        shared["results"][target_file] = exec_res
```

#### 💡 核心忠告

1. **优先使用 Shared Store**: 除非是 `Batch` 任务需要区分 ID，否则请将业务数据全部放入 Shared Store。
2. **避免在 Flow 中硬编码 Params**: 节点参数应由外层调用者或 Batch 逻辑动态注入。

**总结：Shared Store 负责“存什么”，Params 负责“我是谁/我处理哪条”，两者结合实现了极简且强大的并发与循环能力。**

## Batch

Pocket Flow 的批处理核心在于两种维度的循环：**BatchNode（节点内循环）**与 **BatchFlow（流程间循环）**。

### 1. BatchNode：单节点数据拆分

用于在一个节点内处理大量数据（如长文本分片）。

- **prep**: 将数据切分为可迭代对象（如 `list`）。
- **exec**: 对每个切片调用一次（执行逻辑）。
- **post**: 汇总所有 `exec` 的结果列表，更新 `shared` 存储。
- **一句话总结**：**节点内部的任务并行/迭代**，处理的是“数据分片”。

### 2. BatchFlow：流程级多次重放

用于多次运行同一个子图（Flow），每次配置不同。

- **prep**: 返回 **参数字典列表**（例如 `[{"id": 1}, {"id": 2}]`）。
- **self.params**: 子流程中的节点必须通过 `self.params` 获取输入，而不是 `shared`。
- **独立性**: 每次运行互不干扰，适合处理多个文件或独立请求。
- **一句话总结**：**整个流程的参数化循环**，处理的是“任务实例”。

### 核心区别对比

| 特性               | BatchNode                             | BatchFlow                          |
| :----------------- | :------------------------------------ | :--------------------------------- |
| **操作对象**       | 单个节点                              | 整个子 Flow / 图                   |
| **输入来源**       | `shared` 数据分片                     | `prep` 返回的 `params` 字典        |
| **子节点访问方式** | 方法参数                              | `self.params`                      |
| **主要用途**       | 解决单节点处理上限（如 Context 限制） | 解决重复任务的批处理（如批量读档） |

### 3. 多级嵌套 (Nested Batch)

`BatchFlow` 可以层层嵌套（如：目录 -> 文件 -> 块）。

- **参数合并**：子层级会自动继承并合并父层级的 `params`。
- **上下文保留**：最内层节点可以同时访问到所有父层注入的参数（如 `self.params["directory"]` 和 `self.params["filename"]`）。

## Async

Pocket Flow 的异步处理利用 Python 的 `asyncio` 来优化 I/O 密集型任务（如 LLM 调用、API 请求、文件读取）。

### 1. 核心组件

- **AsyncNode**：异步逻辑的载体。
  - `prep_async()`：异步获取数据（读文件、查库）。
  - `exec_async()`：异步核心任务（最常用于 **LLM 调用**）。
  - `post_async()`：异步后续操作（等待用户反馈、多 Agent 协作）。
- **AsyncFlow**：异步流程的容器。
  - **必须**使用 `AsyncFlow` 来运行 `AsyncNode`。
  - **向下兼容**：`AsyncFlow` 中可以包含普通的同步 `Node`。

### 2. 关键语法

- **运行方式**：使用 `await flow.run_async(shared)` 启动。
- **强制匹配**：`AsyncNode` 只能在 `AsyncFlow` 中运行；普通的 `Node` 在异步流中会被自动处理。

### 3. 一针见血的适用场景

- **高并发**：同时调用多个 LLM 接口。
- **等待型任务**：需要 `await` 用户输入或外部系统回调。
- **非阻塞 I/O**：在处理大数据流时，不阻塞主线程。

### 4. 代码本质 (对比)

| 节点类型      | 方法名                                   | 启动入口                 | 适用场景                     |
| :------------ | :--------------------------------------- | :----------------------- | :--------------------------- |
| **Node**      | `prep`, `exec`, `post`                   | `flow.run()`             | 简单逻辑、本地计算           |
| **AsyncNode** | `prep_async`, `exec_async`, `post_async` | `await flow.run_async()` | LLM 调用、API 交互、人机协同 |

**总结：** `AsyncNode` 是为了不让程序在等待 LLM 响应或用户点击时“卡死”，必须配合 `AsyncFlow` 和 `await` 使用。

## Parallel

Pocket Flow 的并行处理通过并发执行 `AsyncNode` 或 `AsyncFlow` 来压榨 I/O 效率（尤其是 LLM 调用）。

### 1. 核心机制

- **非阻塞并发**：利用 `asyncio.gather` 同时发起多个异步任务。
- **适用场景**：**I/O 密集型**（LLM 请求、API 调用、数据库查询）。
- **局限性**：受 Python GIL 限制，无法加速 **CPU 密集型**（如复杂计算）任务。

### 2. 两种并发模式

#### **AsyncParallelBatchNode（节点级并行）**

将一个节点内的多个任务并发化。

- **逻辑**：`prep_async` 返回列表 $\rightarrow$ **所有 `exec_async` 同时启动** $\rightarrow$ `post_async` 汇总。
- **场景**：单步骤在大数据量下的横向扩展（例如：同时总结 10 段文本）。

#### **AsyncParallelBatchFlow（流程级并行）**

将整个子流程（Flow）并发化。

- **逻辑**：`prep_async` 返回参数字典列表 $\rightarrow$ **多个子 Flow 实例同时运行**。
- **场景**：多步骤复杂任务的横向扩展（例如：同时下载、解析、总结 10 个不同的文件）。

### 3. 一针见血的注意事项

- **任务独立性**：被并行的任务之间严禁有数据依赖（必须能同时运行）。
- **频率限制 (Rate Limits)**：并发极易撑爆 LLM 供应商的 RPM/TPM 限制，需设置**信号量 (Semaphore)** 或使用**官方 Batch API**。
- **继承参数**：`AsyncParallelBatchFlow` 会像 `BatchFlow` 一样将参数传递给 `self.params` 以保持独立上下文。

### 4. 选型指南

| 需求                                                  | 推荐方案                            |
| :---------------------------------------------------- | :---------------------------------- |
| **同时启动 10 个 LLM 请求**                           | `AsyncParallelBatchNode`            |
| **同时处理 10 个独立文件（包含读图、OCR、总结多步）** | `AsyncParallelBatchFlow`            |
| **需要按顺序一步步处理**                              | 使用普通 `BatchNode` 或 `BatchFlow` |

**总结**：如果你觉得循环处理太慢，且任务之间互不干扰，就把 `Batch` 换成 `ParallelBatch`。
