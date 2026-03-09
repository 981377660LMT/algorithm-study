# @mariozechner/pi-ai 代码功能全面讲解

> 结合 README.md 与 30 个测试用例，从架构到细节拆解这个统一 LLM API 库。

**九大章节覆盖全部功能**：

| 章节              | 内容                                                               |
| ----------------- | ------------------------------------------------------------------ |
| 一、项目定位      | 统一 20+ LLM 提供商的 API 层                                       |
| 二、核心架构      | 5 个关键模块：types / api-registry / stream / models / EventStream |
| 三、功能详解      | 15 项功能，每项都关联具体测试用例                                  |
| 四、提供商适配    | Google 系列(6个)、OpenAI 系列(3个)、Copilot、Zen 等的特殊处理      |
| 五、OpenAI 兼容层 | `compat` 字段如何解决各小 provider 的差异                          |
| 六、OAuth 认证    | 5 种 OAuth provider 的登录流程                                     |
| 七、测试架构      | 30 个测试分三类：E2E / Mock / 纯单元                               |
| 八、请求生命周期  | 从用户调用到返回结果的完整流程图                                   |
| 九、核心设计理念  | Provider 无关、Context 可序列化/迁移、流式优先、类型安全           |

重点讲解了跨 provider 消息变换（thinking → text、tool call ID 规范化）、上下文溢出检测（15+ 种正则匹配）、Unicode 安全处理等这些仅从 README 看不到的深层实现细节。

---

## 一、项目定位：一句话概括

**pi-ai 是一个「统一 LLM API 层」**——用同一套接口调 20+ 家 AI 提供商（OpenAI、Anthropic、Google、Bedrock、xAI、Groq……），自动处理：

- 模型发现与注册
- 流式/完整请求
- 工具调用（Function Calling）
- Token 计数与费用追踪
- 跨提供商上下文切换
- 上下文序列化/反序列化

---

## 二、核心架构（5 个关键模块）

```
┌─────────────────────────────────────────────────────┐
│                    用户代码                           │
│   stream() / complete() / streamSimple()            │
└───────────────┬─────────────────────────────────────┘
                │
┌───────────────▼──────────────────────────────────────┐
│              api-registry.ts                         │
│   注册表：api名 → { stream, streamSimple }           │
└───────────────┬──────────────────────────────────────┘
                │
┌───────────────▼──────────────────────────────────────┐
│           providers/                                  │
│  anthropic.ts  openai-responses.ts  google.ts  ...   │
│  每个 provider 实现 StreamFunction，处理 SSE 解析、   │
│  消息格式转换、错误映射等                              │
└───────────────┬──────────────────────────────────────┘
                │
┌───────────────▼──────────────────────────────────────┐
│           utils/                                      │
│  event-stream.ts  ← 生产者-消费者异步迭代              │
│  transform-messages.ts  ← 跨 provider 消息变换       │
│  overflow.ts  ← 上下文溢出检测                        │
│  validation.ts  ← AJV 工具参数校验                    │
│  sanitize-unicode.ts  ← Unicode 代理对清洗            │
│  json-parse.ts  ← 部分 JSON 流式解析                  │
└──────────────────────────────────────────────────────┘
```

### 2.1 types.ts — 类型系统

定义了整个库的核心类型：

| 类型                    | 说明                                                               |
| ----------------------- | ------------------------------------------------------------------ |
| `Model<TApi>`           | 模型描述：id、provider、api、费用、上下文窗口等                    |
| `Context`               | 对话上下文：systemPrompt + messages + tools                        |
| `Message`               | 三种角色：`UserMessage` / `AssistantMessage` / `ToolResultMessage` |
| `AssistantMessageEvent` | 流式事件 DFA：start → text/thinking/toolcall 系列 → done/error     |
| `Tool`                  | 工具定义：name + description + TypeBox schema                      |
| `StreamOptions`         | 请求选项：temperature、maxTokens、signal、apiKey、headers 等       |

### 2.2 api-registry.ts — 提供商注册表

```typescript
// 核心机制：一个 Map<apiName, {stream, streamSimple}>
const apiProviderRegistry = new Map<string, RegisteredApiProvider>()

// 注册：每个 provider 模块在 import 时自注册
registerApiProvider({ api: 'anthropic-messages', stream: streamAnthropic, ... })

// 查找：stream() 调用时按 model.api 查表
function resolveApiProvider(api) { return apiProviderRegistry.get(api) }
```

### 2.3 stream.ts — 统一入口（仅 60 行）

