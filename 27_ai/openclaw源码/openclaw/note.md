## OpenClaw 源码整体解读

### 一句话概括

OpenClaw 是一个 **本地优先的个人 AI 助手网关**，用 TypeScript (ESM) 编写，通过 WebSocket 控制面连接 20+ 消息平台，将用户消息路由给 AI 模型，并赋予模型工具调用能力。

---

### 整体架构

```
消息平台 (WhatsApp/Telegram/Discord/Slack/Signal/iMessage/IRC/...)
                │
                ▼
┌─────────────────────────────────────┐
│         Gateway (控制面)              │
│    ws://127.0.0.1:18789             │
│                                     │
│  ┌──────────┐  ┌──────────────────┐ │
│  │ Channels │  │  Agent Runtime   │ │
│  │ 消息收发  │→ │ (PI嵌入式Agent)  │ │
│  └──────────┘  └───────┬──────────┘ │
│                        │            │
│  ┌─────────┐  ┌───────▼──────────┐ │
│  │ Sessions │  │   Tools System   │ │
│  │ 会话状态  │  │ exec/browser/web │ │
│  └─────────┘  └──────────────────┘ │
│                                     │
│  ┌──────────┐  ┌──────────────────┐ │
│  │ Plugins  │  │  Model Providers │ │
│  │ 扩展系统  │  │ OpenAI/Claude/…  │ │
│  └──────────┘  └──────────────────┘ │
└─────────────────────────────────────┘
                │
     ┌──────────┼──────────┐
     ▼          ▼          ▼
   macOS     iOS/Android   CLI
   桌面App    移动Node      命令行
```

---

### 10 大子系统详解

#### 1. 入口 & CLI (`src/entry.ts` → `src/index.ts` → `src/cli/`)

- src/entry.ts — 进程启动入口，处理 Node 参数、环境变量、进程 respawn
- src/index.ts — 加载 `.env`、初始化日志、调用 `buildProgram()` 构建命令树
- src/cli/program.ts — 基于 Commander.js 的命令体系

**主要命令**：`openclaw gateway start` | `config` | `channels` | `plugins` | `skills` | `models` | `nodes` | `daemon` | `agent` | `message send` | `onboard` | `doctor`

#### 2. Gateway 网关 (`src/gateway/`)

核心控制面，职责：

- **WebSocket 服务** — server.impl.ts 启动 WS + HTTP 服务器
- **RPC 方法** — server-methods.ts 注册 `chat`、`config`、`plugin` 等方法
- **消息路由** — server-chat.ts 将入站消息路由到 Agent
- **频道管理** — server-channels.ts 创建 ChannelManager
- **认证** — auth.ts 支持 token/password/device 模式
- **Web UI** — control-ui.ts 提供控制台前端

#### 3. 配置系统 (`src/config/`)

以 `openclaw.yml` 为核心的配置管线（7 步处理）：

```
原始 YAML → 解析 → 文件 include → 环境变量替换 → CLI 覆盖 → 默认值填充 → Zod 校验
```

- config.ts — `loadConfig()` 主入口
- schema.ts — 完整 Zod schema，定义所有配置项
- types.ts — `OpenClawConfig` 主类型（agents/channels/models/tools/plugins）

#### 4. Channel 消息频道 (`src/channels/` + `extensions/`)

每个消息平台是一个 `ChannelPlugin`，实现统一接口：

```typescript
interface ChannelPlugin {
  connect(config): Promise<void>
  disconnect(): Promise<void>
  send(message): Promise<SendResult>
  setTyping(sessionKey, isTyping): Promise<void>
}
```

**内置频道**：Telegram、WhatsApp (Baileys)、Discord (discord.js)、Slack (Bolt)、Signal (signal-cli)、iMessage、IRC、Google Chat、LINE 等 20+ 个平台

**消息流向**：`平台 → Channel Transport → SessionEnvelope 标准化 → 路由到 Agent → AI 回复 → channel.send() 发出`

#### 5. Agent 运行时 (`src/agents/`)

AI 推理核心，包装了嵌入式 PI Agent：

- pi-embedded-runner/run.ts — `runEmbeddedPiAgent()` 主推理函数
- pi-embedded-runner/system-prompt.ts — 系统提示词构建
- pi-embedded-runner/compact.ts — 上下文压缩（摘要化）
- pi-embedded-runner/history.ts — 消息历史管理

