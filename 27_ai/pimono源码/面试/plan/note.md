# Pi-Mono Plan Mode 深度源码分析

## 1. 概述

Plan Mode 是 pi-mono (coding-agent) 的一个 **Extension（扩展）** 实现，提供**只读探索 → 生成计划 → 执行计划**的完整工作流。它是一个典型的 Agent 模式切换设计，核心思想是：

> **先规划、再执行；规划阶段禁止一切写操作，执行阶段跟踪进度。**

源码位置：`packages/coding-agent/examples/extensions/plan-mode/`

---

## 2. 架构总览

```
┌─────────────────────────────────────────────────┐
│                 Extension API                    │
│  (pi-mono 提供给扩展的标准接口)                    │
│  ┌─────────┬────────────┬──────────────────────┐ │
│  │ Events  │ Commands   │ Tools / UI / State   │ │
│  └─────────┴────────────┴──────────────────────┘ │
└──────────────────────┬──────────────────────────┘
                       │
            ┌──────────▼──────────┐
            │   Plan Mode Ext     │
            │                     │
            │  三种状态机:         │
            │  Normal → Plan → Execution
            └─────────────────────┘
```

### 三个模式状态

| 状态 | `planModeEnabled` | `executionMode` | 可用工具 |
|------|:-:|:-:|------|
| **Normal** | `false` | `false` | `read, bash, edit, write` |
| **Plan** (只读) | `true` | `false` | `read, bash, grep, find, ls, questionnaire` |
| **Execution** (追踪) | `false` | `true` | `read, bash, edit, write` + 进度追踪 |

---

## 3. 核心模块分析

### 3.1 入口注册 (`index.ts`)

```typescript
export default function planModeExtension(pi: ExtensionAPI): void
```

扩展工厂函数，接收 `ExtensionAPI` 对象。注册了：

| 注册项 | 方法 | 说明 |
|--------|------|------|
| CLI Flag | `pi.registerFlag("plan", ...)` | 启动时 `--plan` 进入计划模式 |
| 命令 `/plan` | `pi.registerCommand("plan", ...)` | 切换计划模式 |
| 命令 `/todos` | `pi.registerCommand("todos", ...)` | 显示当前计划进度 |
| 快捷键 | `pi.registerShortcut(Key.ctrlAlt("p"), ...)` | Ctrl+Alt+P 切换 |
| 6个事件钩子 | `pi.on(...)` | 核心逻辑所在 |

### 3.2 六大事件钩子（核心逻辑）

这是 Plan Mode 最精妙的部分——**完全基于事件驱动**，没有直接修改 Agent 循环代码：

#### ① `tool_call` — 拦截危险命令

```typescript
pi.on("tool_call", async (event) => {
    if (!planModeEnabled || event.toolName !== "bash") return;
    const command = event.input.command as string;
    if (!isSafeCommand(command)) {
        return { block: true, reason: `Plan mode: command blocked...` };
    }
});
```

**关键设计**：不是简单地禁用 bash 工具，而是**细粒度地过滤每条命令**，允许 `git status`、`ls -la` 等只读操作，阻止 `rm`、`git push` 等写操作。

#### ② `context` — 过滤上下文

```typescript
pi.on("context", async (event) => {
    if (planModeEnabled) return; // 计划模式下保留所有上下文
    // 退出计划模式后，过滤掉含 [PLAN MODE ACTIVE] 的历史消息
    return { messages: event.messages.filter(...) };
});
```

**为什么需要过滤？** 避免计划模式的 system prompt 污染正常模式的上下文窗口。

#### ③ `before_agent_start` — 注入系统提示

```typescript
pi.on("before_agent_start", async () => {
    if (planModeEnabled) {
        return {
            message: {
                customType: "plan-mode-context",
                content: `[PLAN MODE ACTIVE]
You are in plan mode - a read-only exploration mode...
Create a detailed numbered plan under a "Plan:" header:
Plan:
1. First step description
...`,
                display: false, // 不在UI显示，仅发给LLM
            },
        };
    }
    if (executionMode && todoItems.length > 0) {
        // 注入剩余待办列表，指导LLM按计划执行
        return {
            message: {
                content: `[EXECUTING PLAN]
