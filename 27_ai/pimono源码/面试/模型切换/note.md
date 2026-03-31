# Pi-Mono 模型切换深度分析

## 一、架构全景

模型切换涉及 4 个核心包的协作：

```
┌─────────────────────────────────────────────────────────────┐
│                     coding-agent (TUI层)                     │
│  InteractiveMode ──→ AgentSession ──→ ModelRegistry          │
│  (Ctrl+P/L UI)       (状态协调)        (发现+鉴权)            │
├─────────────────────────────────────────────────────────────┤
│                     agent (核心Agent)                        │
│  Agent.setModel() → agent-loop → streamFn → LLM             │
├─────────────────────────────────────────────────────────────┤
│                     ai (LLM抽象层)                           │
│  stream() → ApiProviderRegistry → 具体 Provider 实现          │
│  transformMessages() → 跨 Provider 消息格式转换               │
├─────────────────────────────────────────────────────────────┤
│                     tui (终端UI)                             │
│  ModelSelectorComponent (fuzzy search列表)                    │
│  ScopedModelsSelectorComponent (启用/禁用管理)                │
└─────────────────────────────────────────────────────────────┘
```

---

## 二、类型系统

### Model<TApi> — 模型的完整描述

```ts
// packages/ai/src/types.ts
interface Model<TApi extends Api> {
  id: string // "claude-opus-4-6"
  name: string // 显示名
  api: TApi // ★ 决定路由到哪个 Provider 实现
  provider: Provider // "anthropic" — 标识提供商（用于 auth 查找）
  baseUrl: string // API 端点
  reasoning: boolean // 是否支持推理/thinking
  input: ('text' | 'image')[]
  cost: { input; output; cacheRead; cacheWrite }
  contextWindow: number
  maxTokens: number
  headers?: Record<string, string>
  compat?: OpenAICompletionsCompat | OpenAIResponsesCompat // 兼容性开关
}
```

### Api — 枚举所有支持的协议

```ts
type Api =
  | 'openai-completions' // ChatGPT 风格
  | 'openai-responses' // GPT-5 新API
  | 'azure-openai-responses'
  | 'openai-codex-responses' // Codex 专用
  | 'anthropic-messages' // Claude
  | 'bedrock-converse-stream' // AWS Bedrock
  | 'google-generative-ai' // Gemini
  | 'google-gemini-cli'
  | 'google-vertex'
```

**核心设计**：`model.api` 是路由键，决定调用哪个 streaming 实现。`model.provider` 是鉴权键，决定用哪个 API key。

---

## 三、模型注册与发现

### 3.1 内置模型注册

```ts
// packages/ai/src/models.ts
// 模块加载时从 MODELS（自动生成的巨大常量）初始化
const modelRegistry: Map<string, Map<string, Model<Api>>> = new Map()
for (const [provider, models] of Object.entries(MODELS)) {
  // provider → { modelId → Model }
}
```

### 3.2 API Provider 注册（策略模式）

```ts
// packages/ai/src/api-registry.ts
const apiProviderRegistry = new Map<string, RegisteredApiProvider>()

// 注册内置 providers（模块加载时自动执行）
// packages/ai/src/providers/register-builtins.ts
registerApiProvider({ api: 'anthropic-messages', stream: streamAnthropic, streamSimple: streamSimpleAnthropic })
registerApiProvider({ api: 'openai-completions', stream: streamOpenAICompletions, ... })
registerApiProvider({ api: 'openai-responses', stream: streamOpenAIResponses, ... })
registerApiProvider({ api: 'google-generative-ai', stream: streamGoogle, ... })
// ... 共 9 种协议
```

### 3.3 ModelRegistry — coding-agent 层的聚合管理器

