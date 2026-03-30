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

---

这份 README.md 介绍了一个名为 `@mariozechner/pi-ai` 的工具库。它的核心目标是**统一管理各种 AI 模型（LLM）**，让你用一套代码就能调用 OpenAI、Claude、Gemini 等不同厂商的服务。

以下是对该文档的透彻拆解，分为核心理念、核心概念和开发流程三部分：

### 1. 核心理念：为什么要用它？

通常，每个 AI 厂商的 API 参数和返回格式都不一样（有的叫 `messages`，有的叫 `contents`；有的流式输出叫 `delta`，有的叫 `chunk`）。
`pi-ai` 做了三件事：

- **归一化 (Unified API)**：无论底层是哪个厂家，你调用的函数名（如 `complete`）和参数格式都是统一的。
- **自动发现 (Model Discovery)**：内置了各大主流模型的配置（价格、Token 限制、是否支持图片等）。
- **工具调用 (Tool Calling)**：强制要求模型必须支持“函数调用”功能，这对于构建 AI Agent（代理）至关重要。

---

### 2. 核心概念 (核心对象)

#### A. 模型 (Model)

通过 `getModel(provider, modelId)` 获取。

- **Provider (提供商)**: `openai`, `anthropic`, `google` 等。
- **Model ID**: 具体型号，如 `gpt-4o-mini` 或 `claude-3-5-sonnet`。
- **提示**：它提供了完整的 TypeScript 类型定义，你在编辑器里输入时会有自动补全。

#### B. 上下文 (Context)

这是对话的“存储箱”，包含：

- `systemPrompt`: AI 的角色设定（例如“你是一个资深算法专家”）。
- `messages`: 历史对白列表。
- `tools`: AI 能够调用的工具清单。

#### C. 工具 (Tools)

使用 **TypeBox** 定义。这非常聪明，因为它能自动校验 AI 返回的参数是否合法。

- 例如：你定义一个 `get_weather` 工具，要求传 `city` 字符串。如果 AI 传了个数字，库会报错并告诉 AI 重新传。

---

### 3. 如何使用 (开发流程)

文档展示了两种主要的交互方式：

#### 方式 1：普通模式 (`complete`)

简单直接，发送请求，等待 AI 全部写完后一次性返回。

```typescript
const response = await complete(model, context)
```

#### 方式 2：流式模式 (`stream`) — 推荐

像 ChatGPT 网页版那样，文字一个一个跳出来。通过监听不同的 `event.type` 来处理：

- `text_delta`: 收到一批新文字。
- `toolcall_delta`: AI 正在思考要传什么参数给工具。
- `done`: 生成结束，包含最终的消耗 (Usage) 和费用 (Cost)。

---

### 4. 特色高级功能

- **思考/推理过程 (Thinking/Reasoning)**:
  支持带“思维链”的模型（如 OpenAI 的 o1/o3，或 Claude 的思考模式）。你可以看到 AI 在给出答案前的“内心独白”。
- **多模态 (Image Input)**:
  可以识别图片。如果模型不支持（如老旧模型），库会自动忽略图片，不会导致程序奔溃。
- **全自动费用计算**:
  每个响应都会返回 `usage.cost`，帮你实时监控烧了多少钱。

### 总结建议

如果你是新手：

1.  **先跑通 README.md**：准备好一个 OpenAI 或 Anthropic 的 API Key。
2.  **关注类型安全**：借助 TS 的提示，你会发现写 AI 逻辑像写普通函数一样简单。
3.  **理解工具链**：重点看 `Tool` 的定义，这是把 LLM 变成“能干活的程序”的关键。

---

这份 README.md 是一个名为 `@mariozechner/pi-ai` 的 npm 库的官方文档。由于你是新手，我将采用“**原版内容提要 + 概念通俗解释 + 代码核心对照**”的方式，为你逐节透彻讲解，确保不遗漏任何重要内容。

---

### 1. 核心定位 (Title & Intro)

**原内容**：Unified LLM API with automatic model discovery, provider configuration, token and cost tracking, and simple context persistence and hand-off to other models mid-session.
**新手讲解**：
这是这个库存在的根本原因。目前市面上有大把的 AI 大模型（比如 ChatGPT、Claude、Gemini），但它们各自的调用格式、报错方式和计费统计都不一样。
这个库帮我们做了一件事：**统一接口 (Unified LLM API)**。

