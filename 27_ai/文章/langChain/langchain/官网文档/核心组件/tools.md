这份文档介绍了 **LangChain.js 中的 Tools (工具)** 模块。如果说模型是 Agent 的“大脑”，那么工具就是它的“手”，让 AI 能够与外部世界（API、数据库、文件系统）进行交互。

以下是对文档核心内容的详细解读：

### 1. 什么是工具 (Tools)？
*   **定义**：工具是一个封装了**可调用函数**和**输入 Schema** 的组件。
*   **作用**：它允许 LLM 生成符合特定格式的参数，从而调用外部系统。
*   **服务端工具**：文档特别提到，某些模型（如 OpenAI, Anthropic）内置了服务端执行的工具（如联网搜索、代码解释器），这些不需要你手动编写代码，只需配置即可。

### 2. 如何创建工具？
最简单的方法是使用 `tool` 函数搭配 `zod` 库。

*   **逻辑函数**：第一个参数是实际执行的函数。
*   **元数据**：第二个参数包含 `name`（名称）、`description`（给 AI 看的说明书）和 `schema`（参数验证规则）。

```typescript
import * as z from "zod";
import { tool } from "langchain";

const searchDatabase = tool(
  ({ query }) => `Searching for ${query}...`, // 实际逻辑
  {
    name: "search_database",
    description: "查询客户数据库",
    schema: z.object({
      query: z.string().describe("查询关键词"), // Zod 定义参数
    }),
  }
);
```

### 3. 核心高级特性：访问上下文 (Accessing Context)
这是文档中最具技术含量的部分。工具不仅仅是孤立的函数，它可以通过第二个参数 `config` 访问 Agent 的运行时状态。

#### A. 运行时上下文 (Runtime Context)
*   **场景**：工具需要知道当前是谁在调用（User ID），或者获取当前的配置信息，但你不想使用全局变量。
*   **方法**：通过 `config.context` 访问。
*   **配置**：在 `createAgent` 时需要定义 `contextSchema`。

#### B. 记忆存储 (Store)
*   **场景**：工具需要**跨会话**存取数据（例如：记住用户的名字、偏好设置）。
*   **方法**：通过 `config.store` 访问（通常配合 `InMemoryStore` 或数据库存储）。
*   **能力**：工具可以直接读写这个 Store，实现“长时记忆”。

#### C. 流式写入 (Stream Writer)
*   **场景**：工具执行时间较长，你想给用户实时反馈（例如：“正在连接数据库...”、“已获取数据，正在处理...”）。
*   **方法**：使用 `config.streamWriter`。
*   **效果**：这些中间状态信息会实时流式传输给前端，提升用户体验。

### 总结
这份文档告诉我们，LangChain 中的工具不仅仅是简单的函数映射。通过 **Context**、**Store** 和 **StreamWriter**，工具变成了**有状态的**、**可感知的**组件，这对于构建复杂的生产级 Agent 至关重要。