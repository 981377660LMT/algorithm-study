# @mariozechner/pi-agent-core 全面讲解

**源码部分：**

- **types.ts** — `AgentMessage` 的声明合并设计、`AgentTool` 工具抽象、`AgentEvent` 事件类型、`AgentState` 状态定义
- **agent-loop.ts** — 核心双层循环（外层处理 follow-up，内层处理 tool calls + steering）、`transformContext → convertToLlm` 双层转换管道、工具执行与 steering 中断机制
- **agent.ts** — `Agent` 类的 `prompt()`/`continue()` 三种重载、steering/follow-up 队列管理、并发保护、事件订阅
- **proxy.ts** — 代理流的客户端 partial message 重建

**测试部分（逐个用例讲解）：**

- **agent-loop.test.ts** — 7 个测试：基本事件流、自定义消息过滤、上下文裁剪、工具调用循环、steering 中断跳过剩余工具、continue 基本行为、自定义消息 continue
- **agent.test.ts** — 9 个测试：默认状态、自定义初始化、事件订阅、状态修改器、队列机制、并发保护、continue 处理队列消息、one-at-a-time 语义、sessionId 传递
- **e2e.test.ts** — 5 个场景 × 6 个 Provider 的端到端集成
- **bedrock-models.test.ts** — 模型兼容性差异处理

末尾附了整体架构图和关键设计思想对照表。

## 一、项目定位

这是一个 **有状态的 AI Agent 框架**，核心功能：

1. 管理与 LLM 的多轮对话
2. 支持工具调用（Tool Use）
3. 通过事件流（Event Stream）实时通知 UI
4. 支持在 Agent 运行中插入"转向消息"（Steering）和"后续消息"（Follow-up）

构建在 `@mariozechner/pi-ai` 之上，后者提供底层的 LLM 流式调用能力。

---

## 二、文件结构总览

```
src/
  types.ts        → 所有类型定义（AgentMessage, AgentState, AgentTool, AgentEvent 等）
  agent-loop.ts   → 核心循环逻辑（agentLoop / agentLoopContinue）
  agent.ts        → Agent 类（面向用户的高层封装）
  proxy.ts        → 代理流函数（通过服务器中转 LLM 请求）
  index.ts        → 统一导出

test/
  agent-loop.test.ts   → agent-loop 的单元测试（用 Mock 模拟 LLM）
  agent.test.ts        → Agent 类的单元测试
  e2e.test.ts          → 端到端集成测试（真实调用各大 LLM 提供商）
  bedrock-models.test.ts → Amazon Bedrock 模型兼容性测试
  utils/
    calculate.ts       → 示例工具：计算器
    get-current-time.ts → 示例工具：获取当前时间
```

---

## 三、核心类型 (types.ts)

### 3.1 AgentMessage — 消息的统一抽象

```typescript
// LLM 原生消息 + 自定义消息的联合类型
type AgentMessage = Message | CustomAgentMessages[keyof CustomAgentMessages]
```

**关键设计**：LLM 只认识 `user`/`assistant`/`toolResult` 三种消息，但应用可能需要自定义消息（如通知、UI 状态等）。`AgentMessage` 通过 TypeScript 的**声明合并（Declaration Merging）**让用户扩展：

```typescript
// 用户代码中:
declare module '@mariozechner/pi-agent-core' {
  interface CustomAgentMessages {
    notification: { role: 'notification'; text: string; timestamp: number }
  }
}
```

这样 `AgentMessage` 就自动包含了 `notification` 类型。

### 3.2 AgentTool — 工具定义

```typescript
interface AgentTool<TParameters, TDetails> extends Tool<TParameters> {
  label: string // 显示名称
  execute: (
    // 执行函数
    toolCallId: string,
    params: Static<TParameters>,
    signal?: AbortSignal,
    onUpdate?: AgentToolUpdateCallback<TDetails> // 进度回调
  ) => Promise<AgentToolResult<TDetails>>
}
```

工具返回 `AgentToolResult`，包含：

- `content`: 文本/图片内容数组（给 LLM 看的）
- `details`: 给 UI 展示的详细信息

### 3.3 AgentEvent — 事件类型

```
agent_start / agent_end          → Agent 生命周期
turn_start / turn_end            → 每轮（一次 LLM 调用 + 工具执行）
message_start / update / end     → 消息生命周期
tool_execution_start / update / end → 工具执行生命周期
```

