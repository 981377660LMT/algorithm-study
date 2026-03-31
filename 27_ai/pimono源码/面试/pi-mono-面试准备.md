# pi-mono 深度源码分析 & MiniMax Agent 岗位模拟面试

---

## 一、项目总览

pi-mono 是一个 monorepo 结构的 **AI Coding Agent 框架**，核心由三个包组成：

| 包                              | 定位                | 职责                                                                     |
| ------------------------------- | ------------------- | ------------------------------------------------------------------------ |
| `@mariozechner/pi-ai`           | LLM Provider 抽象层 | 统一多厂商 API（OpenAI/Anthropic/Google/MiniMax 等）、流式推理、工具校验 |
| `@mariozechner/pi-agent-core`   | Agent 循环核心      | Agent Loop（ReAct 模式）、消息管理、事件流、Proxy 转发                   |
| `@mariozechner/pi-coding-agent` | 完整 Coding Agent   | CLI/TUI/RPC 三模式、Session 持久化树、扩展系统、Compaction、Skills       |

**架构层级：**

```
┌─────────────────────────────────────────────────────────────┐
│          coding-agent (CLI/TUI/RPC/SDK)                    │
│    InteractiveMode / PrintMode / RPCMode                   │
├─────────────────────────────────────────────────────────────┤
│          AgentSession (核心中间层)                           │
│    Session 持久化 · 扩展系统 · Compaction · Skills          │
├─────────────────────────────────────────────────────────────┤
│          agent-core (Agent Loop)                           │
│    Tool 执行 · Steering/FollowUp · EventStream            │
├─────────────────────────────────────────────────────────────┤
│          pi-ai (LLM 抽象层)                                │
│    Provider 注册 · Stream 协议 · Tool Schema 校验          │
└─────────────────────────────────────────────────────────────┘
```

---

## 二、核心源码深度分析

### 2.1 EventStream — 生产者-消费者异步迭代器

**文件：** `packages/ai/src/utils/event-stream.ts`

```ts
export class EventStream<T, R = T> implements AsyncIterable<T> {
  private queue: T[] = []
  private waiting: ((value: IteratorResult<T>) => void)[] = []
  private done = false
  private finalResultPromise: Promise<R>
  private resolveFinalResult!: (result: R) => void

  constructor(
    private isComplete: (event: T) => boolean,  // 判断流结束的谓词
    private extractResult: (event: T) => R       // 从终止事件提取结果
  ) {
    this.finalResultPromise = new Promise(resolve => {
      this.resolveFinalResult = resolve
    })
  }

  push(event: T): void { ... }     // 生产者推事件
  end(result?: R): void { ... }    // 结束流
  async *[Symbol.asyncIterator]()  // 消费者拉事件
  result(): Promise<R>             // 获取最终结果
}
```

**关键设计点：**

1. **双缓冲机制**：`queue` 缓冲未消费事件，`waiting` 缓冲等待中的消费者。如果有消费者在等（`waiting.shift()`），直接投递；否则入队。
2. **完成语义**：通过 `isComplete` 谓词判断哪个事件表示流结束，并从该事件中 `extractResult` 得到最终值。
3. **`result()` Promise**：允许不迭代事件，直接 await 最终结果（用于 `complete()` 这种只要最终消息的场景）。

**面试考点：** 这是一个自实现的 **异步可迭代协议**（`Symbol.asyncIterator`），本质是 CSP channel 模式。类似 Go 的 channel，但用 Promise resolve 代替阻塞。

---

### 2.2 API Provider 注册表 — 策略模式 + 注册表模式

**文件：** `packages/ai/src/api-registry.ts`

```ts
const apiProviderRegistry = new Map<string, RegisteredApiProvider>()

export function registerApiProvider<TApi, TOptions>(
  provider: ApiProvider<TApi, TOptions>,
  sourceId?: string
): void { ... }

export function getApiProvider(api: Api): ApiProviderInternal | undefined {
  return apiProviderRegistry.get(api)?.provider
}
```

**设计要点：**

- 每个 Provider（Anthropic/OpenAI/Google 等）实现 `stream()` 和 `streamSimple()` 两个接口
- `sourceId` 用于扩展注册的 Provider，可通过 `unregisterApiProviders(sourceId)` 卸载
- `register-builtins.ts` 在模块加载时自动注册所有内置 Provider（side-effect import）

**面试考点：** 这是经典的 **策略模式（Strategy Pattern）** + **服务定位器（Service Locator）**。新增 Provider 无需修改核心代码，只需 `registerApiProvider()`。

---

### 2.3 Agent Loop — ReAct 模式的核心实现

**文件：** `packages/agent/src/agent-loop.ts`

这是整个框架最核心的文件。Agent Loop 实现了 **ReAct（Reasoning + Acting）** 范式：

```
用户消息 → LLM 推理 → [有工具调用?]
                           ├─ 是 → 执行工具 → 检查 Steering → 继续循环
                           └─ 否 → 检查 FollowUp → [有?] → 继续 : 结束
```

**核心函数 `runLoop()` 的双层循环设计：**

```ts
async function runLoop(currentContext, newMessages, config, signal, stream, streamFn) {
  let pendingMessages = (await config.getSteeringMessages?.()) || []

  // 外层循环：处理 FollowUp 消息
  while (true) {
    let hasMoreToolCalls = true

    // 内层循环：处理工具调用 + Steering 消息
    while (hasMoreToolCalls || pendingMessages.length > 0) {
      // 1. 注入 pending 消息
      if (pendingMessages.length > 0) { ... }

      // 2. 流式获取 LLM 响应
      const message = await streamAssistantResponse(...)

      // 3. 错误处理（提前退出）
      if (message.stopReason === "error" || message.stopReason === "aborted") { return }

      // 4. 提取并执行工具调用
      const toolCalls = message.content.filter(c => c.type === "toolCall")
      if (toolCalls.length > 0) {
        const { toolResults, steeringMessages } = await executeToolCalls(...)
        // 工具结果放入上下文
      }

      // 5. 检查 Steering（用户中途打断）
      pendingMessages = (await config.getSteeringMessages?.()) || []
    }

    // 外层检查：FollowUp 消息
    const followUpMessages = (await config.getFollowUpMessages?.()) || []
    if (followUpMessages.length > 0) {
      pendingMessages = followUpMessages
      continue  // 重新进入内层循环
    }
    break
  }
}
```

