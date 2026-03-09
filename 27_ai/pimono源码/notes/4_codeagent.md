# Pi Coding Agent 深度源码分析

**Pi Coding Agent 的十大关键洞见：**

1. **三层架构分离** — agent-core（纯循环引擎）→ coding-agent/core（领域粘合）→ modes（I/O 展现），使 Agent 核心可嵌入任何运行时

2. **Steer/FollowUp 双轨队列** — 用户可在 Agent 执行中途"插嘴"（steer 打断工具链）或"追加任务"（followUp 等当前轮完成），解决了"必须等完才能说话"的问题

3. **重试 Promise 提前创建** — 在同步事件回调中创建 Promise，异步队列中 resolve，精确避免 `waitForRetry()` 与事件处理的竞态条件

4. **Edit 的 Fuzzy Match** — 先精确匹配，失败后归一化 Unicode（智能引号→ASCII、特殊空格/破折号→普通字符）再匹配，容忍 LLM 输出的微小偏差

5. **追加制 JSONL + id/parentId 树** — 永不删改已有条目，分支只需移动 leafId 指针，崩溃安全且支持无限分支

6. **增量压缩摘要** — 后续压缩不是重写摘要，而是在前次摘要基础上增量更新，防止信息衰减

7. **溢出恢复的"一次机会"** — context overflow 时压缩→重试，但只给一次机会，避免死循环

8. **工具双层包装** — 扩展注册的工具先获得事件钩子，然后所有工具（含内建）再被 `tool_call/tool_result` 拦截层包装，实现"任何工具都可审计/替换"

9. **工具感知的 System Prompt** — 可用工具集变化时，指导原则自动调整（如"有 grep 时别用 bash grep"）

10. **模式无关的 AgentSession** — Interactive/Print/RPC 共享同一个 AgentSession，模式只管 I/O，核心逻辑零重复

> 基于 `@mariozechner/pi-coding-agent` v0.55.4 的完整架构剖析

---

## 一、全局架构鸟瞰

```
┌─────────────────────────────────────────────────────────┐
│                     CLI Entry (main.ts)                 │
│  args → config → resources → session → model → mode    │
└─────────────────────┬───────────────────────────────────┘
                      │
          ┌───────────┼───────────────┐
          ▼           ▼               ▼
   Interactive    Print Mode     RPC Mode
   (TUI+Editor)  (stdin→stdout)  (JSON-RPC)
          │           │               │
          └───────────┼───────────────┘
                      ▼
            ┌─────────────────┐
            │  AgentSession   │  ← 核心协调层,所有模式共享
            │  (会话+事件+工具)│
            └────────┬────────┘
                     │
          ┌──────────┼──────────┐
          ▼          ▼          ▼
       Agent     SessionMgr  Extensions
    (pi-agent-core)  (JSONL)   (loader+runner)
       │
       ▼
    Agent Loop
    (prompt → stream → tools → steer/followUp)
```

### 关键洞见 #1：三层架构分离

Pi 的设计核心理念是**关注点分离**：

1. **pi-agent-core**：纯粹的 Agent 循环引擎，不知道文件系统、UI 或会话持久化的存在
2. **coding-agent/core**：领域层，将 Agent 循环与文件工具、会话存储、上下文压缩粘合
3. **coding-agent/modes**：展现层，Interactive/Print/RPC 三种 I/O 策略

这种分层使得 Agent 核心可以被嵌入任何运行时（SDK、Web、Electron），而不被终端 UI 所绑架。

---

## 二、Agent 循环引擎 (pi-agent-core)

### 2.1 三层循环模型

```
┌─ agentLoop ──────────────────────────────────────────┐
│                                                       │
│  ┌─ Outer Loop (Follow-up Queue) ──────────────────┐ │
│  │  while (hasFollowUp) {                          │ │
│  │    deliverFollowUp()                            │ │
│  │                                                  │ │
│  │    ┌─ Inner Loop (Tool + Steering) ───────────┐ │ │
│  │    │  while (true) {                          │ │ │
│  │    │    transformContext()    // 清洗消息      │ │ │
│  │    │    convertToLlm()       // AgentMsg→LLM  │ │ │
│  │    │    streamResponse()     // LLM 推理      │ │ │
│  │    │                                          │ │ │
│  │    │    if (hasToolCalls) {                    │ │ │
│  │    │      for each toolCall {                 │ │ │
│  │    │        executeTool()                     │ │ │
│  │    │        if (hasSteering) break → inject   │ │ │
│  │    │      }                                   │ │ │
│  │    │      continue // 回到 LLM               │ │ │
│  │    │    } else {                              │ │ │
│  │    │      break // 无工具调用,结束内循环       │ │ │
│  │    │    }                                      │ │ │
│  │    │  }                                        │ │ │
│  │    └──────────────────────────────────────────┘ │ │
│  │                                                  │ │
│  │  } // 检查下一个 followUp                       │ │
│  └──────────────────────────────────────────────────┘ │
│                                                       │
└───────────────────────────────────────────────────────┘
```

### 关键洞见 #2：Steer 与 FollowUp 的队列双轨制

这是 Pi 并发交互模型的精髓：

- **Steer（转向）**：在 Agent 执行**工具调用之间**检查。效果是"打断当前工作流"——跳过剩余工具调用，将用户的新消息注入，让 LLM 重新规划。
- **FollowUp（续接）**：在 Agent 完成一轮**全部工具调用后、即将结束前**检查。效果是"追加新任务"——Agent 完成当前任务后自动接上。

两种队列都支持 `"all"` 和 `"one-at-a-time"` 两种投递模式。这意味着用户可以批量排队多条指令或逐条投递。

**为什么重要？** 这解决了传统 Chat UI 中"必须等 Agent 说完才能发消息"的尴尬。用户可以在 Agent 执行中途随时"插嘴"或"追加任务"。

---

## 三、AgentSession：核心协调层

`AgentSession`（约 2500 行）是整个系统的**中枢神经**。它负责：

### 3.1 事件驱动的会话持久化

```typescript
// Agent 发出事件 → AgentSession 监听 → 写入 JSONL
agent.subscribe(event => {
  if (event.type === 'message_end') {
    sessionManager.appendMessage(event.message) // 持久化到 JSONL
  }
})
```

核心设计：**事件流是真理之源**。AgentSession 通过订阅 Agent 事件来驱动一切副作用——存盘、UI 刷新、扩展通知、自动压缩。

`_processAgentEvent` 方法是整个事件处理的中枢：

```
event → emitExtensionEvent() → emit() to listeners → persistence
      → if agent_end: checkRetry() → checkCompaction()
```

### 3.2 事件队列的顺序保证

```typescript
this._agentEventQueue = this._agentEventQueue.then(
  () => this._processAgentEvent(event),
  () => this._processAgentEvent(event)
)
```

所有事件处理被串行化进 Promise 链。这保证了即使 Agent 快速连续触发事件，处理顺序也与触发顺序一致。两个 then 参数（fulfilled + rejected）确保链条不会因为某个处理器抛错而断裂。

### 关键洞见 #3：重试的 Promise 提前创建

```typescript
private _createRetryPromiseForAgentEnd(event: AgentEvent): void {
    // 在同步事件处理中就创建 Promise，而不是在异步队列中
}
```

Agent 的 `prompt()` 在 `agent.prompt()` resolve 后立即调用 `waitForRetry()`。如果 retry Promise 是在异步的事件队列中才创建的，那么当事件队列还没处理到 `agent_end` 时，`waitForRetry()` 就会错过正在进行的重试。

**解法**：在同步的事件回调中就立即创建 Promise，然后异步的事件处理去 resolve 它。这是一个**竞态条件的经典防御**。

---

## 四、工具系统

### 4.1 七把利刃

| 工具    | 功能                      | 默认启用 |
| ------- | ------------------------- | -------- |
| `read`  | 读取文件内容              | ✅       |
| `bash`  | 执行 shell 命令           | ✅       |
| `edit`  | 精确文本替换              | ✅       |
| `write` | 创建/覆盖文件             | ✅       |
| `grep`  | 内容搜索(尊重 .gitignore) | ❌       |
| `find`  | 文件查找(glob)            | ❌       |
| `ls`    | 列出目录                  | ❌       |

### 4.2 工具注册表的三层架构

```
baseToolRegistry     → 内建 7 个工具
 + extensionTools    → 扩展注册的自定义工具
 = toolRegistry      → 全量工具表 (extensionRunner wrap 过)
   → setActiveToolsByName() → agent.setTools()  → 实际生效的工具子集
```

### 关键洞见 #4：Edit 工具的 Fuzzy Match 容错

LLM 生成的 `oldText` 经常有微小偏差（Unicode 引号、尾部空格、特殊破折号）。Pi 的 Edit 工具实现了一个优雅的容错机制：

```typescript
export function fuzzyFindText(content: string, oldText: string): FuzzyMatchResult {
  // 1. 先尝试精确匹配
  const exactIndex = content.indexOf(oldText)
  if (exactIndex !== -1) return exactMatch

  // 2. 失败后进行模糊匹配：
  //    - 去除行尾空格
  //    - 智能引号 → ASCII 引号
  //    - Unicode 破折号 → ASCII 连字符
  //    - 特殊空格 → 普通空格
  const fuzzyContent = normalizeForFuzzyMatch(content)
  const fuzzyOldText = normalizeForFuzzyMatch(oldText)
  const fuzzyIndex = fuzzyContent.indexOf(fuzzyOldText)
  // ...
}
```

