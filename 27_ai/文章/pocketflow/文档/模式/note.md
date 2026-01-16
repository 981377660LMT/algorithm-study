## Agent

**Agent 模式一针见血讲解：**

**Agent 本质 = 动态路由 (Dynamic Router) + 状态循环 (State Loop)**

Agent 不再是线性的流水线（Pipeline），而是一个**有“大脑”的循环图谱**。核心在于一个**决策节点（Agent Node）**，它不断观察现状，决定下一步走哪条路，或者回头重做。

### 1. 核心机制：图谱实现 (Graph Implementation)

- **大脑 (The Brain):** **Agent Node**。
  - **输入：** 压缩的上下文 (Context) + 明确的动作空间 (Action Space)。
  - **思考：** LLM 输出 YAML（包含 `thinking` 推理过程 + `action` 决策）。
  - **输出：** 决定控制流走向（是去搜索？是去回答？还是报错？）。
- **身体 (The Body):** **Action Nodes**。
  - 只负责执行具体任务（如 `SearchWeb`），执行完**必须**将结果写回共享状态（Shared Store），并通常**循环回 (Loop Back)** 决策节点。

### 2. 两个成败关键

**A. 上下文管理 (Context Management) —— 少即是多**

- **痛点：** LLM 有“中间丢失”效应（Lost in the Middle），塞太多历史记录会变傻。
- **解法：** 不要把整个 Chat History 扔进去。用 RAG 提取最相关的几条，或者只保留最近的 N 次操作结果。

**B. 动作空间设计 (Action Space) —— 像是设计 API**

- **原子化：** 动作之间不要重叠（不要同时有 `read_csv` 和 `read_db`，应该统一为 `query_data`）。
- **参数化：** 让 LLM 填空。与其给它 10 个固定按钮，不如给它一个写 SQL 的框。
- **可回溯：** 允许 Agent 撤销（Undo）。死胡同不要重开游戏，退一步即可。

### 3. 代码逻辑解析 (从示例看模式)

提供的代码展示了一个经典的 **ReAct (Reason + Act)** 循环：

1.  **`DecideAction` (大脑):**
    - 看 `shared["query"]` 和 `shared["context"]`。
    - LLM 决定是 `"search"` 还是 `"answer"`。
    - 如果是 search，提取 search_term。
2.  **`SearchWeb` (手脚):**
    - 执行搜索。
    - **关键点：** `return "decide"` -> **强行把流程拉回 `DecideAction`**。这就是 Agent 能够自我修正、多步执行的根本原因。
3.  **`DirectAnswer` (出口):**
    - 当大脑认为信息足够时，跳出循环，生成最终答案。

**一句话总结：**
不要写死 IF/ELSE。把 IF/ELSE 的判定权交给 LLM，并通过“循环连接”让 LLM 有机会反复尝试，直到做对为止。

## Workflow

**Workflow（工作流）模式一针见血讲解：**

**Workflow 本质 = 确定性任务分解 (Deterministic Task Decomposition)**

如果说 Agent 是“自动驾驶”，Workflow 就是**“工业流水线”**。它适用于那些**步骤固定、顺序明确**的复杂任务。

### 1. 核心逻辑：拆解的艺术 (Decomposition)

Workflow 的成败完全取决于你怎么切分任务。你需要在“太粗”和“太细”之间找到那个**“黄金平衡点” (Sweet Spot)**：

- **太粗 (Too Coarse):** “写一篇关于 AI 安全的论文”。
  - _后果：_ 任务太复杂，LLM 顾头不顾尾，产出质量低，幻觉多。
- **太细 (Too Granular):** “写第一句”、“写第二句”。
  - _后果：_ LLM 丢失整体上下文，段落之间逻辑不连贯，且调用成本极高。
- **黄金平衡点:** “写大纲” -> “写草稿” -> “润色”。
  - _特点：_ 每个节点输入清晰，输出可控。

### 2. 代码实现模式

Workflow 在代码上体现为**线性链式调用**。

