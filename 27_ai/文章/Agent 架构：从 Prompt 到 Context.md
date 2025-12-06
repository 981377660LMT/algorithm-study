# Agent 架构：从 Prompt 到 Context

基于您提供的文章《Agent 架构：从 Prompt 到 Context》，这是一篇非常有深度的技术综述，它揭示了 AI 应用开发从“调教模型（Prompt Engineering）”向“构建系统（Context Engineering）”的范式转变。

以下是对该文章的**深入、详细讲解**，我将其拆解为五个核心逻辑层次，帮助您从架构师的视角理解这一演进。

---

### 核心论点：从“艺术”到“工程”的跨越

文章开篇点明了一个关键分歧：

- **Prompt Engineering (提示工程)**：被视为一种“艺术”或“战术”。它关注如何通过措辞、示例（Few-Shot）和思维链（CoT）来优化单次交互的输出。它的局限在于**脆弱性**（改个词结果就不一样）、**不可扩展**且**无状态**。
- **Context Engineering (上下文工程)**：被视为一种“科学”或“战略”。它不只是写 Prompt，而是构建一个**自动化系统**，负责在正确的时间、以正确的格式，为模型提供正确的信息（Context）。

**一句话总结：Prompt 告诉模型“如何思考”，而 Context 赋予模型完成工作所需的“知识和工具”。**

---

### 第一部分：Context Engineering 的基石——RAG 的进阶

RAG（检索增强生成）不仅仅是给模型挂载知识库，它是 Context Engineering 的核心架构。文章将 RAG 的发展分为了三个阶段：

1.  **Naive RAG（朴素 RAG）**：检索 -> 生成。适合简单问答，但容易查不准。
2.  **Advanced RAG（高级 RAG）**：
    - **检索前**：查询重写（Query Transformation）。
    - **检索后**：重排序（Reranking）、上下文压缩。
3.  **Modular RAG（模块化 RAG）**：这是迈向 Agent 的关键。
    - **CRAG (Corrective RAG)**：带有“自我反思”机制。检索完后，模型先评估“这资料有用吗？”如果没用，触发网络搜索。
    - **Self-RAG**：模型自主决定何时检索，而不是每次都检索。
    - **Agentic RAG**：将 RAG 放入智能体循环中，支持多步推理和工具调用。

**关键概念：Context Stack（上下文技术栈）**
文章提出了一个新的抽象层。以前我们只关注数据库，现在形成了一个完整的栈：

- **数据层**：Vector DB (Milvus, Pinecone 等)。
- **编排层**：LangChain, LlamaIndex。
- **服务层**：Reranking as a Service (Cohere, Jina)。
- **选型建议**：如果追求极致扩展性选 Milvus；如果追求快速上手零运维选 Pinecone。

---

### 第二部分：上下文的“提纯”工艺（如何解决 Lost in the Middle）

这是 Context Engineering 中最硬核的技术细节。LLM 有一个致命缺陷：**Lost in the Middle（中间迷失）**。即当上下文很长时，模型容易忽略中间的信息，只记得开头和结尾。

为了解决这个问题，必须对上下文进行精细化管理：

1.  **高级分块策略 (Chunking)**：

    - 不要只按字符数切分（容易切断语义）。
    - **语义分块**：使用 Embedding 模型检测语义转折点，话题变了再切分。
    - **智能体分块**：让 LLM 自己决定怎么切。

2.  **重排序 (Reranking)**：

    - **双编码器 (Bi-encoder)**：速度快，用于第一轮海量召回（粗排）。
    - **交叉编码器 (Cross-encoder)**：精度高，计算慢。它将 Query 和 Document 同时输入模型进行深度交互，用于第二轮精排。**这是金融/法律等高精度场景的必选项。**

3.  **压缩与摘要**：
    - **过滤**：用 LLM 判断检索回来的文档是否真的相关，不相关直接丢弃。
    - **摘要**：对于长对话历史，不直接塞入 Context，而是让模型生成一个摘要存入记忆。

---

### 第三部分：智能体系统的上下文管理策略

从 **HITL (Human-in-the-Loop)** 转向 **SITL (System-in-the-Loop)**。
以前是人写 Prompt，现在是系统自动生成 Context。LangChain 提出了四个关键策略：

1.  **Write (持久化)**：使用 Scratchpads（草稿本）记录中间思考步骤，使用 Memory 记录长期偏好。
2.  **Select (检索)**：动态路由。不仅检索知识，甚至可以检索“工具”（当工具太多时，只把相关的工具定义放入 Context）。
3.  **Compress (压缩)**：防止 Token 溢出。
4.  **Isolate (隔离)**：
    - **多智能体**：将复杂任务拆解，每个子智能体只维护自己那部分“小而精”的上下文，避免干扰。
    - **沙盒环境**：工具执行的结果（比如巨大的 JSON）不直接给 LLM，只返回关键结果。

---

### 第四部分：Agent 架构与工作流编排

这是文章最高阶的部分，讨论了数据如何在系统中流动。

#### 1. Workflow vs. Agent

- **Workflow (工作流)**：类似“专家系统”。路径是写死的（A -> B -> C）。适合风控、合规等`容错率低`的场景。
- **Agent (智能体)**：路径是动态的。LLM 自主决定下一步干什么。适合`开放式问题`。
- **混合模式**：宏观上是 Workflow，微观节点上是 Agent。

#### 2. 核心设计模式

- **Prompt Chaining**：链式调用，一步接一步。
- **Routing (路由)**：LLM 充当路由器，分析意图后分发给不同的下游模块。
- **Orchestrator-Workers (编排器-工人类)**：一个“包工头”负责拆解任务，分发给多个“专才”去执行，最后汇总。

#### 3. 决策机制：ReAct 与 Reflection

- **ReAct (Reason + Act)**：
  - 循环：**思考 (Thought)** -> **行动 (Action)** -> **观察 (Observation)** -> **再思考**。
  - 这是 Agent 能够自主使用工具的根本逻辑。
- **Reflection (反思)**：
  - 执行完动作后，增加一个“评估器”模块。
  - 如果评估结果不好，Agent 会自我修正，重新规划。这是 DeepResearch 类应用的核心。

#### 4. 实现工具：LangGraph

文章强烈推荐 **LangGraph**。它用“图论”的概念来定义 Agent：

- **State (状态)**：所有节点共享的内存（数据总线）。
- **Nodes (节点)**：执行具体任务的函数（思考、调用工具）。
- **Edges (边)**：定义流转逻辑（条件边实现路由）。
- **Checkpointer**：保存状态，支持“断点续传”和人工介入。

---

### 第五部分：未来展望

Context Engineering 的终极目标是**GraphRAG** 和 **自主智能体**。

- **GraphRAG**：利用知识图谱解决“实体间关系”的推理问题（例如：A 是 B 的父亲，B 是 C 的同事，问 A 和 C 的关系）。
- **终极形态**：`Context Engineering 本质上是一种“补偿机制”，用来弥补当前 LLM 记忆力（Context Window）和逻辑力的不足。随着模型进化，这种外部支架可能会逐渐内化为模型的能力。`

### 总结

这篇文章的核心价值在于它**拔高了 AI 开发的维度**：
不要再沉迷于打磨一句完美的 Prompt 了。作为工程师，你的精力应该放在设计**数据流 (Data Flow)**、构建**上下文栈 (Context Stack)** 以及编排**智能体工作流 (Agent Workflow)** 上。这才是构建企业级、高可靠 AI 应用的正确路径。