Remaining steps: ...
After completing a step, include a [DONE:n] tag`,
            },
        };
    }
});
```

**核心思想**：通过注入隐藏消息来**引导 LLM 行为**，而不是修改 agent 核心循环。

#### ④ `turn_end` — 每轮追踪完成状态

```typescript
pi.on("turn_end", async (event, ctx) => {
    if (!executionMode || todoItems.length === 0) return;
    const text = getTextContent(event.message);
    if (markCompletedSteps(text, todoItems) > 0) {
        updateStatus(ctx); // 更新进度条UI
    }
    persistState(); // 持久化状态
});
```

通过正则 `[DONE:n]` 匹配来追踪进度——LLM 在回复中写 `[DONE:1]`，Extension 自动勾选。

#### ⑤ `agent_end` — 计划完成后的交互决策

```typescript
pi.on("agent_end", async (event, ctx) => {
    // 执行模式下：检查是否全部完成
    if (executionMode && todoItems.every(t => t.completed)) {
        // 发送完成消息，清理状态
        executionMode = false;
        todoItems = [];
        return;
    }
    // 计划模式下：提取计划步骤、弹出选择菜单
    const extracted = extractTodoItems(getTextContent(lastAssistant));
    const choice = await ctx.ui.select("Plan mode - what next?", [
        "Execute the plan (track progress)",
        "Stay in plan mode",
        "Refine the plan",
    ]);
    // 根据用户选择切换状态
});
```

**关键交互流程**：Plan 结束后用 `ctx.ui.select()` 给用户三个选择，形成闭环。

#### ⑥ `session_start` — 状态恢复

```typescript
pi.on("session_start", async (_event, ctx) => {
    // 恢复 CLI flag
    if (pi.getFlag("plan") === true) planModeEnabled = true;
    // 从 session entries 恢复持久化状态
    const planModeEntry = entries.filter(e => e.customType === "plan-mode").pop();
    // 恢复后重新扫描消息，重建完成状态
});
```

支持 **session 恢复**：关闭再打开后，从持久化的 entries 中恢复 todo 状态，并重新扫描 `[DONE:n]` 标记。

---

## 4. 工具函数分析 (`utils.ts`)

### 4.1 命令安全过滤 — 双重检查

```typescript
export function isSafeCommand(command: string): boolean {
    const isDestructive = DESTRUCTIVE_PATTERNS.some(p => p.test(command));
    const isSafe = SAFE_PATTERNS.some(p => p.test(command));
    return !isDestructive && isSafe; // 必须同时：不在黑名单 && 在白名单
}
```

**安全设计要点**：
- **白名单制**：未知命令默认拒绝（`my-script.sh` → 拒绝）
- **黑名单兜底**：即使在白名单中，含有 `>` 重定向、`sudo` 等也会被阻止
- 覆盖了文件操作、git 写操作、包管理器安装、系统命令、编辑器等

### 4.2 计划提取

```typescript
export function extractTodoItems(message: string): TodoItem[]
```

- 匹配 `Plan:` 或 `**Plan:**` 标题
- 提取 `1.` 或 `1)` 格式的编号步骤
- 过滤太短（≤5字符）、以代码/命令开头的条目
- `cleanStepText` 清理 markdown 格式、去除动作词前缀、截断到 50 字符

### 4.3 进度标记

```typescript
export function markCompletedSteps(text: string, items: TodoItem[]): number
```

从 LLM 回复中提取 `[DONE:1]`、`[DONE:2]` 等标记，标记对应步骤为完成。

---

## 5. ExtensionAPI 架构（pi-mono 的扩展系统）

Plan Mode 依赖的 `ExtensionAPI` 是 pi-mono 的核心抽象：

```typescript
interface ExtensionAPI {
    // 事件订阅（14+种生命周期事件）
    on(event: "tool_call" | "context" | "before_agent_start" | "turn_end" | "agent_end" | "session_start" | ..., handler): void;

    // 工具注册
    registerTool(tool: ToolDefinition): void;
    setActiveTools(toolNames: string[]): void;    // ← Plan Mode 用这个切换工具集

    // 命令/快捷键/Flag
    registerCommand(name, options): void;
    registerShortcut(shortcut, options): void;
    registerFlag(name, options): void;

    // 消息系统
    sendMessage(message, options): void;           // 发送自定义消息
    sendUserMessage(content, options): void;       // 模拟用户消息
    appendEntry(customType, data): void;           // 持久化自定义状态

    // UI 控制
    // (通过 ExtensionContext.ui 访问)
}
```

