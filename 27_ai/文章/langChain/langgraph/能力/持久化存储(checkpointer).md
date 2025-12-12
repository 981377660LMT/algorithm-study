这份文档详细介绍了 **LangGraph** 中的 **Persistence (持久化)** 机制。
持久化是 LangGraph 实现记忆、人机协同、时间旅行（Time Travel）和容错能力的基础。

以下是核心概念的总结：

### 1. 核心组件：Checkpointer (检查点保存器)
LangGraph 通过 **Checkpointer** 在图执行的每一步（Super-step）自动保存状态快照。
*   **功能**：保存图的状态（State）、配置（Config）和元数据（Metadata）。
*   **实现**：
    *   `MemorySaver`: 内存存储，仅用于测试/开发。
    *   `SqliteSaver`: 基于 SQLite，适合本地开发。
    *   `PostgresSaver`: 基于 Postgres，适合生产环境。
*   **注意**：如果你使用 LangGraph Cloud (Agent Server)，它会自动处理持久化，无需手动配置。

### 2. 核心概念：Thread (线程)
*   **定义**：`thread_id` 是持久化的唯一标识符。它将一系列的执行步骤串联起来，形成一个会话。
*   **作用**：
    *   **记忆**：通过同一个 `thread_id` 调用图，Agent 可以“记住”之前的对话上下文。
    *   **恢复**：如果程序中断，可以通过 `thread_id` 恢复状态并继续执行。
*   **用法**：
    ```typescript
    const config = { configurable: { thread_id: "conversation-1" } };
    await graph.invoke({ input: "hi" }, config);
    ```

### 3. 核心概念：Checkpoint (检查点)
*   **定义**：某个特定时间点的图状态快照。
*   **包含内容**：
    *   `values`: 当前状态的值。
    *   `next`: 下一步要执行的节点。
    *   `tasks`: 待执行的任务（包含错误信息）。
*   **获取状态**：
    *   `graph.getState(config)`: 获取最新状态。
    *   `graph.getStateHistory(config)`: 获取历史状态列表。

### 4. 高级功能

#### A. 时间旅行 (Time Travel) & 重放 (Replay)
你可以通过指定 `checkpoint_id` 来“回到过去”。
*   **重放**：从某个旧的检查点重新开始执行。LangGraph 会跳过该检查点之前已执行的步骤，只执行之后的步骤。
    ```typescript
    const config = {
      configurable: {
        thread_id: "1",
        checkpoint_id: "old-checkpoint-uuid", // 指定回到哪个时刻
      },
    };
    await graph.invoke(null, config);
    ```

#### B. 修改状态 (Update State)
你可以人为地修改图的当前状态，这对于人机协同（纠正 Agent 的错误）非常有用。
*   **方法**：`graph.updateState(config, values, asNode?)`
*   **行为**：修改状态就像是一个节点输出了新的值。如果状态字段定义了 Reducer，它会按 Reducer 逻辑合并；如果没有，则直接覆盖。

#### C. 跨线程记忆 (Memory Store)
Checkpointer 只能在同一个 Thread 内保存状态。如果你需要**跨 Thread**（例如跨对话）共享数据（如用户偏好），需要使用 **Store**。
*   **Store**：一个键值对存储，支持命名空间。
*   **语义搜索**：Store 支持配置 Embedding 模型，从而实现对记忆的语义搜索。
    ```typescript
    // 写入记忆
    await store.put(["user_1", "memories"], "mem_id", { food: "pizza" });
    
    // 搜索记忆
    const memories = await store.search(["user_1", "memories"], { query: "food" });
    ```

### 总结
在 LangGraph 中，**Checkpointer** 负责“短期记忆”和会话恢复（Thread 级别），而 **Store** 负责“长期记忆”和跨会话知识共享（Global 级别）。两者结合构成了完整的持久化方案。