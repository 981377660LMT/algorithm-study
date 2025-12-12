这份文档详细介绍了 **LangGraph** 的底层运行时——**Pregel**。它负责管理 LangGraph 应用程序的执行逻辑。

以下是核心概念和机制的总结：

### 1. 核心架构
Pregel 基于 Google 的同名图算法（Bulk Synchronous Parallel 模型），将应用程序建模为 **Actors**（执行者）和 **Channels**（通道）的组合。

*   **Actors (PregelNode)**:
    *   相当于图中的节点。
    *   订阅通道，读取数据，执行逻辑，并将结果写入通道。
    *   实现了 LangChain 的 `Runnable` 接口。
*   **Channels**:
    *   用于在 Actors 之间传递数据或存储状态。
    *   **`LastValue`**: 默认类型，存储最后发送的值（常用于输入/输出）。
    *   **`Topic`**: 发布/订阅模式，用于在步骤间传递多个值或累积输出。
    *   **`BinaryOperatorAggregate`**: 使用归约函数（Reducer）更新值（例如累加、拼接）。

### 2. 执行周期 (The Step)
Pregel 将执行过程组织为一系列的“步骤”（Steps），每个步骤包含三个阶段：

1.  **Plan (计划)**: 确定当前步骤需要执行哪些 Actors（基于它们订阅的通道是否有更新）。
2.  **Execution (执行)**: 并行执行所有被选中的 Actors。**注意**：在此阶段，Actors 写入的数据对其他 Actors 不可见。
3.  **Update (更新)**: 步骤结束后，统一将 Actors 写入的数据应用到 Channels 中。

这个循环会一直重复，直到没有 Actors 需要执行或达到最大步数限制。

### 3. 使用方式
虽然开发者通常使用高级 API，但 LangGraph 允许直接操作 Pregel：

*   **高级 API (推荐)**:
    *   **StateGraph (Graph API)**: 通过定义节点和边来构建图，编译时自动生成 Pregel 实例。
    *   **Functional API**: 使用 `entrypoint` 装饰器定义工作流，底层同样编译为 Pregel 实例。
*   **低级 API (直接使用)**:
    *   可以直接实例化 `Pregel` 类，手动配置 `nodes`（节点逻辑）和 `channels`（数据通道）。

### 代码示例 (直接使用 Pregel)

以下是一个直接使用 `Pregel` 构建简单节点的示例：

```typescript
import { EphemeralValue } from "@langchain/langgraph/channels";
import { Pregel, NodeBuilder } from "@langchain/langgraph/pregel";

// 定义节点：订阅 "a"，将输入重复一次，写入 "b"
const node1 = new NodeBuilder()
  .subscribeOnly("a")
  .do((x: string) => x + x)
  .writeTo("b");

// 定义 Pregel 应用
const app = new Pregel({
  nodes: { node1 },
  channels: {
    a: new EphemeralValue<string>(),
    b: new EphemeralValue<string>(),
  },
  inputChannels: ["a"],
  outputChannels: ["b"],
});

// 运行
const result = await app.invoke({ a: "foo" });
console.log(result); // { b: 'foofoo' }
```