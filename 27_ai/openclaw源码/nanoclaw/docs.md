# NanoClaw 文档体系深度解读

> NanoClaw 的 `docs/` 目录包含 10 篇文档，覆盖了从产品设计哲学、安全模型、SDK 内部原理到技能分发架构的完整知识体系。本文按**主题层次**逐篇深度解析。

---

## 目录

1. [REQUIREMENTS.md — 产品哲学与设计决策](#1-requirementsmd--产品哲学与设计决策)
2. [SPEC.md — 系统规格说明书](#2-specmd--系统规格说明书)
3. [SECURITY.md — 安全模型](#3-securitymd--安全模型)
4. [SDK_DEEP_DIVE.md — Claude Agent SDK 逆向深度分析](#4-sdk_deep_divemd--claude-agent-sdk-逆向深度分析)
5. [nanoclaw-architecture-final.md — 技能引擎架构（三路合并方案）](#5-nanoclaw-architecture-finalmd--技能引擎架构三路合并方案)
6. [skills-as-branches.md — 技能即分支（最终方案）](#6-skills-as-branchesmd--技能即分支最终方案)
7. [nanorepo-architecture.md — 精简版技能架构](#7-nanorepo-architecturemd--精简版技能架构)
8. [docker-sandboxes.md — Docker Sandbox 部署指南](#8-docker-sandboxesmd--docker-sandbox-部署指南)
9. [APPLE-CONTAINER-NETWORKING.md — Apple Container 网络配置](#9-apple-container-networkingmd--apple-container-网络配置)
10. [DEBUG_CHECKLIST.md — 调试与排障手册](#10-debug_checklistmd--调试与排障手册)

---

## 1. REQUIREMENTS.md — 产品哲学与设计决策

这篇文档是创始人的"设计宣言"，解释了 NanoClaw 为什么存在以及核心设计理念。

### 1.1 起源：对 OpenClaw 的反思

> "That project became a monstrosity — 4-5 different processes running different gateways, endless configuration files... It's a security nightmare where agents don't run in isolated processes."

NanoClaw 是对 OpenClaw（前身 ClawBot）的**从零重写**，核心诉求是：

| OpenClaw 的问题            | NanoClaw 的解法 |
| -------------------------- | --------------- |
| 4-5 个进程 + 网关          | 单 Node.js 进程 |
| 无数配置文件               | 改代码即配置    |
| 应用层权限防护（漏洞百出） | OS 级容器隔离   |
| 50 万行代码                | 几千行代码      |
| 通用框架设计               | 个人定制软件    |

### 1.2 六大设计哲学

1. **小到能理解**：代码量小到一个人能通读。无微服务、无消息队列、无抽象层。
2. **隔离即安全**：Agent 运行在 Linux 容器中，安全边界在操作系统级别。
3. **为一个人设计**：不是框架/平台。创始人用 WhatsApp 和 Email，就支持这两个。不追求全覆盖。
4. **定制 = 改代码**：极少配置项（如触发词），其余直接改源码。
5. **AI 原生开发**：不需要安装向导（Claude Code 引导）、不需要监控仪表盘（问 Claude）、不需要日志 UI（让 Claude 读）。
6. **Skill 而非 Feature**：不在代码里堆功能，而是通过 Skill（代码变换命令）让用户 fork + 定制。

### 1.3 核心架构决策

| 决策领域 | 选择                                               |
| -------- | -------------------------------------------------- |
| 消息路由 | SQLite 轮询，触发词 `@Andy`，只处理注册群组        |
| 记忆系统 | `CLAUDE.md` 文件层次结构：全局 → 分组              |
| 会话管理 | 每组独立 Claude session，自动 compact 保持上下文   |
| 容器隔离 | 每次调用独立容器，只挂载显式路径                   |
| 定时任务 | cron/interval/once，容器内执行，结果可发消息给用户 |
| 群组管理 | 通过 Main 频道注册，SQLite + IPC 双通道            |

### 1.4 Skill 体系的核心理念

> "Skills Over Features" — 用户 fork 仓库，运行 Skill 定制代码库，最终得到只含自己需要功能的干净代码。

```bash
claude          # 进入 Claude Code
/setup          # 首次安装
/add-whatsapp   # 添加渠道
/customize      # 自定义行为
```

**关键洞察**：NanoClaw 的 Skill 不是运行时插件，而是**编译期代码变换**。每个 Skill 是指导 Claude Code 如何修改源代码的指令集。

---

## 2. SPEC.md — 系统规格说明书

这是最完整的技术规格文档（约 460 行），覆盖了系统的每个子系统。

### 2.1 架构总览

文档用 ASCII 图展示了两层架构：

```
宿主机 (Main Node.js Process):
  Channels (自注册) → SQLite → Message Loop → Group Queue → Container

容器 (Linux VM):
  Agent Runner → Claude Agent SDK → 工具集 → IPC 文件通信
```

### 2.2 渠道系统的完整设计

渠道系统是 NanoClaw 可扩展性的核心。SPEC 详细定义了：

**工厂注册表**：`Map<string, ChannelFactory>`，每个渠道在模块加载时自注册。

**Channel 接口**（8 个方法）：

| 方法            | 说明                           | 必选 |
| --------------- | ------------------------------ | ---- |
| `connect()`     | 建立连接（登录、WebSocket 等） | ✓    |
| `sendMessage()` | 发送文本消息到指定 JID         | ✓    |
| `ownsJid()`     | 判断某 JID 是否属于此渠道      | ✓    |
| `isConnected()` | 连接状态查询                   | ✓    |
| `disconnect()`  | 断开连接                       | ✓    |
| `setTyping()`   | 显示"正在输入"状态             | ✗    |
| `syncGroups()`  | 同步群组列表                   | ✗    |

**自注册流程**：

1. Skill 在 `src/channels/` 添加渠道文件，模块顶层调用 `registerChannel()`
2. `src/channels/index.ts` barrel import 触发注册
3. `main()` 中遍历注册表，有凭据的渠道自动连接，无凭据的跳过

### 2.3 记忆系统

三层文件层次的记忆系统：

```
groups/
  CLAUDE.md              ← 全局记忆（所有组可读，仅 Main 可写）
  whatsapp_main/
    CLAUDE.md            ← Main 组记忆
    notes.md             ← Agent 创建的文件
  whatsapp_family/
    CLAUDE.md            ← 家庭组记忆
```

**工作原理**：Agent 的 `cwd` 设为 `groups/{name}/`，Claude Agent SDK 的 `settingSources: ['project']` 自动加载当前目录和父目录的 `CLAUDE.md`。

### 2.4 会话管理

每个群组维护独立的 Claude session：

- session ID 存在 SQLite `sessions` 表
- 会话文件在 `data/sessions/{group}/.claude/`（JSONL 格式）
- 上下文过长时自动 compact（压缩保留关键信息）
- `resumeSessionAt` 参数确保恢复到正确的对话分支点

### 2.5 配置体系

SPEC 明确了所有配置项和约束：

- 路径必须用绝对路径（容器 volume mount 要求）
- 挂载语法注意事项：只读挂载用 `--mount "type=bind,...,readonly"` 而非 `:ro` 后缀
- `.env` 中只有认证相关变量被提取到 `data/env/env` 并挂载到容器

---

## 3. SECURITY.md — 安全模型

这篇文档定义了 NanoClaw 的**五层安全边界**。

### 3.1 信任模型

| 实体            | 信任级别 | 理由                   |
| --------------- | -------- | ---------------------- |
| Main group      | 受信     | 私人自聊天，管理员控制 |
| Non-main groups | 不受信   | 其他用户可能是恶意的   |
| 容器内 agents   | 沙箱化   | 隔离执行环境           |
| 消息内容        | 用户输入 | 可能包含提示注入攻击   |

### 3.2 五层安全边界详解

**第一层：容器隔离（主安全边界）**

- 进程隔离 — 容器内进程无法影响宿主机
- 文件系统隔离 — 只能看到显式挂载的目录
- 非 root 执行 — `node` 用户 (uid 1000)
- 临时容器 — 每次调用新建 (`--rm`)，无状态残留

**第二层：挂载安全**

- 白名单文件在项目外 (`~/.config/nanoclaw/`)，容器无法访问
- 默认阻止 17 种敏感路径模式（`.ssh`, `.gnupg`, `.aws`, `credentials` 等）
- 符号链接解析后再验证（防遍历攻击）
- `nonMainReadOnly` — 非 Main 组的额外挂载强制只读
- **关键设计**：项目根目录只读挂载。Agent 可写路径（group 目录、IPC、`.claude/`）单独挂载。**这防止 Agent 修改宿主机代码（`src/`、`dist/`、`package.json`），否则下次重启就绕过沙箱了。**

**第三层：会话隔离**

- 每组独立的 `.claude/` 会话目录
- 组之间无法看到彼此的对话历史

**第四层：IPC 授权**

```
              发消息到自己 | 发消息到他组 | 调度任务(自) | 调度任务(他) | 查看所有任务
Main group      ✓           ✓             ✓             ✓             ✓
Non-main        ✓           ✗             ✓             ✗             仅自己的
```

**第五层：凭据隔离（Credential Proxy）**

- 真实 API 密钥永不进入容器
- 容器只看到 `ANTHROPIC_API_KEY=placeholder`
- SDK 请求发到 `http://host.docker.internal:3001`（代理）
- 代理剥离假凭据，注入真密钥，转发到 `api.anthropic.com`
- `.env` 在容器内被 `/dev/null` 遮蔽
- WhatsApp session (`store/auth/`) 也不挂载

### 3.3 Main vs Non-Main 权限对比

| 能力         | Main                        | Non-Main                   |
| ------------ | --------------------------- | -------------------------- |
| 项目根目录   | `/workspace/project` (只读) | 无                         |
| Group 文件夹 | `/workspace/group` (读写)   | `/workspace/group` (读写)  |
| 全局记忆     | 隐式（通过项目挂载）        | `/workspace/global` (只读) |
| 额外挂载     | 可配置                      | 强制只读（除非白名单允许） |
| 网络访问     | 无限制                      | 无限制                     |

---

## 4. SDK_DEEP_DIVE.md — Claude Agent SDK 逆向深度分析

**这是整个文档体系中最有价值的一篇**。它通过逆向工程 `@anthropic-ai/claude-agent-sdk` v0.2.29–0.2.34 的混淆代码，揭示了 SDK 的内部架构、Agent 循环机制、以及一个关键 bug 的发现与修复过程。

### 4.1 三层架构

```
Agent Runner (我们的代码)
  └── query() → SDK (sdk.mjs)
        └── spawns CLI subprocess (cli.js)
              └── Claude API calls, tool execution
              └── Task tool → spawns subagent subprocesses
```

**关键发现**：SDK 的 `query()` 并不直接调用 Claude API。它 spawn 了一个 `cli.js` 子进程，通过 stdin/stdout 的 JSON-lines 协议通信。所有复杂逻辑（Agent 循环、工具执行、子代理编排）都在 CLI 子进程内运行。

### 4.2 核心 Agent 循环 — EZ() 递归生成器

CLI 内部的 agentic loop 是一个**递归异步生成器 `EZ()`**（而非迭代 while 循环）：

```
EZ({ messages, systemPrompt, canUseTool, maxTurns, turnCount=1, ... })
```

**每一轮的流程**：

1. 准备消息（裁剪上下文，必要时 compact）
2. 调用 Anthropic API（`mW1` 流式函数）
3. 从响应中提取 `tool_use` 块
4. **分支**：
   - **无 tool_use** → 停止（跑停止钩子，返回结果）
   - **有 tool_use** → 执行工具 → `turnCount++` → **递归调用 EZ()**

**停止条件决策表**：

| 条件                           | 动作              | 结果类型               |
| ------------------------------ | ----------------- | ---------------------- |
| 响应有 `tool_use` 块           | 执行工具，递归 EZ | 继续                   |
| 响应无 `tool_use` 块           | 跑停止钩子，返回  | `success`              |
| `turnCount > maxTurns`         | 产出错误          | `error_max_turns`      |
| `totalCost >= maxBudgetUsd`    | 产出错误          | `error_max_budget_usd` |
| `abortController` 触发         | 产出中断消息      | 取决于上下文           |
| `stop_reason === "max_tokens"` | 重试最多 3 次     | 继续                   |

### 4.3 query() 完整选项表

文档整理了 SDK 的 ~40 个选项，包括官方文档和源码中发现的非公开选项：

| 选项              | 类型                              | 说明                                 |
| ----------------- | --------------------------------- | ------------------------------------ |
| `prompt`          | `string \| AsyncIterable`         | 输入提示（字符串 vs 流式）           |
| `resume`          | `string`                          | 恢复会话的 session ID                |
| `resumeSessionAt` | `string`                          | 从特定消息 UUID 恢复                 |
| `systemPrompt`    | `string \| preset`                | 系统提示（可追加到预设后面）         |
| `allowedTools`    | `string[]`                        | 允许的工具列表                       |
| `permissionMode`  | `PermissionMode`                  | 权限模式（含 bypassPermissions）     |
| `settingSources`  | `SettingSource[]`                 | 加载的设置来源（user/project/local） |
| `mcpServers`      | `Record<string, McpServerConfig>` | MCP 服务器配置                       |
| `hooks`           | `Record<HookEvent, ...[]>`        | 12 种生命周期钩子                    |
| `maxBudgetUsd`    | `number`                          | 预算上限                             |
| `betas`           | `SdkBeta[]`                       | Beta 功能（如 1M 上下文窗口）        |

**隐藏发现**：`settingSources` 为空数组时，SDK 不加载任何文件系统设置（默认隔离）。必须包含 `'project'` 才会加载 `CLAUDE.md`。

### 4.4 SDKMessage 类型体系

SDK 产出的 16 种消息类型（官方文档只列了 7 种）：

```
system/init              → 会话初始化，包含 session_id, tools, model
system/task_notification → 后台 Agent 完成/失败/停止
system/compact_boundary  → 对话被压缩
system/status            → 状态变更
system/hook_*            → 钩子执行状态
assistant                → Claude 的响应（文本 + tool 调用）
result                   → 最终结果（success / error_* 多种变体）
stream_event             → 流式部分消息
tool_progress            → 长时间运行的工具进度
```

**result 消息的有用字段**：`total_cost_usd`（总费用），`duration_ms`（耗时），`num_turns`（轮次），`modelUsage`（每模型用量明细）。

### 4.5 子代理执行模式 — 三种完全不同的行为

这是 SDK 中最复杂的部分：

**模式 1：同步子代理** (`run_in_background: false`)

- 父 Agent 调 Task 工具 → `VR()` 为子代理运行 `EZ()` → 父等待结果 → 返回
- 有"提升"机制：同步子代理可通过 `Promise.race()` 被提升为后台运行

**模式 2：后台任务** (`run_in_background: true`)

- Bash → spawn 命令后立即返回（空结果 + `backgroundTaskId`）
- Task/Agent → 在 `g01()` 包装器中 fire-and-forget 启动，返回 `status: "async_launched"`
- **完全不等待后台任务**就发 type: "result"
- 完成后单独发 `SDKTaskNotificationMessage`

**模式 3：Agent Teams**（最精妙的）

- Leader 正常跑 EZ 循环（包括 spawn teammates）
- Leader 的 EZ 结束后发第一个 `result`
- **然后进入后续轮询循环**：
  ```
  while (true) {
    检查是否无活跃 teammates + 无运行中任务 → break
    检查 teammates 的未读消息 → 作为新 prompt 重新进入 EZ
    如果 stdin 关闭但有活跃 teammates → 注入关闭提示
    每 500ms 轮询
  }
  ```
- 从消费者角度：收到初始 result 后，AsyncGenerator 可能**继续产出更多消息**（Leader 处理 teammate 响应后重新进入 EZ）

### 4.6 isSingleUserTurn 问题 — 核心 Bug 的发现

**这是本文档最重要的发现**。

当 `prompt` 参数是字符串时：

```javascript
QK = typeof X === 'string' // isSingleUserTurn = true
```

`isSingleUserTurn = true` 时，SDK 在收到第一个 `result` 后立即关闭 CLI 的 stdin：

```
SDK closes stdin → CLI detects stdin close
  → 轮询循环发现 stdin 关闭 + 有活跃 teammates
  → 注入关闭提示：
    "You MUST shut down your team before preparing your final response:
     1. Use requestShutdown to ask each team member to shut down gracefully
     2. Wait for shutdown approvals
     3. Use the cleanup operation to clean up the team"
  → Leader 发 shutdown_request → **Teammates 被中途杀死**
```

**实际后果**：Leader spawn 了 teammates → 他们开始研究 → Leader 完成第一轮 → result 发出 → SDK 关 stdin → teammates 可能才跑了 10 秒就被强制停止。

### 4.7 修复方案：AsyncIterable 输入流

```typescript
// 有 bug 的写法（Agent Teams 子代理会被杀死）：
query({ prompt: 'do something' })

// 修复后（保持 CLI 存活）：
query({ prompt: asyncIterableOfMessages })
```

当 `prompt` 是 `AsyncIterable` 时：

- `isSingleUserTurn = false`
- SDK **不会**在收到 result 后关闭 stdin
- CLI 保持存活，继续处理后台 Agent
- NanoClaw 控制何时结束 iterable

**额外好处**：可以在 Agent 运行期间，通过 push 新消息到 iterable 来实现**流式追加消息**（用户在 Agent 思考时发的新消息直接注入，而不用等容器退出再开新容器）。

### 4.8 V1 vs V2 API 对比

| 方面               | V1 `query()`             | V2 `createSession()`          |
| ------------------ | ------------------------ | ----------------------------- |
| `isSingleUserTurn` | 字符串 prompt 时 true    | 始终 false                    |
| 多轮对话           | 需自己管理 AsyncIterable | `send()` / `stream()` 交替    |
| stdin 生命周期     | result 后自动关闭        | `close()` 前一直开着          |
| Agent 循环         | 相同的 EZ()              | 相同的 EZ()                   |
| API 稳定性         | 稳定                     | 不稳定预览（`unstable_v2_*`） |

**关键发现**：轮次行为**零差别**。两者用同一个 CLI 子进程、同一个 EZ 递归生成器、同一套决策逻辑。唯一区别是 stdin 管理。

### 4.9 Hook 事件系统

12 种生命周期钩子事件：

```
PreToolUse / PostToolUse / PostToolUseFailure  — 工具执行前/后
Notification                                    — 通知消息
UserPromptSubmit                                — 用户提交提示
SessionStart / SessionEnd                       — 会话生命周期
Stop                                            — Agent 停止
SubagentStart / SubagentStop                    — 子代理生命周期
PreCompact                                      — 压缩前（NanoClaw 用来归档对话）
PermissionRequest                               — 权限请求
```

### 4.10 混淆标识符参考

文档还整理了 SDK 和 CLI 中关键的混淆变量名，便于后续调试：

| 混淆名 | 用途                          |
| ------ | ----------------------------- |
| `EZ`   | 核心递归 Agent 循环           |
| `XX`   | ProcessTransport（spawn CLI） |
| `$X`   | Query 类（JSON-line 路由）    |
| `QX`   | AsyncQueue（输入流缓冲）      |
| `mW1`  | Anthropic API 流式调用        |
| `PU1`  | 流式工具执行器                |
| `bd1`  | stdin 读取器                  |

---

## 5. nanoclaw-architecture-final.md — 技能引擎架构（三路合并方案）

这是 NanoClaw 的**第一代技能分发架构设计文档**（约 1100 行），后来被 `skills-as-branches.md` 取代。尽管未最终采用，但其设计思想极具参考价值。

### 5.1 核心原则：三级解决模型

所有操作都遵循这个升级链：

```
Level 1: Git（确定性、程序化）
  git merge-file 合并，git rerere 重放缓存，结构化操作直接应用。
  无 AI 参与。处理绝大多数情况。
    ↓ 无法解决时
Level 2: Claude Code（读 SKILL.md、.intent.md、state.yaml 理解上下文）
  解决 Git 无法处理的冲突。通过 git rerere 缓存解决方案，下次不再需要 AI。
    ↓ 缺乏上下文时
Level 3: 用户（当两个功能在应用层面真正冲突时，需要人类决策）
```

**关键洞察**："Clean merge ≠ working code"。即使三路合并无冲突，语义冲突（改了变量名、移了引用、改了函数签名）仍可能产生运行时错误。**所以测试必须在每次操作后运行**，无论合并是否干净。

### 5.2 共享基线（Shared Base）

`.nanoclaw/base/` 保存干净的核心代码副本（安装技能前的原始代码）。这是所有三路合并的**公共祖先**，只在核心更新时改变。

**为什么需要它**：`git merge-file` 需要三个完整文件：base（公共祖先）、current（当前文件）、other（技能修改后的文件）。base 让系统能准确计算"用户改了什么"和"技能想改什么"两个 diff 并合并。

### 5.3 两类变更

**代码文件（三路合并）**：源代码中技能需要织入逻辑的文件。用 `git merge-file` 合并。技能包中携带**完整的修改后文件**（不是 patch/diff）。

**结构化数据（确定性操作）**：`package.json`、`docker-compose.yml`、`.env.example` 等。多个技能加 npm 依赖不应该做文本合并，而是声明式聚合：

```yaml
structured:
  npm_dependencies:
    whatsapp-web.js: '^2.1.0'
  env_additions:
    - WHATSAPP_TOKEN
  docker_compose_services:
    whatsapp-redis:
      image: redis:alpine
```

所有结构化操作是**隐式**和**批量**的：收集所有技能的声明 → 写一次 `package.json` → 跑一次 `npm install`。

### 5.4 技能包结构

```
skills/add-whatsapp/
  SKILL.md                          # 技能描述和意图
  manifest.yaml                     # 元数据、依赖、结构化操作
  tests/whatsapp.test.ts            # 集成测试
  add/src/channels/whatsapp.ts      # 新增文件（直接复制）
  modify/src/server.ts              # 修改文件（完整的修改后版本）
  modify/src/server.ts.intent.md    # 意图文件
```

**为什么用完整文件而非 diff**：

- `git merge-file` 需要三个完整文件输入
- Git 三路合并用**上下文匹配**，即使用户移动了代码也能工作（line-number-based diff 会失败）
- 可审计：`diff base/file skill/modify/file` 清楚显示修改内容
- 确定性：相同输入永远相同输出

**Intent 文件**（`.intent.md`）— 结构化的意图描述：

```markdown
## What this skill adds

## Key sections

## Invariants（不变量 — 绝对不能违反的约束）

## Must-keep sections（必须保留的部分）
```

这些结构化标题给 Claude Code 在冲突解决时提供精确指导，而非让它从非结构化文本推断。

### 5.5 Apply 流程（10 步）

```
① 预检（兼容性、依赖、未跟踪变更）
② 备份（所有将被修改的文件）
③ 文件操作（rename/delete/move，在合并前执行）
④ 复制新文件
⑤ 三路合并修改文件（git merge-file）
⑥ 冲突解决（共享缓存 → rerere → Claude → 用户）
⑦ 结构化操作（npm deps + env vars + docker-compose，批量）
⑧ post_apply + 更新 state.yaml
⑨ 运行测试（必须！即使所有合并都干净）
⑩ 清理（成功删除备份，失败恢复备份）
```

### 5.6 共享解决缓存

`.nanoclaw/resolutions/` 目录随项目分发预计算的、已验证的冲突解决方案：

```yaml
# meta.yaml — 严格的哈希验证
input_hashes:
  base: 'aaa...'
  current_after_whatsapp: 'bbb...'
  telegram_modified: 'ccc...'
output_hash: 'ddd...'
```

**只有当所有输入哈希完全匹配时才应用缓存**。这确保了解决方案的正确性。

### 5.7 rerere 适配器 — 关键的工程细节

`git rerere` 不能直接与 `git merge-file` 配合工作。原因不是冲突标记格式（rerere 只哈希冲突体），而是 **rerere 需要 unmerged index entries（stage 1/2/3）来检测冲突**。`git merge-file` 只操作文件系统，不碰 index。

**适配器方案**：手动为 merge-file 的结果创建 Git index 状态：

```bash
# 1. merge-file 产出冲突后，创建三个版本的 blob
base_hash=$(git hash-object -w base.ts)
ours_hash=$(git hash-object -w current.ts)
theirs_hash=$(git hash-object -w skill.ts)

# 2. 注入 stage 1/2/3 到 index
printf '100644 %s 1\tfile.ts\0' "$base_hash" | git update-index --index-info
printf '100644 %s 2\tfile.ts\0' "$ours_hash" | git update-index --index-info
printf '100644 %s 3\tfile.ts\0' "$theirs_hash" | git update-index --index-info

# 3. 设置合并状态（rerere 检查 MERGE_HEAD）
echo "$(git rev-parse HEAD)" > .git/MERGE_HEAD

# 4. 现在 rerere 可以识别冲突并自动解决/记录
git rerere

# 5. 解决后清理
rm .git/MERGE_HEAD .git/MERGE_MSG && git reset HEAD
```

**已验证的属性**（33 个测试）：

- `merge-file` 和 `git merge` 对相同输入产生**相同的冲突体** → rerere 哈希互通
- 哈希确定性 → 相同冲突永远相同哈希
- 解决方案可移植 → 复制 `.git/rr-cache/` 到另一个 repo 即可使用
- 相邻 ~3 行内的修改会被合为一个冲突块

### 5.8 状态跟踪（state.yaml）

完整记录安装状态：已安装技能（含每文件哈希）、结构化操作结果、定制修改、路径重映射。使漂移检测即时、重放确定。

### 5.9 核心更新与迁移机制

**大多数更新**自动通过三路合并传播。

**破坏性变更**需要一个**迁移技能**——一个保留旧行为的常规技能，基于新核心编写。在更新时自动应用：

```
Core v0.7: WhatsApp 从核心移到技能
  → 发布 add-whatsapp@2.0.0 迁移技能
  → 用户更新时自动应用，保持 WhatsApp 正常运行
  → 用户可以选择 /remove-skill 来接受新默认行为
```

**更新流程全 13 步**：预检 → 预览（只用 git 命令，不改文件）→ 备份 → 文件操作 → 三路合并 → 冲突解决 → 重施定制补丁 → **更新 base 到新核心** → 应用迁移技能 → 重施更新的技能 → 重跑结构化操作 → 测试 → 清理。

### 5.10 卸载 = 不带目标技能的重放

不是反向 patch，而是：读 `state.yaml` → 去掉目标技能 → 从干净 base 重放剩余技能。这保证了状态一致性。

---

## 6. skills-as-branches.md — 技能即分支（最终方案）

这是**取代了**上述三路合并方案的**最终技能分发架构**。核心洞察：Git 本身就是最好的技能合并引擎。

### 6.1 核心思想

```
旧方案: 自定义技能引擎 + manifest + state.yaml + replay + merge-file
新方案: git branch + git merge + Claude 解决冲突
```

上游仓库维护：

- `main` — 纯核心代码（无技能代码）
- `skill/discord` — main + Discord 集成
- `skill/telegram` — main + Telegram 集成
- 以此类推

### 6.2 技能操作

**安装技能**：

```bash
git fetch upstream skill/discord
git merge upstream/skill/discord
```

**安装多个技能**：

```bash
git merge upstream/skill/discord
git merge upstream/skill/telegram
```

Git 处理组合。如果两个技能修改同一行，是真正的冲突 → Claude 解决。

**更新核心**：

```bash
git fetch upstream main
git merge upstream/main
```

**卸载技能**：

```bash
git log --merges --oneline | grep discord
git revert -m 1 <merge-commit>
```

### 6.3 为什么从三路合并方案转向分支方案

旧方案的问题：

- 需要维护 `.nanoclaw/` 状态目录（base、state.yaml、resolutions、backup）
- 需要自定义的合并引擎代码（`skills-engine/`）
- rerere 适配器增加复杂度
- 重放逻辑复杂

新方案的优势：

- **零自定义基础设施** — 全用 Git 原生操作
- **不需要 state** — Git 历史就是状态
- **不需要 replay** — 变更已在 Git 历史中
- **核心更新后不需要重施技能** — 技能变更在用户的 Git 历史中，`git merge upstream/main` 自然处理
- **冲突解决由 Claude 处理** — 一两年前不可行，但现在 AI 足够强大

### 6.4 Marketplace 机制

技能分为两类：

**运维技能**（在 `main` 上，始终可用）：

- `/setup`, `/debug`, `/update-nanoclaw`, `/customize`, `/update-skills`
- 纯指令文件（SKILL.md），无代码变更

**功能技能**（通过 marketplace，按需安装）：

- `/add-discord`, `/add-telegram` 等
- 通过 Claude Code Plugin marketplace 分发
- 安装命令：`claude plugin install nanoclaw-skills@nanoclaw-skills --scope project`
- 热加载 — 安装后立即可用，无需重启

### 6.5 CI：保持技能分支同步

GitHub Action 在每次 push 到 `main` 时：

1. 列出所有 `skill/*` 分支
2. 对每个技能分支，merge `main` 进去（merge-forward，不是 rebase）
3. 运行构建和测试
4. 通过则 push 更新的技能分支
5. 失败则开 GitHub Issue 人工处理

**为什么 merge-forward 而非 rebase**：

- 无 force-push — 保留已合并技能的用户的历史
- 用户可以 re-merge 技能分支获取更新
- Git 有正确的公共祖先

### 6.6 Flavors（风味版本）

Flavor 是**策展的 NanoClaw fork**——针对特定用途预组合技能和配置：

```yaml
# flavors.yaml（upstream 仓库中）
flavors:
  - name: NanoClaw for Sales
    repo: alice/nanoclaw
    description: Gmail + Slack + CRM, 每日管道摘要
  - name: NanoClaw Minimal
    repo: bob/nanoclaw
    description: 仅 Telegram, 无容器开销
```

`/setup` 时提供选择。安装 flavor = `git merge <flavor>/main`。

### 6.7 社区市场

任何人都可以维护自己的技能分支和 marketplace 仓库：

```json
// .claude/settings.json
{
  "extraKnownMarketplaces": {
    "nanoclaw-skills": { "source": { "source": "github", "repo": "qwibitai/nanoclaw-skills" } },
    "alice-nanoclaw-skills": { "source": { "source": "github", "repo": "alice/nanoclaw-skills" } }
  }
}
```

**特点**：无门槛（任何人可创建）、多市场共存、用同样的 merge 模式安装。

### 6.8 技能依赖

```
skill/voice-transcription → 基于 skill/whatsapp（包含 WhatsApp 的所有变更）
skill/local-whisper → 基于 skill/voice-transcription
```

依赖**隐含在 Git 历史中**：`git merge-base --is-ancestor` 判断一个分支是否是另一个的祖先。无需单独的依赖文件。

### 6.9 贡献者流程

**贡献技能**的人只需：1. Fork → 2. 从 main 分支 → 3. 改代码 → 4. 开 PR。

**维护者**：从 PR 创建 `skill/<name>` 分支 → 把 PR 精简为 CONTRIBUTORS.md 变更并合入 main → 在 marketplace 仓库添加 SKILL.md。

---

## 7. nanorepo-architecture.md — 精简版技能架构

这篇文档是 `nanoclaw-architecture-final.md` 的**浓缩版**（约 300 行 vs 1100 行），保留了所有核心概念但去除了实现细节。它的价值在于作为**快速参考**。

核心内容与 §5 相同，但更简洁：三级解决模型、共享基线、两类变更、技能包结构、Apply 流程、状态跟踪、核心更新、技能卸载、重放、19 条设计原则。

---

## 8. docker-sandboxes.md — Docker Sandbox 部署指南

这是一篇**面向特殊部署场景的实操手册**（约 300 行），指导如何在 Docker Sandbox（微型 VM）中运行 NanoClaw，实现**两层隔离**。

### 8.1 架构

```
Host (macOS / Windows WSL)
└── Docker Sandbox (micro VM，有独立内核)
    ├── NanoClaw process (Node.js)
    │   ├── Channel adapters
    │   └── Container spawner → 嵌套的 Docker daemon
    └── Docker-in-Docker
        └── nanoclaw-agent containers
            └── Claude Agent SDK
```

两层隔离：Agent 容器 + VM 边界。Sandbox 提供 MITM 代理（`host.docker.internal:3128`）处理网络和 API key 注入。

### 8.2 需要的 6 个补丁

在 Docker Sandbox 中运行 NanoClaw 需要 6 处代码修改：

| 补丁位置                 | 问题                                     | 解决方案                                    |
| ------------------------ | ---------------------------------------- | ------------------------------------------- |
| container/Dockerfile     | `npm install` 遇到 MITM 自签证书         | 添加 proxy build args + 临时禁用 strict-ssl |
| container/build.sh       | Docker build 不透传代理环境变量          | 添加 `--build-arg` 参数                     |
| src/container-runner.ts  | `/dev/null` 挂载被 Sandbox 拒绝          | 用空文件遮蔽 `.env`                         |
|                          | Agent 容器不走代理                       | 透传 HTTP_PROXY 等环境变量                  |
|                          | 容器内缺 CA 证书                         | 复制并挂载 CA cert                          |
| src/container-runtime.ts | `cleanupOrphans()` 可能杀死 Sandbox 自身 | 过滤 `os.hostname()`                        |
| src/credential-proxy.ts  | 上游 API 请求需走 Sandbox 代理           | 添加 `HttpsProxyAgent`                      |
| setup/container.ts       | build 时缺代理参数                       | 同 build.sh 的修改                          |

### 8.3 网络模型

```
Agent container → DinD bridge → Sandbox VM → host.docker.internal:3128 → Host proxy → api.anthropic.com
```

**"Bypass" 不是跳过代理**，而是代理不做 MITM 检查（WhatsApp 的 Noise 协议需要）。Node.js 不自动用 `HTTP_PROXY`，必须显式配置 `HttpsProxyAgent`。

### 8.4 渠道特定注意事项

- **Telegram**：不需要代理 bypass，但 grammy Bot 构造函数需要配置 `baseFetchConfig.agent`
- **WhatsApp**：需要 proxy bypass（`*.whatsapp.com`, `*.whatsapp.net`），WebSocket 需要 `HttpsProxyAgent`
- Clone 必须先到 `~` 再移到 workspace（virtiofs 可能在 clone 时损坏 git pack 文件）

---

## 9. APPLE-CONTAINER-NETWORKING.md — Apple Container 网络配置

这是 macOS 26 上使用 Apple Container 的**网络配置手册**。Apple Container 的 vmnet 网络默认不支持外部访问。

### 9.1 快速配置

```bash
# 1. IP 转发
sudo sysctl -w net.inet.ip.forwarding=1

# 2. NAT
echo "nat on en0 from 192.168.64.0/24 to any -> (en0)" | sudo pfctl -ef -
```

### 9.2 IPv6 DNS 问题

DNS 默认返回 IPv6（AAAA）记录在 IPv4（A）之前。由于 NAT 只处理 IPv4，Node.js 会先尝试 IPv6 然后失败。解决方案：

```bash
NODE_OPTIONS=--dns-result-order=ipv4first
```

同时设在 Dockerfile 和 container-runner.ts 的 `-e` flag 中。

### 9.3 网络拓扑

```
Container VM (192.168.64.x)
  → bridge100 (192.168.64.1，vmnet 创建的宿主机桥接)
    → IP forwarding (sysctl) 路由 bridge100 → en0
      → NAT (pfctl) 伪装 192.168.64.0/24 → en0 的 IP
        → en0 (WiFi/Ethernet) → Internet
```

---

## 10. DEBUG_CHECKLIST.md — 调试与排障手册

这篇文档是运维参考手册，记录了已知问题和诊断步骤。

### 10.1 三个已知问题（截至 2026-02-08）

**Issue 1 [已修复]：Resume 分支到陈旧位置**

- 原因：Agent Teams spawn 子代理 CLI，写入同一个 session JSONL。Resume 时 CLI 可能选了子代理活动之前的陈旧分支顶端
- 修复：传 `resumeSessionAt`，用最后一个 assistant 消息的 UUID 显式锚定

**Issue 2：IDLE_TIMEOUT == CONTAINER_TIMEOUT**

- 两个都是 30 分钟 → 同时触发 → 容器总是 SIGKILL (exit code 137) 退出
- 建议：IDLE_TIMEOUT 应缩短到 5 分钟

**Issue 3：游标在 Agent 成功前就推进了**

- `processGroupMessages` 在 Agent 运行前推进了 `lastAgentTimestamp`
- 如果容器超时 → 重试找不到消息（游标已越过）→ 消息永久丢失

### 10.2 诊断命令集

**快速状态检查**：

```bash
launchctl list | grep nanoclaw        # 服务状态
container ls | grep nanoclaw           # 运行中的容器
grep -E 'ERROR|WARN' logs/nanoclaw.log | tail -20   # 最近错误
```

**Session 分支诊断**：

```bash
# 检查 transcript JSONL 中的 parentUuid 分支
python3 -c "
import json
lines = open('.../session.jsonl').read().strip().split('\n')
for i, line in enumerate(lines):
  d = json.loads(line)
  if d.get('type') == 'user' and d.get('message'):
    parent = d.get('parentUuid', 'ROOT')[:8]
    content = str(d['message'].get('content', ''))[:60]
    print(f'L{i+1} parent={parent} {content}')
"
```

**容器超时调查**、**Agent 不响应**、**挂载问题**、**WhatsApp 认证问题**等各有对应的诊断命令集。

---

## 文档间的演进关系

```
REQUIREMENTS.md (哲学 & 愿景)
     │
     ├──→ SPEC.md (技术规格)
     │       │
     │       ├──→ SECURITY.md (安全模型细化)
     │       └──→ DEBUG_CHECKLIST.md (运维经验)
     │
     ├──→ SDK_DEEP_DIVE.md (独立的 SDK 逆向研究)
     │
     └──→ 技能分发架构（经历了三次演进）：
           │
           ├─ v1: nanoclaw-architecture-final.md (三路合并 + 自定义引擎)
           │       │
           │       └─ nanorepo-architecture.md (v1 的精简版)
           │
           └─ v2: skills-as-branches.md (Git 分支 + Claude 冲突解决)
               └─ 彻底取代 v1，删除了 skills-engine/ 全部代码

部署方案（容器运行时特定）：
  ├── docker-sandboxes.md (Docker Sandbox / DinD 部署)
  └── APPLE-CONTAINER-NETWORKING.md (Apple Container 网络)
```

---

## 核心收获

### 1. 安全设计

NanoClaw 的安全不是"加验证"而是"减攻击面"：

- **容器隔离** → Agent 只看到挂载的文件
- **凭据代理** → 密钥永不进容器
- **项目只读** → Agent 无法修改宿主机代码
- **IPC 目录即身份** → 无法伪造来源

### 2. SDK 工程

SDK 的 `isSingleUserTurn` 问题是一个精彩的逆向工程案例：

- 表面现象：Agent Teams 子代理被莫名中止
- 根因：`typeof prompt === "string"` 触发了 stdin 自动关闭
- 影响面：所有使用字符串 prompt + Agent Teams 的应用
- 修复：改用 AsyncIterable prompt

### 3. 技能架构演进

从三路合并到 Git 分支的演进体现了一个重要认知转变：

- **v1 思路**：AI 不够可靠，需要确定性的程序化合并 + 缓存
- **v2 思路**：AI（Claude）现在足够强大，可以直接解决合并冲突
- **结果**：删除了整个 skills-engine（上千行），用 `git merge` + Claude 替代

### 4. "小到能理解"的实践

NanoClaw 通过极致简约实现了可理解性：

- 单进程而非微服务
- SQLite 而非 PostgreSQL + Redis
- 文件 IPC 而非网络协议
- 改代码而非改配置
- `try ALTER catch` 而非 migration 文件