这意味着当 LLM 把 `'` 输出成 `'`（智能引号）时，编辑仍然能成功。**收窄模糊范围**（只处理已知的常见偏差）比正则全匹配更安全。

### 4.3 Bash 工具的流式截断

```
命令 → spawn(shell) → stdout/stderr 合流 → 写临时文件(大输出)
                                          → 保留滚动缓冲区
                                          → 截断为末尾 N 行 / N KB
```

- 默认截断到最后 256 行或 32KB
- 完整输出保存到临时文件，路径返回给 LLM
- 这防止了 `cat` 一个大文件把上下文窗口撑爆

---

## 五、会话存储：追加制 JSONL 树

### 5.1 数据模型

```
session.jsonl:
────────────────────────────
{type:"session", version:3, id:"abc", cwd:"/project"}  ← 头部
{type:"message", id:"1", parentId:null, message:{role:"user",...}}
{type:"message", id:"2", parentId:"1", message:{role:"assistant",...}}
{type:"message", id:"3", parentId:"2", message:{role:"user",...}}
{type:"compaction", id:"4", parentId:"3", summary:"...", firstKeptEntryId:"3"}
{type:"message", id:"5", parentId:"4", message:{role:"user",...}}
...
```

### 关键洞见 #5：追加制 + 树结构 = 零数据丢失的分支

传统做法是「改写文件」或「SQL 数据库」。Pi 选择了**追加制 JSONL + id/parentId 链**：

- **永不删除、永不修改**已有条目。创建分支就是把 `leafId` 指向某个旧条目，然后追加新条目。
- **读取路径**：从 leafId 沿 parentId 链回溯到根，反转即为当前分支的顺序。
- **优势**：崩溃安全（只要追加成功就不丢数据）、可审计、支持无限分支。
- **代价**：文件会持续增长（包含所有分支的全部数据）。

### 5.2 buildSessionContext()

这个方法将 JSONL 树**重建为 LLM 可用的消息列表**：

```
1. 构建 id→entry 索引
2. 从 leafId 沿 parentId 链回溯到根，收集路径
3. 找到路径上最新的 compaction 节点
4. 如果有 compaction：
   - 先输出 compaction summary（作为特殊消息）
   - 输出 compaction 节点标记的 firstKeptEntryId 之后的消息
5. 如果没有 compaction：输出全部消息
6. 跳过 label、custom（非消息类）、model_change 等元数据条目
   （它们只影响 model/thinkingLevel 的恢复，不进入 LLM 上下文）
```

---

## 六、上下文工程（Context Engineering）深度剖析

6.1 上下文构建 — buildSessionContext() 如何从树形会话中提取线性消息序列，以及 convertToLlm() 的类型转换
6.2 Token 估算 — 双层策略（Usage 精确值 + chars/4 启发式），以及为什么不用 tiktoken
6.3 压缩触发 — 溢出触发 vs 阈值触发的完整决策流程，包括模型切换和旧错误的防御逻辑
6.4 切割点算法 — findCutPoint() 三步走，正常切割 vs Split Turn 的图解
6.5 prepareCompaction() — 纯函数的压缩准备阶段
6.6 增量摘要 — 两套提示词的设计意图和信息衰减防护
6.7 compact() 执行 — Split Turn 的并行双摘要合并
6.8 \_runAutoCompaction() 编排 — 从事件发出到扩展钩子到消息重建的完整流程
6.9 文件操作累积追踪 — 跨多次压缩的 read/modified 文件累积机制
6.10 分支摘要 — 公共祖先查找、token 预算管理、文件操作不受预算限制的设计
6.11 溢出恢复的一次机会 — 防死循环设计
6.12 扩展自定义 — session_before_compact 和 session_compact 钩子
6.13 配置参数 — 三个参数的含义和调优建议
6.14 端到端流程图 — 从用户消息到压缩完成的完整链路
6.15 与其他 Agent 对比 — Pi vs Claude Code vs Cursor 的压缩策略差异

> **核心源码文件：**
>
> - `src/core/compaction/compaction.ts` — 自动压缩逻辑
> - `src/core/compaction/branch-summarization.ts` — 分支摘要
> - `src/core/compaction/utils.ts` — 共享工具（文件追踪、序列化）
> - `src/core/session-manager.ts` — 上下文构建（`buildSessionContext`）
> - `src/core/messages.ts` — 消息类型转换与 LLM 适配
> - `src/core/agent-session.ts` — 压缩触发与编排
> - `src/core/extensions/types.ts` — 扩展钩子类型定义

Pi 的"上下文工程"不只是压缩。它是一整套**管理 LLM 有限上下文窗口的机制**，涵盖：

1. **上下文构建** — 从树形会话中提取 LLM 可见的消息序列
2. **自动压缩 (Compaction)** — 当上下文超限时摘要化旧内容
3. **分支摘要 (Branch Summarization)** — 在树形会话中切换分支时保留上下文
4. **溢出恢复** — LLM 报上下文溢出错误时自动压缩重试
5. **扩展自定义** — 允许扩展接管或定制压缩行为

### 6.1 上下文构建：`buildSessionContext()`

这是整个上下文工程的**起点**——每次调用 LLM 前，都要从会话树中构建出一个线性消息序列。

```typescript
// session-manager.ts
function buildSessionContext(entries, leafId, byId): SessionContext {
  // 1. 从 leaf 沿 parentId 链走到根，收集路径
  const path: SessionEntry[] = []
  let current = byId.get(leafId)
  while (current) {
    path.unshift(current)
    current = current.parentId ? byId.get(current.parentId) : undefined
  }

  // 2. 沿路径提取设置变更（thinkingLevel, model）
  //    和最新 compaction 节点
  let compaction: CompactionEntry | null = null
  for (const entry of path) {
    if (entry.type === "compaction") compaction = entry
    if (entry.type === "thinking_level_change") thinkingLevel = entry.thinkingLevel
    if (entry.type === "model_change") model = { provider, modelId }
  }

  // 3. 构建消息列表
  if (compaction) {
    // 先输出压缩摘要（作为 user 角色的特殊消息）
    messages.push(createCompactionSummaryMessage(compaction.summary, ...))
    // 再输出 firstKeptEntryId 之后到 compaction 之间的消息
    // 最后输出 compaction 之后的消息
  } else {
    // 无压缩：输出路径上所有消息
  }

  return { messages, thinkingLevel, model }
}
```

**关键设计：只用最新的 compaction 节点**。树路径上可能有多个 compaction 条目（多次压缩），但只有最新的那一个生效——它的摘要已经累积了所有之前的信息。

消息发送给 LLM 前，还需要经过 `convertToLlm()` 转换：

```typescript
// messages.ts - 将自定义消息类型转换为标准 LLM 消息
function convertToLlm(messages: AgentMessage[]): Message[] {
  return messages.map(m => {
    switch (m.role) {
      case 'bashExecution':
        return { role: 'user', content: bashExecutionToText(m) }
      case 'compactionSummary':
        return {
          role: 'user',
          content: COMPACTION_SUMMARY_PREFIX + m.summary + COMPACTION_SUMMARY_SUFFIX
        }
      // 包裹在 <summary>...</summary> 标签中
      case 'branchSummary':
        return { role: 'user', content: BRANCH_SUMMARY_PREFIX + m.summary + BRANCH_SUMMARY_SUFFIX }
      case 'custom':
        return { role: 'user', content: m.content }
      default:
        return m // user, assistant, toolResult 原样传递
    }
  })
}
```

LLM 看到的最终消息序列结构：

```
┌────────┬─────────────────────┬─────────────────────────┐
│ system │ [summary as user]   │ kept messages...        │
│ prompt │ (from compaction)   │ (from firstKeptEntryId) │
└────────┴─────────────────────┴─────────────────────────┘
```

注意：**压缩摘要被注入为 user 角色消息**，而不是 system prompt 的一部分。这是因为摘要内容是动态的、会话特定的，不适合放在系统提示中。

### 6.2 Token 估算机制

Pi 使用两层 token 估算策略：

**第一层：基于 LLM 返回的 Usage 数据（精确）**

```typescript
// 取最后一个非中止/非错误的 assistant 消息的 usage
function calculateContextTokens(usage: Usage): number {
  return usage.totalTokens || usage.input + usage.output + usage.cacheRead + usage.cacheWrite
}
```

**第二层：字符估算（保守，用于无 usage 数据时）**

```typescript
function estimateTokens(message: AgentMessage): number {
  let chars = 0
  // 累加所有文本内容的字符数
  // 图片估算为 4800 字符 ≈ 1200 tokens
  return Math.ceil(chars / 4) // chars/4 启发式
}
```

**混合策略 `estimateContextTokens()`** — 如果有最近的 assistant usage 就用它作为基准，再加上后续消息的估算值：

```typescript
function estimateContextTokens(messages: AgentMessage[]): ContextUsageEstimate {
  const usageInfo = getLastAssistantUsageInfo(messages)
  if (!usageInfo) {
    // 没有任何 usage 数据，全部用估算
    return { tokens: sumEstimates(messages), ... }
  }
  // 有 usage：精确值 + 后续消息的估算
  const usageTokens = calculateContextTokens(usageInfo.usage)
  let trailingTokens = 0
  for (let i = usageInfo.index + 1; i < messages.length; i++) {
    trailingTokens += estimateTokens(messages[i])
  }
  return { tokens: usageTokens + trailingTokens, ... }
}
```

