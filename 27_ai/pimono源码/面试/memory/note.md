# Pi-Mono Memory (Compaction) 实现深度分析

## 一、核心问题：LLM 的上下文窗口有限

Agent 在长对话中会不断积累消息（用户输入、助手回复、工具调用结果等），最终超出模型的 context window。Pi 的 Memory 方案就是 **Compaction（压缩）**——用 LLM 对旧消息生成结构化摘要，替换掉原始消息，释放上下文空间。

## 二、整体架构

```
┌──────────────────────────────────────────────────────────────┐
│                      Agent Session                           │
│  ┌─────────────┐    ┌──────────────┐    ┌────────────────┐  │
│  │ Agent Loop   │───▶│ Session Mgr  │───▶│ .jsonl 持久化   │  │
│  │ (pi-agent)  │    │ (树状结构)     │    │ (append-only)  │  │
│  └──────┬──────┘    └──────┬───────┘    └────────────────┘  │
│         │                  │                                 │
│         ▼                  ▼                                 │
│  ┌──────────────┐   ┌──────────────┐                        │
│  │ Context      │   │ Compaction   │                        │
│  │ Overflow     │   │ Engine       │                        │
│  │ Detection    │   │ (摘要生成)    │                        │
│  └──────────────┘   └──────┬───────┘                        │
│                            │                                 │
│                    ┌───────┴────────┐                        │
│                    │ Branch Summary │                        │
│                    │ (分支摘要)      │                        │
│                    └────────────────┘                        │
└──────────────────────────────────────────────────────────────┘
```

**核心文件：**
- `packages/coding-agent/src/core/compaction/compaction.ts` — 自动压缩逻辑
- `packages/coding-agent/src/core/compaction/branch-summarization.ts` — 分支摘要
- `packages/coding-agent/src/core/compaction/utils.ts` — 序列化、文件追踪
- `packages/coding-agent/src/core/session-manager.ts` — Session 树状结构 + 持久化
- `packages/coding-agent/src/core/agent-session.ts` — 触发 compaction 的控制逻辑
- `packages/coding-agent/src/core/messages.ts` — 消息类型 + `convertToLlm()`
- `packages/ai/src/utils/overflow.ts` — 上下文溢出检测

## 三、两大 Memory 机制

| 机制 | 触发条件 | 目的 |
|------|---------|------|
| **Compaction（自动压缩）** | `contextTokens > contextWindow - reserveTokens` 或手动 `/compact` | 摘要旧消息，释放上下文空间 |
| **Branch Summarization（分支摘要）** | `/tree` 导航切换分支 | 保留离开分支的上下文到新分支 |

## 四、Compaction 详细流程

### 4.1 触发条件（两种情况）

在 `agent-session.ts` 的 `_checkCompaction()` 中，每次 assistant 回复后检查：

**Case 1：Context Overflow（溢出恢复）**
- LLM 返回 context overflow 错误（由 `isContextOverflow()` 检测）
- 移除错误的 assistant message
- 执行 compaction + **自动重试**（`willRetry = true`）
- 一次性尝试，防止无限循环（`_overflowRecoveryAttempted` 标志）

**Case 2：Threshold（阈值压缩）**
- `contextTokens > contextWindow - reserveTokens`（默认 reserveTokens = 16384）
- 执行 compaction，**不自动重试**，用户继续手动输入

### 4.2 Compaction 核心算法

```
prepareCompaction() → compact() → sessionManager.appendCompaction() → reload
```

#### Step 1: 找切割点 — `findCutPoint()`

从最新消息向前遍历，累加 token 估算值（`chars / 4` 启发式），直到累计 >= `keepRecentTokens`（默认 20000）。

- 合法切割点：user / assistant / bashExecution / custom 消息
- **永远不能在 toolResult 处切割**（tool result 必须跟在 tool call 后面）
- 如果单个 turn 超过 keepRecentTokens → **Split Turn**（在 turn 中间切割）

#### Step 2: 提取消息

```
boundaryStart = prevCompactionIndex + 1  // 从上次 compaction 后面开始
boundaryEnd = entries.length

messagesToSummarize = entries[boundaryStart .. cutPoint]  // 要摘要的旧消息
keptMessages = entries[cutPoint .. end]                   // 保留的新消息
```

