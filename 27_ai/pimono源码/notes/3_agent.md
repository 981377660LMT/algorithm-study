# @mariozechner/pi-agent-core 深度精读

**架构核心：无状态引擎 + 有状态壳**

整个包只有 ~700 行代码，分三层：`Agent`（状态管理）→ `agentLoop`（纯函数循环）→ `streamFn`（LLM 通信）。`agentLoop` 是无状态的，可以独立使用，Agent 类只是可选的便利层。

**最精妙的 5 个设计：**

1. **声明合并扩展消息类型** — `CustomAgentMessages` 空接口让消费者零侵入地添加自定义消息类型（notification、artifact 等），库代码不用改
2. **Steering vs Follow-up 双队列** — steering 在工具执行间隙打断 agent 改方向，follow-up 在 agent 完成后追加任务，精确对应两种真实用户行为
3. **串行工具执行** — 有意牺牲并行性能，换取每个工具执行后都能检查用户中断的能力；在 agent 场景下，用户控制权比速度更重要
4. **`transformContext` + `convertToLlm` 两阶段管道** — 上下文管理（裁剪/压缩）和协议转换（自定义消息→LLM 消息）关注点分离
5. **`skipInitialSteeringPoll`** — 防止 steering 消息在 `continue()` 时被重复消费的状态泄漏防护，只有真正实现过 agent 循环的人才会踩到这个坑

## 一、全局架构：三层分离

整个 agent 包只有 **4 个源文件**，约 700 行代码，但架构非常清晰：

```
┌─────────────────────────────────────────────────┐
│  Agent (agent.ts)         — 有状态的外壳         │
│  ┌─────────────────────────────────────────────┐ │
│  │  agentLoop (agent-loop.ts) — 无状态纯函数引擎 │ │
│  │  ┌─────────────────────────────────────────┐ │ │
│  │  │  streamSimple / streamProxy  — LLM 通信  │ │ │
│  │  └─────────────────────────────────────────┘ │ │
│  └─────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────┘
```

| 层          | 职责                     | 有无状态 |
| ----------- | ------------------------ | -------- |
| `Agent` 类  | 状态管理、事件分发、队列 | 有状态   |
| `agentLoop` | turn 调度、工具执行      | 无状态   |
| `streamFn`  | LLM 底层通信             | 无状态   |

**洞见：** Agent 类本质上是一个 **状态管理壳**，核心逻辑全在无状态的 `agentLoop` 中。这种设计意味着你可以完全绕过 Agent 类，直接用 `agentLoop` 构建自己的状态管理方案（比如接入 Redux 或 Svelte store），而不会失去任何核心能力。

---

## 二、types.ts — 类型系统设计哲学

### 2.1 AgentMessage：用声明合并实现开放扩展

```typescript
// 空接口 — 等待消费者通过 declaration merging 注入
export interface CustomAgentMessages {}

// AgentMessage = LLM 原生消息 | 任何自定义消息
export type AgentMessage = Message | CustomAgentMessages[keyof CustomAgentMessages]
```

这是整个包最精妙的设计：

- LLM 只认 `user / assistant / toolResult` 三种角色
- 但真实 UI 需要 notification、artifact、status 等自定义消息
- 传统做法是搞一个大 union type，每加一种消息类型就要改库代码
- 这里用 **TS 声明合并（declaration merging）** 让消费者在自己的代码里扩展类型，库代码零改动

```typescript
// 消费者侧
declare module '@mariozechner/agent' {
  interface CustomAgentMessages {
    artifact: { role: 'artifact'; html: string; timestamp: number }
  }
}
// 现在 AgentMessage 自动包含 artifact 类型，convertToLlm 里过滤掉即可
```

**洞见：** 这是 **开放-封闭原则** 在类型系统层面的体现。库对修改封闭，对扩展开放 — 不是通过继承，而是通过 TS 结构类型的声明合并。

### 2.2 convertToLlm + transformContext：两阶段管道

```
AgentMessage[] → transformContext() → AgentMessage[] → convertToLlm() → Message[] → LLM
                    (可选)                                (必需)
```

为什么分两步？

