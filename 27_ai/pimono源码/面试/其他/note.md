## 已覆盖 vs 未覆盖

| 已有                                                                             | 未覆盖（重要度 ★ 数）         |
| -------------------------------------------------------------------------------- | ----------------------------- |
| tool, plan, subAgent, memory, session, skill, systemPrompt, checkpoint, 模型切换 | **Compaction** ★★★★★          |
|                                                                                  | **Extensions 插件系统** ★★★★★ |
|                                                                                  | **Agent Loop 核心循环** ★★★★  |
|                                                                                  | **Event Bus 事件驱动** ★★★★   |
|                                                                                  | **Streaming 流式输出** ★★★★   |
|                                                                                  | **RPC 协议 & SDK** ★★★★       |
|                                                                                  | **Auto-Retry & 容错** ★★★     |
|                                                                                  | **Thinking Levels** ★★★       |

---

## 一、Compaction（上下文压缩）★★★★★

**这是 Agent 最核心的工程挑战之一，面试必问。**

### 整体流程

```
长对话 token 超窗口
         │
         ▼
┌─────────────────────────┐
│  自动检测 Context Overflow │
│  (agent-session 触发)     │
└────────┬────────────────┘
         ▼
┌─────────────────────────┐
│  Compaction Summary      │  ← 用 LLM 对历史对话做摘要
│  - 保留文件操作记录       │  ← 哪些文件被读/写/创建
│  - split turn 处理       │  ← 超大 turn 切片压缩
│  - turn prefix 保留      │  ← 保留关键上下文前缀
└────────┬────────────────┘
         ▼
┌─────────────────────────┐
│  Branch Summarization    │  ← /tree 分支导航时的上下文保留
│  - 从其他分支跳过来时     │
│  - 需要摘要之前分支做了啥 │
│  - token budget 16384    │
└─────────────────────────┘
```

### 关键源码解析

**1. 触发判定** — `compaction.ts: shouldCompact()`

只需一行：当前 token 数是否超过 `contextWindow - reserveTokens`

```typescript
// compaction.ts
export function shouldCompact(
  contextTokens: number,
  contextWindow: number,
  settings: CompactionSettings
): boolean {
  if (!settings.enabled) return false
  return contextTokens > contextWindow - settings.reserveTokens
}

export const DEFAULT_COMPACTION_SETTINGS: CompactionSettings = {
  enabled: true,
  reserveTokens: 16384, // 留给 system prompt + 新 turn 的空间
  keepRecentTokens: 20000 // 保留最近 20k token 不压缩
}
```

**2. 两种自动触发场景** — `agent-session.ts: _checkCompaction()`

```typescript
// agent-session.ts
private async _checkCompaction(assistantMessage: AssistantMessage): Promise<void> {
  // Case 1: Overflow — LLM 返回了 context overflow 错误
  //   → 删除错误消息 → compact → 自动重试当前请求
  if (sameModel && !errorIsFromBeforeCompaction && isContextOverflow(assistantMessage, contextWindow)) {
    if (this._overflowRecoveryAttempted) {
      // 防止无限循环：只允许一次 compact-and-retry
      return;
    }
    this._overflowRecoveryAttempted = true;
    this.agent.replaceMessages(messages.slice(0, -1)); // 删除错误消息
    await this._runAutoCompaction('overflow', true);    // compact + willRetry=true
    return;
  }

  // Case 2: Threshold — turn 成功但 context 快满了
  //   → compact → 不自动重试（用户手动继续）
  const contextTokens = calculateContextTokens(assistantMessage.usage);
  if (shouldCompact(contextTokens, contextWindow, settings)) {
    await this._runAutoCompaction('threshold', false);
  }
}
```

**3. 切割点算法** — `compaction.ts: findCutPoint()`

核心思想：**从最新消息往回走**，累积 token 直到超过 `keepRecentTokens`，在该处切割。