- 你只需要写一套代码，就能无缝切换各大公司的 AI。
- 它能自动计算你花了多少 Token（字数）和钱。
- 支持“中途换人”（比如让 Claude 思考，然后把结果无缝传给 GPT-4 接着聊）。
- **注**：它只收录支持“工具调用”(Tool Calling) 的模型，因为这对于做复杂的 AI 智能体 (Agent) 至关重要。

---

### 2. 支持的供应商 (Supported Providers)

**新手讲解**：
这里列出了一大串各大 AI 厂商，包括：

- 国外的巨头：OpenAI, Anthropic (Claude), Google, Azure (微软), Amazon 等。
- 开源及新兴势力：Mistral, Groq, xAI (马斯克的 Grok), MiniMax (国内的高维)。
- 集成平台：OpenRouter, 硅基流动等。
- 甚至支持本地部署的 AI（如 Ollama）。
  **意义**：只要你能搞到相关厂商的 API Key（密钥），用这个库就能随时调遣它们。

---

### 3. 安装与快速上手 (Installation & Quick Start)

**原内容**：`npm install @mariozechner/pi-ai` 和一大段代码。
**新手讲解**：
这是最核心的演示部分。

1.  **选定模型**：`getModel('openai', 'gpt-4o-mini')` —— 这是告诉库：“我要用 OpenAI 家的 gpt-4o-mini 模型”。
2.  **定义工具 (Tools)**：给 AI 配备“手和脚”。AI 本身不能联网看时间，所以在这里用 `TypeBox`（一个定义数据格式的工具）告诉 AI：“我这里有个叫 `get_time` 的函数，如果你需要查时间，你可以返回这个指令给我，我帮你查”。
3.  **构造上下文 (Context)**：相当于**聊天记录**。`systemPrompt` 是系统设定（比如“你是个助手”），`messages` 是用户说的话。
4.  **调用方式（两种）**：
    - **方法一：流式传输 (Stream)**。就像 ChatGPT 网页版一样，字是一个一个蹦出来的。文档里用 `switch (event.type)` 处理了各种情况：比如 `text_delta`（文字陆续输出）、`thinking_start`（AI 开始深度思考）等。
    - **执行工具回调**：当 AI 判定它需要查时间时，它会输出一个 `toolCall` 信号。代码中紧接着会真的去电脑上获取当前时间，然后把结果(`toolResult`)再塞回聊天记录（Context），并再次调用 `complete` 让 AI 根据时间总结出最终回答。
    - **方法二：完整返回 (Complete)**。不蹦字，等 AI 彻底想清楚了，一次性把整段回复交给你。

---

### 4. 工具调用 (Tools)

**新手讲解**：
工具调用让大模型不再是“瞎子和聋子”。

- **Defining Tools (定义工具)**：你需要严格定义工具的名字、描述和需要的参数，AI 才能理解怎么用。
- **Handling Tool Calls (处理调用)**：AI 请求调用工具，你执行代码后，必须用特定格式（带上 `toolCallId`）把结果还给 AI。
- **Streaming Tool Calls (流式接收工具参数)**：当工具参数很长时，这部分代码教你怎么在参数只接收到一半时（比如正在生成 JSON）就提前在界面上给用户做进度展示。
- **Validating Tool Arguments (验证参数)**：防止 AI “胡言乱语”乱传参数导致程序崩溃，收到指令必先校验。

---

### 5. 图片输入 (Image Input)

**新手讲解**：
如果你选的是支持视觉的模型（比如 GPT-4o 中带 Vision 的），你可以把本地图片转成 `base64`（一种把图片编码成字符串的文本格式），然后直接和文字一起放进 `Context` 传给 AI，AI 就能帮你“看图说话”。

---

### 6. 思考与推理 (Thinking/Reasoning)

**新手讲解**：
最近的新模型（如 DeepSeek-R1, OpenAI o1, Claude 3.7）都有“慢思考”能力，即先在内部生成一大段思考过程，再给出结论。

