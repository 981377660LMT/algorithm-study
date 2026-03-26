## GitHub Copilot Chat 源码整体解读

### 一、项目概览

这是 VS Code 的 GitHub Copilot Chat 扩展源码，一个典型的 **VS Code Extension + LLM Agent 架构**。核心是一个 **工具调用循环（Tool Calling Loop）**驱动的 AI 编程助手。

---

### 二、目录结构（三层架构）

```
src/
├── extension/     # 🎯 扩展层：VS Code UI 集成、业务逻辑
│   ├── extension/     # 入口激活
│   ├── chat/          # Chat UI 集成
│   ├── conversation/  # 会话管理、Chat Participant 注册
│   ├── prompt/        # 请求处理器、意图路由
│   ├── prompts/       # 系统提示词（prompt-tsx 组件化）
│   ├── intents/       # 意图检测 & Agent 配置
│   ├── tools/         # 工具注册 & 调用
│   ├── agents/        # 内置 Agent 定义（Explore/Plan/Ask/Edit）
│   ├── mcp/           # MCP 服务器集成
│   ├── completions/   # Inline 代码补全
│   ├── inlineEdits/   # 编辑器内联编辑
│   └── context/       # 上下文收集
├── platform/      # 🔧 平台层：跨环境抽象
│   ├── openai/        # LLM API 调用
│   ├── authentication/# GitHub 认证
│   ├── networking/    # 网络请求
│   ├── prompts/       # Prompt 路径处理
│   ├── mcp/           # MCP 协议实现
│   └── ...            # 30+ 平台服务
├── util/          # 🛠️ 工具层
│   ├── common/        # 依赖注入框架
│   └── vs/            # 从 VS Code 核心移植的工具
└── lib/           # 📦 外部依赖
```

---

### 三、核心数据流（一个请求的完整生命周期）

```
用户输入
  │
  ▼
① VS Code Chat API → createChatParticipant() 注册了 9 个参与者
  │                   (default/vscode/terminal/edits/editingSession...)
  ▼
② ChatParticipantRequestHandler  [prompt/node/chatParticipantRequestHandler.ts]
  │  - 认证检查
  │  - 变量清洗（过滤 .gitignore 文件）
  │  - 意图选择（edit/ask/explain/generate...）
  ▼
③ DefaultIntentRequestHandler   [prompt/node/defaultIntentRequestHandler.ts]
  │  - 构建 Prompt
  │  - 启动工具调用循环
  ▼
④ ToolCallingLoop              [intents/node/toolCallingLoop.ts]  ← 🔥 核心引擎
  │  循环（默认最多 15 轮）:
  │    1. 收集可用工具 + 构建 Prompt
  │    2. 发送给 LLM（带 tools 定义）
  │    3. 流式接收响应
  │    4. 如果 LLM 请求调用工具 → 执行工具 → 把结果拼回上下文 → 继续循环
  │    5. 如果 LLM 没有工具调用 → 检查停止条件 → 结束
  ▼
⑤ 响应流式输出 → ChatResponseStream → VS Code UI
```

---

### 四、六大核心模块详解

#### 1. 扩展激活入口

extension/extension/vscode/extension.ts 中 `baseActivate()`:

- 创建 **依赖注入容器**（`InstantiationService`）
- 注册 **100+ 服务**（认证、遥测、搜索、Chat 等）
- 等待实验服务初始化
- 激活所有 **Contribution**（插件化模块系统）
- 暴露 `CopilotExtensionApi` 给其他扩展

**设计模式**：VS Code 风格的 DI（`IInstantiationService` + `SyncDescriptor` 延迟实例化）

#### 2. 工具系统

tools/common/toolNames.ts 定义了所有内置工具：

| 类别         | 工具                                                                                           |
| ------------ | ---------------------------------------------------------------------------------------------- |
| **文件操作** | `read_file`, `create_file`, `replace_string_in_file`, `apply_patch`, `list_dir`, `file_search` |
| **搜索**     | `grep_search`, `semantic_search`, `search_subagent`                                            |
| **终端**     | `run_in_terminal`, `get_terminal_output`, `terminal_selection`                                 |
| **编辑器**   | `get_errors`, `get_changed_files`, `get_vscode_api`                                            |
| **测试**     | `test_failure`, `runTests`, `test_search`                                                      |
| **Agent**    | `runSubagent`, `switch_agent`, `manage_todo_list`                                              |
| **其他**     | `fetch_webpage`, `memory`, `github_repo`, `install_extension`                                  |

