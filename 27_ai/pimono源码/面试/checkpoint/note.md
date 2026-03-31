# Pi-mono Checkpoint & 分支实现深度分析

## 一、整体架构

Pi-mono 的 session 系统围绕**树结构**构建，核心设计哲学是：

1. **Append-only JSONL 文件**：所有 session 数据以 JSONL 格式存储，只追加不修改
2. **树结构而非线性链**：每个 entry 有 `id` + `parentId`，形成树而非链表
3. **Leaf 指针**：`leafId` 追踪当前位置，分支/导航只需移动指针，不修改历史
4. **两级上下文管理**：Compaction（压缩）保留近期上下文 + Branch Summary（分支摘要）跨分支保留上下文

```
核心文件：
├── session-manager.ts          # Session 树结构、JSONL 存储、分支操作
├── agent-session.ts            # 上层编排（compaction、tree 导航、fork）
├── compaction/
│   ├── compaction.ts           # 自动/手动 compaction 逻辑
│   ├── branch-summarization.ts # 分支摘要生成
│   └── utils.ts                # 序列化、文件追踪工具
└── messages.ts                 # 消息类型定义（BranchSummaryMessage 等）
```

---

## 二、Session 树结构（核心数据模型）

### 2.1 存储格式

JSONL 文件，每行一个 JSON 对象：

```
{"type":"session","version":3,"id":"uuid","timestamp":"...","cwd":"..."}
{"type":"message","id":"a1b2","parentId":null,"message":{...}}
{"type":"message","id":"c3d4","parentId":"a1b2","message":{...}}
{"type":"compaction","id":"e5f6","parentId":"c3d4","summary":"...","firstKeptEntryId":"a1b2",...}
```

### 2.2 Entry 类型

```typescript
type SessionEntry =
  | SessionMessageEntry       // 用户/助手/工具消息
  | ThinkingLevelChangeEntry  // 思考级别变更
  | ModelChangeEntry          // 模型切换
  | CompactionEntry           // 压缩摘要
  | BranchSummaryEntry        // 分支摘要
  | CustomEntry               // 扩展自定义数据（不参与 LLM 上下文）
  | CustomMessageEntry        // 扩展自定义消息（参与 LLM 上下文）
  | LabelEntry                // 用户书签/标记
  | SessionInfoEntry          // Session 元数据
```

### 2.3 树的构建

```typescript
// SessionManager.getTree()：把扁平 entries 构建为树
getTree(): SessionTreeNode[] {
  // 1. 为每个 entry 创建 node
  // 2. 根据 parentId 建立父子关系
  // 3. 孤儿节点作为根
  // 4. 子节点按 timestamp 排序（最旧在前）
}
```

### 2.4 上下文构建（buildSessionContext）

关键函数 —— 从 leaf 到 root 走一条路径，提取 LLM 需要的消息：

```typescript
function buildSessionContext(entries, leafId, byId): SessionContext {
  // 1. 从 leafId 沿 parentId 向上走到 root，收集 path
  // 2. 沿 path 提取最新的 thinkingLevel、model
  // 3. 找到路径上最后一个 compaction entry
  // 4. 如果有 compaction：
  //    - 先输出 compaction summary（作为 user 消息）
  //    - 再输出 firstKeptEntryId 之后的消息
  // 5. 如果没有 compaction：输出所有消息
  // 6. branch_summary 和 custom_message 也被转为消息
}
```

**核心洞察**：LLM 看到的上下文 = compaction summary + 保留的消息 + 分支摘要。三者都被转为 `user` 角色消息注入。

---

## 三、Compaction（上下文压缩/Checkpoint）

### 3.1 触发条件

```
两种触发方式：
1. 自动触发：contextTokens > contextWindow - reserveTokens（默认 16384）
2. 手动触发：/compact [instructions]

特殊情况：
- Overflow recovery：LLM 返回 context overflow 错误时，自动 compact + retry
- 仅尝试一次 overflow recovery，失败后提示用户
```