### 3.4 AgentState — Agent 状态

```typescript
interface AgentState {
  systemPrompt: string
  model: Model<any>
  thinkingLevel: ThinkingLevel // "off" | "minimal" | "low" | "medium" | "high" | "xhigh"
  tools: AgentTool[]
  messages: AgentMessage[] // 完整对话历史
  isStreaming: boolean // 是否正在流式输出
  streamMessage: AgentMessage | null // 当前部分消息（流式中）
  pendingToolCalls: Set<string> // 正在执行的工具 ID
  error?: string
}
```

---

## 四、核心循环 (agent-loop.ts)

### 4.1 消息流转管道

这是整个框架最重要的数据流设计：

```
AgentMessage[] → transformContext() → AgentMessage[] → convertToLlm() → Message[] → LLM
```

- **transformContext**：可选。在 `AgentMessage` 层面做操作，如裁剪旧消息、注入外部上下文
- **convertToLlm**：必需。把 `AgentMessage[]` 转为 LLM 能理解的 `Message[]`，过滤自定义类型

### 4.2 两个入口函数

#### `agentLoop(prompts, context, config, signal?, streamFn?)`

用用户消息启动一轮新对话。流程：

```
1. 发射 agent_start、turn_start
2. 发射 message_start/end（用户消息）
3. 进入 runLoop（主循环）
```

#### `agentLoopContinue(context, config, signal?, streamFn?)`

从现有上下文继续，不添加新消息。用于**重试**场景。前置条件：最后一条消息不能是 `assistant`。

### 4.3 runLoop — 主循环详解

```
外层 while(true)：处理 follow-up 消息
  └─ 内层 while(hasMoreToolCalls || pendingMessages)：处理工具调用和转向消息
       ├─ 注入 pending 消息
       ├─ 流式获取 LLM 回复（streamAssistantResponse）
       ├─ 如果有工具调用 → 执行工具（executeToolCalls）
       ├─ 检查转向消息
       └─ turn_end
  └─ 检查 follow-up 消息，有则继续外层循环
```

**关键细节**：

- 每次工具执行后都会检查 steering 消息
- 检测到 steering → 跳过剩余工具，注入 steering 消息，让 LLM 响应中断
- 所有工具执行完、无 steering → 检查 follow-up，有则继续

### 4.4 工具执行 (executeToolCalls)

```typescript
for (每个 toolCall) {
  1. 发射 tool_execution_start
  2. 查找工具 → 验证参数 → 执行
  3. 发射 tool_execution_end
  4. 生成 toolResultMessage
  5. 检查 steering 消息 → 有则跳过剩余工具
}
```

跳过的工具会生成错误结果：`"Skipped due to queued user message."`

---

## 五、Agent 类 (agent.ts)

### 5.1 设计思路

`Agent` 是面向用户的高层封装，职责：

1. 管理 `AgentState`
2. 将 `prompt()` / `continue()` 转发给 `agentLoop` / `agentLoopContinue`
3. 消费事件流，更新内部状态
4. 提供事件订阅机制
5. 管理 steering/follow-up 队列

### 5.2 prompt() — 发送消息

```typescript
// 三种重载：
await agent.prompt("Hello");                              // 纯文本
await agent.prompt("看图", [{ type: "image", ... }]);     // 带图片
await agent.prompt({ role: "user", content: "...", ... }); // 自定义 AgentMessage
```

内部做了什么：

1. 检查是否已在 streaming（是则抛错）
2. 构造 `AgentMessage[]`
3. 调用 `_runLoop(messages)`

### 5.3 continue() — 继续对话

用于重试或处理队列中的消息。智能判断：

- 最后一条是 `assistant` → 尝试出队 steering / follow-up 消息
- 最后一条是 `user` / `toolResult` → 直接继续

### 5.4 \_runLoop() — 内部运行

```
1. 创建 AbortController
2. 设置 isStreaming = true
3. 构建 AgentLoopConfig（包含 convertToLlm、getSteeringMessages、getFollowUpMessages）
4. 调用 agentLoop / agentLoopContinue
5. 消费事件流：
   - message_start → 设置 streamMessage
   - message_update → 更新 streamMessage
   - message_end → 追加到 messages，清空 streamMessage
   - tool_execution_start/end → 更新 pendingToolCalls
   - agent_end → 设置 isStreaming = false
6. 每个事件都 emit 给外部订阅者
```