#### Step 3: 生成摘要 — `generateSummary()`

- 消息先经 `convertToLlm()` 转为 LLM 格式，再经 `serializeConversation()` 序列化为文本
- 序列化格式：`[User]: ...`, `[Assistant]: ...`, `[Tool result]: ...`
- **为什么序列化？** 防止 LLM 把要摘要的内容当作持续对话来回复
- 如果有上一次 compaction 的 summary → 使用 **增量更新 prompt**（`UPDATE_SUMMARIZATION_PROMPT`），而非全量重新生成
- 使用 `SUMMARIZATION_SYSTEM_PROMPT` 让模型只输出结构化摘要

#### Step 4: 结构化摘要格式

```markdown
## Goal
[用户要做什么]

## Constraints & Preferences
- [约束条件和偏好]

## Progress
### Done
- [x] [已完成]
### In Progress
- [ ] [进行中]
### Blocked
- [阻塞项]

## Key Decisions
- **[决定]**: [理由]

## Next Steps
1. [下一步]

## Critical Context
- [继续工作所需的关键信息]

<read-files>
path/to/file.ts
</read-files>

<modified-files>
path/to/changed.ts
</modified-files>
```

#### Step 5: 追加 CompactionEntry 并 Reload

- `sessionManager.appendCompaction()` 写入 JSONL
- `buildSessionContext()` 重建消息列表：**summary 消息 + firstKeptEntryId 之后的消息**

### 4.3 Compaction 前后 LLM 看到的变化

```
Before:
  [system] [msg1] [msg2] [msg3] [msg4] [msg5] [msg6] [msg7] [msg8]

After:
  [system] [summary of msg1-msg4] [msg5] [msg6] [msg7] [msg8]
                                   ↑
                          firstKeptEntryId
```

**Summary 以 user 角色注入 LLM**，包裹在 `<summary>` 标签中：
```
The conversation history before this point was compacted into the following summary:
<summary>
...结构化摘要...
</summary>
```

### 4.4 Split Turn（大 Turn 分割）

当单个 turn（一个 user 消息 + 后续所有 assistant/tool 消息）超过 `keepRecentTokens` 时：

```
entry:  [hdr] [usr] [ass] [tool] [ass] [tool] [tool] [ass] [tool]
               ↑                                      ↑
        turnStartIndex=1                     firstKeptEntryId=7

turnPrefixMessages = [usr, ass, tool, ass, tool, tool]  // 要摘要的前缀
keptMessages = [ass, tool]                               // 保留的最近部分
```

生成两个摘要并合并：
1. **History Summary**（历史摘要）— 之前的完整 turn 们
2. **Turn Prefix Summary**（轮次前缀摘要）— 当前 turn 被切掉的前半部分

## 五、Branch Summarization

### 5.1 Session 的树状结构

Session 不是线性链表，而是**树**：
- 每个 entry 有 `id` 和 `parentId`
- 分支通过 `branch(entryId)` 在历史节点创建新子节点
- 持久化格式：JSONL（append-only，不修改历史行）

```
[user] ─── [assistant] ─── [user] ─── [assistant] ─┬─ [user] ← 当前叶子
                                                    │
                                                    └─ [branch_summary] ─── [user] ← 另一个分支
```

### 5.2 导航到不同分支时生成摘要

`collectEntriesForBranchSummary()` → `prepareBranchEntries()` → `generateBranchSummary()`

1. 找 common ancestor（LCA）
2. 收集旧叶子到 LCA 之间的所有 entries
3. 从最新到最旧遍历，在 token budget 内选择消息
4. 生成摘要，追加 `BranchSummaryEntry`

### 5.3 累积文件追踪

每次 compaction/branch summary 都从以下提取文件操作：
- assistant 消息中的 tool call（read/write/edit）
- 前一次 compaction/branch summary 的 `details.readFiles` / `details.modifiedFiles`

→ 文件追踪跨多次 compaction **累积**传递

## 六、Context Overflow 检测

`packages/ai/src/utils/overflow.ts` 中的 `isContextOverflow()`：

**两种检测方式：**

