# pi-mono SubAgent 深度分析

## 1. 整体架构概览

pi-mono 的 subagent 系统采用 **进程隔离 + 扩展注册** 的架构。核心思路：每个子 agent 是一个独立的 `pi` 子进程，拥有隔离的上下文窗口（context window），通过 JSON 模式流式传输结果给父 agent。

```
┌────────────────────────────────────────────────┐
│ 主 Agent (coding-agent interactive mode)        │
│                                                 │
│  Extension System (runner.ts)                   │
│    └── registerTool("subagent", ...)            │
│         │                                       │
│         ├── Single: spawn 1 子进程               │
│         ├── Parallel: spawn N 子进程(≤8,并发≤4)  │
│         └── Chain: 顺序 spawn, {previous} 传递   │
│                                                 │
│  每个子进程: pi --mode json -p --no-session      │
│    └── 独立 context window                       │
│    └── 独立 system prompt (来自 agent .md)       │
│    └── 独立 tools/model 配置                     │
└────────────────────────────────────────────────┘
```

## 2. 三层架构

### 2.1 底层: Agent Loop (`packages/agent/src/agent-loop.ts`)

agent-loop 是整个系统的核心循环引擎，**不区分"主 agent"还是"子 agent"**——所有 agent 共用同一套循环逻辑。

```typescript
// 双层循环结构
async function runLoop(...) {
  // 外层循环: 处理 follow-up messages
  while (true) {
    // 内层循环: 处理 tool calls + steering messages
    while (hasMoreToolCalls || pendingMessages.length > 0) {
      // 1. 注入 pending messages（steering 中断消息）
      // 2. streamAssistantResponse → LLM 流式响应
      // 3. 如果有 tool calls → executeToolCalls（顺序执行）
      // 4. 每次 tool 执行后，检查 getSteeringMessages（用户中断）
    }
    // 内层结束后，检查 getFollowUpMessages（追加对话）
    // 没有更多消息 → break 退出
  }
}
```

**关键设计点：**

- **Steering Messages**: 用户在 agent 工作时输入的中断消息。在**每次 tool 执行后**检查，如果有 steering 消息，**跳过剩余 tool calls**，立即注入到下一轮
- **Follow-up Messages**: 在 agent 完成所有工作后才检查的追加消息。允许在 agent "想停下来"时追加任务
- **AgentMessage 抽象**: 使用 `AgentMessage`（LLM Message + 自定义消息的联合类型），通过 `convertToLlm` 在 LLM 调用边界才转换为标准 Message[]

### 2.2 中层: Agent Class (`packages/agent/src/agent.ts`)

Agent 类是对 agent-loop 的状态封装，提供响应式 API：

```typescript
class Agent {
  // 状态管理
  private _state: AgentState  // systemPrompt, model, tools, messages, isStreaming...
  
  // 消息队列
  private steeringQueue: AgentMessage[]  // 中断队列
  private followUpQueue: AgentMessage[]  // 追加队列
  
  // 核心方法
  prompt(input)   // 发送新 prompt，启动 agent loop
  continue()      // 从当前状态继续（重试/恢复）
  steer(msg)      // 入队 steering 消息
  followUp(msg)   // 入队 follow-up 消息
  abort()         // 通过 AbortController 中止
  waitForIdle()   // 等待 agent 完成
  
  // 事件系统
  subscribe(fn)   // 订阅 AgentEvent（agent_start/end, turn_start/end, message_*, tool_execution_*）
}
```

**事件流：**
```
agent_start
  → turn_start
    → message_start (user)
    → message_end (user)
    → message_start (assistant, 流式)
    → message_update (assistant, 多次)  
    → message_end (assistant)
    → tool_execution_start
    → tool_execution_update (可选, 部分结果)
    → tool_execution_end
    → message_start (toolResult)
    → message_end (toolResult)
  → turn_end
  → turn_start ... (如果有更多 tool calls)
agent_end
```

### 2.3 上层: SubAgent Extension (`examples/extensions/subagent/index.ts`)

SubAgent 作为一个 **Extension Tool** 注册到主 agent。它不是 agent-loop 内部的概念，而是通过扩展系统（Extension API）在上层实现的。

