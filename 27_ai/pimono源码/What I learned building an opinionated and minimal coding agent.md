# [What I learned building an opinionated and minimal coding agent](https://mariozechner.at/posts/2025-11-30-pi-coding-agent/) 解读

> 作者：Mario Zechner，游戏引擎 libGDX 的创造者，资深开发者/教练/演讲者。
> 发表于 2025-11-30。本文是他构建自己的极简编码代理（coding agent）工具 **pi** 的完整心路历程和技术总结。

文章最核心的几个洞察：

**1. 架构分层清晰**：pi 分为 4 层——`pi-ai`（统一 LLM API）→ `pi-agent-core`（代理循环）→ `pi-tui`（终端 UI）→ `pi-coding-agent`（CLI 壳），每层职责单一。

**2. "用文件代替功能"哲学**：是全文最独特的设计模式。内置 todo → `TODO.md`，计划模式 → `PLAN.md`，MCP → CLI+README，后台进程 → tmux，子代理 → bash 生成新实例。核心理念：**不在代理内部管理状态，让文件系统成为真实数据来源（Single Source of Truth）**。

**3. 极简提示 + 极简工具的可行性**：系统提示不到 1000 tokens，只有 4 个工具（read/write/edit/bash），但在 Terminal-Bench 2.0 上击败了许多更复杂的工具。原因是前沿模型经过大量 RL 训练，天生理解编码代理应该做什么。

**4. 可观测性 > 一切**：这是贯穿全文的第一主题。子代理不可观测所以不要、后台 bash 不可观测所以用 tmux、Claude Code 的 plan mode 不可观测所以用文件——所有设计决策的底层逻辑都是"我能看到发生了什么吗？"

**5. 上下文工程的务实观**：跨提供商的上下文切换（Claude → GPT → Gemini）、MCP 的 token 浪费（13-18k tokens 工具描述）、子代理的上下文传递损失——作者对"什么进入上下文窗口"有极致的控制欲。

## 一、文章背景与动机

### 1.1 作者的 LLM 编码演化史

作者过去三年的编码辅助工具演化路径：

```
ChatGPT 复制粘贴 → Copilot 自动补全（对他没用）→ Cursor → Claude Code / Codex / Amp 等新一代编码代理
```

他偏爱 **Claude Code**，但随着时间推移，Claude Code 变成了"太空飞船"——80% 的功能他用不上，系统提示和工具每次发布都在变，破坏了他的工作流，而且界面闪烁。

### 1.2 核心痛点（为什么自己造轮子）

| 痛点                                        | 说明                                                    |
| ------------------------------------------- | ------------------------------------------------------- |
| **上下文工程（Context Engineering）不可控** | 现有工具会在背后注入你看不到的内容，UI 也不展示这些内容 |
| **可观测性（Observability）差**             | 几乎没有工具让你检查与模型交互的每一个细节              |
| **会话格式不透明**                          | 缺乏文档化的、可后处理的会话格式                        |
| **自托管（Self-hosting）支持差**            | Vercel AI SDK 对自托管模型（特别是 tool calling）不友好 |
| **API 有机演化的包袱**                      | 现有工具为了向后兼容积累了大量历史包袱                  |

### 1.3 他的设计哲学

> **"If I don't need it, it won't be built."**（不需要的东西就不造。）

这是一种极端的 YAGNI（You Aren't Gonna Need It）原则。

---

## 二、架构总览

pi 由四个包组成：

```
┌─────────────────────────────────────────────────┐
│                 pi-coding-agent                  │  ← CLI，把一切串起来
│  (会话管理、自定义工具、主题、项目上下文文件)      │
├─────────────────────┬───────────────────────────┤
│      pi-tui         │     pi-agent-core         │
│  (终端 UI 框架)      │  (代理循环、工具执行)       │
├─────────────────────┴───────────────────────────┤
│                    pi-ai                         │  ← 统一 LLM API
│  (多提供商、流式、工具调用、思维链、上下文切换)     │
└─────────────────────────────────────────────────┘
```