- `transformContext` 工作在 **AgentMessage 层面**：你可以在这裁剪上下文窗口、注入 RAG 检索结果，操作的是你的业务消息
- `convertToLlm` 工作在 **LLM 协议层面**：把你的自定义消息过滤/转换为标准 LLM 消息

**洞见：** 这是 **关注点分离** — 上下文管理的逻辑（多少消息、怎么压缩）和协议转换的逻辑（怎么跟 LLM 说话）是两个不同的关注点，不应耦合在一起。

### 2.3 AgentTool：从 Tool 到可执行

```typescript
export interface AgentTool<TParameters extends TSchema, TDetails = any> extends Tool<TParameters> {
  label: string // UI 展示用
  execute: (
    toolCallId: string,
    params: Static<TParameters>, // 类型安全的参数
    signal?: AbortSignal, // 可取消
    onUpdate?: AgentToolUpdateCallback<TDetails> // 流式更新
  ) => Promise<AgentToolResult<TDetails>>
}
```

注意 `params: Static<TParameters>` — 用 `@sinclair/typebox` 的 `Static` 从 JSON Schema 类型推导出 TS 类型。这意味着：

1. 工具的 parameters 用 JSON Schema 定义（LLM 需要）
2. execute 的参数自动获得 TypeScript 类型安全
3. 一份 schema 两头用，不需要手写两套类型

### 2.4 AgentEvent：细粒度的生命周期事件

```typescript
export type AgentEvent =
  | { type: 'agent_start' } // 整个运行开始
  | { type: 'agent_end' } // 整个运行结束
  | { type: 'turn_start' } // 一轮开始（一次 LLM 调用 + 工具执行）
  | { type: 'turn_end' } // 一轮结束
  | { type: 'message_start' } // 任何消息开始
  | { type: 'message_update' } // assistant 消息流式更新
  | { type: 'message_end' } // 消息结束
  | { type: 'tool_execution_start' } // 工具开始
  | { type: 'tool_execution_update' } // 工具流式更新
  | { type: 'tool_execution_end' } // 工具结束
```

三层嵌套的生命周期：`Agent > Turn > Message/Tool`，使得 UI 可以在任意粒度做出响应动画或状态更新。

---

## 三、agent-loop.ts — 无状态的核心引擎

### 3.1 双层循环：外循环处理 follow-up，内循环处理 tool call + steering

```
runLoop()
┌─ while (true)                          // 外循环：follow-up
│  ┌─ while (hasMoreToolCalls || pending) // 内循环：tool call + steering
│  │   1. 注入 pendingMessages
│  │   2. streamAssistantResponse()  → 调用 LLM
│  │   3. 如果有 tool call → executeToolCalls()
│  │   4. 检查 steering 消息
│  │   5. 如果 steering → 跳过剩余 tool call，注入 steering
│  └─ end while
│  检查 follow-up 消息 → 有则 continue，无则 break
└─ end while
```

**洞见：** 这个双层循环是 "agentic loop" 模式的标准实现，但它的精妙之处在于 **steering 和 follow-up 的区分**：

- **Steering**（转向）：在工具执行间隙检查，**会打断**当前工具链，让 agent 改变方向
- **Follow-up**（追问）：只在 agent 觉得"完事了"才检查，**不会打断**工作，只会追加新任务

这正好对应了真实 chat UI 的两种用户行为：

1. "停下！换个方向" → steering
2. "做完了？再帮我做个 xxx" → follow-up

### 3.2 streamAssistantResponse()：LLM 调用边界

```typescript
async function streamAssistantResponse(...) {
  // 1. transformContext（可选）
  let messages = context.messages;
  if (config.transformContext) {
    messages = await config.transformContext(messages, signal);
  }

  // 2. convertToLlm（必需）
  const llmMessages = await config.convertToLlm(messages);

  // 3. 构造 LLM context
  const llmContext: Context = {
    systemPrompt: context.systemPrompt,
    messages: llmMessages,
    tools: context.tools,
  };

  // 4. 动态解析 API key（应对过期 token）
  const resolvedApiKey = config.getApiKey
    ? await config.getApiKey(config.model.provider)
    : config.apiKey;

  // 5. 调用 LLM 流
  const response = await streamFunction(config.model, llmContext, { ...config, apiKey: resolvedApiKey, signal });

  // 6. 处理流式事件
  for await (const event of response) {
    // ... 更新 partial message，推送事件
  }
}
```

