## 函数式API
这份文档介绍了 **LangGraph Functional API**。这是一种允许你使用标准的编程结构（如 `if` 语句、`for` 循环、函数调用）来构建 Agent 工作流，同时又能享受 LangGraph 提供的持久化、流式传输和“人在回路” (Human-in-the-loop) 功能的方法。

相比于声明式的 Graph API，Functional API 更适合习惯命令式编程或希望对现有代码进行最小改动的场景。

以下是核心概念和使用指南的总结：

### 1. 两个核心构建块

*   **`entrypoint` (入口点)**:
    *   封装整个工作流逻辑的函数。
    *   负责管理执行流、处理中断和长运行任务。
    *   **必须**接受一个位置参数作为输入。
    *   通常配置 `checkpointer` 以启用持久化。

*   **`task` (任务)**:
    *   代表一个离散的工作单元（如 API 调用、数据处理）。
    *   **关键特性**: 任务的执行结果会被保存到 Checkpoint 中。当工作流暂停并恢复时，已完成的任务**不会重新执行**，而是直接读取缓存的结果。
    *   必须在 `entrypoint` 或其他 `task` 内部调用。

### 2. 代码示例

以下示例展示了一个写文章并请求人工审核的工作流：

```typescript
import { MemorySaver, entrypoint, task, interrupt, Command } from "@langchain/langgraph";

// 1. 定义任务 (Task)
// 任务结果会被持久化，恢复时不会重新运行
const writeEssay = task("writeEssay", async (topic: string) => {
  await new Promise((resolve) => setTimeout(resolve, 1000)); // 模拟耗时操作
  return `An essay about topic: ${topic}`;
});

// 2. 定义入口点 (Entrypoint)
const workflow = entrypoint(
  { checkpointer: new MemorySaver(), name: "workflow" },
  async (topic: string) => {
    // 调用任务
    const essay = await writeEssay(topic);

    // 3. 人在回路 (Interrupt)
    // 暂停执行，等待外部输入
    const isApproved = interrupt({
      essay,
      action: "Please approve/reject the essay",
    });

    return {
      essay,
      isApproved,
    };
  }
);

// 4. 运行与恢复
async function run() {
  const config = { configurable: { thread_id: "thread-1" } };

  // 第一次运行：会执行 writeEssay，然后在 interrupt 处暂停
  const stream1 = await workflow.stream("cat", config);
  for await (const chunk of stream1) {
    console.log("Run 1:", chunk);
  }

  // 第二次运行：提供 resume 值来恢复
  // writeEssay 不会重新运行，直接使用上次结果
  const stream2 = await workflow.stream(new Command({ resume: true }), config);
  for await (const chunk of stream2) {
    console.log("Run 2:", chunk);
  }
}
```

### 3. 关键特性与机制

*   **持久化与恢复**:
    *   Functional API 不像 Graph API 那样在每个“超级步”后保存，而是将 `task` 的结果保存到与 `entrypoint` 关联的 Checkpoint 中。
    *   恢复时，LangGraph 会重放逻辑，遇到已执行过的 `task` 直接返回存储的结果，从而跳过重复计算。

*   **短期记忆 (`getPreviousState`)**:
    *   在同一个 `thread_id` 的连续调用之间，可以使用 `getPreviousState()` 获取上一次调用的返回值。
    *   使用 `entrypoint.final({ value, save })` 可以分离“返回给调用者的值”和“保存到 Checkpoint 的值”。

*   **确定性 (Determinism)**:
    *   **重要规则**: 任何包含随机性（如生成随机数、获取当前时间）或副作用（如写文件、API 调用）的操作，**必须**封装在 `task` 中。
    *   如果不封装在 `task` 中，恢复执行时这些代码会重新运行，可能导致逻辑分叉或副作用重复执行。

### 4. Functional API vs Graph API

| 特性 | Functional API | Graph API |
| :--- | :--- | :--- |
| **风格** | **命令式** (标准代码, if/else) | **声明式** (节点, 边, 状态图) |
| **状态管理** | 隐式 (函数作用域变量) | 显式 (State Schema, Reducers) |
| **可视化** | 不支持 (运行时动态生成) | 支持 (可生成图表) |
| **适用场景** | 快速原型、线性逻辑、现有代码改造 | 复杂多分支、多智能体协作、需要可视化 |

### 5. 常见陷阱

1.  **序列化**: `entrypoint` 的输入输出以及 `task` 的输出必须是 **JSON 可序列化**的，否则无法进行 Checkpoint。
2.  **副作用外泄**: 不要直接在 `entrypoint` 主体中写 `fs.writeFileSync` 或 `Math.random()`。如果工作流恢复，这些代码会再次执行。请务必将它们包裹在 `task` 中。