### 3.2 Compaction 流程

```
┌──────────────────────────────────────────────────────┐
│ 1. prepareCompaction() - 纯函数，计算压缩方案        │
│    ├─ 找到上一次 compaction 位置                      │
│    ├─ 从 compaction 之后开始估算 token               │
│    ├─ findCutPoint() 反向扫描找到切分点               │
│    │   └─ 保留最近 keepRecentTokens（默认 20000）    │
│    ├─ 分离 messagesToSummarize + 保留消息             │
│    └─ 提取文件操作记录                                │
│                                                      │
│ 2. compact() - 调用 LLM 生成摘要                     │
│    ├─ 普通情况：一次 LLM 调用 generateSummary()      │
│    ├─ Split Turn：并行两次 LLM 调用                  │
│    │   ├─ generateSummary() —— 历史摘要              │
│    │   └─ generateTurnPrefixSummary() —— Turn 前缀   │
│    └─ 附加文件操作列表到摘要末尾                      │
│                                                      │
│ 3. sessionManager.appendCompaction() - 保存           │
│    └─ 追加 CompactionEntry 到 JSONL                  │
│                                                      │
│ 4. agent.replaceMessages() - 刷新 LLM 上下文         │
│    └─ buildSessionContext() 重新走 root→leaf 路径     │
└──────────────────────────────────────────────────────┘
```

### 3.3 Cut Point 算法

```typescript
function findCutPoint(entries, startIndex, endIndex, keepRecentTokens) {
  // 1. 找到所有合法切分点（user/assistant/bash/custom 消息，不能切 toolResult）
  // 2. 从最新消息反向累加 token 估算值（chars/4）
  // 3. 累加 >= keepRecentTokens 时停止
  // 4. 找到距离停止点最近的合法切分点
  // 5. 检测是否为 split turn（切分点不是 user 消息）
}
```

**Split Turn 场景**：当单个 turn 超过 keepRecentTokens 时（例如超长的工具调用链），切分点落在 turn 中间的 assistant 消息上。此时生成两个摘要并合并。

### 3.4 迭代 Compaction

多次 compaction 时使用 **UPDATE_SUMMARIZATION_PROMPT**，将新消息合并到上一次的摘要中，而非重新生成：

```
首次 compaction：SUMMARIZATION_PROMPT（从零生成）
后续 compaction：UPDATE_SUMMARIZATION_PROMPT（在上次摘要基础上更新）
```

### 3.5 CompactionEntry 结构

```typescript
interface CompactionEntry<T = unknown> {
  type: "compaction";
  id: string;                 // 唯一 ID
  parentId: string;           // 树中的父节点
  timestamp: string;
  summary: string;            // LLM 生成的结构化摘要
  firstKeptEntryId: string;   // 保留消息的起点 ← 关键字段
  tokensBefore: number;       // 压缩前的 token 数
  details?: T;                // 默认为 {readFiles, modifiedFiles}
  fromHook?: boolean;         // 是否由扩展提供
}
```

### 3.6 LLM 看到的内容

```
压缩后 LLM 视角：

┌────────┬──────────────────┬─────────────────────────────┐
│ system │ compaction       │ 保留的消息                    │
│ prompt │ summary          │ (从 firstKeptEntryId 起)     │
│        │ (注入为 user msg)│                               │
└────────┴──────────────────┴─────────────────────────────┘
```

### 3.7 摘要格式（Checkpoint 格式）

```markdown
## Goal
[用户目标]

## Constraints & Preferences
- [约束和偏好]

## Progress
### Done
- [x] [完成的任务]
### In Progress
- [ ] [进行中的工作]
### Blocked
- [阻塞项]

## Key Decisions
- **[决策]**: [理由]

## Next Steps
1. [下一步]

## Critical Context
- [关键上下文]

<read-files>
path/to/file1.ts
</read-files>

<modified-files>
path/to/changed.ts
</modified-files>
```

---

## 四、分支（Branching）

### 4.1 两种分支方式

