## Kimi Code CLI 源码整体解读

### 1. 项目定位

Kimi Code CLI 是 月之暗面(MoonshotAI) 开源的 **终端 AI Agent**，类似于 Claude Code / Codex CLI，能在终端中读写代码、执行 shell、搜索网页，并自主规划执行。

---

### 2. 整体架构 (分层)

```
┌─────────────────────────────────────────────────────┐
│                  CLI / UI 层                         │
│   cli/  ──  typer 命令行入口                         │
│   ui/shell/  ──  终端交互式 Shell (prompt-toolkit)    │
│   ui/acp/  ──  ACP 协议服务端                        │
│   ui/print/  ──  非交互式输出                        │
│   web/  ──  Web UI (FastAPI + WebSocket)             │
├─────────────────────────────────────────────────────┤
│                  Wire 通讯层                         │
│   wire/  ── Soul↔UI 之间的 spmc 消息通道             │
│   wire/types.py  ── 消息类型定义(Turn/Step/ToolCall) │
├─────────────────────────────────────────────────────┤
│                  Soul 层 (Agent 核心)                │
│   soul/kimisoul.py  ── 主循环：turn → step → tool    │
│   soul/agent.py  ── Agent 加载、Runtime、子 Agent     │
│   soul/context.py  ── 对话历史管理 (jsonl 持久化)     │
│   soul/compaction.py  ── 上下文压缩                  │
│   soul/approval.py  ── 工具执行审批机制              │
│   soul/toolset.py  ── 工具集管理 (含 MCP)            │
│   soul/denwarenji.py  ── D-Mail 时间回溯（检查点回退）│
├─────────────────────────────────────────────────────┤
│                  Tools 层                            │
│   tools/shell/  ── Shell 命令执行                    │
│   tools/file/  ── ReadFile/WriteFile/Grep/Glob/...   │
│   tools/web/  ── SearchWeb / FetchURL                │
│   tools/multiagent/  ── Task (子 Agent 委托)          │
│   tools/ask_user/  ── 向用户提问                     │
│   tools/todo/  ── Todo 列表管理                      │
│   tools/think/  ── 内部推理                          │
│   tools/dmail/  ── D-Mail (回退到旧检查点)           │
├─────────────────────────────────────────────────────┤
│               底层库 (workspace packages)            │
│   kosong  ── LLM 抽象层 (generate/step/ChatProvider) │
│   kaos (pykaos)  ── 文件系统抽象 (本地/SSH)          │
│   kimi-code  ── (扩展包)                             │
└─────────────────────────────────────────────────────┘
```

---

### 3. 核心模块详解

#### 3.1 入口 — `cli/__init__.py`

