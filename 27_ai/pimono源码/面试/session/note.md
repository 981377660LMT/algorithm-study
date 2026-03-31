# Pi-Mono Session & Compaction 深度分析

## 一、Session 架构总览

### 1.1 存储格式：Append-Only JSONL Tree

Session 的核心设计思想是 **append-only 的树结构**，存储为 JSONL（JSON Lines）文件：

```
~/.pi/agent/sessions/--<encoded-cwd>--/<timestamp>_<uuid>.jsonl
```

**为什么选择 append-only？**

- 文件只追加不修改，天然防止并发写入导致的数据损坏
- 恢复简单：断电/崩溃最多丢失最后一行
- 结合 JSONL（每行一个 JSON），解析时跳过损坏行即可

**核心数据结构**：每行是一个 `FileEntry`，第一行是 `SessionHeader`，后续是 `SessionEntry`。

```typescript
// 所有 entry 共享的基类
interface SessionEntryBase {
  type: string;
  id: string;           // 8-char hex (从 UUID 截取，碰撞检测)
  parentId: string | null;  // 形成树结构的关键
  timestamp: string;
}
```

### 1.2 Entry 类型体系

| Entry Type | 用途 | 是否参与 LLM 上下文 |
|------------|------|---------------------|
| `session` (header) | 文件元数据，cwd/version/parentSession | 否 |
| `message` | 包装 AgentMessage (user/assistant/toolResult/bashExecution) | 是 |
| `compaction` | 压缩摘要，记录 firstKeptEntryId | 是（作为 summary 注入） |
| `branch_summary` | 分支摘要，树导航时保留离开分支的上下文 | 是（作为 user message 注入） |
| `custom_message` | 扩展注入的消息 | 是（转为 user message） |
| `custom` | 扩展状态持久化 | 否 |
| `model_change` | 模型切换记录 | 否（但影响 context 构建） |
| `thinking_level_change` | 思考级别变更 | 否（但影响 context 构建） |
| `label` | 用户书签/标签 | 否 |
| `session_info` | 会话显示名称 | 否 |

### 1.3 树结构 vs 线性历史

**关键设计决策：同文件内树结构**。

传统方案（如 Claude Code 早期）每次分支创建新文件。Pi 选择了单文件树：

```
entry: [user1] ──► [asst1] ──► [user2] ──► [asst2] ──┬──► [user3a] ← leaf
                                                       │
                                                       └──► [branch_summary] ──► [user3b]
```

**优势**：

1. 单文件便于备份和传输
2. append-only 保证历史不丢失
3. 分支只需移动 leaf 指针，无需复制数据

**`branch()` 实现极其简洁**：

```typescript
branch(branchFromId: string): void {
  this.leafId = branchFromId;  // 就这一行
}
```

## 二、SessionManager 核心实现

### 2.1 上下文构建：buildSessionContext()

这是连接 session 存储和 LLM 调用的桥梁。算法：

```
1. 从 leafId 通过 parentId 链回溯到 root → 得到 path[]
2. 遍历 path，提取 thinkingLevel / model 设置
3. 如果路径上有 CompactionEntry：
   a. 先注入 compaction.summary 作为第一条消息
   b. 然后注入 firstKeptEntryId 到 compaction 之间的 "kept messages"
   c. 最后注入 compaction 之后的消息
4. 如果没有 Compaction，直接遍历所有 entry 转换为 AgentMessage
```

**关键细节**：路径上可能有多种 entry 类型，但只有 `message`、`custom_message`、`branch_summary` 会被转换为发送给 LLM 的消息。`compaction` entry 被特殊处理——它不是直接作为消息发送，而是作为 **`compactionSummary` 角色的消息**注入到上下文开头。

### 2.2 延迟持久化策略

**一个巧妙的优化**：Session 文件不会在第一条 user message 时创建，而是**等到第一条 assistant message 出现后才写入磁盘**。

```typescript
_persist(entry: SessionEntry): void {
  const hasAssistant = this.fileEntries.some(
    e => e.type === "message" && e.message.role === "assistant"
  );
  if (!hasAssistant) {
    this.flushed = false;
    return;  // 还没有 assistant 回复，不写文件
  }
  if (!this.flushed) {
    // 第一次 flush：把所有累积的 entry 一次性写入
    for (const e of this.fileEntries) {
      appendFileSync(this.sessionFile, JSON.stringify(e) + "\n");
    }
    this.flushed = true;
  } else {
    appendFileSync(this.sessionFile, JSON.stringify(entry) + "\n");
  }
}
```

