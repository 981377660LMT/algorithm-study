# [What makes Claude Code so damn good](https://minusx.ai/blog/decoding-claude-code/#31-llm-search---rag-based-search)

这篇文章是对 **Claude Code (CC)** 的一次深度逆向工程和架构剖析。作者（MinusX 团队）通过拦截和分析 Claude Code 的网络请求，揭示了它为何比许多复杂的 Agent 框架更好用。

核心结论令人惊讶：**Claude Code 的成功不在于复杂的架构，而在于极致的简单性、对小模型的巧妙运用以及对 RAG（检索增强生成）的摒弃。**

以下是对这篇文章的深度详细剖析：

### 1. 核心哲学：KISS 原则 (Keep It Simple, Dummy)

文章反复强调一个观点：**拒绝过度设计**。
目前业界流行复杂的“多智能体协作”（Multi-Agent Systems）、复杂的图结构（Graph）和自动切换。但作者指出，这些复杂性让调试变得极其困难（难了 10 倍）。

- **Claude Code 的选择：** 单一主循环，单一消息历史，没有复杂的 Agent 握手流程。
- **启示：** 在构建 Agent 时，先在一个文件里写死循环，直到它跑不通了再考虑拆分。

### 2. 架构设计：控制流 (Control Loop)

Claude Code 的控制流设计非常“反直觉”，因为它太简单了：

- **单一主线程 (One Main Loop):** 尽管多智能体很火，CC 坚持使用单线程。它维护一个扁平的消息列表。
- **有限的分支:** 只有在处理极其复杂的任务时，它才会生成一个“子 Agent”（Clone 自身），但这个子 Agent 不能再生成孙子 Agent（最大深度为 1）。
- **自我管理的 Todo List:** 这是防止 Agent “变蠢”的关键。
  - CC 维护一个由**模型自己管理**的 Todo List。
  - 这解决了长上下文中的“遗忘”问题（Context Rot）。模型每一步都会参考这个列表，决定是继续、划掉任务还是新增任务。这是一种动态的“思维链”（Chain of Thought）。

### 3. 模型策略：大小模型搭配 (The Haiku Strategy)

这是降低成本和提高速度的商业机密：

- **50% 以上的调用使用的是 Claude 3.5 Haiku (小模型)。**
- **分工明确：**
  - **Sonnet (大模型):** 负责核心逻辑推理、写代码、做决策。
  - **Haiku (小模型):** 负责“脏活累活”——读取大文件、解析网页、处理 Git 历史、总结长对话、生成简单的标签。
- **启示：** 不要把所有任务都扔给 GPT-4 或 Claude 3.5 Sonnet。构建 Agent 时，必须实现模型路由，让便宜的模型处理上下文这一层。

### 4. 工具与搜索：LLM Search >>> RAG

这是文章中最具颠覆性的观点：**在代码场景下，传统的 RAG（向量检索）是糟糕的。**

- **为什么 RAG 不行？**
  - 代码逻辑是精确的，向量相似度是模糊的。
  - RAG 引入了隐藏的故障模式（切片切坏了、检索召回率低、重排错误）。
- **Claude Code 怎么做？**
  - **LLM Search:** 它像人类程序员一样，使用 `grep` (ripgrep), `find`, `ls` 等工具。
  - **原理:** LLM 非常懂代码结构。它能写出极高精度的正则表达式来定位代码块。它先看文件列表，再 `grep` 关键词，再 `read` 具体行数。
  - **优势:** 这种方式是确定性的，且利用了模型本身的智力，而不是依赖外部的黑盒检索器。

### 5. 提示词工程 (Prompt Engineering)

Claude Code 的 System Prompt 长达 2800 tokens，工具定义长达 9400 tokens。其中的技巧包括：

- **`claude.md` (上下文文件):**
  - 这是用户偏好的核心载体。类似于 Cursor 的 `.cursorrules`。
  - CC 每次请求都会带上这个文件的内容，强制模型遵守用户的编码规范（如“不要用 TypeScript，只用 JS”）。
- **XML 标签与 Markdown:**
  - 大量使用 XML 标签（如 `<system-reminder>`, `<good-example>`, `<bad-example>`）来结构化提示词。
  - 使用 Markdown 标题区分不同板块（风格、任务管理、工具策略）。
- **"Shouting" (大写强调):**
  - 虽然听起来很土，但 `IMPORTANT`, `MUST`, `NEVER` 依然是控制 LLM 行为最有效的手段。
  - 例如：`IMPORTANT: DO NOT ADD ANY COMMENTS unless asked`。
- **算法化提示 (Write the Algorithm):**
  - 不要只告诉模型“做什么”，要告诉它“怎么做”（算法）。
  - 在 Prompt 中写出伪代码或流程图逻辑，告诉模型在遇到 A 情况时判断 B，然后执行 C。

### 6. 工具设计：高低搭配

- **低级工具:** `Bash`, `Read`, `Write`。赋予模型最大的灵活性。
- **高级工具:** `WebFetch`, `mcp__ide__getDiagnostics`。对于确定性很强的任务（如获取报错信息），封装成高级工具，避免模型手动敲命令出错。
- **策略:** 常用且容易出错的操作封装为高级工具；不常用或需要灵活性的操作保留低级工具（Bash）。

### 总结：如何复刻 Claude Code 的魔力？

如果你想构建一个好用的 Agent，文章给出的建议是：

1.  **保持单线程循环**，别搞复杂的多智能体图。
2.  **让 Agent 自己维护 Todo List**，作为它的短期记忆和规划器。
3.  **放弃代码 RAG**，教模型使用 `grep` 和 `find`。
4.  **大量使用小模型** (Haiku/Flash) 来处理阅读和总结任务。
5.  **Prompt 要写得像算法说明书**，包含大量正反示例（Few-Shot）。
6.  **支持 `claude.md` / `rules` 文件**，让用户能控制 Agent。

这篇文章实际上是在告诉我们：**现阶段的 Agent 胜利，属于工程化细节的胜利，而不是架构复杂度的胜利。**
