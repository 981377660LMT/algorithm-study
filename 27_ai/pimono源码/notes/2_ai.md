# `@mariozechner/pi-ai` 深度源码分析

**pi-ai 的本质**：一个统一 LLM API 层，用三层洋葱模型（调用层 → 注册表 → Provider 适配器）屏蔽 10+ 家 API 差异。

**最重要的 11 个洞见**：

1. **注册表模式** — `registerApiProvider()` 在模块加载时自执行，支持运行时替换/扩展 provider
2. **AssistantMessage 携带产出元信息** — `provider`/`api`/`model` 字段使跨 provider 上下文传递成为可能，`transformMessages()` 据此决定是否丢弃签名、降级 thinking block
3. **每个流事件带 `partial` 快照** — 消费者不需要自己拼装状态，直接用 `partial.content` 渲染
4. **手写 EventStream** — 比 ReadableStream 更轻量，push/pull 模型无需背压（LLM 速率远低于内存分配速率）
5. **OAuth 模式伪装 Claude Code** — `user-agent: claude-cli/2.1.62`，利用 Anthropic 对 Claude Code 的特殊配额
6. **usage 在 `message_start` 就捕获** — abort 后仍保留 input token 计数
7. **配置优于继承** — `OpenAICompletionsCompat` 的 12 个 flag 解决 15+ 家兼容 API 的碎片化差异
8. **孤儿 tool call 修补** — 自动插入合成 toolResult，防止中断后的 API 错误
9. **15 个正则匹配上下文溢出** — 每家 provider 错误格式不同，z.ai 甚至默默接受溢出
10. **`StringEnum` helper** — 解决 Google 不支持 TypeBox `Type.Enum` 的 `anyOf/const` 模式
11. **跨平台条件加载** — Node API 动态 import，AJV 在 Chrome 扩展中降级，Web Crypto PKCE

**一句话**：pi-ai 的价值不在抽象有多漂亮，而在于对每个 provider 的脏活累活处理得多彻底。

> pi-ai 是一个统一 LLM API 层，用**一套类型系统**和**一个流式协议**屏蔽了 OpenAI / Anthropic / Google / Bedrock / GitHub Copilot 等 10+ 家不同 API 的差异。它是 pi-mono 代码编辑器的"AI 后端脊柱"。

---

## 1. 全局架构：三层洋葱模型

```
调用者 (stream / complete / streamSimple / completeSimple)
        │
        ▼
  ┌──────────────────────┐
  │  API Registry        │  ← 运行时注册表，key = Api 字符串
  │  (api-registry.ts)   │     value = { stream, streamSimple }
  └──────┬───────────────┘
         │ resolve
         ▼
  ┌──────────────────────────┐
  │  Provider 实现            │  ← anthropic.ts / openai-completions.ts / google.ts …
  │  (每个导出 stream +       │     各自处理鉴权、参数转换、SSE 解析
  │   streamSimple 函数)      │
  └──────┬───────────────────┘
         │ push events
         ▼
  ┌──────────────────────────┐
  │  EventStream<T, R>       │  ← 通用异步迭代器 + Promise<R>
  │  AssistantMessageEvent   │     统一的 13 种事件类型
  │  Stream                  │
  └──────────────────────────┘
```

**洞见 1：注册表模式解耦了"谁能调用"和"怎么实现"。**
`registerApiProvider()` 在模块加载时自动执行（register-builtins.ts 末尾的 `registerBuiltInApiProviders()`），
调用者只需要一个 `model.api` 字符串就够了。这意味着你可以在运行时：

- 替换某个 provider 的实现（测试 mock）
- 注册自定义 provider（比如 Ollama）
- 用 `unregisterApiProviders(sourceId)` 按来源批量卸载

---

## 2. 类型系统：最核心的设计决策

### 2.1 Message 三元组

```ts
type Message = UserMessage | AssistantMessage | ToolResultMessage
```

关键设计：

| 角色         | 内容类型                                         | 特殊字段                                          |
| ------------ | ------------------------------------------------ | ------------------------------------------------- |
| `user`       | `string \| (TextContent \| ImageContent)[]`      | 无                                                |
| `assistant`  | `(TextContent \| ThinkingContent \| ToolCall)[]` | `api`, `provider`, `model`, `usage`, `stopReason` |
| `toolResult` | `(TextContent \| ImageContent)[]`                | `toolCallId`, `toolName`, `isError`               |

**洞见 2：AssistantMessage 携带了"谁产生了它"的元信息。**
这不是多余的——`transformMessages()` 依赖 `assistantMsg.provider === model.provider` 来判断是否跨 provider 切换，
从而决定是否要：