**为什么？** 避免创建大量只有 user message 的"废弃"session 文件（比如用户输入了一句话但 LLM 还没回复就退出了）。

### 2.3 版本迁移

三代版本，自动迁移：

- **v1 → v2**：给每个 entry 添加 `id`/`parentId`，将 `firstKeptEntryIndex` (数字索引) 转为 `firstKeptEntryId` (字符串 ID)
- **v2 → v3**：将 `hookMessage` 角色重命名为 `custom`

迁移**原地修改内存数据**并**重写整个文件**（`_rewriteFile()`）。

### 2.4 Fork vs Tree Navigate

| 维度 | Fork (`/fork`) | Tree Navigate (`/tree`) |
|------|---------------|------------------------|
| 文件操作 | 创建**新的** session 文件 | 同文件内移动 leaf 指针 |
| 历史保留 | 新文件只包含选中路径 | 完整历史保留在原文件 |
| 摘要 | 无 | 可选（生成 BranchSummaryEntry） |
| parentSession | 记录在新文件 header 中 | 不涉及 |
| 实现方式 | `createBranchedSession()` | `branch()` / `branchWithSummary()` |

## 三、Compaction（上下文压缩）

### 3.1 触发条件

两种触发方式：

**自动触发（agent_end 事件后）**：

```typescript
// Case 1: Overflow - LLM 返回上下文溢出错误
if (isContextOverflow(assistantMessage, contextWindow)) {
  // 移除错误消息，compact，自动重试
}

// Case 2: Threshold - 距离上限较近
if (contextTokens > contextWindow - reserveTokens) {
  // compact，不自动重试（用户手动继续）
}
```

**手动触发**：`/compact [instructions]`

### 3.2 核心算法：prepareCompaction()

```
输入: pathEntries (当前分支的所有 entry), settings

1. 找到上一次 compaction 的位置 (prevCompactionIndex)
   - 如果有前序 compaction，我们只需要处理它之后的内容
   
2. 估算当前上下文 token 数 (estimateContextTokens)
   - 优先使用最后一次 assistant message 的 usage 数据
   - 对 usage 之后的 trailing 消息，用 chars/4 启发式估算

3. 找切割点 (findCutPoint)
   - 从最新消息往前累加 token，直到 >= keepRecentTokens (默认 20k)
   - 只能在有效位置切割（user/assistant/custom/bash，不能在 toolResult）

4. 区分"完整 turn 切割"和"split turn 切割"
   - 如果切在 user message 上 → 完整 turn 切割
   - 如果切在 assistant message 上 → split turn（一个 turn 太大装不下）
   
5. 提取:
   - messagesToSummarize: 要被摘要掉的消息
   - turnPrefixMessages: (仅 split turn) 被切断 turn 的前半段
   - previousSummary: 前序 compaction 的摘要（用于增量更新）
   - fileOps: 从 tool call 中提取的文件操作记录
```

### 3.3 摘要生成：compact()

```
输入: preparation, model, apiKey

1. 如果是 split turn:
   - 并行生成两个摘要:
     a. 历史摘要 (generateSummary) - 使用增量更新 prompt
     b. 轮次前缀摘要 (generateTurnPrefixSummary) - 为保留的后缀提供上下文
   - 合并为一个 summary

2. 如果是完整 turn 切割:
   - 只生成历史摘要

3. 追加文件操作列表到 summary 末尾
   - <read-files> 和 <modified-files> XML 标签

4. 返回 CompactionResult: { summary, firstKeptEntryId, tokensBefore, details }
```

### 3.4 增量摘要 vs 全量摘要

**关键优化**：如果已经有前序 compaction 的摘要 (`previousSummary`)，使用 **UPDATE_SUMMARIZATION_PROMPT** 进行增量更新，而非重新摘要所有内容。

```
全量 prompt: "Summarize this conversation..."
增量 prompt: "Update the existing summary with new information. 
              PRESERVE all existing info, ADD new progress..."
```

这显著降低了多次 compaction 的累积信息损失。

### 3.5 文件操作追踪

Compaction 不仅摘要对话内容，还**累积追踪文件操作**：