## 使用方法
这份文档指南介绍了如何使用 **LangGraph Functional API**。这是一种通过标准的编程方式（如函数、循环）来构建工作流的方法，同时赋予代码持久化、流式传输、记忆和“人在回路”（Human-in-the-loop）的能力。

以下是核心功能和使用模式的总结：

### 1. 核心概念

*   **`entrypoint` (入口点)**:
    *   封装工作流的主函数。
    *   **输入限制**: 函数只能接受**一个**位置参数作为输入（如果需要多个输入，请使用对象/字典）。
    *   **持久化**: 通过配置 `checkpointer`（如 `MemorySaver`）来启用状态保存。
*   **`task` (任务)**:
    *   被 `entrypoint` 调用的原子工作单元（如 LLM 调用、API 请求）。
    *   **特性**: 任务结果会被 Checkpoint 保存。当工作流暂停（如等待人工输入）并恢复时，已完成的任务**不会重新执行**，而是直接读取缓存结果。

### 2. 基础工作流示例

```typescript
import { entrypoint, task, MemorySaver } from "@langchain/langgraph";

// 1. 定义任务
const double = task("double", async (x: number) => {
  return x * 2;
});

// 2. 定义入口点
const checkpointer = new MemorySaver();

const workflow = entrypoint(
  { checkpointer, name: "workflow" },
  async (input: { value: number }) => {
    // 调用任务
    const result = await double(input.value);
    return { result };
  }
);

// 3. 运行
const config = { configurable: { thread_id: "1" } };
await workflow.invoke({ value: 5 }, config); // Output: { result: 10 }
```

### 3. 并行执行 (Parallel Execution)

利用 `Promise.all` 可以并发执行多个任务，特别适合 I/O 密集型操作（如并发调用 LLM）。

```typescript
const addOne = task("addOne", async (number: number) => {
  return number + 1;
});

const graph = entrypoint(
  { checkpointer, name: "graph" },
  async (numbers: number[]) => {
    // 并行执行所有任务
    return await Promise.all(numbers.map(addOne));
  }
);
```

### 4. 人在回路 (Human-in-the-loop)

使用 `interrupt` 暂停工作流等待人工反馈，使用 `Command` 恢复执行。

```typescript
import { interrupt, Command } from "@langchain/langgraph";

const workflow = entrypoint(
  { checkpointer, name: "hitl_workflow" },
  async (input: string) => {
    // 1. 执行一些逻辑...
    
    // 2. 暂停并等待输入
    const feedback = interrupt(`Please review: ${input}`);
    
    // 3. 恢复后继续执行
    return `Processed with feedback: ${feedback}`;
  }
);

// 运行流程：
// 第一次运行 (暂停)
await workflow.stream("test", config); 

// 第二次运行 (恢复)
// 使用 Command 提供 resume 值
await workflow.stream(new Command({ resume: "Approved" }), config);
```

### 5. 短期记忆与状态管理

Functional API 允许在同一 `thread_id` 的不同调用之间保持状态。

*   **读取上一轮状态**: `entrypoint` 的第二个参数 `previous` 包含上一轮的返回值。
*   **分离返回值与保存值**: 使用 `entrypoint.final({ value, save })`。`value` 返回给调用者，`save` 保存到 Checkpoint 供下一轮使用。

**累加器示例：**
```typescript
const accumulate = entrypoint(
  { checkpointer, name: "accumulate" },
  async (n: number, previous?: number) => {
    const prev = previous || 0;
    const total = prev + n;
    // 返回 prev 给调用者，但保存 total 到状态中
    return entrypoint.final({ value: prev, save: total });
  }
);

// invoke(1) -> 返回 0, 保存 1
// invoke(2) -> 返回 1, 保存 3
```

### 6. 可靠性配置

可以在定义 `task` 时配置重试策略和缓存策略。

*   **重试 (Retry)**: 处理网络抖动或临时错误。
*   **缓存 (Cache)**: 避免重复执行昂贵的操作（设置 TTL）。

```typescript
import { RetryPolicy } from "@langchain/langgraph";

const retryPolicy: RetryPolicy = { retryOn: (error) => error instanceof Error };

const unstableTask = task(
  {
    name: "unstable",
    retry: retryPolicy, // 配置重试
    cache: { ttl: 120 } // 配置缓存 (2分钟)
  },
  async () => { /* ... */ }
);
```

### 7. 互操作性
Functional API 和 Graph API 共享相同的运行时，因此：
*   可以在 `entrypoint` 中调用编译好的 Graph (`await graph.invoke(...)`)。
*   可以在 `entrypoint` 中调用其他的 `entrypoint`。