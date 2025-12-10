这份文档详细介绍了 **LangChain.js 中的 Short-term Memory (短期记忆)** 机制。短期记忆是 Agent 在单次会话（Thread）中记住上下文、用户偏好和交互历史的关键能力。

以下是对文档核心内容的详细解读：

### 1. 核心概念：为什么需要短期记忆？
*   **定义**：短期记忆允许应用程序在单个线程或对话中记住之前的交互。
*   **痛点**：
    *   LLM 的上下文窗口（Context Window）是有限的。
    *   即使窗口够大，过长的上下文会导致模型“分心”（注意力分散）、响应变慢且成本增加。
*   **解决方案**：LangChain 通过 **Checkpointer** 将状态持久化到数据库中，并在每次交互时加载相关上下文。

### 2. 基础用法 (Basic Usage)
要启用短期记忆，只需在 `createAgent` 时传入一个 `checkpointer`。

*   **开发环境**：使用 `MemorySaver`（内存存储，重启即失）。
*   **生产环境**：使用 `PostgresSaver` 等数据库适配器（持久化存储）。
*   **调用方式**：在 `invoke` 时传入 `configurable: { thread_id: "..." }` 来指定当前对话的 ID。

```typescript
import { createAgent } from "langchain";
import { MemorySaver } from "@langchain/langgraph";

const agent = createAgent({
    model: "gpt-4o",
    tools: [],
    checkpointer: new MemorySaver(), // 启用记忆
});

// 传入 thread_id 以恢复之前的对话状态
await agent.invoke(
    { messages: [{ role: "user", content: "Hi!" }] },
    { configurable: { thread_id: "1" } }
);
```

### 3. 自定义状态 (Customizing State)
默认情况下，Agent 只记住 `messages`。你可以通过 `middleware` 和 `stateSchema` 扩展状态，让 Agent 记住更多业务数据（如 `userId`, `preferences`）。

```typescript
const customStateSchema = z.object({
    userId: z.string(),
    preferences: z.record(z.string(), z.any()),
});
// ...在 createAgent 中配置 middleware
```

### 4. 管理上下文窗口 (Common Patterns)
为了防止对话历史撑爆 LLM 的 Token 限制，文档提供了三种策略：

1.  **Trim Messages (裁剪消息)**：
    *   使用 `trimMessages` 工具函数。
    *   策略：保留最后 N 个 Token，或最后 N 条消息。
    *   实现：通过 `beforeModel` 中间件在调用模型前执行裁剪。

2.  **Delete Messages (删除消息)**：
    *   使用 `RemoveMessage` 类。
    *   策略：从 Graph 状态中永久删除特定消息。
    *   实现：在 `postModelHook` 或中间件中返回 `new RemoveMessage({ id: ... })`。

3.  **Summarize Messages (总结消息)**：
    *   使用 `summarizationMiddleware`。
    *   策略：当消息达到一定数量时，调用 LLM 将旧的历史记录压缩成一段摘要（Summary），替换掉原始消息。

### 5. 访问与修改记忆 (Access Memory)
文档展示了如何在不同组件中读写记忆：

*   **在 Tools 中**：
    *   **读**：通过 `config.context` 读取（如获取 `userId`）。
    *   **写**：返回 `Command` 对象来更新状态（如工具执行后更新用户的 `userName`）。
*   **在 Prompt 中**：
    *   动态生成系统提示词，例如 `Address the user as ${config.context?.userName}`。
*   **在 Middleware 中**：
    *   **Before Model**：在模型调用前修改状态（如裁剪历史）。
    *   **After Model**：在模型返回后验证或修改结果（如敏感词过滤）。

### 总结
LangChain 的短期记忆机制不仅仅是简单的“保存聊天记录”。它是一个完整的**状态管理系统**，结合 **LangGraph** 的能力，允许开发者精细控制：
1.  **存什么**（自定义 Schema）。
2.  **存哪里**（Checkpointer）。
3.  **怎么存**（裁剪、总结、删除）。
4.  **怎么用**（在 Tool 和 Prompt 中动态读取）。