---

## 三、pi-ai：统一 LLM API 层（深度解析）

### 3.1 四种核心 API

与所有 LLM 提供商通信，本质上只需要对接 **四种 API**：

| API                          | 提供商                                          | 说明                   |
| ---------------------------- | ----------------------------------------------- | ---------------------- |
| **OpenAI Completions API**   | OpenAI、xAI、Groq、Cerebras、Mistral、Chutes... | 最广泛但各家实现有差异 |
| **OpenAI Responses API**     | OpenAI（新版）                                  | OpenAI 的下一代 API    |
| **Anthropic Messages API**   | Anthropic                                       | Claude 系列            |
| **Google Generative AI API** | Google                                          | Gemini 系列            |

### 3.2 提供商兼容性的"噩梦"细节

即便它们都声称支持 Completions API，但不同提供商的行为差异很大：

```typescript
// 在 openai-completions.ts 中的各种兼容性处理：
// - Cerebras、xAI、Mistral、Chutes 不支持 store 字段
// - Mistral 和 Chutes 使用 max_tokens 而不是 max_completion_tokens
// - Cerebras、xAI、Mistral、Chutes 不支持 developer role 的系统提示
// - Grok 模型不接受 reasoning_effort
// - 不同提供商把推理内容放在不同字段 (reasoning_content vs reasoning)
```

**核心教训**：构建统一 LLM API 不是火箭科学，但处理各家的"特色"才是真正的工作量。

### 3.3 Token 和成本追踪的痛苦

- Anthropic 的方式最合理
- 有些提供商在 SSE 流开始时报告 token 数，有些只在结束时报告
- **如果请求被中断，准确的成本追踪几乎不可能**
- 你无法提供唯一 ID 来与账单 API 关联
- Google 至今不支持 tool call streaming（作者的吐槽："which is extremely Google"）

### 3.4 上下文切换（Context Handoff）

这是 pi-ai **从一开始就设计进去**的核心特性：你可以在一个会话中**跨提供商切换模型**。

```typescript
// 用 Claude 开始对话
const claude = getModel('anthropic', 'claude-sonnet-4-5')
context.messages.push({ role: 'user', content: 'What is 25 * 18?' })
const claudeResponse = await complete(claude, context, { thinkingEnabled: true })

// 切换到 GPT —— Claude 的思维链会被转换为 <thinking> 标签的文本
const gpt = getModel('openai', 'gpt-5.1-codex')
context.messages.push({ role: 'user', content: 'Is that correct?' })

// 再切换到 Gemini
const gemini = getModel('google', 'gemini-2.5-flash')

// 上下文可以序列化/反序列化为 JSON
const serialized = JSON.stringify(context)
```

**技术挑战**：

- 不同提供商的思维链（thinking traces）格式不同
- 一些提供商会在事件流中插入**签名 blob**，后续请求必须回放
- 切换同一提供商内的不同模型也需要处理这些转换

### 3.5 类型安全的模型注册表

作者从 OpenRouter 和 models.dev 解析数据，生成 `models.generated.ts`，包含 token 成本和能力信息。可以类型安全地引用模型：

```typescript
const model = getModel('anthropic', 'claude-sonnet-4-5') // 类型安全！
```

也支持自定义/自托管模型：

```typescript
const ollamaModel: Model<'openai-completions'> = {
  id: 'llama-3.1-8b',
  api: 'openai-completions',
  provider: 'ollama',
  baseUrl: 'http://localhost:11434/v1'
  // ...
}
```

### 3.6 中止（Abort）支持

许多统一 LLM API 完全忽略了中止请求的能力——作者认为这"完全不可接受"。

pi-ai 从一开始就设计了全链路中止支持：

- 使用标准的 `AbortController` / `AbortSignal`
- 中止后仍然返回**部分结果**（很多 API 做不到）
- 支持整个管道的中止，包括工具调用

### 3.7 结构化拆分工具结果（Structured Split Tool Results）