**三种消息流机制：**

| 机制         | 时机               | 效果                                                  |
| ------------ | ------------------ | ----------------------------------------------------- |
| **Steering** | Agent 执行工具期间 | 中断剩余工具调用，立即注入用户消息，下一轮 LLM 会看到 |
| **FollowUp** | Agent 自然结束后   | 追加消息，Agent 继续执行（不中断）                    |
| **Prompt**   | Agent 空闲时       | 正常启动新一轮                                        |

**Steering 的中断机制：**

```ts
async function executeToolCalls(tools, assistantMessage, signal, stream, getSteeringMessages) {
  for (let index = 0; index < toolCalls.length; index++) {
    // 执行当前工具...
    const result = await tool.execute(...)

    // 每个工具执行完后检查 Steering
    if (getSteeringMessages) {
      const steering = await getSteeringMessages()
      if (steering.length > 0) {
        // 跳过剩余工具，标记为 "Skipped due to queued user message"
        for (const skipped of toolCalls.slice(index + 1)) {
          results.push(skipToolCall(skipped, stream))
        }
        break
      }
    }
  }
}
```

**面试考点：**

1. 为什么是双层循环？外层处理 FollowUp，内层处理 ToolCalls + Steering。
2. Steering vs FollowUp 的语义区别：Steering 是"紧急打断"，FollowUp 是"任务追加"。
3. 工具调用是**串行**的（非并行），每个工具执行后都检查 Steering，实现了**可中断的协作式调度**。

---

### 2.4 streamAssistantResponse — LLM 调用边界的消息转换

```ts
async function streamAssistantResponse(context, config, signal, stream, streamFn) {
  // 1. AgentMessage[] → AgentMessage[]（可选 transformContext）
  let messages = context.messages
  if (config.transformContext) {
    messages = await config.transformContext(messages, signal)
  }

  // 2. AgentMessage[] → Message[]（convertToLlm）
  const llmMessages = await config.convertToLlm(messages)

  // 3. 构建 LLM Context（与具体 Agent 消息类型解耦）
  const llmContext: Context = {
    systemPrompt: context.systemPrompt,
    messages: llmMessages,
    tools: context.tools
  }

  // 4. 动态 API Key 获取（支持 OAuth 等短期 Token）
  const resolvedApiKey = config.getApiKey
    ? await config.getApiKey(config.model.provider)
    : undefined

  // 5. 流式调用 LLM
  const response = await streamFunction(config.model, llmContext, {
    ...config,
    apiKey: resolvedApiKey,
    signal
  })

  // 6. 事件驱动的消息组装
  for await (const event of response) {
    switch (event.type) {
      case 'start':
        // 将 partial 消息立即加入 context（后续事件更新同一位置）
        context.messages.push(partialMessage)
        break
      case 'text_delta':
      case 'toolcall_delta':
        // 原地更新 context 最后一条消息
        context.messages[context.messages.length - 1] = partialMessage
        break
      case 'done':
      case 'error':
        // 用最终消息替换 partial
        return finalMessage
    }
  }
}
```

**关键设计：**

- **两层消息转换**：`transformContext`（Agent 级别，如裁剪）→ `convertToLlm`（过滤非 LLM 消息）
- **原地更新**：流式过程中通过直接修改 `context.messages` 数组的最后一个元素实现"边流边更新"
- **动态 API Key**：支持 OAuth Token 过期刷新（长时间工具执行后 token 可能过期）

---

### 2.5 Agent 类 — 状态管理 + 观察者模式

**文件：** `packages/agent/src/agent.ts`

Agent 类是 agent-loop 的高级封装，提供：

```ts
export class Agent {
  private _state: AgentState // 核心状态
  private listeners: Set<fn> // 观察者列表
  private abortController // 取消控制
  private steeringQueue // Steering 消息队列
  private followUpQueue // FollowUp 消息队列

  // 消息投递
  async prompt(input) // 启动新对话轮次
  steer(m) // 队列 Steering 消息
  followUp(m) // 队列 FollowUp 消息
  async continue() // 从当前上下文继续
  abort() // 取消当前执行

  // 事件订阅（观察者模式）
  subscribe(fn): () => void // 返回取消订阅函数

  // 等待空闲
  waitForIdle(): Promise<void>
}
```

**队列出队策略：**

```ts
private dequeueSteeringMessages(): AgentMessage[] {
  if (this.steeringMode === 'one-at-a-time') {
    // 每次只取一条，保证 Agent 能对每条 Steering 独立响应
    return this.steeringQueue.length > 0 ? [this.steeringQueue.shift()!] : []
  }
  // 'all' 模式：一次性取出全部
  const all = this.steeringQueue.slice()
  this.steeringQueue = []
  return all
}
```

**面试考点：**

1. 为什么 `prompt()` 要检查 `isStreaming`？防止并发 prompt。
2. `continue()` 的使用场景：Retry（上一条是 error 的 assistant 消息）或 Steering 后续处理。
3. `waitForIdle()` 的实现：内部用 Promise + resolveRunningPrompt，每次 `_runLoop` 结束时 resolve。

---

### 2.6 消息类型系统 — 可扩展的联合类型

**文件：** `packages/agent/src/types.ts`