**洞见：** `getApiKey` 的设计解决了一个真实的工程问题 — OAuth token 过期。工具执行可能耗时很长（比如运行代码、访问 API），等工具执行完再调 LLM 时 token 可能已经过期了。每次 LLM 调用前重新拿 key，这是从实际踩过坑后的设计。

### 3.3 executeToolCalls() 的 steering 中断机制

```typescript
for (let index = 0; index < toolCalls.length; index++) {
  // ... 执行第 i 个 tool call

  // 每执行完一个 tool，检查 steering
  if (getSteeringMessages) {
    const steering = await getSteeringMessages()
    if (steering.length > 0) {
      steeringMessages = steering
      // 跳过剩余 tool call，标记为 "Skipped"
      const remainingCalls = toolCalls.slice(index + 1)
      for (const skipped of remainingCalls) {
        results.push(skipToolCall(skipped, stream))
      }
      break
    }
  }
}
```

**关键设计：** tool call 是**串行执行**的（不是并行），这样才能在每个 tool 之间插入 steering 检查点。被跳过的 tool call 会收到一个 `isError: true` 的结果消息 "Skipped due to queued user message"，这样 LLM 就知道这些工具没执行。

**洞见：** 选择串行而非并行执行 tool call 是一个有意的权衡 — 牺牲了并行性能，换来了用户随时中断的能力。在 AI agent 场景下，用户控制权比执行速度更重要。

### 3.4 EventStream：异步生成器 + Promise 的完美融合

`agentLoop` 返回 `EventStream<AgentEvent, AgentMessage[]>`，这个类来自 `pi-ai`，它同时支持：

- `for await...of` 逐个消费事件
- `.result()` 获取最终结果

```typescript
const stream = agentLoop([userMessage], context, config)
for await (const event of stream) {
  // 实时处理每个事件
}
const allMessages = await stream.result() // 拿到完整结果
```

内部实现是一个 push-based 的异步队列，生产者（循环逻辑）调 `stream.push(event)`，消费者通过 async iterator 拉取。`createAgentStream()` 用 终止判定函数 + 结果提取函数 构造：

```typescript
new EventStream<AgentEvent, AgentMessage[]>(
  event => event.type === 'agent_end', // 终止条件
  event => (event.type === 'agent_end' ? event.messages : []) // 提取结果
)
```

---

## 四、agent.ts — 有状态的外壳

### 4.1 Agent 的本质：状态容器 + 事件总线 + 队列管理

```
Agent
├── _state: AgentState          // 唯一真相源
├── listeners: Set<Function>     // 观察者模式
├── steeringQueue: AgentMessage[] // 中断队列
├── followUpQueue: AgentMessage[] // 追问队列
├── abortController              // 取消控制
└── runningPrompt                // 空闲等待
```

Agent 类自己不调 LLM，不执行 tool，它做的事情是：

1. 维护 `AgentState`（messages、streaming 状态、error）
2. 把 `AgentLoopConfig` 组装起来，传给 `agentLoop`
3. 消费 `agentLoop` 发出的事件，更新状态并转发给外部 listener
4. 管理 steering 和 follow-up 队列

### 4.2 prompt() 与 continue() 的对称设计

```typescript
// prompt: 新增用户消息 → 开始循环
async prompt(input) {
  msgs = ... // 构造用户消息
  await this._runLoop(msgs);
}

// continue: 不新增消息 → 从当前上下文继续
async continue() {
  // 特殊处理：如果最后是 assistant，尝试从队列取消息
  if (messages[last].role === "assistant") {
    const steering = this.dequeueSteeringMessages();
    if (steering.length > 0) {
      await this._runLoop(steering, { skipInitialSteeringPoll: true });
      return;
    }
    const followUp = this.dequeueFollowUpMessages();
    if (followUp.length > 0) {
      await this._runLoop(followUp);
      return;
    }
    throw new Error("Cannot continue from assistant");
  }
  await this._runLoop(undefined); // 直接继续
}
```

