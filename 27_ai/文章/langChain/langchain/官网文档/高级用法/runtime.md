这份文档介绍了 **LangChain.js 的运行时 (Runtime)** 机制。它是 Agent 执行的基础环境（底层基于 LangGraph），主要用于在 Agent、工具和中间件之间传递数据和状态，避免使用全局变量。

以下是核心知识点总结：

### 1. Runtime 对象包含什么？
工具函数的第二个参数即为 `runtime` 对象
Runtime 对象主要包含三个部分：
*   **Context (上下文)**：静态信息，如用户 ID、数据库连接、会话配置等。这些数据在单次调用（Invocation）中是只读且共享的。
*   **Store (存储)**：用于读写**长期记忆 (Long-term Memory)** 的 `BaseStore` 实例。
*   **Stream Writer (流式写入器)**：用于发送自定义流式更新（例如工具内部的进度条）。

### 2. 如何配置 Context
你需要先定义结构，然后在调用时传入数据。

1.  **定义 Schema**：在 `createAgent` 中使用 Zod 定义 `contextSchema`。
2.  **传入数据**：在 `agent.invoke` 的第二个配置参数中传入具体的 `context` 对象。

```typescript
// 1. 定义 Schema
const contextSchema = z.object({
  userName: z.string(),
});

const agent = createAgent({
  // ...
  contextSchema, // 绑定 Schema
});

// 2. 调用时传入数据
await agent.invoke(
  { messages: [...] },
  { context: { userName: "John Smith" } } // 注入 Context
);
```

### 3. 如何访问 Runtime 信息

#### 在工具 (Tools) 中
工具函数的第二个参数即为 `runtime` 对象。你可以用它来获取用户信息或读取长期记忆。

```typescript
const myTool = tool(
  async (input, runtime) => { // runtime 包含 context, store, writer
    const user = runtime.context.userName;
    // ...
  },
  // ...
);
```

#### 在中间件 (Middleware) 中
中间件 Hooks（如 `beforeModel`, `afterModel`）的第二个参数也是 `runtime` 对象。你可以利用它来动态修改 Prompt 或记录特定用户的日志。

```typescript
const myMiddleware = createMiddleware({
  // ...
  beforeModel: (state, runtime) => {
    console.log(`User: ${runtime.context.userName}`);
    // ...
  },
});
```

### 总结
Runtime 机制实现了**依赖注入**模式，使得 Agent 保持无状态（Stateless）、易于测试且可复用。通过 Context，你可以安全地将业务数据（如 User ID）穿透传递到深层的工具逻辑中。