### 5.5 Steering 和 Follow-up

**Steering（转向）**：在 Agent 运行中插入消息，打断当前工作。

```typescript
agent.steer({ role: 'user', content: '算了，换个方向', timestamp: Date.now() })
```

**Follow-up（后续）**：在 Agent 完成后追加消息。

```typescript
agent.followUp({ role: 'user', content: '顺便总结一下', timestamp: Date.now() })
```

两种模式：

- `"one-at-a-time"`（默认）：每次只取一条
- `"all"`：一次全部取出

### 5.6 错误处理

执行出错时，Agent 会：

1. 创建一个 `stopReason: "error"` 的 assistant 消息
2. 追加到消息历史
3. 设置 `state.error`
4. 发射 `agent_end`

中止（abort）时类似，但 `stopReason: "aborted"`。

---

## 六、代理流 (proxy.ts)

`streamProxy` 用于通过中间服务器转发 LLM 请求，典型场景：

```
客户端 → 代理服务器（管理 API Key） → LLM Provider
```

关键设计：

- 服务器发送 **SSE 流**，但**省略了 `partial` 字段**以节省带宽
- 客户端在本地**重建 partial message**
- 支持 `text_delta`、`thinking_delta`、`toolcall_delta` 等增量事件

---

## 七、测试用例详解

### 7.1 agent-loop.test.ts — 循环逻辑的单元测试

用 `MockAssistantStream` 模拟 LLM 响应，不发真实 API 请求。

#### 测试 1: 基本事件流

```
输入: "Hello"
Mock LLM 返回: "Hi there!"
验证: 事件序列包含 agent_start → turn_start → message_start → message_end → turn_end → agent_end
结果: 2 条新消息（用户 + 助手）
```

**学到什么**：最简单的一次交互会产生完整的事件生命周期。

#### 测试 2: 自定义消息类型 + convertToLlm

```
上下文包含: [notification 消息]  ← 自定义类型
输入: "Hello"
convertToLlm: 过滤掉 notification，只保留 user/assistant/toolResult
验证: LLM 只收到 1 条消息（用户消息），notification 被过滤
```

**学到什么**：`convertToLlm` 是自定义消息类型的关键桥梁，决定 LLM 看到什么。

#### 测试 3: transformContext 上下文裁剪

```
上下文: [旧消息1, 旧回复1, 旧消息2, 旧回复2]
输入: "new message"
transformContext: 只保留最后 2 条
验证: LLM 最终收到 2 条消息（裁剪后的）
```

**学到什么**：`transformContext` 在 `convertToLlm` 之前执行，用于上下文窗口管理。

#### 测试 4: 工具调用

```
输入: "echo something"
Mock LLM 第 1 次调用: 返回 toolCall{name: "echo", args: {value: "hello"}}
工具执行: 返回 "echoed: hello"
Mock LLM 第 2 次调用: 返回 "done"
验证: 工具被执行，事件包含 tool_execution_start/end
```

**学到什么**：工具调用会触发额外的 LLM 轮次，LLM 先调用工具 → 得到结果 → 再次回复。

#### 测试 5: Steering 中断 + 跳过剩余工具

```
输入: "start"
Mock LLM 返回: 两个 toolCall [echo("first"), echo("second")]
Steering: 执行完第一个工具后，返回 "interrupt" 消息
验证:
  - 只有第一个工具被真正执行
  - 第二个工具被跳过，结果标记为 isError: true
  - 中断消息出现在上下文中
```

**学到什么**：Steering 机制让用户可以实时打断 Agent 的工作流程。

#### 测试 6: agentLoopContinue 基本行为

```
上下文: [已有用户消息]
不添加新消息，直接继续
验证:
  - 只返回新的助手消息（不含已有消息）
  - 不发射用户消息事件
```

**学到什么**：`agentLoopContinue` 用于从现有上下文恢复，常用于重试。

#### 测试 7: 自定义消息做为最后消息的 continue

```
上下文: [custom 消息]
convertToLlm: 将 custom 转为 user 消息
验证: 不报错，LLM 收到转换后的 user 消息
```

**学到什么**：`continue` 只检查原始 role 不是 `assistant` 即可，具体转换交给 `convertToLlm`。

---

### 7.2 agent.test.ts — Agent 类的单元测试

#### 测试 1: 默认状态初始化