| 操作 | `/fork` | `/tree` |
|------|---------|---------|
| 含义 | 从某个 user 消息分叉 | 在树中自由导航 |
| 文件 | **创建新 session 文件** | **同一文件内移动 leaf** |
| 摘要 | 无 | 可选（用户决定） |
| 场景 | 想从头重来某个方向 | 想在历史分支间切换 |

### 4.2 Fork 实现

```typescript
async fork(entryId: string) {
  // 1. 找到选中的 user 消息
  // 2. 触发 session_before_fork 扩展事件（可取消）
  // 3. 如果选中的是根消息（无 parentId）：
  //    → newSession({ parentSession: previousFile })
  // 4. 否则：
  //    → createBranchedSession(entry.parentId)
  //       提取 root→parent 的路径到新文件
  // 5. 刷新 agent 消息
  // 6. 返回选中消息的文本（预填充到编辑器）
}
```

`createBranchedSession()` 的关键点：
- 从树中提取一条路径写入新的 JSONL 文件
- 新 session header 中记录 `parentSession` 指向原文件
- label 也被迁移

### 4.3 Tree Navigation 实现

```typescript
async navigateTree(targetId, options) {
  // 1. 计算公共祖先（oldLeaf 和 target 两条路径的最深交叉点）
  // 2. 收集待摘要的 entries（从 oldLeaf 回溯到公共祖先）
  // 3. 触发 session_before_tree 扩展事件
  // 4. 若用户选择摘要，调用 generateBranchSummary()
  // 5. 确定新 leaf 位置：
  //    - user 消息 → leaf = parent，文本进编辑器
  //    - 非 user 消息 → leaf = 选中节点
  // 6. 调用 branch() 或 branchWithSummary() 移动 leaf
  // 7. 刷新 agent 消息
}
```

### 4.4 分支摘要生成

```
寻找公共祖先：
         ┌─ B ─ C ─ D (old leaf)
    A ───┤
         └─ E ─ F (target)

公共祖先 = A
待摘要 entries = [B, C, D]

摘要后的树：
         ┌─ B ─ C ─ D
    A ───┤
         └─ E ─ F ─ [branch_summary of B,C,D] (new leaf)
```

collectEntriesForBranchSummary()：
```typescript
function collectEntriesForBranchSummary(session, oldLeafId, targetId) {
  // 1. 获取 oldLeafId 路径上所有节点 ID (Set)
  // 2. 获取 targetId 的路径，反向查找最深的公共祖先
  // 3. 从 oldLeaf 回溯到公共祖先，收集 entries
  // 4. reverse 得到时间顺序
}
```

prepareBranchEntries()：
```typescript
function prepareBranchEntries(entries, tokenBudget) {
  // 第一遍：从所有 entries 收集文件操作（包括嵌套的 branch_summary.details）
  // 第二遍：从最新到最旧添加消息，直到达到 token 预算
  //   ├─ compaction/branch_summary 类型优先保留（即使超预算90%也加入）
  //   └─ 跳过 toolResult（上下文在 assistant 的 toolCall 中）
}
```

### 4.5 BranchSummaryEntry 结构

```typescript
interface BranchSummaryEntry<T = unknown> {
  type: "branch_summary";
  id: string;
  parentId: string;         // 新 leaf 位置
  timestamp: string;
  fromId: string;           // 被放弃的旧 leaf
  summary: string;          // LLM 生成的摘要
  details?: T;              // 默认 {readFiles, modifiedFiles}
  fromHook?: boolean;       // 是否由扩展提供
}
```

### 4.6 SessionManager.branch() / branchWithSummary()

```typescript
// 简单分支：仅移动 leaf 指针
branch(branchFromId: string): void {
  this.leafId = branchFromId;
  // 不修改任何 entry，不删除任何东西
}

// 带摘要的分支：移动 leaf + 追加摘要 entry
branchWithSummary(branchFromId, summary, details, fromHook): string {
  this.leafId = branchFromId;  // 或 null
  // 追加 BranchSummaryEntry 作为 branchFromId 的子节点
  this._appendEntry(entry);
  return entry.id;
}

// 重置 leaf 到空（用于导航到根 user 消息）
resetLeaf(): void {
  this.leafId = null;
  // 下一次 appendXXX 将创建新的 root entry
}
```