- 丢弃 `thinkingSignature`（OpenAI 加密推理签名不能给 Anthropic 用）
- 丢弃 `redacted` thinking（只对同 model 有效）
- 规范化 `toolCallId`（不同 API 的 ID 格式天差地别）
- 把 `ThinkingContent` 降级为 `TextContent`（跨模型时保留思考内容但去掉签名）

### 2.2 AssistantMessageEvent：统一的 13 种流事件

```ts
type AssistantMessageEvent =
  | { type: 'start' }
  | { type: 'text_start' | 'text_delta' | 'text_end' }
  | { type: 'thinking_start' | 'thinking_delta' | 'thinking_end' }
  | { type: 'toolcall_start' | 'toolcall_delta' | 'toolcall_end' }
  | { type: 'done' }
  | { type: 'error' }
```

**洞见 3：每个事件都带 `partial: AssistantMessage`，这是一个"可查询的快照"。**
消费者不需要自己拼装状态——每次收到 delta，`partial` 已经是当前为止的完整 AssistantMessage。
这简化了 UI 层的渲染逻辑：只需 `partial.content` 即可渲染当前状态。

### 2.3 Model 与 Api 的泛型绑定

```ts
interface Model<TApi extends Api> {
  api: TApi
  compat?: TApi extends 'openai-completions'
    ? OpenAICompletionsCompat
    : TApi extends 'openai-responses'
      ? OpenAIResponsesCompat
      : never
}
```

`Model` 是泛型的，**类型参数约束了 `compat` 的形状**。这意味着在 TypeScript 层面：

- `getModel("openai", "gpt-4o")` 返回 `Model<"openai-completions">`，其 `compat` 类型是 `OpenAICompletionsCompat`
- `getModel("anthropic", "claude-sonnet-4-5")` 返回 `Model<"anthropic-messages">`，其 `compat` 是 `never`

---

## 3. EventStream：生产者-消费者的精巧实现

```ts
class EventStream<T, R = T> implements AsyncIterable<T> {
  private queue: T[] = []
  private waiting: ((value: IteratorResult<T>) => void)[] = []
  private done = false
  private finalResultPromise: Promise<R>
}
```

这是一个 **手写的异步队列**，同时支持两种消费模式：

1. **流式消费**：`for await (const event of stream) { ... }`
2. **等待最终结果**：`const msg = await stream.result()`

**设计精髓：**

- `push()` 时检查是否有 waiter——如果有，直接 resolve 它（零延迟传递）；否则入队
- `result()` 返回的 Promise 在 `push()` 检测到 `isComplete(event)` 时 resolve
- `end()` 关闭流并通知所有等待的消费者
- `done` 和 `error` 事件**既 push 给流迭代器，又 resolve `result()` promise**

**洞见 4：这比 Node.js `Readable` 或 Web `ReadableStream` 更轻量，且避免了背压复杂性。**
LLM 流的速率远低于内存分配速率，所以不需要背压。简单的 push/pull 模型就够了。

---

## 4. Provider 实现：以 Anthropic 为例的"适配器模式"

每个 provider 导出两个函数：

- `streamXxx`：暴露 provider 原生选项（如 `AnthropicOptions.thinkingEnabled`）
- `streamSimpleXxx`：只接收 `SimpleStreamOptions`（统一的 reasoning level），内部映射到原生选项

### 4.1 Anthropic 的关键实现细节

**a) 三种认证模式（createClient）**

```
1. GitHub Copilot → Bearer auth + 选择性 beta headers
2. OAuth Token (sk-ant-oat) → Bearer auth + Claude Code 伪装 headers
3. API Key → 标准 API key 认证
```

**洞见 5：OAuth 模式伪装成 Claude Code 客户端。**
`"user-agent": "claude-cli/2.1.62"` + `"x-app": "cli"` + Claude Code 工具名映射。
这不只是 cosmetic，而是利用了 Anthropic 对 Claude Code 的特殊配额/路由策略。
工具名 `Read`/`Write`/`Bash` 等在 OAuth 模式下会被映射到 Claude Code 的标准命名。

**b) 自适应思考 vs 预算思考**

```ts
if (supportsAdaptiveThinking(model.id)) {
  // Opus 4.6 / Sonnet 4.6: 模型自己决定何时思考、思考多少
  params.thinking = { type: 'adaptive' }
  params.output_config = { effort: 'high' }
} else {
  // 旧模型: 固定预算
  params.thinking = { type: 'enabled', budget_tokens: 1024 }
}
```

**c) 流事件的状态机处理**

```
message_start → 捕获初始 usage（input tokens）
content_block_start → 创建新的 text/thinking/toolCall block
content_block_delta → 增量追加内容
  └─ signature_delta → 追加思考签名（不是独立事件类型）
content_block_stop → 清理、发出 end 事件
message_delta → 更新 stopReason 和 usage
```