> **为什么用 chars/4 而不是 tiktoken？** 因为 Pi 支持多种模型提供商（Anthropic、OpenAI、Google、Ollama 等），每家的 tokenizer 不同。chars/4 是一个**保守的过估（overestimate）**，确保不会因为低估而导致意外溢出。

### 6.3 压缩触发：双通道机制

压缩的触发在 `agent-session.ts` 的 `_checkCompaction()` 中，它检查**两种触发条件**：

```
┌──────────────────────────────────────────────────────┐
│              _checkCompaction() 决策流程              │
│                                                      │
│  agent_end 事件 → _lastAssistantMessage              │
│       │                                              │
│       ├─ 先检查是否需要重试（retryable error）        │
│       │    是 → 重试，不走压缩                        │
│       │                                              │
│       ├─ Case 1: 溢出触发                            │
│       │    条件: isContextOverflow(msg) &&            │
│       │          sameModel &&                        │
│       │          !errorIsFromBeforeCompaction         │
│       │    动作: 删除错误消息 → 压缩 → 自动重试      │
│       │    限制: _overflowRecoveryAttempted 一次机会  │
│       │                                              │
│       └─ Case 2: 阈值触发                            │
│            条件: contextTokens > contextWindow        │
│                  - reserveTokens                     │
│            动作: 压缩（不自动重试）                    │
└──────────────────────────────────────────────────────┘
```

**两种触发的关键差异：**

|              | 溢出触发 (overflow)                           | 阈值触发 (threshold)      |
| ------------ | --------------------------------------------- | ------------------------- |
| 触发时机     | LLM 返回 context overflow 错误                | 正常回复后检查 token 用量 |
| 是否重试     | **是**，压缩后 `agent.continue()`             | **否**，用户继续手动操作  |
| 防御措施     | 只允许一次，`_overflowRecoveryAttempted`      | 无（每次超限都会触发）    |
| 错误消息处理 | 从 agent state 中**移除**（但保留在 session） | 无需处理                  |

溢出触发还有两个精巧的防御逻辑：

```typescript
// 防御 1: 模型切换后不误触发
// 从 opus (小窗口) 切换到 codex (大窗口) 后，opus 的溢出错误不应触发压缩
const sameModel =
  assistantMessage.provider === this.model.provider && assistantMessage.model === this.model.id

// 防御 2: 压缩后残留的旧错误不应再触发压缩
// 场景: opus 失败 → 切到 codex → 压缩 → 切回 opus → opus 旧错误仍在 kept 区域
const errorIsFromBeforeCompaction =
  compactionEntry !== null &&
  assistantMessage.timestamp < new Date(compactionEntry.timestamp).getTime()
```

### 6.4 切割点算法详解

切割点决定了哪些消息被摘要、哪些被保留。这是压缩的核心算法。

**`findCutPoint()` 三步走：**

```
Step 1: 找出所有合法切割点
  ← 遍历 [startIndex, endIndex) 范围内的 entry
  ← 合法类型: user, assistant, bashExecution, custom, branchSummary, customMessage
  ← 不合法: toolResult（必须跟工具调用在一起）

Step 2: 从最新消息向前累计 token
  ← 每个 message entry 估算 token
  ← 累计 ≥ keepRecentTokens 时，找到当前位置最近的合法切割点

Step 3: 判断是否Split Turn
  ← 切割点是 user 消息？ → 正常切割（isSplitTurn = false）
  ← 切割点是 assistant 消息？ → 这轮太大了，需要拆分（isSplitTurn = true）
  ← 向前找到这个 turn 的起始 user 消息 → turnStartIndex
```

**图解正常切割 vs Split Turn：**

```
正常切割（turn 边界）:
  ┌─────┬─────┬──────┬─────┬─────┬──────┬─────┐
  │ usr │ ass │ tool │ usr │ ass │ tool │ ass │
  └─────┴─────┴──────┴─────┴─────┴──────┴─────┘
  └──── summarize ───┘ └──── keep ──────────┘
                       ↑ cutPoint (user msg)

Split Turn（单轮超大）:
  ┌─────┬─────┬──────┬─────┬──────┬─────┬──────┐
  │ usr │ ass │ tool │ ass │ tool │ ass │ tool │
  └─────┴─────┴──────┴─────┴──────┴─────┴──────┘
  ↑ turnStart                     ↑ cutPoint (assistant msg)
  └───── turnPrefixMessages ──────┘ └── keep ──┘
```

**Split Turn 的处理策略：** 当一个 turn 大到超过 `keepRecentTokens` 时，Pi 不会把整个 turn 都扔掉。它会：

1. 把 turn 的前半段（从 user 消息到切割点之前）作为 `turnPrefixMessages`
2. turn 前面的完整 turns 作为 `messagesToSummarize`
3. **并行生成两个摘要**然后合并：

```typescript
const [historyResult, turnPrefixResult] = await Promise.all([
  messagesToSummarize.length > 0
    ? generateSummary(messagesToSummarize, model, ...)
    : "No prior history.",
  generateTurnPrefixSummary(turnPrefixMessages, model, ...),
])
summary = `${historyResult}\n\n---\n\n**Turn Context (split turn):**\n\n${turnPrefixResult}`
```

Turn prefix 有专用的精简提示词（`TURN_PREFIX_SUMMARIZATION_PROMPT`），只关注三件事：

- 原始请求是什么
- 前半段做了什么
- 后半段需要什么上下文

### 6.5 `prepareCompaction()` —— 压缩准备阶段

这个纯函数完成所有压缩前的计算，不触碰任何 I/O：

```typescript
function prepareCompaction(pathEntries, settings): CompactionPreparation | undefined {
  // 1. 如果最后一个 entry 就是 compaction，跳过（刚压缩过）
  if (pathEntries[pathEntries.length - 1].type === 'compaction') return undefined

  // 2. 找到上一次 compaction 的位置（前向搜索边界）
  let prevCompactionIndex = -1
  for (let i = pathEntries.length - 1; i >= 0; i--) {
    if (pathEntries[i].type === 'compaction') {
      prevCompactionIndex = i
      break
    }
  }
  const boundaryStart = prevCompactionIndex + 1 // 从上次压缩之后开始

  // 3. 估算当前 token 用量
  const tokensBefore = estimateContextTokens(usageMessages).tokens

  // 4. 找切割点
  const cutPoint = findCutPoint(pathEntries, boundaryStart, boundaryEnd, settings.keepRecentTokens)

  // 5. 分离 messagesToSummarize 和 turnPrefixMessages
  const historyEnd = cutPoint.isSplitTurn ? cutPoint.turnStartIndex : cutPoint.firstKeptEntryIndex

  // 6. 提取前次摘要（用于增量更新）
  let previousSummary = prevCompaction?.summary

  // 7. 提取文件操作（从前次 details + 当前消息）
  const fileOps = extractFileOperations(messagesToSummarize, pathEntries, prevCompactionIndex)

  return {
    firstKeptEntryId,
    messagesToSummarize,
    turnPrefixMessages,
    isSplitTurn,
    tokensBefore,
    previousSummary,
    fileOps,
    settings
  }
}
```

返回的 `CompactionPreparation` 是一个**不可变的快照**，扩展可以检查它来决定是否自定义压缩行为。

### 6.6 增量摘要 vs 全量摘要

Pi 使用**两套不同的提示词**来生成摘要：

**首次压缩 `SUMMARIZATION_PROMPT`：** 从零生成结构化摘要

```
The messages above are a conversation to summarize.
Create a structured context checkpoint summary...
Use this EXACT format:
## Goal / ## Constraints / ## Progress / ## Key Decisions / ## Next Steps / ## Critical Context
```

**后续压缩 `UPDATE_SUMMARIZATION_PROMPT`：** 在前次摘要基础上增量更新

```
The messages above are NEW conversation messages to incorporate
into the existing summary provided in <previous-summary> tags.

Update the existing structured summary with new information. RULES:
- PRESERVE all existing information from the previous summary
- ADD new progress, decisions, and context from the new messages
- UPDATE the Progress section: move items from "In Progress" to "Done"
- UPDATE "Next Steps" based on what was accomplished
```

增量策略的好处：**防止信息衰减**。每次压缩都是基于前次摘要 + 新消息，关键决策和上下文不断累积，而不是从头重写导致旧信息丢失。

**消息序列化**在发给 LLM 摘要前，对话需要被"扁平化"为文本（防止模型继续对话）：

```typescript
// utils.ts - serializeConversation()
function serializeConversation(messages: Message[]): string {
  // [User]: message text
  // [Assistant thinking]: thinking content
  // [Assistant]: response text
  // [Assistant tool calls]: read(path="foo.ts"); edit(path="bar.ts", ...)
  // [Tool result]: output text
}
```

**摘要系统提示词**也有专门设计（`SUMMARIZATION_SYSTEM_PROMPT`）：

```
You are a context summarization assistant. Your task is to read a conversation
between a user and an AI coding assistant, then produce a structured summary
following the exact format specified.

Do NOT continue the conversation. Do NOT respond to any questions in the
conversation. ONLY output the structured summary.
```

> **为什么要强调"不要继续对话"？** 因为摘要模型看到的是一段完整的对话文本，很容易被误导去回答对话中的问题。这条指令确保模型只做摘要。

### 6.7 `compact()` —— 压缩执行