```ts
// packages/coding-agent/src/core/model-registry.ts
class ModelRegistry {
  private models: Model<Api>[] = []

  constructor(authStorage, modelsJsonPath?) {
    this.loadModels()
  }

  private loadModels() {
    // 1. 从 models.json 加载自定义模型 + provider override
    const { models: customModels, overrides, modelOverrides } = this.loadCustomModels(path)
    // 2. 加载内置模型，应用 override (baseUrl/headers/apiKey/modelOverride)
    const builtInModels = this.loadBuiltInModels(overrides, modelOverrides)
    // 3. 合并（custom 同 provider+id 覆盖 built-in）
    let combined = this.mergeCustomModels(builtInModels, customModels)
    // 4. 让 OAuth providers 修改模型（如动态 baseUrl）
    for (const oauthProvider of this.authStorage.getOAuthProviders()) {
      combined = oauthProvider.modifyModels(combined, cred)
    }
    this.models = combined
  }

  // ★ 只返回有 auth 的模型（用于 Ctrl+P 循环）
  getAvailable(): Model<Api>[] {
    return this.models.filter(m => this.authStorage.hasAuth(m.provider))
  }

  // Extensions 可动态注册自定义 Provider
  registerProvider(name, config) { ... }
}
```

**加载流程图**：

```
models.generated.ts (编译时生成, 200+ 内置模型)
        ↓
getProviders() → getModels(provider)
        ↓
   models.json (~/.pi/models.json, 用户自定义)
        ↓ merge
   OAuth providers modifyModels()
        ↓ filter(hasAuth)
   getAvailable() → 可用模型列表
```

---

## 四、模型解析（Resolver）

### 4.1 初始模型选择 — 5 级优先级

```ts
// packages/coding-agent/src/core/model-resolver.ts
async function findInitialModel(options): Promise<InitialModelResult> {
  // 1. CLI 参数 --provider + --model（最高优先级）
  if (cliProvider && cliModel) → resolveCliModel()

  // 2. --models 标志的第一个模型（如果不是 resume/continue）
  if (scopedModels.length > 0 && !isContinuing) → scopedModels[0]

  // 3. settings.json 保存的默认模型
  if (defaultProvider && defaultModelId) → modelRegistry.find()

  // 4. 已知 provider 的默认模型（硬编码表）
  for (const provider of knownProviders) → defaultModelPerProvider[provider]

  // 5. 第一个可用模型
  return availableModels[0]
}
```

### 4.2 CLI模型解析 — 模糊匹配算法

```ts
function resolveCliModel({ cliProvider, cliModel, modelRegistry }) {
  // 步骤1：解析 "provider/model" 格式
  // 步骤2：精确匹配 model.id（case-insensitive）
  // 步骤3：模糊匹配（id 或 name 包含 pattern）
  //   - 优先 alias（如 claude-sonnet-4-5）而非带日期的版本
  //   - 多个 alias 取字母序最高的
  // 步骤4：解析 thinking level 后缀（如 "sonnet:high"）
  //   - 处理 model ID 中含冒号的情况（如 OpenRouter 的 model:exacto）
  // 步骤5：fallback 构建自定义模型对象
}
```

### 4.3 Scope 解析 — 支持 Glob 和 thinking level

```ts
async function resolveModelScope(patterns: string[], modelRegistry): Promise<ScopedModel[]> {
  for (const pattern of patterns) {
    if (含有 glob 字符 *, ?, [) {
      // minimatch 匹配 "provider/modelId" 或 "modelId"
      // 支持 "*sonnet*:high" 格式
    } else {
      parseModelPattern(pattern, availableModels)
      // 递归解析: 先尝试整个 pattern，失败则在最后一个冒号处拆分
    }
  }
}
```

---

## 五、运行时模型切换（核心！）

### 5.1 快捷键绑定

```ts
// packages/coding-agent/src/core/keybindings.ts
const DEFAULT_APP_KEYBINDINGS = {
  cycleModelForward: 'ctrl+p', // 向前循环
  cycleModelBackward: 'shift+ctrl+p', // 向后循环
  selectModel: 'ctrl+l' // 打开搜索选择器
}
```

### 5.2 UI 层事件注册

