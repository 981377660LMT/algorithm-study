# 也许是 Coding Agent 的参考答案 — Claude Code 深度解析

https://bytedance.larkoffice.com/docx/XszbdUrjJoHAhLxGgfwcBi3pnPS

这是一个非常深入且切中当前 AI 编程前沿的话题。Claude（特别是 Claude 3.5 Sonnet/Opus）目前在编程领域表现卓越，其背后的工程实践、Prompt 技巧以及架构设计非常值得借鉴。

### 1. Claude 如何用大小模型精准协作？ (Orchestration & Routing)

在复杂的编程任务中，单纯依赖最强模型（如 Opus）成本过高且速度慢，单纯依赖小模型（如 Haiku）智力不足。精准协作的核心在于 **“路由（Routing）”** 和 **“分层处理”**。

- **意图识别与路由 (The Router):**
  - 前端通常会有一个轻量级模型（或专门的分类器），用于分析用户的 Prompt。
  - 如果任务是简单的“解释这段代码”或“补全一行”，路由给小模型（Haiku/Flash）。
  - 如果任务是“重构架构”或“解决复杂 Bug”，路由给大模型（Sonnet/Opus）。
- **上下文预处理 (Context Pre-processing):**
  - 小模型可以作为“秘书”。它先扫描巨大的代码库，提取出相关的摘要、函数签名或文件列表。
  - 它将这些“精炼后”的上下文喂给大模型，从而节省大模型的 Context Window 和推理成本。
- **监督与修正 (Oversight):**
  - **Draft-Verify 模式**：小模型快速生成代码草稿，大模型负责 Code Review 和逻辑修正。或者反过来，大模型制定高层计划（Plan），小模型负责填充具体的样板代码（Implementation）。

### 2. Claude Code 如何设计 MultiAgent 系统？

Claude 的超长上下文（200k+）和对 XML 标签的极佳遵循能力，使其非常适合构建 MultiAgent 系统。

- **角色定义 (Role Definition):**
  - 系统通常包含多个专门的 Agent：`Planner`（规划者）、`Coder`（编码者）、`Reviewer`（审查者）、`Executor`（执行者/工具调用者）。
- **共享状态 (Shared State / Blackboard):**
  - 不同于通过 API 互相发消息，Claude 的 MultiAgent 往往共享一个不断更新的 `Context`（黑板模式）。
  - 例如：`Planner` 将步骤写入 `<plan>` 标签，`Coder` 读取计划并写入 `<code_change>`，`Executor` 运行测试并将结果写入 `<test_output>`。
- **工具使用 (Tool Use / MCP):**
  - Claude Code 的核心是 **Model Context Protocol (MCP)** 的思想。Agent 不仅仅是生成文本，而是被赋予了“手脚”（读取文件、运行 Shell、Git 操作）。
  - Agent 的循环通常是：`思考 (Thinking)` -> `调用工具 (Tool Call)` -> `观察结果 (Observation)` -> `修正 (Refinement)`。
- **反思循环 (Reflection Loop):**
  - 当测试失败时，系统不会立即停止。`Reviewer` Agent 会分析错误日志，将错误原因反馈给 `Coder`，形成一个自动修复的闭环。

### 3. 它的 Prompt 设计有什么关键？

Claude 的 Prompt Engineering 与 GPT 系列略有不同，它极其依赖 **结构化数据**。

- **XML 标签隔离 (XML Tagging):**
  - 这是 Claude Prompt 的灵魂。将不同类型的信息用标签包裹，例如 `<context>`, `<instructions>`, `<code_style>`, `<examples>`。
  - 这能防止指令注入，并帮助模型明确区分“参考资料”和“执行命令”。
- **Chain of Thought (CoT) 显式引导:**
  - 强制模型在输出代码前先进行思考。Prompt 中会包含：`Before writing code, please think step-by-step inside <thinking> tags.`
  - 这能显著提高复杂逻辑代码的准确率。
- **System Prompt 的角色沉浸:**
  - 不仅仅是 "You are a coder"，而是详细定义能力边界："You are an expert in Rust and TypeScript. You prefer functional programming patterns. You must handle all edge cases."