这是作者引以为豪的创新：**工具结果分为两部分**：

```
┌────────────────────┐    ┌────────────────────┐
│   给 LLM 的内容     │    │   给 UI 的内容      │
│   (text/JSON)       │    │   (结构化数据)       │
└────────────────────┘    └────────────────────┘
```

```typescript
const weatherTool: AgentTool = {
  execute: async (toolCallId, args) => {
    return {
      output: `Temperature: ${temp}°C`, // 给 LLM
      details: { temp } // 给 UI（结构化数据）
    }
  }
}
```

工具还可以返回图片附件，会自动转换为各提供商的原生格式。

### 3.8 部分 JSON 解析（Partial JSON Parsing）

在工具调用流式传输过程中，pi-ai 会**渐进式解析 JSON 参数**，这样 UI 可以在调用完成之前就显示部分结果（比如文件 diff 的流式渲染）。

### 3.9 代理循环（Agent Loop）

`pi-agent-core` 提供的代理循环：

```
用户消息 → LLM 响应 → 有工具调用？ ─是→ 执行工具 → 结果反馈给 LLM → 循环
                                      │
                                      否→ 返回最终响应
```

**特点**：

- **没有 max_steps 限制**——"我从来没找到这个功能的用例"
- 支持消息队列：在代理工作时排队新消息
- 两种消息注入模式：逐条或批量
- 通过事件流驱动 UI

### 3.10 为什么不用 Vercel AI SDK？

作者引用了 Armin（opencode 作者）的博客文章，大意是：直接基于提供商 SDK 构建给你**完全控制权**，更小的 API 表面积，以及你想要的 API 设计自由度。

---

## 四、pi-tui：终端 UI 框架（深度解析）

### 4.1 两种 TUI 模式的权衡

| 方式     | 全屏 TUI                               | 追加式 TUI                        |
| -------- | -------------------------------------- | --------------------------------- |
| **原理** | 接管终端视口，当作像素缓冲区           | 像普通 CLI 程序一样写入终端       |
| **代表** | Amp、opencode                          | Claude Code、Codex、Droid、**pi** |
| **优点** | 完全控制渲染                           | 保留滚动缓冲区、原生搜索和滚动    |
| **缺点** | 丢失滚动缓冲区、需要自己实现搜索和滚动 | 渲染能力有限                      |

作者选择了**追加式 TUI**，因为编码代理本质上是线性的聊天界面，完美契合终端原生行为。

### 4.2 保留模式（Retained Mode）UI

pi-tui 采用简单的保留模式：

```
Component {
  render(width) → string[]     // 返回渲染后的行数组（含 ANSI 转义码）
  handleInput(data)            // 处理键盘输入
}

Container {
  children: Component[]        // 垂直排列的子组件
  → 收集所有子组件的渲染行
}

TUI extends Container {
  → 管理整个终端界面
}
```

**缓存优化**：已完成流式传输的助手消息不需要每次重新解析 Markdown 和生成 ANSI 序列——直接返回缓存的行。

### 4.3 差分渲染（Differential Rendering）

渲染算法非常简单：

```
1. 首次渲染：直接输出所有行
2. 宽度变化：清屏并完全重新渲染（因为软换行变了）
3. 普通更新：找到第一个不同的行，从该行开始重新渲染到末尾
4. 特殊情况：如果第一个变化行在可视区域上方（用户滚动了），全部清屏重渲染
```

**防闪烁**：使用 Synchronized Output 转义序列（`CSI ?2026h` 和 `CSI ?2026l`），告诉终端**缓冲所有输出并原子性地显示**。Ghostty 和 iTerm2 支持良好，VS Code 内置终端支持较差。

**开销分析**：

- 存储整个滚动缓冲区的"后备缓冲"——在现代计算机上只有几百 KB
- 每次渲染都要比较所有行——V8 引擎处理毫无压力
- 换来的是**极简的编程模型**

---