## 3. SubAgent Tool 详细设计

### 3.1 三种执行模式

| 模式 | 参数 | 行为 |
|------|------|------|
| **Single** | `{ agent, task }` | 单个子进程执行任务 |
| **Parallel** | `{ tasks: [{agent, task}, ...] }` | 最多 8 个任务，并发度≤4 |
| **Chain** | `{ chain: [{agent, task}, ...] }` | 顺序执行，`{previous}` 占位符传递上游输出 |

### 3.2 进程隔离机制

每个子 agent 通过 `spawn("pi", args)` 创建独立进程：

```typescript
async function runSingleAgent(...) {
  const args = ["--mode", "json", "-p", "--no-session"];
  // --mode json: 输出 JSON 事件流（非交互 TUI）
  // -p: prompt mode（非交互式）
  // --no-session: 不持久化 session
  
  if (agent.model) args.push("--model", agent.model);
  if (agent.tools) args.push("--tools", agent.tools.join(","));
  
  // 如果有 system prompt, 写入临时文件传递（安全: mode 0o600）
  if (agent.systemPrompt.trim()) {
    const tmp = writePromptToTempFile(agent.name, agent.systemPrompt);
    args.push("--append-system-prompt", tmp.filePath);
  }
  
  args.push(`Task: ${task}`);
  
  const proc = spawn("pi", args, { 
    cwd: cwd ?? defaultCwd, 
    shell: false,           // 不经过 shell，安全
    stdio: ["ignore", "pipe", "pipe"] 
  });
  
  // 解析 stdout 的 NDJSON 事件流
  proc.stdout.on("data", (data) => {
    // 逐行解析 JSON 事件
    // event.type === "message_end" → 收集消息、统计 usage
    // event.type === "tool_result_end" → 收集工具结果
  });
}
```

**为什么用进程隔离而非内存隔离？**

1. **Context Window 完全隔离**：每个子 agent 有独立的上下文窗口，不会污染主 agent 的对话
2. **资源限制**：并发子进程数受 MAX_CONCURRENCY=4 控制
3. **安全边界**：子进程的 tool 权限独立配置
4. **故障隔离**：子进程崩溃不影响主进程
5. **模型隔离**：每个子 agent 可以用不同模型（scout 用 Haiku 省钱，worker 用 Sonnet 做复杂任务）

### 3.3 Agent 发现机制 (`agents.ts`)

```typescript
// 两级 agent 配置目录
// 用户级: ~/.pi/agent/agents/*.md  （始终加载）
// 项目级: .pi/agents/*.md          （需显式启用 agentScope: "both"）

interface AgentConfig {
  name: string;          // frontmatter 中的 name
  description: string;   // frontmatter 中的 description
  tools?: string[];      // 限制可用工具列表
  model?: string;        // 指定模型
  systemPrompt: string;  // markdown body 作为 system prompt
  source: "user" | "project";
  filePath: string;
}

// 发现逻辑：
// 1. 扫描目录中的 .md 文件
// 2. 解析 YAML frontmatter + markdown body
// 3. scope="both" 时项目级覆盖同名用户级 agent
```

### 3.4 安全模型

1. **默认只加载用户级 agents**：`agentScope` 默认 `"user"`
2. **项目级 agents 需确认**：`confirmProjectAgents: true`（默认），交互式弹窗确认
3. **Prompt 隔离**：system prompt 写入临时文件（mode 0o600），用完立即删除
4. **无 shell 注入**：`spawn("pi", args, { shell: false })`

## 4. Chain 模式详解——最有趣的设计

Chain 模式实现了 **agent 间的信息流水线**：

```typescript
// Chain 执行核心逻辑
for (let i = 0; i < params.chain.length; i++) {
  const step = params.chain[i];
  
  // 关键: {previous} 占位符被替换为上一步的输出
  const taskWithContext = step.task.replace(/\{previous\}/g, previousOutput);
  
  const result = await runSingleAgent(..., taskWithContext, ...);
  
  // 错误时中断 chain
  if (result.exitCode !== 0 || result.stopReason === "error") {
    return { isError: true, ... };
  }
  
  // 提取最终输出作为下一步输入
  previousOutput = getFinalOutput(result.messages);
}
```