```typescript
// 4 个函数，本质都是：查注册表 → 委托给对应 provider
export function stream(model, context, options)     { ... }
export function complete(model, context, options)    { ... }  // = stream + .result()
export function streamSimple(model, context, options){ ... }  // 简化版，reasoning: 'medium'
export function completeSimple(...)                  { ... }
```

### 2.4 models.ts — 模型注册

```typescript
// 自动生成的 models.generated.ts 包含所有已知模型数据
// getModel('openai', 'gpt-4o-mini') → 带完整类型推导的 Model 对象
// getProviders() / getModels('anthropic') → 枚举可用提供商和模型
```

### 2.5 EventStream — 生产者-消费者异步迭代

```typescript
// 核心：支持 for-await-of 的事件流
class AssistantMessageEventStream extends EventStream<AssistantMessageEvent, AssistantMessage> {
  // push(event)  — provider 往里塞事件
  // for await (const event of stream) — 用户消费事件
  // await stream.result()  — 获取最终 AssistantMessage
}
```

---

## 三、功能详解 + 对应测试用例

### 3.1 流式文本生成

**功能**: 向 LLM 发送消息，逐 token 流式返回。

**事件流 DFA**:

```
start → text_start → text_delta* → text_end → done
```

**测试**: `stream.test.ts` → `basicTextGeneration`

- 发 "What is TypeScript?" 验证收到有效文本
- 追加 "What did I just ask?" 验证多轮对话记忆

```typescript
// 使用示例
const s = stream(model, { messages: [{ role: 'user', content: '你好' }] })
for await (const event of s) {
  if (event.type === 'text_delta') process.stdout.write(event.delta)
}
```

---

### 3.2 工具调用（Function Calling）

**功能**: LLM 请求调用外部工具 → 用户执行 → 返回结果给 LLM。

**事件流**:

```
start → toolcall_start → toolcall_delta* → toolcall_end → done(reason='toolUse')
```

**测试**: `stream.test.ts` → `handleToolCall`

- 定义 `calculator` 工具（expression 参数）
- 问 "What is 123 \* 456?" → 模型调用 calculator → 返回 56088
- 验证 toolcall_start/delta/end 事件完整性

**参数校验** — `validation.ts`:

```typescript
// AJV 自动校验 + 类型强制转换
const validatedArgs = validateToolCall(tools, toolCall)
// 校验失败 → 错误信息返回给模型让其重试
```

---

### 3.3 思考/推理（Thinking/Reasoning）

**功能**: 支持模型展示内部推理过程（Anthropic Thinking、OpenAI Reasoning、Gemini Thinking）。

**事件流**:

```
start → thinking_start → thinking_delta* → thinking_end → text_start → ...
```

**测试**:

- `stream.test.ts` → `handleThinking`: 验证 thinking 事件存在
- `interleaved-thinking.test.ts`: 两轮对话都有 thinking（Anthropic Opus 4.5/4.6）
- `xhigh.test.ts`: xhigh reasoning effort 的支持与拒绝行为
- `supports-xhigh.test.ts`: `supportsXhigh()` 模型能力检测

**简化接口** `streamSimple` — 自动映射 reasoning 级别:

```typescript
// 用户只需说 reasoning: 'medium'，库自动适配各 provider
await completeSimple(model, context, { reasoning: 'medium' })
// Anthropic → thinkingEnabled + budgetTokens
// OpenAI → reasoningEffort: 'medium'
// Google → thinking: { enabled: true, budgetTokens: 8192 }
```

---

### 3.4 图像输入（Vision）

**功能**: 向支持视觉的模型发送图片。

**测试**: `stream.test.ts` → `handleImage`

- 发送红色圆圈 PNG（test/data/red-circle.png）
- 验证模型回复包含 "red" 和 "circle"

```typescript
const msg = {
  role: 'user',
  content: [
    { type: 'text', text: '描述这张图片' },
    { type: 'image', data: base64Data, mimeType: 'image/png' }
  ]
}
```

---

### 3.5 图片工具结果

**功能**: 工具执行结果可以包含图片（给模型"看"）。

**测试**: `image-tool-result.test.ts`

- 纯图片工具结果：工具返回红色圆圈 PNG，模型能识别
- 文本+图片混合：工具同时返回描述文字和图片

```typescript
context.messages.push({
  role: 'toolResult',
  toolCallId: call.id,
  toolName: 'generate_chart',
  content: [
    { type: 'text', text: '直径 100 像素' },
    { type: 'image', data: base64, mimeType: 'image/png' }
  ],
  isError: false,
  timestamp: Date.now()
})
```

---

### 3.6 中止请求（Abort）