```ts
// packages/coding-agent/src/modes/interactive/interactive-mode.ts
this.defaultEditor.onAction('cycleModelForward', () => this.cycleModel('forward'))
this.defaultEditor.onAction('cycleModelBackward', () => this.cycleModel('backward'))
this.defaultEditor.onAction('selectModel', () => this.showModelSelector())
```

### 5.3 Ctrl+P 循环切换 — 完整调用链

```
用户按 Ctrl+P
  │
  ▼
InteractiveMode.cycleModel("forward")
  │
  ▼
AgentSession.cycleModel("forward")
  │
  ├── if scopedModels.length > 0 → _cycleScopedModel()
  │     │
  │     ▼
  │   getScopedModelsWithApiKey()     // 过滤有 API key 的 scoped models
  │   if (只有1个) return undefined   // 无法循环
  │   找到当前模型在列表中的索引
  │   nextIndex = (currentIndex + 1) % length
  │   agent.setModel(next.model)      // ★ 更新 Agent 内部 _state.model
  │   sessionManager.appendModelChange()  // 写入 session 持久化
  │   settingsManager.setDefaultModel()   // 更新全局默认
  │   setThinkingLevel(next.thinkingLevel) // scoped model 可以带 thinking level
  │   _emitModelSelect(next, prev, 'cycle') // 通知 extensions
  │
  └── else → _cycleAvailableModel()
        │
        ▼
      modelRegistry.getAvailable()    // 所有有 auth 的模型
      同样的 modulo 循环逻辑
      modelRegistry.getApiKey(next)   // 校验 API key 有效
      agent.setModel(nextModel)
      setThinkingLevel(currentLevel)  // 保持当前 thinking level，re-clamp
  │
  ▼
返回 ModelCycleResult { model, thinkingLevel, isScoped }
  │
  ▼
InteractiveMode:
  footer.invalidate()             // 更新底栏模型名显示
  updateEditorBorderColor()       // 更新边框颜色（不同 provider 不同颜色）
  showStatus("Switched to ...")   // 状态提示
```

### 5.4 Ctrl+L 选择器 — 搜索 + 可视列表

```ts
// InteractiveMode.showModelSelector()
private showModelSelector(initialSearchInput?: string) {
  // 创建 ModelSelectorComponent overlay
  new ModelSelectorComponent(
    tui,
    currentModel,           // 当前模型（标 ✓）
    settingsManager,
    modelRegistry,
    scopedModels,           // 受限范围
    async (model) => {      // onSelect
      await this.session.setModel(model)   // ★ 和 cycleModel 不同，用 setModel
      footer.invalidate()
      showStatus(`Model: ${model.id}`)
    },
    () => cancel,
    initialSearchInput
  )
}
```

**ModelSelectorComponent 内部**：

```
┌─────────────────────────────────────────┐
│ ──────── border ────────                │
│                                         │
│ Scope: all | scoped    (Tab切换)        │
│ Only showing models with configured ... │
│                                         │
│ [搜索框: fuzzy filter]                  │
│                                         │
│ → claude-opus-4-6        [anthropic] ✓  │ ← 当前模型
│   claude-sonnet-4-5      [anthropic]    │
│   gpt-5.1-codex          [openai]      │
│   gemini-2.5-pro         [google]      │
│               (4/32)                    │
│                                         │
│   Model Name: Claude Opus 4.6           │
│ ──────── border ────────                │
└─────────────────────────────────────────┘
```

功能：

- **Tab 键**：切换 all / scoped 范围
- **上下箭头**：选择，wrap around
- **Enter**：确认选择，自动保存为默认模型
- **模糊搜索**：`fuzzyFilter(models, query, ({id, provider}) => ...)`
- **当前模型排最前**，按 provider 分组排序

### 5.5 setModel vs cycleModel 的区别

|                 | `setModel(model)`                      | `cycleModel(dir)`                                    |
| --------------- | -------------------------------------- | ---------------------------------------------------- |
| 调用场景        | Ctrl+L 选择器、Extensions、/model 命令 | Ctrl+P / Shift+Ctrl+P                                |
| API key 校验    | 前置校验，无 key 抛异常                | 列表已过滤掉无 key 的                                |
| thinking level  | 保持当前，re-clamp                     | scoped: 用 scoped 的 level; available: 保持+re-clamp |
| 持久化          | 写 session + settings                  | 写 session + settings                                |
| Extensions 事件 | `source: 'set'`                        | `source: 'cycle'`                                    |