```ts
// 基础 LLM 消息
export type Message = UserMessage | AssistantMessage | ToolResultMessage

// 扩展接口（通过 TS 声明合并扩展）
export interface CustomAgentMessages {
  // 空默认 — 应用通过 declare module 扩展
}

// Agent 消息 = LLM 消息 + 自定义消息的联合
export type AgentMessage = Message | CustomAgentMessages[keyof CustomAgentMessages]
```

**扩展方式：**

```ts
// 应用代码中：
declare module '@mariozechner/agent' {
  interface CustomAgentMessages {
    artifact: ArtifactMessage
    notification: NotificationMessage
  }
}
```

**面试考点：** 这是 TypeScript 的 **声明合并（Declaration Merging）** 技巧，类似于 Express 的 `Request` 扩展。好处是保持类型安全的同时允许第三方扩展。

---

### 2.7 AssistantMessageEvent — DFA 状态机

**文件：** `packages/ai/src/types.ts`

```ts
export type AssistantMessageEvent =
  | { type: 'start'; partial: AssistantMessage }
  | { type: 'text_start'; contentIndex: number; partial: AssistantMessage }
  | { type: 'text_delta'; contentIndex: number; delta: string; partial: AssistantMessage }
  | { type: 'text_end'; contentIndex: number; content: string; partial: AssistantMessage }
  | { type: 'thinking_start'; ... }
  | { type: 'thinking_delta'; ... }
  | { type: 'thinking_end'; ... }
  | { type: 'toolcall_start'; ... }
  | { type: 'toolcall_delta'; ... }
  | { type: 'toolcall_end'; ... }
  | { type: 'done'; reason: 'stop' | 'length' | 'toolUse'; ... }
  | { type: 'error'; reason: 'aborted' | 'error'; ... }
```

这是一个**确定性有限状态自动机（DFA）**：

```
start → [text_start → text_delta* → text_end]*
      → [thinking_start → thinking_delta* → thinking_end]*
      → [toolcall_start → toolcall_delta* → toolcall_end]*
      → done | error
```

**设计优势：**

1. 每个事件携带完整的 `partial` 消息快照 → UI 可以在任意时刻渲染完整状态
2. `contentIndex` 标识当前操作的内容块索引 → 支持多内容块的交错流式
3. 事件类型是判别联合（Discriminated Union）→ TypeScript 可以在 switch 中精确推断类型

---

### 2.8 Proxy Stream — 服务端代理架构

**文件：** `packages/agent/src/proxy.ts`

```ts
export function streamProxy(model, context, options: ProxyStreamOptions) {
  // 1. 发起 HTTP POST 到代理服务器
  const response = await fetch(`${options.proxyUrl}/api/stream`, {
    headers: { Authorization: `Bearer ${options.authToken}` },
    body: JSON.stringify({ model, context, options: { temperature, maxTokens, reasoning } })
  })

  // 2. SSE 流式解析
  while (true) {
    const { done, value } = await reader.read()
    buffer += decoder.decode(value, { stream: true })
    const lines = buffer.split('\n')
    buffer = lines.pop() || ''

    for (const line of lines) {
      if (line.startsWith('data: ')) {
        const proxyEvent = JSON.parse(line.slice(6).trim())
        // 3. 从 ProxyEvent 重建完整的 partial 消息（服务器为省带宽不传 partial）
        const event = processProxyEvent(proxyEvent, partial)
        stream.push(event)
      }
    }
  }
}
```

**ProxyEvent vs AssistantMessageEvent：**

- 服务器发送的 ProxyEvent 不含 `partial` 字段（节省带宽）
- 客户端通过 `processProxyEvent` 维护 `partial` 消息，与标准 AssistantMessageEvent 对齐

---

### 2.9 跨 Provider 消息兼容性

**文件：** `packages/ai/src/providers/transform-messages.ts`

当用户在对话中切换模型（如 Claude → GPT）时，历史消息格式需要转换：

```ts
export function transformMessages(messages, model, normalizeToolCallId?) {
  return messages.map(msg => {
    if (msg.role === 'assistant') {
      const isSameModel = msg.provider === model.provider && msg.api === model.api
      return {
        ...msg,
        content: msg.content.flatMap(block => {
          if (block.type === 'thinking') {
            if (block.redacted) return isSameModel ? block : [] // 加密思考只对同模型有效
            if (isSameModel && block.thinkingSignature) return block // 签名思考保留
            if (!block.thinking?.trim()) return [] // 空思考丢弃
            if (isSameModel) return block
            return { type: 'text', text: block.thinking } // 跨模型：thinking 降级为 text
          }
          if (block.type === 'toolCall' && !isSameModel) {
            // 移除 thoughtSignature，标准化 toolCallId
            delete normalizedToolCall.thoughtSignature
            normalizedToolCall.id = normalizeToolCallId(toolCall.id, model, msg)
          }
        })
      }
    }
  })
}
```

**关键兼容问题：**

- OpenAI tool call ID 可达 450+ 字符含特殊字符，Anthropic 要求 `^[a-zA-Z0-9_-]+$` 最多 64 字符
- Google 的 `thoughtSignature` 是不透明签名，跨模型无效
- 加密的 `redacted` thinking 只对同一模型有效

---

### 2.10 Context Overflow 检测

**文件：** `packages/ai/src/utils/overflow.ts`

```ts
const OVERFLOW_PATTERNS = [
  /prompt is too long/i, // Anthropic
  /exceeds the context window/i, // OpenAI
  /input token count.*exceeds the maximum/i, // Google
  /context window exceeds limit/i // MiniMax
  // ... 15+ patterns
]

export function isContextOverflow(message: AssistantMessage, contextWindow?: number): boolean {
  // Case 1: 错误消息匹配
  if (message.stopReason === 'error' && message.errorMessage) {
    if (OVERFLOW_PATTERNS.some(p => p.test(message.errorMessage!))) return true
  }
  // Case 2: 静默溢出（z.ai 等不报错的 provider）
  if (contextWindow && message.usage.input > contextWindow) return true
  return false
}
```

