## OpenCode 源码整体解读

### 一、项目定位

OpenCode 是一个 **100% 开源的 AI 编程 Agent**，类似 Claude Code，但不绑定任何 LLM 厂商。核心特点：

- TUI 终端界面 + Web + 桌面端（Tauri）三端支持
- Client/Server 分离架构（可远程驱动）
- 内置 LSP（代码智能）、MCP（外部工具协议）支持
- 基于 **Bun 运行时** + **TypeScript** 开发

---

### 二、Monorepo 结构

```
packages/
├── opencode/       ← 🧠 核心引擎（CLI + Server + Agent 逻辑）
├── app/            ← 🖥️ Web 前端（SolidJS + Vite）
├── desktop/        ← 桌面端（Tauri）
├── console/        ← 控制台 Web 应用
├── ui/             ← 共享 UI 组件
├── sdk/            ← JS SDK
├── plugin/         ← 插件系统
├── util/           ← 工具函数
└── web/            ← 官网/Landing page
```

包管理器为 **Bun**，构建编排用 **Turborepo**。

---

### 三、核心引擎 `packages/opencode/src/`

这是最核心的部分，大约 **15+ 子模块**：

#### 1. 入口 & CLI（index.ts）

使用 **yargs** 构建命令行，核心命令：

| 命令               | 功能                |
| ------------------ | ------------------- |
| `opencode`（默认） | 启动 TUI 交互 Agent |
| `opencode serve`   | 启动 HTTP API 服务  |
| `opencode web`     | 启动 Web 界面       |
| `opencode session` | 管理对话会话        |
| `opencode models`  | 列出可用模型        |
| `opencode mcp`     | 管理 MCP 工具服务   |

#### 2. Agent 系统（agent/）

Agent 是核心执行单元，内置 7 个 Agent：

| Agent          | 模式     | 功能                                 |
| -------------- | -------- | ------------------------------------ |
| **build**      | primary  | 默认 Agent，全能开发，读写执行皆可   |
| **plan**       | primary  | 只读分析模式，禁止编辑，适合代码探索 |
| **general**    | subagent | 复杂搜索/多步任务的子 Agent          |
| **explore**    | subagent | 快速只读代码探索                     |
| **compaction** | internal | 上下文压缩（对话太长时自动触发）     |
| **title**      | internal | 自动生成会话标题                     |
| **summary**    | internal | 生成代码变更摘要                     |

Agent 的关键属性：`name`、`model`（使用的LLM）、`permission`（权限规则）、`prompt`（系统提示词）。用户可在 `opencode.json` 中自定义 Agent。

#### 3. Session 会话管理（session/）

Session 是对话的容器，支持：

- **创建/fork/列表/删除** — 完整的 CRUD
- **消息持久化** — 通过 Drizzle ORM + SQLite（WAL 模式）
- **快照恢复** — 支持 revert 到任意检查点
- **事件驱动** — 通过 Bus 发布 `Created/Updated/Deleted/Diff` 事件

#### 4. LLM 流式调用（session/llm.ts）

这是 Agent 与 LLM 交互的核心，处理流程：

```
构造 System Prompt → 插件钩子变换 → 过滤工具（按权限）→
选择模型（含 variant）→ 发起 SSE 流 → 解析 Tool Calls → 循环执行
```

关键特性：

- 使用 **Vercel AI SDK**（`ai` 包）统一各厂商接口
- 支持 **缓存感知 prompt 结构**（两段式 system prompt）
- 集成 **OpenTelemetry** 遥测
- **Doom Loop 检测**：同一工具连续被调用 3 次，自动暂停请求确认

#### 5. Tool 工具系统（tool/）

Agent 可调用的工具集，每个工具定义：`id` + `parameters`（Zod 校验）+ `execute` 函数。

| 工具                     | 功能                    |
| ------------------------ | ----------------------- |
| `read`                   | 读取文件                |
| `write`                  | 写入文件                |
| `edit` / `multiedit`     | 编辑文件（精确替换）    |
| `apply_patch`            | 应用 diff patch         |
| `bash`                   | 执行 Shell 命令         |
| `glob` / `grep`          | 文件搜索                |
| `websearch` / `webfetch` | 网页搜索/抓取           |
| `codesearch`             | 代码搜索                |
| `task`                   | 启动子 Agent 执行子任务 |
| `question`               | 向用户提问              |
| todo                     | 任务列表管理            |
| `lsp`                    | LSP 代码智能（实验性）  |
| `skill`                  | 调用技能脚本            |

工具注册表还支持从 `{tool,tools}/*.{js,ts}` 加载**自定义工具**和**插件工具**。

#### 6. Provider 模型提供方（provider/）

抽象了 **15+ LLM 厂商**：