1. **Error-based**：匹配各 provider 的错误消息模式（15+ 种正则）
   - Anthropic: `prompt is too long`
   - OpenAI: `exceeds the context window`
   - Google: `input token count ... exceeds the maximum`
   - xAI / Groq / OpenRouter / llama.cpp / LM Studio 等

2. **Silent overflow**：某些 provider（如 z.ai）不报错但 `usage.input > contextWindow`
   - 通过 `contextWindow` 参数比较

**溢出恢复策略：**
- 检测到 overflow → 移除错误 assistant message → 执行 compaction → 自动 `agent.continue()` 重试
- 只尝试一次，避免无限循环
- 只对当前模型的错误触发（切换模型后不误触发）
- 跳过 compaction 之前的历史错误

## 七、Token 估算机制

`estimateTokens()` — `chars / 4` 的启发式方法：
- 文本消息：`text.length / 4`
- 图片：固定 1200 tokens（4800 chars）
- tool call：`name.length + JSON.stringify(args).length`
- thinking：`thinking.length / 4`

`estimateContextTokens()` — 优先使用最后一次 assistant message 的真实 `Usage` 数据 + 后续消息的估算

`calculateContextTokens()` — 从 Usage 对象计算：`totalTokens || input + output + cacheRead + cacheWrite`

## 八、Session 持久化（JSONL 格式）

```
~/.pi/agent/sessions/--<path>--/<timestamp>_<uuid>.jsonl
```

- **Append-only**：新消息通过 `appendFileSync` 追加
- **延迟写入**：等第一条 assistant message 后才 flush 到磁盘（避免空 session 文件）
- **迁移**：v1(线性) → v2(树状+id/parentId) → v3(hookMessage→custom 重命名)

Session 读取时 `buildSessionContext()` 从叶子遍历到根，构建 LLM 能用的消息列表。遇到 CompactionEntry 时，先注入 summary，再注入 firstKeptEntryId 之后的消息。

## 九、Extension Hooks

两个核心 hook 允许扩展自定义 memory 行为：

### `session_before_compact`
- 在 compaction 执行前触发
- 可以 **取消**（`{ cancel: true }`）
- 可以提供**自定义 summary**（替代默认 LLM 生成的摘要）
- 接收 `CompactionPreparation` 对象（含待摘要消息、文件操作、设置等）

### `session_before_tree`
- 在分支导航前触发
- 可以取消导航或提供自定义分支摘要

## 十、面试关键点 & 设计亮点

### 1. 为什么不用向量数据库（RAG）做 memory？
- Coding agent 的上下文是**高度连贯的对话序列**，不是离散知识点
- 需要保留的是**goal、progress、decisions** 这些结构化状态
- 摘要提取比向量检索更适合保持任务连续性

### 2. 增量 vs 全量摘要
- 首次 compaction：全量摘要（`SUMMARIZATION_PROMPT`）
- 后续 compaction：**增量更新**（`UPDATE_SUMMARIZATION_PROMPT`），基于前一次 summary + 新消息
- 好处：信息不会在多次 compaction 中丢失

### 3. Split Turn 处理
- 单个超长 turn 不能简单丢弃
- 并行生成 history summary + turn prefix summary，然后合并
- 保证即使单 turn 超出 budget 也不丢失关键上下文

### 4. 树状 Session 结构
- 不是链表（线性覆盖），而是 **树**（保留所有历史分支）
- 支持分支导航、回溯、fork
- `append-only` JSONL 确保数据安全性

### 5. 防御性的 Overflow 恢复
- 多 provider 兼容的正则匹配
- 静默溢出检测（z.ai 场景）
- 单次重试限制，避免死循环
- 模型切换后不误触发

### 6. 文件操作的累积追踪
- 每次 compaction 记录 `readFiles` + `modifiedFiles`
- 下次 compaction 从前一次的 details 继承
- LLM 始终知道整个 session 中操作过哪些文件

### 7. 序列化防止对话延续
- `serializeConversation()` 将多轮对话序列化为 `[User]: ... [Assistant]: ...` 文本格式
- 防止 summarization 模型把要压缩的内容"继续回复"下去

### 8. Extension 可扩展性
- `session_before_compact` hook 允许自定义压缩策略（如使用不同模型、不同摘要格式）
- `fromHook` 标记区分 pi-native 和 extension-generated 的 compaction
