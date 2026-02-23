# 再聊无状态 Agent：系统提示词到底该不该存？

你好！看起来你正在深入研究 Coding Agent 的架构设计，特别是关于 **System Prompt（系统提示词）的状态管理** 这一关键决策。

这篇文章的核心观点非常有价值：**将 Agent 设计为无状态（Stateless），并利用 Prompt Cache 解决成本问题。** 这是一个从简单的 Demo 走向生产级应用的关键转变。

我将按照你的要求，基于这篇文章的内容，结合实际工程视角，为你深度拆解这一设计模式。我们不仅要看“怎么做”，还要看“为什么要这样做”以及“潜在的深坑”。

### 核心论点拆解：为什么要“抛弃”有状态设计？

很多初学者容易写出“有状态”的 Agent，原因在于符合直觉：只有一份系统提示词，初始化一次，一直用到死。

#### 1. 有状态设计的致命缺陷 (The "Stateful" Trap)

正文中提到了一段典型的有状态代码：

```typescript
class Agent {
  private systemPrompt?: string // 状态存储在实例里
  history: Message[] = [] // 历史也存储在实例里
  // ...
  async chat() {
    // 只有第一次才把 systemPrompt 加入 history
    if (this.history.length === 0) this.history.push(system)
    // ...
  }
}
```

**深度解读：**

- **上下文污染 (Context Pollution):** 当对话变长（例如超过 10 轮），早期的 `system` 消息在 Token 序列中距离当前生成位置越来越远。LLM（特别是某些 Attention 机制实现）对开头内容的注意力可能会衰减，导致 Agent "忘记" 它的角色，开始胡言乱语。
- **灵活性丧失 (Lost Flexibility):** 如果你在第 5 轮对话时，想让 Agent 从 "写代码模式" 切换到 "代码审查模式"，你无法做到。因为 `system` 消息已经在第 1 轮被写死在 `history` 数组里了。你必须清空历史重开，这意味着上下文丢失。
- **序列化噩梦 (Serialization Nightmare):** 如果你要把 Agent 的状态存入 Redis 以便下次继续，你得保存整个 `history` 数组。如果 System Prompt 很大（例如 5k tokens），每次存取都要带着这坨巨大的静态文本，浪费带宽和存储。

#### 2. 无状态设计的架构变革 (The "Stateless" Revolution)

文章提出的方案是将 **数据（State）** 与 **行为（Behavior）** 分离。

**架构图解：**

- **Agent (Behavior):** 变成了一个纯函数。`Output = f(SystemPrompt, History, UserInput)`。它不持有任何数据，只负责执行。
- **Context (State):** 每次请求时动态组装。

```typescript
// 伪代码：纯函数式思维
const context = {
  systemPrompt: loadPrompt('coding_mode'), // 动态加载，随时可换
  history: loadHistoryFromDB(sessionId), // 纯粹的对话记录
  userInput: 'Fix this bug'
}
agent.execute(context)
```

**深度优势：**

- **动态调整 (Dynamic Injection):** 第 1 轮用 "规划专家" 的 Prompt，第 2 轮用 "Python 专家" 的 Prompt，历史记录 `history` 是共享的，但 "大脑" (System Prompt) 可以随时换。
- **上下文刷新 (Context Refresh):** 每次请求，System Prompt 都被放在 Token 序列的最前面（如果是 Anthropic 模型，甚至有专用字段）。这意味着模型每一轮都能接收到**最新、最完整、未被稀释**的指令。

---

### 关键技术点：Prompt Cache (成本与性能的救星)

`这是让无状态设计在经济上可行的核心技术。没有它，每次发 10k system prompt 会让你破产。`

#### 1. 缓存的工作原理

LLM 的推理本质是计算 Next Token Probabilities。对于固定的前缀文本（Prefix），其计算出的 KV Cache (Key-Value pairs) 是完全一样的。

- **无缓存：**这也是大多数 API 的默认行为。每次都要从头计算 System Prompt 的 KV Cache。
- **有缓存：** 服务端保留了 System Prompt 的计算结果。当你的请求头带有 `cacheControl: { type: 'ephemeral' }` 且文本匹配时，它可以直接跳过这部分的计算。

#### 2. 代码实现细节 (重点)

在实现 `AnthropicChatService` 时，以下细节容易出错：

**构建可缓存的消息结构：**