```typescript
async function compact(preparation, model, apiKey, customInstructions?, signal?): Promise<CompactionResult> {
  if (isSplitTurn && turnPrefixMessages.length > 0) {
    // 并行生成两个摘要（历史 + turn 前缀）
    const [historyResult, turnPrefixResult] = await Promise.all([...])
    summary = `${historyResult}\n\n---\n\n**Turn Context (split turn):**\n\n${turnPrefixResult}`
  } else {
    // 只生成历史摘要
    summary = await generateSummary(messagesToSummarize, model, ...)
  }

  // 在摘要末尾追加文件操作列表
  const { readFiles, modifiedFiles } = computeFileLists(fileOps)
  summary += formatFileOperations(readFiles, modifiedFiles)
  // 生成格式:
  // <read-files>
  // path/to/file1.ts
  // </read-files>
  // <modified-files>
  // path/to/changed.ts
  // </modified-files>

  return { summary, firstKeptEntryId, tokensBefore, details: { readFiles, modifiedFiles } }
}
```

### 6.8 `_runAutoCompaction()` —— 完整编排流程

`agent-session.ts` 中的编排逻辑是理解压缩全貌的关键：

```
_runAutoCompaction(reason, willRetry)
│
├─ 1. 发出 auto_compaction_start 事件（通知 UI 显示"正在压缩..."）
├─ 2. 创建 AbortController（用户可取消）
├─ 3. 获取 API key 和 branch entries
├─ 4. prepareCompaction(pathEntries, settings)  → CompactionPreparation
│
├─ 5. 扩展钩子: session_before_compact
│     ├─ 扩展返回 { cancel: true } → 中止压缩
│     └─ 扩展返回 { compaction: {...} } → 使用扩展提供的摘要
│
├─ 6. 如果扩展没接管 → compact(preparation, model, apiKey)
│
├─ 7. 保存到 session: sessionManager.appendCompaction(summary, firstKeptEntryId, ...)
├─ 8. 重建消息: sessionManager.buildSessionContext() → agent.replaceMessages()
│
├─ 9. 扩展钩子: session_compact（通知压缩完成）
├─ 10. 发出 auto_compaction_end 事件
│
└─ 11. 后续动作:
       ├─ willRetry (溢出恢复) → 删除残留错误消息 → agent.continue()
       └─ hasQueuedMessages → agent.continue()（处理等待中的消息）
```

**重要细节：** 步骤 11 的 `setTimeout(() => agent.continue(), 100)` 中的 100ms 延迟，是为了让事件循环有时间处理压缩完成后的 UI 更新。

### 6.9 文件操作的累积追踪

文件追踪是压缩中容易被忽略但极其重要的机制。它确保 LLM 在多次压缩后仍然知道"项目中哪些文件被动过"。

**追踪来源：** 从 assistant 消息中的 `toolCall` 提取

```typescript
function extractFileOpsFromMessage(message, fileOps) {
  // 只看 assistant 消息中的 toolCall block
  for (const block of message.content) {
    if (block.type !== 'toolCall') continue
    const path = block.arguments.path
    switch (block.name) {
      case 'read':
        fileOps.read.add(path)
        break
      case 'write':
        fileOps.written.add(path)
        break
      case 'edit':
        fileOps.edited.add(path)
        break
    }
  }
}
```

**累积机制：** 每次压缩时，先从前次 compaction 的 `details` 中恢复旧的文件列表，再从新消息中提取追加：

```typescript
function extractFileOperations(messages, entries, prevCompactionIndex) {
  const fileOps = createFileOps()

  // 从前次 compaction details 累积
  if (prevCompactionIndex >= 0) {
    const prev = entries[prevCompactionIndex] as CompactionEntry
    if (!prev.fromHook && prev.details) {
      for (const f of prev.details.readFiles) fileOps.read.add(f)
      for (const f of prev.details.modifiedFiles) fileOps.edited.add(f)
    }
  }

  // 从新消息中提取
  for (const msg of messages) extractFileOpsFromMessage(msg, fileOps)

  return fileOps
}
```

**最终计算：** `computeFileLists()` 将 read/written/edited Set 合并为两个列表：

- `readFiles` = 只读过没改过的文件（`read - (edited ∪ written)`）
- `modifiedFiles` = 被修改过的文件（`edited ∪ written`）

### 6.10 分支摘要（Branch Summarization）

当使用 `/tree` 导航到不同分支时，Pi 会提供离开分支的摘要，注入到新分支的上下文中。

**导航示例：**

```
         ┌─ B ─ C ─ D (当前位置，即将离开)
    A ───┤
         └─ E ─ F (目标)

公共祖先: A
需要摘要: B, C, D

导航后:
         ┌─ B ─ C ─ D ─ [B,C,D 的摘要]
    A ───┤
         └─ E ─ F (新的当前位置)
```

**核心算法 `collectEntriesForBranchSummary()`：**

```typescript
function collectEntriesForBranchSummary(session, oldLeafId, targetId) {
  // 1. 取 oldLeaf 的完整路径 Set
  const oldPath = new Set(session.getBranch(oldLeafId).map(e => e.id))

  // 2. 取 target 的路径，从后往前找第一个也在 oldPath 中的 = 公共祖先
  const targetPath = session.getBranch(targetId)
  let commonAncestorId = null
  for (let i = targetPath.length - 1; i >= 0; i--) {
    if (oldPath.has(targetPath[i].id)) {
      commonAncestorId = targetPath[i].id
      break
    }
  }

  // 3. 从 oldLeaf 沿 parentId 链回溯到公共祖先，收集 entry
  const entries = []
  let current = oldLeafId
  while (current && current !== commonAncestorId) {
    entries.push(session.getEntry(current))
    current = entry.parentId
  }
  entries.reverse() // 恢复时间顺序

  return { entries, commonAncestorId }
}
```

**Token 预算管理 `prepareBranchEntries()`：**

当分支太长时，不能全部摘要化。Pi 从**最新到最旧**遍历 entry，直到超出 token 预算：

```typescript
function prepareBranchEntries(entries, tokenBudget = 0): BranchPreparation {
  // 第一遍: 从所有 entry 提取文件操作（即使不在 token 预算内）
  //   → 确保文件追踪不丢失
  for (const entry of entries) {
    if (entry.type === 'branch_summary' && entry.details) {
      // 累积嵌套分支摘要的文件追踪
      for (const f of entry.details.readFiles) fileOps.read.add(f)
      for (const f of entry.details.modifiedFiles) fileOps.edited.add(f)
    }
  }

  // 第二遍: 从最新到最旧添加消息，直到超出预算
  for (let i = entries.length - 1; i >= 0; i--) {
    const message = getMessageFromEntry(entries[i])
    const tokens = estimateTokens(message)
    if (tokenBudget > 0 && totalTokens + tokens > tokenBudget) {
      // 特殊处理: 如果是摘要类 entry 且预算还有 90% 未用，强行塞入
      if (
        (entry.type === 'compaction' || entry.type === 'branch_summary') &&
        totalTokens < tokenBudget * 0.9
      ) {
        messages.unshift(message) // 摘要是重要上下文，优先保留
      }
      break
    }
    messages.unshift(message)
    totalTokens += tokens
  }

  return { messages, fileOps, totalTokens }
}
```

> **关键洞见：文件操作和消息内容使用不同的预算策略。** 文件操作总是全量收集（成本极低），而消息内容受 token 预算限制。这确保了即使分支很长，文件追踪也不会丢失。

### 6.11 溢出恢复的"一次机会"原则

```typescript
// agent-session.ts
private async _checkCompaction(assistantMessage, skipAbortedCheck = true) {
  // Case 1: 溢出
  if (sameModel && !errorIsFromBeforeCompaction && isContextOverflow(msg, contextWindow)) {
    if (this._overflowRecoveryAttempted) {
      // 已经试过了，放弃
      emit({ type: "auto_compaction_end", errorMessage:
        "Context overflow recovery failed after one compact-and-retry attempt." })
      return
    }
    this._overflowRecoveryAttempted = true

    // 移除错误消息（agent state 中移除，session 历史保留）
    const messages = this.agent.state.messages
    if (messages.last.role === "assistant") {
      this.agent.replaceMessages(messages.slice(0, -1))
    }

    // 压缩 → 自动重试 (willRetry = true)
    await this._runAutoCompaction("overflow", true)
  }
}
```

**为什么只给一次机会？** 如果压缩后仍然溢出，说明 `keepRecentTokens` 的消息本身就超过了上下文窗口，压缩无法解决——需要用户切换到更大上下文窗口的模型。无限重试只会浪费 API 费用。

### 6.12 扩展自定义压缩

Pi 的压缩系统完全可被扩展接管。这是通过两个事件钩子实现的：

**`session_before_compact`** —— 压缩前，可拦截或替换

```typescript
// 扩展可以:
pi.on('session_before_compact', async (event, ctx) => {
  const { preparation, branchEntries, signal } = event

  // 1. 取消压缩
  return { cancel: true }

  // 2. 提供自定义摘要
  return {
    compaction: {
      summary: '自定义摘要...',
      firstKeptEntryId: preparation.firstKeptEntryId,
      tokensBefore: preparation.tokensBefore,
      details: {
        /* 任意 JSON 数据 */
      }
    }
  }
})
```

**`session_compact`** —— 压缩后通知

```typescript
pi.on('session_compact', async event => {
  // event.compactionEntry - 保存的压缩条目
  // event.fromExtension - 是否由扩展生成
  // 可用于更新扩展内部状态
})
```

扩展自定义压缩的 `details` 字段可以存储任意 JSON 数据（如 artifact 索引、版本标记等），Pi 会原样保存和传递。

### 6.13 配置参数