```typescript
// compaction.ts
export function findCutPoint(entries, startIndex, endIndex, keepRecentTokens): CutPointResult {
  // 从后往前走，累计 token
  for (let i = endIndex - 1; i >= startIndex; i--) {
    const messageTokens = estimateTokens(entry.message)
    accumulatedTokens += messageTokens

    if (accumulatedTokens >= keepRecentTokens) {
      // 在最近的合法切割点切
      // 合法切割点 = user/assistant/custom/bashExecution 消息
      // 永远不在 toolResult 处切（它必须跟着 toolCall）
      break
    }
  }

  return {
    firstKeptEntryIndex: cutIndex, // 保留的起始位置
    turnStartIndex, // 如果切在 turn 中间，turn 的起始位置
    isSplitTurn: !isUserMessage // 是否是 split turn
  }
}
```

**4. 摘要生成** — 结构化 Prompt + 增量更新

```typescript
// compaction.ts — 摘要格式（不是简单截断，是结构化摘要）
const SUMMARIZATION_PROMPT = `Create a structured context checkpoint summary:
## Goal
## Constraints & Preferences
## Progress (Done / In Progress / Blocked)
## Key Decisions
## Next Steps
## Critical Context`

// 增量更新：有 previousSummary 时用 UPDATE_SUMMARIZATION_PROMPT
// "PRESERVE all existing information, ADD new progress"
const summary = previousSummary
  ? await generateSummary(
      messages,
      model,
      reserveTokens,
      apiKey,
      signal,
      customInstructions,
      previousSummary
    )
  : await generateSummary(messages, model, reserveTokens, apiKey, signal, customInstructions)
```

**5. 文件操作追踪** — 跨 compaction 不丢失

```typescript
// compaction.ts — CompactionDetails
export interface CompactionDetails {
  readFiles: string[] // 读过的文件列表
  modifiedFiles: string[] // 修改过的文件列表
}

// 从 tool calls 中提取 + 从上一次 compaction 的 details 中继承
// 确保 "哪些文件被改过" 永远不丢
summary += formatFileOperations(readFiles, modifiedFiles)
```

**6. Split Turn 处理** — 超大 turn 的特殊逻辑

当一个 turn（user→assistant→tool→assistant→tool...）太大时，切割点可能落在 turn 中间。
此时需要分别对 prefix 和 history 做摘要，然后合并：

```typescript
// compaction.ts: compact()
if (isSplitTurn && turnPrefixMessages.length > 0) {
  // 并行生成两个摘要
  const [historyResult, turnPrefixResult] = await Promise.all([
    generateSummary(messagesToSummarize, ...),  // 历史摘要
    generateTurnPrefixSummary(turnPrefixMessages, ...), // turn 前缀摘要
  ]);
  summary = `${historyResult}\n\n---\n\n**Turn Context (split turn):**\n\n${turnPrefixResult}`;
}
```

**7. Branch Summarization** — `/tree` 导航时

从一个分支跳到另一个分支时，需要摘要离开的分支做了什么：

```typescript
// branch-summarization.ts
export function collectEntriesForBranchSummary(session, oldLeafId, targetId) {
  // 找到两个分支的共同祖先
  const oldPath = new Set(session.getBranch(oldLeafId).map(e => e.id))
  const targetPath = session.getBranch(targetId)
  // 从后往前找最深公共祖先
  for (let i = targetPath.length - 1; i >= 0; i--) {
    if (oldPath.has(targetPath[i].id)) {
      commonAncestorId = targetPath[i].id
      break
    }
  }
  // 收集从 oldLeaf → commonAncestor 的所有 entries
  // 再做 summarization
}
```

### 面试回答模板

> "Compaction 解决的核心问题是：长对话超过 context window 后 Agent 就死了。
> pi-mono 的方案是**两级压缩**：
>
> 1. **自动触发**：监测每次 assistant 回复的 usage，当 token 超过 `contextWindow - reserveTokens` 时触发
> 2. **结构化摘要**：不是简单截断，而是用 LLM 生成 Goal/Progress/Key Decisions/Next Steps 格式的摘要
> 3. **文件操作追踪**：跨 compaction 传递 readFiles/modifiedFiles，避免'忘记改过什么文件'
> 4. **Split Turn**：超大 turn 切在中间时，对 prefix 和 history 分别摘要再合并
> 5. **Overflow Recovery**：如果 LLM 直接返回 overflow 错误，先删错误消息再 compact 再自动重试（只允许一次）
> 6. **Branch Summarization**：树形会话跳转时，找公共祖先，摘要离开的分支"

---

## 二、Extensions 插件系统 ★★★★★

### 架构概览

```
Extension = TypeScript 模块 (index.ts)
  │
  ├── 生命周期钩子 (event handlers)
  ├── 自定义工具 (registerTool)
  ├── 自定义命令 (registerCommand)
  ├── 快捷键 (registerShortcut)
  ├── 消息渲染器 (setMessageRenderer)
  └── Provider 注册 (registerProvider)