```typescript
const agent = new Agent()
// systemPrompt: ""
// model: gemini-2.5-flash-lite
// thinkingLevel: "off"
// tools: [], messages: []
// isStreaming: false, streamMessage: null
```

**学到什么**：Agent 有合理的默认值，可以零配置创建。

#### 测试 2: 自定义初始状态

```typescript
const agent = new Agent({
  initialState: {
    systemPrompt: 'You are a helpful assistant.',
    model: getModel('openai', 'gpt-4o-mini'),
    thinkingLevel: 'low'
  }
})
```

**学到什么**：通过 `initialState` 部分覆盖默认值。

#### 测试 3: 事件订阅

```typescript
const unsub = agent.subscribe((event) => { ... });
unsub(); // 取消订阅
```

**学到什么**：状态修改器（setSystemPrompt 等）不触发事件，只有 prompt/continue 等流式操作才会。

#### 测试 4: 状态修改器

测试了所有 setter：`setSystemPrompt`、`setModel`、`setThinkingLevel`、`setTools`、`replaceMessages`（创建副本）、`appendMessage`、`clearMessages`。

#### 测试 5: Steering/Follow-up 队列

```typescript
agent.steer(message) // 消息入队，不进 messages
agent.followUp(message) // 同上
```

**学到什么**：队列消息只在 Agent 运行时被消费。

#### 测试 6: 并发保护

```
agent.prompt("First")  → 开始流式
agent.prompt("Second") → 抛错 "Agent is already processing"
agent.continue()       → 抛错 "Agent is already processing"
```

**学到什么**：Agent 禁止并发 prompt，要用 `steer()` 或 `followUp()` 排队。

#### 测试 7: continue() 处理队列消息

```
messages: [user, assistant]  ← 最后是 assistant
followUp 队列: ["Queued follow-up"]
agent.continue() → 出队 follow-up → 新一轮对话
验证: follow-up 消息出现在历史中，最后一条是 assistant
```

**学到什么**：`continue()` 会自动处理队列中的 steering/follow-up 消息。

#### 测试 8: one-at-a-time steering 语义

```
messages: [user, assistant]
steering 队列: ["Steering 1", "Steering 2"]
agent.continue()
验证: 产生了 2 次 LLM 调用（每次取一条 steering）
消息顺序: [user, assistant, user, assistant]
```

**学到什么**：`one-at-a-time` 模式下，每轮只消费一条 steering 消息。

#### 测试 9: sessionId 传递

```typescript
agent.sessionId = 'session-abc'
agent.prompt('hello')
// streamFn 收到 options.sessionId === "session-abc"
```

**学到什么**：`sessionId` 透传给 LLM Provider，用于缓存等场景。

---

### 7.3 e2e.test.ts — 端到端集成测试

用真实 API 测试多个 LLM Provider（需要 API Key）。

#### 公用测试场景

| 场景                  | 输入                                    | 验证                                |
| --------------------- | --------------------------------------- | ----------------------------------- |
| basicPrompt           | "What is 2+2?"                          | 回复包含 "4"                        |
| toolExecution         | "Calculate 123 \* 456"                  | 工具执行返回 56088                  |
| abortExecution        | 发送后 100ms abort                      | 最后消息 stopReason: "aborted"      |
| stateUpdates          | "Count from 1 to 5"                     | 事件包含 message_update（流式增量） |
| multiTurnConversation | "My name is Alice" → "What is my name?" | 回复包含 "alice"                    |

#### 覆盖的 Provider

- Google (gemini-2.5-flash)
- OpenAI (gpt-4o-mini)
- Anthropic (claude-haiku-4-5)
- xAI (grok-3)
- Groq (openai/gpt-oss-20b)
- Cerebras (gpt-oss-120b)

每个 Provider 都会跑上面 5 个场景。

---

### 7.4 bedrock-models.test.ts — Bedrock 兼容性测试

测试 Amazon Bedrock 上的各种模型，主要处理模型间的差异：

| 问题类别         | 说明                              |
| ---------------- | --------------------------------- |
| 需要推理配置文件 | 部分模型不支持按需调用            |
| 无效模型 ID      | 某些区域没有的模型                |
| maxTokens 超限   | 配置的 token 上限超过模型实际限制 |
| 不支持 reasoning | 重放对话时拒绝推理内容            |
| 签名格式验证     | Anthropic 新模型验证签名格式      |

---

