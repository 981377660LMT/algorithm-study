这份文档介绍了如何在 **LangChain.js** 中利用 **LangGraph** 的持久化层（Persistence Layer）来实现 Agent 的 **长期记忆 (Long-term Memory)**。

与短期记忆（State，通常随对话结束而消失）不同，长期记忆允许 Agent 跨对话存储和检索用户偏好、历史记录或特定知识。

以下是核心概念和使用指南：

### 1. 存储机制 (Memory Store)
LangGraph 使用 **Store** 来保存数据。数据以 JSON 文档的形式存储，结构类似于文件系统：

*   **Namespace (命名空间)**：类似于文件夹路径，是一个字符串数组（例如 `["users", "user_123"]`）。通常用于按用户、组织或应用上下文隔离数据。
*   **Key (键)**：类似于文件名，是一个字符串（例如 `"profile"`）。
*   **Value (值)**：实际存储的 JSON 对象。

**Store 的主要操作：**
*   `put(namespace, key, value)`: 写入或更新数据。
*   `get(namespace, key)`: 读取数据。
*   `search(namespace, options)`: 搜索数据（支持基于内容的过滤，如果配置了 Embedding 模型，还支持语义搜索）。

### 2. 集成步骤
要在 Agent 中使用长期记忆，通常需要三个步骤：
1.  初始化一个 `Store` 实例（如 `InMemoryStore` 用于测试，生产环境应使用数据库支持的 Store）。
2.  在 `createAgent` 时将 `store` 传入。
3.  在 **Tools (工具)** 中通过 `runtime.store` 进行读写。

### 3. 代码示例

#### A. 初始化与配置
```typescript
import { InMemoryStore } from "@langchain/langgraph";
import { createAgent } from "langchain";

// 初始化存储 (生产环境请使用持久化数据库)
const store = new InMemoryStore();

const agent = createAgent({
  model: "gpt-4o",
  tools: [/* ... */],
  // 关键：将 store 传递给 agent
  store: store, 
  // 定义上下文 Schema 以便传递 userId
  contextSchema: z.object({ userId: z.string() }),
});
```

#### B. 在工具中读取记忆 (Read)
工具可以通过 `runtime.store.get` 获取存储的数据。

```typescript
const getUserInfo = tool(
  async (_, runtime: ToolRuntime) => {
    // 1. 从上下文中获取用户 ID
    const userId = runtime.context.userId;
    
    // 2. 从 Store 中读取数据
    // namespace 为 ["users"]，key 为 userId
    const memory = await runtime.store.get(["users"], userId);
    
    return memory ? JSON.stringify(memory.value) : "User not found";
  },
  {
    name: "get_user_info",
    description: "Get user profile from long-term memory",
    schema: z.object({}),
  }
);
```

#### C. 在工具中写入记忆 (Write)
工具可以通过 `runtime.store.put` 保存或更新数据。

```typescript
const saveUserInfo = tool(
  async ({ name, language }, runtime: ToolRuntime) => {
    const userId = runtime.context.userId;
    
    // 3. 将数据写入 Store
    await runtime.store.put(
      ["users"], // Namespace
      userId,    // Key
      { name, language } // Value (JSON)
    );
    
    return "User info saved.";
  },
  {
    name: "save_user_info",
    description: "Save user profile to long-term memory",
    schema: z.object({
      name: z.string(),
      language: z.string(),
    }),
  }
);
```

### 总结
*   **长期记忆** 依赖于 LangGraph 的 `Store` 接口。
*   数据通过 **Namespace** 和 **Key** 进行组织。
*   Agent 的 **Tools** 是读写长期记忆的主要场所，通过 `runtime.store` 访问。
*   在调用 `agent.invoke` 时，务必通过 `context` 传入必要的标识符（如 `userId`），以便工具知道去哪里查找数据。