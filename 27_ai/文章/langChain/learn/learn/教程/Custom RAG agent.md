这份文档详细介绍了如何使用 **LangGraph** 构建一个自定义的 **RAG (检索增强生成) Agent**。与标准的 RAG 链不同，这个 Agent 具备决策能力：它能决定是否需要检索、评估检索到的文档是否相关，并在必要时重写查询。

以下是基于文档内容的完整 TypeScript 实现代码。

### 1. 安装依赖

```bash
npm install @langchain/langgraph @langchain/openai @langchain/community @langchain/textsplitters cheerio zod
```

### 2. 完整代码实现

这段代码将文档加载、索引、工具创建以及 LangGraph 的节点和边定义整合在一起。

```typescript
import { CheerioWebBaseLoader } from "@langchain/community/document_loaders/web/cheerio";
import { RecursiveCharacterTextSplitter } from "@langchain/textsplitters";
import { MemoryVectorStore } from "langchain/vectorstores/memory";
import { OpenAIEmbeddings, ChatOpenAI } from "@langchain/openai";
import { createRetrieverTool } from "langchain/tools/retriever";
import { HumanMessage, AIMessage, BaseMessage } from "@langchain/core/messages";
import { StateGraph, START, END, Annotation } from "@langchain/langgraph";
import { ToolNode } from "@langchain/langgraph/prebuilt";
import { ChatPromptTemplate } from "@langchain/core/prompts";
import * as z from "zod";

// --- 1. 预处理文档与创建索引 ---
async function setupRetriever() {
  const urls = [
    "https://lilianweng.github.io/posts/2023-06-23-agent/",
    "https://lilianweng.github.io/posts/2023-03-15-prompt-engineering/",
    "https://lilianweng.github.io/posts/2023-10-25-adv-attack-llm/",
  ];

  const docs = await Promise.all(
    urls.map((url) => new CheerioWebBaseLoader(url).load())
  );
  const docsList = docs.flat();

  const textSplitter = new RecursiveCharacterTextSplitter({
    chunkSize: 500,
    chunkOverlap: 50,
  });
  const docSplits = await textSplitter.splitDocuments(docsList);

  const vectorStore = await MemoryVectorStore.fromDocuments(
    docSplits,
    new OpenAIEmbeddings()
  );

  return vectorStore.asRetriever();
}

// --- 2. 定义 Graph State ---
// 定义图的状态结构，这里只需要存储消息列表
const GraphState = Annotation.Root({
  messages: Annotation<BaseMessage[]>({
    reducer: (x, y) => x.concat(y),
    default: () => [],
  }),
});

// --- 3. 定义节点逻辑 ---

async function runAgent() {
  const retriever = await setupRetriever();
  
  // 创建检索工具
  const tool = createRetrieverTool(retriever, {
    name: "retrieve_blog_posts",
    description:
      "Search and return information about Lilian Weng blog posts on LLM agents, prompt engineering, and adversarial attacks on LLMs.",
  });
  const tools = [tool];
  const toolNode = new ToolNode(tools);

  // 节点 1: 决定查询或直接回复
  async function generateQueryOrRespond(state: typeof GraphState.State) {
    const { messages } = state;
    const model = new ChatOpenAI({ model: "gpt-4o", temperature: 0 }).bindTools(tools);
    const response = await model.invoke(messages);
    return { messages: [response] };
  }

  // 节点 2: 评估文档相关性
  async function gradeDocuments(state: typeof GraphState.State) {
    const { messages } = state;
    // 获取检索到的内容（通常是 ToolMessage）
    const lastToolMessage = messages[messages.length - 1]; 
    const context = lastToolMessage.content as string;
    
    // 获取原始问题（假设是第一条消息，实际场景可能需要更复杂的逻辑查找最近的问题）
    const question = messages[0].content as string;

    const gradeDocumentsSchema = z.object({
      binaryScore: z.string().describe("Relevance score 'yes' or 'no'"),
    });

    const model = new ChatOpenAI({
      model: "gpt-4o",
      temperature: 0,
    }).withStructuredOutput(gradeDocumentsSchema);

    const prompt = ChatPromptTemplate.fromTemplate(
      `You are a grader assessing relevance of retrieved docs to a user question.
      Here are the retrieved docs:
      \n ------- \n
      {context}
      \n ------- \n
      Here is the user question: {question}
      If the content of the docs are relevant to the users question, score them as relevant.
      Give a binary score 'yes' or 'no' score to indicate whether the docs are relevant to the question.
      Yes: The docs are relevant to the question.
      No: The docs are not relevant to the question.`
    );

    const score = await prompt.pipe(model).invoke({ question, context });

    // 我们不返回消息，而是返回一个特殊的标记给条件边使用
    // 注意：在实际 LangGraph 中，通常通过更新 State 中的某个字段来传递决策，
    // 或者直接在 Conditional Edge 中执行此逻辑。
    // 为了匹配文档逻辑，我们将结果作为临时状态处理，或者直接返回给 Edge 判断。
    // 这里我们简化处理：将评分结果作为一条不可见的系统消息或直接由 Edge 函数调用。
    
    // *修正*: 文档中的 gradeDocuments 返回了字符串 "generate" 或 "rewrite"。
    // 但在 LangGraph 中，Node 必须返回 State 的更新。
    // 因此，我们将逻辑稍微调整：gradeDocuments 仅负责评分，
    // 真正的路由逻辑放在 Conditional Edge 中，或者我们在这里返回一个特殊的 AIMessage 供后续判断。
    
    return { messages: [new AIMessage({ content: score.binaryScore === "yes" ? "generate" : "rewrite" })] };
  }

  // 节点 3: 重写问题
  async function rewrite(state: typeof GraphState.State) {
    const { messages } = state;
    const question = messages[0].content;
    const model = new ChatOpenAI({ model: "gpt-4o", temperature: 0 });
    const rewritePrompt = ChatPromptTemplate.fromTemplate(
      `Look at the input and try to reason about the underlying semantic intent / meaning. \n
      Here is the initial question: {question} \n
      Formulate an improved question:`
    );
    const response = await rewritePrompt.pipe(model).invoke({ question });
    // 将重写后的问题作为新的 HumanMessage 添加，以便再次触发检索
    // 注意：实际生产中可能需要替换原问题或标记为新一轮
    return { messages: [response] }; 
  }

  // 节点 4: 生成最终答案
  async function generate(state: typeof GraphState.State) {
    const { messages } = state;
    const question = messages[0].content;
    // 找到最近的一条 ToolMessage 作为上下文
    const contextMsg = messages.slice().reverse().find(m => m._getType() === "tool");
    const context = contextMsg ? contextMsg.content : "";

    const prompt = ChatPromptTemplate.fromTemplate(
      `You are an assistant for question-answering tasks.
      Use the following pieces of retrieved context to answer the question.
      If you don't know the answer, just say that you don't know.
      Question: {question}
      Context: {context}`
    );
    const llm = new ChatOpenAI({ model: "gpt-4o", temperature: 0 });
    const response = await prompt.pipe(llm).invoke({ context, question });
    return { messages: [response] };
  }

  // --- 4. 构建图 ---

  // 条件边逻辑：是否检索
  function shouldRetrieve(state: typeof GraphState.State) {
    const { messages } = state;
    const lastMessage = messages[messages.length - 1];
    // @ts-ignore
    if (lastMessage.tool_calls?.length) {
      return "retrieve";
    }
    return END;
  }

  // 条件边逻辑：评分后的路由
  function checkRelevance(state: typeof GraphState.State) {
    const { messages } = state;
    const lastMessage = messages[messages.length - 1];
    // gradeDocuments 节点返回了 "generate" 或 "rewrite" 作为内容
    if (lastMessage.content === "generate") {
      return "generate";
    }
    return "rewrite";
  }

  const workflow = new StateGraph(GraphState)
    .addNode("generateQueryOrRespond", generateQueryOrRespond)
    .addNode("retrieve", toolNode)
    .addNode("gradeDocuments", gradeDocuments)
    .addNode("rewrite", rewrite)
    .addNode("generate", generate)
    
    .addEdge(START, "generateQueryOrRespond")
    .addConditionalEdges("generateQueryOrRespond", shouldRetrieve)
    .addEdge("retrieve", "gradeDocuments")
    .addConditionalEdges("gradeDocuments", checkRelevance)
    .addEdge("generate", END)
    .addEdge("rewrite", "generateQueryOrRespond"); // 重写后重新尝试检索

  const app = workflow.compile();

  // --- 5. 运行测试 ---
  const inputs = {
    messages: [new HumanMessage("What does Lilian Weng say about types of reward hacking?")],
  };

  console.log("Starting Agent...");
  for await (const output of await app.stream(inputs)) {
    for (const [key, value] of Object.entries(output)) {
      const lastMsg = value.messages[value.messages.length - 1];
      console.log(`\n--- Node: ${key} ---`);
      console.log(lastMsg.content);
    }
  }
}

runAgent().catch(console.error);
```