- 使用 **typer** 框架，入口函数 `kimi()` 定义在 cli/**init**.py
- pyproject.toml 中 `kimi = "kimi_cli.cli:cli"` 注册为命令行脚本
- 支持 `--model`、`--yolo`、`--session`、`--mcp-config-file` 等参数
- 子命令组：`kimi mcp`（MCP 管理）、`kimi web`（Web UI）、`kimi info`（信息查询）

#### 3.2 应用初始化 — app.py → `KimiCLI`

- `KimiCLI.create()` 是核心工厂方法，负责：
  1. 加载配置 (`Config` from TOML)
  2. 创建 OAuth 管理器
  3. 解析模型/Provider（支持 kimi/openai/anthropic/gemini/vertexai）
  4. 调用 `create_llm()` 创建 LLM 实例
  5. 创建 `Runtime` → 加载 Agent → 构建 `KimiSoul`
  6. 启动对应的 UI（Shell/ACP/Print/Web）

#### 3.3 配置系统 — config.py

- Pydantic BaseModel 层级结构：`Config` → `LLMProvider` / `LLMModel` / `LoopControl` / `Services`
- 支持 TOML 配置文件 + 环境变量覆盖 (`KIMI_API_KEY` 等)
- `LoopControl`：控制 `max_steps_per_turn`(100)、`max_retries_per_step`(3)、`reserved_context_size`(50000) 等

#### 3.4 LLM 抽象 — llm.py + `kosong` 包

- llm.py 定义 `LLM` dataclass，封装 `ChatProvider` + 能力集合 (`image_in/video_in/thinking`)
- **kosong 包**是独立的 LLM 抽象层：
  - `kosong.generate()` — 流式调用 LLM，合并 streamed parts 为完整 Message
  - `kosong.step()` — 在 generate 基础上加入工具调度
  - `ChatProvider` 接口 — 支持 Kimi/OpenAI Legacy/OpenAI Responses/Anthropic/Gemini/VertexAI
  - `Toolset` / `SimpleToolset` — 工具集抽象

#### 3.5 Agent 系统 — `soul/agent.py` + agentspec.py

- **AgentSpec**：YAML 声明式 Agent 定义，支持 `extend` 继承
  - 默认 Agent 包含 14 个工具 + 1 个 coder 子 Agent
- **Runtime**：Agent 运行时环境，包含 config/llm/session/approval/skills 等
  - 支持 `copy_for_fixed_subagent()` / `copy_for_dynamic_subagent()` 派生子 Agent Runtime
- **LaborMarket**：管理固定子 Agent（YAML 声明）和动态子 Agent（运行时创建）
- **load_agent()**：从 YAML → 加载系统提示词（Jinja2 模板渲染）→ 注册工具 → 加载子 Agent

#### 3.6 Soul 主循环 — `soul/kimisoul.py`

这是整个 Agent 的 **核心引擎**：

```python
run(user_input)
  → TurnBegin
  → _turn(user_message)
      → checkpoint()
      → append_message(user_message)
      → _agent_loop()          # 核心循环
          while steps < max_steps:
              → compaction 检测 (上下文超限时压缩)
              → kosong.step()  # 调 LLM + 工具执行
              → 如果无 tool_calls → 结束
              → 处理 steer 消息（用户实时干预）
              → D-Mail 检测（检查点回退）
  → TurnEnd
```

关键机制：

- **Compaction**：当 context_tokens + reserved ≥ max_context_size 时，用 LLM 生成摘要压缩历史
- **Approval**：Shell/WriteFile 等危险操作需用户审批，支持 yolo 模式跳过
- **D-Mail (DenwaRenji)**：Agent 可发送 "D-Mail" 回退到之前的检查点重新执行（受 Steins;Gate 启发的命名）
- **Steer**：用户可在 Agent 运行中注入即时指令
- **Slash Commands**：`/skill:xxx`、`/flow:xxx` 等斜杠命令

#### 3.7 工具系统 — `tools/`

| 工具              | 文件                     | 功能                              |
| ----------------- | ------------------------ | --------------------------------- |
| `Shell`           | tools/shell/             | 执行 bash/powershell 命令，需审批 |
| `ReadFile`        | tools/file/read.py       | 读文件（支持行号范围）            |
| `WriteFile`       | tools/file/write.py      | 写文件                            |
| `StrReplaceFile`  | tools/file/replace.py    | 字符串替换编辑文件                |
| `Glob`            | tools/file/glob.py       | 文件模式搜索                      |
| `Grep`            | tools/file/grep_local.py | 正则搜索（用 ripgrep）            |
| `ReadMediaFile`   | tools/file/read_media.py | 读图片/视频                       |
| `SearchWeb`       | tools/web/search.py      | 网络搜索                          |
| `FetchURL`        | tools/web/fetch.py       | 抓取网页                          |
| `Task`            | tools/multiagent/task.py | 委托子 Agent 执行子任务           |
| `AskUserQuestion` | tools/ask_user/          | 向用户提问                        |
| `SetTodoList`     | tools/todo/              | 管理 Todo 列表                    |
| `Think`           | tools/think/             | 内部推理                          |
| `SendDMail`       | tools/dmail/             | 回退到旧检查点                    |

工具基于 kosong 的 `CallableTool2[Params]` 抽象，每个工具有配套的 `.md` 描述文件作为 LLM prompt。

#### 3.8 Wire 通讯层 — `wire/`

- **Wire**：Soul ↔ UI 之间的 **spmc (单生产者多消费者)** 消息通道
- `WireSoulSide`：Soul 端发送消息（支持合并相邻的 ContentPart）
- `WireUISide`：UI 端订阅消息
- `WireFile`：所有消息持久化到 `wire.jsonl`
- 消息类型：`TurnBegin/TurnEnd`、`StepBegin`、`ContentPart`、`ToolCall`、`ApprovalRequest` 等

#### 3.9 Session 管理 — session.py

- 每个工作目录有独立的 session 目录
- Session 包含 `context.jsonl`（对话历史）和 `wire.jsonl`（UI 消息日志）
- 支持 `SessionState` 持久化（审批状态、动态子 Agent、附加目录等）

#### 3.10 Skill 系统 — `skill/`

- 支持三级 Skill 发现：内置 → 用户级 (`~/.config/agents/skills/`) → 项目级 (`.agents/skills/`)
- 兼容多种目录约定（`.kimi/`、`.claude/`、`.codex/`）
- 支持两种类型：`standard`（文本提示词）和 `flow`（Mermaid/D2 流程图驱动的多步 Prompt Flow）

#### 3.11 ACP 支持 — `acp/`

- 实现 [Agent Client Protocol](https://github.com/agentclientprotocol/agent-client-protocol)
- 可作为 ACP Server 被 Zed、JetBrains 等 IDE 调用

#### 3.12 底层包

- **kosong** (`packages/kosong`)：LLM 抽象层，独立可用。核心是 `generate()` (流式生成) + `step()` (生成+工具调度)，支持 Kimi/OpenAI/Anthropic/Gemini 等 Provider
- **kaos (pykaos)** (`packages/kaos`)：文件系统抽象，支持本地 + SSH 远程文件操作，通过 `KaosPath` 统一路径

---

### 4. 数据流概览

```
用户输入 → Shell UI (prompt-toolkit)
         → KimiSoul.run()
         → Context 追加消息
         → kosong.step(ChatProvider, system_prompt, toolset, history)
             → ChatProvider.generate() 流式调 LLM API
             → 解析 tool_calls
             → Toolset.handle() 执行工具 (需要时走 Approval)
         → Wire.send() 向 UI 推送结果
         → 继续循环直到无 tool_calls 或达到 max_steps
         → UI 渲染输出 (rich)
```

---

### 5. 关键设计亮点

1. **声明式 Agent 定义**：YAML + Jinja2 模板系统提示词，支持继承和子 Agent
2. **D-Mail 时间回溯**：类似 Steins;Gate 概念，Agent 可以"回到过去"重新从某个检查点执行
3. **Prompt Flow**：Skill 支持 Mermaid/D2 流程图语法定义多步执行流程
4. **Wire 消息协议**：UI 和核心完全解耦，同一个 Soul 可以挂接 Shell/ACP/Web 等不同前端
5. **上下文压缩**：自动检测 token 使用率，超限时用 LLM 生成摘要压缩历史
6. **多 Provider 支持**：通过 kosong 抽象层支持 Kimi/OpenAI/Anthropic/Gemini/VertexAI
