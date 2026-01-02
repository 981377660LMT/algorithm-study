这篇文章的核心在于**解构 Claude Code 的工作原理，并指导如何利用现有的 LLM（大语言模型）和工具链，构建一个类似的、能够自主编写和执行代码的 AI 编程助手**。

一针见血地讲，它主要包含以下几个关键步骤：

1.  **核心大脑（LLM 选择）**：

    - 选择一个具备强逻辑推理和代码生成能力的模型（如 Claude 3.5 Sonnet, GPT-4o）。
    - 这是系统的“大脑”，负责理解需求、规划任务和生成代码。

2.  **执行环境（REPL/Sandbox）**：

    - 构建一个安全的沙箱环境（如 Docker 容器或受限的 Shell）。
    - 让 AI 不仅能写代码，还能**运行代码**、查看报错、读取文件系统。这是 Claude Code 区别于普通聊天机器人的关键——它有“手”。

3.  **工具链集成（Tool Use/Function Calling）**：

    - 定义一套工具集（Tools），例如：`read_file`（读文件）、`write_file`（写文件）、`run_command`（运行终端命令）、`search_code`（搜索代码）。
    - 利用 LLM 的 Function Calling 能力，让模型自主决定何时调用什么工具。

4.  **循环反馈机制（Agent Loop）**：

    - 实现 **思考 -> 行动 -> 观察 -> 修正** 的循环（ReAct 模式）。
    - 模型写代码 -> 运行代码 -> 发现报错 -> 读取错误日志 -> 修改代码 -> 再次运行，直到任务完成。

5.  **上下文管理（Context Management）**：
    - 如何高效地将项目结构、相关文件内容喂给模型，同时避免 Token 溢出（例如使用 RAG 技术或智能截断）。

**总结：**
实现自己的 Claude Code，本质上就是写一个 **Agent（智能体）**，它通过 **Shell/文件系统接口** 连接 LLM 与你的本地代码库，让 LLM 变成一个能干活的“全栈工程师”，而不仅仅是一个“问答机器人”。

---

Claude Code 的设计理念非常清晰：

- 无状态 Agent：每次对话创建新实例，简单可靠
- Agentic Loop：工具调用 → 执行 → 反馈 → 继续，直到任务完成
- 权限控制：读操作自动执行，写操作需要确认，安全可控

## 什么是 Coding Agent、Claude Code 的设计理念、核心架构概览

### 1.1 什么是 Coding Agent

#### 1.1.1 从 Chatbot 到 Agent

传统 Chatbot 是“有丰富知识的军师”，提供建议但需人工执行；**Coding Agent** 则是“能听懂指令的实习生”，能自主翻阅代码、修改文件、执行命令并修复错误。
其核心公式为：**Coding Agent = LLM + System Prompt + Context + Tools**。

- **LLM**：大脑，负责逻辑推理。
- **System Prompt**：性格与行为边界。
- **Context**：当前任务的感知。
- **Tools**：伸向物理世界的手（读写文件、执行 Shell 等）。

当内置工具不够用时，我们会通过 MCP 等机制把更多外部能力挂接进来；再往后，又演化出了更高层的 Skills，把常见目标和操作封装成“可复用的一键能力”，让 Agent 调用起来更智能、更高效。当任务足够复杂，一个“实习生”难以招架时，Sub-agents 登场了，我们把大任务拆解开，由专注于规划的 Planner、专注于编码的 Coder、专注于测试的 Tester 组成“虚拟开发团队”。它们各司其职，共享上下文，协同作战。

从这个角度看，前半场是 Prompt 工程：设计好性格、目标和边界；后半程是 Context 工程：持续构建、维护和更新上下文；而不断丰富的 Tools 与 Skills 一方面拓宽了 Agent 的能力边界，另一方面，它们的返回结构又在反过来充实 Context，让 Agent 在一个闭环中不断自我感知、自我决策和自我纠偏。

#### 1.1.2 核心特征

1.  **自主决策**：理解高层需求并分解步骤。
2.  **工具调用**：通过 Read/Write/Bash 等工具与外部交互。
3.  **循环执行**：遵循 **观察 -> 思考 -> 行动** 的 Agentic Loop。
4.  **上下文感知**：理解项目结构与历史。

#### 1.1.3 Coding Agent vs 传统 IDE 插件

传统插件是确定性的程序逻辑，而 Agent 的核心是**不确定性**，其解决路径是动态规划的，具备全局感知和学习能力。

### 1.2 Claude Code 的设计理念

> Claude Code 是 Anthropic 官方推出的 AI 编程助手 CLI 工具。Claude Code 90% 的代码是它自己写的

1.  **"On Distribution" 技术栈**：选择模型熟悉的技术（TS, React, Ink），让模型能“自己构建自己”。
2.  **极简架构**：将逻辑交给模型，而非复杂的规则引擎。
3.  **本地优先**：直接访问本地文件系统，追求低延迟。
4.  **权限控制**：默认只读，写入需确认，危险操作强提醒。

### 1.3 Blade 项目介绍

Blade 借鉴了 Claude Code 的设计理念，使用相同技术栈（TS + Ink + Bun），并支持 OpenAI 兼容 API。其结构清晰，涵盖了从 Agent 核心到工具系统、UI 渲染及 MCP 协议的完整实现。

| 模块       | 主要文件             | 职责                 |
| ---------- | -------------------- | -------------------- |
| CLI 入口   | blade.tsx            | 命令解析、启动流程   |
| Agent 核心 | Agent.ts             | LLM 交互、执行循环   |
| 工具系统   | tools/               | 工具定义、注册、执行 |
| 执行管道   | ExecutionPipeline.ts | 7 阶段执行流程       |
| 上下文管理 | ContextManager.ts    | 消息历史、压缩       |
| UI 系统    | ui/                  | 终端界面渲染         |
| MCP 集成   | mcp/                 | 协议实现             |
| 配置管理   | config/              | 配置加载、合并       |

### 1.4 Agentic Loop 核心概念

#### 1.4.1 什么是 Agentic Loop

这是 Agent 的心脏，描述了 Agent 如何持续与 LLM 和工具交互。

#### 1.4.2 关键流程

1.  **构建消息**：整合 System Prompt、历史记录与用户输入。
2.  **调用 LLM**：获取回复或工具调用请求。
3.  **执行工具**：通过执行管道运行工具并获取结果。
4.  **注入反馈**：将工具执行结果返回给 LLM，进入下一轮循环。

```ts
// src/agent/Agent.ts (简化版)
private async executeLoop(
  message: string,
  context: ChatContext,
  options?: LoopOptions
): Promise<LoopResult> {
  // 1. 构建消息历史
  const messages: Message[] = [
    { role: 'system', content: systemPrompt },
    ...context.messages,
    { role: 'user', content: message }
  ];

  // 2. 获取可用工具
  const tools = this.registry.getFunctionDeclarations();

  // 3. 核心循环
  while (true) {
    // 3.1 调用 LLM
    const response = await this.chatService.chat(messages, tools);

    // 3.2 检查是否完成（无工具调用）
    if (!response.toolCalls || response.toolCalls.length === 0) {
      return { success: true, finalMessage: response.content };
    }

    // 3.3 执行工具调用
    for (const toolCall of response.toolCalls) {
      const result = await this.executionPipeline.execute(
        toolCall.function.name,
        JSON.parse(toolCall.function.arguments)
      );

      // 3.4 将结果注入消息历史
      messages.push({
        role: 'tool',
        tool_call_id: toolCall.id,
        content: result.llmContent
      });
    }

    // 继续下一轮循环...
  }
}
```

#### 1.4.3 终止条件

任务完成（无工具调用）、达到轮次上限（如 100 次）、用户中止或发生不可恢复的错误。

## 技术栈选择、项目结构设计、开发环境搭建

### 2.1 技术栈选择

#### 2.1.1 核心技术栈

构建一个 Coding Agent，我们需要选择合适的技术栈。参考 Claude Code 的选择，Blade 使用以下技术：

- **语言**: TypeScript (模型最熟悉的语言)
- **UI 框架**: Ink (使用 React 组件构建 CLI 界面)
- **运行时/构建**: Bun (极速运行时与构建工具)
- **状态管理**: Zustand (解耦 UI 与逻辑的状态管理)
- **LLM SDK**: OpenAI SDK (兼容多种模型后端)

#### 2.1.2 为什么选择 TypeScript

作为前端，最熟悉的语言是 TS，且前端生态提供了构建 CLI 所需的所有工具。

#### 2.1.3 为什么选择 Ink

Ink 允许使用 React 的组件化思维来组织 CLI 代码，支持智能 Diff 更新界面，避免了传统 CLI 的全屏重绘。

#### 2.1.4 为什么选择 Zustand

Blade 最初使用 React Context，但遇到了双轨数据源不一致、非 React 环境（Agent 类）无法访问 Context 等问题。
通过 Zustand，我们将状态管理与 React 解耦。UI 只是消费者，而 Agent 核心逻辑可以通过 `vanillaStore.getState()` 同步获取最新状态，完美解决了状态同步难题。

### 2.2 项目结构设计

#### 2.2.1 目录结构