```json
// ~/.pi/agent/settings.json 或 <project>/.pi/settings.json
{
  "compaction": {
    "enabled": true, // 是否启用自动压缩
    "reserveTokens": 16384, // 预留给 LLM 回复的 token 数
    "keepRecentTokens": 20000 // 保留最近多少 token 不压缩
  }
}
```

| 参数               | 默认值  | 含义                                                                                |
| ------------------ | ------- | ----------------------------------------------------------------------------------- |
| `enabled`          | `true`  | 自动压缩开关。关闭后仍可手动 `/compact`                                             |
| `reserveTokens`    | `16384` | 上下文窗口中预留的 token。触发条件: `contextTokens > contextWindow - reserveTokens` |
| `keepRecentTokens` | `20000` | 压缩时保留最近的 token 量。决定切割点位置                                           |

**参数调优建议：**

- `keepRecentTokens` 太小 → 保留的上下文太少，LLM 可能丢失工作记忆
- `keepRecentTokens` 太大 → 留给摘要的空间太少，可能紧接着又触发压缩
- `reserveTokens` 太小 → LLM 回复空间不够，可能被截断
- `reserveTokens` 太大 → 上下文窗口浪费过多

### 6.14 端到端流程图

```
用户发送消息
    │
    ▼
AgentSession.prompt()
    │
    ▼
Agent 循环 (pi-agent-core)
    │ ← model.complete() 调用 LLM
    ▼
agent_end 事件
    │
    ├── _checkCompaction()
    │     │
    │     ├── 溢出? ──→ _runAutoCompaction("overflow", willRetry=true)
    │     │                │
    │     │                ├── prepareCompaction()
    │     │                ├── 扩展钩子 session_before_compact
    │     │                ├── compact() 或扩展提供
    │     │                ├── sessionManager.appendCompaction()
    │     │                ├── buildSessionContext() → replaceMessages()
    │     │                ├── 扩展钩子 session_compact
    │     │                └── agent.continue() (自动重试)
    │     │
    │     └── 超阈值? ──→ _runAutoCompaction("threshold", willRetry=false)
    │                      │
    │                      └── (同上，但不自动重试)
    │
    ▼
等待下一次用户输入（或手动 /compact）
```

### 6.15 与其他 Coding Agent 的压缩策略对比

| 设计维度   | Pi                                    | Claude Code       | Cursor       |
| ---------- | ------------------------------------- | ----------------- | ------------ |
| 触发方式   | 自动阈值 + 溢出恢复 + 手动 `/compact` | 仅手动 `/compact` | 服务器端自动 |
| 摘要策略   | 增量更新（保留前次摘要）              | 全量重写          | 不透明       |
| 文件追踪   | 累积追踪 `read/modified`              | 无                | 不透明       |
| 分支感知   | 树形分支摘要                          | 无分支概念        | 无分支概念   |
| 扩展定制   | 完全可拦截/替换                       | 无                | 无           |
| Split Turn | 并行双摘要合并                        | 无                | 不透明       |
| Token 估算 | Usage + chars/4 混合                  | 未公开            | 未公开       |

Pi 的压缩设计最核心的差异化在于：**增量摘要 + 文件累积追踪 + 扩展可定制**。这使得长会话的上下文质量不会随着压缩次数增加而衰减。

---

## 七、扩展系统 (Extensions)

### 7.1 加载机制

```
~/.pi/agent/extensions/*.ts     ← 用户全局扩展
.pi/extensions/*.ts             ← 项目级扩展
npm:package-name                ← npm 包扩展
```

通过 `jiti`（运行时 TypeScript 编译器）加载，使得扩展可以直接用 TypeScript 编写而无需预编译。

对于 Bun 编译的二进制分发，使用 `virtualModules` 将依赖包注入到 jiti 的模块解析中：

```typescript
const VIRTUAL_MODULES = {
  '@sinclair/typebox': _bundledTypebox,
  '@mariozechner/pi-agent-core': _bundledPiAgentCore,
  '@mariozechner/pi-tui': _bundledPiTui
  // ...
}
```

### 7.2 Extension API 表面

扩展通过工厂函数接收 `ExtensionAPI (pi)` 对象：

```typescript
export default function (pi: ExtensionAPI) {
  // 注册工具
  pi.registerTool({ name, description, parameters, execute })
  // 注册命令（slash command）
  pi.registerCommand({ name, description, handler })
  // 注册消息渲染器（自定义 UI）
  pi.registerMessageRenderer({ customType, render })
  // 监听事件
  pi.on('turn_start', handler)
  pi.on('tool_call', handler) // 可拦截/修改工具调用
  pi.on('tool_result', handler) // 可修改工具返回
  pi.on('input', handler) // 拦截用户输入
  // 会话控制
  pi.sendMessage(msg)
  pi.sendUserMessage(content)
  // ...
}
```

### 关键洞见 #8：工具的双层包装

```
原始工具 → wrapRegisteredTools()  → 扩展注册的工具获得事件钩子
         → wrapToolsWithExtensions() → 所有工具（含内建）获得 tool_call/tool_result 拦截
```

扩展可以**拦截任何工具调用**的输入和输出。这使得"审计"、"沙盒"、"远程执行"等场景成为可能。例如，一个 SSH 扩展可以把所有 bash 命令重定向到远程服务器。

---

## 八、技能系统 (Skills)

### 8.1 设计理念

Skills 不是自动触发的——它们是**用户主动调用的知识包**：

```
/skill:frontend-design 创建一个登录页
```

调用时，Pi 会读取 SKILL.md 文件内容，包装成 XML 块注入到用户消息前：

```xml
<skill name="frontend-design" location="/path/to/SKILL.md">
References are relative to /path/to/.
[skill content...]
</skill>

创建一个登录页
```

### 8.2 发现机制

```
~/.pi/agent/skills/*/SKILL.md   ← 用户级
.pi/skills/*/SKILL.md           ← 项目级
直接 .md 文件                    ← 也可以
```

在系统提示词中，Skills 被格式化为列表呈现给 LLM，让 LLM 知道有哪些技能可供用户调用。但**LLM 自身不会主动调用 Skill**——这是一个有意的设计约束，避免 LLM 在不必要时加载大量技能内容。

---

## 九、System Prompt 构建

### 9.1 动态组装

系统提示词不是静态字符串，而是根据当前状态动态构建的：

```
基础模板（角色定义 + 可用工具列表 + 指导原则）
+ 工具相关提示片段（toolSnippets / promptGuidelines）
+ 追加系统提示（appendSystemPrompt）
+ 项目上下文文件（AGENTS.md / CLAUDE.md）
+ 技能列表（当 read 工具可用时）
+ 当前时间 + 工作目录
```

### 关键洞见 #9：工具感知的指导原则

系统提示词中的**指导原则随可用工具集动态调整**：

```typescript
if (hasBash && !hasGrep && !hasFind && !hasLs) {
  addGuideline('Use bash for file operations like ls, rg, find')
} else if (hasBash && (hasGrep || hasFind || hasLs)) {
  addGuideline('Prefer grep/find/ls tools over bash (faster, respects .gitignore)')
}
```

当只有 bash 时，告诉 LLM 用 bash 做文件操作。当有专用工具时，引导 LLM 优先用专用工具。这**避免了 LLM 在有更好选择时仍然 `cat file.py | grep`**。

---

## 十、三种运行模式

### 10.1 Interactive Mode（TUI）

基于自研 `pi-tui` 库，包含 34+ UI 组件：

- 消息流容器（聊天、待处理、状态）
- 自定义编辑器（多行输入、自动补全、粘贴图片）
- 底部栏（模型信息、token 用量、快捷键）
- 各种选择器（模型、会话、主题、扩展）

### 10.2 Print Mode

标准 I/O 模式，适合脚本和管道：

```bash
echo "Explain this codebase" | pi --print
cat prompt.md | pi --print --model claude-3-sonnet
```

### 10.3 RPC Mode

JSON-RPC 协议，允许外部程序（IDE 插件、Web UI）驱动 Agent。

### 关键洞见 #10：模式无关的 AgentSession

三种模式共享同一个 `AgentSession` 实例。模式只负责 I/O（怎么展示、怎么接收输入），不负责业务逻辑。这意味着：

- Print 模式的 Agent 和 Interactive 模式的 Agent 行为**完全一致**
- 可以在 RPC 模式中实现 Web UI，而核心逻辑零改动
- 测试可以直接对 AgentSession 编写，不需要模拟 UI

---

## 十一、自动重试

### 11.1 处理策略

Pi 区分两类错误：

| 错误类型                 | 处理方式              |
| ------------------------ | --------------------- |
| 上下文溢出               | 压缩 → 重试（仅一次） |
| 可重试错误(429/5xx/过载) | 指数退避重试          |
| 其他错误                 | 直接报给用户          |

### 11.2 指数退避

```typescript
const delayMs = settings.baseDelayMs * 2 ** (this._retryAttempt - 1)
// 第1次: baseDelay
// 第2次: baseDelay * 2
// 第3次: baseDelay * 4
// ...
```

重试期间可被用户中止（`abortRetry()`），中止会立即 resolve retry Promise。

---

## 十二、设计理念总结

### 12.1 核心原则

1. **事件驱动，非回调耦合**：Agent 循环不知道 UI/持久化的存在，通过事件解耦
2. **追加制不可变存储**：会话文件只追加不修改，保证崩溃安全和完整审计
3. **可组合的工具链**：工具注册表 + 扩展包装 = 任何工具都可被拦截/替换
4. **渐进式降级**：从精确匹配到模糊匹配(edit)，从阈值压缩到溢出恢复(compaction)
5. **用户永远在控制台**：Steer/FollowUp 双轨制确保用户随时可以介入