## 五、pi-coding-agent：核心设计哲学（深度解析）

这是全文最精华的部分——一系列**反主流**的设计决策。

### 5.1 极简系统提示（Minimal System Prompt）

pi 的完整系统提示**不到 1000 tokens**（加上工具定义）：

```
You are an expert coding assistant. You help users with coding tasks by
reading files, executing commands, editing code, and writing new files.

Available tools:
- read: Read file contents
- bash: Execute bash commands
- edit: Make surgical edits to files
- write: Create or overwrite files

Guidelines:
- Use bash for file operations like ls, grep, find
- Use read to examine files before editing
- Use edit for precise changes (old text must match exactly)
- Use write only for new files or complete rewrites
- Be concise in your responses
- Show file paths clearly when working with files
```

对比：

- **Claude Code** 的系统提示：巨大（数千 tokens）
- **opencode**：从 Claude Code 复制并裁剪
- **Codex**：类似 pi 的极简风格

**为什么这行得通？** 现代前沿模型已经被 RL 训练到天际——它们**天生理解**什么是编码代理。不需要一万个 token 的系统提示来教它们。

### 5.2 极简工具集（Minimal Toolset）

只有 **4 个工具**：

| 工具    | 用途                              |
| ------- | --------------------------------- |
| `read`  | 读取文件内容（支持文本和图片）    |
| `write` | 写文件（创建或覆盖）              |
| `edit`  | 精确文本替换（oldText → newText） |
| `bash`  | 执行 bash 命令                    |

对比 Claude Code 的几十个工具定义。作者认为**这四个工具就是编码代理需要的全部**——模型知道怎么用 bash，也被训练过类似的 read/write/edit 模式。

### 5.3 默认 YOLO 模式（无权限检查）

pi **没有任何安全围栏**：

- 不弹出文件操作权限对话框
- 不用 Haiku 预检 bash 命令
- 完全的文件系统访问
- 可以用你的用户权限执行任何命令

**作者的论点**：

> 其他编码代理的安全措施大多是"安全剧场"（security theater）。一旦你的代理可以写代码和运行代码，游戏就结束了。防止数据泄露的唯一方法是切断执行环境的所有网络访问，但这让代理几乎无法使用。

引用 Simon Willison 的"双 LLM"模式——连 Simon 自己都承认"这个解决方案很糟糕"。

**核心三元悖论**（不可能三角）：

```
      读取数据
       / \
      /   \
     /     \
执行代码 ── 网络访问

这三个能力同时存在时，安全就是不可能的。
既然无法解决，不如放弃假装安全。
```

### 5.4 不内置 To-Do（No Built-in To-Dos）

作者的观点：**内置的 to-do 列表通常让模型更困惑**，增加了模型需要追踪和更新的状态。

替代方案：让模型读写一个 `TODO.md` 文件——外部状态化、可见、可控。

```markdown
# TODO.md

- [x] Implement user authentication
- [ ] Write API documentation
```

### 5.5 不内置计划模式（No Plan Mode）

告诉代理"先想想这个问题，不要修改文件"就够了。如果需要持久化的计划，写到 `PLAN.md` 文件中：

- 跨会话持久
- 可以版本控制
- 可以和代理协同编辑
- **最重要的：完全可观测**

作者对 Claude Code Plan Mode 的批评：它最终也是写一个 markdown 文件到磁盘，但过程中你必须批准大量命令调用，而且 Claude Code 会催生子代理，你完全看不到子代理做了什么。

### 5.6 不支持 MCP（No MCP Support）

这是最"反潮流"的决定之一。作者的理由：

**成本问题**：

- Playwright MCP：21 个工具，13.7k tokens
- Chrome DevTools MCP：26 个工具，18k tokens
- 这些在每个会话开始就注入上下文，占用 7-9% 的上下文窗口

**替代方案**：构建 CLI 工具 + README 文件