blade/
├── src/ # 源代码
│ ├── agent/ # Agent 核心
│ │ ├── Agent.ts # 主 Agent 类
│ │ ├── ExecutionEngine.ts # 执行引擎
│ │ ├── types.ts # 类型定义
│ │ └── subagents/ # 子 Agent 系统
│ │
│ ├── tools/ # 工具系统
│ │ ├── builtin/ # 内置工具
│ │ │ ├── file/ # 文件工具 (Read, Write, Edit)
│ │ │ ├── search/ # 搜索工具 (Glob, Grep)
│ │ │ ├── shell/ # Shell 工具 (Bash)
│ │ │ └── web/ # 网络工具 (WebFetch)
│ │ ├── registry/ # 工具注册表
│ │ ├── execution/ # 执行管道
│ │ └── types/ # 工具类型
│ │
│ ├── ui/ # UI 系统
│ │ ├── components/ # React 组件
│ │ ├── hooks/ # 自定义 Hooks
│ │ ├── themes/ # 主题管理
│ │ └── App.tsx # UI 入口
│ │
│ ├── config/ # 配置管理
│ │ ├── ConfigManager.ts # 配置管理器
│ │ ├── types.ts # 配置类型
│ │ └── defaults.ts # 默认配置
│ │
│ ├── context/ # 上下文管理
│ │ ├── ContextManager.ts # 上下文管理器
│ │ ├── CompactionService.ts # 压缩服务
│ │ └── storage/ # 存储实现
│ │
│ ├── services/ # 服务层
│ │ ├── ChatServiceInterface.ts
│ │ └── OpenAIChatService.ts
│ │
│ ├── mcp/ # MCP 协议
│ ├── prompts/ # 提示词管理
│ ├── logging/ # 日志系统
│ ├── store/ # Zustand Store
│ │
│ ├── cli/ # CLI 相关
│ │ ├── config.ts # CLI 配置
│ │ └── middleware.ts # 中间件
│ │
│ └── blade.tsx # 主入口
│
├── tests/ # 测试文件
│ ├── unit/ # 单元测试
│ ├── integration/ # 集成测试
│ └── e2e/ # 端到端测试
│
├── docs/ # 文档
│ ├── public/ # 用户文档
│ ├── development/ # 开发文档
│ └── tutorial/ # 本系列教程
│
├── scripts/ # 构建脚本
├── dist/ # 构建产物
├── package.json
├── tsconfig.json
└── biome.json # 代码风格配置

#### 2.2.2 模块职责划分

- **入口层**: 命令解析、启动流程。
- **UI 层**: 界面渲染、用户交互。
- **业务层 (Agent)**: 核心逻辑，完全无状态，负责接收上下文、调用 LLM 并决定下一步行动。
- **服务层**: 通用服务（ChatService, ConfigManager）。

### 2.3 核心依赖详解

#### 2.3.1 关键依赖

- **OpenAI SDK**: 用于与 LLM API 通信。
- **Zod**: 用于运行时参数验证，确保工具调用参数准确。
- **js-tiktoken**: 用于精确估算 Token 数量，辅助上下文压缩。

### 2.4 Hello World Agent

#### 2.4.1 SimpleAgent.ts (核心逻辑)

```typescript
import OpenAI from 'openai'

export class SimpleAgent {
  private client: OpenAI
  constructor(config: { apiKey: string; baseURL?: string }) {
    this.client = new OpenAI({ apiKey: config.apiKey, baseURL: config.baseURL })
  }

  async chat(message: string): Promise<string> {
    const response = await this.client.chat.completions.create({
      model: 'gpt-4',
      messages: [{ role: 'user', content: message }]
    })
    return response.choices[0]?.message?.content || ''
  }
}
```

#### 2.4.2 App.tsx (UI 界面)

```tsx
import React, { useState } from 'react'
import { Box, Text } from 'ink'
import TextInput from 'ink-text-input'
import Spinner from 'ink-spinner'

export const App: React.FC<{ apiKey: string }> = ({ apiKey }) => {
  const [input, setInput] = useState('')
  const [response, setResponse] = useState('')
  const [isLoading, setIsLoading] = useState(false)

  const handleSubmit = async (value: string) => {
    setIsLoading(true)
    // 调用 Agent 逻辑...
    setIsLoading(false)
  }

  return (
    <Box flexDirection="column" padding={1}>
      <Text bold color="cyan">
        🗡️ My Coding Agent
      </Text>
      {isLoading ? <Spinner type="dots" /> : <Text>{response}</Text>}
      <TextInput value={input} onChange={setInput} onSubmit={handleSubmit} />
    </Box>
  )
}
```

## yargs CLI 框架、中间件机制、启动流程详解

### 3.1 CLI 架构概览

Blade 的启动是一个**并行且分层**的过程。当用户输入 `blade` 命令时，系统会同时启动版本检查与参数解析，确保在进入交互界面前，环境已就绪。

- **入口文件 (`blade.tsx`)**: 负责 Shebang 声明、早期 Debug 解析及 yargs 实例创建。
- **中间件链**: 在命令执行前，完成配置加载、权限校验和参数转换。
- **UI 桥接**: 通过 `AppWrapper` 将 CLI 参数注入 React 上下文，并处理版本更新逻辑。

### 3.2 yargs 配置与中间件

#### 3.2.1 全局选项分组

通过 yargs 的 `group` 属性，我们将选项分为：调试选项、输出选项、安全选项和 AI 选项，提升 `--help` 的可读性。

#### 3.2.2 核心中间件机制

中间件是 CLI 的“过滤器”，Blade 实现了三个关键中间件：

1.  **`loadConfiguration`**: 最核心的中间件。它初始化 `ConfigManager`，并将配置同步到 Zustand Store。这样，后续的 Agent 类（非 React 环境）也能通过 `getState()` 获取配置。
2.  **`validatePermissions`**: 处理 `--yolo` 等快捷参数，并检测 `allowedTools` 与 `disallowedTools` 的冲突。
3.  **`validateOutput`**: 确保输出格式（如 JSON）仅在非交互的 `--print` 模式下生效。

### 3.3 启动流程详解

#### 3.3.1 默认命令 ($0)

当不带子命令启动时，执行默认逻辑：

1.  解析初始消息（如 `blade "帮我写个脚本"`）。
2.  调用 Ink 的 `render` 函数启动 React UI。
3.  配置 `patchConsole` 和自定义的退出清理函数。

#### 3.3.2 UI 初始化 (`AppWrapper`)

在进入主界面前，`AppWrapper` 会执行以下步骤：

- **版本检查**: 等待启动时开启的并行 Promise。若有更新，显示 `UpdatePrompt`。
- **状态同步**: 调用 `initializeStoreState`，若未配置模型，则引导进入 `needsSetup` 状态。
- **环境预加载**: 加载主题、注册 Subagents、初始化 Hooks 系统。

### 3.4 版本检查与自动升级

为了保持工具的先进性，Blade 内置了智能版本管理：

- **缓存机制**: 结果缓存 1 小时，避免频繁请求 NPM Registry。
- **跳过逻辑**: 支持“跳过当前版本”，直到下一个新版本发布。
- **交互式升级**: 用户选择 `Update now` 后，程序会自动执行 `npm install -g` 并优雅退出。

### 3.5 错误处理

- **CLI 层**: 利用 `cli.fail` 捕获参数错误，并提供“Did you mean”建议。
- **UI 层**: 使用 React `ErrorBoundary` 捕获渲染崩溃，防止终端直接卡死，并输出友好的错误堆栈。

## Agent 类设计、无状态架构、核心执行循环

### 4.1 Agent 设计理念

#### 4.1.1 无状态设计原则

Blade 的 Agent 采用**无状态设计（Stateless Design）**。Agent 本身不保存任何会话状态（如消息历史），所有状态通过 `context` 参数传入。这使得 Agent 实例可以随用随弃，天然支持并发安全且易于测试。

```ts
/**
 * Agent核心类 - 无状态设计
 *
 * 设计原则：
 * 1. Agent 本身不保存任何会话状态（sessionId, messages 等）
 * 2. 所有状态通过 context 参数传入
 * 3. Agent 实例可以每次命令创建，用完即弃
 * 4. 历史连续性由外部 SessionContext 保证
 */
export class Agent {
  // 配置（只读）
  private config: BladeConfig
  private runtimeOptions: AgentOptions

  // 初始化状态
  private isInitialized = false
  private activeTask?: AgentTask

  // 核心组件
  private executionPipeline: ExecutionPipeline
  private systemPrompt?: string
  private chatService!: IChatService
  private executionEngine!: ExecutionEngine
  private attachmentCollector?: AttachmentCollector

  constructor(
    config: BladeConfig,
    runtimeOptions: AgentOptions = {},
    executionPipeline?: ExecutionPipeline
  ) {
    this.config = config
    this.runtimeOptions = runtimeOptions
    this.executionPipeline = executionPipeline || this.createDefaultPipeline()
  }
}
```

### 4.2 Agent 类结构

#### 4.2.1 核心组件

Agent 类通过静态工厂方法 `create()` 初始化，其核心组件包括：

- **ExecutionPipeline**: 负责工具的调度与执行。
- **ChatService**: 封装 LLM API 调用。
- **ExecutionEngine**: 处理复杂的执行逻辑。
- **AttachmentCollector**: 处理 `@` 文件提及。

### 4.3 Agentic Loop 详解

#### 4.3.1 执行循环入口 (`chat`)

`chat` 方法是 Agent 的主入口，负责解析 `@` 提及并根据模式选择 `runLoop` 或 `runPlanLoop`。

#### 4.3.2 核心执行循环 (`executeLoop`)

这是实现 **ReAct (Reasoning + Acting)** 模式的核心：

1.  **准备**: 获取可用工具，构建消息历史。
2.  **调用 LLM**: 获取思考过程及工具调用请求。
3.  **执行工具**: 通过 `ExecutionPipeline` 执行工具（含权限检查）。
4.  **反馈**: 将工具结果注入消息历史，进入下一轮循环。