**功能**: 随时取消进行中的 LLM 请求，保留已接收的部分内容。

**测试**: `abort.test.ts`

- **流式中止**: 接收 50 字符后 abort → stopReason === 'aborted'，部分内容保留
- **立即中止**: 请求前就 abort → 立即返回 aborted
- **中止后继续**: abort 后的 context 加入新消息 → 可正常对话

```typescript
const controller = new AbortController()
setTimeout(() => controller.abort(), 2000)
const s = stream(model, context, { signal: controller.signal })
// 中止后 response.content 包含已接收的部分内容
```

**Token 统计** — `tokens.test.ts` 验证:

- Anthropic/Google: 中止后仍有 token 统计（早期发送 usage）
- OpenAI: 中止后无统计（usage 在最终 chunk 中）

---

### 3.7 跨提供商切换（Cross-Provider Handoff）

**功能**: 同一对话可以中途换模型/提供商，上下文自动适配。

**核心机制** — `transform-messages.ts`:

```
1. user 消息 → 原样传递
2. 同 provider 的 assistant 消息 → 保留 thinking blocks + signatures
3. 不同 provider 的 assistant 消息 → thinking 转为 <thinking> 标签的纯文本
4. tool call ID → 跨 provider 时规范化（去掉 OpenAI 的超长 ID）
5. redacted thinking → 只对同模型保留，跨模型直接丢弃
```

**测试**: `cross-provider-handoff.test.ts`

- 为每个 provider 生成含 thinking + tool call + tool result 的完整 context
- **两两交叉测试**: Anthropic ↔ OpenAI ↔ Google ↔ Bedrock ↔ xAI ↔ ...（27+ 组合）

**测试**: `transform-messages-copilot-openai-to-anthropic.test.ts`

- Copilot (OpenAI) → Anthropic 时 thinking 转 text
- 跨模型迁移时清除 `thoughtSignature`

**测试**: `openai-responses-reasoning-replay-e2e.test.ts`

- 中止后的 reasoning-only 记录处理
- 同 provider 但不同模型间切换
- 跨 provider 切换（Anthropic → OpenAI Codex）时工具 ID 兼容

---

### 3.8 工具调用 ID 规范化

**功能**: 不同 provider 生成的 tool call ID 格式不同，跨 provider 使用需要兼容处理。

**问题**: OpenAI Responses API 生成 450+ 字符的 `{call_id}|{base64}` 格式 ID，传给 Anthropic（要求 ≤64 字符 + `^[a-zA-Z0-9_-]+$`）会报错。

**测试**: `tool-call-id-normalization.test.ts`

- GitHub Copilot → OpenRouter 管道符 ID 规范化
- GitHub Copilot → OpenAI Codex ID 规范化
- 使用真实 issue #1022 的失败 ID 复现测试

---

### 3.9 Anthropic 工具名称规范化

**功能**: Anthropic OAuth（Claude Code）要求工具名 PascalCase，但用户定义可能是 snake_case。

**测试**: `anthropic-tool-name-normalization.test.ts`

- `todowrite` → 发给 API 作 `TodoWrite` → 返回时还原为 `todowrite`
- pi 内置工具 `read` → 不变成 `Read`
- 自定义工具 `my_custom_tool` → 直接透传

---

### 3.10 上下文溢出检测

**功能**: 发送超长内容时，自动检测各 provider 的溢出错误。

**核心** — `overflow.ts`:

```typescript
// 15+ 种正则模式匹配不同 provider 的错误信息
const OVERFLOW_PATTERNS = [
  /prompt is too long/i,           // Anthropic
  /exceeds the context window/i,   // OpenAI
  /input token count.*exceeds/i,   // Google
  /maximum prompt length is \d+/i, // xAI
  // ...
];

function isContextOverflow(message, contextWindow?) → boolean
```

**测试**: `context-overflow.test.ts`

- 构造超过 contextWindow + 10000 tokens 的消息
- 验证 `isContextOverflow()` 对每个 provider 都返回 true

---

### 3.11 空消息处理

**功能**: 优雅处理各种"空"输入，不崩溃。

**测试**: `empty.test.ts`

- 空数组 `[]`
- 空字符串 `""`
- 纯空白 `"   \n\t  "`
- 对话中的空 assistant 消息

---

### 3.12 缓存管理

**功能**: 通过 `PI_CACHE_RETENTION` 环境变量或选项控制 prompt 缓存时间。

**测试**: `cache-retention.test.ts`

- **Anthropic**: 默认 `ephemeral`（5 min）→ `long` 时 TTL 1h → `none` 无缓存
- **OpenAI Responses**: 默认无缓存 → `long` 时 `prompt_cache_retention: "24h"` → 代理 URL 不设

