# Pi Coding Agent 源码功能全面讲解

**核心架构：**

- `AgentSession` — 所有模式共享的核心会话抽象
- `createAgentSession()` SDK 入口
- 树形 JSONL 会话存储 + 分支/导航

**关键功能（逐一结合测试用例讲解）：**

- **Compaction**（手动/自动/Extension 钩子/Thinking 模型兼容）
- **工具系统**（7 个内置工具 + 动态注册机制）
- **并发控制**（prompt 互斥 + steer/followUp 队列）
- **自动重试**（overloaded_error 等瞬时错误）
- **Extension 系统**（发现、加载、快捷键冲突检测）
- **Skills**（SKILL.md 加载验证）
- **Prompt Templates**（$ARGUMENTS / $@ / Bash 风格切片）
- **System Prompt 构建**（工具列表 + 上下文文件 + Skills）
- **Settings 管理**（全局/项目级 + 外部编辑保护）
- **Model 解析**（模式匹配 + thinking level 语法）
- **RPC 模式**、**图像处理**、**Package 管理**等

每个功能点都附有对应测试用例的关键断言代码，方便对照源码理解。

> 基于 `packages/coding-agent` 的 README.md + 50+ 测试用例的综合分析

---

## 一、项目概览：Pi 是什么？

Pi 是一个**极简的终端 AI 编程助手**（Terminal Coding Harness），核心理念是：

> "让 Pi 适配你的工作流，而不是相反。"

它不内置子代理、计划模式等"大而全"的功能，而是提供了一套强大的**扩展机制**：Extensions、Skills、Prompt Templates、Themes，以及可分享的 Pi Packages。

**四种运行模式：**

| 模式        | 用途                        |
| ----------- | --------------------------- |
| Interactive | 交互式终端 UI（默认模式）   |
| Print/JSON  | 输出结果后退出（脚本友好）  |
| RPC         | stdin/stdout 进程间通信     |
| SDK         | 嵌入到自己的 Node.js 应用中 |

---

## 二、核心架构（src/core/）

### 2.1 AgentSession — 核心会话抽象

**文件：** `src/core/agent-session.ts`

`AgentSession` 是所有模式（interactive / print / rpc）共享的核心类，封装了：

- **Agent 状态管理**：当前模型、thinking level
- **事件总线**：subscribe 订阅事件，自动持久化到 session 文件
- **Compaction（上下文压缩）**：手动和自动两种
- **Bash 执行**
- **会话切换和分支**

```
┌─────────────────────┐
│   Interactive Mode   │   Print Mode   │   RPC Mode   │   SDK
│   (TUI / Ink)        │                │              │
└──────────┬──────────┘                                │
           │                                            │
           ▼                                            ▼
┌──────────────────────────────────────────────────────────┐
│                      AgentSession                         │
│  ┌─────────┐  ┌────────────┐  ┌──────────┐  ┌────────┐ │
│  │  Agent   │  │ SessionMgr │  │ Settings │  │ Models │ │
│  │ (pi-core)│  │ (.jsonl)   │  │ Manager  │  │Registry│ │
│  └─────────┘  └────────────┘  └──────────┘  └────────┘ │
│  ┌───────────────┐  ┌───────────┐  ┌────────┐          │
│  │  Extensions   │  │   Tools   │  │ Skills │          │
│  │   Runner      │  │  (built-in│  │        │          │
│  └───────────────┘  │  + custom)│  └────────┘          │
│                      └───────────┘                       │
└──────────────────────────────────────────────────────────┘
```

### 2.2 SDK 入口 — createAgentSession()

**文件：** `src/core/sdk.ts`

一切的起点。支持的选项：

```typescript
const { session } = await createAgentSession({
  sessionManager: SessionManager.inMemory(),  // 或 .create(dir)
  model: getModel('anthropic', 'claude-sonnet-4-5'),
  thinkingLevel: 'high',
  tools: codingTools,  // [read, bash, edit, write]
  customTools: [...],  // 扩展自定义工具
  resourceLoader: ..., // 加载 extensions/skills/prompts/themes
});
```

**测试佐证（sdk-skills.test.ts）：**

```typescript
// 默认自动发现 skills
it('should discover skills by default and expose them on session.skills')

// --no-skills 模式
it('should have empty skills when resource loader returns none')

// 自定义 skills
it('should use provided skills when resource loader supplies them')
```

---

## 三、会话管理（Session Management）

### 3.1 SessionManager — JSONL 树形结构

**文件：** `src/core/session-manager.ts`

会话以 **JSONL 格式** 存储，每一行是一个 JSON entry，通过 `id` 和 `parentId` 形成**树结构**：

