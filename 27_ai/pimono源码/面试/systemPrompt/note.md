# Pi-mono System Prompt 深度分析

## 一、整体架构概览

Pi-mono 的 System Prompt 系统是一个**多层组合、可扩展、可覆盖**的架构。核心思想是：

```
最终 SystemPrompt = buildSystemPrompt(资源层内容 + 工具信息 + Skills + Context Files + 自定义追加)
                     ↓ (每个 turn 前)
                   Extension beforeAgentStart hook 可动态修改
                     ↓
                   传递给 Agent → AgentLoop → LLM Context
```

涉及的核心文件：

| 文件                                         | 职责                                                    |
| -------------------------------------------- | ------------------------------------------------------- |
| `coding-agent/src/core/system-prompt.ts`     | `buildSystemPrompt()` - 核心拼装函数                    |
| `coding-agent/src/core/resource-loader.ts`   | 资源发现与加载（SYSTEM.md、AGENTS.md、Skills、Prompts） |
| `coding-agent/src/core/agent-session.ts`     | Session 层：调用 build、注入扩展 hook、传给 Agent       |
| `coding-agent/src/core/skills.ts`            | Skills 加载与 XML 格式化                                |
| `coding-agent/src/core/prompt-templates.ts`  | Prompt 模板（`/` 命令）机制                             |
| `coding-agent/src/core/extensions/runner.ts` | Extension 的 `before_agent_start` hook                  |
| `agent/src/agent.ts`                         | Agent 类，持有 `state.systemPrompt`                     |
| `agent/src/agent-loop.ts`                    | Agent 主循环，将 systemPrompt 放入 LLM Context          |

---

## 二、System Prompt 的构建流程

### 2.1 `buildSystemPrompt()` — 核心拼装函数

位于 `system-prompt.ts`，接受 `BuildSystemPromptOptions`：

```typescript
interface BuildSystemPromptOptions {
  customPrompt?: string // 自定义完整替换（来自 SYSTEM.md）
  selectedTools?: string[] // 当前激活的工具列表
  toolSnippets?: Record<string, string> // 工具的一行描述
  promptGuidelines?: string[] // 额外的 guidelines
  appendSystemPrompt?: string // 追加内容（来自 APPEND_SYSTEM.md）
  cwd?: string // 工作目录
  contextFiles?: Array<{ path: string; content: string }> // AGENTS.md 文件
  skills?: Skill[] // Skills 列表
}
```

**两条路径**：

#### 路径 A：有 `customPrompt`（用户提供了 SYSTEM.md）

```
customPrompt（原始内容）
  + appendSystemPrompt（追加内容）
  + Project Context（AGENTS.md 文件列表）
  + Skills（XML 格式的 skill 描述，仅当 read 工具可用时）
  + 日期时间 + 工作目录
```

此路径下，工具列表和 guidelines 不会自动生成，用户需在 SYSTEM.md 中自行定义。

#### 路径 B：无 `customPrompt`（使用默认 prompt）

```
默认角色定义："You are an expert coding assistant operating inside pi..."
  + Available tools: 根据 selectedTools 动态生成
  + Guidelines: 根据工具集自动推导 + 用户额外 guidelines
  + Pi documentation 路径指引
  + appendSystemPrompt
  + Project Context（AGENTS.md）
  + Skills（XML 格式）
  + 日期时间 + 工作目录
```

### 2.2 Guidelines 的智能推导

Guidelines 不是写死的，而是**根据当前激活的工具集动态生成**：

```typescript
// 有 bash 但没有 grep/find/ls → 用 bash 做文件操作
if (hasBash && !hasGrep && !hasFind && !hasLs) {
  addGuideline('Use bash for file operations like ls, rg, find')
}
// 有 bash 也有 grep/find/ls → 优先使用专用工具
else if (hasBash && (hasGrep || hasFind || hasLs)) {
  addGuideline('Prefer grep/find/ls tools over bash for file exploration')
}
// 有 read + edit → 先读后改
if (hasRead && hasEdit) {
  addGuideline('Use read to examine files before editing')
}
```

**设计意图**：避免 prompt 浪费 token 描述当前不可用的工具。不同工具组合下，LLM 获得的指引完全不同。

### 2.3 工具描述在 System Prompt 中的表示

工具描述分两层：