**设计亮点**：扩展不需要 fork 或修改 agent 核心代码，仅通过事件 hook + API 调用即可实现完整的模式切换功能。

---

## 6. 面试关键问题总结

### Q1: Plan Mode 的核心架构思想是什么？

**A**: 基于 Extension System 的**事件驱动状态机**。三个状态（Normal/Plan/Execution）通过事件钩子切换，不侵入 agent 核心循环。关键手段：
- `setActiveTools()` 控制可用工具集
- `before_agent_start` 注入隐藏系统提示引导 LLM
- `tool_call` 拦截危险操作
- `turn_end` / `agent_end` 追踪和管理进度

### Q2: 如何保证 Plan 模式的安全性？

**A**: **白名单 + 黑名单双重校验**。
- 工具级：`setActiveTools` 直接移除 `edit`、`write` 工具
- 命令级：bash 命令通过 `isSafeCommand` 逐条过滤，即使 bash 工具可用，也只允许安全命令
- 未知命令默认拒绝

### Q3: LLM 如何知道自己在 Plan 模式？

**A**: 通过 `before_agent_start` 钩子注入一条 `display: false` 的隐藏消息，包含 `[PLAN MODE ACTIVE]` 标记和行为指导（输出格式要求、限制说明）。LLM 看得到，用户看不到。

### Q4: 进度追踪如何实现？

**A**: 
1. Plan 阶段：从 LLM 回复中正则提取 `Plan:` 标题下的编号步骤 → `TodoItem[]`
2. Execution 阶段：每个 `turn_end` 扫描回复中的 `[DONE:n]` 标记 → 更新完成状态
3. UI 展示：`ctx.ui.setWidget()` 显示进度 widget，`ctx.ui.setStatus()` 更新状态栏

### Q5: Session 恢复如何处理？

**A**: 
- `pi.appendEntry("plan-mode", state)` 在每次状态变更后持久化
- `session_start` 事件中恢复状态
- 恢复后重新扫描 `plan-mode-execute` 标记之后的消息，重建 `[DONE:n]` 完成状态
- 避免旧 plan 的 `[DONE:n]` 干扰当前 plan

### Q6: 与其他 Agent 框架（LangChain、AutoGPT）的 Plan-and-Execute 有什么区别？

**A**: 
- **LangChain ReAct/Plan-and-Execute**: 在 agent loop 层面硬编码规划逻辑，需要修改 agent 代码
- **Pi-Mono Plan Mode**: 纯 Extension 实现，零侵入 agent 核心。通过事件 hook + prompt injection + tool filtering 实现完整的规划执行流程
- 启示：**好的扩展系统 ≈ 不需要改框架就能实现任何模式**

### Q7: 这个设计有什么局限性？

**A**:
- 依赖 LLM 输出特定格式（`Plan:` 标题、`[DONE:n]` 标记），格式不匹配则功能失效
- 命令白名单需要人工维护，可能遗漏安全命令或新增危险命令
- 状态管理用闭包变量（`let planModeEnabled`），单例模式，不支持并发 session
- 计划步骤提取的正则比较脆弱，复杂的 markdown 格式可能提取失败

---

## 7. 代码文件索引

| 文件 | 行数 | 说明 |
|------|------|------|
| `examples/extensions/plan-mode/index.ts` | 341 | 主逻辑：状态机 + 6个事件钩子 |
| `examples/extensions/plan-mode/utils.ts` | 168 | 工具函数：命令过滤、计划提取、进度标记 |
| `examples/extensions/plan-mode/README.md` | 66 | 文档 |
| `test/plan-mode-utils.test.ts` | 260 | 工具函数单元测试 |
| `src/core/extensions/types.ts` | 1200+ | ExtensionAPI 接口定义 |
| `src/core/extensions/loader.ts` | 270+ | ExtensionAPI 实现（代理到 runtime） |
