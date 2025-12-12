## 选择 API
这份文档详细对比了 LangGraph 提供的两种构建工作流的 API：**Graph API (图 API)** 和 **Functional API (函数式 API)**。两者共享相同的底层运行时，但设计理念和适用场景不同。

以下是选择指南和核心差异的总结：

### 1. 快速决策指南

| 特性 | **Graph API (图 API)** | **Functional API (函数式 API)** |
| :--- | :--- | :--- |
| **核心范式** | **声明式** (定义节点、边、状态) | **命令式** (标准代码控制流) |
| **最佳场景** | 复杂的多分支、并行处理、团队协作 | 快速原型、线性逻辑、现有代码迁移 |
| **状态管理** | **显式共享状态** (Schema 定义) | **函数作用域** (局部变量传递) |
| **可视化** | 极佳 (清晰展示决策树和路径) | 较弱 (逻辑隐藏在代码中) |
| **控制流** | 通过条件边 (Conditional Edges) 实现 | 通过 `if/else`、循环实现 |

### 2. 详细对比

#### ✅ 何时使用 Graph API
当你的应用具有以下特征时，首选 Graph API：
*   **复杂的分支逻辑**: 需要清晰地可视化多个决策点和路径。
*   **并行执行**: 需要同时运行多个操作（如并发获取新闻、天气、股票），然后同步合并结果。
*   **全局状态共享**: 多个节点需要访问和修改同一个共享状态对象。
*   **团队协作**: 不同的团队成员负责不同的节点逻辑，通过图结构进行组装。

**代码风格 (Python 示例):**
```python
# 显式定义节点和边
workflow = StateGraph(AgentState)
workflow.add_node("call_llm", call_llm_node)
workflow.add_conditional_edges("call_llm", should_continue)
```

#### ✅ 何时使用 Functional API
当你的需求符合以下情况时，首选 Functional API：
*   **现有代码集成**: 希望在现有的过程式代码中添加 LangGraph 功能（如持久化、流式输出），且改动最小。
*   **标准控制流**: 逻辑主要是线性的，或者习惯使用标准的 `if/else` 和 `for` 循环。
*   **快速原型**: 不想编写繁琐的状态 Schema 定义，只想快速验证想法。
*   **局部状态**: 状态只需要在函数之间传递，不需要全局共享。

**代码风格 (Python 示例):**
```python
# 使用装饰器和标准 Python 逻辑
@entrypoint(checkpointer=checkpointer)
def workflow(user_input: str) -> str:
    processed = process_user_input(user_input).result()
    if "urgent" in processed:
        return handle_urgent(processed).result()
    return handle_normal(processed).result()
```

### 3. 混合使用与迁移
*   **混合使用**: 你可以在同一个应用中结合使用两者。例如，使用 Graph API 处理复杂的多智能体编排，而在某个节点内部使用 Functional API 处理简单的数据转换。
*   **迁移**:
    *   **Functional -> Graph**: 当工作流变得过于复杂，条件分支难以维护时，迁移到 Graph API。
    *   **Graph -> Functional**: 当发现图结构过于“过度设计”，逻辑其实只是简单的线性调用时，简化为 Functional API。

### 总结
*   选择 **Graph API** 以获得结构化、可视化和精细的控制。
*   选择 **Functional API** 以获得开发速度、简洁性和对现有代码的兼容性。

## 图API

这份文档详细介绍了 **LangGraph** 的核心构建模块——**Graph API**。它采用图论的方法来建模 Agent 工作流，通过节点（Nodes）和边（Edges）来定义逻辑和流转。

以下是核心概念和高级特性的总结：

### 1. 核心组件 (The Trinity)

*   **State (状态)**:
    *   图的共享数据结构。
    *   **Schema**: 通常使用 Zod 定义。支持定义 `Input`（输入）、`Output`（输出）和 `Private`（内部）Schema，以控制数据的可见性。
    *   **Reducers (归约器)**: 定义状态如何更新的关键。
        *   *默认行为*: 覆盖（Overwrite）。
        *   *自定义行为*: 例如 `messages` 字段通常需要追加（Append）而不是覆盖。LangGraph 提供了 `MessagesZodMeta` 来自动处理消息的序列化和 ID 追踪。
*   **Nodes (节点)**:
    *   执行具体工作的函数（可以是 LLM 调用，也可以是普通代码）。
    *   接收 `state` 和 `config` 作为参数。
    *   返回状态的**更新量**（Partial update），而不是整个状态。
    *   **Caching**: 支持通过 `cachePolicy` 对节点结果进行缓存（基于输入 key 和 TTL）。
*   **Edges (边)**:
    *   控制流转逻辑。
    *   **Normal Edges**: 固定流转 (A -> B)。
    *   **Conditional Edges**: 动态流转，根据函数返回的字符串决定下一个节点。

### 2. 高级控制流

*   **Command**:
    *   一种在节点内部**同时**进行“状态更新”和“路由跳转”的机制。
    *   替代条件边的更紧凑写法。
    *   **跨子图跳转**: 支持 `graph: Command.PARENT`，用于在多 Agent 系统中跳出子图，将控制权交还给父图（Handoffs）。
    *   **HITL**: 用于在中断后恢复执行 (`resume`)。
*   **Send**:
    *   用于 **Map-Reduce** 模式。
    *   允许从一个节点动态分发多个任务到同一个下游节点（但输入状态不同），实现并行处理。

### 3. 运行时与配置