```ts
// src/agent/Agent.ts:401-1080 (简化版)

private async executeLoop(
  message: string,
  context: ChatContext,
  options?: LoopOptions,
  systemPrompt?: string
): Promise<LoopResult> {
  const startTime = Date.now();

  // === 1. 准备阶段 ===

  // 获取可用工具 (基于权限模式过滤)
  const registry = this.executionPipeline.getRegistry();
  const tools = registry.getFunctionDeclarationsByMode(context.permissionMode);

  // 构建消息历史
  const messages: Message[] = [];
  if (systemPrompt) {
    messages.push({ role: 'system', content: systemPrompt });
  }
  // 注入历史消息和当前用户消息
  messages.push(...context.messages);
  messages.push({ role: 'user', content: message });

  // === 2. 循环配置 ===

  const TURN_LIMIT = 100; // 硬性轮次上限，防止无限循环
  const isYoloMode = context.permissionMode === PermissionMode.YOLO; // YOLO 模式无限制

  // 优先级: CLI参数 > 调用参数 > 配置文件 > 默认值(-1)
  const configuredMaxTurns = this.runtimeOptions.maxTurns ??
                             options?.maxTurns ??
                             this.config.maxTurns ?? -1;

  // 特殊值处理：0 = 禁用对话
  if (configuredMaxTurns === 0) {
    return { success: false, error: { type: 'chat_disabled', message: '...' } };
  }

  // 应用安全上限 (YOLO 模式除外)
  const maxTurns = configuredMaxTurns === -1
    ? TURN_LIMIT
    : Math.min(configuredMaxTurns, TURN_LIMIT);

  let turnsCount = 0;
  const allToolResults: ToolResult[] = [];

  // === 3. 核心循环 ===
  // Agentic Loop: 循环调用直到任务完成
  while (true) {
    // 3.1 检查中断信号
    if (options?.signal?.aborted) {
      return { success: false, error: { type: 'aborted' } };
    }

    // 3.2 检查并压缩上下文
    // 基于 Token 使用量自动压缩历史消息
    await this.checkAndCompactInLoop(context, turnsCount, lastPromptTokens);

    // 3.3 轮次计数
    turnsCount++;
    options?.onTurnStart?.({ turn: turnsCount, maxTurns });

    // 3.4 调用 LLM
    const turnResult = await this.chatService.chat(
      messages,
      tools,
      options?.signal
    );

    // 3.5 通知 UI 显示内容 (流式或完整内容)
    if (turnResult.content && options?.onContent) {
      options.onContent(turnResult.content);
    }
    if (turnResult.reasoningContent && options?.onThinking) {
      options.onThinking(turnResult.reasoningContent); // DeepSeek R1 等推理模型支持
    }

    // 3.6 检查是否完成（无工具调用）
    if (!turnResult.toolCalls || turnResult.toolCalls.length === 0) {
      // 意图未完成检测 (如 "让我来..." 但未调用工具)
      if (this.detectIncompleteIntent(turnResult.content)) {
          messages.push({ role: 'user', content: '请执行你提到的操作...' });
          continue;
      }

      return {
        success: true,
        finalMessage: turnResult.content,
        metadata: { turnsCount, toolCallsCount: allToolResults.length }
      };
    }

    // 3.7 添加 assistant 消息到历史
    messages.push({
      role: 'assistant',
      content: turnResult.content || '',
      tool_calls: turnResult.toolCalls,
    });

    // 3.8 执行每个工具调用
    for (const toolCall of turnResult.toolCalls) {
      if (toolCall.type !== 'function') continue;

      // 解析参数 (含自动修复逻辑)
      const params = JSON.parse(toolCall.function.arguments);

      // 通过 ExecutionPipeline 执行 (包含权限检查、确认、执行、Hook)
      const result = await this.executionPipeline.execute(
        toolCall.function.name,
        params,
        {
          sessionId: context.sessionId,
          signal: options?.signal,
          confirmationHandler: context.confirmationHandler, // 传递确认处理器
          permissionMode: context.permissionMode,
        }
      );
      allToolResults.push(result);

      // 通知 UI 更新状态
      options?.onToolResult?.(toolCall, result);

      // 检查是否需要退出循环 (如 ExitPlanMode 工具)
      if (result.metadata?.shouldExitLoop) {
          return { success: result.success, finalMessage: '循环已退出', ... };
      }

      // 添加工具结果到消息历史
      messages.push({
        role: 'tool',
        tool_call_id: toolCall.id,
        name: toolCall.function.name,
        content: result.llmContent || result.displayContent || '',
      });
    }

    // === 4. 轮次上限处理 (非 YOLO 模式) ===
    // 达到 maxTurns 后暂停，询问用户是否继续
    if (turnsCount >= maxTurns && !isYoloMode) {
      if (options?.onTurnLimitReached) {
        const response = await options.onTurnLimitReached({ turnsCount });
        if (response.continue) {
          // 用户选择继续：压缩上下文，重置计数器，继续循环
          await this.forceCompact(context);
          turnsCount = 0;
          continue; // 直接继续当前 while(true) 循环
        }
      }

      // 用户选择停止或非交互模式
      return {
        success: false, // 注意：此处根据设计可能是 true 或 false，视具体需求而定，代码中返回 false + error
        error: { type: 'max_turns_exceeded', message: '...' },
        metadata: { ... }
      };
    }
  }
}
```

### 4.4 消息格式与工具调用

Blade 遵循 OpenAI Chat 格式，定义了 `Message` 和 `ToolCall` 接口，支持 DeepSeek R1 等模型的思维链（`reasoning_content`）。

```ts
export interface Message {
  role: 'system' | 'user' | 'assistant' | 'tool'
  content: string
  // DeepSeek R1 等推理模型产生的思维链内容
  reasoning_content?: string
  // assistant 消息专用：发起的工具调用列表
  tool_calls?: ToolCall[]
  // tool 消息专用：关联的调用 ID
  tool_call_id?: string
  // tool 消息专用：工具名称
  name?: string
}

export interface ToolCall {
  id: string
  type: 'function'
  function: {
    name: string
    arguments: string // JSON 字符串
  }
}
```

### 4.5 循环控制机制

- **轮次限制**: 设置硬性安全上限（100 轮），防止无限循环。
- **中断信号**: 全程检查 `AbortSignal`，支持 Ctrl+C 立即停止。
- **意图未完成检测**: 识别模型“只说不做”的幻觉并自动重试。
  某些模型（尤其是较小的模型或中文模型）可能会回复"让我来..."或以冒号结尾，但没有实际生成 `tool_calls`。Blade 实现了启发式检测来处理这种情况。

  1. 直接停止，人工介入，减少重试 token 消耗
  2. 多模型架构，可以自动切更大/更强的模型，需要考虑：fallback 触发条件、模型切换的上下文传递、成本控制。
  3. System Prompt 增强，比如【禁止使用"让我来..."、"我将..."等表达后不跟随工具调用。】这个我尝试过，小模型的 instruction following 能力本身就弱，堆 prompt 帮助有限
  4. tool_choice 改成 required，很僵硬，不推荐
  5. 更激进的重试提示，加入 system-reminder。可以试，但如果陷入死循环了，换个措辞大概率也没用。

  ```ts
  // src/agent/Agent.ts:680-713

  const INCOMPLETE_INTENT_PATTERNS = [
    /：\s*$/, // 中文冒号结尾
    /:\s*$/, // 英文冒号结尾
    /\.\.\.\s*$/, // 省略号结尾
    /让我(先|来|开始|查看|检查|修复)/, // 中文意图词
    /Let me (first|start|check|look|fix)/i // 英文意图词
  ]

  const content = turnResult.content || ''
  const isIncompleteIntent = INCOMPLETE_INTENT_PATTERNS.some(p => p.test(content))

  // 最多重试 2 次
  if (isIncompleteIntent && recentRetries < 2) {
    messages.push({
      role: 'user',
      content: '请执行你提到的操作，不要只是描述。'
    })
    continue
  }
  ```

### 4.6 Plan 模式

一种只读调研模式。通过专用系统提示词和工具过滤策略，让 Agent 仅使用只读工具进行任务规划。

### 4.7 @ 文件提及处理

支持在消息中使用 `@path/to/file` 语法。Agent 会自动读取文件内容并以 `<file>` 标签形式注入上下文。

### 4.8 Token 管理与压缩

当 Token 使用量达到窗口的 80% 时触发自动压缩，通过过滤历史消息或摘要技术确保对话持续。

### 4.9 踩坑记录

1.  **上下文丢失**: 必须通过外部 Session 维护历史。
2.  **循环误判**: Plan 模式应跳过内容重复检测。
    移除 LoopDetectionService，改用更简单的策略：
    - maxTurns 配置项控制最大轮次
    - 硬性安全上限 SAFETY_LIMIT = 100
    - 让 LLM 在系统提示中自己判断是否陷入循环
3.  **空转问题**: 需处理模型生成了描述但未生成 `tool_calls` 的情况。
4.  **孤儿 tool 消息**: 压缩上下文时需确保 `tool` 消息与其 `assistant` 消息成对存在。

## 系统提示词、Plan 模式提示词、工具提示词、压缩提示词

### 5.1 提示词工程概述

System Prompt 是 Coding Agent 的“灵魂”，它定义了 Agent 的身份、能力边界、行为准则和输出风格。Blade 的设计理念是：**更少地规定性指导，更多地信任模型的判断力**，同时保留必要的硬性约束。

- Agent 是谁？
  身份定义
- 能做什么？不能做什么？
  安全边界
- 如何输出？
  风格约束
- 如何使用工具？
  工具策略
- 如何处理复杂任务？
  任务管理

### 5.2 Blade 提示词架构

Blade 采用固定顺序的分层架构，确保优先级清晰：

1.  **环境上下文**：动态生成，包含当前日期、工作目录、系统信息。
2.  **基础提示词**：定义身份（Blade Code）与核心逻辑。
3.  **项目配置 (BLADE.md)**：项目特有的约定与命令。
4.  **追加内容**：用户自定义的临时约束。

#### 5.2.1 提示词构建器实现

```typescript
export async function buildSystemPrompt(options: BuildOptions): Promise<string> {
  const parts = [
    getEnvironmentContext(), // 1. 环境
    options.mode === 'plan' ? PLAN_PROMPT : DEFAULT_PROMPT, // 2. 基础
    await loadBladeConfig(options.projectPath), // 3. 项目 (BLADE.md)
    options.append // 4. 追加
  ]
  return parts.filter(Boolean).join('\n\n---\n\n')
}
```

### 5.3 核心行为准则

#### 5.3.1 输出风格（4 行原则）

为了提升交互效率，Blade 规定：除了生成代码或复杂调试外，所有解释、确认和状态更新应控制在 **4 行以内**。