### 核心流程解析

这个 Agent 的工作流比简单的 RAG 链更健壮，因为它包含了一个**自我修正循环**：

1.  **generateQueryOrRespond**: LLM 接收用户问题。如果它认为需要外部知识，它会生成一个 `tool_call`（调用检索器）。如果不需要，直接结束。
2.  **retrieve**: 执行检索，获取文档内容。
3.  **gradeDocuments**: 这是一个关键的**评估步骤**。另一个 LLM 实例会检查检索到的文档是否真的回答了用户的问题。
    *   **相关 (Yes)** -> 转到 `generate` 节点生成最终答案。
    *   **不相关 (No)** -> 转到 `rewrite` 节点。
4.  **rewrite**: 如果文档不相关，Agent 认为可能是原始问题写得不好。它会重写问题，使其更适合检索，然后跳回第一步重新开始。
5.  **generate**: 使用经过验证的相关文档生成最终回答。

### 关键技术点

*   **Conditional Edges (条件边)**: `shouldRetrieve` 和 `checkRelevance` 函数决定了图的流向，实现了逻辑分支。
*   **ToolNode**: LangGraph 提供的预构建节点，用于自动执行工具调用。
*   **State Management**: 所有的节点都通过共享的 `GraphState` (主要是 `messages` 数组) 进行通信。