```
SessionHeader   { type: "session", id, cwd, timestamp }
  ├── MessageEntry    { type: "message", id, parentId, message }
  ├── CompactionEntry { type: "compaction", summary, firstKeptEntryId }
  ├── BranchSummaryEntry { type: "branch_summary", summary }
  ├── ModelChangeEntry   { type: "model_change", provider, modelId }
  ├── ThinkingLevelChangeEntry { type: "thinking_level_change" }
  ├── CustomEntry    { type: "custom", customType, data }  // 不参与 LLM 上下文
  ├── CustomMessageEntry { type: "custom_message", content }  // 参与 LLM 上下文
  ├── LabelEntry     { type: "label", targetId, label }
  └── SessionInfoEntry { type: "session_info", name }
```

支持两种模式：

- **文件持久化**：`SessionManager.create(dir)` → 写入 `~/.pi/agent/sessions/`
- **内存模式**：`SessionManager.inMemory()` → 用于 `--no-session` 或测试

### 3.2 分支（Branching）— /tree 和 /fork

**测试佐证（agent-session-branching.test.ts）：**

```typescript
it('should allow forking from single message', async () => {
  await session.prompt('Say hello')
  const userMessages = session.getUserMessagesForForking()
  const result = await session.fork(userMessages[0].entryId)
  expect(result.selectedText).toBe('Say hello')
  expect(session.messages.length).toBe(0) // fork 后对话为空
})

it('should fork from middle of conversation', async () => {
  await session.prompt('Say one')
  await session.prompt('Say two')
  await session.prompt('Say three')
  // fork 从第二条开始，保留前一条 + 回复
  const result = await session.fork(userMessages[1].entryId)
  expect(session.messages.length).toBe(2) // user + assistant
})

it('should support in-memory forking in --no-session mode')
```

**核心设计**：fork 不创建新文件，而是在同一个 JSONL 文件中创建新的分支。通过 parentId 链实现原地分支，所有历史保留在单一文件中。

### 3.3 树导航（Tree Navigation）

**测试佐证（agent-session-tree-navigation.test.ts）：**

```typescript
// 导航到用户消息，文本放入编辑器
it('should navigate to user message and put text in editor', async () => {
  const result = await session.navigateTree(rootNode.entry.id, { summarize: false })
  expect(result.editorText).toBe('First message')
})

// 导航到助手消息时，没有编辑器文本
it('should navigate to non-user message without editor text')

// 带摘要的导航：跳转时自动生成分支摘要
it('should create branch summary when navigating with summarize=true', async () => {
  const result = await session.navigateTree(rootNode.entry.id, { summarize: true })
  expect(result.summaryEntry?.type).toBe('branch_summary')
  expect(result.summaryEntry?.summary.length).toBeGreaterThan(0)
})

// 摘要挂载到正确的父节点
it('should attach summary to correct parent when navigating to nested user message')

// 中止摘要生成
it('should handle abort during summarization')
```

---

## 四、Compaction（上下文压缩）

长对话会耗尽上下文窗口。Compaction 将旧消息压缩成摘要，保留最近的消息。

### 4.1 核心算法

**文件：** `src/core/compaction/`

```
完整对话: [U1, A1, U2, A2, U3, A3, U4, A4]
                ↓ compaction
压缩后:   [SummaryOf(U1-A2), U3, A3, U4, A4]
             ↑ compactionSummary    ↑ 保留最近的 N tokens
```

**测试佐证（compaction.test.ts）：** 纯单元测试

```typescript
// token 计算
it('should calculate total context tokens from usage') // input + output + cacheRead + cacheWrite

// 判断何时压缩
it('should return true when context exceeds threshold')

// 获取最后一个非中断助手消息的 usage
it('should skip aborted messages') // stopReason === "aborted" 跳过
```

### 4.2 手动压缩（/compact）

**测试佐证（agent-session-compaction.test.ts）：** E2E 真实 LLM 调用

```typescript
it('should trigger manual compaction via compact()', async () => {
  await session.prompt('What is 2+2?')
  await session.prompt('What is 3+3?')
  const result = await session.compact()
  expect(result.summary.length).toBeGreaterThan(0)
  // 压缩后第一条变成摘要
  expect(session.messages[0].role).toBe('compactionSummary')
})

it('should maintain valid session state after compaction') // 压缩后仍可继续使用
it('should persist compaction to session file') // JSONL 中写入 compaction entry
it('should work with --no-session mode (in-memory only)')
```

### 4.3 自动压缩