*   **Recursion Limit (递归限制)**:
    *   防止无限循环，默认限制为 25 个 Super-steps。
    *   可以通过 `config.metadata.langgraph_step` 获取当前步数，从而在达到限制前**主动**进行降级处理（Proactive handling）。
*   **Runtime Context**:
    *   通过 `ContextSchema` 定义运行时配置（如模型名称、API Key 等），这些配置不属于图的状态，但在执行时通过 `config` 传入。

### 代码示例

以下示例展示了如何定义状态、使用 Reducer、以及使用 `Command` 进行路由：

```typescript
import { StateGraph, START, END, Command, MessagesZodMeta } from "@langchain/langgraph";
import { registry } from "@langchain/langgraph/zod";
import { z } from "zod";
import { AIMessage, BaseMessage, HumanMessage } from "@langchain/core/messages";

// 1. 定义状态 (使用 Reducer 处理消息追加)
const State = z.object({
  messages: z.array(z.custom<BaseMessage>())
    .register(registry, MessagesZodMeta), // 自动处理追加和序列化
  count: z.number().default(0),
});

// 2. 定义节点
const agentNode = (state: z.infer<typeof State>) => {
  const currentCount = state.count;
  
  // 使用 Command 同时更新状态和决定路由
  if (currentCount < 3) {
    return new Command({
      update: { 
        messages: [new AIMessage(`Count is ${currentCount}`)],
        count: currentCount + 1 
      },
      goto: "agent", // 循环调用自己
    });
  }
  
  return new Command({
    update: { messages: [new AIMessage("Done!")] },
    goto: END,
  });
};

// 3. 构建图
const workflow = new StateGraph(State)
  .addNode("agent", agentNode, {
    ends: ["agent", END] // 显式声明可能的跳转目标
  })
  .addEdge(START, "agent");

// 4. 编译
const app = workflow.compile();

// 5. 运行
await app.invoke({ messages: [new HumanMessage("Start")] });
```

### 总结
LangGraph 的 Graph API 提供了比简单的链式调用更强大的能力，特别是通过 **Reducers** 管理复杂状态，以及通过 **Command** 和 **Send** 处理复杂的动态流转和并行任务。

## 使用方法

这份文档详细指南介绍了如何使用 **LangGraph** 的 **Graph API** 来构建、控制和可视化 Agent 工作流。

以下是核心功能的总结：

### 1. 状态管理 (State Management)
*   **定义**: 使用 **Zod** Schema 定义图的状态。
*   **Reducers**: 默认情况下状态更新会覆盖旧值。使用 `reducer`（如数组拼接）可以实现增量更新（例如追加消息）。
*   **MessagesState**: 内置 `MessagesZodMeta` 专门用于处理聊天记录，支持自动追加消息和序列化。
*   **私有状态**: 支持定义仅在特定节点间传递的私有状态，或分离输入/输出 Schema。

### 2. 构建图结构
*   **基本组件**: 使用 `addNode` 添加处理逻辑（函数），使用 `addEdge` 连接节点（如 `START` -> `Node A` -> `END`）。
*   **顺序执行**: 简单的线性步骤。
*   **并行执行 (Branching)**: 一个节点连接多个下游节点时，它们会并行运行（Fan-out）。
*   **条件分支**: 使用 `addConditionalEdges` 根据状态动态决定下一个节点。
*   **循环**: 通过条件边路由回之前的节点，需配合 **Recursion Limit**（递归限制）防止死循环。

### 3. 高级控制流
*   **Map-Reduce (`Send` API)**: 允许从一个节点动态分发多个任务（不同的状态输入）到同一个下游节点，用于并行处理列表数据。
*   **Command API**: 允许在节点函数内部**同时**返回状态更新和路由指令（`goto`）。
    *   支持跳转到父图节点 (`graph: Command.PARENT`)。
    *   支持在工具（Tools）内部直接更新图状态。

### 4. 运行时配置与可靠性
*   **运行时配置**: 通过 `config` 参数在调用时传入动态参数（如模型选择、系统提示词），而不污染图状态。
*   **重试策略**: 为节点配置 `retryPolicy`，自动处理特定类型的错误（如网络波动）。

### 5. 可视化
*   支持将编译后的图导出为 **Mermaid** 语法或 **PNG** 图片，便于调试和文档化。

### 代码示例
以下是一个包含状态定义、节点逻辑和图构建的基础示例：

```typescript
import { StateGraph, START, END } from "@langchain/langgraph";
import { AIMessage, BaseMessage, HumanMessage } from "@langchain/core/messages";
import { MessagesZodMeta } from "@langchain/langgraph";
import { registry } from "@langchain/langgraph/zod";
import * as z from "zod";

// 1. 定义状态 (使用 MessagesZodMeta 处理消息追加)
const State = z.object({
  messages: z
    .array(z.custom<BaseMessage>())
    .register(registry, MessagesZodMeta),
  count: z.number().default(0),
});

// 2. 定义节点
const node = (state: z.infer<typeof State>) => {
  const newMessage = new AIMessage("Hello!");
  // 返回部分状态更新
  return { 
    messages: [newMessage], 
    count: state.count + 1 
  };
};

// 3. 构建图
const graph = new StateGraph(State)
  .addNode("bot", node)
  .addEdge(START, "bot")
  .addEdge("bot", END)
  .compile();

// 4. 运行
const result = await graph.invoke({ 
  messages: [new HumanMessage("Hi")] 
});

console.log(result);
```