```

### 关键源码解析

**1. Extension 类型定义** — `types.ts`

```typescript
// extensions/types.ts — 完整的生命周期事件
// Session 事件
SessionStartEvent → SessionBeforeSwitchEvent → SessionSwitchEvent
SessionBeforeForkEvent → SessionForkEvent
SessionBeforeCompactEvent → SessionCompactEvent
SessionBeforeTreeEvent → SessionTreeEvent
SessionShutdownEvent

// Agent 事件
BeforeAgentStartEvent → AgentStartEvent → TurnStartEvent
  → ToolExecutionStartEvent → ToolExecutionUpdateEvent → ToolExecutionEndEvent
  → MessageStartEvent → MessageUpdateEvent → MessageEndEvent
  → TurnEndEvent → AgentEndEvent

// 拦截事件（可修改/阻止）
ToolCallEvent   — 返回 { block: true } 可阻止工具调用
ToolResultEvent — 可修改工具返回结果
InputEvent      — 可转换/拦截用户输入
ContextEvent    — 可修改发给 LLM 的消息列表
```

**2. ExtensionRunner 核心** — `runner.ts`

Runner 是扩展的执行引擎，管理所有 Extension 的生命周期：

```typescript
// runner.ts
export class ExtensionRunner {
  private extensions: Extension[]

  // 通用事件发射 — 依次调用所有 extension 的 handler
  async emit<TEvent>(event: TEvent): Promise<Result> {
    for (const ext of this.extensions) {
      const handlers = ext.handlers.get(event.type)
      for (const handler of handlers) {
        try {
          const result = await handler(event, ctx)
          // session_before_* 事件：如果返回 { cancel: true } 则短路
          if (this.isSessionBeforeEvent(event) && result?.cancel) {
            return result
          }
        } catch (err) {
          // 错误隔离：一个 extension 崩溃不影响其他
          this.emitError({ extensionPath: ext.path, event: event.type, error: message })
        }
      }
    }
  }

  // tool_call 拦截 — 可以 block 工具调用
  async emitToolCall(event: ToolCallEvent): Promise<ToolCallEventResult | undefined> {
    for (const ext of this.extensions) {
      const result = await handler(event, ctx)
      if (result?.block) return result // 阻止执行
    }
  }

  // tool_result 拦截 — 可以修改工具返回值（链式调用）
  async emitToolResult(event: ToolResultEvent): Promise<ToolResultEventResult | undefined> {
    let currentEvent = { ...event }
    for (const handler of handlers) {
      const result = await handler(currentEvent, ctx)
      if (result?.content) currentEvent.content = result.content // 链式修改
    }
  }

  // context 事件 — 修改发给 LLM 的完整消息列表
  async emitContext(messages: AgentMessage[]): Promise<AgentMessage[]> {
    let currentMessages = structuredClone(messages) // 深拷贝防止意外修改
    for (const handler of handlers) {
      const result = await handler({ type: 'context', messages: currentMessages }, ctx)
      if (result?.messages) currentMessages = result.messages
    }
    return currentMessages
  }
}
```

**3. ToolDefinition — 自定义工具注册**

```typescript
// extensions/types.ts
export interface ToolDefinition<TParams, TDetails> {
  name: string
  label: string // UI 显示名
  description: string // 给 LLM 看的描述
  promptSnippet?: string // system prompt 中的工具说明（一行）
  promptGuidelines?: string[] // 追加到 system prompt Guidelines 区域
  parameters: TParams // TypeBox Schema（JSON Schema）

