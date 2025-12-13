这份文档介绍了 **Deep Agents** (`deepagents`)，这是一个基于 LangGraph 构建的独立库，旨在构建能够处理复杂、多步骤任务的智能体（Agent）。

以下是该库的核心概览：

### 1. 核心定位
*   **目标**: 构建具备规划能力、能使用子 Agent、并利用文件系统处理复杂任务的 Agent。
*   **灵感来源**: Claude Code, Deep Research, Manus。
*   **适用场景**: 当你需要处理复杂的长流程任务、管理大量上下文、需要任务委派或持久化记忆时使用。对于简单任务，建议仍使用标准的 LangChain Agent。

### 2. 核心能力
*   **规划与任务拆解**:
    *   内置 `write_todos` 工具。
    *   允许 Agent 将大任务分解为离散步骤，并根据新信息动态调整计划。
*   **上下文管理 (文件系统)**:
    *   提供 `ls`, `read_file`, `write_file`, `edit_file` 等工具。
    *   允许 Agent 将大量上下文“卸载”到文件系统中，防止上下文窗口溢出，有效处理变长的工具结果。
*   **子 Agent 生成 (Subagents)**:
    *   内置 `task` 工具。
    *   允许生成专门的子 Agent 来处理特定子任务。这有助于**隔离上下文**，保持主 Agent 的上下文整洁。
*   **长期记忆**:
    *   结合 LangGraph Store 实现跨对话和线程的信息保存与检索。

### 3. 生态系统关系
*   **LangGraph**: 提供底层的图执行和状态管理机制。
*   **LangChain**: 所有的工具和模型集成都可以无缝配合 Deep Agents 使用。
*   **LangSmith**: 提供可观测性、评估和部署支持。

---

这份文档介绍了如何使用 `deepagents` 库快速构建一个具备**规划能力**、**文件系统工具**和**子 Agent 能力**的深度智能体（Deep Agent）。

以下是整合了文档步骤的完整代码示例。该示例创建了一个研究 Agent，它可以使用 Tavily 搜索引擎进行资料搜集并撰写报告。

### 前置要求
请确保已安装依赖并设置环境变量：
```bash
npm install deepagents @langchain/tavily zod
export ANTHROPIC_API_KEY="your-api-key"
export TAVILY_API_KEY="your-tavily-api-key"
```

### 完整代码示例

```typescript
import { tool } from "@langchain/core/tools";
import { TavilySearch } from "@langchain/tavily";
import { z } from "zod";
import { createDeepAgent } from "deepagents";

// 1. 定义搜索工具
const internetSearch = tool(
  async ({
    query,
    maxResults = 5,
    topic = "general",
    includeRawContent = false,
  }: {
    query: string;
    maxResults?: number;
    topic?: "general" | "news" | "finance";
    includeRawContent?: boolean;
  }) => {
    const tavilySearch = new TavilySearch({
      maxResults,
      tavilyApiKey: process.env.TAVILY_API_KEY,
      includeRawContent,
      topic,
    });
    return await tavilySearch.invoke({ query });
  },
  {
    name: "internet_search",
    description: "Run a web search",
    schema: z.object({
      query: z.string().describe("The search query"),
      maxResults: z
        .number()
        .optional()
        .default(5)
        .describe("Maximum number of results to return"),
      topic: z
        .enum(["general", "news", "finance"])
        .optional()
        .default("general")
        .describe("Search topic category"),
      includeRawContent: z
        .boolean()
        .optional()
        .default(false)
        .describe("Whether to include raw content"),
    }),
  }
);

// 2. 定义系统提示词
const researchInstructions = `You are an expert researcher. Your job is to conduct thorough research and then write a polished report.

You have access to an internet search tool as your primary means of gathering information.

## \`internet_search\`

Use this to run an internet search for a given query. You can specify the max number of results to return, the topic, and whether raw content should be included.
`;

// 3. 创建 Deep Agent
const agent = createDeepAgent({
  tools: [internetSearch],
  systemPrompt: researchInstructions,
});

// 4. 运行 Agent
async function main() {
  console.log("Starting research...");
  const result = await agent.invoke({
    messages: [{ role: "user", content: "What is langgraph?" }],
  });

  // 打印最终回复
  const lastMessage = result.messages[result.messages.length - 1];
  console.log("\nFinal Report:\n");
  console.log(lastMessage.content);
}

main().catch(console.error);
```

### 核心特性说明

当你运行这段代码时，`deepagents` 会自动处理以下复杂逻辑：

1.  **自动规划**: 使用内置的 `write_todos` 工具将 "What is langgraph?" 分解为具体的搜索和总结步骤。
2.  **上下文管理**: 如果搜索结果过长，Agent 会自动使用内置的文件系统工具（`write_file`, `read_file`）将内容写入文件，避免超出上下文窗口。
3.  **子任务委派**: 如果任务足够复杂，它可能会生成子 Agent 来并行处理特定的搜索任务。
  
---

这份文档介绍了如何自定义 **Deep Agents** 的核心配置，包括模型、系统提示词和工具。

以下是自定义选项的总结：

### 1. 模型 (Model)
默认使用 `claude-sonnet-4-5-20250929`。你可以通过传递 LangChain 模型对象（如 `ChatOpenAI`、`ChatAnthropic`）或模型标识符字符串来更改底层模型。

````typescript
import { ChatOpenAI } from "@langchain/openai";
import { createDeepAgent } from "deepagents";

const agent = createDeepAgent({
  model: new ChatOpenAI({
    model: "gpt-4o", // 替换为 OpenAI 模型
    temperature: 0,
  }),
});
````

### 2. 系统提示词 (System Prompt)
Deep Agents 内置了受 Claude Code 启发的默认提示词，涵盖了规划、文件系统操作和子 Agent 的使用说明。你可以通过 `systemPrompt` 参数传入针对特定用例的指令。

### 3. 工具 (Tools)
你可以通过 `tools` 参数传入自定义工具（例如网络搜索工具）。

除了自定义工具外，Deep Agents 始终会自动包含以下**内置工具**：
*   **规划**: `write_todos` (更新待办事项列表)
*   **文件系统**: `ls`, `read_file`, `write_file`, `edit_file` (管理文件上下文)
*   **子智能体**: `task` (生成子 Agent 处理特定任务)

````typescript
import { createDeepAgent } from "deepagents";
// ... 假设 internetSearch 已定义 ...

const agent = createDeepAgent({
  tools: [internetSearch], // 添加自定义工具
  systemPrompt: "You are an expert researcher.", // 自定义提示词
});
````