1. **内建工具**：硬编码在 `toolDescriptions` map 中
2. **自定义工具**（来自 Extension）：通过 `toolSnippets` 传入

```typescript
const toolDescriptions: Record<string, string> = {
  read: 'Read file contents',
  bash: 'Execute bash commands (ls, grep, find, etc.)',
  edit: 'Make surgical edits to files (find exact text and replace)',
  write: 'Create or overwrite files',
  grep: 'Search file contents for patterns',
  find: 'Find files by glob pattern',
  ls: 'List directory contents'
}
```

每个工具还可以附带 `promptGuidelines`——当工具被激活时，它的指引会自动注入到 system prompt：

```typescript
// agent-session.ts 中工具注册时收集
const toolGuidelines = this._toolPromptGuidelines.get(name)
if (toolGuidelines) {
  promptGuidelines.push(...toolGuidelines)
}
```

---

## 三、资源发现机制（Resource Loader）

### 3.1 SYSTEM.md — 完全自定义 System Prompt

`ResourceLoader` 负责发现 `SYSTEM.md` 文件：

```typescript
// 优先级：项目级 > 全局级
private discoverSystemPromptFile(): string | undefined {
  // 1. 项目级: cwd/.pi/SYSTEM.md
  const projectPath = join(this.cwd, CONFIG_DIR_NAME, "SYSTEM.md");
  if (existsSync(projectPath)) return projectPath;
  // 2. 全局级: ~/.pi/agent/SYSTEM.md
  const globalPath = join(this.agentDir, "SYSTEM.md");
  if (existsSync(globalPath)) return globalPath;
  return undefined;
}
```

还支持 `APPEND_SYSTEM.md`（同样的优先级逻辑），用于在默认 prompt 后追加内容而非完全替换。

### 3.2 AGENTS.md / CLAUDE.md — 项目上下文文件

`loadProjectContextFiles()` 遍历**从 cwd 到根目录**的所有祖先目录，寻找 `AGENTS.md` 或 `CLAUDE.md`：

```typescript
function loadProjectContextFiles(options) {
  // 1. 全局 context: ~/.pi/agent/AGENTS.md
  const globalContext = loadContextFileFromDir(resolvedAgentDir);

  // 2. 从 cwd 向上遍历每一级目录
  let currentDir = resolvedCwd;
  while (true) {
    const contextFile = loadContextFileFromDir(currentDir);
    // 收集所有找到的 AGENTS.md / CLAUDE.md
    ...
    currentDir = resolve(currentDir, "..");
  }
}
```

这些文件内容被注入到 system prompt 的 `# Project Context` 部分，格式为：

```
# Project Context

Project-specific instructions and guidelines:

## /path/to/AGENTS.md

（文件内容）
```

### 3.3 Skills 系统