```
MCP 方式：                           CLI 方式：
每次会话都加载所有工具描述              按需读取 README（渐进披露）
即使大多数不会用到                     只在需要时付出 token 成本
                                    可组合（管道、链式命令）
                                    易扩展（只需添加脚本）
```

如果非要用 MCP，可以用 `mcporter` 工具把 MCP 服务器包装成 CLI 工具。

### 5.7 不支持后台 Bash（No Background Bash）

pi 的 bash 工具同步执行。不存在后台进程管理。

**替代方案**：使用 **tmux**！

- 可以在 tmux 中运行开发服务器、调试器
- 完全的可观测性——你可以跳进去和代理一起调试
- tmux 提供 CLI 参数来列出所有活跃会话
- Claude Code 也能用 tmux，不需要内置后台 bash

### 5.8 不支持子代理（No Sub-Agents）

Claude Code 的子代理问题：

- **零可观测性**——"黑盒中的黑盒"
- 上下文传递效果差
- 调试子代理的错误很痛苦

**作者的根本观点**：

> 在会话中使用子代理做上下文收集，说明你没有提前规划。应该先在独立会话中收集上下文，创建一个有用的制品（artifact），然后在新会话中使用它。

```
❌ 反模式：会话中途催生子代理收集上下文
✅ 正确做法：独立会话收集上下文 → 创建制品 → 新会话使用制品
```

如果确实需要子代理，直接让 pi 通过 bash 生成自己的新实例：

```bash
pi --print "review the code for bugs"
```

或者在 tmux 中催生，获得完全可观测性。

**对并行子代理的看法**：催生多个子代理并行实现不同功能是"反模式"，会让代码库变成垃圾堆。

---

## 六、Benchmark 结果

作者在 **Terminal-Bench 2.0** 上运行了测试（每个任务 5 次试验），pi + Claude Opus 4.5 的表现：

- 与 Codex、Cursor、Windsurf 等竞争
- 在排行榜上表现优秀
- 注意到**美国西海岸上线后（PST 时间）错误率上升**，所以单独做了一次"只在 CET 时间"运行的测试

**有趣发现**：Terminal-Bench 团队自己的 **Terminus 2**（一个极简代理，只给模型一个 tmux 会话，没有任何花哨工具）也在排行榜上表现不错——进一步证明极简方法可以做得同样好。

---

## 七、关键思想总结

### 7.1 上下文工程是核心

> **"Context engineering is paramount."**

精确控制什么进入模型的上下文窗口，比任何花哨功能都重要。现有工具在背后注入太多内容，而且不透明。

### 7.2 可观测性 > 功能丰富

贯穿全文的主题：作者宁要**可观测的简单系统**，也不要封装良好但不透明的复杂系统。

### 7.3 模型已经足够聪明

前沿模型不需要冗长的系统提示和复杂的工具定义。1000 tokens 以下的提示 + 4 个工具就够了。

### 7.4 "用文件代替功能"模式

```
内置 To-Do  → TODO.md 文件
计划模式     → PLAN.md 文件
MCP 服务器   → CLI 工具 + README 文件
后台进程     → tmux 会话
子代理       → bash 催生新 pi 实例
```

这个模式的哲学是：**不要在代理中内置状态管理，让文件系统成为状态的来源**。

### 7.5 安全的诚实观

与其搞无效的安全剧场，不如承认：只要代理能读数据+执行代码+访问网络，就没有真正的安全。要安全就用容器隔离。

---

## 八、对我们的启示

1. **构建自己的工具是最好的学习方式**——通过造轮子理解每一层抽象的代价
2. **极简不是偷懒，是工程纪律**——每个"不加"的功能都是经过深思熟虑的
3. **可观测性是所有复杂系统的基石**——看不到的东西无法调试
4. **文件系统是最好的持久化 API**——简单、通用、可版本控制、跨会话
5. **不要迷信 benchmark**——作者自己说"real proof is in the pudding"（实践出真知）
6. **统一 API 的价值和代价**——标准化节省了时间，但"leaky abstractions"（抽象泄漏）是不可避免的