---

## 六、跨 Provider 消息格式转换（最难的部分）

### 6.1 问题：切换模型后，历史消息格式可能不兼容

例如：从 Claude 切到 GPT，历史中有 Anthropic 的 thinking block 和 tool call ID 格式。

### 6.2 stream() 路由机制

```ts
// packages/ai/src/stream.ts — 极简路由
function stream(model, context, options) {
  const provider = getApiProvider(model.api) // 按 model.api 查注册表
  return provider.stream(model, context, options)
}
```

### 6.3 transformMessages() — 跨 Provider 兼容核心

```ts
// packages/ai/src/providers/transform-messages.ts
function transformMessages(messages, model, normalizeToolCallId?) {
  // ★ 判断是否同模型
  const isSameModel = msg.provider === model.provider
                   && msg.api === model.api
                   && msg.model === model.id

  // === 第一遍扫描：逐消息转换 ===

  // 1. Thinking Block 处理
  if (block.type === 'thinking') {
    if (block.redacted && !isSameModel)
      → 跳过（加密的 redacted thinking 仅同模型有效）
    if (isSameModel && block.thinkingSignature)
      → 保留（用于 replay）
    if (empty thinking)
      → 跳过
    if (!isSameModel)
      → { type: 'text', text: block.thinking }  // ★ 降级为普通文本
  }

  // 2. Tool Call ID 归一化
  //    OpenAI Responses API 生成 450+ 字符含 | 的 ID
  //    Anthropic 要求 ^[a-zA-Z0-9_-]+$ 最长 64 字符
  if (!isSameModel && normalizeToolCallId) {
    const normalizedId = normalizeToolCallId(toolCall.id, model, source)
    toolCallIdMap.set(originalId, normalizedId)
  }

  // 3. toolResult 的 toolCallId 同步更新
  if (msg.role === 'toolResult') {
    const normalizedId = toolCallIdMap.get(msg.toolCallId)
    if (normalizedId) msg.toolCallId = normalizedId
  }

  // 4. thoughtSignature 清理（Google 特有的 opaque signature）
  if (!isSameModel) delete toolCall.thoughtSignature

  // === 第二遍扫描：孤儿 Tool Call 修复 ===

  // 如果 assistant 消息有 tool_call 但没有对应的 tool_result
  // （可能是中断或错误），插入合成的错误结果
  if (!existingToolResultIds.has(tc.id)) {
    → 插入 { role: 'toolResult', isError: true, content: 'No result provided' }
  }

  // 跳过 stopReason === 'error' | 'aborted' 的 assistant 消息
  // 避免重放不完整的 turn
}
```

### 6.4 各 Provider 如何调用 transformMessages

每个具体 provider 实现（如 `anthropic.ts`, `openai-completions.ts`）在发送请求前调用：

```ts
const messages = transformMessages(context.messages, model, normalizeToolCallId)
```

这确保了：

- 切换到 Anthropic 时，OpenAI 的长 tool call ID 被截断归一化
- 切换到 OpenAI 时，Anthropic 的 thinking block 被降级为文本
- 加密的 thinking（redacted/signed）被正确跳过
- 孤儿 tool call 不会导致 API 报错

---

## 七、模型配置持久化

### 7.1 三层持久化

```
1. Agent._state.model        ← 内存中的当前模型，立即生效
2. sessionManager            ← session JSON，记录模型变更历史
3. settingsManager           ← ~/.pi/settings.json，全局默认
```

### 7.2 Session 恢复