**面试考点：** 为什么不用一个统一的错误码？因为各 Provider 的错误格式不统一，有的返回 400，有的返回正常响应但 content 为空。这是防御性编程的典型案例。

---

### 2.11 Session 持久化 — JSONL + 树结构

**Session 文件格式：**

```jsonl
{"type":"session","version":3,"id":"uuid","timestamp":"...","cwd":"/path"}
{"type":"message","id":"a1b2c3d4","parentId":null,"message":{"role":"user",...}}
{"type":"message","id":"e5f6g7h8","parentId":"a1b2c3d4","message":{"role":"assistant",...}}
{"type":"compaction","id":"i9j0k1l2","parentId":"e5f6g7h8","summary":"...","firstKeptEntryId":"..."}
```

**树结构实现：** 每个 Entry 有 `id`（8字符十六进制）和 `parentId`，形成 DAG。

**分支操作：**

- **In-place Branch**：修改当前 leaf 指针到另一个 entry（不创建新文件）
- **Fork**：创建新的 JSONL 文件，header 中有 `parentSession` 引用

**Compaction 流程：**

1. 从最新消息向前遍历，累计 token 直到超过 `keepRecentTokens`（默认 20KB）
2. 截断点之前的消息提取出来
3. LLM 生成摘要（结构化，跟踪文件操作记录）
4. 写入 `CompactionEntry`，包含 `summary` + `firstKeptEntryId`
5. 重载 Session：system prompt + 摘要 + firstKeptEntryId 之后的消息

---

### 2.12 Extension 系统 — 完整的插件生命周期

**事件钩子体系：**

```
session_start → before_agent_start → turn_start
  → tool_call (可拦截/修改)
  → tool_result (可修改)
  → message_start → message_update → message_end
  → turn_end
→ session_before_compact → session_before_fork → session_shutdown
```

**Extension API 核心能力：**

```ts
export default function(pi: ExtensionAPI) {
  pi.on("tool_call", (event, ctx) => {
    if (event.input.dangerous) return { block: true, reason: "已拦截" }
  })
  pi.registerTool({ name: "my-tool", execute: async (...) => ... })
  pi.registerCommand("hello", { handler: async (args, ctx) => ... })
  pi.registerShortcut("ctrl+shift+g", { ... })
}
```

---

## 三、模拟面试题 & 深度解答

### Q1：请描述 pi-mono 的整体架构分层，以及各层的核心职责

**参考答案：**

pi-mono 采用三层架构：

**底层 `pi-ai`**：LLM Provider 统一抽象层。核心抽象是 `ApiProvider` 接口，通过注册表模式（`registerApiProvider`）支持多厂商接入。每个 Provider 实现 `stream()` 和 `streamSimple()` 两个函数。`streamSimple` 多了 `reasoning` 参数处理。这一层还负责：

- 消息类型定义（`UserMessage`, `AssistantMessage`, `ToolResultMessage`）
- `EventStream` 异步可迭代器
- 工具参数的 JSON Schema 校验（基于 AJV + TypeBox）
- 跨 Provider 消息兼容（`transformMessages`）
- Context Overflow 检测（多 Provider 正则匹配 + 静默溢出检测）

**中间层 `agent-core`**：Agent 循环引擎。核心是 `runLoop()` 的双层循环：

- 内层循环处理 Tool Calls + Steering（用户中断）
- 外层循环处理 FollowUp（任务追加）
- 工具串行执行，每次执行后检查 Steering 实现可中断调度
- 消息系统支持通过 TypeScript 声明合并扩展自定义消息类型

**顶层 `coding-agent`**：完整产品。提供：

- 三种运行模式（Interactive TUI / Print 流式输出 / RPC JSON 协议）
- `AgentSession` 作为核心中间层，管理 Session 持久化、Extension 生命周期、Compaction
- JSONL 树结构 Session 文件，支持分支和 Fork
- Extension 系统（事件钩子 + 工具注册 + 命令注册）
- Skills 系统（Markdown + YAML frontmatter 的能力包）
- Prompt Template 系统

---

### Q2：EventStream 是如何实现异步生产-消费模式的？与 Node.js 的 Readable Stream 有什么区别？

**参考答案：**

`EventStream` 实现了 `AsyncIterable<T>` 接口，核心是一个**双缓冲**机制：

```
生产者 push(event)
  ├─ 有等待中的消费者？→ 直接 resolve 其 Promise（零拷贝投递）
  └─ 没有？→ 入缓冲队列

消费者 for await...of
  ├─ 缓冲队列有数据？→ 直接 yield
  └─ 没有？→ 创建 Promise 挂入 waiting 数组，等生产者 push
```

与 Node.js Readable Stream 的**关键区别**：

| 维度     | EventStream                    | Node.js Readable                                |
| -------- | ------------------------------ | ----------------------------------------------- |
| 协议     | `AsyncIterable` (for-await-of) | `EventEmitter` + `data`/`end` 事件              |
| 背压     | 无背压（所有事件都缓冲）       | 有背压（`highWaterMark`, `pause()`/`resume()`） |
| 结果提取 | `result()` 返回终止事件的结果  | 无此概念                                        |
| 完成判定 | 通过 `isComplete` 谓词函数     | 固定的 `end` 事件                               |
| 类型安全 | 泛型 `<T, R>` 完全类型化       | 需要手动类型断言                                |

**设计取舍**：EventStream 故意不实现背压，因为 LLM 流式响应的速率相对 consumer（UI 渲染）很低，丢帧不如缓冲全量事件。而 `result()` 方法使得可以选择"只要最终结果"的使用模式。