**洞见：** `continue()` + `steer()/followUp()` 的组合覆盖了所有恢复场景：

- **错误重试**：最后是 toolResult → `continue()` 直接重试 LLM 调用
- **用户追加**：最后是 assistant → `followUp()` + `continue()` 追加问题
- **用户中断**：agent 运行中 → `steer()` 立即中断

### 4.3 \_runLoop()：桥接状态和无状态

`_runLoop` 是 Agent 最核心的方法，做了三件事：

1. **组装 config**：把 Agent 内部的队列方法暴露给 agentLoop

   ```typescript
   const config = {
     ...
     getSteeringMessages: async () => this.dequeueSteeringMessages(),
     getFollowUpMessages: async () => this.dequeueFollowUpMessages(),
   };
   ```

2. **消费事件流**：把 agentLoop 的事件映射到内部状态更新

   ```typescript
   for await (const event of stream) {
     switch (event.type) {
       case 'message_end':
         this.appendMessage(event.message)
         break
       case 'tool_execution_start':
         this._state.pendingToolCalls.add(event.toolCallId)
         break
       // ...
     }
     this.emit(event) // 转发给外部 listener
   }
   ```

3. **处理边界情况**：partial message 清理、错误兜底、AbortController 生命周期

### 4.4 skipInitialSteeringPoll 的巧妙设计

当 `continue()` 从 steering 队列取出消息传给 `_runLoop` 时，会设置 `skipInitialSteeringPoll: true`：

```typescript
getSteeringMessages: async () => {
  if (skipInitialSteeringPoll) {
    skipInitialSteeringPoll = false;
    return [];
  }
  return this.dequeueSteeringMessages();
},
```

为什么？因为 `runLoop` 开头会立即调一次 `getSteeringMessages()`，如果不跳过，刚注入的 steering 消息还没被 LLM 处理，又会被当作新的 steering 消息取出来——死循环。

**洞见：** 这种细节只有在真正实现过 agent 循环的人才会意识到。这是一个典型的 **"状态泄漏"防护**。

---

## 五、proxy.ts — 带宽优化的 SSE 代理

### 5.1 核心思路：服务端剥离 `partial`，客户端重建

标准的 LLM streaming 事件每个 delta 都带一个完整的 `partial` 消息（包含到目前为止的全部内容）。对于代理场景，这意味着巨大的带宽浪费。

```typescript
// ProxyAssistantMessageEvent — 没有 partial 字段
export type ProxyAssistantMessageEvent = { type: 'text_delta'; contentIndex: number; delta: string }
// ...
```

服务端只发 delta，客户端自己维护一个 `partial: AssistantMessage` 对象，收到 delta 就拼接，然后构造出标准的 `AssistantMessageEvent`（带 partial）。

### 5.2 SSE 解析

```typescript
buffer += decoder.decode(value, { stream: true })
const lines = buffer.split('\n')
buffer = lines.pop() || '' // 最后一行可能不完整，保留

for (const line of lines) {
  if (line.startsWith('data: ')) {
    const data = line.slice(6).trim()
    const proxyEvent = JSON.parse(data)
    const event = processProxyEvent(proxyEvent, partial)
    stream.push(event)
  }
}
```

这是标准的 SSE（Server-Sent Events）协议解析，手写的原因是浏览器的 `EventSource` API 不支持 POST 请求和自定义 headers。

---

## 六、设计模式与工程洞见总结

### 6.1 "无状态核心 + 有状态壳" 模式

这是我认为这个包**最重要的架构决策**。`agentLoop` 是纯函数（接受输入、产生 EventStream），没有副作用、没有内部状态。`Agent` 类是可选的便利层。

好处：

- 测试时直接测 `agentLoop`，用 mock streamFn，不需要 Agent 实例
- 可以在不同的状态管理框架中复用同一个循环逻辑
- 循环逻辑的正确性和状态管理的正确性可以独立验证