### 12.2 值得学习的工程实践

- **Promise 链串行化事件处理**：解决异步事件的顺序问题
- **同步创建 + 异步 resolve 的 Promise 模式**：解决竞态条件
- **jiti + virtualModules**：实现"在编译后的二进制中运行 TypeScript 扩展"
- **增量压缩摘要**：避免信息随压缩次数增加而衰减
- **工具感知的 System Prompt**：根据可用工具动态调整 LLM 行为指导

### 12.3 架构取舍

| 选择               | 得到                 | 放弃                 |
| ------------------ | -------------------- | -------------------- |
| JSONL 而非 SQLite  | 崩溃安全、易于调试   | 查询性能、文件膨胀   |
| 追加制而非 CRUD    | 数据不丢失、支持分支 | 存储效率             |
| jiti 加载扩展      | 用户直接写 TS        | 启动速度、安全沙盒   |
| chars/4 估算 token | 快速、零依赖         | 精确性（但偏保守）   |
| Fuzzy edit match   | 容忍 LLM 的微小偏差  | 极端情况下可能误匹配 |

---

## 十三、代码结构速查

```
src/
├── main.ts              # CLI 入口：参数解析 → 配置 → 模型 → 模式分发
├── cli.ts               # CLI 路由
├── config.ts            # 安装检测、目录解析
├── index.ts             # 公共 API 导出
├── migrations.ts        # JSONL 格式迁移
│
├── core/
│   ├── agent-session.ts # ⭐ 核心协调层 (~2500 行)
│   ├── session-manager.ts # JSONL 追加制会话存储 + 树操作
│   ├── settings-manager.ts # 用户设置（JSONL）
│   ├── system-prompt.ts # 动态 System Prompt 构建
│   ├── model-registry.ts # 模型注册与 API Key 管理
│   ├── model-resolver.ts # 模型名称解析
│   ├── bash-executor.ts # Bash 执行 + 流式截断
│   ├── resource-loader.ts # 资源懒加载（扩展/技能/提示/主题/上下文）
│   ├── prompt-templates.ts # 提示词模板
│   ├── skills.ts        # 技能发现与加载
│   ├── event-bus.ts     # 通用事件总线
│   ├── messages.ts      # 自定义消息类型定义
│   │
│   ├── tools/           # 七个内建工具
│   │   ├── bash.ts      # shell 执行 + 输出截断
│   │   ├── read.ts      # 文件读取
│   │   ├── edit.ts      # 精确文本替换 + fuzzy match
│   │   ├── write.ts     # 文件创建/覆盖
│   │   ├── grep.ts      # 内容搜索
│   │   ├── find.ts      # 文件查找
│   │   ├── ls.ts        # 目录列表
│   │   ├── edit-diff.ts # 差异计算 + fuzzy 匹配
│   │   ├── truncate.ts  # 输出截断策略
│   │   └── path-utils.ts
│   │
│   ├── compaction/      # 上下文压缩
│   │   ├── compaction.ts  # 核心算法：切割点、摘要生成
│   │   ├── branch-summarization.ts # 分支摘要
│   │   └── utils.ts     # 文件操作追踪、序列化
│   │
│   └── extensions/      # 扩展系统
│       ├── loader.ts    # TypeScript 扩展加载(jiti)
│       ├── runner.ts    # 扩展生命周期管理
│       ├── wrapper.ts   # 工具包装（拦截层）
│       └── types.ts     # 扩展 API 类型定义
│
├── modes/
│   ├── interactive/     # TUI 交互模式
│   │   ├── interactive-mode.ts  # 主循环 + 事件处理
│   │   ├── components/  # 34+ TUI 组件
│   │   └── theme/       # 主题系统
│   ├── print-mode.ts    # 标准 I/O 模式
│   └── rpc-mode.ts      # JSON-RPC 模式
│
├── cli/                 # CLI 辅助
│   ├── args.ts          # 参数解析
│   ├── session-picker.ts # 会话选择
│   └── list-models.ts   # 模型列表
│
└── utils/               # 工具函数
    ├── shell.ts         # Shell 配置与环境
    ├── git.ts           # Git 操作
    ├── clipboard*.ts    # 剪贴板（含图片）
    ├── frontmatter.ts   # YAML frontmatter 解析
    └── image-*.ts       # 图片处理/转换
```

---

## 十四、关键功能实现深潜

**新增洞见 #11~#23：**

- **#11 流式消息就地突变**：push on start → replace reference on delta，避免高频数组操作
- **#12 被跳过的工具仍入上下文**：LLM API 要求每个 tool_call 必有 tool_result
- **#13 延迟刷新**：只在有助手回复后才写盘，避免空会话文件
- **#14 分支只是移动指针**：`branch()` 只需一行 `this.leafId = branchFromId`
- **#15 tool_call 拦截是阻塞的**：扩展返回 block → 直接 throw error
- **#16 tool_result 链式修改**：多扩展中间件管道模式
- **#17 双重身份消息转换**：bashExecution → user, compactionSummary → `<context_checkpoint>` 标签
- **#18 transformContext 扩展注入点**：每次 LLM 调用前扩展可增删改消息
- **#19 AGENTS.md/CLAUDE.md 互为后备**：务实的行业兼容
- **#20 扩展冲突不阻塞**：宽容加载 + 诊断报告
- **#21 Bash 模式实时检测**：输入 `!` 瞬间边框变色
- **#22 Extension UI 双向 RPC 桥**：RPC 模式下扩展 UI 请求的 Promise 桥接
- **#23 API Key 动态解析**：支持 OAuth token 自动刷新

还涵盖了 **5 个关键设计模式**（事件队列串行化、同步创建+异步 resolve、就地突变、中间件管道、延迟批量写入）和与其他 Coding Agent 的设计对比表。

### 14.1 Agent Loop：三层嵌套控制流的精确实现

pi-agent-core 中的 Agent Loop 是整个系统的心脏。它的精妙在于用**两层 while 循环 + 逐工具中断检查**实现了复杂的并发交互语义。

```typescript
// 伪代码还原核心循环
async function agentLoop(prompts, context, config) {
  // 添加用户消息
  for (const prompt of prompts) {
    context.messages.push(prompt)
  }
  emit('agent_start')

  // ═══ 外层循环：Follow-Up 队列驱动 ═══
  while (true) {
    // ═══ 内层循环：Tool Call + Steering ═══
    while (true) {
      emit('turn_start')

      // ① 消息转换管线
      const transformed = await transformContext(context.messages)
      const llmMessages = convertToLlm(transformed)

      // ② 流式推理
      const assistant = await streamAssistantResponse(model, llmMessages)
      // 流式过程中：partial message 就地 mutate（push on start, update on delta）

      // ③ 检查终止条件
      if (assistant.stopReason === 'error' || assistant.stopReason === 'aborted') {
        emit('turn_end', { message: assistant, toolResults: [] })
        emit('agent_end')
        return
      }

      // ④ 提取工具调用
      const toolCalls = extractToolCalls(assistant)
      if (toolCalls.length === 0) {
        emit('turn_end', { message: assistant, toolResults: [] })
        break // 无工具调用 → 退出内层循环
      }

      // ⑤ 逐工具执行 + 逐工具中断检查（关键！）
      const results = []
      for (let i = 0; i < toolCalls.length; i++) {
        const result = await tool.execute(toolCalls[i], signal)
        results.push(result)
        context.messages.push(toolResult)

        // ⑥ 每个工具执行完后检查 Steering
        const steering = await getSteeringMessages()
        if (steering.length > 0) {
          // 跳过剩余工具（标记为 skipped error）
          for (const skipped of toolCalls.slice(i + 1)) {
            results.push(skipToolCall(skipped))
          }
          // 注入 steering 消息
          for (const msg of steering) {
            context.messages.push(msg)
          }
          break // 回到内层循环顶部 → 新的 LLM 调用
        }
      }

      emit('turn_end', { message: assistant, toolResults: results })
      // 内层循环继续 → 下一轮 LLM 调用处理工具结果
    }

    // ⑦ 内层循环结束 → 检查 Follow-Up
    const followUps = await getFollowUpMessages()
    if (followUps.length === 0) break // 无后续 → 退出外层循环

    // 注入 follow-up 消息，继续外层循环
    for (const msg of followUps) {
      context.messages.push(msg)
    }
  }

  emit('agent_end')
}
```

#### 洞见 #11：流式消息的就地突变策略

```typescript
// 不是：push → pop → push（浪费）
// 而是：push on start → mutate in-place on delta → replace on done
case "start":
    partialMessage = event.partial;
    context.messages.push(partialMessage);        // 添加
case "text_delta":
    partialMessage = event.partial;
    context.messages[context.messages.length - 1] = partialMessage;  // 替换引用
case "done":
    context.messages[context.messages.length - 1] = finalMessage;    // 替换为最终版
```

Partial message 在流式过程中被**添加一次，然后就地更新**。这允许订阅者通过引用持有的对象看到实时更新，同时避免了频繁的数组操作。当 delta 到来时，替换的是数组最后一个元素的引用，而不是 push/pop。

#### 洞见 #12：被跳过的工具调用仍然进入上下文

当 Steering 打断工具链时，剩余的工具调用不是简单丢弃——它们被标记为 **"skipped" error result** 加入上下文：