**洞见 6：usage 在 `message_start` 就捕获了 input tokens。**
即使流被中途 abort，也能保留 input token 计数。这对用量追踪至关重要。

**d) Tool Call 的流式 JSON 解析**

```ts
block.partialJson += event.delta.partial_json
block.arguments = parseStreamingJson(block.partialJson)
```

`parseStreamingJson` 先尝试 `JSON.parse`（完整 JSON 最快），失败后用 `partial-json` 库解析不完整的 JSON。
这使得工具调用的参数可以在流式传输中**渐进式解析**，UI 可以实时展示正在构建的参数。

### 4.2 OpenAI Completions 的兼容性迷宫

`OpenAICompletionsCompat` 有 **12 个兼容性选项**，覆盖了 15+ 个 OpenAI 兼容 API 的差异：

```ts
interface OpenAICompletionsCompat {
  supportsStore?: boolean // 有些 provider 不支持 store 字段
  supportsDeveloperRole?: boolean // system vs developer 角色
  supportsReasoningEffort?: boolean // 推理强度控制
  maxTokensField?: 'max_completion_tokens' | 'max_tokens'
  requiresToolResultName?: boolean // 工具结果是否需要 name 字段
  requiresAssistantAfterToolResult?: boolean // 工具结果后是否需要插入 assistant 消息
  requiresThinkingAsText?: boolean // thinking block → <thinking> 文本
  requiresMistralToolIds?: boolean // Mistral 要求刚好 9 个字符的 ID
  thinkingFormat?: 'openai' | 'zai' | 'qwen' // 三种不同的思考参数格式
  // ...
}
```

**洞见 7：兼容性不是通过继承或多态解决的，而是通过配置驱动。**
每个模型可以在 `model.compat` 中覆盖默认行为。默认值通过 URL 自动检测（如 `baseUrl.includes("api.mistral.ai")`）。
这是务实的工程决策——API 兼容性的差异太碎片化了，继承层次会爆炸。

---

## 5. transformMessages：跨 Provider 上下文传递的核心

这是整个库最关键的函数之一。当用户在对话中切换模型时，历史消息需要被转换：

```ts
function transformMessages(messages, model, normalizeToolCallId?) {
  // 第一遍：转换内容块
  //   - 同模型：保留 thinking signatures（用于推理回放）
  //   - 跨模型：thinking → text（保留内容、丢弃签名）
  //   - 跨模型：规范化 toolCallId
  //   - 丢弃 redacted thinking（仅同模型有效）
  //   - 丢弃空 thinking blocks
  //   - 丢弃 error/aborted 的 assistant messages
  // 第二遍：修补孤儿 tool calls
  //   - 如果 assistant 消息有 toolCall 但没有对应的 toolResult
  //   - 插入合成的 "No result provided" toolResult
  //   - 这满足了每个 API 的"每个 toolCall 必须有 toolResult"约束
}
```

**洞见 8：孤儿 tool call 修补是一个极其实际的 edge case。**
场景：用户在工具执行中中断、模型 abort、或者手动编辑了上下文。
没有这个修补，下一次调用会因为缺少 toolResult 而报 API 错误。

---

## 6. 上下文溢出检测：防御性编程的典范

`isContextOverflow()` 用 **15 个正则模式** 匹配各家 provider 的溢出错误消息：

```ts
const OVERFLOW_PATTERNS = [
  /prompt is too long/i, // Anthropic
  /exceeds the context window/i, // OpenAI
  /input token count.*exceeds/i, // Google
  /maximum prompt length is \d+/i, // xAI
  /reduce the length of the messages/i // Groq
  // ... 更多
]
```

而且处理了两个特殊情况：

- **Cerebras/Mistral**：返回 400/413 但没有 body → 通过状态码匹配
- **z.ai**：默默接受溢出请求 → 通过 `usage.input > contextWindow` 检测

**洞见 9：这是"现实世界的 API 没有标准化"的缩影。**
每家 provider 的错误格式都不同，甚至同一家的错误格式也会变。这种模式匹配列表需要持续维护。

---

## 7. 工具参数验证：AJV + TypeBox

```ts
const calculatorSchema = Type.Object({
  a: Type.Number({ description: 'First number' }),
  operation: StringEnum(['add', 'subtract', 'multiply', 'divide'])
})
```

- **TypeBox** 生成 JSON Schema（LLM 需要）的同时提供 TypeScript 类型推导
- **AJV** 在运行时验证 LLM 返回的工具参数
- **coerceTypes** 开启：LLM 返回字符串 `"42"`，AJV 自动转为数字 `42`

