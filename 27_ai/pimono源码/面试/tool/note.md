# Pi-mono Tool 系统深度分析

## 一、整体架构：三层抽象

```
┌──────────────────────────────────────────────────────┐
│  ai 包 (最底层)                                       │
│  Tool { name, description, parameters }              │
│  → 纯 LLM schema 描述，无执行能力                      │
└──────────┬───────────────────────────────────────────┘
           │ extends
┌──────────▼───────────────────────────────────────────┐
│  agent 包 (中间层)                                    │
│  AgentTool extends Tool {                            │
│    + label: string        (UI 展示名)                 │
│    + execute(id, params, signal?, onUpdate?)          │
│      → Promise<AgentToolResult<TDetails>>             │
│  }                                                    │
└──────────┬───────────────────────────────────────────┘
           │ 实例化 / 扩展
┌──────────▼───────────────────────────────────────────┐
│  coding-agent 包 (最上层)                             │
│  7 个内置 tool + Extension 系统动态注册 tool            │
│  bash / read / edit / write / grep / find / ls       │
└──────────────────────────────────────────────────────┘
```

### 1.1 类型层次

```typescript
// ai 包: 纯描述型，给 LLM 看
interface Tool<TParameters extends TSchema> {
  name: string
  description: string
  parameters: TParameters // TypeBox schema
}

// agent 包: 加了执行能力
interface AgentTool<TParams, TDetails> extends Tool<TParams> {
  label: string
  execute(
    toolCallId: string,
    params: Static<TParams>,
    signal?: AbortSignal,
    onUpdate?: AgentToolUpdateCallback<TDetails>
  ): Promise<AgentToolResult<TDetails>>
}

// AgentToolResult: 统一的返回格式
interface AgentToolResult<T> {
  content: (TextContent | ImageContent)[] // 给 LLM 看的内容
  details: T // 给 UI 看的结构化数据
}
```

**关键设计**: `content` 和 `details` 的分离非常精妙——content 会被序列化发给 LLM 作为 toolResult，而 details 只用于 TUI 展示(如 diff、截断信息等)。LLM 只看到精简文本，UI 可以做丰富渲染。

### 1.2 参数校验：TypeBox

所有 tool 用 `@sinclair/typebox` 定义参数 schema：

```typescript
const editSchema = Type.Object({
  path: Type.String({ description: 'Path to the file to edit' }),
  oldText: Type.String({ description: 'Exact text to find and replace' }),
  newText: Type.String({ description: 'New text to replace the old text with' })
})
```

- schema 既是 TypeScript 类型(通过 `Static<typeof editSchema>` 获取)
- 也是 JSON Schema(发给 LLM 做 function calling)
- agent-loop 中通过 `validateToolArguments(tool, toolCall)` 做运行时校验

---

## 二、Tool 注册与发现

### 2.1 内置 tool (coding-agent)

工厂模式，每个 tool 有两种形态：

```typescript
// 1. 默认实例 (用 process.cwd())
export const editTool = createEditTool(process.cwd());

// 2. 工厂函数 (指定 cwd，支持 options)
export function createEditTool(cwd: string, options?: EditToolOptions): AgentTool<...>

// 组合导出
export const codingTools: Tool[] = [readTool, bashTool, editTool, writeTool];
export const readOnlyTools: Tool[] = [readTool, grepTool, findTool, lsTool];
```

### 2.2 Extension 动态注册 (ToolDefinition)

Extension 系统提供了更强大的 `ToolDefinition`，比 AgentTool 多了：

- `promptSnippet` / `promptGuidelines`: 自动注入 system prompt
- `renderCall()` / `renderResult()`: 自定义 TUI 渲染
- `execute` 多了 `ctx: ExtensionContext` 参数

注册流程：

```
extension.registerTool(toolDef)
  → loader 存入 extension.tools Map
  → runtime.refreshTools() 通知 session
  → wrapRegisteredTools() 将 ToolDefinition 包装成 AgentTool
```

### 2.3 Extension 拦截内置 tool (wrapToolsWithExtensions)

Extension 还能 hook 内置 tool 的执行，做 before/after 拦截。

---

## 三、Agent Loop 中的 Tool 执行流程