#### 5.3.2 执行效率（Action over narration）

禁止在工具调用前输出冗余文本。

- **Bad**: "我将为你读取 package.json..." [调用工具]
- **Good**: [直接调用工具]

#### 5.3.3 语言要求

强制要求模型使用 **简体中文** 进行所有响应，包括解释、错误提示和代码注释。

### 5.4 Plan 模式：双层防护机制

Plan 模式是 Blade 的只读调研阶段，通过以下两层确保安全：

1.  **提示词约束**：注入 `PLAN_MODE_SYSTEM_PROMPT`，明确禁止修改文件。
2.  **工具过滤**：在权限层自动拒绝所有非只读工具。

#### 5.4.1 五阶段检查点

Plan 模式引导模型遵循：**探索 (Explore) -> 设计 (Design) -> 评审 (Review) -> 展示计划 (Present) -> 退出 (Exit)**。模型必须调用 `ExitPlanMode` 工具并提交完整计划供用户审批。

### 5.5 上下文压缩策略

当对话接近 Token 窗口限制（80%）时，Blade 会触发智能压缩：

- **9 部分结构**：要求模型总结原始意图、技术概念、已修改文件、待办事项等。
- **孤儿消息处理**：压缩时自动过滤掉失去关联的 `tool` 消息，防止 LLM 报错。

### 5.6 工具描述与 BLADE.md

- **结构化描述**：不同于 Claude Code 的纯文本，Blade 使用 TypeScript 接口定义工具的 `usageNotes`、`examples` 和 `important` 约束，便于程序化处理。
- **项目大脑 (BLADE.md)**：类似于 `CLAUDE.md`，用于存储项目特定的构建命令、代码风格和架构约定，Agent 会在每轮对话中感知这些信息。

## 工具系统设计与实现：工具抽象、内置工具与注册机制

### 6.1 工具系统概述

#### 6.1.1 什么是工具

在 Coding Agent 中，**工具 (Tool)** 是 Agent 与外部世界交互的桥梁。LLM 本身只能生成文本，但通过工具，它可以获得“手”和“眼”：

- **读取和修改文件**：直接操作代码库。
- **执行 Shell 命令**：编译、测试、安装依赖。
- **搜索代码库**：通过 Glob 或 Grep 快速定位。
- **获取网页内容**：查阅最新的 API 文档。

#### 6.1.2 工具系统架构

Blade 的工具系统采用分层设计，确保了高度的可扩展性：

- `builtin/`：核心内置工具（文件、Shell、搜索等）。
- `registry/`：工具注册表，管理工具的发现与检索。
- `execution/`：执行管道，处理权限校验与结果封装。
- `core/`：提供 `createTool` 工厂函数，标准化工具定义。

src/tools/
├── builtin/ # 内置工具实现
│ ├── file/ # 文件工具 (Read, Write, Edit)
│ ├── search/ # 搜索工具 (Glob, Grep)
│ ├── shell/ # Shell 工具 (Bash)
│ ├── web/ # 网络工具 (WebFetch, WebSearch)
│ ├── plan/ # Plan 工具 (EnterPlanMode, ExitPlanMode)
│ ├── todo/ # TODO 工具 (TodoWrite)
│ ├── task/ # 任务工具 (Task)
│ └── system/ # 系统工具 (Skill, SlashCommand)
├── registry/ # 工具注册表
│ └── ToolRegistry.ts # 工具管理中心
├── execution/ # 执行管道
│ └── ExecutionPipeline.ts
├── core/ # 核心工具
│ └── createTool.ts # 工具工厂函数
├── types/ # 类型定义
│ └── ToolTypes.ts
└── validation/ # 参数验证
└── zodSchemas.ts

### 6.2 工具类型定义

```ts
// src/tools/types/ToolTypes.ts

export interface Tool<TParams = unknown> {
  // 基本信息
  readonly name: string // 工具唯一名称
  readonly displayName: string // 显示名称
  readonly kind: ToolKind // 工具类型
  readonly isReadOnly: boolean // 是否只读
  readonly isConcurrencySafe: boolean // 是否并发安全
  readonly strict: boolean // 是否启用结构化输出

  // 描述和元数据
  readonly description: ToolDescription
  readonly version: string
  readonly category?: string
  readonly tags: string[]

  // 核心方法
  getFunctionDeclaration(): FunctionDeclaration // 生成 LLM 工具定义
  build(params: TParams): ToolInvocation<TParams> // 构建执行实例
  execute(params: TParams, signal?: AbortSignal): Promise<ToolResult>

  // 权限相关
  extractSignatureContent?: (params: TParams) => string
  abstractPermissionRule?: (params: TParams) => string
}
```

#### 6.2.1 ToolKind 枚举

工具根据其副作用程度分为三种类型，这直接决定了系统的权限审批策略：

| 类型     | 枚举值     | 说明                       | 示例工具               |
| :------- | :--------- | :------------------------- | :--------------------- |
| **只读** | `ReadOnly` | 无副作用，不修改系统状态   | `Read`, `Glob`, `Grep` |
| **写入** | `Write`    | 修改文件内容               | `Edit`, `Write`        |
| **执行** | `Execute`  | 执行命令，可能有持久化影响 | `Bash`, `Task`         |

#### 6.2.2 ToolResult 接口

为了平衡 LLM 的理解需求与用户的阅读体验，我们将执行结果分离：

- **`llmContent`**：传递给 LLM 的完整数据（可能包含数千行代码）。
- **`displayContent`**：显示在终端的简洁摘要（如：`✅ 成功读取 app.ts (1000 行)`）。

### 6.3 工具工厂函数与参数验证

#### 6.3.1 createTool 工厂函数

我们不直接编写工具类，而是通过 `createTool` 函数进行声明式定义。它强制要求提供 `name`、`schema` 和 `description`，并自动生成 LLM 所需的 JSON Schema。

```typescript
// src/tools/core/createTool.ts
export function createTool<TSchema extends z.ZodType>(
  config: ToolConfig<TSchema>
): Tool<z.infer<TSchema>> {
  const jsonSchema = zodToJsonSchema(config.schema)
  return {
    ...config,
    getFunctionDeclaration() {
      return {
        name: config.name,
        description: config.description.short,
        parameters: jsonSchema
      }
    },
    async execute(params, context) {
      const validated = config.schema.parse(params) // 运行时验证
      return config.execute(validated, context)
    }
  }
}
```

#### 6.3.2 使用 Zod 进行参数验证

通过 Zod，我们实现了“定义即验证”：

1. **类型推断**：自动生成 TypeScript 类型。
2. **运行时拦截**：防止 LLM 传入格式错误的参数。
3. **文档生成**：Zod 的描述信息会直接进入提示词。

### 6.4 核心内置工具实现

#### 6.4.1 Read 工具：带行号的读取

为了方便 LLM 定位，`Read` 工具返回的内容会自动带上行号（`cat -n` 格式），并支持 `offset` 和 `limit` 分页读取大文件。

#### 6.4.2 Edit 工具：精确字符串替换

不同于全量覆盖，`Edit` 工具采用 `old_string` -> `new_string` 的替换模式。这能极大地节省 Token，并减少模型因重写整个文件而引入无关错误的风险。

#### 6.4.3 Bash 工具：持久化 Shell 会话

`Bash` 工具允许 Agent 执行任意命令。我们支持后台运行模式（`run_in_background`），方便启动开发服务器或长时间运行的测试任务。

### 6.5 工具注册表与发现

`ToolRegistry` 是工具管理的中心。它不仅管理内置工具，还支持通过 **MCP (Model Context Protocol)** 协议动态发现外部工具。它维护了分类索引和标签索引，使得系统可以根据当前模式（如 Plan 模式）快速筛选出可用的工具集。

### 6.6 踩坑记录：工具系统开发中的陷阱

#### 1. Edit 工具的“唯一性”陷阱

**问题**：如果 `old_string` 在文件中出现多次，直接替换可能会改错地方。
**对策**：强制唯一性检查。如果匹配项 > 1 且未开启 `replace_all`，直接报错并要求 LLM 提供更多上下文。`编辑操作宁可失败也不要猜测用户意图`。

#### 2. Read-Before-Write 验证

**问题**：LLM 可能会凭“记忆”修改文件，但文件可能已被外部修改。
**对策**：引入 `FileAccessTracker`。在执行 `Edit` 前，检查该文件在当前会话中是否已被 `Read`。如果没有，强制要求模型先读取。

#### 3. 智能引号导致的匹配失败

**问题**：从网页复制的代码可能包含“智能引号”（`“`），导致 `old_string` 匹配失败。
**对策**：在匹配前进行引号标准化（Normalize Quotes），将所有变体统一转换为标准单/双引号。

#### 4. Bash 输出截断

**问题**：`npm install` 等命令可能产生数 MB 的日志，撑爆上下文。
**对策**：实现智能截断策略。保留输出的前 60% 和后 30%，中间部分替换为 `... [X characters truncated] ...`。

## 执行管道、权限模型、确认机制

### 7.1 执行管道概述

如果说工具是 Agent 的“手脚”，那么执行管道就是它的“中枢神经系统”。它确保了 Agent 的每一次行动都是一个有组织、有纪律、可预测的过程。

#### 7.1.1 为什么需要执行管道

在设计之初，为了处理安全、灵活性和一致性，Blade 将工具执行解耦为七个独立阶段：

- **安全与信任**：通过权限检查和用户确认确保操作安全。
- **灵活性**：通过工具发现和钩子（Hook）机制实现高度扩展。
- **一致性**：通过参数验证和标准执行确保结果健壮。

#### 7.1.2 管道架构

[`ExecutionPipeline`](src/tools/execution/ExecutionPipeline.ts) 类负责管理整个生命周期，将关注点分离为独立的阶段。