  execute(toolCallId, params, signal, onUpdate, ctx): Promise<AgentToolResult>

  renderCall?(args, theme): Component // 自定义工具调用渲染
  renderResult?(result, options, theme): Component // 自定义结果渲染
}
```

**4. ExtensionContext — 传给 handler 的上下文**

```typescript
// extensions/types.ts
export interface ExtensionContext {
  ui: ExtensionUIContext // UI 方法（select, confirm, input, notify...）
  hasUI: boolean // print/RPC 模式下为 false
  cwd: string
  sessionManager: ReadonlySessionManager // 只读！
  modelRegistry: ModelRegistry
  model: Model<any> | undefined
  isIdle(): boolean // agent 是否空闲
  abort(): void // 中断当前操作
  compact(options?): void // 触发 compaction
  getSystemPrompt(): string // 获取当前 system prompt
}

// 扩展的命令处理上下文（比普通 context 多了会话控制能力）
export interface ExtensionCommandContext extends ExtensionContext {
  waitForIdle(): Promise<void>
  newSession(options?): Promise<{ cancelled: boolean }>
  fork(entryId): Promise<{ cancelled: boolean }>
  navigateTree(targetId, options?): Promise<{ cancelled: boolean }>
  switchSession(sessionPath): Promise<{ cancelled: boolean }>
  reload(): Promise<void>
}
```

**5. 安全机制 — 快捷键冲突检测**

```typescript
// runner.ts — 保留操作不允许 extension 覆盖
const RESERVED_ACTIONS: ReadonlyArray<KeyAction> = [
  "interrupt", "clear", "exit", "suspend",
  "cycleThinkingLevel", "cycleModelForward", "submit", "copy", ...
];