---

### Q3：Agent Loop 的 Steering 和 FollowUp 机制是如何实现的？为什么要分两种？

**参考答案：**

**实现：**

Steering 和 FollowUp 都是通过回调函数注入 Loop 的：

```ts
interface AgentLoopConfig {
  getSteeringMessages?: () => Promise<AgentMessage[]> // 每次工具执行后检查
  getFollowUpMessages?: () => Promise<AgentMessage[]> // Agent 自然结束后检查
}
```

在 `Agent` 类中，两者分别对应内部队列：

```ts
steer(m: AgentMessage) { this.steeringQueue.push(m) }
followUp(m: AgentMessage) { this.followUpQueue.push(m) }
```

**Steering 的中断语义：**

- 在 `executeToolCalls` 中，每执行完一个工具后调用 `getSteeringMessages()`
- 如果返回非空，**跳过剩余工具调用**（标记为 "Skipped due to queued user message"）
- Steering 消息被注入到 context 中，下一轮 LLM 调用会看到
- 效果：用户说"停下来，改做这个" → Agent 放弃未执行的工具，转向新指令

**FollowUp 的续航语义：**

- 在 `runLoop` 的外层循环中，Agent 自然结束后调用 `getFollowUpMessages()`
- 如果返回非空，设为 `pendingMessages`，重新进入内层循环
- 效果：用户在 Agent 运行时排队消息 → Agent 完成当前任务后自动处理

**为什么要分两种？**

这反映了两种不同的**交互意图**：

1. **Steering = 紧急打断**：用户发现 Agent 走偏了，需要立刻纠正方向。代价是丢弃未执行的工具调用。
2. **FollowUp = 排队追加**：用户想到了后续任务，但不急，等 Agent 做完当前的再说。不打断当前流程。

此外两种队列还可以独立配置出队策略：`one-at-a-time`（逐条）vs `all`（一次性全部）。

---

### Q4：跨 Provider 消息兼容性是如何处理的？举个具体例子

**参考答案：**

`transformMessages()` 函数在每次 LLM 调用前执行，负责将历史消息转换为当前模型可接受的格式。
**具体例子：用户先用 Claude，再切到 GPT-4o**
假设历史消息中有一条 Claude 的 assistant 消息包含 thinking block：

```ts
{
  role: "assistant",
  content: [
    { type: "thinking", thinking: "Let me analyze...", thinkingSignature: "sig_abc" },
    { type: "text", text: "Here's the answer..." },
    { type: "toolCall", id: "toolu_01X...", name: "read", arguments: {...} }
  ],
  provider: "anthropic",
  api: "anthropic-messages",
  model: "claude-sonnet-4-20250514"
}
```

切换到 GPT-4o 后，`transformMessages` 检测到 `!isSameModel`，执行：

1. **thinking block** → 降级为 text block：

   ```ts
   { type: "thinking", thinking: "Let me analyze..." }
   // 变为
   { type: "text", text: "Let me analyze..." }
   ```

   因为 OpenAI 不理解 `thinking` 类型。如果 thinking 是 `redacted`（加密），则直接丢弃。

2. **textSignature** → 被剥离（OpenAI 不认识 Anthropic 的签名）

3. **toolCall ID** → 标准化：
   - Anthropic 生成 `toolu_01XjY...`（40+ 字符）
   - 某些 Provider 要求特定格式（如 Mistral 要求恰好 9 位字母数字）
   - 通过 `normalizeToolCallId` 回调映射，并同步更新对应 `toolResult` 的 `toolCallId`

4. **thoughtSignature**（Google 特有）→ 跨模型时移除
   **关键难点：** ID 映射必须是**双向**的——改了 toolCall 的 ID，对应的 toolResult 也必须同步改。所以第一遍遍历 assistant 消息建立 `toolCallIdMap`，第二遍处理 toolResult 时查表替换。

---

### Q5：Compaction（上下文压缩）的完整流程是什么？如果一轮对话就超过了 keepRecentTokens 会怎样？

**参考答案：**

**触发条件：** 上下文 token 数 > `contextWindow - reserveTokens`（默认保留 16KB 给 response）

**完整流程：**

1. **定位截断点**：从最新消息向前遍历，累计 token 直到超过 `keepRecentTokens`（默认 20KB），该位置即为截断点。
2. **提取旧消息**：从上一次 compaction 的 `firstKeptEntryId`（或消息列表开头）到截断点之间的消息。
3. **LLM 摘要**：用当前模型生成结构化摘要，包括：
   - 对话要点
   - 文件操作记录（`readFiles`, `modifiedFiles`）
4. **写入 CompactionEntry**：
   ```jsonl
   {"type":"compaction","id":"xxx","parentId":"yyy",
    "summary":"...", "firstKeptEntryId":"zzz",
    "tokensBefore":150000, "details":{"readFiles":[...],"modifiedFiles":[...]}}
   ```
5. **重载 Session**：Agent 的上下文变为：
   ```
   [system prompt] + [compaction summary 作为 user 消息] + [firstKeptEntryId 之后的消息]
   ```

**Split Turn 问题：**

如果某个 Turn（一次 assistant 回复 + 其工具调用/结果）的 token 数就超过了 `keepRecentTokens`，截断点会落在 Turn 的**中间**。

处理方式：

1. 截断点之前的部分生成 Summary A（历史摘要）
2. 截断点到 Turn 结束的部分生成 Summary B（不完整 Turn 的摘要）
3. `合并 A + B 作为最终 compaction summary`

这确保了即使单个 Turn 很长（比如 Agent 用 bash 读了大量文件），也能正确压缩。

---

### Q6：Extension 系统的 tool_call 拦截机制是怎么工作的？

**参考答案：**

Extension 通过 `pi.on("tool_call", handler)` 注册钩子。在 `executeToolCalls` 实际调用工具之前，`ExtensionRunner` 会：