```typescript
// 被跳过的工具会产生一个 error toolResult
function skipToolCall(toolCall): ToolResultMessage {
  return {
    role: 'toolResult',
    toolCallId: toolCall.id,
    content: [{ type: 'text', text: 'Tool execution skipped (user interrupted)' }],
    isError: true
  }
}
```

**为什么？** LLM API 要求每个 tool_call 必须有对应的 tool_result。如果只跳过不加结果，LLM 会报错或产生幻觉。这是对 API 协议的严格遵守。

---

### 14.2 会话树：从 JSONL 到可导航的树形结构

#### 树遍历核心算法（buildSessionContext）

```
JSONL 文件条目（时序追加）:
  [session] → [msg:1] → [msg:2] → [msg:3] → [compaction:4] → [msg:5] → [msg:6]
                  ↑                                                ↑
                  └── branch point                          leafId 指向这里

路径构建（从 leaf 到 root）:
  leaf(msg:6) → parentId → compaction:4 → parentId → msg:3 → ... → msg:1 → null
  反转得到: [msg:1, msg:2, msg:3, compaction:4, msg:5, msg:6]

消息构建（如果有 compaction）:
  Phase 1: 发出 compaction summary 消息
  Phase 2: 发出 firstKeptEntryId 到 compaction 之间的"保留"消息
  Phase 3: 发出 compaction 之后的消息

最终 LLM 看到:
  [compaction_summary, kept_msg_3, msg_5, msg_6]
  （msg_1, msg_2 被摘要替代）
```

#### 洞见 #13：延迟刷新——只在有助手回复后才写盘

```typescript
_persist(entry: SessionEntry): void {
    const hasAssistant = this.fileEntries.some(
        e => e.type === "message" && e.message.role === "assistant"
    );
    if (!hasAssistant) {
        this.flushed = false;
        return; // 不写！
    }
    if (!this.flushed) {
        // 首次：批量写入所有积压条目
        for (const e of this.fileEntries) {
            appendFileSync(this.sessionFile, JSON.stringify(e) + "\n");
        }
        this.flushed = true;
    } else {
        appendFileSync(this.sessionFile, JSON.stringify(entry) + "\n");
    }
}
```

只有当 LLM **确实产出了回复**，才把整个会话写入磁盘。这避免了：

- 用户输入一句话就退出，产生大量空会话文件
- 会话头和用户消息的无意义磁盘写入
- 一旦有助手回复，**批量追加**所有积压条目（原子性更好）

#### 洞见 #14：分支只是移动指针

```typescript
branch(branchFromId: string): void {
    this.leafId = branchFromId; // 就这一行！
}
```

创建分支不需要复制数据、创建新文件或修改任何已有条目。只需要把 `leafId` 指回某个历史条目。下一次 `append()` 时，新条目的 `parentId` 指向这个历史条目，自然形成分支。

如果需要记录"为什么分支"，用 `branchWithSummary()`，它会额外追加一个 `BranchSummaryEntry`，其中的 summary 告诉 LLM "之前这条路你做了什么、为什么换方向了"。

---

### 14.3 工具拦截管线：从注册到执行的完整流水线

```
                    工具注册表构建
                    ─────────────
内建工具 (7个)                扩展注册的工具
  read, bash,                pi.registerTool(def)
  edit, write,                    │
  grep, find, ls                  ▼
       │                  wrapRegisteredTool()
       │                  给每个扩展工具注入运行时上下文
       │                          │
       ▼                          ▼
    baseToolRegistry ──────── toolRegistry (合并)
                                  │
                                  ▼
                       wrapToolsWithExtensions()
                       ┌──────────────────────┐
                       │ 对每个工具创建拦截层：  │
                       │                      │
                       │ execute前:            │
                       │   emit("tool_call")   │
                       │   → 扩展可 block      │
                       │                      │
                       │ execute 原始工具       │
                       │                      │
                       │ execute后:            │
                       │   emit("tool_result") │
                       │   → 扩展可修改返回值   │
                       └──────────────────────┘
                                  │
                                  ▼
                    setActiveToolsByName(names)
                    → 只激活子集 → agent.setTools()
```

#### 洞见 #15：tool_call 拦截是同步阻塞的

```typescript
// wrapper.ts 中
const callResult = await runner.emitToolCall({
  type: 'tool_call',
  toolName: tool.name,
  toolCallId,
  input: params
})

if (callResult?.block) {
  throw new Error(callResult.reason || 'Blocked by extension')
}
```

当扩展返回 `{ block: true, reason: "..." }` 时，工具执行**直接抛出错误**。这个错误会变成 `toolResult` 中的 `isError: true` 返回给 LLM。LLM 看到的是"工具执行失败"，而不是"工具被阻止"——从 LLM 视角它们是等价的。

#### 洞见 #16：tool_result 链式修改

```typescript
// runner.ts 中
for (const handler of handlers) {
  const handlerResult = await handler(currentEvent, ctx)
  if (handlerResult?.content !== undefined) {
    currentEvent.content = handlerResult.content // 覆盖内容
    modified = true
  }
  // 下一个 handler 看到的是上一个 handler 修改后的结果
}
```

多个扩展可以**链式修改**工具返回结果。扩展 A 修改了 content，扩展 B 看到的 event.content 已经是 A 修改后的。这是一个**中间件管道**模式。

---

### 14.4 消息转换管线：从领域消息到 LLM 消息

Pi 维护了两层消息类型系统：

```
AgentMessage (领域层)              Message (LLM 层)
─────────────────                 ──────────────
user          ─────→             user
assistant     ─────→             assistant
toolResult    ─────→             toolResult
bashExecution ──┐                (无对应)
custom        ──┤ convertToLlm()
branchSummary ──┤──→             user (包装为文本)
compactionSummary─┘
```

#### 洞见 #17："双重身份"消息的转换策略

```typescript
export function convertToLlm(messages: AgentMessage[]): Message[] {
  return messages
    .map(m => {
      switch (m.role) {
        case 'bashExecution':
          if (m.excludeFromContext) return undefined // !! 前缀
          return {
            role: 'user',
            content: [
              {
                type: 'text',
                text:
                  `$ ${m.command}\n${m.output}\n` +
                  (m.exitCode !== undefined ? `[exit code: ${m.exitCode}]` : '') +
                  (m.cancelled ? ' [cancelled]' : '') +
                  (m.truncated ? ' [truncated]' : '')
              }
            ]
          }

        case 'compactionSummary':
          return {
            role: 'user',
            content: [
              {
                type: 'text',
                text:
                  '<context_checkpoint>\n' +
                  'Below is a summary of our conversation so far...\n' +
                  m.summary +
                  '\n</context_checkpoint>'
              }
            ]
          }

        case 'branchSummary':
          return {
            role: 'user',
            content: [
              {
                type: 'text',
                text:
                  '<branch_context>\n' +
                  'The user navigated to a different point in our conversation...\n' +
                  m.summary +
                  '\n</branch_context>'
              }
            ]
          }

        default:
          return m // user/assistant/toolResult 直接透传
      }
    })
    .filter(Boolean)
}
```

**关键设计**：

1. **Bash 执行 → user 消息**：用户在终端执行的命令和输出被格式化为文本，注入到对话中，让 LLM 知道用户做了什么
2. **`!!` 双叹号 → 完全跳过**：`excludeFromContext` 标记的 bash 执行不会进入 LLM 上下文（隐私/无关命令）
3. **压缩摘要 → `<context_checkpoint>` 标签**：特殊标签让 LLM 明确知道这是之前对话的摘要，不是用户的新消息
4. **分支摘要 → `<branch_context>` 标签**：让 LLM 知道对话发生了方向切换

#### 洞见 #18：transformContext 是扩展的"上下文注入点"

```typescript
// SDK 创建 Agent 时
const agent = new Agent({
  transformContext: async messages => {
    const runner = extensionRunnerRef.current
    if (!runner) return messages
    return runner.emitContext(messages) // 扩展可以增删改消息
  }
})
```

每次 LLM 调用前，`transformContext` 被调出。扩展可以：

- 注入额外上下文（如实时代码分析结果）
- 过滤敏感信息
- 重排消息优先级
- 限制上下文窗口用量

---

### 14.5 资源发现：AGENTS.md 向上遍历算法

```
项目结构:
  /home/user/                       ~/.pi/agent/AGENTS.md  (全局)
  /home/user/project/               AGENTS.md              (祖先级)
  /home/user/project/packages/      (无)
  /home/user/project/packages/foo/  AGENTS.md              (项目级)
  /home/user/project/packages/foo/src/  ← CWD

发现顺序（从 CWD 向上）:
  src/ → 无
  foo/ → 发现 AGENTS.md ✓ (project)
  packages/ → 无
  project/ → 发现 AGENTS.md ✓ (ancestor)
  user/ → 无
  home/ → 无
  / → 无

加载到 system prompt 的顺序:
  [~/.pi/agent/AGENTS.md, /project/AGENTS.md, /project/packages/foo/AGENTS.md]
  从最外层到最内层，内层可覆盖外层指导
```

#### 洞见 #19：`AGENTS.md` 和 `CLAUDE.md` 互为后备

在每个目录中，先查找 `AGENTS.md`，找不到再查找 `CLAUDE.md`。这是对行业生态的务实兼容——很多项目已经有 `CLAUDE.md`，Pi 不强制用户改名。

#### 洞见 #20：扩展冲突检测