getShortcuts(effectiveKeybindings) {
  for (const [key, shortcut] of ext.shortcuts) {
    const builtIn = builtinKeybindings[normalizedKey];
    if (builtIn?.restrictOverride === true) {
      addDiagnostic(`conflicts with built-in shortcut. Skipping.`);
      continue; // 不允许覆盖
    }
  }
}
```

### 面试回答模板

> "pi-mono 的 Extension 系统是代码级的运行时插件。核心设计：
>
> 1. **生命周期钩子**：覆盖 session/agent/turn/tool 全生命周期
> 2. **拦截链**：tool_call 可以 block（阻止执行），tool_result 可以 chain modify（链式修改），context 可以全量替换
> 3. **错误隔离**：每个 handler 都 try-catch，一个 extension 崩溃不影响其他
> 4. **工具注入**：registerTool 把自定义工具注册到 LLM 的 tool list 里
> 5. **安全**：保留快捷键不可覆盖，命令名冲突自动跳过
> 6. 与 Skill 的区别：Extension 是 TS 代码可以执行网络请求/文件操作/UI，Skill 是声明式的 Prompt + Tool 组合"

---

## 三、Agent Loop 核心循环 ★★★★

### 源码解析 — `packages/agent/src/agent-loop.ts`

这是 Agent 的心脏，整个 ReAct 循环的实现：

```typescript
// agent-loop.ts — 主循环（简化版）
async function runLoop(currentContext, newMessages, config, signal, stream) {
  // 外层循环：处理 follow-up 消息
  while (true) {
    let hasMoreToolCalls = true
    let pendingMessages = await config.getSteeringMessages() // 检查用户打断

    // 内层循环：处理 tool calls + steering
    while (hasMoreToolCalls || pendingMessages.length > 0) {
      stream.push({ type: 'turn_start' })

      // 1. 注入 pending 消息（用户在 agent 工作时输入的）
      for (const message of pendingMessages) {
        currentContext.messages.push(message)
      }

      // 2. 调用 LLM 获取 assistant 回复（流式）
      const message = await streamAssistantResponse(currentContext, config, signal, stream)

      // 3. 如果 LLM 出错或被中断，退出
      if (message.stopReason === 'error' || message.stopReason === 'aborted') {
        stream.push({ type: 'agent_end', messages: newMessages })
        return
      }

      // 4. 执行工具调用
      const toolCalls = message.content.filter(c => c.type === 'toolCall')
      if (toolCalls.length > 0) {
        const { toolResults, steeringMessages } = await executeToolCalls(
          tools,
          message,
          signal,
          stream,
          config.getSteeringMessages
        )
        // 工具结果加入 context
        for (const result of toolResults) currentContext.messages.push(result)
      }

      stream.push({ type: 'turn_end', message, toolResults })

      // 5. 检查 steering 消息（用户打断）
      pendingMessages = await config.getSteeringMessages()
    }

    // 6. Agent 想停了 → 检查 follow-up 队列
    const followUpMessages = await config.getFollowUpMessages()
    if (followUpMessages.length > 0) {
      pendingMessages = followUpMessages
      continue // 继续外层循环
    }
    break // 无更多消息，退出
  }
}
```

**Steering 中断机制** — 用户在 Agent 执行工具时输入新消息

```typescript
// agent-loop.ts: executeToolCalls()
async function executeToolCalls(tools, assistantMessage, signal, stream, getSteeringMessages) {
  for (let index = 0; index < toolCalls.length; index++) {
    // 执行工具...
    const result = await tool.execute(toolCall.id, validatedArgs, signal, onUpdate)

    // 每执行完一个工具，检查是否有 steering 消息
    if (getSteeringMessages) {
      const steering = await getSteeringMessages()
      if (steering.length > 0) {
        // 用户打断了！跳过剩余工具调用
        const remainingCalls = toolCalls.slice(index + 1)
        for (const skipped of remainingCalls) {
          results.push(skipToolCall(skipped)) // "Skipped due to queued user message."
        }
        break
      }
    }
  }
}
```

**LLM 调用边界** — AgentMessage → Message 的转换

```typescript
// agent-loop.ts: streamAssistantResponse()
async function streamAssistantResponse(context, config, signal, stream) {
  // transformContext 钩子（Extension 的 context 事件）
  let messages = context.messages
  if (config.transformContext) {
    messages = await config.transformContext(messages, signal)
  }
  // 转换为 LLM 格式（AgentMessage[] → Message[]）
  const llmMessages = await config.convertToLlm(messages)
  // 动态获取 API Key（支持 OAuth token 过期刷新）
  const resolvedApiKey = await config.getApiKey(config.model.provider)
  // 流式调用 LLM
  const response = await streamFunction(config.model, llmContext, {
    apiKey: resolvedApiKey,
    signal
  })
  // 逐 event 转发
  for await (const event of response) {
    switch (event.type) {
      case 'start':
        stream.push({ type: 'message_start', message: event.partial })
        break
      case 'text_delta':
      case 'thinking_delta':
      case 'toolcall_delta':
        stream.push({ type: 'message_update', assistantMessageEvent: event })
        break
      case 'done':
      case 'error':
        stream.push({ type: 'message_end', message: finalMessage })
        return finalMessage
    }
  }
}
```

### 面试回答模板

> "Agent Loop 是双层 while 循环实现的 ReAct 模式：
>
> - **内层循环**：LLM 调用 → 工具执行 → 加入 context → 再调 LLM，直到没有 toolCall
> - **外层循环**：内层结束后检查 follow-up 队列，有新消息则继续
> - **Steering 机制**：每执行完一个工具就检查用户是否有新输入，有则跳过剩余工具
> - **关键边界**：AgentMessage（内部格式）只在 LLM 调用处转为 Message（LLM 格式），Extension 的 context 事件可以在转换前修改"

---

## 四、Event Bus 事件驱动架构 ★★★★

### 源码解析 — `event-bus.ts`

极简但精巧的实现：

```typescript
// event-bus.ts — 完整源码（只有 33 行）
export interface EventBus {
  emit(channel: string, data: unknown): void
  on(channel: string, handler: (data: unknown) => void): () => void // 返回 unsubscribe
}