**测试佐证（agent-session-auto-compaction-queue.test.ts）：**

```typescript
// 阈值触发后恢复排队消息
it('should resume after threshold compaction when only agent-level queued messages exist')

// overflow 恢复已尝试过一次后不再重复压缩
it('should not compact repeatedly after overflow recovery already attempted')
```

自动压缩有两种触发方式：

1. **threshold**：接近上下文限制时主动压缩
2. **overflow**：上下文溢出时紧急恢复压缩后重试

### 4.4 带 Thinking 模型的压缩

**测试佐证（compaction-thinking-model.test.ts）：**

```typescript
// 修复了 maxTokens < thinkingBudget 时 compact 失败的问题
it('should compact successfully with claude-opus-4-5-thinking and thinking level high')
it('should compact successfully with claude-sonnet-4-5 (non-thinking) for comparison')
```

### 4.5 Extension 钩子自定义压缩

**测试佐证（compaction-extensions.test.ts）：**

```typescript
// 扩展可以取消压缩
it('should allow extensions to cancel compaction')

// 扩展可以提供自定义摘要（替换默认 LLM 生成）
it('should allow extensions to provide custom compaction', async () => {
  // 在 session_before_compact 事件中返回自定义 compaction
  const result = await session.compact()
  expect(result.summary).toBe('Custom summary from extension')
})
```

---

## 五、工具系统（Tools）

### 5.1 内置工具

**文件：** `src/core/tools/`

| 工具    | 功能               | 默认启用 |
| ------- | ------------------ | -------- |
| `read`  | 读取文件内容       | ✅       |
| `bash`  | 执行 bash 命令     | ✅       |
| `edit`  | 精确替换文件内容   | ✅       |
| `write` | 创建/覆盖文件      | ✅       |
| `grep`  | 搜索文件内容（rg） | ❌       |
| `find`  | 按模式查找文件     | ❌       |
| `ls`    | 列出目录           | ❌       |

**测试佐证（tools.test.ts）：**

```typescript
// read 工具
it('should read file contents that fit within limits')
it('should truncate files exceeding line limit') // >2000 行截断
it('should truncate when byte limit exceeded') // >50KB 截断
it('should handle offset parameter') // 分页读取
it('should handle limit parameter')
it('should detect image MIME type from file magic') // 通过魔数识别图片

// write 工具
it('should write file contents')
it('should create parent directories')

// edit 工具 - 精确匹配替换
it('should edit file with exact match')

// bash 工具
it('should execute commands')
```

### 5.2 动态工具注册

**测试佐证（agent-session-dynamic-tools.test.ts）：**

扩展可以在 `session_start` 事件中动态注册工具：

```typescript
it('refreshes tool registry when tools are registered after initialization', async () => {
  // 扩展在 session_start 中注册 "dynamic_tool"
  pi.registerTool({
    name: 'dynamic_tool',
    description: 'Tool registered from session_start',
    promptSnippet: 'Run dynamic test behavior',
    promptGuidelines: ['Use dynamic_tool when the user asks for...'],
    parameters: Type.Object({}),
    execute: async () => ({ content: [{ type: 'text', text: 'ok' }] })
  })

  // bindExtensions 后工具可用
  await session.bindExtensions({})
  expect(session.getAllTools().map(t => t.name)).toContain('dynamic_tool')
  expect(session.systemPrompt).toContain('- dynamic_tool: Run dynamic test behavior')
})
```

---

## 六、并发控制与重试

### 6.1 并发 Prompt 保护

**测试佐证（agent-session-concurrent.test.ts）：**

```typescript
// 流式处理中再次调用 prompt() 会抛错
it('should throw when prompt() called while streaming', async () => {
  session.prompt('First message') // 开始流式
  await expect(session.prompt('Second message')).rejects.toThrow('Agent is already processing...')
})

// steer() 和 followUp() 则可以在流式中调用
it('should allow steer() while streaming') // 转向消息：插入当前工具执行后
it('should allow followUp() while streaming') // 追加消息：等 agent 完成后发送
```

**消息队列设计：**

- `steer()` — 中断当前工具，立即插入
- `followUp()` — 等待当前轮次完成后追加
- 可配置 `"one-at-a-time"` 或 `"all"` 模式

### 6.2 自动重试

**测试佐证（agent-session-retry.test.ts）：**