- **One-Shot / Few-Shot Examples:**
  - 在 Prompt 中提供一个 `<example>`，展示“输入”和理想的“输出”格式，特别是对于特定的代码风格或 JSON 结构，效果立竿见影。

### 4. 如何不基于 Embedding 做代码的检索？

Embedding（向量检索）在代码领域有缺陷：它擅长语义相似（"猫" vs "狗"），但不擅长精确匹配（变量名 `user_id` vs `userId`）或结构依赖。

不基于 Embedding 的检索通常结合了 **传统搜索** 和 **智能代理**：

- **基于 LSP (Language Server Protocol) 的符号跳转:**
  - 这是 IDE 的原生能力。通过静态分析，构建引用图（Reference Graph）。
  - 当 AI 需要修改函数 `A` 时，系统自动检索 `A` 的定义、所有调用 `A` 的地方、以及 `A` 调用的其他函数。这是最精准的上下文。
- **基于 AST (抽象语法树) 的结构化搜索:**
  - 不把代码看作文本，而是看作树。搜索“所有继承自 `BaseController` 的类”或“所有使用了 `useEffect` 的文件”。
- **Agentic Search (代理式搜索):**
  - 这是目前最先进的做法（也是 Claude Code 类似工具的做法）。
  - 不预先建立索引。而是给 AI 一个 `grep` 或 `find` 工具。
  - AI 像人类一样：先 `ls` 看目录 -> 猜测文件名 -> `grep` 搜索关键词 -> `read_file` 读取内容。如果没找到，它会换个关键词再搜。这种方式最灵活，能找到 Embedding 找不到的隐式逻辑。
- **BM25 / 关键词加权:**
  - 传统的全文检索，对变量名、特定错误码的匹配效果远好于向量检索。

### 5. 如何把 Claude Code 用到极致？（通用 AI 编程方法论）

要将 AI 编程工具用到极致，不能把它当作“搜索引擎”，而要把它当作一个“刚入职的高级工程师”。

- **上下文管理 (Context Curation):**
  - **不要**把整个项目扔进去。
  - **要**手动挑选相关文件（`@File`）。如果你要改 `auth.ts`，请务必同时把 `user_model.ts` 和 `auth_interface.ts` 喂给它。精准的上下文 = 精准的代码。
- **测试驱动开发 (TDD) 增强:**
  - 先让 AI 写测试用例（Unit Test），运行失败。
  - 再让 AI 写实现代码，直到测试通过。
  - 这种流程能极大降低 AI 产生“幻觉代码”的概率，因为有测试作为客观的验证标准。
- **迭代式对话 (Iterative Refinement):**
  - 不要期望 One-shot 完美。
  - 第一轮：“生成一个大概的方案”。
  - 第二轮：“优化一下错误处理”。
  - 第三轮：“把变量命名改成符合我们项目规范的风格”。
- **利用 Lint/Compiler 反馈:**
  - 把 IDE 的报错信息（TypeScript 错误、Rust 借用检查错误）直接复制给 AI。Claude 非常擅长根据编译器报错修复代码。
- **生成文档与解释:**
  - 对于复杂的遗留代码（Legacy Code），先让 AI 生成一份 `<explanation>` 或伪代码，确认它理解正确了，再让它动手修改。

**总结：**
Claude 及其同类工具的强大在于 **Context（上下文）** 和 **Reasoning（推理）** 的结合。通过结构化的 Prompt（XML）、代理式的工具使用（Agentic Workflow）以及精准的上下文投喂，你可以让它从一个简单的“代码补全工具”变成一个真正的“架构师级助手”。

---

## 简介

Claude Code 的设计理念可以概括为：**“极简的架构 + 极致的上下文工程”**。它不追求复杂的 Multi-Agent 编排，而是通过对模型能力的深度压榨，定义了 Coding Agent 的“参考答案”。

以下是其核心设计理念的深度拆解：

### 1. 架构哲学：薄 UI，厚执行 (Thin UI, Thick Execution)

Claude Code 没有采用传统的复杂 Agent 框架，其核心是一个 **“大模型-工具”主循环**。