---

## 五、Compaction vs Branch Summary 对比

| 维度 | Compaction | Branch Summary |
|------|-----------|---------------|
| 触发 | 自动（token 超限）或 `/compact` | `/tree` 导航时用户选择 |
| 目的 | 压缩历史，释放上下文窗口 | 保留被放弃分支的上下文 |
| Entry 类型 | `CompactionEntry` | `BranchSummaryEntry` |
| 关键字段 | `firstKeptEntryId`（保留起点） | `fromId`（离开的位置） |
| LLM 注入 | `compactionSummary` → user msg | `branchSummary` → user msg |
| 迭代更新 | 支持（UPDATE_SUMMARIZATION_PROMPT） | 不支持（每次独立生成） |
| 文件追踪 | 累积 readFiles + modifiedFiles | 累积 readFiles + modifiedFiles |

---

## 六、扩展机制

两个关键扩展点允许自定义 compaction/branch summary：

```typescript
// 1. session_before_compact：拦截压缩
pi.on("session_before_compact", async (event, ctx) => {
  const { preparation, branchEntries, signal } = event;
  // 取消：return { cancel: true }
  // 自定义摘要：return { compaction: { summary, firstKeptEntryId, tokensBefore, details } }
});

// 2. session_before_tree：拦截树导航
pi.on("session_before_tree", async (event, ctx) => {
  const { preparation, signal } = event;
  // 取消：return { cancel: true }
  // 自定义摘要：return { summary: { summary, details } }
  // 覆盖 instructions/label：return { customInstructions, label }
});
```

---

## 七、面试要点总结

### 核心设计决策

1. **树结构 + append-only** → 永不丢失历史，分支是移动指针而非复制
2. **JSONL 存储** → 追加写入高效，crash-safe，易于调试
3. **Compaction = LLM 生成的 checkpoint** → 不是简单截断，而是结构化摘要
4. **buildSessionContext = root→leaf 路径遍历** → 只发送当前分支的消息给 LLM
5. **Branch Summary 注入位置** → 在新 leaf 处追加，而非在旧分支标记
6. **迭代更新 vs 全量摘要** → 多次 compaction 时合并而非重新生成

### 典型面试问题

**Q：为什么用树结构而不是线性列表？**
A：支持分支（branching），用户可以回到任何历史节点重新开始，不丢失旧分支。导航只需改 leafId 指针。

**Q：Compaction 如何保证不丢失重要上下文？**
A：结构化摘要（Goal/Progress/Decisions/NextSteps），保留近期 keepRecentTokens，累积文件追踪列表，迭代更新合并旧摘要。

**Q：为什么 Branch Summary 和 Compaction Summary 都注入为 user msg？**
A：LLM API 只接受 user/assistant/toolResult，这些自定义类型需要映射。注入为 user 消息确保 LLM 将其视为上下文信息而非自己的输出。

**Q：Split Turn 是什么？如何处理？**
A：当单个 turn（user→多个 assistant+tool）超过 keepRecentTokens 时，切分点落在 turn 中间。此时并行生成两个摘要：历史摘要 + turn 前缀摘要，合并后作为 compaction。

**Q：Fork 和 Tree Navigate 的区别？**
A：Fork 创建新 session 文件（从历史路径提取），tree navigate 在同一文件内移动 leaf 指针。Fork 不生成摘要，tree navigate 可选摘要。

**Q：Overflow recovery 如何工作？**
A：LLM 返回 context overflow 错误 → 移除错误消息 → 自动 compact → retry。仅尝试一次，失败后提示用户。同时检查错误是否来自已切换的模型或已 compact 过的区域，避免重复触发。