典型 workflow `/implement <query>`：
```
scout (Haiku, 快速) → 输出 context 摘要
  ↓ {previous}
planner (Sonnet, 规划) → 输出实施计划
  ↓ {previous}
worker (Sonnet, 执行) → 实际修改文件
```

## 5. 流式更新与 UI 渲染

SubAgent 支持**流式中间状态更新**：

```typescript
// onUpdate 回调
const emitUpdate = () => {
  onUpdate({
    content: [{ type: "text", text: getFinalOutput(messages) || "(running...)" }],
    details: makeDetails([currentResult]),  // 包含 usage 统计、消息历史
  });
};

// 每收到一个 message_end 事件就 emitUpdate()
// Parallel 模式：每个 task 的进度独立更新
```

渲染系统：
- **Collapsed view**: 状态图标 + agent 名 + 最近 10 条操作 + usage 统计
- **Expanded view** (Ctrl+O): 完整 task 文本 + 所有 tool calls + Markdown 渲染输出
- **Parallel**: 显示 "2/3 done, 1 running" 实时状态

## 6. 与其他系统的对比

| 特性 | pi-mono SubAgent | OpenAI Swarm | LangGraph |
|------|-----------------|-------------|-----------|
| 隔离方式 | 进程隔离 | 内存隔离(函数调用) | 内存隔离(图节点) |
| Agent 定义 | Markdown + YAML frontmatter | Python 代码 | Python 代码 |
| 信息传递 | stdout JSON 流 + {previous} | 函数返回值 | State graph |
| 并行支持 | ✅ (max 8, concurrency 4) | ❌ | ✅ |
| 流式更新 | ✅ NDJSON 事件流 | 部分 | 部分 |
| 安全模型 | 用户/项目分级+确认 | 无 | 无 |
| 模型混用 | ✅ 每个 agent 独立模型 | ✅ | ✅ |

## 7. 面试关键讨论点

### Q: 为什么选择进程隔离而不是线程/协程？
**A**: 核心目标是 **context window 隔离**。LLM 对话是有状态的，子 agent 需要完全独立的对话上下文。进程隔离是最自然的方式——直接复用 `pi` CLI 的完整能力（工具链、模型切换、session 管理），无需在内存中维护多份状态。代价是 fork 开销更大，但 LLM API 延迟远大于进程创建开销。

### Q: Chain 模式的 {previous} 有什么局限？
**A**: 
1. 上下文膨胀：前一步输出全量传递给下一步，可能超出 context window
2. 无法回溯：chain 是单向的，后续 agent 无法请求前置 agent 补充信息
3. 错误传播：一步失败整个 chain 中断，没有重试单个 step 的机制

### Q: agent-loop 的 steering message 机制有什么用？
**A**: 实现了 **Human-in-the-loop** 模式。用户在 agent 执行复杂多工具操作时可以随时干预——steering 在每个 tool 执行完成后检查，如果用户输入了新指令，立即**跳过剩余 tool calls**并注入新指令。这实现了"中途纠正"而不是"中止重来"。

### Q: 扩展系统如何实现 subagent？  
**A**: SubAgent 不是 agent-loop 内建功能，而是通过 Extension API 实现的 **AgentTool**。Extension 通过 `pi.registerTool()` 注册一个名为 `subagent` 的工具。当主 agent（LLM）在对话中调用 `subagent` 工具时，Extension 的 `execute()` 方法被触发，它 spawn 子进程、收集结果、返回 `AgentToolResult`。这种设计将 subagent 能力解耦为可选插件。

### Q: 如何保证安全性？
**A**:
1. **spawn 不经过 shell**：`{ shell: false }` 防止命令注入
2. **用户/项目 agent 分级**：项目级 agent 需显式启用 + 弹窗确认
3. **临时文件安全**：system prompt 临时文件用 mode `0o600` 写入，用完 unlink
4. **Tool 白名单**：每个 agent 可独立限制可用工具（如 reviewer 只有只读工具）
5. **并发限制**：MAX_PARALLEL_TASKS=8, MAX_CONCURRENCY=4 防止资源耗尽