```typescript
it('retries after a transient error and succeeds', async () => {
  // 第 1 次失败 (overloaded_error)，第 2 次成功
  expect(getCallCount()).toBe(2)
  expect(events).toEqual(['start:1', 'end:success=true'])
})

it('exhausts max retries and emits failure', async () => {
  // 设置 maxRetries=2，99 次失败 → 共尝试 3 次 (1 + 2 retries)
  expect(getCallCount()).toBe(3)
  expect(events).toContain('end:success=false')
})

it('prompt waits for retry completion even when assistant message_end handling is delayed')
```

配置：`settings.retry = { enabled: true, maxRetries: 3, baseDelayMs: 1 }`

---

## 七、扩展系统（Extensions）

### 7.1 发现与加载

**测试佐证（extensions-discovery.test.ts）：**

```typescript
it('discovers direct .ts files in extensions/')
it('discovers subdirectory with index.ts')
it('prefers index.ts over index.js')
```

扩展放在 `~/.pi/agent/extensions/`、`.pi/extensions/` 或 Pi Package 中。

### 7.2 快捷键冲突检测

**测试佐证（extensions-runner.test.ts）：**

```typescript
// 扩展注册 Ctrl+C 会被拦截（保留给核心功能）
it('warns when extension shortcut conflicts with built-in')
it('blocks shortcuts for reserved actions even when rebound')

// 非保留快捷键可以覆盖
it('warns but allows when extension uses non-reserved built-in shortcut')
```

### 7.3 扩展能力总览

扩展可以做到：

- 注册自定义工具（或替换内置工具）
- 注册命令（`/command`）
- 注册快捷键
- 监听事件（`tool_call`, `message_end`, `session_start` 等）
- 自定义 UI 组件（状态栏、页头、页脚、覆盖层）
- 自定义压缩逻辑
- 权限管控和路径保护

---

## 八、Skills（技能）

### 8.1 加载与验证

**文件：** `src/core/skills.ts`

**测试佐证（skills.test.ts）：**

```typescript
it('should load a valid skill') // 读取 SKILL.md + frontmatter
it("should warn when name doesn't match parent directory")
it('should warn when name contains invalid characters')
it('should warn when name exceeds 64 characters')
it('should warn and skip skill when description is missing') // 缺描述→跳过
it('should load nested skills recursively')
it('should parse disable-model-invocation frontmatter field')
```

Skill 文件示例：

```markdown
---
name: my-skill
description: 当用户问 X 时使用这个技能
---

# 步骤

1. 做 A
2. 做 B
```

调用方式：`/skill:my-skill` 或由模型自动加载。

---

## 九、Prompt Templates（提示模板）

**测试佐证（prompt-templates.test.ts）：**

```typescript
// 参数替换
it('should replace $ARGUMENTS with all args joined')
it('should replace $@ with all args joined') // $@ 等价于 $ARGUMENTS
it('should NOT recursively substitute patterns in argument values') // 防止注入

// Bash 风格切片
it('should slice from index (${@:2})') // 从第 2 个参数开始
it('should slice with length (${@:2:2})') // 从第 2 个开始取 2 个
```

使用方式：在 `~/.pi/agent/prompts/` 放置 `.md` 文件，输入 `/name` 展开模板。

---

## 十、System Prompt 构建

**测试佐证（system-prompt.test.ts）：**

```typescript
it('shows (none) for empty tools list') // 无工具时显示 "(none)"
it('includes all default tools') // read, bash, edit, write
it('includes custom tools in available tools section')
it('appends promptGuidelines to default guidelines')
it('deduplicates and trims promptGuidelines') // 去重 + 去空白
```

系统提示的组成：

1. 基础指令（或 `SYSTEM.md` 自定义替换）
2. 可用工具列表
3. 工具使用指南（根据启用的工具动态生成）
4. 上下文文件（`AGENTS.md` / `CLAUDE.md`）
5. Skills 列表
6. 当前日期时间 + 工作目录

---

## 十一、Settings 管理

**测试佐证（settings-manager.test.ts）：**

```typescript
// 保留外部编辑的设置
it('should preserve enabledModels when changing thinking level')
it('should preserve custom settings when changing theme')
it('should let in-memory changes override file changes for same key')
```

设置层级：

- `~/.pi/agent/settings.json` — 全局
- `.pi/settings.json` — 项目级（覆盖全局）

**关键设计**：修改任何一个设置时，不会覆盖文件中其他手动编辑的字段。采用"读取 → 合并 → 写入"策略。

---

## 十二、Model 解析

**测试佐证（model-resolver.test.ts）：**

```typescript
it('exact match returns model with undefined thinking level')
it('partial match returns best model') // "sonnet" → "claude-sonnet-4-5"
it('sonnet:high returns sonnet with high thinking level') // 支持 model:thinking 语法
```