- **节点 (Node):** 每个节点只做一件事（单一职责原则）。
  - `GenerateOutline`: 获取 Topic -> 产出 Outline。
  - `WriteSection`: 获取 Outline -> 产出 Draft。
  - `ReviewAndRefine`: 获取 Draft -> 产出 Final Article。
- **连接 (Connection):** 使用 `>>` 操作符定义确定性的数据流向。
  - `outline >> write >> review`
  - 这意味着：上一步的输出必然是下一步的输入基础（通过 Shared Store 传递）。

### 3. Workflow vs Agent：何时切换？

这是系统设计的关键分界线：

- **使用 Workflow:** 当流程是**线性**的，你完全知道第一步之后一定是第二步（如：文章写作、数据 ETL、每日报表）。
- **切换到 Agent:** 当流程充满**分支**和**边缘情况 (Edge Cases)**。如果你发现你在 Workflow 里写了太多的 `if-else` 来处理异常，或者下一步做什么取决于上一步的结果（如：搜索结果不满意需要重搜），请立即改用 **Agent** 模式。

**一句话总结：**
Workflow 是**刚性**的流水线，负责稳健输出；Agent 是**柔性**的大脑，负责应对变化。不要用 Workflow 处理复杂决策，也不要用 Agent 处理简单流程。

## RAG

**RAG (Retrieval Augmented Generation) 一针见血讲解：**

**RAG 本质 = 外挂大脑 (External Brain) + 实时查阅 (Real-time Lookup)**

LLM 的知识只停留在它训练结束的那一天。RAG 通过把私有数据“喂”给 LLM，解决了两大问题：**知识过时**和**私有数据不可知**。

它在架构上严格分为**冷热两条链路**。

### 1. 离线链路 (Offline Stage) —— 建索引（冷启动）

这是一次性的（或定时的）重型任务。核心是**“把文本变成数字”**。

- **切片 (ChunkDocs):** 把长文档切碎。切多大？看你的 LLM 窗口和语义完整性（通常 500-1000 tokens）。
- **向量化 (EmbedDocs):** **最关键一步**。调用 Embedding 模型，把文本切片转化为高维向量（即“数字列表”）。这让计算机能理解语义相似度。
- **存储 (StoreIndex):** 存入向量数据库（如 FAISS, Milvus, Qdrant）。
- **PocketFlow 特性 (BatchNode):** `ChunkDocs` 和 `EmbedDocs` 继承自 `BatchNode`，意味着框架会自动并行处理大量文件，提高吞吐量。

### 2. 在线链路 (Online Stage) —— 查答案（热请求）

这是用户每次提问都会触发的实时链路。核心是**“语义相似度匹配”**。

- **问题向量化 (EmbedQuery):** 用户的提问也必须变成向量。**注意：必须使用和离线阶段完全相同的 Embedding 模型。**
- **检索 (RetrieveDocs):** 在数据库里找和“问题向量”距离最近的“文档向量”。这是单纯的数学计算（余弦相似度等）。
- **生成 (GenerateAnswer):** 只有这一步才真正调用 LLM（GPT/Claude/Gemini）。
  - **公式：** `Answer = LLM(Prompt + Content + Question)`
  - LLM 此时充当的是**阅读理解者**，而不是知识库。

### 3. 代码模式解析

- **BatchNode 的应用：**
  在离线阶段，`EmbedDocs` 使用了 `BatchNode`。这意味着如果 `Chunks` 有 1000 个，你可以一次性并行处理，或者利用 Python 列表推导式高效映射，而不是写死循环一个个调接口。
- **Context 注入：**
  `GenerateAnswer` 节点展示了 Prompt Engineering 的核心技巧——**上下文注入**：
  ```python
  prompt = f"Question: {question}\nContext: {chunk}\nAnswer:"
  ```
  你并没问 LLM “你知道这事吗？”，而是说“根据这段文字（Context），回答这个问题”。这就是 RAG 能够减少幻觉的根本原因。

**一句话总结：**
RAG 就是考试时的“开卷作弊”：Offline 阶段负责把书整理好做成索引，Online 阶段负责根据考题快速翻到书的那一页，然后抄给老师（用户）看。

## MapReduce

## StructuredOutput

## MultiAgent
