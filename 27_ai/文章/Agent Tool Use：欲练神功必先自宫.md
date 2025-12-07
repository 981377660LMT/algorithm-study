# [Agent Tool Use：欲练神功必先自宫](https://blog.zihanjian.com/article/27c9db82-fa86-807e-821e-c8911d112faa)

这篇文章的核心观点非常犀利，它借用“欲练神功必先自宫”的比喻，指出**为了构建生产级、高可靠的 Agent 系统，开发者必须放弃（或限制使用）LLM 原生提供的、看似方便的 Function Calling（函数调用）黑盒模式，转而采用“解耦式架构”（Plan-then-Execute, PtE）。**

以下是对这篇文章的详细剖析和讲解：

### 1. 核心冲突：原生 Function Calling (FC) vs. 解耦式 PtE

- **原生 Function Calling (FC):**

  - **模式:** `Think -> Act -> Observe -> Think -> Act ...`
  - **特点:** 这是 OpenAI 等厂商推崇的模式。LLM 自己决定下一步调什么函数，执行完后再决定下一步。
  - **问题:** 这是一个“黑盒”。你很难控制它何时停止，很难并行执行任务（因为它通常是串行的），且一旦中间出错，整个链条容易崩溃。作者称之为“单体式 Agent”。

- **解耦式 PtE (Plan-then-Execute):**
  - **模式:** `Planner (LLM) -> Plan (JSON) -> Executor (Code) -> Result`
  - **特点:** 将“大脑”（规划）和“手脚”（执行）分开。
    - **Planner:** 只负责生成一个完整的计划（比如：第一步查天气，第二步查机票，这两步可以并行）。
    - **Executor:** 这是一个确定性的代码程序（非 LLM），负责解析计划、并发调用工具、处理重试和错误。
  - **优势:** 可控、可调试、支持高并发、更安全。

### 2. 决策框架：什么时候用什么？(4.1 节解读)

作者给出了一个非常实用的技术选型指南：

| 维度           | **Function Calling (FC)** | **Plan-then-Execute (PtE)**   |
| :------------- | :------------------------ | :---------------------------- |
| **项目阶段**   | 快速原型验证 (Demo)       | 生产环境 (Production)         |
| **任务复杂度** | 简单、单回合 (查天气)     | 复杂、多步骤、有依赖关系      |
| **安全性**     | 低风险 (内部工具)         | 高风险 (处理敏感数据、资金)   |
| **性能要求**   | 对延迟不敏感              | 高并发、低延迟 (需要并行处理) |
| **工具数量**   | 少 (<20 个)               | 大量工具                      |
| **可控性**     | 依赖模型“幻觉”            | 行为必须可预测、可调试        |

**一句话总结：** 如果你是做 Demo 或者简单的 Chatbot，用 FC；如果你是做企业级业务流程自动化，必须用 PtE。

### 3. 实施 PtE 的最佳实践 (4.2 节解读)

如果你决定采用 PtE 架构，作者给出了具体的工程建议，这与我们之前讨论的 **LangGraph** 的设计理念不谋而合：

1.  **定义形式化的计划模式 (Formal Plan Schema):**

    - *不要*让 LLM 输出一段自然语言说“我打算先做 A 再做 B”。
    - *要*强制 LLM 输出结构化的 **JSON/YAML**。定义好字段：`task_id`, `dependencies`, `tool_name`, `arguments`。这相当于定义了 Planner 和 Executor 之间的 API 接口。

2.  **构建健壮的执行器 (Resilient Executor):**

    - Executor 是纯代码逻辑。它必须能处理 HTTP 超时、API 报错。
    - **关键点:** 单个工具失败不应导致整个 Agent 崩溃。Executor 应该捕获错误并决定是重试还是报告给 Planner。

3.  **隔离规划器 (Isolate the Planner):**

    - **零信任原则:** 不要相信 LLM 输出的计划是安全的。在 Executor 执行之前，必须校验参数（例如：防止 SQL 注入，防止删除系统文件）。

4.  **拥抱异步与并行 (Embrace Asynchronicity):**

    - 这是 PtE 最大的优势。如果计划显示 Task A 和 Task B 没有依赖关系，Executor 应该同时启动它们，而不是像 FC 那样傻傻地排队。

5.  **实现重规划循环 (Re-planning Loop):**
    - 计划不是一成不变的。如果 Executor 发现环境变了或者工具彻底失败了，它需要把错误信息反馈给 Planner，让 Planner 生成一个新的计划（Re-plan）。

### 4. 深度思考：为什么说这是“技术懒惰”？(4.3 节解读)

作者在最后提出了一个非常有深度的观点：

- **现状:** 很多开发者直接使用 OpenAI 的 Assistants API 或 LangChain 的 `AgentExecutor`，因为这最简单，代码最少。
- **批判:** 作者认为这是一种“技术懒惰”。因为这种架构在面对复杂业务时极其脆弱。
- **未来:**
  - **LLM = CEO (战略):** 负责高层决策、意图理解、生成计划。
  - **Code/System = COO (战术):** 负责具体的调度、执行、状态管理、容错。
- **结论:** 虽然未来 AGI 可能强大到不需要解耦，但在通往 AGI 的路上，**解耦（PtE）是必经的“弯路”**。现在就依赖全自动的 FC 是不负责任的。

### 5. 与 LangGraph 的联系

这篇文章实际上是在为 **LangGraph** 这种架构背书。

- **LangChain 的旧版 Agent** 大多是基于 Function Calling 的单体循环，容易陷入死循环，且难以并行。
- **LangGraph** 允许你显式地构建 PtE 架构：
  - 你可以创建一个 `Planner Node`（调用 LLM 生成 JSON）。
  - 创建一个 `Executor Node`（解析 JSON 并并行执行工具）。
  - 通过 `Conditional Edges` 决定是结束还是回到 `Planner Node` 进行重规划。

**总结：** 这篇文章告诫我们，不要迷信 LLM 的全能性。在构建严肃应用时，`要用传统的软件工程思想（解耦、接口定义、错误处理）去约束和驾驭 LLM，而不是把控制权全盘交给它。`
