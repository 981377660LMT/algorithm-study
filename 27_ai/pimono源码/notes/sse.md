## SSE 在 Pimono 中的具体使用方式

### 什么是 SSE？

**SSE (Server-Sent Events)** 是一种基于 HTTP 的单向流式协议。客户端发送一个普通 HTTP 请求，服务器保持连接不关闭，持续以特定文本格式推送事件。

SSE 协议格式非常简单：

```
data: {"type":"text_delta","text":"Hello"}

data: {"type":"text_delta","text":" World"}

data: [DONE]

```

每条消息以 `data:` 开头，**每条消息之间以两个换行符 `\n\n` 分隔**。

---

### Pimono 中 SSE 的两种使用方式

#### 方式一：SDK 封装 SSE（OpenAI / Anthropic）

对于 OpenAI 和 Anthropic，Pimono 使用它们的官方 SDK，SDK 内部帮你处理了 SSE 解析。代码只需传 `stream: true`，然后 `for await` 迭代即可：

**Anthropic** (anthropic.ts):

```typescript
// SDK 内部发 HTTP 请求，header 带 Accept: text/event-stream
// SDK 自动解析 SSE 格式，暴露为 AsyncIterable
const anthropicStream = client.messages.stream(
  { ...params, stream: true },
  { signal: options?.signal }
)

// 直接 for-await 迭代 SDK 解析好的事件对象
for await (const event of anthropicStream) {
  if (event.type === 'content_block_delta') {
    // event.delta.text → 一小段增量文本
  }
}
```

**OpenAI Completions** (openai-completions.ts):

```typescript
const openaiStream = await client.chat.completions.create(
  { ...params, stream: true }, // stream: true 触发 SSE
  { signal: options?.signal }
)

for await (const chunk of openaiStream) {
  // chunk.choices[0].delta.content → 增量文本
}
```

这些 SDK 底层都是发 `POST` 请求，服务器返回 `Content-Type: text/event-stream`，SDK 负责解析 `data: ...` 行。

#### 方式二：手动 fetch + 解析 SSE（Codex / Google Gemini CLI）

对于一些供应商，Pimono 自己用 `fetch` 发请求，**手动解析 SSE 流**。这是最能揭示 SSE 工作原理的部分：

**HTTP 请求** (openai-codex-responses.ts):

```typescript
response = await fetch(url, {
  method: 'POST',
  headers: {
    accept: 'text/event-stream', // ← 告知服务器我要 SSE
    'content-type': 'application/json',
    Authorization: `Bearer ${token}`
  },
  body: JSON.stringify(body),
  signal: options?.signal
})
```

**手动解析 SSE** (openai-codex-responses.ts):

```typescript
async function* parseSSE(response: Response): AsyncGenerator<Record<string, unknown>> {
  if (!response.body) return

  const reader = response.body.getReader() // ← 获取底层字节流 reader
  const decoder = new TextDecoder()
  let buffer = ''

  while (true) {
    const { done, value } = await reader.read() // ← 逐块读取
    if (done) break
    buffer += decoder.decode(value, { stream: true }) // ← 字节 → 字符串

    // SSE 协议：两个换行 \n\n 分隔一条消息
    let idx = buffer.indexOf('\n\n')
    while (idx !== -1) {
      const chunk = buffer.slice(0, idx)
      buffer = buffer.slice(idx + 2)

      // 提取 data: 开头的行
      const dataLines = chunk
        .split('\n')
        .filter(l => l.startsWith('data:'))
        .map(l => l.slice(5).trim())

      if (dataLines.length > 0) {
        const data = dataLines.join('\n').trim()
        if (data && data !== '[DONE]') {
          // ← [DONE] 是流结束标志
          yield JSON.parse(data) // ← yield 解析后的 JSON 对象
        }
      }
      idx = buffer.indexOf('\n\n')
    }
  }
}
```

**Google Gemini CLI** (google-gemini-cli.ts) 也是类似的手动解析，URL 甚至直接带了 `?alt=sse` 参数：

```typescript
const requestUrl = `${endpoint}/v1internal:streamGenerateContent?alt=sse`;
const response = await fetch(requestUrl, {
  method: "POST",
  headers: { "Accept": "text/event-stream", ... },
  body: requestBodyJson,
});

// 手动从 response.body 读取并解析 data: 行
const reader = response.body.getReader();
// ... 同样的 buffer + split("\n") + startsWith("data:") 逻辑
```

---

### 完整数据流图

```
   HTTP POST (stream: true)
   Accept: text/event-stream
          │
          ▼
   ┌─────────────────────┐
   │   AI 服务器 (SSE)    │
   │ Content-Type:        │
   │ text/event-stream    │
   └──────┬──────────────┘
          │  持续推送文本块:
          │  data: {"type":"content_block_delta","delta":{"text":"He"}}  \n\n
          │  data: {"type":"content_block_delta","delta":{"text":"llo"}} \n\n
          │  data: {"type":"message_stop"}                               \n\n
          ▼
   ┌─────────────────────┐
   │  SSE 解析层          │  SDK 内置 / 手动 parseSSE()
   │  response.body       │  ReadableStream → reader.read() 逐块
   │  .getReader()        │  buffer 累积 → \n\n 切割 → data: 提取 → JSON.parse
   └──────┬──────────────┘
          │  yield 解析后的 JS 对象
          ▼
   ┌─────────────────────┐
   │  Provider 适配层     │  for await (event of sseStream)
   │  anthropic.ts 等     │  将供应商事件 → 归一化为 AssistantMessageEvent
   │                     │  stream.push({ type: 'text_delta', delta: 'He' })
   └──────┬──────────────┘
          │
          ▼
   ┌─────────────────────┐
   │  EventStream         │  生产者-消费者 AsyncIterable
   │  (event-stream.ts)   │  queue[] + waiting[] 缓冲机制
   └──────┬──────────────┘
          │
          ▼
   ┌─────────────────────┐
   │  上层消费者           │  for await (const e of stream) { ... }
   │  CLI / UI            │  实时展示增量文本
   └─────────────────────┘
```