```
runLoop()
  │
  ├─ streamAssistantResponse()
  │    ├─ transformContext (AgentMessage[] → AgentMessage[])   // 可选：压缩上下文
  │    ├─ convertToLlm (AgentMessage[] → Message[])           // 转换为 LLM 格式
  │    └─ streamSimple(model, context, options)                // 流式调用 LLM
  │
  ├─ 检查 toolCalls
  │
  └─ executeToolCalls()
       └─ for each toolCall (顺序执行，非并行):
            ├─ emit tool_execution_start
            ├─ validateToolArguments(tool, toolCall) // TypeBox 校验
            ├─ tool.execute(id, args, signal, onUpdate)
            │    └─ onUpdate → emit tool_execution_update (流式进度)
            ├─ emit tool_execution_end
            ├─ 构造 ToolResultMessage 加入 context
            └─ getSteeringMessages() → 如果用户中途打断，skip 后续 tool
```

**重要细节**:

1. **顺序执行**: tool 是一个一个执行的（for 循环），不是并行
2. **Steering 机制**: 每执行完一个 tool 就检查用户是否有新消息；如有则 skip 剩余 tool，立即响应用户
3. **Follow-up 机制**: 当所有 tool 执行完且无 steering 时，检查 `getFollowUpMessages()`
4. **信号传递**: AbortSignal 贯穿整个链路，支持随时取消

---

## 四、精读 Edit Tool (最核心的 tool)

### 4.1 为什么选 edit 精读？

Edit 是 coding agent 最高频使用的 tool，也是实现最复杂的：

- 需要处理精确文本匹配 + 模糊回退
- 需要处理各种编码边界(BOM, CRLF, Unicode)
- 需要生成 diff 给 UI 预览
- 是理解 Agent "如何改代码" 的关键

### 4.2 参数设计

```typescript
{
  path: string,     // 文件路径
  oldText: string,  // 要替换的精确文本
  newText: string   // 替换后的文本
}
```

这是 **search-and-replace** 模式，而非行号模式。优点：

- 对 LLM 更友好：LLM 输出的上下文天然是文本片段
- 避免行号漂移问题
- 要求 oldText 在文件中唯一，防止误编辑

### 4.3 核心执行流程

```
execute(path, oldText, newText)
  │
  ├─ 1. resolveToCwd(path, cwd)              // 路径解析
  │      处理 ~、@前缀、Unicode空格、NFD、macOS截图路径
  │
  ├─ 2. ops.access(absolutePath)              // 检查读写权限
  │
  ├─ 3. ops.readFile(absolutePath)            // 读文件
  │
  ├─ 4. stripBom(rawContent)                  // 剥离 BOM
  │      LLM 不会生成 BOM 字符，但文件可能有
  │
  ├─ 5. detectLineEnding(content)             // 检测 CRLF vs LF
  │      normalizeToLF(content/oldText/newText) // 统一为 LF 处理
  │
  ├─ 6. fuzzyFindText(content, oldText)       // ⭐ 匹配算法
  │      ├─ 先试 精确匹配 content.indexOf(oldText)
  │      └─ 失败则 模糊匹配:
  │           ├─ 去除行尾空白
  │           ├─ 智能引号 → ASCII 引号
  │           ├─ Unicode 破折号 → ASCII 连字符
  │           └─ 特殊空格 → 普通空格
  │
  ├─ 7. 检查唯一性: occurrences > 1 → 报错
  │      "Found N occurrences... Please provide more context"
  │
  ├─ 8. 执行替换:
  │      newContent = before + newText + after
  │
  ├─ 9. 检查是否真的改了: oldContent === newContent → 报错
  │
  ├─ 10. restoreLineEndings(newContent, originalEnding)
  │       恢复原始行尾 + 补回 BOM
  │
  ├─ 11. ops.writeFile(absolutePath, finalContent)
  │
  └─ 12. generateDiffString(oldContent, newContent)
          返回 { content: ["Successfully replaced..."], details: { diff, firstChangedLine } }
```

### 4.4 Fuzzy Match：两阶段匹配

这是 edit tool 最精巧的部分：