1. 将工具调用信息封装成 `ToolCallEvent`
2. 依次调用所有注册了 `tool_call` 的 extension handler
3. Handler 可以返回：
   - `{ block: true, reason: "..." }` → **拦截**，工具不执行，返回 reason 作为 tool result（isError=true）
   - `{ modify: { arguments: {...} } }` → **修改**参数后再执行
   - `undefined` / 无返回 → 放行

**实际应用场景：**

- 安全审计：拦截危险的 `bash` 命令（如 `rm -rf /`）
- 参数注入：给 `read` 工具自动注入项目路径前缀
- 权限控制：在 RPC 模式下限制某些工具只能被特定客户端调用

**Extension 的工具包装（wrapper.ts）：**

```ts
function wrapToolWithExtensions(tool: AgentTool, extensionRunner: ExtensionRunner): AgentTool {
  return {
    ...tool,
    execute: async (toolCallId, params, signal, onUpdate) => {
      // 1. 触发 tool_call 事件
      const decision = await extensionRunner.emitToolCall(tool.name, params)
      if (decision?.block) {
        return { content: [{ type: 'text', text: decision.reason }], details: {} }
      }
      // 2. 使用可能被修改的参数执行
      const actualParams = decision?.modify?.arguments ?? params
      const result = await tool.execute(toolCallId, actualParams, signal, onUpdate)
      // 3. 触发 tool_result 事件
      await extensionRunner.emitToolResult(tool.name, result)
      return result
    }
  }
}
```

---

### Q7：pi-mono 的 Model 类型系统是如何实现类型安全的？

**参考答案：**

`Model` 接口被 `Api` 类型参数化：

```ts
export interface Model<TApi extends Api> {
  id: string
  name: string
  api: TApi              // 类型参数决定了可用的 Provider
  provider: Provider
  baseUrl: string
  reasoning: boolean     // 是否支持思考模式
  cost: { input, output, cacheRead, cacheWrite }
  contextWindow: number
  maxTokens: number
  compat?: TApi extends 'openai-completions' ? OpenAICompletionsCompat : ...
}
```

`getModel()` 函数使用 **条件类型 + 字面量类型推断**：

```ts
export function getModel<
  TProvider extends KnownProvider,
  TModelId extends keyof (typeof MODELS)[TProvider]
>(provider: TProvider, modelId: TModelId): Model<ModelApi<TProvider, TModelId>>
```

这意味着：

```ts
const model = getModel('google', 'gemini-2.5-flash-lite-preview-06-17')
// model 的类型是 Model<'google-generative-ai'>
// 编译时就知道它使用 Google 的 API
```

`compat` 字段使用条件类型：

```ts
compat?: TApi extends 'openai-completions'
  ? OpenAICompletionsCompat
  : TApi extends 'openai-responses'
    ? OpenAIResponsesCompat
    : never
```

只有 `openai-completions` 的模型才能访问 `compat` 中的 `supportsReasoningEffort` 等字段。

---

### Q8：请分析 Tool 参数校验的实现，以及 TypeBox 在其中的作用

**参考答案：**

**TypeBox** 是一个 JSON Schema 生成库，它的 `Type.Object(...)` 既生成 TS 类型又生成 JSON Schema：

```ts
import { Type, type Static, type TSchema } from '@sinclair/typebox'

const ReadToolParams = Type.Object({
  file_path: Type.String({ description: 'Absolute path to file' }),
  offset: Type.Optional(Type.Number({ description: 'Start line' })),
  limit: Type.Optional(Type.Number({ description: 'Max lines' }))
})

// Static<typeof ReadToolParams> = { file_path: string, offset?: number, limit?: number }
```

**校验流程（`validateToolArguments`）：**

```ts
export function validateToolArguments(tool: Tool, toolCall: ToolCall): any {
  if (!ajv || isBrowserExtension) {
    return toolCall.arguments // 浏览器扩展环境跳过校验（CSP 限制）
  }

  const validate = ajv.compile(tool.parameters) // 编译 JSON Schema
  const args = structuredClone(toolCall.arguments) // 深拷贝！

  if (validate(args)) {
    return args // AJV 会原地修改 args（类型强制转换）
  }

  throw new Error(`Validation failed for tool "${toolCall.name}": ...`)
}
```

**关键细节：**

1. **`structuredClone`**：因为 AJV 的 `coerceTypes: true` 会原地修改对象（如 string "123" → number 123），所以先深拷贝。
2. **CSP 兼容**：Chrome Extension Manifest V3 禁止 `eval`/`Function`，而 AJV 内部使用它们编译 schema。检测到此环境时跳过校验。
3. **AJV singleton**：全局只创建一个实例，避免重复编译 schema 的性能开销。

---

### Q9：proxy.ts 的 SSE 解析有什么值得注意的工程细节？

**参考答案：**

```ts
while (true) {
  const { done, value } = await reader.read()
  buffer += decoder.decode(value, { stream: true }) // 关键：stream: true
  const lines = buffer.split('\n')
  buffer = lines.pop() || '' // 关键：保留不完整的最后一行

  for (const line of lines) {
    if (line.startsWith('data: ')) {
      const data = line.slice(6).trim()
      if (data) {
        const proxyEvent = JSON.parse(data)
        stream.push(processProxyEvent(proxyEvent, partial))
      }
    }
  }
}
```

**工程细节：**

1. **`TextDecoder({ stream: true })`**：多字节 UTF-8 字符可能跨越 chunk 边界，`stream: true` 确保解码器不会将不完整的多字节序列当作错误。

2. **`buffer = lines.pop() || ""`**：SSE 数据可能在 chunk 边界处被截断。`split("\n")` 后最后一个元素可能是不完整的行，必须保留到下一个 chunk 拼接。