### 关键要点

| 层面            | 说明                                                                               |
| --------------- | ---------------------------------------------------------------------------------- |
| **协议**        | HTTP/1.1 + `Content-Type: text/event-stream`（SSE 标准）                           |
| **消息格式**    | `data: {JSON}\n\n`，以 `[DONE]` 结束                                               |
| **底层读取**    | `response.body.getReader()` → `ReadableStream` 逐 chunk 读取                       |
| **解析方式**    | 累积 buffer → 按 `\n\n` 切分事件 → 过滤 `data:` 行 → `JSON.parse`                  |
| **备选传输**    | OpenAI Codex 还支持 WebSocket（`transport: 'websocket'`），将 `https:` 转为 `wss:` |
| **AbortSignal** | 全链路传递 `signal`，支持中断：`reader.cancel()` + `socket.close()`                |

---

## Pimono AI 流式传输架构

Pimono 采用了一套**统一的流式传输抽象层**，核心协议是 **SSE (Server-Sent Events)**，同时支持 WebSocket。

### 1. 传输协议

在 types.ts 中定义了支持的传输类型：

```typescript
export type Transport = 'sse' | 'websocket' | 'auto'
```

默认使用 **SSE**，各供应商 SDK（OpenAI、Anthropic 等）底层都通过 HTTP 长连接 + SSE 协议推送增量数据。

### 2. 核心流式架构：`EventStream`（生产者-消费者模式）

核心在 src/utils/event-stream.ts，实现了一个**异步可迭代的事件流**：

```typescript
class EventStream<T, R> implements AsyncIterable<T> {
  private queue: T[] = []           // 事件缓冲队列
  private waiting: ((value) => void)[] = []  // 等待中的消费者
  private done = false

  push(event: T): void { ... }     // 生产者推入事件
  end(result?: R): void { ... }    // 标记流结束
  async *[Symbol.asyncIterator]()   // 消费者通过 for-await-of 消费
  result(): Promise<R>              // 获取最终完整结果
}
```

关键设计：

- **生产者**：各 Provider 异步接收上游 SSE 事件，调用 `push()` 推入统一格式的事件
- **消费者**：通过 `for await (const event of stream)` 消费
- 如果消费者慢于生产者，事件会缓存在 `queue` 中
- 如果消费者等待中，新事件直接投递给 waiting 的 Promise

`AssistantMessageEventStream` 是它的特化版本，当收到 `done` 或 `error` 事件时自动完成。

### 3. 统一事件协议（DFA 状态机）

定义在 types.ts，所有供应商的流被归一化为统一的事件序列：

```
start → [text_start → text_delta* → text_end]*
      → [thinking_start → thinking_delta* → thinking_end]*
      → [toolcall_start → toolcall_delta* → toolcall_end]*
      → done | error
```

| 事件类型         | 说明                                        |
| ---------------- | ------------------------------------------- |
| `start`          | 流开始，携带初始 partial message            |
| `text_delta`     | 文本增量片段                                |
| `thinking_delta` | 推理思考增量片段                            |
| `toolcall_delta` | 工具调用参数增量 JSON                       |
| `*_end`          | 某个内容块完成                              |
| `done`           | 整条消息完成，携带完整的 `AssistantMessage` |
| `error`          | 出错或中止                                  |

### 4. 各供应商适配（以 Anthropic 为例）

每个供应商将其私有 SSE 事件映射为统一事件。以 Anthropic provider 为例：

```
Anthropic SSE 事件                 →  统一事件
─────────────────────────────────────────────
message_start                      →  捕获 input token usage
content_block_start (text)         →  text_start
content_block_delta (text_delta)   →  text_delta
content_block_delta (thinking_delta) → thinking_delta
content_block_delta (input_json_delta) → toolcall_delta
content_block_stop                 →  text_end / thinking_end / toolcall_end
message_delta                      →  更新 usage + stopReason
stream 结束                        →  done
```

OpenAI Responses API 类似，通过 `processResponsesStream()` 做同样的归一化。

### 5. 调用流程总结

```
stream(model, context, options)
  │
  ├── 1. resolveApiProvider(model.api)    // 从注册表查找 Provider
  │
  ├── 2. new AssistantMessageEventStream() // 创建统一事件流
  │
  ├── 3. (async) {                        // 异步 IIFE 作为生产者
  │       创建供应商 SDK client
  │       发起带 stream: true 的 HTTP 请求 (SSE)
  │       for await (event of 供应商SSE流) {
  │         归一化为 AssistantMessageEvent
  │         stream.push(event)            // 推入统一流
  │       }
  │       stream.push({ type: 'done' })
  │       stream.end()
  │     }
  │
  └── 4. return stream                    // 立即返回，消费者异步迭代
```

### 6. 总结

| 层次             | 实现                                                                  |
| ---------------- | --------------------------------------------------------------------- |
| **传输协议**     | SSE（HTTP 流式响应），可选 WebSocket                                  |
| **流抽象**       | `EventStream` — 基于 Promise 的生产者-消费者异步迭代器                |
| **统一事件协议** | `AssistantMessageEvent` 联合类型，DFA 状态机式的事件序列              |
| **供应商适配**   | 每个 Provider 将私有 SSE 事件映射为统一格式推入 `EventStream`         |
| **消费方式**     | `for await (const event of stream)` 或 `stream.result()` 获取完整结果 |