```typescript
// src/tools/execution/ExecutionPipeline.ts
export class ExecutionPipeline extends EventEmitter {
  private stages: PipelineStage[]

  constructor(private registry: ToolRegistry) {
    super()
    this.stages = [
      new DiscoveryStage(this.registry), // 1. 工具发现
      new PermissionStage(), // 2. 权限检查
      new HookStage(), // 3. Pre-Tool Hooks
      new ConfirmationStage(), // 4. 用户确认
      new ExecutionStage(), // 5. 实际执行
      new PostToolUseHookStage(), // 6. Post-Tool Hooks
      new FormattingStage() // 7. 结果格式化
    ]
  }
}
```

### 7.2 七阶段详解

#### 7.2.1 Stage 1: Discovery（工具发现）

[`DiscoveryStage`](src/tools/execution/PipelineStages.ts) 负责从注册表中查找对应的工具实例。如果工具不存在，则立即中止执行。

#### 7.2.2 Stage 2: Permission（权限检查）

[`PermissionStage`](src/tools/execution/PipelineStages.ts) 是最核心的安全层。它会：

1. **参数验证**：执行 Zod Schema 校验。
2. **规则匹配**：检查 `allow/ask/deny` 规则。
3. **模式覆盖**：根据当前的 [`PermissionMode`](src/config/types.ts)（如 YOLO 或 PLAN）调整决策。

#### 7.2.3 Stage 3: Hook (Pre)

[`HookStage`](src/hooks/HookStage.ts) 允许在工具执行前注入逻辑。Hook 可以决定 `allow`（允许）、`ask`（需确认）或 `deny`（拒绝），甚至可以修改工具的输入参数。

#### 7.2.4 Stage 4: Confirmation（用户确认）

对于被标记为 `needsConfirmation` 的危险操作，[`ConfirmationStage`](src/tools/execution/PipelineStages.ts) 会挂起执行并请求用户审批。

#### 7.2.5 Stage 5: Execution（实际执行）

[`ExecutionStage`](src/tools/execution/PipelineStages.ts) 调用工具的 `execute` 方法，处理异步逻辑并捕获执行错误。

#### 7.2.6 Stage 6: Hook (Post)

[`PostToolUseHookStage`](src/hooks/PostToolUseHookStage.ts) 在执行完成后运行，用于清理资源或记录审计日志。

#### 7.2.7 Stage 7: Formatting（结果格式化）

[`FormattingStage`](src/tools/execution/PipelineStages.ts) 统一输出格式，确保 `llmContent`（给模型看）和 `displayContent`（给用户看）都符合预期。

### 7.3 权限模式

Blade 定义了四种 [`PermissionMode`](src/config/types.ts) 来平衡安全与效率：

| 模式          | ReadOnly 工具 | Write 工具  | Execute 工具 | 说明                         |
| :------------ | :------------ | :---------- | :----------- | :--------------------------- |
| **DEFAULT**   | ✅ 自动批准   | ❓ 需确认   | ❓ 需确认    | 标准安全模式                 |
| **AUTO_EDIT** | ✅ 自动批准   | ✅ 自动批准 | ❓ 需确认    | 信任文件修改，不信任命令执行 |
| **YOLO**      | ✅ 自动批准   | ✅ 自动批准 | ✅ 自动批准  | 完全信任模式，风险自担       |
| **PLAN**      | ✅ 自动批准   | ❌ 拒绝     | ❌ 拒绝      | 只读调研模式，禁止任何副作用 |

### 7.4 权限规则配置

用户可以在配置文件中通过字符串模式定义细粒度的权限规则：

```json
{
  "permissions": {
    "allow": ["Read(**/*.ts)", "Bash(npm:*)"],
    "deny": ["Bash(rm -rf:*)", "Write(/etc/*)"]
  }
}
```

[`PermissionChecker`](src/config/PermissionChecker.ts) 会根据这些规则对工具调用的“签名”（如 `Bash(npm install)`）进行匹配。

### 7.5 确认机制

当需要用户确认时，系统会通过 [`ConfirmationHandler`](src/tools/types/ExecutionTypes.ts) 弹出交互界面。[`ConfirmationPrompt`](src/ui/components/ConfirmationPrompt.tsx) 组件会展示：

- **操作预览**：如 `Edit` 工具的 Diff 预览。
- **风险提示**：明确告知该操作可能带来的副作用。
- **批准范围**：支持“仅此一次”或“在本会话中记住”。

### 7.6 敏感文件检测

为了防止模型意外读取或修改敏感信息，[`SensitiveFileDetector`](src/tools/validation/SensitiveFileDetector.ts) 会对路径进行扫描：

- **高敏感**：`.env`、私钥（`.pem`）、凭证文件。匹配时通常直接拒绝。
- **中敏感**：数据库文件（`.sqlite`）、日志文件。匹配时强制要求用户确认。

## 消息历史、Token 统计、自动压缩策略

- 多层存储架构（内存、持久化、缓存）
- Token 计算与阈值检测
- 智能压缩策略
- JSONL 格式的会话持久化

### 8.1 上下文管理的挑战

Coding Agent 的对话通常包含大量的文件内容和工具输出，这会导致 Token 消耗极快。为了支持长对话并节省成本，我们需要解决四个核心问题：**Token 溢出处理**、**跨会话持久化**、**流式读取性能**以及**上下文的逻辑连续性**。

为什么需要上下文管理？
Coding Agent 与 LLM 的交互有几个特点：

1. 对话可能很长 - 一次代码重构可能涉及几十个文件操作
2. 工具调用产生大量输出 - 读取文件、执行命令的结果都需要注入上下文
3. Token 有上限 - 超出限制会导致 API 调用失败
4. 需要跨会话延续 - 用户可能希望"继续昨天的工作"

### 8.2 整体架构设计

Blade 采用三层存储架构来管理上下文：

- **内存层 (MemoryStore)**：存放当前活跃会话，提供极速读写。
- **持久层 (PersistentStore)**：使用 **JSONL** 格式按项目隔离存储，支持追加写入和故障恢复。
- **缓存层 (CacheStore)**：缓存工具执行结果和压缩摘要，减少重复计算。

### 8.3 为什么选择 JSONL 格式？

不同于传统的 JSON 或 SQLite，Blade 选择 **JSON Lines (JSONL)** 作为持久化格式：

1. **追加友好**：新消息直接 `append` 到文件末尾，无需读取和重写整个文件。
2. **流式读取**：支持逐行解析，适合处理超大历史记录。
3. **人类可读**：每行都是一个独立的 JSON 对象，便于开发者直接打开调试。

### 8.4 Token 计算与阈值检测

准确的 Token 统计是触发压缩的前提。Blade 使用 [`js-tiktoken`](src/context/TokenCounter.ts) 进行计算：

- **模型对齐**：优先使用当前模型的 Encoding，降级方案为 GPT-4 Encoding、estimateTokens。
- **开销计算**：除了消息内容，还需计算 `role` 字段和 `tool_calls` 的固定开销。
- **阈值检测**：当 Token 占用达到窗口上限的 **80%** 时，系统会自动触发压缩流程。

### 8.5 智能压缩策略 (`CompactionService`)

压缩的核心是调用 LLM 生成高质量的总结
压缩不仅仅是截断，而是一个“总结 + 保留”的过程：

1.  **文件分析**：[`FileAnalyzer`](src/context/FileAnalyzer.ts) 扫描历史，找出被修改最频繁或最近提及的 5 个核心文件。
2.  **生成总结**：调用 LLM 对历史对话进行深度总结，涵盖原始意图、技术决策、已修复的错误和待办事项。
3.  **保留窗口**：保留最近 **20%** 的原始消息，确保 Agent 记得刚刚发生的细节。
4.  **孤儿消息过滤**：这是最容易踩坑的地方。压缩时必须确保保留的 `tool` 响应消息在上下文中仍能找到对应的 `assistant` 调用 ID，否则 LLM 会报错。

```typescript
// src/context/CompactionService.ts 核心逻辑片段
const availableToolCallIds = new Set(
  candidateMessages.flatMap(m => m.tool_calls?.map(tc => tc.id) || [])
)
const retainedMessages = candidateMessages.filter(
  msg => msg.role !== 'tool' || availableToolCallIds.has(msg.tool_call_id)
)

// 推荐配置
const contextConfig = {
  maxContextTokens: 128000, // 模型上下文窗口
  compressionThreshold: 0.8, // 80% 触发压缩
  retainPercent: 0.2 // 保留最近 20% 消息
}

// 计算
// 128000 * 0.8 = 102400 tokens 时触发
// 保留约 20% 的最近消息 + 总结

// 优先级排序
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// 1. 被修改的文件（Write/Edit）
// 2. 提及次数多的文件
// 3. 最近提及的文件
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// 最多包含 5 个文件
// 每个文件最多 1000 行
```

### 8.6 降级与手动触发

- **降级策略**：如果 LLM 总结失败，系统会回退到“简单截断”模式，保留最近 30% 的消息并注入警告。
- **手动压缩**：用户可以通过 `/compact` 命令随时手动触发压缩，以优化当前的 Token 占用。

## UI 系统与终端渲染：Ink 框架、Markdown 渲染与焦点管理

### 9.1 为什么选择 Ink？

传统 CLI 开发依赖于手动拼接 ANSI 转义码，布局计算复杂且状态管理混乱。**Ink** 是 React 在终端的渲染器，它允许我们使用熟悉的组件模型、Hooks 和 Flexbox 布局来构建 CLI。

| 特性         | 传统 CLI          | Ink                       |
| :----------- | :---------------- | :------------------------ |
| **组件复用** | ❌ 手动复制字符串 | ✅ React 组件化           |
| **状态管理** | ❌ 全局变量       | ✅ useState / Zustand     |
| **布局系统** | ❌ 手动计算行列   | ✅ Flexbox (Yoga)         |
| **样式系统** | ❌ 原始 ANSI 码   | ✅ 声明式 color/bold 属性 |

### 9.2 应用初始化 (`AppWrapper`)

