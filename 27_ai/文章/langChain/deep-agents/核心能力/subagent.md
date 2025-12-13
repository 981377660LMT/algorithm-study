这份文档详细介绍了 **Deep Agents** 中的 **Subagents (子智能体)** 机制。子智能体是实现任务委派和保持主智能体上下文整洁（Context Quarantine）的关键功能。

以下是核心概念和使用指南的总结：

### 1. 为什么要使用子智能体？
*   **解决上下文膨胀 (Context Bloat)**: 复杂的任务（如网络搜索、数据库查询）会产生大量的中间步骤和工具输出。子智能体在隔离的环境中执行这些步骤，只将最终结果返回给主智能体。
*   **适用场景**: 多步骤任务、需要特定领域指令的任务、需要不同模型能力的场景。
*   **不适用场景**: 简单的单步任务，或需要保留中间过程上下文的场景。

### 2. 配置子智能体
可以通过 `subagents` 参数传递子智能体列表。支持两种定义方式：

#### A. 字典配置 (推荐用于大多数场景)
定义一个包含配置信息的对象。
*   **必填**: `name` (唯一标识), `description` (主智能体据此决定是否调用), `system_prompt` (指令), `tools` (工具集)。
*   **可选**: `model` (覆盖主模型), `middleware`, `interrupt_on`。

```typescript
const researchSubagent = {
  name: "research-agent",
  description: "Used to research more in depth questions", // 描述要具体
  systemPrompt: "You are a great researcher...", // 指令要包含输出格式要求
  tools: [internetSearch],
  model: new ChatAnthropic({ model: "claude-sonnet-4-5-20250929" }), // 可选：使用特定模型
};

const agent = createDeepAgent({
  // ...
  subagents: [researchSubagent],
});
```

#### B. CompiledSubAgent (高级场景)
传入一个预编译的 LangGraph 图 (`Runnable`)。适用于需要自定义复杂工作流的场景。

```typescript
const customSubagent = {
  name: "data-analyzer",
  description: "Specialized agent for complex data analysis",
  runnable: customGraph, // 预先编译好的 LangGraph 实例
};
```

### 3. 内置通用子智能体 (General-Purpose)
Deep Agents 始终包含一个名为 `general-purpose` 的子智能体。
*   **特性**: 继承主智能体的系统提示词、工具和模型。
*   **用途**: 适用于不需要特殊指令，但需要隔离上下文的通用多步任务。

### 4. 最佳实践
1.  **描述 (Description) 要清晰**: 主智能体依赖描述来路由任务。避免 "helps with stuff"，应使用 "Analyzes financial data..."。
2.  **提示词 (Prompt) 要详细**: 指导子智能体如何使用工具，并强制要求返回**摘要**而非原始数据。
3.  **最小化工具集**: 只给子智能体完成任务所需的工具，提高聚焦度和安全性。
4.  **返回简洁结果**: 明确指示子智能体不要返回中间过程的详细日志，只返回最终结论。

### 5. 常见问题排查
*   **不委派任务**: 检查 `description` 是否不够具体，或在主提示词中显式要求使用 `task()` 工具。
*   **上下文依然膨胀**: 检查子智能体的提示词，确保它被要求返回“摘要”。对于大数据量，指示其写入文件系统而非直接返回文本。
*   **选错子智能体**: 区分不同子智能体的 `description`，使其职责边界更清晰。