```
OpenAI | Anthropic | Google (Gemini/Vertex) | Azure | AWS Bedrock |
XAI | Mistral | Groq | Cohere | Perplexity | OpenRouter |
DeepInfra | Together AI | GitHub Copilot | ...
```

- 通过 `@ai-sdk/*` 系列包实现统一接口
- 从 **models.dev** 自动拉取模型元数据（每小时刷新）
- 支持 OAuth、LiteLLM 代理、自定义加载器

#### 7. Permission 权限系统（permission/）

基于规则的细粒度权限控制：

```typescript
{ permission: "edit", pattern: "~/**", action: "deny" }   // 禁止编辑 Home 目录
{ permission: "bash", pattern: "*",    action: "ask"  }   // bash 命令需确认
```

三种 Action：`allow`（直接允许）、`deny`（直接拒绝）、`ask`（运行时询问用户）。

#### 8. Server HTTP API（server/）

基于 **Hono** 框架，暴露 REST API：

```
/session/*      会话 CRUD
/message/*      消息管理
/project/*      项目管理
/file/*         文件操作
/mcp/*          MCP 工具
/provider/*     模型提供方
/config/*       配置
/pty/*          终端伪终端
/ws             WebSocket（TUI 同步）
```

支持 CORS、Basic Auth、OpenAPI 文档。这使得 TUI、Web、桌面端都通过同一 API 交互。

#### 9. MCP 集成（mcp/）

实现 **Model Context Protocol**，支持接入外部工具：

- 传输层：Stdio / SSE / HTTP
- 支持 OAuth 认证
- 工具自动转换为 AI SDK Tool 格式

#### 10. LSP 集成（lsp/）

内置 **Language Server Protocol** 客户端：

- 支持 TypeScript、Python (Pyright)、Go、Rust、Ruby、Vue、Deno 等
- 自动检测项目根目录
- 提供符号查找、诊断等代码智能能力

#### 11. Event Bus（bus/）

实例级发布/订阅总线，类型安全（Zod 定义事件 Schema），各模块通过事件通信。

#### 12. Config 配置系统（config/）

**7 层配置合并**（低→高优先级）：

```
远程 .well-known/opencode → 全局 ~/.config/opencode/ →
OPENCODE_CONFIG 环境变量 → 项目 opencode.json →
.opencode/ 目录 → OPENCODE_CONFIG_CONTENT → 托管配置
```

#### 13. Storage 持久化（storage/）

- **SQLite** + **Drizzle ORM**（WAL 模式，高并发性能好）
- 数据库路径：`~/.opencode/opencode.db`
- 支持从旧版 JSON 文件迁移

---

### 四、核心数据流

```
用户输入 (TUI/Web/API)
    │
    ▼
  Session.create() ──────→ Database (SQLite)
    │
    ▼
  LLM.stream()
    ├─ Provider.getModel()     // 获取 LLM 实例
    ├─ ToolRegistry.tools()    // 按权限过滤可用工具
    ├─ Plugin.trigger()        // 插件钩子
    │
    ▼
  Processor (流式处理)
    ├─ reasoning-delta         // 思考过程
    ├─ tool-call → execute()   // 调用工具
    │   ├─ Permission.check()  // 权限校验
    │   └─ Tool.execute()      // 工具执行
    ├─ tool-result             // 工具返回
    └─ 循环直到 LLM 结束
    │
    ▼
  Bus.publish() ──→ TUI 渲染 / Web 推送 / 事件存储
```

---

### 五、关键设计模式

| 模式                   | 说明                                                      |
| ---------------------- | --------------------------------------------------------- |
| **Instance 单例**      | `Instance.state()` 按目录缓存实例，确保同一项目共享上下文 |
| **Effect 库**          | 大量使用 Effect-TS 进行错误处理和依赖注入                 |
| **Zod 校验**           | 所有配置、事件、工具参数都有运行时类型校验                |
| **事件总线**           | 模块间通过 Bus 解耦通信                                   |
| **分层权限**           | 默认→用户配置→Agent 专属，逐层覆盖                        |
| **插件化**             | 工具、Agent、LSP、命令都可通过插件扩展                    |
| **Client/Server 分离** | 核心逻辑在 Server，TUI/Web/Desktop 都是客户端             |

---

### 六、快速上手建议

1. **从 index.ts 开始**：理解 CLI 入口和命令注册
2. **看 agent/agent.ts**：理解 Agent 定义和内置 Agent
3. **跟踪 session/llm.ts + session/processor.ts**：这是 Agent 执行的核心循环
4. **阅读 tool/ 下任意一个工具**（如 `bash.ts`）：理解工具如何定义和执行
5. **看 config/**：了解配置如何加载和合并
6. **运行 `bun run dev`** 启动本地开发