- **TUI (Terminal UI)**：`选择终端而非 IDE 插件，是为了获得最高权限的系统控制力（文件、Git、进程）。`
- **执行层**：本质上是一个自主的 React 循环。模型通过 `tool_use` 决策，系统执行后返回 `tool_result`，循环往复直到任务完成。
- **三种模式**：
  - **Manual**：人类作为“安全阀”。
  - **Auto**：追求效率的“自动驾驶”。
  - **Plan**：纯思维链（CoT）推演，不产生副作用，用于架构设计。

### 2. 多模型协作：Haiku 护航，Sonnet 攻坚

Claude Code 巧妙地利用了不同尺寸模型的成本与速度优势：

- **Haiku (小模型)**：负责“边缘计算”。包括 TUI 渲染、对话标题生成、任务摘要、请求预处理。这保证了交互的丝滑感，且不占用主模型的上下文配额。
- **Sonnet/Opus (大模型)**：负责“核心推理”。专注于代码逻辑理解、复杂工具调用和长程规划。

### 3. 上下文工程 (Context Engineering) 的极致优化

这是 Claude Code 区别于普通“套壳”产品的核心护城河：

- **不依赖 Embedding 的检索**：它更倾向于通过 `ls -R`、`grep` 和 `read_file` 让模型自主探索，而不是预先通过向量数据库切片。这种“模型自主检索”比静态向量检索更具语义准确性。
- **渐进式披露 (Progressive Disclosure)**：通过 **Skill** 机制，初始只提供能力索引，当模型需要时才加载详细的 `SKILL.md`。这极大地节省了 Token，同时保持了高精度的指令遵循。
- **项目记忆 (CLAUDE.md)**：将项目规范、技术栈偏好等持久化在本地文件中，作为“外部系统提示词”，实现了零成本的个性化定制。

### 4. 确定性与随机性的平衡

Claude Code 意识到 Agent 的决策是随机的，因此它在工具端提供了极强的**确定性反馈**：

- **Bash 工具**：允许模型直接运行测试、编译命令，通过真实的错误信息（而非模型预测的错误）来驱动修复循环。
- **文件操作**：提供原子化的读写工具，确保代码修改的精确度。

### 5. 总结：它是“船”而非“柱子”

正如你之前提到的隐喻，Claude Code 是那艘**顺着 Agent 能力涨潮而升起的“船”**：

- 它不试图通过复杂的硬编码逻辑来“修补”模型。
- 它提供了一套完美的“外骨骼”（工具集和上下文管理），让模型能力（水）能够无损地转化为生产力。

**核心启示**：最好的 Coding Agent 往往不是逻辑最复杂的，而是最能让模型“看清”代码现状、并能“无碍”操作环境的。

## Context 构成

1. System Prompt
   System Prompt 中主要为静态内容，其中有这几个重点：

   - 主动性（Proactiveness）：平衡好“执行任务”和“不让用户意外”的关系，猜测主要是因为 Claude 3.7 系列模型经常一下子做的太多，过分满足用户意图。
   - 遵循现有规范：优先理解并遵守现有的编码风格和使用的库，所以 Claude Code 明显可以发现更能融入已有的项目，新代码和原有代码非常契合。
   - 重点提及不要加太多注释（看来确实被一行注释一行代码的 ai coding 风格折磨了。。。）
   - 有一个 TODO 系统进行内部任务的规划和执行，现在也比较常见，类似 Manus，Deep Research 的复现也都有

   值得注意的是 System Prompt 动态注入的信息非常少，有且只有：

   - env：自动注入 prompt 的用户环境信息，主要包含日期、当前目录位置等
   - git：若判断为 git 仓库，也会自动注入最近的 git commit message，因此 claude code 对 git 信息处理非常强大，能敏锐感受到最近的变更

   应该是为了 prefix caching 的效率，将大量动态内容全部放到 user prompt 中
   更多关于前缀缓存的内容可以参考 manus 的 context engineering 介绍（https://manus.im/blog/Context-Engineering-for-AI-Agents-Lessons-from-Building-Manus）
   简单可以理解为，**若 prompt 的前缀完全相同，推理成本会大幅降低（约 1/10），首 Token 延迟会大幅缩短**