**洞见 10：`StringEnum` helper 的存在揭示了 Google API 的兼容性问题。**
Google 不支持 TypeBox `Type.Enum` 生成的 `anyOf/const` 模式，只接受 `{ type: "string", enum: [...] }`。
一个小 helper 函数解决了一个实际的跨 provider 兼容性问题。

**浏览器扩展兼容**：在 Chrome Manifest V3 扩展中，CSP 禁止 `eval` / `Function` 构造器，
所以 AJV 初始化会失败。代码优雅地降级为"信任 LLM 输出，跳过验证"。

---

## 8. OAuth 认证体系

```
OAuthProviderInterface {
  login(callbacks) → credentials   // PKCE 流程 + 本地回调服务器
  refreshToken(credentials)         // 刷新过期 token
  getApiKey(credentials)            // 把 credentials 转为 API key 字符串
  modifyModels?(models, creds)      // 可选：修改模型列表(如更新 baseUrl)
}
```

支持 5 种 OAuth 提供者：

1. **Anthropic** → Claude Pro/Max 订阅
2. **GitHub Copilot** → 设备码流程
3. **Google Gemini CLI** → Google Cloud OAuth
4. **Google Antigravity** → Google Cloud 代码助手
5. **OpenAI Codex** → ChatGPT OAuth

PKCE（`pkce.ts`）用 Web Crypto API 实现，确保 Node.js 和浏览器都能用。

---

## 9. 值得注意的工程实践

### 9.1 Unicode Surrogate 净化

```ts
export function sanitizeSurrogates(text: string): string {
  return text.replace(/[\uD800-\uDBFF](?![\uDC00-\uDFFF])|(?<![\uD800-\uDBFF])[\uDC00-\uDFFF]/g, '')
}
```

每次发送文本到 API 前都调用。来源：工具结果中可能包含 emoji 等 BMP 之外的字符，
如果不正确序列化就会产生未配对的代理对，导致 JSON 解析失败。

### 9.2 HTTP Proxy 支持

```ts
// http-proxy.ts - 仅 4 行代码
import('undici').then(({ EnvHttpProxyAgent, setGlobalDispatcher }) => {
  setGlobalDispatcher(new EnvHttpProxyAgent())
})
```

在 `stream.ts` 顶部 `import "./utils/http-proxy.js"` 确保在任何 API 调用前设置代理。
通过 ES Module 缓存保证只执行一次。

### 9.3 环境检测与条件加载

```ts
// env-api-keys.ts
// 动态 import node:fs / node:os / node:path，避免破坏浏览器/Vite 构建
if (typeof process !== 'undefined' && (process.versions?.node || process.versions?.bun)) {
  import('node:fs').then(m => {
    _existsSync = m.existsSync
  })
  // ...
}
```

**这个库同时支持 Node.js、Bun 和浏览器环境。** 所有 Node.js 特有 API 都是条件加载的。

### 9.4 缓存策略

Anthropic 支持三级缓存：

- `"none"`：不缓存
- `"short"`：临时缓存（默认）
- `"long"`：1 小时 TTL（仅 api.anthropic.com 支持）

缓存标记加在 system prompt 和最后一个 user message 上，用 `cache_control: { type: "ephemeral" }` 实现。

---

## 10. 测试设计的洞见

### Cross-Provider Handoff Test

最有价值的测试：为每个 provider/model 生成一个包含 thinking + tool call + tool result 的上下文，
然后**把所有其他 provider 的上下文拼在一起**发给目标 model。如果失败，说明有兼容性问题。

### Context Overflow Test

系统性测试每个 provider 的溢出行为，确保 `isContextOverflow()` 对所有 provider 都返回 true。

### Unicode Surrogate Test

模拟工具结果中包含 emoji 和未配对代理对的场景。

---

## 11. 架构理念总结

| 理念             | 体现                                                                 |
| ---------------- | -------------------------------------------------------------------- |
| **配置优于继承** | `OpenAICompletionsCompat` 12 个 flag 而非 12 个子类                  |
| **渐进式降级**   | AJV 在浏览器扩展中降级、proxy 仅在 Node 中启用                       |
| **现实主义**     | 15 个正则匹配溢出错误、孤儿 tool call 修补、Claude Code 伪装         |
| **流优先**       | EventStream 是一等公民，complete() 只是 `stream().result()` 的语法糖 |
| **可观测性**     | 每个事件带 `partial` 快照、`onPayload` 回调审查原始请求              |
| **跨平台**       | 条件加载 Node API、Web Crypto PKCE、CSP 兼容                         |

**一句话总结：pi-ai 是一个"在 LLM API 碎片化现实中追求统一性"的务实工程作品。
它的价值不在于抽象多漂亮，而在于对每个 provider 的脏活累活处理得多彻底。**