```typescript
function fuzzyFindText(content: string, oldText: string): FuzzyMatchResult {
  // 阶段1: 精确匹配
  const exactIndex = content.indexOf(oldText);
  if (exactIndex !== -1) {
    return { found: true, index: exactIndex, usedFuzzyMatch: false,
             contentForReplacement: content };  // 用原始内容做替换
  }

  // 阶段2: 模糊匹配 - 在标准化空间中匹配
  const fuzzyContent = normalizeForFuzzyMatch(content);
  const fuzzyOldText = normalizeForFuzzyMatch(oldText);
  const fuzzyIndex = fuzzyContent.indexOf(fuzzyOldText);

  if (fuzzyIndex !== -1) {
    return { found: true, index: fuzzyIndex, usedFuzzyMatch: true,
             contentForReplacement: fuzzyContent };  // 用标准化内容做替换!
  }

  return { found: false, ... };
}
```

**关键设计决策**: 当 fuzzy match 命中时，`contentForReplacement` 不是原始文件内容，而是标准化后的版本。这意味着替换操作会在标准化空间中进行，顺带修复了 trailing whitespace / Unicode 引号等格式问题。这是有意为之——"既然在修了，顺便修一下格式"。

### 4.5 Pluggable Operations 模式

所有 tool 都采用 Operations 接口实现可替换的底层操作：

```typescript
interface EditOperations {
  readFile: (absolutePath: string) => Promise<Buffer>
  writeFile: (absolutePath: string, content: string) => Promise<void>
  access: (absolutePath: string) => Promise<void>
}
```

默认实现用 `fs` 模块。但可以注入自定义实现（如 SSH 远程文件系统），完美符合**依赖反转原则**。

### 4.6 Abort 处理模式

每个 tool 都实现了完整的 abort 机制（统一模式）：

```typescript
execute: async (_toolCallId, params, signal?) => {
  return new Promise((resolve, reject) => {
    if (signal?.aborted) {
      reject(new Error('Operation aborted'))
      return
    }

    let aborted = false
    const onAbort = () => {
      aborted = true
      reject(new Error('Operation aborted'))
    }
    signal?.addEventListener('abort', onAbort, { once: true })

    ;(async () => {
      try {
        // ... 每个步骤前检查 if (aborted) return;
        // ... 执行实际操作
        signal?.removeEventListener('abort', onAbort)
        resolve(result)
      } catch (error) {
        signal?.removeEventListener('abort', onAbort)
        if (!aborted) reject(error)
      }
    })()
  })
}
```

### 4.7 Diff 生成

`generateDiffString` 使用 `diff` 库的 `diffLines`：

- 输出带行号的 unified diff 格式 (`+lineNum` / `-lineNum`)
- 上下文行数可配（默认 4 行）
- 返回 `firstChangedLine` 用于 TUI 编辑器跳转

### 4.8 Edit 工具的 TUI 预览 (computeEditDiff)

`computeEditDiff` 是 `edit-diff.ts` 导出的另一个函数，被 `tool-execution.ts` 在 tool **执行前** 调用，用于给用户显示即将发生的变更预览。逻辑与实际执行完全一致（读文件 → fuzzy match → 生成 diff），但不写入文件。

---

## 五、面试要点总结

| 维度                   | 要点                                                                 |
| ---------------------- | -------------------------------------------------------------------- |
| **类型分层**           | ai (schema) → agent (+ execute) → coding-agent (具体实现)            |
| **参数校验**           | TypeBox: 一份 schema 同时服务 TS 类型 + JSON Schema + 运行时校验     |
| **执行模型**           | 顺序执行，非并行；每 tool 后检查 steering message                    |
| **Abort 机制**         | AbortSignal 贯穿全链路，每步骤都检查                                 |
| **匹配策略**           | 两阶段: 精确 → 模糊(Unicode 标准化)；唯一性校验防误编辑              |
| **依赖反转**           | Operations 接口解耦 I/O，支持本地/SSH/mock                           |
| **content vs details** | content 给 LLM，details 给 UI，关注分离                              |
| **Extension**          | ToolDefinition 比 AgentTool 更强: 自定义渲染 + prompt 注入 + context |
| **截断策略**           | 统一的 truncate 工具: 2000 行 / 50KB 双限制                          |
| **编码处理**           | BOM 剥离/恢复, CRLF↔LF 标准化, Unicode 空格/引号/破折号              |