```ts
// model-resolver.ts
async function restoreModelFromSession(savedProvider, savedModelId, currentModel, modelRegistry) {
  const restoredModel = modelRegistry.find(savedProvider, savedModelId)
  const hasApiKey = restoredModel ? await modelRegistry.getApiKey(restoredModel) : false

  if (restoredModel && hasApiKey) return restoredModel

  // Fallback: 模型不存在或没 API key
  // → 用当前模型 → 用已知 provider 默认模型 → 用第一个可用模型
}
```

### 7.3 enabledModels — 用户自定义循环范围

```ts
// settings.json
{
  "enabledModels": ["claude-*:high", "gpt-5*:medium"]  // glob + thinking level
}
```

通过 `ScopedModelsSelectorComponent`（Ctrl+O → Model Configuration）管理：

- Session 级别的启用/禁用
- Ctrl+S 持久化到 settings
- 支持按 provider 批量启用/禁用
- 支持 Alt+↑↓ 调整顺序

---

## 八、Extensions 扩展机制

Extensions 可以：

1. **注册自定义 Provider**：`modelRegistry.registerProvider(name, { models, baseUrl, apiKey, streamSimple, oauth })`
2. **监听 model_select 事件**：`{ type: 'model_select', model, previousModel, source: 'set'|'cycle'|'restore' }`
3. **动态修改模型**：OAuth provider 的 `modifyModels()` 可以在运行时修改模型属性

---

## 九、面试要点总结

### Q1：切换模型时历史消息怎么处理？

`transformMessages()` 在每次 LLM 调用前将历史消息转换为目标 provider 格式：

- **Thinking blocks**：同模型保留（含签名用于 replay），跨模型降级为文本，加密的删除
- **Tool Call ID**：归一化到目标格式（OpenAI 450+ 字符 → Anthropic 64 字符限制）
- **孤儿 Tool Calls**：自动插入合成的 error 结果，避免 API 报错
- **错误消息**：跳过 error/aborted 的 assistant 消息

### Q2：模型路由是怎么实现的？

`model.api` 字段是路由键，通过 `ApiProviderRegistry`（策略模式）分发：

```
model.api → apiProviderRegistry.get(api) → provider.stream(model, context, options)
```

中间没有 if-else，完全基于注册表查找，新 provider 只需 `registerApiProvider()` 即可。

### Q3：Ctrl+P 和 Ctrl+L 有什么区别？

- **Ctrl+P**：循环切换，modulo 索引，快速无 UI overlay，scoped models 可带 thinking level
- **Ctrl+L**：打开搜索选择器，支持模糊搜索、scope 切换、排序，选择后保存为默认

### Q4：自定义模型如何配置？

通过 `~/.pi/models.json`，支持：

- 定义全新 provider + models（需 baseUrl + apiKey + api）
- Override 内置 provider 的 baseUrl/headers/apiKey
- Per-model override（改名、改 maxTokens、改 compat 等）
- 支持 `$(command)` 和 `$ENV_VAR` 动态解析 API key

### Q5：初始模型选择的优先级？

1. `--provider + --model` CLI 参数
2. `--models` 的第一个模型
3. settings.json 保存的默认模型
4. 已知 provider 的硬编码默认模型表
5. 第一个有 API key 的模型

### Q6：tokenizer 兼容性怎么处理？

每个 Model 有 `contextWindow` 和 `maxTokens` 字段。切换模型后 `setThinkingLevel()` 会 **re-clamp** thinking budget 到新模型的能力范围内，避免超出限制。

### Q7：设计亮点

1. **`model.api` vs `model.provider` 分离**：一个 provider（如 github-copilot）可路由到不同的 API 实现（anthropic/openai/google），一个 API 可被多个 provider 使用
2. **策略模式注册表**：零 if-else 的 provider 路由，扩展性极强
3. **三级状态**：内存（秒级）→ Session（会话级）→ Settings（永久），各层关注不同
4. **Scoped Models**：`--models` 标志 + enabledModels 设置 + 模糊 glob 匹配，灵活限定切换范围
5. **transformMessages 的两遍扫描**：第一遍格式转换 + 第二遍孤儿修复，保证语义完整性
6. **Extensions 可插拔**：OAuth、自定义 stream、model_select 事件，完全开放