支持 20+ 提供商（Anthropic、OpenAI、Google、GitHub Copilot 等），通过模式匹配快速切换。

---

## 十三、RPC 模式

**测试佐证（rpc.test.ts）：** 非 Node.js 的进程集成

```typescript
it('should get state') // 获取模型、流状态
it('should save messages to session file')
it('should handle manual compaction')
it('should execute bash command') // 外部进程执行 bash
it('should include bash output in LLM context') // bash 结果注入上下文
```

RPC 通过 stdin/stdout JSON-RPC 协议通信，适用于编辑器集成等场景。

---

## 十四、图像处理

**测试佐证：**

- **image-processing.test.ts**：图片转 PNG、调整大小

  ```typescript
  it('should convert JPEG to PNG')
  it('should resize image exceeding dimension limits')
  it('should resize image exceeding byte limit')
  ```

- **clipboard-image.test.ts**：从剪贴板粘贴图片 (Ctrl+V)
- **block-images.test.ts**：控制是否阻止图片发送

---

## 十五、Package 管理

**测试佐证（package-manager.test.ts）：**

```typescript
it('should resolve local extension paths from settings')
it('should resolve skill paths from settings')
it('should resolve project paths relative to .pi')
```

支持 npm 包和 git 仓库两种安装方式。包可以包含 extensions、skills、prompts、themes 四类资源。

---

## 十六、辅助工具

### Frontmatter 解析

**测试佐证（frontmatter.test.ts）：**

```typescript
it('parses keys, strips quotes, and returns body')
it('normalizes newlines and handles CRLF')
it('throws on invalid YAML frontmatter')
it('parses | multiline yaml syntax')
```

用于解析 Skill 和 Prompt Template 的 YAML 元数据。

---

## 十七、测试用例总表

| 测试文件                                    | 功能领域          | 类型     |
| ------------------------------------------- | ----------------- | -------- |
| agent-session-branching.test.ts             | 会话分支 (fork)   | E2E      |
| agent-session-compaction.test.ts            | 手动压缩          | E2E      |
| agent-session-auto-compaction-queue.test.ts | 自动压缩队列恢复  | Unit     |
| agent-session-concurrent.test.ts            | 并发保护          | Unit     |
| agent-session-dynamic-tools.test.ts         | 动态工具注册      | Unit     |
| agent-session-retry.test.ts                 | 自动重试          | Unit     |
| agent-session-tree-navigation.test.ts       | 树导航 + 分支摘要 | E2E      |
| compaction.test.ts                          | 压缩算法          | Unit     |
| compaction-extensions.test.ts               | 压缩扩展钩子      | E2E      |
| compaction-thinking-model.test.ts           | Thinking 模型压缩 | E2E      |
| compaction-summary-reasoning.test.ts        | 摘要推理          | Unit/E2E |
| tools.test.ts                               | 内置工具          | Unit     |
| extensions-discovery.test.ts                | 扩展发现          | Unit     |
| extensions-runner.test.ts                   | 快捷键冲突        | Unit     |
| extensions-input-event.test.ts              | 输入事件          | Unit     |
| skills.test.ts                              | 技能加载验证      | Unit     |
| sdk-skills.test.ts                          | SDK 技能选项      | Unit     |
| prompt-templates.test.ts                    | 模板参数替换      | Unit     |
| system-prompt.test.ts                       | 系统提示构建      | Unit     |
| settings-manager.test.ts                    | 设置管理          | Unit     |
| model-resolver.test.ts                      | 模型解析          | Unit     |
| model-registry.test.ts                      | 模型注册表        | Unit     |
| rpc.test.ts                                 | RPC 模式          | E2E      |
| package-manager.test.ts                     | 包管理            | Unit     |
| image-processing.test.ts                    | 图像处理          | Unit     |
| frontmatter.test.ts                         | YAML 解析         | Unit     |
| auth-storage.test.ts                        | 认证存储          | Unit     |
| session-selector-\*.test.ts                 | 会话选择器        | Unit     |
| path-utils.test.ts                          | 路径工具          | Unit     |

---

## 十八、核心设计理念

1. **极简核心 + 无限扩展**：不内置子代理、计划模式、权限弹窗等，全部通过 Extension 实现
2. **树形会话**：单文件 JSONL 存储，支持原地分支和导航
3. **渐进式上下文管理**：通过 Compaction 自动/手动压缩历史
4. **工具动态注册**：Extension 可在运行时注册/替换工具
5. **多层配置**：全局 → 项目级设置覆盖，保护外部编辑内容
6. **安全优先**：Pi Packages 有完全系统访问权限，文档明确提醒审查源码