export function createEventBus(): EventBusController {
  const emitter = new EventEmitter()
  return {
    emit: (channel, data) => {
      emitter.emit(channel, data)
    },
    on: (channel, handler) => {
      // 关键：用 safeHandler 包装，错误不会传播到 emitter
      const safeHandler = async (data: unknown) => {
        try {
          await handler(data)
        } catch (err) {
          console.error(`Event handler error (${channel}):`, err)
          // 吞掉错误，不影响其他 handler
        }
      }
      emitter.on(channel, safeHandler)
      return () => emitter.off(channel, safeHandler) // 返回清理函数
    },
    clear: () => {
      emitter.removeAllListeners()
    }
  }
}
```

### AgentSession 中的事件系统

AgentSession 本身也是一个事件发射器，定义了更丰富的事件类型：

```typescript
// agent-session.ts — 事件处理流水线
private _handleAgentEvent = (event: AgentEvent): void => {
  // 1. 同步创建 retry promise（防止竞态）
  this._createRetryPromiseForAgentEnd(event);

  // 2. 串行异步队列（保证事件顺序）
  this._agentEventQueue = this._agentEventQueue.then(
    () => this._processAgentEvent(event),
    () => this._processAgentEvent(event)  // 即使前一个失败也继续
  );
};

private async _processAgentEvent(event: AgentEvent): Promise<void> {
  // 1. 先发给 Extensions
  await this._emitExtensionEvent(event);
  // 2. 再通知所有 listener
  this._emit(event);
  // 3. 处理 session 持久化
  if (event.type === 'message_end') {
    this.sessionManager.appendMessage(event.message);
  }
  // 4. 检查 auto-retry 和 auto-compaction
  if (event.type === 'agent_end') {
    if (this._isRetryableError(msg)) { await this._handleRetryableError(msg); }
    await this._checkCompaction(msg);
  }
}
```

**设计要点：用串行 Promise 队列保证事件顺序**。Agent 事件是同步发出的，但 Extension handler 可能是异步的，用 `_agentEventQueue.then()` 链确保事件按序处理。

---

## 五、Streaming 流式输出 ★★★★

### 源码解析

**1. Stream 抽象层** — `packages/ai/src/stream.ts`

```typescript
// stream.ts — 统一的流式接口
export function stream(model, context, options): AssistantMessageEventStream {
  const provider = resolveApiProvider(model.api) // 根据 model.api 找 provider
  return provider.stream(model, context, options)
}

// 支持同步调用（等待流完成）
export async function complete(model, context, options): Promise<AssistantMessage> {
  const s = stream(model, context, options)
  return s.result() // 等待流结束，返回完整消息
}
```

**2. Agent Loop 中的流式处理**

```typescript
// agent-loop.ts: streamAssistantResponse()
const response = await streamFunction(model, llmContext, { apiKey, signal })

// 逐 event 消费流
for await (const event of response) {
  switch (event.type) {
    case 'start':
      // 拿到 partial message，加入 context（作为占位）
      partialMessage = event.partial
      context.messages.push(partialMessage)
      stream.push({ type: 'message_start', message: { ...partialMessage } })
      break

    case 'text_delta':
    case 'thinking_delta':
    case 'toolcall_delta':
      // 更新 partial message（原地替换 context 中的占位）
      partialMessage = event.partial
      context.messages[context.messages.length - 1] = partialMessage
      // 向上层转发 update 事件（UI 逐字渲染）
      stream.push({ type: 'message_update', assistantMessageEvent: event })
      break

    case 'done':
    case 'error':
      // 替换为最终完整消息
      const finalMessage = await response.result()
      context.messages[context.messages.length - 1] = finalMessage
      stream.push({ type: 'message_end', message: finalMessage })
      return finalMessage
  }
}
```

**3. AbortSignal 取消链** — 全链路贯穿

```
用户按 Escape
  → agent.abort()
    → AbortController.abort()
      → signal 传给 stream() → HTTP 请求取消
      → signal 传给 tool.execute() → bash 子进程 kill
      → signal 传给 compact() → 摘要生成取消