`AppWrapper` 是 UI 的入口，负责在渲染主界面前完成环境准备：加载配置文件、初始化 Zustand Store、设置主题以及启动版本检查。

### 9.3 焦点管理机制

在复杂的 CLI 中，多个组件（如输入框、确认对话框、主题选择器）可能同时存在。Blade 使用基于 Zustand 的焦点系统来确保键盘输入被分发到正确的组件：

```typescript
export enum FocusId {
  MAIN_INPUT = 'main-input',
  CONFIRMATION_PROMPT = 'confirmation-prompt',
  THEME_SELECTOR = 'theme-selector'
}

// 在组件中通过 isActive 控制 useInput 钩子
useInput(
  (input, key) => {
    // 仅在聚焦时处理逻辑
  },
  { isActive: currentFocus === FocusId.MAIN_INPUT }
)
```

### 9.4 消息渲染系统 (`MessageRenderer`)

为了支持完整的 Markdown 格式，Blade 实现了一个轻量级的解析器，将文本拆分为不同的块进行渲染：

- **代码块**：使用 `CodeHighlighter` 进行语法高亮。
- **表格**：支持自动对齐和中文字符宽度计算。
- **列表与标题**：通过不同的缩进和颜色进行区分。

### 9.5 代码高亮与性能优化

Blade 使用 `lowlight`（highlight.js 的虚拟 DOM 版本）实现语法高亮。为了防止 LLM 输出超长代码导致终端卡顿，我们采用了**智能截断策略**：仅对当前终端可见范围内的行进行高亮处理，大幅提升了渲染性能。

### 9.6 主题系统设计

通过 `ThemeManager` 单例，Blade 支持动态切换主题。主题定义涵盖了基础颜色（Primary, Success, Error）以及精细的语法高亮配色（Keyword, String, Comment）。

### 9.7 踩坑记录：UI 开发中的陷阱

#### 1. 光标位置控制

`ink-text-input` 默认不支持程序化控制光标位置。在实现 `@` 文件补全功能时，我们不得不重写了一个自定义的 `TextInput` 组件，通过 `chalk.inverse` 手动渲染光标。

- 第三方 CLI 组件功能有限，复杂需求需自定义
- Claude Code、Gemini CLI 等都实现了自定义 TextInput

#### 2. 表格宽度计算

直接使用 `string.length` 会导致中文字符对齐出错。必须使用 `string-width` 库来计算字符的真实显示宽度，并预先移除 Markdown 标记（如 `**`）后再计算。

```ts
import stringWidth from 'string-width'
export function getPlainTextLength(text: string): number {
  // 1. 移除 Markdown 标记
  const cleanText = text
    .replace(/\*\*(.*?)\*\*/g, '$1') // **粗体**
    .replace(/\*(.*?)\*/g, '$1') // *斜体*
    .replace(/~~(.*?)~~/g, '$1') // ~~删除线~~
    .replace(/`([^`]+)`/g, '$1') // `代码`
    .replace(/\[(.*?)\]\(.*?\)/g, '$1') // [链接](url)

  // 2. 使用 string-width 正确计算中文等宽度
  return stringWidth(cleanText)
}
```

#### 3. 粘贴检测

用户粘贴大段文本时，终端会将其拆分为多个小块快速输入。我们通过“分片合并 + 启发式超时检测”机制，成功识别出粘贴行为并触发专门的预览逻辑。

```ts
// ✅ 粘贴检测机制
const PASTE_CONFIG = {
  TIMEOUT_MS: 100, // 分片合并超时
  RAPID_INPUT_THRESHOLD_MS: 150, // 快速输入阈值
  LARGE_INPUT_THRESHOLD: 300 // 大文本阈值
}

const pasteState = useRef({
  chunks: [] as string[],
  timeoutId: null as NodeJS.Timeout | null,
  firstInputTime: null as number | null
})