```typescript
// 从 tool call 中提取
extractFileOpsFromMessage(message, fileOps);
// 追踪 read/write/edit 操作

// 从前序 compaction 的 details 中继承
if (prevCompaction.details) {
  details.readFiles → fileOps.read
  details.modifiedFiles → fileOps.edited
}

// 最后计算:
// readOnly = read 但没被 modified 的文件
// modified = written + edited 的文件
```

这些信息附加在 summary 末尾，帮助 LLM 知道哪些文件已经处理过。

### 3.6 Overflow 恢复机制

当 LLM 返回上下文溢出错误时：

```
1. 移除错误的 assistant message（从 agent 状态中）
2. 运行 auto compaction
3. 等 100ms 后自动 agent.continue()（重试）
4. 设置 _overflowRecoveryAttempted = true
5. 如果再次 overflow → 放弃（防止无限循环）
```

同模型检查：只在**同一模型**的 overflow 触发恢复。如果用户刚从小上下文模型切到大上下文模型，旧模型的错误不应触发新模型的 compaction。

## 四、Branch Summarization（分支摘要）

### 4.1 触发场景

用户通过 `/tree` 导航到不同分支时，可以选择生成离开分支的摘要。

### 4.2 算法：collectEntriesForBranchSummary()

```
输入: session (只读), oldLeafId, targetId

1. 获取 oldPath (oldLeafId 到 root 的路径)
2. 获取 targetPath (targetId 到 root 的路径)  
3. 从 targetPath 倒序找最深的公共祖先 (commonAncestor)
4. 从 oldLeafId 沿 parentId 回溯到 commonAncestor，收集所有 entry
5. 逆序为时间顺序
```

```
         ┌─ B ─ C ─ D (old leaf)
    A ───┤
         └─ E ─ F (target)

commonAncestor = A
entriesToSummarize = [B, C, D]
```

### 4.3 Token 预算管理

`prepareBranchEntries()` 从**最新到最旧**添加消息，直到达到 token 预算：

```
tokenBudget = contextWindow - reserveTokens (默认 16384)

两遍扫描:
1. 正向遍历：收集所有 entry 的文件操作（即使超出预算也要追踪）
2. 逆向遍历：添加消息直到超出预算
   - 对 compaction 和 branch_summary entry，如果已用 < 90% 预算，强制保留
```

### 4.4 摘要在树中的位置

摘要**不是**附加在旧分支末尾，而是**附加在导航目标位置**：

```
导航前: A → B → C → D (old leaf)
                ↑
              target = C

导航后: A → B → C ──┬── D (旧分支继续存在)
                     └── [branch_summary] (新 leaf)
```

用户后续追加的消息会成为 `branch_summary` 的子节点，而不是 target 的子节点。

## 五、扩展机制

Compaction 和 Branch Summarization 都暴露了扩展点：

| 事件 | 用途 | 可以做什么 |
|------|------|-----------|
| `session_before_compact` | compact 前 | 取消 / 提供自定义摘要 |
| `session_compact` | compact 后 | 读取结果，更新扩展状态 |
| `session_before_tree` | tree 导航前 | 取消 / 自定义摘要 / 修改 instructions |
| `session_tree` | tree 导航后 | 读取结果，更新扩展状态 |
| `session_before_fork` | fork 前 | 取消 / 跳过对话恢复 |
| `session_fork` | fork 后 | 在新 session 中初始化扩展状态 |

扩展可以**完全替换**默认的 compaction 逻辑，例如使用不同的模型、不同的摘要格式、或加入特定领域的结构化信息。

## 六、核心设计要点总结

1. **Append-only JSONL Tree**：结合了 append-only 日志的可靠性和树结构的灵活性
2. **延迟写入**：等到有 assistant 回复才 persist，避免废弃文件
3. **增量摘要**：compaction 时利用前序 summary 进行增量更新，减少信息损失
4. **Split Turn 处理**：当单个 turn 超大时，并行生成历史摘要和轮次前缀摘要
5. **文件操作累积追踪**：跨 compaction/branch summary 累积 read/modified 文件列表
6. **Overflow 自愈**：上下文溢出时自动 compact + retry，一次机会防止循环
7. **同模型检查**：只在同模型 overflow 时触发恢复，跨模型切换不误触发
8. **扩展友好**：所有关键节点都暴露 before/after 事件，支持取消和自定义