---

### 3.13 Token 与费用追踪

**功能**: 自动统计输入/输出/缓存 token 并计算费用。

```typescript
// 每个 AssistantMessage 自带：
message.usage = {
  input: 1234,  output: 567,
  cacheRead: 100,  cacheWrite: 50,
  totalTokens: 1951,
  cost: { input: 0.0012, output: 0.0085, total: 0.0097, ... }
}
```

**测试**: `total-tokens.test.ts`

- 验证 `totalTokens === input + output + cacheRead + cacheWrite`（所有 provider）

---

### 3.14 Unicode 安全处理

**功能**: 防止 Unicode 代理对（surrogate pair）导致崩溃。

**测试**: `unicode-surrogate.test.ts`

- emoji（🙈👍❤️🤔🚀）+ 多语言文本
- 真实世界 LinkedIn 德语评论数据（含大量 emoji）
- 故意构造的未配对高代理字符 `0xD83D`
- 验证不抛出 "no low surrogate in string" 错误

---

### 3.15 孤立工具调用处理

**功能**: 工具调用后用户取消/跳过，没有对应 toolResult → 自动过滤。

**测试**: `tool-call-without-result.test.ts`

- 助手发起 tool call → 用户直接发新消息（不提供 toolResult）
- 验证请求正常成功，不报错

---

## 四、提供商特定适配

### 4.1 Google 系列（6 个测试）

| 测试文件                                           | 关键行为                                                       |
| -------------------------------------------------- | -------------------------------------------------------------- |
| `google-gemini-cli-claude-thinking-header.test.ts` | Gemini CLI 调 Claude thinking 模型时加 `anthropic-beta` header |
| `google-gemini-cli-empty-stream.test.ts`           | 空 SSE 响应 → 自动重试（不重复 start/done 事件）               |
| `google-gemini-cli-retry-delay.test.ts`            | 解析 `Retry-After`、`x-ratelimit-reset` 等延迟 header          |
| `google-shared-gemini3-unsigned-tool-call.test.ts` | 无 `thoughtSignature` 的 tool call → 转为文本防止模型模仿      |
| `google-thinking-signature.test.ts`                | `thought === true` 才判定 thinking，签名保留逻辑               |
| `google-tool-call-missing-args.test.ts`            | 缺少 args 的 `functionCall` → 默认 `{}`                        |

### 4.2 OpenAI 系列

| 测试文件                                        | 关键行为                                                       |
| ----------------------------------------------- | -------------------------------------------------------------- |
| `openai-codex-stream.test.ts`                   | Codex SSE 流解析、session headers、gpt-5.3 minimal→low 钳制    |
| `openai-completions-tool-choice.test.ts`        | `tool_choice` 转发、strict 模式开关、Groq qwen3 reasoning 映射 |
| `openai-completions-tool-result-images.test.ts` | 工具结果中图片 → 打包进末尾 user 消息                          |

### 4.3 GitHub Copilot

| 测试                               | 关键行为                                           |
| ---------------------------------- | -------------------------------------------------- |
| `github-copilot-anthropic.test.ts` | Bearer auth + Copilot headers + thinking beta 控制 |

### 4.4 OpenCode Zen/Go

| 测试          | 关键行为                     |
| ------------- | ---------------------------- |
| `zen.test.ts` | 遍历所有 Zen/Go 模型冒烟测试 |

---

## 五、OpenAI 兼容层（重要设计）

许多小 provider（Groq、xAI、Cerebras、Mistral、MiniMax……）使用 OpenAI 兼容 API，但各有差异。
通过 `compat` 字段处理：

```typescript
interface OpenAICompletionsCompat {
  supportsStore?: boolean // LiteLLM 不支持 store
  supportsDeveloperRole?: boolean // 某些 provider 不支持 developer role
  requiresToolResultName?: boolean // 工具结果需不需要 name 字段
  requiresMistralToolIds?: boolean // Mistral 要求 ID 刚好 9 个字母数字
  thinkingFormat?: 'openai' | 'zai' | 'qwen' // 不同的推理参数格式
  maxTokensField?: 'max_completion_tokens' | 'max_tokens'
  // ...
}
```

库会根据 `baseUrl` 自动检测已知 provider 的 compat 设置，也支持手动覆盖。

---

## 六、OAuth 认证

支持 5 种 OAuth provider：

