这份文档详细介绍了 **LangChain.js** 中的 **Retrieval (检索)** 机制及其在 **RAG (检索增强生成)** 中的应用。

检索旨在解决 LLM 的两大核心限制：**有限的上下文窗口**和**静态的训练数据**。

以下是文档的核心内容总结：

### 1. 检索管道 (Retrieval Pipeline)
构建 RAG 应用通常涉及以下模块化组件：
*   **Document Loaders**: 从外部源（如 PDF, Notion, Slack）加载数据。
*   **Text Splitters**: 将文档分割成适合处理的块。
*   **Embedding Models**: 将文本转换为向量。
*   **Vector Stores**: 存储和搜索向量数据的数据库。
*   **Retrievers**: 根据查询返回相关文档的接口。

### 2. RAG 架构模式
文档对比了三种主要的 RAG 实现架构：

| 架构 | 描述 | 特点 | 适用场景 |
| :--- | :--- | :--- | :--- |
| **2-Step RAG** | 线性流程：先检索，后生成。 | ✅ 控制力高<br>⚡️ 延迟低且可预测 | FAQ 机器人、文档问答 |
| **Agentic RAG** | 由 Agent 决定**何时**以及**如何**检索。检索被封装为工具 (Tool)。 | ✅ 灵活性高<br>⏳ 延迟可变 | 研究助手、复杂推理任务 |
| **Hybrid RAG** | 结合两者，增加了查询优化、结果验证和自我修正循环。 | ⚖️ 平衡控制与灵活<br>🔄 支持迭代优化 | 需要高质量验证的领域问答 |

### 3. 代码示例：Agentic RAG
在 Agentic RAG 中，检索能力被封装为一个工具，Agent 根据用户的问题自主决定是否查阅文档。

```typescript
import { tool, createAgent } from "langchain";
import * as z from "zod";

// 1. 定义检索工具
const fetchDocumentation = tool(
  async (input) => {
    // 模拟获取文档内容的逻辑
    // 在实际应用中，这里通常是查询向量数据库或调用外部 API
    return `Fetched content for url: ${input.url}`;
  },
  {
    name: "fetch_documentation",
    description: "Fetch and convert documentation from a URL",
    schema: z.object({
      url: z.string().describe("The URL of the documentation to fetch"),
    }),
  }
);

// 2. 创建 Agent 并赋予工具
const agent = createAgent({
  model: "claude-sonnet-4-0",
  tools: [fetchDocumentation], // Agent 可以自主决定是否使用此工具
  systemPrompt: "You are a helpful assistant. Use the fetch_documentation tool if you need external info.",
});

// 3. 调用 Agent
const response = await agent.invoke({
  messages: [
    { role: "user", content: "Check the docs at https://example.com/docs and summarize them." }
  ],
});
```

### 总结
*   如果你的任务流程固定且对延迟敏感，选择 **2-Step RAG**。
*   如果任务需要多步推理或不确定是否需要检索，选择 **Agentic RAG**。
*   如果需要极高的准确性和自我纠错能力，选择 **Hybrid RAG**。