useInput(input => {
  const now = Date.now()
  const timeSinceFirst = pasteState.current.firstInputTime
    ? now - pasteState.current.firstInputTime
    : 0

  // 启发式判断是否为粘贴
  const isPaste =
    input.length > PASTE_CONFIG.LARGE_INPUT_THRESHOLD || // 大段文本
    input.includes('\n') || // 包含换行
    (timeSinceFirst < PASTE_CONFIG.RAPID_INPUT_THRESHOLD_MS && pasteState.current.chunks.length > 0) // 快速连续

  if (isPaste) {
    // 收集分片
    pasteState.current.chunks.push(input)
    if (!pasteState.current.firstInputTime) {
      pasteState.current.firstInputTime = now
    }

    // 重置超时
    if (pasteState.current.timeoutId) {
      clearTimeout(pasteState.current.timeoutId)
    }

    pasteState.current.timeoutId = setTimeout(() => {
      const mergedText = pasteState.current.chunks.join('')
      onPaste?.(mergedText) // 触发粘贴回调
      // 重置状态
      pasteState.current = { chunks: [], timeoutId: null, firstInputTime: null }
    }, PASTE_CONFIG.TIMEOUT_MS)
  } else {
    // 普通输入
    setValue(prev => prev + input)
  }
})
```

## MCP 协议集成：协议原理、工具发现与服务器管理

### 10.1 MCP 协议概述

**MCP (Model Context Protocol)** 是由 Anthropic 推出的开放协议，旨在标准化 AI 应用与外部数据源/工具之间的交互。如果说工具系统是 Agent 的“手脚”，那么 MCP 就是 Agent 的“USB 接口”，让它可以无缝连接到 GitHub、Google Drive、本地数据库等各种生态系统。

#### 10.1.1 核心概念

- **Server**：提供工具和资源的服务端（如 GitHub MCP Server）。
- **Client**：调用工具的客户端（即我们的 Blade Agent）。
- **Transport**：通信层，支持 **stdio**（本地进程）、**SSE**（服务器发送事件）和 **HTTP**。

### 10.2 Blade MCP 架构设计

Blade 的 MCP 实现位于 `src/mcp/` 目录下，采用分层架构：

- **`McpClient`**：对官方 SDK 的健壮封装，处理连接生命周期、指数退避重试和自动重连。
- **`McpRegistry`**：中心注册表，管理多个服务器的发现与工具冲突解决。
- **`createMcpTool`**：适配器，负责将 MCP 的 JSON Schema 转换为 Blade 内部使用的 Zod Schema。

### 10.3 McpClient：弹性连接与错误分类

为了应对不稳定的网络或本地进程崩溃，[`McpClient`](src/mcp/McpClient.ts) 实现了复杂的错误分类机制：

- **临时错误**（如超时、503）：触发指数退避重试。
- **永久错误**（如命令未找到、权限拒绝）：立即中止并报错。
- **认证错误**：触发 OAuth 2.0 流程或提示用户检查 Token。

### 10.4 McpRegistry：工具冲突解决策略

当连接多个 MCP 服务器时，可能会出现同名工具（如两个服务器都有 `search` 工具）。`McpRegistry` 采用以下策略：

1. **无冲突**：直接使用原始名称（如 `create_issue`）。
2. **有冲突**：自动添加服务器前缀（如 `github__search` 和 `local__search`）。

### 10.5 MCP Tool 转换器：JSON Schema → Zod

这是集成中最具技术挑战的部分。MCP 服务器返回的是 JSON Schema，而 Blade 的工具系统强依赖 Zod 进行类型推断。我们实现了一个递归转换器：

```typescript
// src/mcp/createMcpTool.ts 核心逻辑
function convertJsonSchemaToZod(jsonSchema: any): z.ZodSchema {
  switch (jsonSchema.type) {
    case 'string':
      return z.string()
    case 'number':
      return z.number()
    case 'object':
      const shape: any = {}
      for (const [key, val] of Object.entries(jsonSchema.properties)) {
        shape[key] = convertJsonSchemaToZod(val)
      }
      return z.object(shape)
    // ... 处理 array, boolean, enum 等
  }
}
```

### 10.6 与 Agent 执行管道的集成

MCP 工具被注册后，会像内置工具一样进入 `ExecutionPipeline`。当 LLM 调用一个 MCP 工具时：

1. **Discovery**：在注册表中找到该 MCP 工具。
2. **Permission**：由于 MCP 工具具有外部副作用，通常会触发**用户确认**。
3. **Execution**：`McpClient` 通过选定的 Transport（如 stdio）向远程服务器发送请求。
4. **Formatting**：将返回的文本、图片或资源内容格式化为 LLM 可理解的上下文。

### 10.7 /mcp 命令：实时监控

用户可以通过 `/mcp` 命令查看当前连接状态：

- **`/mcp`**：显示所有服务器的健康状况（🟢 已连接 / 🔴 错误）。
- **`/mcp tools`**：列出所有动态加载的外部工具及其描述。

### 10.8 常见 MCP Server 示例

通过在 `.blade/config.json` 中配置，你可以立即扩展 Agent 的能力：

- **GitHub**：操作 Issue、PR 和代码搜索。
- **Filesystem**：受限的本地文件访问。
- **SQLite**：直接查询本地数据库。

## 双文件配置、Zustand Store、SSOT 架构

### 11.1 配置系统概述

#### 11.1.1 配置的挑战

一个功能完备的 Coding Agent 需要管理极其复杂的配置矩阵。这些配置不仅涵盖了从底层 API 通信到上层 UI 表现的方方面面，还必须处理多层级的覆盖逻辑。

**1. 配置分类**

| 配置类别       | 包含内容                                      |
| :------------- | :-------------------------------------------- |
| **API 配置**   | 多个 LLM 提供商的 API Key、Base URL、模型名称 |
| **行为配置**   | 权限规则（Allow/Deny）、Hooks 脚本、环境变量  |
| **UI 配置**    | 终端主题、语言偏好、字体与间距                |
| **运行时配置** | CLI 传入的临时参数、当前会话的活跃状态        |

**2. 配置层级与优先级**

为了实现“全局通用、项目定制、本地私有”的灵活性，系统必须遵循严格的优先级顺序（从低到高）：

| 配置类型     | 持久化 | 优先级   | 作用域 | 说明                                 |
| :----------- | :----- | :------- | :----- | :----------------------------------- |
| **默认配置** | 否     | 1 (最低) | 全局   | 代码内置的兜底参数                   |
| **用户配置** | 是     | 2        | 全局   | `~/.blade/config.json`               |
| **项目配置** | 是     | 3        | 项目   | `.blade/settings.json` (随 Git 共享) |
| **本地配置** | 是     | 4        | 项目   | `settings.local.json` (不提交 Git)   |
| **CLI 参数** | 否     | 5 (最高) | 会话   | 如 `--yolo` 或 `--model`             |

这种多层级架构带来的核心挑战在于：**如何实现高效的深度合并（Deep Merge）**，并确保在 Agent 的非 React 核心逻辑中，能够实时、同步地获取到最终生效的配置值。
为了平衡灵活性与安全性，Blade 采用了**双文件配置架构**，将配置按用途分离：

- **`config.json` (基础配置)**：存储 API Key、模型列表、UI 主题等。通常在全局（`~/.blade/`）定义。
- **`settings.json` (行为配置)**：存储权限规则、Hooks、环境变量等。支持项目级（`./.blade/`）覆盖，并支持 `settings.local.json` 用于存放不希望提交到 Git 的本地私有设置。

### 11.2 ConfigManager：智能加载与合并

[`ConfigManager`](src/config/ConfigManager.ts) 负责启动时的初始化。它遵循 **“本地 > 项目 > 用户”** 的优先级进行深度合并：

- **数组字段**（如权限规则）：去重追加。
- **对象字段**（如环境变量）：深度合并。
- **插值支持**：支持在 JSON 中使用 `${VAR:-default}` 语法引用环境变量。

### 11.3 ConfigService：安全的持久化

配置持久化最怕“竞态条件”和“数据丢失”。[`ConfigService`](src/config/ConfigService.ts) 通过以下机制确保安全：

1. **字段路由表**：定义每个字段该写入哪个文件、是否允许持久化。
2. **Read-Modify-Write 原子操作**：写入前先读取磁盘，合并后再通过 `write-file-atomic` 写入，防止覆盖掉配置文件中的未知字段。
3. **文件互斥锁 (Mutex)**：确保同一时间只有一个操作在修改特定的配置文件。

### 11.4 Zustand 状态管理：解耦 UI 与逻辑

Blade 放弃了 React Context，转而使用 **Zustand Vanilla Store**。这是整个系统的“单一事实来源 (SSOT)”：

- **多 Slice 架构**：将状态划分为 `session`（消息）、`config`（运行时配置）、`focus`（UI 焦点）、`command`（命令队列）。
- **非 React 环境访问**：Agent 核心逻辑可以通过 `vanillaStore.getState()` 同步获取配置，无需通过 React 组件层层传递。
- **命令队列系统**：支持在 Agent 执行任务时，用户继续输入命令并进入 `pendingCommands` 队列，实现任务的流水线执行。

```typescript
// src/store/vanilla.ts 核心架构
export const vanillaStore = createStore<BladeStore>()(
  subscribeWithSelector((...a) => ({
    session: createSessionSlice(...a),
    config: createConfigSlice(...a),
    command: createCommandSlice(...a)
    // ...
  }))
)
```

### 11.5 权限检查器 (PermissionChecker)

Blade 实现了三级权限模型：**Deny (拒绝) > Allow (允许) > Ask (询问)**。

- 支持 **Glob 模式** 匹配：例如 `Read(**/*.env)` 或 `Bash(npm *)`。
- **默认安全**：未匹配任何规则的操作默认触发用户确认。

### 11.6 踩坑记录：状态管理中的陷阱

#### 1. Store 初始化竞态

在 CLI 模式下，Store 可能还未完成配置加载就被访问。我们通过 `ensureStoreInitialized` 模式，利用共享 Promise 确保所有入口点（UI、CLI、Agent）在访问数据前 Store 已就绪。

#### 2. 选择器导致的无限重渲染

在 React 中订阅对象或数组时，必须使用 `useShallow` 或常量空引用（如 `const EMPTY = []`），否则每次状态更新都会因引用变化导致 UI 剧烈闪烁。

#### 3. 配置持久化丢失

早期版本直接覆盖文件导致手动编辑的注释或未知字段丢失。现在的“读取-合并-原子写入”流程完美解决了向前兼容问题。

## 12. 高级功能与扩展

在前面的章节中，我们构建了一个功能完整的 Coding Agent。本章将探讨六种高级扩展机制：Hooks 系统、Slash Commands、Subagent 机制、Skills 系统、IDE 集成和多端扩展架构。这些特性让 Agent 从一个单一工具演变为一个可扩展的多平台开发者工具。

### 12.1 Hooks 系统

Hooks 是 Blade 最强大的扩展机制之一，允许用户在特定事件点注入自定义 Shell 命令。通过 Hooks，你可以拦截工具调用、后处理输出、注入动态上下文或自动化工作流。

Blade 支持 11 种 Hook 事件，分为四大类：

- **工具执行类**：`PreToolUse`, `PostToolUse`, `PostToolUseFailure`, `PermissionRequest`
- **会话生命周期类**：`UserPromptSubmit`, `SessionStart`, `SessionEnd`
- **控制流类**：`Stop`, `SubagentStop`
- **其他**：`Notification`, `Compaction`

### 12.2 Slash Commands

Slash Commands 是用户与 Agent 交互的快捷方式。Blade 提供了内置命令（如 `/help`, `/clear`, `/init`, `/model`）以及强大的 `/git` 命令，结合 AI 能力提供智能提交和代码审查。

### 12.3 自定义 Slash Commands

除了内置命令，Blade 支持用户通过简单的 Markdown 文件定义自己的命令。这些命令支持参数插值、Bash 嵌入和文件引用，并能通过 `SlashCommand` 工具让 AI 主动调用。

### 12.4 Subagent 机制

Subagent 是专门化的子 Agent，拥有专属的系统提示词和限定工具集。通过 `Task` 工具，主 Agent 可以委派独立子任务（如专门的代码审查或测试编写）给 Subagent 并行执行。

### 12.5 Skills 系统

Skills 允许开发者定义专业能力和工作流规范。它采用**渐进式披露**设计：初始仅加载元数据以节省 Token，当 AI 判断需要时再按需加载完整的 `SKILL.md` 指令。

### 12.6 扩展模式总结

| 特性          | Hooks            | Slash Commands   | Subagents  | Skills          |
| :------------ | :--------------- | :--------------- | :--------- | :-------------- |
| **触发方式**  | 自动（事件驱动） | 手动（用户输入） | 手动/自动  | 自动（AI 判断） |
| **主要用途**  | 拦截/增强/自动化 | 系统操作         | 专业化任务 | 最佳实践/方法论 |
| **AI 可调用** | ❌               | ❌               | ✅         | ✅              |

### 12.7 IDE 集成

Coding Agent 不应局限于终端。Blade 通过 WebSocket 和 JSON-RPC 协议与 IDE（如 VS Code）通信，支持在编辑器中打开文件、获取选中代码、查看诊断信息以及预览 Diff。同时，Blade 正在规划对 **ACP (Agent Client Protocol)** 标准协议的支持。

### 12.8 多端扩展展望

基于核心层的抽象，Blade 的能力可以轻松扩展到 Web 端（SSE 流式输出）、原生 IDE 插件以及 Electron 应用，实现“一套核心，多端运行”的架构。

## 总结

至此，我们完成了全部内容，从零开始，一步步构建了一个功能完整的 CLI Coding Agent。

在实现 Blade 的过程中，最深刻的体会是：**提示词和模型，决定了 Agent 的本质。**

- **提示词定性格**：System Prompt 定义了 Agent 的行为边界和决策倾向。
- **模型定大脑**：工程可以无限逼近模型的天花板，但很难突破它。

**工程的价值在于：让模型的能力充分释放，而不是被糟糕的设计所浪费。** 无论是 Agentic Loop、工具系统、权限控制还是上下文管理，每一个模块都是在帮助模型更好地发挥其潜力。

从 0 到 1 实现 Blade 项目，让我彻底摆脱了“调 API”式的浅层认知。我不再把 AI 看作一个黑盒，而是开始思考如何为它设计合适的工具系统和执行管道。编程的范式确实在改变，Blade 项目也将持续演进，探索更多有趣的方向。

---

- 对于编程领域，如果压缩的话会不会把一些关键代码给压缩，导致下轮对话反而缺失关键信息了？平衡是怎么做到的
  原则是：保留 结构化信息（文件路径、函数签名、关键决策），丢弃 冗余内容（重复读取、verbose 输出）， 用自然语言总结 上下文要点

---

## JSONL

在构建 Coding Agent 时，上下文的持久化存储至关重要。Blade 选择 **JSONL (JSON Lines)** 而非传统的 JSON 或数据库，是基于性能、内存效率和鲁棒性的深思熟虑。

### 8.3 深度解析：为什么选择 JSONL 格式？

#### 8.3.1 什么是 JSONL？

JSONL（也称为 ndjson）是一种文本格式，要求每一行都是一个有效的 JSON 对象，并以换行符（`\n`）分隔。

```jsonl
{"role": "user", "content": "Hello", "timestamp": 1704240000}
{"role": "assistant", "content": "Hi there!", "timestamp": 1704240005}
```

#### 8.3.2 核心优势

1. **极速追加 (Append-Only Performance)**

   - **JSON**: 存储为数组 `[...]`。每次新增消息，必须先将整个文件读入内存，解析成对象，`push` 新元素，再序列化并覆盖写入。当历史记录达到数万行代码时，这种 I/O 开销是不可接受的。
   - **JSONL**: 新消息只需以 `a` (append) 模式直接写入文件末尾。无论文件多大，写入操作的时间复杂度始终是 $O(1)$。

2. **内存效率与流式处理 (Streaming)**

   - **JSON**: 必须一次性加载整个文件。如果会话包含大量 `Read` 工具输出，内存占用会迅速飙升。
   - **JSONL**: 支持逐行读取。我们可以使用 Node.js 的 `readline` 接口，像流水一样处理历史记录，仅在内存中保留当前需要的片段。

3. **容错性 (Robustness)**
   - **JSON**: 如果文件末尾因为程序崩溃少了一个 `]`，整个文件将无法解析。
   - **JSONL**: 即使某一行损坏，解析器只需跳过该行即可恢复，不会导致整个会话历史丢失。

#### 8.3.3 与其他格式对比

| 特性         | JSON        | SQLite   | **JSONL**               |
| :----------- | :---------- | :------- | :---------------------- |
| **写入速度** | 慢 (需重写) | 快       | **极快 (直接追加)**     |
| **读取方式** | 全量加载    | SQL 查询 | **流式/逐行**           |
| **人类可读** | 是          | 否       | **是**                  |
| **并发支持** | 差          | 强       | **中 (支持多进程追加)** |
| **复杂度**   | 极低        | 高       | **低**                  |

### 8.3.4 Blade 中的实现示例

在 `PersistentStore.ts` 中，我们利用 Node.js 的文件流实现高效存取：

````typescript
// ...existing code...
### 8.3 为什么选择 JSONL 格式？

// ...existing code...
#### 8.3.4 代码实现参考

```typescript
import fs from 'node:fs/promises';
import readline from 'node:readline';

export class JsonlStore {
  /**
   * 追加写入：O(1) 复杂度
   */
  async append(filePath: string, data: any) {
    const line = JSON.stringify(data) + '\n';
    await fs.appendFile(filePath, line, 'utf-8');
  }

  /**
   * 流式读取：内存占用极低
   */
  async readLastN(filePath: string, n: number): Promise<any[]> {
    const fileStream = await fs.open(filePath, 'r');
    const rl = readline.createInterface({
      input: fileStream.createReadStream(),
      crlfDelay: Infinity
    });

    const lines: any[] = [];
    for await (const line of rl) {
      if (line.trim()) {
        lines.push(JSON.parse(line));
      }
    }
    return lines.slice(-n);
  }
}
````