- 传统的库需要给不同的模型写不同的配置代码。这个库提供了统一接口（`streamSimple` / `completeSimple`），你只需要设置 `reasoning: 'medium'`，它就在底层自动帮你转换成各家厂商需要的对应参数格式。你可以单独抓取 `thinking_delta` 和 `text_delta`，把“思考过程”和“最终结果”分离开来显示。

---

### 7. 停止原因与错误处理 (Stop Reasons & Error Handling)

**新手讲解**：

- **Stop Reasons**：AI 为什么停下不说了？可能是正常说完 (`"stop"`)，或者是字数达到上限爆栈了 (`"length"`)，或者是想用工具等待你回复 (`"toolUse"`)。
- **Error Handling**：程序写飞了或者断网了怎么办。
- **Aborting (强行中止)**：给了一个打断 AI 的方法（利用 `AbortController`）。并且告诉你，即便是被打断说到一半的残缺句子，也可以存进聊天记录里，下次接着这段残缺的继续聊（Continuing After Abort）。

---

### 8. API、模型和提供商架构 (APIs, Models, and Providers)

**新手讲解**：
这部分是对库内部设计理念的解释。

- **架构关系**：Provider(微软/谷歌等) 提供 API 通信格式(比如 `openai-completions` 通用格式)，API 里跑着具体的 Model 模型。
- **Custom Models (自定义模型)**：手把手教你如何接入本地部署的 Ollama（开源模型），或者通过国内第三方代理转发的 API 服务。你只需要自己造一个字典（配置 URL 和计费标准），这个库依然能完美兼容。

---

### 9. 跨厂商上下文无缝切换 (Cross-Provider Handoffs)

**新手讲解（杀手级特性）**：
不同的 AI 脾气和格式不一样，把 Claude 的聊天记录直接喂给 GPT 经常会导致报错。
这个工具帮你做好了底层的“翻译”工作。你可以前一句通过 API 和 Claude 聊并产生了一段“思考过程”，下一句马上切换成 GPT 模型审视这句话，所有聊天记录、工具调用的历史都会被自动转换为新模型认识的兼容格式，中间不会报错断开。

---

### 10. 上下文序列化 (Context Serialization)

**新手讲解**：
AI 是没有记忆的！每次请求你都得把历史聊天打包发过去。
`Context`（聊天记录）在这里被设计成了标准的 JSON 格式，你可以随时用 `JSON.stringify` 把它转成字符串，存在本地硬盘、数据库或者浏览器的缓存中。下次打开程序，读出来接着聊。

---

### 11. 浏览器使用与环境变量 (Browser Usage & Environment Variables)

**新手讲解**：

- **在后端运行 (Node.js)**：推荐方式。你可以把密码 (API KEY) 写在电脑系统的环境变量里（比如 `OPENAI_API_KEY=xxx`），代码里不用写明文密码，库会自动感知、抓取并使用。安全省心。
- **在前端运行 (Browser)**：警告！严禁在网页前端写死你的 API Key 发送请求，哪怕再掩饰也很容易被黑客在这个页面 F12 抓包盗用你的 Key 刷额度。只能在内部演示用。

---

### 12. OAuth 授权登录 (OAuth Providers)

**新手讲解**：
有些大公司的服务不提供传统的“一串字母密钥”，而是需要弹出一个网页让你去点授权（比如登录 GitHub Copilot 或者 谷歌云）。这部分主要给高级用户展示如何在这个纯代码环境里拉起网页弹窗、获取 OAuth Token 操作。新手在使用常规的如 OpenAI, 硅基流动等服务时，不需要看这节。

### 13. 添加新供应商 (Development - Adding a New Provider)

**新手讲解**：
这是写给想给这个开源库提交代码的贡献者 (Contributor) 看的。罗列了如果要接入一家新的 AI 厂商，需要在核心配置、测试文件和文档说明等 8 个地方修改源码的清单。你只需当用户的，也可以跳过不看。

---

**总结**：这个库就是程序员调用各大 AI 提供的一个“万能插座+翻译官”，你了解它的核心逻辑就是：**准备语境 -> 发起调用（或流式监听） -> 配合完成工具指令 -> 拿到最终文本 -> 保存在历史记录。**
