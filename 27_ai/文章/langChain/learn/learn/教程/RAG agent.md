这份文档介绍了如何使用 LangChain 构建 **RAG (检索增强生成)** 应用。RAG 结合了检索（从外部数据源查找信息）和生成（LLM 回答问题）的能力。

文档主要展示了两种构建 RAG 的方式：
1.  **RAG Agent**: 使用工具（Tool）进行检索，LLM 自主决定何时搜索以及搜索什么。
2.  **RAG Chain**: 固定的两步流程（检索 -> 生成），速度更快，适合简单问答。

以下是基于文档的核心步骤和代码实现总结：

### 1. 环境准备与数据索引 (Indexing)

在构建 RAG 之前，必须先将数据加载并存入向量数据库。这里以抓取一篇博客文章为例。

**安装依赖**:
```bash
npm install langchain @langchain/community @langchain/textsplitters @langchain/openai cheerio zod
```

**索引代码 (通用基础)**:
```typescript
import "cheerio";
import { CheerioWebBaseLoader } from "@langchain/community/document_loaders/web/cheerio";
import { RecursiveCharacterTextSplitter } from "@langchain/textsplitters";
import { OpenAIEmbeddings } from "@langchain/openai";
import { MemoryVectorStore } from "langchain/vectorstores/memory";

// 1. 加载数据
const loader = new CheerioWebBaseLoader(
  "https://lilianweng.github.io/posts/2023-06-23-agent/",
  { selector: "p" }
);
const docs = await loader.load();

// 2. 切分文本
const splitter = new RecursiveCharacterTextSplitter({
  chunkSize: 1000,
  chunkOverlap: 200,
});
const allSplits = await splitter.splitDocuments(docs);

// 3. 存入向量数据库
const embeddings = new OpenAIEmbeddings();
const vectorStore = new MemoryVectorStore(embeddings);
await vectorStore.addDocuments(allSplits);

console.log("Indexing completed.");
```

---

### 2. 方法一：RAG Agent (推荐用于复杂任务)

这种方式将检索功能封装为一个 **Tool**。Agent 可以进行多轮对话，甚至在回答一个问题时进行多次不同的搜索。

```typescript
import { createAgent, tool } from "langchain";
import { ChatOpenAI } from "@langchain/openai";
import { SystemMessage } from "@langchain/core/messages";
import * as z from "zod";

// 定义检索工具
const retrieveTool = tool(
  async ({ query }) => {
    const retrievedDocs = await vectorStore.similaritySearch(query, 2);
    const serialized = retrievedDocs
      .map((doc) => `Source: ${doc.metadata.source}\nContent: ${doc.pageContent}`)
      .join("\n");
    return serialized;
  },
  {
    name: "retrieve",
    description: "Retrieve information related to a query.",
    schema: z.object({ query: z.string() }),
  }
);

// 创建 Agent
const model = new ChatOpenAI({ model: "gpt-4" });
const agent = createAgent({
  model,
  tools: [retrieveTool],
  systemPrompt: new SystemMessage(
    "You have access to a tool that retrieves context from a blog post. Use the tool to help answer user queries."
  ),
});

// 运行 Agent
const result = await agent.invoke({
  messages: [{ role: "user", content: "What is Task Decomposition?" }],
});

console.log(result.messages[result.messages.length - 1].content);
```

---

### 3. 方法二：RAG Chain (推荐用于简单问答)

这种方式通过 **Middleware (中间件)** 在调用 LLM 之前拦截消息，执行搜索，并将结果注入到 System Prompt 中。这只需要一次 LLM 调用，速度更快。

```typescript
import { createAgent, dynamicSystemPromptMiddleware } from "langchain";
import { SystemMessage } from "@langchain/core/messages";

const ragChainAgent = createAgent({
  model,
  tools: [], // 不需要工具，因为检索是自动发生的
  middleware: [
    dynamicSystemPromptMiddleware(async (state) => {
      // 1. 获取用户最后一条消息
      const lastQuery = state.messages[state.messages.length - 1].content;

      // 2. 执行搜索
      const retrievedDocs = await vectorStore.similaritySearch(lastQuery as string, 2);
      const docsContent = retrievedDocs.map((doc) => doc.pageContent).join("\n\n");

      // 3. 动态构建 System Message 并注入上下文
      const systemMessage = new SystemMessage(
        `You are a helpful assistant. Use the following context in your response:\n\n${docsContent}`
      );

      // 4. 返回新的消息列表 (System Msg + 原有消息)
      return [systemMessage, ...state.messages];
    }),
  ],
});

// 运行 Chain
const result = await ragChainAgent.invoke({
  messages: [{ role: "user", content: "What is Task Decomposition?" }],
});

console.log(result.messages[result.messages.length - 1].content);
```

### 总结：如何选择？

| 特性 | RAG Agent (工具模式) | RAG Chain (链式模式) |
| :--- | :--- | :--- |
| **灵活性** | 高。LLM 决定是否搜索、搜索几次。 | 低。每次必搜，且只搜一次。 |
| **成本/速度** | 较慢。可能需要多次 LLM 调用（思考->搜索->回答）。 | 快。仅需一次 LLM 调用。 |
| **适用场景** | 复杂推理、多步问题、闲聊与搜索混合。 | 明确的知识库问答、对延迟敏感的场景。 |