2. User Prompt
   为了补充 System Prompt 的简单，Claude Code 使用了大量的 Inplace Prompt ：
   在用户实际对话的 query 中穿插系统生成的 prompt
   以下方一个 case 为例：

   - 黄色：系统插入的 inplace prompt
   - 绿色：用户输入的内容
   - 剩余内容为工具输入和输出，可以暂时忽略，下方会详细分析

   把最需要模型注意的内容”尽量放在末尾予以强调，避免内容太长导致 LLM 失焦以往前文

   <system-reminder > 的内容和 user message 并无差别，和加权的做法类似，只是 xml 标签这种强格式更容易引起 Attention，理论上和 Claude Code 写 IMPORTANT, VERY IMPORTANT 并无差异 ...

3. Context Compact
   由于代码内容极大，非常容易超过上下文 200k，Claude Code 提供了自动压缩上下文的能力
   总结会把过去的历史消息截断，使得前缀缓存的命中率较高。

   - 总结的时机为条件判断，保证两次总结间的大模型对话，能够命中前缀缓存
   - 反例是常见的“滑动窗口形式”，每轮对话都过滤最老的一轮消息，该方法会让缓存无法命中
   - 删除的方式是：将工具的输出由原始内容变为[Old tool result content cleared]

4. 自定义指令

## 工具、应用场景

Code Rag 真的没有前途了吗？
个人的判断不是。Claude Code 也有自己的立场

1. 代码没有很好的 embedding 模型
2. Rag 需要有 Index 过程，工具使用成本更高
3. 用户也有合规的考虑，不愿意上传代码仓库

---

在创始人的采访中也明确表示过，并没有采用任何的 embedding based 方法进行代码召回
并且原因是：他们内部尝试过 embedding 的方法，但是效果不如规则匹配好
source：https://www.youtube.com/watch?v=zDmW5hJPsvQ
**而并行子 agent 的方式很好的弥补了顺序执行规则匹配耗时上的劣势，大大提升了代码检索效率。**
When you are doing an open ended search that may require multiple rounds of globbing and grepping, use the Agent tool instead （摘录自 Glob 工具的描述）

多个 LLM 并行 grep + ls 各种工具调用。

- Plan Mode
  Agent 执行错误很多时候不是因为模型能力不行，而是人类的意图没有很好的表达。
  规划模式有将信息对齐步骤前置，通过多次对话，将用户的意图同步给 Agent。
  我个人的经验是，复杂需求甚至需要 plan mode 后写下一个 doc，再执行该 doc。大方向对齐后，可以放心让 Agent 在背后运行超长时间（15min+）

---

Claude Code 的核心理念确实体现了“The Simple Thing that Works”（简单即有效）的工程哲学，避开了过度复杂的 RAG 架构，专注于上下文窗口的利用。

### 1. Claude 如何用大小模型精准协作？

Claude Code 采用了典型的**路由器（Router）模式**来协调不同能力的模型，既节省成本又保证质量。

- **轻量级模型（Stubber/Planning）：**
  在此阶段，通常使用较小、响应更快的模型（如 Claude 3.5 Haiku 或 Sonnet）。它们的任务不是编写最终代码，而是进行“瘦身”和“规划”。

  - **文件修剪：** 当需要把大量文件放入上下文时，小模型负责移除不相关的函数体、注释或测试数据，生成代码的“骨架”（Stub）。这使得上下文能容纳更多文件的结构信息。
  - **初步意图识别：** 快速判断用户的请求是需要搜索、运行测试还是仅仅是询问。

- **重量级模型（Coder/Executor）：**
  当真正需要进行复杂的逻辑推理、生成具体的代码补丁或修复 Bug 时，系统会切换到最强的模型（Claude 3.5 Sonnet 或 Opus）。
  - **精准打击：** 大模型接收经过“瘦身”的精简上下文以及具体的指令，生成的代码准确率极高。

### 2. Claude Code 如何设计 MultiAgent 系统？

最好不讨论`multiagent`

不是 multiagent(去中心化)，而是中心化

Claude Code 的 Agent 设计并非复杂的网状拓扑，而是一种**线性的工具链循环（Loop of Tools）**。