```

---

## 六、RPC 协议 & SDK ★★★★

### RPC 协议设计 — `rpc-types.ts`

stdin/stdout JSON Lines 协议，支持关联 ID：

```typescript
// rpc-types.ts — 命令类型（stdin 输入）
export type RpcCommand =
  // 提示
  | { id?: string; type: "prompt"; message: string; images?: ImageContent[] }
  | { id?: string; type: "steer"; message: string }
  | { id?: string; type: "follow_up"; message: string }
  | { id?: string; type: "abort" }
  // 状态
  | { id?: string; type: "get_state" }
  // 模型
  | { id?: string; type: "set_model"; provider: string; modelId: string }
  | { id?: string; type: "cycle_model" }
  // Compaction
  | { id?: string; type: "compact"; customInstructions?: string }
  // 重试
  | { id?: string; type: "abort_retry" }
  // Bash
  | { id?: string; type: "bash"; command: string }
  // Session
  | { id?: string; type: "fork"; entryId: string }
  | { id?: string; type: "switch_session"; sessionPath: string }
  | { id?: string; type: "get_messages" }
  // ...

// 响应类型（stdout 输出）— 与命令 ID 关联
{ id: "cmd-1", type: "response", command: "prompt", success: true }
{ id: "cmd-1", type: "response", command: "get_state", success: true, data: RpcSessionState }
// 错误也带 ID
{ id: "cmd-1", type: "response", command: "set_model", success: false, error: "Model not found" }
```

### SDK — `sdk.ts`: `createAgentSession()`

```typescript
// sdk.ts — 工厂函数（核心 API）
export async function createAgentSession(options: CreateAgentSessionOptions = {}): Promise<CreateAgentSessionResult> {
  // 1. 初始化基础设施
  const authStorage = options.authStorage ?? AuthStorage.create(authPath);
  const modelRegistry = new ModelRegistry(authStorage, modelsPath);
  const settingsManager = SettingsManager.create(cwd, agentDir);
  const sessionManager = SessionManager.create(cwd);
  const resourceLoader = new DefaultResourceLoader({ cwd, agentDir, settingsManager });
  await resourceLoader.reload();

  // 2. Model 解析（支持 session 恢复 → settings default → provider default 三级 fallback）
  let model = options.model;
  if (!model && hasExistingSession && existingSession.model) {
    model = modelRegistry.find(existingSession.model.provider, existingSession.model.modelId);
  }
  if (!model) {
    const result = await findInitialModel({ ... });
    model = result.model;
  }

  // 3. Thinking Level 解析（session → settings → default medium）
  let thinkingLevel = options.thinkingLevel ?? settingsManager.getDefaultThinkingLevel() ?? 'medium';
  if (!model?.reasoning) thinkingLevel = 'off'; // 非推理模型强制关闭

  // 4. 创建 Agent 核心
  const agent = new Agent({
    initialState: { systemPrompt: "", model, thinkingLevel, tools: [] },
    convertToLlm: convertToLlmWithBlockImages,
    transformContext: async (messages) => runner.emitContext(messages), // Extension 钩子
    getApiKey: async (provider) => modelRegistry.getApiKeyForProvider(provider),
  });

  // 5. 如果有 existing session，恢复消息
  if (hasExistingSession) {
    agent.replaceMessages(existingSession.messages);
  }

  return { session: new AgentSession({ agent, sessionManager, ... }) };
}
```

### 三大运行模式共享 Core

```
                    ┌── Interactive Mode (TUI)
AgentSession ──────├── RPC Mode (headless JSON)
  (core)           └── Print Mode (单次输出)