| Provider          | 说明                             |
| ----------------- | -------------------------------- |
| Anthropic         | Claude Pro/Max 订阅              |
| OpenAI Codex      | ChatGPT Plus/Pro 订阅            |
| GitHub Copilot    | Copilot 订阅                     |
| Google Gemini CLI | Cloud Code Assist                |
| Antigravity       | 免费 Gemini 3 / Claude / GPT-OSS |

```bash
npx @mariozechner/pi-ai login              # 交互式登录
npx @mariozechner/pi-ai login anthropic    # 指定 provider
```

---

## 七、测试架构总览

```
30 个测试文件
├── E2E 测试（需要真实 API Key）
│   ├── stream.test.ts        ← 最核心的多 provider 测试
│   ├── abort.test.ts         ← 中止信号处理
│   ├── cross-provider-handoff.test.ts  ← 跨 provider 上下文交换
│   ├── context-overflow.test.ts  ← 溢出检测
│   ├── empty.test.ts         ← 空输入处理
│   ├── tokens.test.ts        ← 中止时的 token 统计
│   ├── total-tokens.test.ts  ← totalTokens 正确性
│   ├── image-tool-result.test.ts  ← 图片工具结果
│   ├── interleaved-thinking.test.ts  ← 交错思考
│   ├── xhigh.test.ts         ← xhigh reasoning
│   ├── unicode-surrogate.test.ts  ← Unicode 安全
│   ├── tool-call-without-result.test.ts  ← 孤立工具调用
│   ├── tool-call-id-normalization.test.ts  ← ID 规范化
│   └── zen.test.ts / bedrock-models.test.ts  ← 冒烟测试
│
├── Mock 测试（不需要 API Key）
│   ├── cache-retention.test.ts       ← 缓存策略验证
│   ├── github-copilot-anthropic.test.ts  ← Copilot 请求构造
│   ├── openai-codex-stream.test.ts   ← Codex SSE 解析
│   ├── openai-completions-tool-choice.test.ts  ← 兼容性处理
│   ├── openai-completions-tool-result-images.test.ts  ← 图片转换
│   ├── google-gemini-cli-*.test.ts   ← Gemini CLI 系列
│   ├── google-shared-*.test.ts       ← Google 消息转换
│   ├── google-thinking-signature.test.ts  ← 签名检测
│   └── transform-messages-*.test.ts  ← 消息变换
│
└── 纯单元测试
    ├── supports-xhigh.test.ts  ← 模型能力检测
    └── anthropic-tool-name-normalization.test.ts  ← 名称映射
```

---

## 八、一图胜千言：请求生命周期

```
用户代码
  │
  ▼
stream(model, context, options)
  │
  ├─ 1. api-registry 查找 provider
  │
  ├─ 2. transform-messages: 跨 provider 消息适配
  │     ├─ thinking → text (跨 provider)
  │     ├─ tool call ID 规范化
  │     ├─ redacted thinking 丢弃 (跨模型)
  │     └─ 孤立 tool call 过滤
  │
  ├─ 3. provider 构造请求
  │     ├─ 设置 auth headers
  │     ├─ 转换消息格式 (Anthropic/OpenAI/Google 各不同)
  │     ├─ 应用 compat 兼容设置
  │     └─ 缓存策略 (cache_control / prompt_cache_retention)
  │
  ├─ 4. 发起 HTTP/SSE 请求
  │
  ├─ 5. 解析流式响应 → push 事件到 EventStream
  │     ├─ text_delta → 累积文本
  │     ├─ toolcall_delta → partial-json 增量解析
  │     ├─ thinking_delta → 累积思考
  │     └─ 错误/中止 → error 事件
  │
  ├─ 6. done 事件 → 计算 usage + cost
  │
  └─ 7. 返回 AssistantMessage
        ├─ content: TextContent[] | ThinkingContent[] | ToolCall[]
        ├─ usage: { input, output, cacheRead, cacheWrite, cost }
        ├─ stopReason: 'stop' | 'toolUse' | 'length' | 'error' | 'aborted'
        └─ api, provider, model, timestamp
```

---

## 九、核心设计理念

1. **Provider 无关**: 用户写一套代码，换个 `getModel()` 参数就切 provider
2. **Context 可序列化**: `JSON.stringify(context)` 即可持久化，跨服务传递
3. **Context 可迁移**: 同一上下文可从 Claude 切到 GPT-5 再切到 Gemini
4. **流式优先**: 所有接口都基于 `EventStream`，`complete()` 只是 `stream().result()` 的语法糖
5. **渐进式复杂度**: `completeSimple` → `complete` → `stream` → 直调 `streamAnthropic` 等
6. **类型安全**: TypeBox schema 定义工具参数，编译期 + 运行时双重校验
