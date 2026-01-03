# 也许是 Coding Agent 的参考答案 — Claude Code 深度解析

https://bytedance.larkoffice.com/docx/XszbdUrjJoHAhLxGgfwcBi3pnPS

1. Claude 如何用大小模型精准协作？
2. Claude Code 如何设计 MultiAgent 系统？
3. 它的 Prompt 设计有什么关键？
4. 如何不基于 Embedding 做代码的检索？
5. 如何把 Claude Code 用到极致？（方法同样适用于其他 ai 编程产品）

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

## 感想