### 8.3.5 适用场景建议

- **会话日志**：完美契合。
- **工具执行缓存**：非常适合。
- **结构化配置**：建议继续使用标准 **JSON**，因为配置通常较小且需要频繁的随机读写。

---

# Claude Code 实现原理深度分析

这篇文章详细剖析了如何从零构建一个类似 Claude Code 的 AI 编程助手。以下是其中最有价值的技术亮点和难点分析：

---

## 一、架构设计亮点

### 1.1 无状态 Agent 设计

```typescript
// Agent 本身不保存任何会话状态
export class Agent {
  // 配置（只读）
  private config: BladeConfig
  // 所有状态通过 context 参数传入
  async executeLoop(message: string, context: ChatContext): Promise<LoopResult>
}
```

**价值所在**：

- **可测试性**：每次测试传入独立 context，无需 mock 复杂状态
- **并发安全**：多个会话可共享同一 Agent 实例
- **容错性**：Agent 崩溃不影响历史状态，重启后可继续

### 1.2 七阶段执行管道

```
Discovery → Permission → Hook(Pre) → Confirmation → Execution → Hook(Post) → Formatting
```

**技术难点**：将工具执行解耦为独立阶段，每个阶段职责单一：

| 阶段         | 职责     | 难点                    |
| :----------- | :------- | :---------------------- |
| Discovery    | 工具查找 | MCP 动态工具的延迟加载  |
| Permission   | 权限校验 | Glob 模式匹配的性能优化 |
| Confirmation | 用户确认 | 异步挂起与恢复机制      |
| Execution    | 实际执行 | 超时控制与信号中断      |

---

## 二、核心技术难点

### 2.1 意图未完成检测（幻觉处理）

这是**最关键的工程难题之一**。小模型经常"只说不做"：

```typescript
// 检测模型说了"让我来..."但没有实际调用工具的情况
const INCOMPLETE_INTENT_PATTERNS = [
  /：\s*$/, // 中文冒号结尾
  /让我(先|来|开始)/, // 中文意图词
  /Let me (first|start|check)/i
]

if (isIncompleteIntent && recentRetries < 2) {
  messages.push({
    role: 'user',
    content: '请执行你提到的操作，不要只是描述。'
  })
  continue
}
```

**解决方案权衡**：

| 方案                    | 优点     | 缺点                   |
| :---------------------- | :------- | :--------------------- |
| 重试提示                | 简单直接 | 可能陷入死循环         |
| 多模型 fallback         | 效果好   | 成本高、上下文传递复杂 |
| `tool_choice: required` | 强制调用 | 过于僵硬               |

### 2.2 上下文压缩中的孤儿消息问题

**问题本质**：压缩历史时，如果保留了 `tool` 响应但丢失了对应的 `assistant` 调用，LLM 会报错。

```typescript
// 必须成对保留
const availableToolCallIds = new Set(
  candidateMessages.flatMap(m => m.tool_calls?.map(tc => tc.id) || [])
)

// 过滤掉孤立的 tool 消息
const retainedMessages = candidateMessages.filter(
  msg => msg.role !== 'tool' || availableToolCallIds.has(msg.tool_call_id)
)
```

**这是容易被忽视但极其关键的细节**——很多开源实现会在这里出 bug。

### 2.3 Edit 工具的唯一性陷阱

```typescript
// old_string 在文件中出现多次怎么办？
const matches = content.split(old_string).length - 1

if (matches > 1 && !replace_all) {
  return {
    success: false,
    error: `Found ${matches} matches. Provide more context or use replace_all.`
  }
}
```

**设计哲学**：`编辑操作宁可失败也不要猜测用户意图`。这比"智能选择第一个匹配"更安全。

---

## 三、性能优化亮点

### 3.1 JSONL 格式选择

| 操作          | JSON                  | JSONL               |
| :------------ | :-------------------- | :------------------ |
| 追加消息      | O(n) - 需重写整个文件 | **O(1) - 直接追加** |
| 读取最后 N 条 | O(n) - 全量加载       | **O(N) - 流式读取** |
| 容错恢复      | 单个 `]` 丢失全废     | **仅跳过损坏行**    |

```typescript
// 极速追加
async append(filePath: string, data: any) {
  const line = JSON.stringify(data) + '\n';
  await fs.appendFile(filePath, line, 'utf-8'); // O(1)
}
```

### 3.2 Token 预算管理

```typescript
const contextConfig = {
  maxContextTokens: 128000,
  compressionThreshold: 0.8, // 80% 时触发
  retainPercent: 0.2 // 保留最近 20%
}

// 压缩后保留：总结 + 最近 20% 原始消息 + 5 个核心文件内容
```

**关键洞察**：不是简单截断，而是**结构化保留**：

1. 被修改最多的文件（优先级最高）
2. 最近提及的技术决策
3. 待办事项列表

---

## 四、安全设计

### 4.1 四层权限模式

| 模式          | ReadOnly | Write     | Execute   |
| :------------ | :------- | :-------- | :-------- |
| **DEFAULT**   | ✅ 自动  | ❓ 需确认 | ❓ 需确认 |
| **AUTO_EDIT** | ✅ 自动  | ✅ 自动   | ❓ 需确认 |
| **YOLO**      | ✅ 全开  | ✅ 全开   | ✅ 全开   |
| **PLAN**      | ✅ 只读  | ❌ 拒绝   | ❌ 拒绝   |

### 4.2 Read-Before-Write 校验

```typescript
// 防止 LLM 凭"记忆"修改已被外部更改的文件
class FileAccessTracker {
  hasReadInSession(filePath: string, sessionId: string): boolean
}

// Edit 工具执行前检查
if (!tracker.hasReadInSession(filePath, sessionId)) {
  return { error: '请先读取文件再修改' }
}
```

---

## 五、状态管理创新

### 5.1 Zustand Vanilla Store 解耦 UI 与逻辑

**问题**：Agent 核心逻辑运行在非 React 环境，无法使用 Context。

```typescript
// 创建独立于 React 的 Store
export const vanillaStore = createStore<BladeStore>()(
  subscribeWithSelector((...a) => ({
    session: createSessionSlice(...a),
    config: createConfigSlice(...a)
  }))
)

// Agent 类中同步获取配置
const config = vanillaStore.getState().config
```

**这解决了"双轨数据源不一致"的经典难题**——UI 用 Hook 订阅，Agent 用 `getState()` 同步读取。

---

## 六、可扩展性设计

### 6.1 MCP 协议集成：JSON Schema → Zod 转换

**难点**：MCP Server 返回 JSON Schema，但 Blade 内部依赖 Zod。

```typescript
function convertJsonSchemaToZod(jsonSchema: any): z.ZodSchema {
  switch (jsonSchema.type) {
    case 'string':
      return z.string()
    case 'object':
      const shape = {}
      for (const [key, val] of Object.entries(jsonSchema.properties)) {
        shape[key] = convertJsonSchemaToZod(val) // 递归转换
      }
      return z.object(shape)
  }
}
```

### 6.2 工具冲突自动解决

当多个 MCP Server 有同名工具时：

```typescript
// 无冲突：create_issue
// 有冲突：github__search, local__search（自动添加前缀）
```

---

## 七、总结：核心设计哲学

文章揭示的最重要洞察：

> **工程的价值在于：让模型的能力充分释放，而不是被糟糕的设计所浪费。**

这体现在：

1. **极简架构**：将复杂决策交给 LLM，而非硬编码规则引擎
2. **安全默认**：宁可多问一次确认，也不自作主张
3. **渐进式披露**：Skills 系统初始只加载元数据，按需展开
4. **容错优先**：JSONL 单行损坏不影响整体、孤儿消息过滤

这套系统展示了如何把一个"能聊天的 LLM"变成一个"能干活的全栈工程师"——关键不在于模型本身，而在于围绕它构建的工具系统、权限控制和执行管道。