```typescript
// 假设你有一个这样的转换函数
import { Message } from '../types' // 假设类型定义在此

function convertToSDKMessages(messages: Message[], systemPromptContent: string) {
  const sdkMessages = []

  // 1. 系统提示词作为独立对象，并打上缓存标记
  // 注意：Anthropic SDK 最新版通常将 system 作为顶层参数，而不是 messages 数组的一部分
  // 但这里为了演示 Cache Control 的逻辑，我们看 API 结构

  // 正确的做法通常是在 API 调用参数中构造 system
  const systemBlock = [
    {
      type: 'text',
      text: systemPromptContent,
      cache_control: { type: 'ephemeral' } // 关键：打标
    }
  ]

  // 2. 转换历史消息，不需要 cache_control (除非你想缓存长对话的中间点)
  const historyBlocks = messages.map(msg => ({
    role: msg.role,
    content: msg.content
  }))

  return { system: systemBlock, messages: historyBlocks }
}
```

**注意事项：**

- **最少 Token 数限制：** Anthropic 的 Cache 通常要求缓存块至少包含 1024 个 Token。如果你的 System Prompt 只有 500 个 Token，缓存可能不会生效（或者不划算）。
- **缓存生命周期：** `ephemeral` 通常只有 5-10 分钟。这意味着如果用户发呆了 15 分钟再发下一条消息，缓存可能失效，你会再次支付写入成本。如果你有极高的并发，这很有用；如果是低频对话，成本优势可能不明显。

---

### 落地方案：你需要修改的代码

基于文章提到的改动清单，我建议你在当前工作区按照以下步骤进行重构：

#### 1. 修改 `ChatContext` 定义

这是基础，确保 context 能够承载动态的 system prompt。

```typescript
export interface Message {
  role: 'user' | 'assistant' | 'system' // 注意：system 可能不再需要存在于这个数组里，或者仅作为显示用
  content: string
  // ...
}

export interface ChatContext {
  messages: Message[] // 纯对话历史 (User + Assistant)
  workspaceRoot: string
  sessionId: string
  // 新增：动态注入的系统提示词
  systemPrompt?: string
  // ...
}
```

#### 2. 重构 Agent 类

剥离状态，使其“纯净化”。

```typescript
import { ChatContext } from '../types'
// ... imports

export class Agent {
  // 删除 private systemPrompt: string;
  // 删除 history: Message[];
  // Agent 不再持有这些状

  constructor(private llm: LLMInterface) {}

  async chat(context: ChatContext) {
    // 1. 获取或生成 System Prompt
    // 如果 context 里没传，就在这里现生成，但不再存到 this 上
    const finalSystemPrompt = context.systemPrompt || (await this.defaultSystemPromptBuilder())

    // 2. 拼接环境上下文 (Environment Context)
    const envInfo = `Current File: ${context.currentFile}\nOS: Mac`
    const fullSystemContent = `${finalSystemPrompt}\n\nEnvironment:\n${envInfo}`

    // 3. 构造请求给 LLM
    // 注意：这里不需要把 system push 到 context.messages 里
    // 而是传给 LLM 接口，让 LLM 接口决定怎么组装 (比如 Anthropic 放在 system 参数里)
    const response = await this.llm.chat({
      system: fullSystemContent,
      messages: context.messages
    })

    return response
  }

  private async defaultSystemPromptBuilder() {
    return 'You are a helpful coding assistant...'
  }
}
```

#### 3. 适配 LLM 服务层 (支持 Cache)

这是为了省钱的关键一步。

```typescript
// ... 假设这是你的 LLM 适配器

export class AnthropicLLM {
  async chat(params: { system: string; messages: Message[] }) {
    const { system, messages } = params

    // 调用 Anthropic SDK
    const response = await this.client.messages.create({
      model: 'claude-3-5-sonnet-20240620',
      max_tokens: 4096,
      system: [
        {
          type: 'text',
          text: system,
          // 开启缓存！
          cache_control: { type: 'ephemeral' }
        }
      ],
      messages: messages.map(m => ({
        role: m.role as any,
        content: m.content
      }))
    })

    // 记录一下 token 使用情况，看看 cache 是否命中
    console.log(`Cache Read: ${response.usage.cache_read_input_tokens}`)
    console.log(`Cache Create: ${response.usage.cache_creation_input_tokens}`)

    return response.content[0].text
  }
}
```

### 总结

这篇文章提出的**“System Prompt 动态注入 + Prompt Cache”** 方案，是目前构建高级 Coding Agent 的最佳实践。

它解决了：

1.  **一致性问题：** 永远使用最新的 System Prompt。
2.  **上下文污染问题：** System Prompt 永远在最前，权重最高。
3.  **成本问题：** 通过 Cache 复用重复的文本输入。
4.  **架构耦合问题：** Agent 变成无状态执行器，利于扩展。

---

这套模式的本质是关注点分离：

- Agent：纯粹的执行引擎，无状态
- systemPrompt：通过 context 动态注入
- History：外部存储，持久化

这与微服务架构中的 Stateless Service + External State Store 设计思路完全一致。