### 7.5 测试工具 (utils/)

#### calculate.ts — 计算器工具

```typescript
const calculateTool: AgentTool = {
  name: 'calculate',
  label: 'Calculator',
  description: 'Evaluate mathematical expressions',
  parameters: Type.Object({
    expression: Type.String({ description: '...' })
  }),
  execute: async (_id, args) => {
    const result = new Function(`return ${args.expression}`)()
    return {
      content: [{ type: 'text', text: `${args.expression} = ${result}` }],
      details: undefined
    }
  }
}
```

#### get-current-time.ts — 时间工具

```typescript
const getCurrentTimeTool: AgentTool = {
  name: "get_current_time",
  label: "Current Time",
  parameters: Type.Object({
    timezone: Type.Optional(Type.String({ description: "..." }))
  }),
  execute: async (_id, args) => {
    return { content: [...], details: { utcTimestamp: Date.now() } };
  },
};
```

---

## 八、整体架构图

```
┌─────────────────────────────────────────────┐
│                用户代码                      │
│  agent.prompt("Hello")                      │
│  agent.subscribe(event => ...)              │
│  agent.steer(...) / agent.followUp(...)     │
└───────────────┬─────────────────────────────┘
                │
                ▼
┌─────────────────────────────────────────────┐
│              Agent 类 (agent.ts)             │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  │
│  │  State   │  │ Queues   │  │  Events  │  │
│  │ 管理     │  │ steering │  │  订阅    │  │
│  │          │  │ followUp │  │  发射    │  │
│  └──────────┘  └──────────┘  └──────────┘  │
└───────────────┬─────────────────────────────┘
                │
                ▼
┌─────────────────────────────────────────────┐
│         Agent Loop (agent-loop.ts)           │
│                                              │
│  外层循环: follow-up                         │
│   └─ 内层循环: tool calls + steering         │
│       ├─ streamAssistantResponse()           │
│       │   ├─ transformContext()              │
│       │   ├─ convertToLlm()                 │
│       │   └─ streamFn() → LLM               │
│       └─ executeToolCalls()                  │
│           └─ 每个工具后检查 steering          │
└───────────────┬─────────────────────────────┘
                │
                ▼
┌─────────────────────────────────────────────┐
│     Stream Function (可替换)                 │
│  ┌──────────────┐  ┌──────────────────┐     │
│  │ streamSimple │  │  streamProxy     │     │
│  │ (直接调用)    │  │  (代理服务器)    │     │
│  └──────────────┘  └──────────────────┘     │
└───────────────┬─────────────────────────────┘
                │
                ▼
         LLM Provider (OpenAI / Anthropic / Google / ...)
```

---

## 九、关键设计思想总结

| 设计点             | 做法                                   | 好处                                          |
| ------------------ | -------------------------------------- | --------------------------------------------- |
| 消息抽象           | `AgentMessage` = LLM 消息 + 自定义消息 | 应用可加自定义消息类型，同时保持对 LLM 的兼容 |
| 双层转换           | `transformContext` → `convertToLlm`    | 上下文管理和消息格式转换解耦                  |
| 事件流             | 细粒度的事件类型                       | UI 可精确响应每个生命周期阶段                 |
| Steering/Follow-up | 队列 + 轮询机制                        | 用户可实时打断或追加任务，无需等待            |
| 可替换 StreamFn    | `streamSimple` / `streamProxy`         | 直连或代理，灵活切换部署方式                  |
| 声明合并           | TypeScript `declare module`            | 类型安全地扩展消息类型，无需修改库代码        |
| 单次并发           | `isStreaming` 锁                       | 防止并发 prompt 导致状态混乱                  |

---

这个库叫做 `@mariozechner/pi-agent-core`，它的核心功能是**帮助你快速搭建一个“聪明、带记忆、能使用工具、且能像打字机一样流式回复”的 AI 智能体（Agent）**。
下面我带你**逐条段落对照**，用通俗易懂的语言把这份 README 彻底拆解一遍，绝不遗漏每个细节。

---

### 1. 标题与简介 (Title Series)

**原文**：Stateful agent with tool execution and event streaming. Built on `@mariozechner/pi-ai`.
**讲解**：