**推理流程**：`构建 system prompt → 收集工具列表 → 调用模型 API → 流式输出 → 处理 tool_use → 循环直到 end_turn`

#### 6. Model Provider 模型提供商 (`src/agents/models-config.providers.*`)

统一多家 AI 提供商接口：

| 提供商                         | 说明                   |
| ------------------------------ | ---------------------- |
| OpenAI                         | GPT-4o, o1             |
| Anthropic                      | Claude 3.5 Sonnet/Opus |
| Google Gemini                  | Flash/Pro              |
| Azure OpenAI                   | 企业部署               |
| Ollama / vLLM / LM Studio      | 本地模型               |
| OpenRouter                     | 聚合代理               |
| Amazon Bedrock                 | AWS 托管               |
| 字节/Minimax/通义千问/百度千帆 | 国内模型               |

支持 **API Key 轮换**（auth-profiles.ts）、**模型 Failover**（model-fallback.ts）和**自动发现**（Ollama/Bedrock）。

#### 7. Tools 工具系统 (`src/agents/tools/` + `src/agents/pi-tools.*`)

Agent 可调用的能力：

| 工具                              | 功能               |
| --------------------------------- | ------------------ |
| `exec`                            | 执行命令行         |
| `browser`                         | 自动化浏览器       |
| `web_search` / `web_fetch`        | 网络搜索和抓取     |
| `memory`                          | 语义记忆检索       |
| `canvas`                          | 可视化画布（A2UI） |
| `sessions_send` / `sessions_list` | 跨 session 通讯    |
| `image` / `pdf` / `tts`           | 多媒体处理         |

**工具策略管线** 8 层过滤：Provider 兼容 → 消息类型 → 模型兼容 → 权限控制 → 白名单 → 黑名单 → Profile → 循环检测

#### 8. Sessions 会话系统 (`src/sessions/`)

- 每个对话是一个 `SessionEntry`（sessionKey + lastChannel + to）
- YAML 持久化到 `~/.openclaw/sessions/`
- 支持上下文裁剪、会话压缩、跨频道回复

#### 9. Plugin 插件系统 (`extensions/` + `src/plugins/`)

插件可扩展 5 大维度：

| 维度         | 接口                     |
| ------------ | ------------------------ |
| 频道         | `ChannelPlugin`          |
| 工具         | `ToolFactory`            |
| 模型提供商   | `ProviderPlugin`         |
| Gateway 方法 | `GatewayRequestHandlers` |
| CLI 命令     | Commander.js 注册        |

插件生命周期：`发现 → 加载 package.json → 校验配置 → 动态 import → 调用 register* 函数 → 注册到 PluginRegistry`

Runtime API（`PluginRuntime`）提供 config/subagent/system/media/tts/stt/tools/channel/events/logging 等能力。

#### 10. Skills 技能系统 (`skills/`)

~40 个预置 CLI 工具包（如 `apple-notes`、`nano-pdf`、`openai-whisper`、`github`、`spotify` 等），Agent 可通过 `exec` 工具调用。

技能以 `skill.yml` 描述元数据，系统自动生成 system prompt 段告知 Agent 可用技能。

---

### 核心数据流总结

```
1. 用户在 WhatsApp/Telegram/等平台发消息
2. Channel Plugin 收到消息，标准化为 SessionEnvelope
3. Gateway 路由到对应 Agent 的 Session
4. Agent Runtime 构建 system prompt + tools + 历史消息
5. 调用 Model Provider API（OpenAI/Claude/…）
6. 模型返回文本或 tool_use
7. 若 tool_use → 执行工具 → 结果送回模型 → 循环
8. 最终回复通过 Channel Plugin 发回平台
```

### 技术栈

- **语言**：TypeScript (ESM)，Node ≥ 22
- **包管理**：pnpm monorepo（pnpm-workspace.yaml）
- **构建**：tsdown (基于 esbuild)
- **测试**：Vitest + V8 coverage
- **Lint**：Oxlint + Oxfmt
- **配置**：YAML (JSON5 兼容) + Zod 校验

这是一个设计精良的 monorepo 项目，核心思想是 **Gateway 作为单一控制面**，所有消息平台、AI 模型、工具能力都是可插拔的。理解了 Gateway → Channel → Agent → Tool 这条主线，就掌握了整个项目的骨架。