```typescript
// resource-loader.ts
const conflicts = this.detectExtensionConflicts(extensions)
// 检测：同名工具、同名命令、同名标志
// 结果：不阻止加载，而是作为诊断信息报告
```

冲突不会导致加载失败。**所有扩展都会被加载**，但冲突会被记录为诊断信息。这是"宽容加载、诊断报告"的策略——让用户知道问题但不阻塞工作流。

---

### 14.6 Interactive Mode：TUI 事件循环

#### 初始化链路

```
init()
├─ loadChangelog()          → 版本更新提示
├─ ensureTools()            → 检查 fd/rg 是否安装
├─ setupLayout()            → 创建 UI 容器层级
│   ├─ headerContainer      → Logo + 帮助文字
│   ├─ chatContainer        → 主消息流
│   ├─ pendingMsgContainer  → 排队中的消息显示
│   ├─ statusContainer      → 压缩/重试状态
│   └─ widgetContainers[]   → 扩展 widget
├─ setupKeyHandlers()       → 快捷键绑定
├─ setupEditorSubmit()      → 编辑器提交处理
├─ initExtensions()         → 扩展生命周期启动
├─ renderInitialMessages()  → 恢复历史消息的 UI 渲染
├─ ui.start()               → TUI 事件循环启动
└─ subscribeToAgent()       → 事件订阅（最后一步！）
```

#### 流式渲染三阶段

```
message_start (助手开始说话)
    → 创建 AssistantMessageComponent
    → 添加到 chatContainer
    → 首次渲染

message_update (每个 token 到达)
    → 更新 streamingMessage 引用
    → component.updateContent(message)  // 增量渲染
    → 检测新的 toolCall → 创建 ToolExecutionComponent

message_end (助手说完)
    → 标记所有 pending tool 的参数完成
    → 清空 streamingComponent 引用

tool_execution_start → 创建工具执行 UI 组件
tool_execution_update → 流式更新输出
tool_execution_end → 最终结果 + 清理
```

#### 洞见 #21：Bash 模式的实时检测

```typescript
// 编辑器每次 onChange 都检测
this.defaultEditor.onChange = (text: string) => {
  const wasBashMode = this.isBashMode
  this.isBashMode = text.trimStart().startsWith('!')
  if (wasBashMode !== this.isBashMode) {
    this.updateEditorBorderColor() // 边框颜色实时切换
  }
}
```

用户在编辑器中输入 `!` 的瞬间，边框颜色就变了——告诉用户"你现在是在写 shell 命令，不是给 LLM 的消息"。这种**即时视觉反馈**是优秀 TUI 设计的典范。

---

### 14.7 RPC Mode：面向集成的协议设计

```
外部进程                    Pi (RPC Mode)
  │                              │
  │ ── stdin: JSON command ──────→ handleCommand()
  │                              │   ├─ prompt
  │                              │   ├─ steer
  │                              │   ├─ follow_up
  │                              │   ├─ abort
  │                              │   ├─ get_state
  │                              │   ├─ set_model
  │                              │   ├─ compact
  │                              │   └─ ...
  │                              │
  │ ←── stdout: JSON events ─────│ subscribe(event => JSON.stringify)
  │                              │   ├─ message_start
  │                              │   ├─ message_update
  │                              │   ├─ message_end
  │                              │   ├─ tool_execution_*
  │                              │   └─ ...
  │                              │
  │ ── stdin: extension_ui_resp ─→ pendingExtensionRequests.resolve()
  │ ←── stdout: extension_ui_req─│ （扩展需要 UI 输入时）
```

#### 洞见 #22：Extension UI 的双向 RPC 桥

当扩展在 RPC 模式下需要用户输入（比如 `pi.ui.confirm("确定删除？")`）时：

1. Pi 把请求序列化为 JSON 写到 stdout
2. 创建一个 pending Promise（存入 `pendingExtensionRequests` Map）
3. **等待**外部进程在 stdin 回一个 `extension_ui_response`
4. 收到响应后 resolve Promise

这让远端 UI（比如 Web 前端）可以弹真正的对话框，而不是被 TUI 的 stdin/stdout 困住。

---

### 14.8 SDK：嵌入式使用的 10 步初始化

```typescript
const { session } = await createAgentSession({
  cwd: '/my/project',
  model: someModel,
  customTools: [myTool]
})

// 订阅事件
session.subscribe(event => {
  if (event.type === 'message_update') renderToken(event)
})

// 发送提示
await session.prompt('请帮我重构这个函数')
```

SDK 的 `createAgentSession()` 内部经历了完整的 10 步初始化：

```
1. 路径解析 (cwd, agentDir)
2. 创建存储层 (authStorage, modelRegistry, settingsManager, sessionManager)
3. 资源加载 (resourceLoader.reload() — 异步发现扩展/技能/提示/主题)
4. 会话恢复检查 (已有 session 文件？有消息？)
5. 模型选择 (恢复上次 → 用户指定 → 自动发现 → 报错)
6. ThinkingLevel 确定 (恢复 → 设置默认 → 钳制到模型能力)
7. 消息转换包装 (convertToLlm + blockImages 过滤)
8. Agent 实例创建 (注入 transformContext, convertToLlm, getApiKey)
9. 消息恢复 (replaceMessages 或初始化新会话)
10. AgentSession 创建 (绑定 agent + sessionMgr + settings + extensions)
```

#### 洞见 #23：API Key 的动态解析

```typescript
const agent = new Agent({
  getApiKey: async provider => {
    return modelRegistry.getApiKeyForProvider(provider)
  }
})
```

API Key 不是在初始化时获取一次就缓存——它在**每次 LLM 调用前动态解析**。这支持：

- OAuth token 的自动刷新（GitHub Copilot 等）
- 运行时切换 provider 不需要重启
- 环境变量在运行中被修改后立刻生效

---

## 十五、关键设计模式汇总

### 模式 1：事件队列串行化

```typescript
this._agentEventQueue = this._agentEventQueue.then(
  () => this._processAgentEvent(event),
  () => this._processAgentEvent(event) // rejected 也继续
)
```

**问题**：Agent 事件是同步触发的，但处理（保存磁盘、通知扩展）是异步的。  
**解法**：Promise 链保证顺序。两个回调（fulfilled + rejected）保证链条不被 error 断裂。

### 模式 2：同步创建 Promise + 异步 resolve

```typescript
// 同步路径（事件回调中）
this._retryPromise = new Promise(resolve => {
  this._retryResolve = resolve
})

// 异步路径（事件队列中）
if (retrySucceeded) {
  this._retryResolve() // resolve 之前创建的 Promise
}

// 等待路径（prompt 方法中）
await this.waitForRetry() // 等待 Promise resolve
```

**问题**：`prompt()` 调用 `agent.prompt()` 后立刻需要知道是否要重试。但重试决策在异步事件队列中。  
**解法**：在同步事件回调中就创建 Promise 占位，异步处理时填充结果。

### 模式 3：就地突变 + 引用替换

```typescript
// 流式消息：添加一次，之后只替换引用
context.messages[lastIndex] = updatedPartial
// 而不是 pop + push 或 splice
```

**问题**：流式 token 到达频率极高（每秒数十次），避免频繁数组操作。  
**解法**：固定位置的引用替换。订阅者通过最后一个元素的引用看到更新。

### 模式 4：中间件管道（工具包装）

```typescript
// 链式修改：A 的输出是 B 的输入
for (const handler of handlers) {
  const result = await handler(currentEvent)
  if (result?.content) currentEvent.content = result.content
}
```

### 模式 5：延迟写入 + 批量刷新

```typescript
if (!hasAssistantMessage) return; // 延迟
if (!this.flushed) {
    for (const e of allEntries) appendFileSync(...); // 批量
    this.flushed = true;
} else {
    appendFileSync(..., entry); // 单条追加
}
```

---

## 十六、与其他 Coding Agent 的设计差异

| 设计维度   | Pi                  | Cursor/Windsurf | Claude Code   |
| ---------- | ------------------- | --------------- | ------------- |
| 会话存储   | JSONL 追加制树      | 服务器端        | 本地 JSON     |
| 分支       | 原生树结构          | 无              | 无            |
| 上下文压缩 | 增量摘要 + 文件追踪 | 服务器端        | 手动 /compact |
| 工具拦截   | 扩展双层包装        | 无公开API       | 无            |
| 运行模式   | TUI/Print/RPC       | IDE 插件        | CLI           |
| 扩展系统   | jiti 加载 TS        | IDE 原生        | MCP Server    |
| 用户中断   | Steer/FollowUp 双轨 | IDE 原生        | Ctrl+C        |

Pi 的独特定位：**一个可嵌入、可扩展、面向终端但不限于终端的 Agent 框架**。它的设计目标不只是"一个好用的 CLI 工具"，而是"一个其他产品可以构建其上的 Agent 引擎"。

---

## 十七、总结：Pi 的核心设计哲学

1. **事件驱动，不耦合** — 每一层只通过事件通信，不直接调用对方
2. **追加不修改** — JSONL 只追加；消息只替换引用，不删除
3. **用户永远优先** — Steer 可以在任何工具之间打断；!! 可以隐藏命令；abort 可以随时中止
4. **宽容加载，严格执行** — 扩展冲突不阻塞加载；工具参数由 TypeBox schema 严格验证
5. **渐进式降级** — 精确匹配 → 模糊匹配(edit)；阈值压缩 → 溢出恢复(compaction)；重试 → 放弃
6. **一次实现，多处运行** — AgentSession 是所有模式的公共层；SDK 暴露相同能力