- **REPL 循环模式：** 它本质上是一个 Read-Eval-Print Loop。
  - Agent 观察当前环境。
  - Agent 决定调用一个工具（如 `grep`, `cat`, `symbol_search`）。
  - 工具返回结果。
  - Agent 根据结果决定下一步（继续搜索或开始写代码）。
- **显式状态管理：** Agent 维护一个对话历史，其中不仅包含用户的 Prompt，还严格记录了之前的工具调用结果。如果上一步检索失败，Agent 会在下一个 Prompt 中明确修正搜索策略，而不是盲目重试。

### 3. 它的 Prompt 设计有什么关键？

关键在于**防御性 Prompting** 和 **上下文无关性训练**。

- **XML 标签的严格使用：** Claude 模型被微调为对 XML 结构极其敏感。系统 Prompts 通常会要求模型把思考过程包裹在 `<thinking>` 中，把调用的工具包裹在 `<tool_use>` 中。这实际上强制模型在行动前先进行“思维链”（Chain of Thought）推理。
- **One-Shot/Few-Shot 示例：** 在 System Prompt 中并没有塞入海量规则，而是给出了几个完美的“交互示例”。例如，展示一个“搜索->索引->读取文件->修改代码”的标准流程。
- **权限边界明确：** Prompt 中明确规定了 Agent _不能_ 做什么（例如不能在没有读取文件内容的情况下直接猜测修改），这种否定式的指令对于减少幻觉至关重要。

### 4. 如何不基于 Embedding 做代码的检索？

这是 Claude Code 最具“非共识”也是最有效的一点：**暴力 grep 优于向量检索（Embedding）**。

在代码场景中，Embedding 往往会丢失精确的符号信息（比如变量名 `user_id` 和 `user_idx` 在向量空间可能很近，但在代码逻辑中天差地别）。Claude Code 的做法是：

- **基于 ripgrep 的即时搜索：** 它不预先建立庞大的向量数据库，而是通过 Agent 调用命令行工具（类似 `grep -r "functionName"`）。
- **符号表索引（CTags 类似物）：** 它会快速生成一个轻量级的符号映射表（类/函数定义所在的文件路径），完全基于文本匹配。
- **LLM 驱动的阅读：**
  1.  先搜索文件名或关键字。
  2.  LLM 看到文件名列表，决定读哪几个。
  3.  LLM 读完文件内容，自行提取关联信息。
      _原理：_ 利用 LLM 巨大的上下文窗口（200k+），直接把“搜索”变成“阅读理解”问题，而不是“相似度匹配”问题。

### 5. 如何把 Claude Code 用到极致？

无论是 Claude Code 还是 Cursor/Windsurf，以下方法能显著提升效果：

- **像对人一样说话（Context Dump）：**
  不要只说“修复 bug”，要说：“我在运行 `npm test` 时遇到了 X 错误，我怀疑是 `auth.ts` 里最近改的逻辑导致的，请先通过 grep 确认该逻辑在哪里被调用。”
- **手动提供上下文（Pinning）：**
  如果你知道涉及哪些文件，显式地告诉工具。虽然 AI 能检索，但直接给它 `@FileA` 和 `@FileB` 能省去它两轮的探索步骤，让它把 Token 预算全用在写代码上。
- **要求行动计划（Plan First）：**
  强制 AI 在写代码前输出：“我打算先做 A，再做 B”。如果计划不对，立刻打断。这避免了 AI 写了一大堆错误代码后你才发现方向错了。
- **保持 Session 干净：**
  如果一个对话太长，上下文会被截断或产生很多噪音。解决一个具体任务后，立刻开启新的 Chat Session。

---

**总结：** Claude Code 的成功不在于它发明了什么新的深奥算法，而在于它承认了 LLM 上下文窗口越来越大的趋势，**果断抛弃了为了节省 Token 而设计的复杂检索系统（RAG）**，回归了最符合人类直觉的“搜索-阅读-修改”工作流。

---

system-reminder 的内容和 user message 并无差别，和加权的做法类似，只是 xml 标签这种强格式更容易引起 Attention，理论上和 Claude Code 写 IMPORTANT, VERY IMPORTANT 并无差异 ..