Skills 使用 [Agent Skills 标准](https://agentskills.io) 的 XML 格式注入：

```xml
<available_skills>
  <skill>
    <name>browser-tools</name>
    <description>Interactive browser automation via Chrome DevTools Protocol.</description>
    <location>/path/to/SKILL.md</location>
  </skill>
</available_skills>
```

关键设计：**Skills 的完整内容不会直接放入 system prompt**，只放名称、描述和文件路径。LLM 需要使用 `read` 工具主动读取 SKILL.md 内容。这是一种**延迟加载**模式，节省 token。

条件：只有当 `read` 工具可用时，才注入 Skills section。

### 3.4 覆盖机制

ResourceLoader 提供完整的 override 链：

```typescript
// 源码级覆盖
systemPromptOverride?: (base: string | undefined) => string | undefined;
appendSystemPromptOverride?: (base: string[]) => string[];
skillsOverride?: (base) => { skills, diagnostics };
// ...同理还有 extensions, prompts, themes, agentsFiles 等
```

---

## 四、Agent Session 层的 System Prompt 管理

### 4.1 `_rebuildSystemPrompt()` — 组装最终 prompt

`AgentSession` 的 `_rebuildSystemPrompt()` 是连接 ResourceLoader 和 `buildSystemPrompt()` 的桥梁：

```typescript
private _rebuildSystemPrompt(toolNames: string[]): string {
  // 1. 从注册的工具中收集 snippets 和 guidelines
  const toolSnippets: Record<string, string> = {};
  const promptGuidelines: string[] = [];
  for (const name of validToolNames) {
    // 每个工具可以贡献自己的一行描述
    const snippet = this._toolPromptSnippets.get(name);
    // 每个工具可以贡献额外的 guidelines
    const toolGuidelines = this._toolPromptGuidelines.get(name);
  }

  // 2. 从 ResourceLoader 获取各种资源
  const loaderSystemPrompt = this._resourceLoader.getSystemPrompt();      // SYSTEM.md 内容
  const loaderAppendSystemPrompt = this._resourceLoader.getAppendSystemPrompt(); // APPEND_SYSTEM.md
  const loadedSkills = this._resourceLoader.getSkills().skills;           // Skills 列表
  const loadedContextFiles = this._resourceLoader.getAgentsFiles().agentsFiles;  // AGENTS.md

  // 3. 调用 buildSystemPrompt 组装
  return buildSystemPrompt({
    cwd, skills, contextFiles, customPrompt,
    appendSystemPrompt, selectedTools, toolSnippets, promptGuidelines
  });
}
```

### 4.2 工具变更触发 System Prompt 重建

当激活工具集发生变化时，system prompt 会自动重建：

```typescript
setActiveToolsByName(toolNames: string[]): void {
  this.agent.setTools(tools);
  // 重建 system prompt 以反映新的工具集
  this._baseSystemPrompt = this._rebuildSystemPrompt(validToolNames);
  this.agent.setSystemPrompt(this._baseSystemPrompt);
}
```

### 4.3 Extension Hook — `before_agent_start`

每次用户发送消息前，Extension 可以动态修改 system prompt：

```typescript
// agent-session.ts 中 _promptInternal()
const result = await this._extensionRunner.emitBeforeAgentStart(
  expandedText,
  currentImages,
  this._baseSystemPrompt
)

// Extension 可以返回修改后的 systemPrompt
if (result?.systemPrompt) {
  this.agent.setSystemPrompt(result.systemPrompt)
} else {
  // 未修改则恢复为 base prompt（防止上一轮的修改残留）
  this.agent.setSystemPrompt(this._baseSystemPrompt)
}
```

Extension 的 `before_agent_start` handler 能力：

```typescript
interface BeforeAgentStartEvent {
  type: "before_agent_start";
  prompt: string;            // 用户输入
  images?: ImageContent[];   // 用户图片
  systemPrompt: string;      // 当前 system prompt
}

interface BeforeAgentStartEventResult {
  message?: { ... };                // 注入自定义消息
  systemPrompt?: string;            // 修改后的 system prompt
}
```

**多个 Extension 链式处理**：每个 Extension 的输出 systemPrompt 会作为下一个 Extension 的输入。

---

## 五、System Prompt 到 LLM 的传递

### 5.1 Agent 类持有 systemPrompt

```typescript
// agent.ts
class Agent {
  private _state: AgentState = {
    systemPrompt: '',  // ← 这里存储
    model: ...,
    tools: [],
    messages: [],
    ...
  };

  setSystemPrompt(v: string) {
    this._state.systemPrompt = v;
  }
}
```

### 5.2 Agent Loop 将 systemPrompt 放入 LLM Context

```typescript
// agent-loop.ts → streamAssistantResponse()
const llmContext: Context = {
  systemPrompt: context.systemPrompt,  // ← 直接传入
  messages: llmMessages,               // 经过 convertToLlm 转换的消息
  tools: context.tools
};

const response = await streamFunction(config.model, llmContext, { ... });
```

`Context` 对象直接传给 `streamSimple` 函数，最终由具体的 provider（Anthropic/OpenAI/Google 等）将 systemPrompt 转换为对应 API 的 system 消息格式。

---

## 六、Prompt Template 系统

### 6.1 模板加载

从三个来源加载 `.md` 文件作为 prompt 模板：

```
1. 全局: ~/.pi/agent/prompts/*.md
2. 项目: cwd/.pi/prompts/*.md
3. CLI 显式指定路径
```

### 6.2 模板参数替换

支持 bash 风格的参数占位符：

```markdown
---
description: Code review helper
---

Review the following code: $1

Focus on: $ARGUMENTS
```

替换规则：

| 占位符              | 含义                        |
| ------------------- | --------------------------- |
| `$1`, `$2`, ...     | 位置参数                    |
| `$@` / `$ARGUMENTS` | 所有参数拼接                |
| `${@:N}`            | 从第 N 个参数开始的所有参数 |
| `${@:N:L}`          | 从第 N 个参数开始取 L 个    |

### 6.3 模板展开

用户输入 `/templateName arg1 arg2` → 经过 `expandPromptTemplate()` 展开为模板内容 + 参数替换。

---

## 七、Compaction（上下文压缩）的 System Prompt

当对话过长时，会使用专门的 summarization system prompt 进行上下文压缩：

```typescript
const SUMMARIZATION_SYSTEM_PROMPT = `You are a context summarization assistant. 
Your task is to read a conversation between a user and an AI coding assistant, 
then produce a structured summary following the exact format specified.

Do NOT continue the conversation. Do NOT respond to any questions in the conversation. 
ONLY output the structured summary.`
```

这是一个独立的、固定的 system prompt，与主 system prompt 完全分离。

---

## 八、设计亮点总结

### 8.1 层次化覆盖机制

```
CLI --system-prompt 参数
  ↓ 覆盖
SYSTEM.md 文件发现（项目级 > 全局级）
  ↓ 覆盖
默认 system prompt
  ↓ 追加
APPEND_SYSTEM.md
  ↓ 追加
AGENTS.md / CLAUDE.md（祖先目录遍历）
  ↓ 追加
Skills（XML 格式描述）
  ↓ 动态修改
Extension before_agent_start hook
```

### 8.2 工具感知的 Prompt 生成

System Prompt 不是静态的，而是根据当前工具集动态调整：

- 工具列表随激活状态变化
- Guidelines 根据工具组合智能推导
- 工具自身可以贡献 snippets 和 guidelines

### 8.3 Skills 的延迟加载

不把 SKILL.md 全文塞入 system prompt，只放元数据（名称+描述+路径），让 LLM 按需读取。这在 skill 数量多的场景下大幅节省 token。

### 8.4 Extension 的 Per-turn 修改能力

Extension 可以在每个 turn 开始前修改 system prompt，但修改是临时的——下一个 turn 会恢复到 base prompt（除非 Extension 再次修改）。这避免了 Extension 的修改互相干扰或累积。

### 8.5 Context Files 的祖先目录遍历

AGENTS.md 不仅在 cwd 查找，而是从 cwd 一直向上查找到 `/`。这允许 monorepo 场景下在不同层级放置不同的上下文文件，子目录继承父目录的上下文。

---

## 九、面试关键问答

**Q: Pi 的 system prompt 是如何组装的？**

A: 分两条路径。如果用户提供了 SYSTEM.md（通过文件或 CLI），则用它替换默认 prompt；否则使用内建的默认 prompt。无论哪条路径，都会追加 APPEND_SYSTEM.md、AGENTS.md 上下文文件、Skills XML 描述、以及日期时间和工作目录信息。

**Q: 工具和 system prompt 的关系是什么？**

A: 当前激活的工具集直接影响 system prompt 的内容。工具列表会列在 prompt 中，guidelines 会根据工具组合动态推导（如有 edit 就提示"先 read 再 edit"），每个工具还可以贡献自己的 prompt snippet 和 guidelines。工具变更时自动触发 prompt 重建。

**Q: Extension 如何修改 system prompt？**

A: 通过 `before_agent_start` hook。每次用户发送消息前触发，Extension 收到当前 system prompt，可返回修改后的版本。多个 Extension 链式执行。修改仅对当前 turn 有效，下一 turn 会恢复 base prompt。

**Q: Skills 是怎么集成到 system prompt 中的？**

A: 以 XML 格式注入 skill 的元数据（名称、描述、文件路径），不包含完整内容。LLM 在需要时使用 `read` 工具读取 SKILL.md。这是延迟加载模式，节省 token。仅当 read 工具可用时才注入。

**Q: AGENTS.md 和 SYSTEM.md 有什么区别？**

A: SYSTEM.md 是**替换**整个 system prompt 的机制；AGENTS.md 是**追加**项目上下文的机制。AGENTS.md 会遍历从 cwd 到根目录的所有祖先目录，支持 monorepo 多层级上下文继承。SYSTEM.md 只在 `.pi/` 和全局目录查找。

**Q: Prompt Template 和 System Prompt 的关系是什么？**

A: 完全不同的机制。Prompt Template 是用户输入侧的模板（`/template args`），展开后作为用户消息发送，不影响 system prompt。System Prompt 是固定的系统指令，每次 LLM 调用都会附加。