### 6.2 依赖注入贯穿始终

几乎所有行为都通过配置注入，而非硬编码：

| 注入点                | 作用              |
| --------------------- | ----------------- |
| `streamFn`            | 替换 LLM 通信层   |
| `convertToLlm`        | 自定义消息转换    |
| `transformContext`    | 自定义上下文管理  |
| `getApiKey`           | 动态 API key 解析 |
| `getSteeringMessages` | 中断消息来源      |
| `getFollowUpMessages` | 追问消息来源      |

### 6.3 流式优先

从底层 LLM 响应到顶层 UI 更新，全链路都是流式的：

```
LLM stream → streamAssistantResponse → EventStream → Agent.emit → UI listener
```

没有任何地方会"等全部完成再通知"，`message_update` 事件在每个 delta 都会触发。

### 6.4 相比 LangChain/Vercel AI SDK 的区别

| 维度     | pi-agent-core         | LangChain              | Vercel AI SDK       |
| -------- | --------------------- | ---------------------- | ------------------- |
| 代码量   | ~700 行               | 数万行                 | 数千行              |
| 抽象层级 | 薄封装，贴近 LLM 原语 | 重抽象，Chain/Runnable | 中等，偏 React 集成 |
| 消息类型 | 声明合并，零侵入扩展  | 固定类型               | 固定类型            |
| 状态管理 | 可选的 Agent 类       | 内置 memory module     | 依赖 React state    |
| 工具执行 | 串行 + steering 中断  | 并行                   | 并行                |

**核心理念差异：** pi-agent-core 把自己定位为一个 **最小可用的 agent 循环**，只做 "调 LLM → 执行工具 → 再调 LLM" 这件事，把状态管理、UI 渲染、上下文压缩全部留给消费者通过钩子注入。它不试图成为一个框架，而是一个构建块。

---

## 七、运行流程完整追踪

以 `agent.prompt("你好")` 为例，追踪完整的执行路径：

```
1. Agent.prompt("你好")
   → 构造 UserMessage { role: "user", content: [{ type: "text", text: "你好" }] }
   → 调 _runLoop([userMessage])

2. _runLoop()
   → 创建 AbortController
   → 设 isStreaming = true
   → 构造 AgentContext（copy messages）
   → 构造 AgentLoopConfig（注入 convertToLlm, getSteeringMessages 等）
   → 调 agentLoop([userMessage], context, config, signal, streamFn)

3. agentLoop()
   → push agent_start, turn_start
   → push message_start(userMessage), message_end(userMessage)
   → 进入 runLoop()

4. runLoop() — 外循环第 1 次
   → 检查 steering → 无
   → runLoop() — 内循环第 1 次
     → streamAssistantResponse()
       → transformContext()（如有）
       → convertToLlm() → 过滤出 [UserMessage]
       → streamFn(model, llmContext, options)
       → LLM 返回流式事件
       → push message_start(partial)
       → push message_update(partial) × N
       → push message_end(finalAssistantMessage)
     → 检查 tool calls → 假设无
     → push turn_end
   → 内循环结束
   → 检查 follow-up → 无
   → 外循环结束
   → push agent_end
   → stream.end()

5. 回到 _runLoop()
   → for await 消费事件：
     message_end → appendMessage(assistantMessage) 到 state
     agent_end → isStreaming = false
   → finally: 清理 AbortController，resolveRunningPrompt
```

---

## 八、关键细节备忘

1. **AgentState.pendingToolCalls** 是 `Set<string>`，每次更新都创建新 Set — 这是为了触发 UI 框架的响应式更新（引用变了才会重新渲染）
2. **agentLoopContinue** 严格检查最后一条消息不能是 assistant — 因为 LLM 的对话协议要求 assistant 后面必须跟 user 或 toolResult
3. **错误处理**：tool 抛出异常 → 捕获 → 构造 `isError: true` 的 toolResult → LLM 看到错误信息可以自行决定下一步
4. **partial message** 在 streaming 过程中会被 **就地修改**（`context.messages[last] = partialMessage`），这是一个性能优化，避免每个 delta 都复制整个消息数组