- **Stateful (带状态的)**：意思是这个大模型不是“阅后即焚”的，它记性很好，能自动帮你管理之前的聊天记录。
- **tool execution (工具执行)**：它不仅仅能聊天，还能根据你的配置自动调用代码（比如去帮你读取一个本地文件）。
- **event streaming (事件流)**：它支持像 ChatGPT 网页版那样，一个字一个字地把结果吐给你，而不是让你干等好几秒才给你一整段话。
- 它底层是基于 `@mariozechner/pi-ai` 这个基础库构建的。

### 2. 安装 (Installation)

**讲解**：在你的项目终端里运行 `npm install @mariozechner/pi-agent-core` 即可安装这个包。

### 3. 快速开始 (Quick Start)

**讲解**：这里给出了一个最简单、能跑起来的代码例子，分三步：

1.  **创建 Agent**：用 `new Agent` 实例化一个智能体。在 `initialState` 里给它设定了人设（`systemPrompt: "You are a helpful assistant."`），并指定了要调用的大模型（这里调用了 Anthropic 的 Claude 模型）。
2.  **监听开口说话 (subscribe)**：`agent.subscribe` 相当于给 AI 身上挂了个窃听器。当它检测到 `event.type === "message_update"`（消息更新中）并且是纯文本时，就把刚生成的字（`delta`）打印到屏幕上。这实现了**流式输出**。
3.  **发送问题 (prompt)**：`await agent.prompt("Hello!");` 这就是你对 AI 说的第一句话。

### 4. 核心概念 (Core Concepts)

这里讲了该库设计的核心思想：

#### AgentMessage vs LLM Message (智能体消息 vs 大模型消息)

- **痛点**：通常的大模型（LLM）很笨，它只认识三种角色：`user`(用户)、`assistant`(AI助手)、`toolResult`(工具返回的结果)。但是在你开发稍微复杂点的应用时，你可能需要一些**特殊的系统消息**，比如给页面发一个 `notification`（通知）。
- **解法**：这个库引入了灵活的 `AgentMessage`（UI 和你能看懂的所有消息）。但在把消息发送给大模型之前，有一个必经的关卡叫 `convertToLlm`，它的作用就是把大模型看不懂的特殊消息过滤掉或转换掉，只把纯粹的对话喂给大模型。

#### Message Flow (消息流转图)

当你发出一句话后，它经历了什么：
`你发出的所有消息` -> `【可选】截断缩减 (transformContext)` -> `【必须】转换并过滤给大模型看的消息 (convertToLlm)` -> `发送给大模型处理 (LLM)`

---

### 5. 事件流 (Event Flow，极其重要)

AI 思考和做事的过程会被打碎成一个个“事件广播”出来。这部分教你如何根据广播来刷新你前端的 UI。

#### prompt("Hello") 发出简单聊天时的顺序：

1. `agent_start` (开始干活) -> `turn_start` (开始这一路对话)
2. `message_start/end` (处理你发的 Hello)
3. `message_start/update.../end` (AI 开始思考，`update` 一点点往屏幕蹦字，最后结束)
4. `turn_end` -> `agent_end` (彻底干完闭嘴)

#### 带有工具调用时的顺序 ("Read config.json")：

如果 AI 发现需要调用它手里的工具（比如读取文件工具）：

1. 还是先输出文字（告诉你它要去用工具了）。
2. 发出 `tool_execution_start/update/end` 广播：告诉你它**开始用工具 -> 工具执行中 -> 工具执行完毕拿到结果**。
3. 工具拿到结果后，它会**自动发起新的一轮(turn_start)**，让大模型根据刚偷看的文件结果，再整理成人类语言回复给你。

#### continue() 事件

也就是“接茬继续”。如果不小心网络断了或报错了，你不需要把聊天记录重传一遍，只要调用 `await agent.continue();`，它就会从出问题的地方接着跑。

---

### 6. 配置选项 (Agent Options)

初始化 `new Agent({...})` 时能传哪些高级参数：

- `initialState`：初始配置（上面的快速开始里用过，包括大模型、基础聊天记录等）。
- `convertToLlm`：前文提到的过滤函数，决定哪些消息传给大模型。
- `transformContext`：在消息太多时，用来剔除老旧消息或注入外部信息的函数。
- `steeringMode` / `followUpMode`：中断与追加模式限制配置（后面专门讲）。
- `streamFn` / `sessionId` / `getApiKey`：代理相关的网络请求、缓存和动态 Token 刷新设置。
- `thinkingBudgets`：针对一些“深度思考模型”(像 OpenAI o1/o3 或者 Claude3.7) 设置它们思考使用的 token 额度预算。