三个模式共享同一个 AgentSession，只是 I/O 层不同
```

---

## 七、Auto-Retry 容错机制 ★★★

### 源码解析 — `agent-session.ts`

**1. 错误分类** — 只重试瞬态错误

```typescript
// agent-session.ts: _isRetryableError()
private _isRetryableError(message: AssistantMessage): boolean {
  if (message.stopReason !== 'error') return false;
  // Context overflow 走 compaction，不走 retry
  if (isContextOverflow(message, contextWindow)) return false;
  // 正则匹配瞬态错误
  return /overloaded|rate.?limit|too many requests|429|500|502|503|504|
    service.?unavailable|server error|connection.?error|
    fetch failed|terminated|retry delay/i.test(message.errorMessage);
}
```

**2. 指数退避** — 可中断的 sleep

```typescript
// agent-session.ts: _handleRetryableError()
private async _handleRetryableError(message: AssistantMessage): Promise<boolean> {
  const settings = this.settingsManager.getRetrySettings();
  if (!settings.enabled) return false;

  this._retryAttempt++;
  if (this._retryAttempt > settings.maxRetries) {
    this._emit({ type: 'auto_retry_end', success: false, finalError: message.errorMessage });
    this._retryAttempt = 0;
    return false;
  }

  // 指数退避: baseDelay * 2^(attempt-1)
  const delayMs = settings.baseDelayMs * 2 ** (this._retryAttempt - 1);
  this._emit({ type: 'auto_retry_start', attempt: this._retryAttempt, delayMs, errorMessage });

  // 删除错误消息，保持在 session 历史里但不在 context 里
  this.agent.replaceMessages(messages.slice(0, -1));

  // 可中断的等待（用户可以 abortRetry）
  this._retryAbortController = new AbortController();
  try { await sleep(delayMs, this._retryAbortController.signal); } catch { return false; }

  // 用 setTimeout 断开事件处理链，避免递归
  setTimeout(() => { this.agent.continue(); }, 0);
  return true;
}
```

**3. Retry + Compaction 的配合**

```typescript
// agent-session.ts: _processAgentEvent()
if (event.type === 'agent_end' && this._lastAssistantMessage) {
  const msg = this._lastAssistantMessage
  // 先检查 retry（429/500）
  if (this._isRetryableError(msg)) {
    const didRetry = await this._handleRetryableError(msg)
    if (didRetry) return // retry 启动了，不走 compaction
  }
  // 再检查 compaction（overflow / threshold）
  await this._checkCompaction(msg)
}
```

**关键：成功响应时立即重置重试计数器**

```typescript
// agent-session.ts
if (assistantMsg.stopReason !== 'error' && this._retryAttempt > 0) {
  this._emit({ type: 'auto_retry_end', success: true, attempt: this._retryAttempt })
  this._retryAttempt = 0 // 重置！防止跨 turn 累积
}
```

---

## 八、Thinking Levels（推理预算控制）★★★

```
off → minimal → low → medium(default) → high → xhigh

每个级别对应不同的 thinking token budget
用户可通过 shift+tab 实时切换
```

```typescript
// sdk.ts — Thinking Level 解析逻辑
let thinkingLevel =
  options.thinkingLevel ?? settingsManager.getDefaultThinkingLevel() ?? DEFAULT_THINKING_LEVEL // 'medium'

// 非推理模型强制关闭
if (!model || !model.reasoning) {
  thinkingLevel = 'off'
}

// Agent 创建时传入 thinkingBudgets 配置（来自 settings.json）
const agent = new Agent({
  thinkingBudgets: settingsManager.getThinkingBudgets()
})
```

---

## 总结：面试必补清单

| 优先级 | 特性                   | 一句话                                                     |
| ------ | ---------------------- | ---------------------------------------------------------- |
| **P0** | **Compaction**         | 两级压缩 + 结构化摘要 + 文件追踪 + overflow 自动恢复       |
| **P0** | **Extensions**         | 生命周期钩子 + 拦截链 + 工具注入 + 错误隔离                |
| **P0** | **Agent Loop**         | 双层 while 循环 + steering 打断 + follow-up 队列           |
| **P1** | **Streaming + 取消链** | for-await-of 流式 + AbortSignal 全链路取消                 |
| **P1** | **Event Bus**          | 33 行代码的发布订阅 + 错误隔离 + 串行 Promise 队列保序     |
| **P1** | **RPC/SDK**            | JSON Lines 协议 + createAgentSession 工厂 + 三模式共享核心 |
| **P2** | **Auto-Retry**         | 正则分类瞬态错误 + 指数退避 + 可中断 sleep + 成功即重置    |
| **P2** | **Thinking Levels**    | 6 级推理深度 + 非推理模型强制 off + settings 可配          |