3. **ProxyEvent 带宽优化**：服务器不发送 `partial` 字段（完整的 AssistantMessage 快照），客户端通过 `processProxyEvent` 自行维护 `partial` 状态。一条 30 字的 text_delta 如果带完整 partial 可能有 500+ 字节的开销。

4. **AbortController 集成**：注册了 `abort` 事件监听器，abort 时取消 reader。finally 块中移除监听器防止内存泄漏。

5. **错误恢复**：catch 中区分 `signal.aborted`（用户主动取消，reason="aborted"）和其他错误（网络问题等，reason="error"），不同的 stopReason 会影响上层重试逻辑。

---

### Q10：如果要为 MiniMax 模型添加一个新的 Provider，需要做哪些工作？

**参考答案：**

从代码中可以看到 MiniMax 已经在 `KnownProvider` 中声明了（`'minimax' | 'minimax-cn'`），以及在 overflow 检测中有对应的正则：`/context window exceeds limit/i`。

添加新 Provider 的步骤：

1. **定义模型**：在 `models.generated.ts` 中添加 MiniMax 的模型定义：

   ```ts
   'minimax': {
     'abab7-chat': {
       id: 'abab7-chat', name: 'MiniMax ABAB7', api: 'openai-completions',
       provider: 'minimax', baseUrl: 'https://api.minimax.chat/v1',
       cost: { input: X, output: Y, ... }, contextWindow: 245760, maxTokens: 16384,
       reasoning: false, input: ['text', 'image'],
     }
   }
   ```

2. **如果使用 OpenAI 兼容 API**（最常见路径）：
   - 设置 `api: 'openai-completions'`
   - 在 `compat` 中配置兼容性选项：
     ```ts
     compat: {
       supportsStore: false,
       supportsDeveloperRole: false,
       maxTokensField: 'max_tokens',
       supportsUsageInStreaming: true,
     }
     ```

3. **如果 MiniMax 有独特的 API 格式**（需要自定义 Provider）：
   - 创建 `src/providers/minimax.ts`
   - 实现 `ApiProvider` 接口：
     ```ts
     export const minimaxProvider: ApiProvider<'minimax-api', MinimaxStreamOptions> = {
       api: 'minimax-api',
       stream: (model, context, options) => { ... },
       streamSimple: (model, context, options) => { ... },
     }
     ```
   - 在 `register-builtins.ts` 中注册：
     ```ts
     registerApiProvider(minimaxProvider)
     ```

4. **添加 Overflow 检测**（已有）：在 `overflow.ts` 的 `OVERFLOW_PATTERNS` 中确认有匹配 MiniMax 错误消息的正则。

5. **API Key 环境变量**：在 `env-api-keys.ts` 中添加 `MINIMAX_API_KEY` 的读取。

---

### Q11：从前端工程角度，这个项目有哪些值得借鉴的 TypeScript 高级用法？

**参考答案：**

1. **判别联合（Discriminated Union）**：
   `AgentEvent`、`AssistantMessageEvent` 都是以 `type` 字段做判别的联合类型，配合 switch/case 实现穷举检查。

2. **声明合并（Declaration Merging）**：
   `CustomAgentMessages` 接口使用 `declare module` 扩展，让第三方可以添加自定义消息类型而不修改核心代码。

3. **条件类型 + 索引访问类型**：

   ```ts
   type ModelApi<TProvider, TModelId> = (typeof MODELS)[TProvider][TModelId] extends {
     api: infer TApi
   }
     ? TApi
     : never
   ```

4. **泛型约束传递**：
   `Model<TApi extends Api>` 确保 API 类型在整个调用链中传递，从 `getModel()` 到 `stream()` 到 Provider 实现。

5. **TypeBox 双重身份**：既生成 TypeScript 类型（`Static<T>`）又生成 JSON Schema（运行时校验），实现了"一次定义，编译时 + 运行时双重保障"。

6. **异步迭代器协议**：`EventStream` 实现 `Symbol.asyncIterator`，使得 LLM 流可以直接用 `for await...of` 消费。

7. **工厂函数 + DI**：`createAgentSession()` 接受所有依赖作为参数，便于测试时注入 mock。

---

### Q12：如果让你设计一个 Agent 的"暂停-恢复"功能，基于现有架构你会怎么实现？

**参考答案：**

现有架构其实已经非常接近支持暂停恢复了：

**暂停：**

1. 调用 `agent.abort()`（已有），触发 AbortController cancel
2. 在 `_runLoop` 的 catch 中，当 `stopReason === "aborted"` 时，保存当前 `newMessages` 到 session
3. Session 的 JSONL 树结构天然支持"回溯到任意点"

**恢复：**

1. `agent.continue()`（已有）可以从当前 context 继续
2. 如果最后一条消息是 assistant（中间有工具调用未完成），利用 Steering 机制注入一条"请继续之前未完成的任务"的用户消息
3. Session 的 `navigateTree(entryId)` 可以回到暂停点

**需要新增的能力：**

1. **工具执行状态快照**：如果暂停发生在工具执行中间（比如 bash 正在运行），需要记录哪些工具已完成、哪些被跳过
2. **AbortSignal 传递到工具**：已实现（`execute` 函数接收 `signal` 参数）
3. **持久化暂停元数据**：新增一个 `PauseEntry` 类型到 Session JSONL 中

关键点是 **Agent Loop 的可中断设计（串行工具执行 + Steering 检查）已经为暂停恢复奠定了基础**。

---

### Q13：这个框架的性能瓶颈可能在哪里？你会怎么优化？

**参考答案：**

1. **消息序列化**：每次 LLM 调用前都要 `convertToLlm(messages)` 遍历全部消息。优化：增量转换，只转换新增部分。