### 7. 智能体状态 (Agent State)

你可以随时通过 `agent.state` 偷窥 AI 当前的内部状态。这里面存放了系统提示词、绑定的模型、所有的历史消息（`messages`）、当前是否还在打字（`isStreaming`）、以及正在执行还未完成的工具列表等。

### 8. 方法大全 (Methods)

这里是一本 API 字典，列举了你怎么操作这个 AI：

- **Prompting (提示/发问)**：除了普通文本，还支持传图片、自己拼装原始消息，以及用 `continue()` 继续发话。
- **State Management (状态管理)**：你可以在中途动态换人设 (`setSystemPrompt`)、换模型 (`setModel`)、设置工具、甚至清空整个聊天记忆 (`clearMessages` 或者 `reset()`)。
- **控制 (Control)**：`agent.abort()` 强行打断 AI 说话或者打断它用工具；`waitForIdle()` 等待 AI 彻底歇下来。
- **事件 (Events)**：`subscribe((event) => {...})` 挂载监听，它会返回一个 `unsubscribe` 函数，调用就能取消窃听。

---

### 9. 智能引导与追加任务 (Steering and Follow-up)

这是非常高级也是非常实用的 UI 控制功能，也就是**“中途插嘴”**机制。

- `agent.steer({...})`（**急迫插嘴**）：当 AI 正在慢吞吞地执行某个复杂的工具时，你突然说：“停停停！别查那个了，查这个！” `steer` 会立刻把手里没干完的工具当作“失败/跳过”掐掉，并且把你的新指令塞进去强行让它重新回应你。
- `agent.followUp({...})`（**排队追加**）：当 AI 正在好好干活时，你突然想起什么，不想打断它，但想让它干完手头的后接着做。比如：“对了，这一通操作弄完后，顺便给我写个总结。”它会乖乖等现在的事情跑完，再自动开启新的一轮处理你排队的任务。
- 最后还提供了清空这两个插嘴队列的方法（`clearXXXQueue`）。

### 10. 自定义消息类型 (Custom Message Types)

**讲解**：针对 TypeScript 玩家。如果你代码里一定要存在一种非人类也非大模型的 UI 专用消息（比如页面通知 `role: "notification"`）。你可以利用 TypeScript 的 `declare module` "声明合并"特性，强行把 `notification` 这个合法的角色塞进库的底层定义中。然后在 `convertToLlm` 把这个奇怪的角色过滤掉，就不会报错了。

### 11. 怎么给 AI 打造工具 (Tools)

AI 本身没法读你的电脑文件，需要你帮它写好“工具”。

- 用 `TypeBox` 这个库去严密地定义这个工具需要人类（其实是 AI）传入什么参数（比如必传一个 `[path: 文件路径]`）。
- **核心逻辑 `execute`**：定义工具实际上怎么干活的（比如用了 `fs.readFile` 去读文件）。
- **错误处理规则 (Error Handling)**：**重点！** 如果工具执行出错了（比如文件没找到），不要温和地 `return "没找到出错了"`，而是要**粗暴地 `throw new Error()` 抛出异常**。底层框架抓到这个报错后，会自动转换成“标准错误报告”拿给大模型看，大模型看到这种红牌警告才能更好地进行自我纠错。

### 12. 代理服务用法 (Proxy Usage)

如果你是在网页（浏览器前端）里用这个库，为了防黑客把你的 API 密钥偷走，你肯定不能直接在前端填密钥。这时候用自定义的 `streamProxy`，让前端发请求给你的后端 Node.js，让后端拿着密钥去访问大模型。

### 13. 底层 API (Low-Level API)

如果你觉得 `Agent` 这个大类太笨重了，我就想自己从零捏一个循环咋整？库还提供了原生的生成器函数 `agentLoop` 和 `agentLoopContinue`。用 `for await` 循环就可以直接拿到最纯粹的底层数据流来自己定制逻辑。这是给“高阶老鸟”玩的功能。

---

**总结给你的学习建议**：
既然你是新手，你只需要先按 **2安装 -> 3跑通 Quick Start -> 11照猫画虎写一个工具** 这三个步骤来，你的 AI 就能带上强大的工具库替你打工了，其他的插嘴排队功能（8和9）等你把基本盘玩熟了再看！有哪里没搞清的随时问我。