工具通过 `ToolRegistry.registerTool()` 注册，支持**模型特定工具**（如 Gemini 专用搜索工具可覆盖默认实现）。

#### 3. Prompt 系统（Prompt-TSX）

使用 `@vscode/prompt-tsx` —— **类 React 的 JSX 组件化 Prompt 构建**：

```tsx
// defaultAgentInstructions.tsx 中的系统提示
<InstructionMessage>
  <Tag name="instructions">
    You are a highly sophisticated automated coding agent...
    {tools[ToolName.ReadFile] && <>When using read_file, prefer reading large sections...</>}
  </Tag>
  <Tag name="toolUseInstructions">...根据可用工具动态生成指令...</Tag>
</InstructionMessage>
```

**Prompt 注册表**（promptRegistry.ts）为不同模型匹配最优 Prompt：

- anthropicPrompts.tsx — Claude 优化
- openAIPrompts.tsx — GPT 优化
- geminiPrompts.tsx — Gemini 优化
- xAIPrompts.tsx — xAI/Grok 优化

#### 4. Agent 系统

内置 Agent 通过 `.agent.md` 文件（YAML frontmatter）定义：

| Agent       | 职责                           | 默认模型         |
| ----------- | ------------------------------ | ---------------- |
| **Edit**    | 全功能代码编辑                 | 用户选择         |
| **Ask**     | 只读问答                       | 用户选择         |
| **Explore** | 快速只读代码探索子 Agent       | Claude Haiku 4.5 |
| **Plan**    | 任务规划，拆解后交给其他 Agent | 用户选择         |

Agent 配置（agentTypes.ts）：

```typescript
interface AgentConfig {
  name
  description
  argumentHint
  tools: string[] // 可用工具集
  model?: string // 首选模型
  handoffs?: AgentHandoff[] // Agent 间切换
}
```

**Handoff 机制**：Agent 间可互相交接任务（如 Plan → Edit），保留上下文。

#### 5. 会话存储

`ConversationStore` = **LRU 缓存**（1000 条）：

- `Turn`：单次交互（请求 + 响应 + 工具调用轮次）
- `Conversation`：sessionId + Turn 数组
- 响应类型：Markdown、引用、进度、FileTree、锚点

#### 6. 认证系统

`IAuthenticationService` 管理：

- GitHub OAuth 会话（`user:email` / `read:user` / `repo` 不同权限级别）
- `CopilotToken`（带过期时间 + Feature Flags）
- 权限升级弹窗（需要更高权限时引导用户）

---

### 五、关键设计模式

| 模式                  | 用途                                              |
| --------------------- | ------------------------------------------------- |
| **依赖注入**          | `IInstantiationService` + 装饰器注入，100+ 服务   |
| **Contribution 系统** | 可插拔模块，支持激活阻塞器（`activationBlocker`） |
| **Observable 响应式** | `derived()`, `autorun()` 驱动配置/状态变更        |
| **流式架构**          | `ChatResponseStream` 增量输出                     |
| **Prompt-TSX**        | 类 React 组件化的 Prompt 构建                     |
| **工具注册表**        | 静态注册 + 模型特定变体 + 运行时动态启用          |
| **Agent Handoff**     | Agent 间任务交接，保持上下文连续                  |

---

### 六、快速上手建议

1. **从 Tool Calling Loop 开始**：intents/node/toolCallingLoop.ts 是整个 Agent 模式的核心引擎
2. **看 Prompt 构建**：defaultAgentInstructions.tsx 是系统提示词的源头，理解 Copilot 对 LLM 的所有指令
3. **读官方文档**：docs/tools.md 和 docs/prompts.md 是最好的新手导读
4. **跟踪一个请求**：从 `chatParticipants.ts` → `chatParticipantRequestHandler.ts` → `defaultIntentRequestHandler.ts` → `toolCallingLoop.ts`，完整走一遍