2. **AJV Schema 编译**：每次 `validateToolArguments` 都 `ajv.compile(tool.parameters)`。优化：缓存编译结果（用 WeakMap 以 TSchema 对象为 key）。

3. **structuredClone**：每次校验前深拷贝参数。对于大参数（如文件内容），开销显著。优化：验证通过后直接返回原对象，失败时才需要错误信息中的原始值。

4. **EventStream 无背压**：如果 LLM 响应速度很快但 UI 渲染慢，事件会无限堆积。优化：添加可选的背压机制，或在 UI 层做节流。

5. **JSONL Session 加载**：每次启动需要解析整个 JSONL 文件重建树。对于长 session（数千轮），启动慢。优化：
   - 建立索引文件（记录 entry offset）
   - 懒加载（只加载当前 branch 的 entries）
   - SQLite 替代 JSONL

6. **Compaction 的 LLM 调用**：压缩本身需要一次 LLM 调用（生成摘要），在 token 紧张时可能导致连锁压缩。优化：使用更小的模型做摘要，或本地 embedding + 聚类。

---

### Q14：从 Agent 产品设计的角度，Steering vs FollowUp 两种交互方式分别适合什么场景？

**参考答案：**

**Steering（中断转向）— "我改主意了"**

- Agent 正在读错误的文件 → "不是那个文件，读 src/utils.ts"
- Agent 的 bash 命令可能有破坏性 → "停！不要执行 rm"
- Agent 走偏了 → "你理解错了需求，应该是..."
- 技术特征：跳过未执行的工具调用，立即转向

**FollowUp（排队追加）— "做完这个再做那个"**

- Agent 在修 bug → 用户想到 "改完后跑一下测试"
- Agent 在重构 → 用户追加 "顺便更新文档"
- 批量任务 → 依次排队多个任务
- 技术特征：不打断当前工作，等 Agent 自然结束后再处理

**one-at-a-time vs all 的考量：**

- Steering 默认 `one-at-a-time`：每条打断指令都应该被 Agent 独立处理
- FollowUp 默认 `one-at-a-time`：如果一次给 Agent 5 个任务，它可能混淆优先级
- 但可以设为 `all`：当用户发送的是补充信息（"顺便告诉你，密码是 xxx，环境是 prod"）时，一次性给出更好

---

### Q15：请解释 pi-mono 中的"树形 Session"设计，与线性对话历史相比有什么优势？

**参考答案：**

**数据结构：**

每个 Entry 有 `id`（8 字符 hex）和 `parentId`，形成一棵树（实际是 DAG）：

```
         a1b2 (user: "fix the bug")
           │
         c3d4 (assistant: "I'll try approach A...")
        ╱    ╲
     e5f6     g7h8 (fork: "try approach B instead")
   (toolResult)   (toolResult)
       │            │
     i9j0         k1l2
   (assistant)   (assistant)
```

**核心操作：**

- `getBranch()`: 从当前 leaf 沿 parentId 向上遍历 → 得到线性视图（喂给 LLM）
- `getTree()`: 返回完整树（用于 UI 树形浏览器）
- `switchBranch(entryId)`: 切换到另一个叶子节点（in-place，不创建新文件）
- `fork(entryId)`: 从某个节点分叉创建新的 Session 文件

**优势：**

1. **非破坏性回溯**：线性历史中"撤销"意味着删除消息。树结构中，回到之前的节点只需切换 leaf 指针，所有分支都保留。

2. **方案对比**：从同一个节点分出多个分支，可以让 Agent 分别尝试不同方案，然后对比结果。

3. **Compaction 安全**：压缩只影响当前 branch 的视图，其他分支不受影响。每个分支可以有不同的 compaction 摘要。

4. **协作与审计**：树结构天然记录了"为什么走这条路"的决策轨迹，适合 Code Review。

5. **存储效率**：JSONL append-only，分支不需要复制数据。Fork 创建新文件时也只存增量，通过 `parentSession` 引用原始数据。

---

## 四、高频概念速查

| 概念                  | 位置                                 | 核心要点                                        |
| --------------------- | ------------------------------------ | ----------------------------------------------- |
| EventStream           | ai/utils/event-stream.ts             | 异步可迭代 + 双缓冲 + result() Promise          |
| API Registry          | ai/api-registry.ts                   | Map 注册表 + sourceId 热拔插                    |
| Agent Loop            | agent/agent-loop.ts                  | 双层 while 循环 + Steering 中断 + FollowUp 续航 |
| transformMessages     | ai/providers/transform-messages.ts   | thinking 降级 + toolCallId 映射 + redacted 处理 |
| validateToolArguments | ai/utils/validation.ts               | AJV + TypeBox + structuredClone + CSP 兼容      |
| isContextOverflow     | ai/utils/overflow.ts                 | 15+ 正则 + 静默溢出检测                         |
| AgentMessage          | agent/types.ts                       | 声明合并扩展自定义消息类型                      |
| Session Tree          | coding-agent/core/session-manager.ts | JSONL + id/parentId DAG + getBranch()/getTree() |
| Compaction            | coding-agent/core/compaction/        | 找截断点 → LLM 摘要 → 重载 Session              |
| Extension             | coding-agent/core/extensions/        | 事件钩子 + tool_call 拦截 + 工具/命令注册       |

---

## 五、面试策略建议

1. **从架构分层讲起**：先描述三层（ai → agent-core → coding-agent），展示你对全局的理解
2. **深入核心循环**：Agent Loop 的双层循环 + Steering/FollowUp 是最有技术含量的设计
3. **强调 TypeScript 功力**：声明合并、条件类型、判别联合、异步迭代器协议
4. **展示工程思维**：跨 Provider 兼容、Context Overflow 防御、Session 树结构的优势
5. **给出自己的见解**：性能瓶颈分析、暂停恢复的设计方案等，展示你不只是读懂了代码，还能基于它设计新功能
