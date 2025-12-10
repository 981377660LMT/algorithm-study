这份文档详细介绍了 **LangChain.js 中的 Structured Output (结构化输出)** 机制。它允许 Agent 返回类型安全、格式确定的数据（如 JSON），而不是非结构化的自然语言文本。

以下是对文档核心内容的分析和讲解：

### 1. 核心概念
*   **目的**：将 LLM 的输出从“一段话”转变为“可编程的数据对象”。
*   **实现方式**：在 `createAgent` 中配置 `responseFormat` 参数。
*   **结果存储**：结构化数据会自动解析并存储在 Agent 状态的 `structuredResponse` 字段中。

### 2. 两种实现策略
LangChain 提供了两种策略来生成结构化数据，你可以手动指定，也可以让框架自动选择：

#### A. Provider Strategy (原生策略)
*   **原理**：利用模型提供商（如 OpenAI, Gemini）API 原生支持的 JSON Mode 或 Structured Output 功能。
*   **优点**：最可靠，严格遵循 Schema。
*   **用法**：使用 `providerStrategy(schema)`。
*   **适用**：当底层模型原生支持结构化输出时。

#### B. Tool Calling Strategy (工具调用策略)
*   **原理**：利用模型的 **Tool Calling** 能力。LangChain 会创建一个“虚拟工具”，其参数结构就是你想要的输出格式，强制模型调用这个工具来生成数据。
*   **优点**：通用性强，适用于所有支持工具调用的模型。
*   **用法**：使用 `toolStrategy(schema)`。
*   **特性**：支持自定义 `toolMessageContent`（在历史记录中显示的文本）和 `handleError`（错误处理逻辑）。

### 3. Schema 定义
支持两种定义方式：
1.  **Zod Schema** (推荐)：TypeScript 友好，支持类型推断和验证。
2.  **JSON Schema**：标准的 JSON 对象定义。

### 4. 自动错误修正 (Self-Correction)
这是 LangChain Agent 的一大亮点。如果模型生成的结构不符合 Schema（例如：评分范围要求 1-5，模型输出了 10），Agent 会自动执行以下流程：
1.  捕获验证错误。
2.  将错误信息（如 "Rating must be <= 5"）封装为 `ToolMessage` 发回给模型。
3.  **模型根据错误信息进行自我修正并重试**。
4.  直到生成符合要求的数据或达到重试上限。

### 代码示例
以下是一个使用 **Zod** 和 **Tool Strategy** 提取用户信息的完整示例：

```typescript
import * as z from "zod";
import { createAgent, toolStrategy } from "langchain";

// 1. 定义输出结构 (Zod Schema)
const UserProfile = z.object({
  name: z.string().describe("用户的姓名"),
  age: z.number().describe("用户的年龄"),
  interests: z.array(z.string()).describe("用户的兴趣爱好列表"),
});

// 2. 创建 Agent
const agent = createAgent({
  model: "gpt-4o", // 或其他支持工具调用的模型
  tools: [],
  // 强制要求输出符合 UserProfile 结构
  responseFormat: toolStrategy(UserProfile), 
});

// 3. 调用 Agent
const result = await agent.invoke({
  messages: [
    { role: "user", content: "我叫 Alice，今年 25 岁，喜欢编程和登山。" }
  ]
});

// 4. 获取类型安全的结构化结果
console.log(result.structuredResponse);
/* 输出:
{
  name: "Alice",
  age: 25,
  interests: ["编程", "登山"]
}
*/
```

### 总结
使用 `createAgent` 的 `responseFormat` 是获取结构化数据的最佳实践。它屏蔽了底层模型 API 的差异（原生 vs 工具调用），并内置了强大的**错误重试机制**，确保你拿到的数据一定是符合 Schema 定义的。