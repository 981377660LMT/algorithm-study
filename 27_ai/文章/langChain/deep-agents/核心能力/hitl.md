这份文档介绍了 **Deep Agents** 的“人在回路” (Human-in-the-loop) 功能。
该功能允许在执行敏感工具操作（如删除文件、发送邮件）之前暂停 Agent，等待人工批准、修改或拒绝。

以下是配置和使用指南的总结：

### 1. 核心配置 (`interrupt_on`)

通过 `interrupt_on` 参数配置哪些工具需要中断。
**注意**：使用此功能**必须**配置 `checkpointer`（如 `MemorySaver`）以保存中断时的状态。

配置方式支持：
*   `true`: 启用默认中断（允许批准、编辑、拒绝）。
*   `false`: 禁用中断。
*   对象配置: 指定允许的决策类型。

```typescript
import { createDeepAgent, MemorySaver } from "deepagents";

const agent = createDeepAgent({
  // ... 其他配置
  tools: [deleteFile, sendEmail],
  checkpointer: new MemorySaver(), // 必须项！
  interruptOn: {
    delete_file: true, // 默认：允许 approve, edit, reject
    read_file: false,  // 不中断
    send_email: { allowedDecisions: ["approve", "reject"] }, // 禁止编辑参数
  },
});
```

### 2. 决策类型

人工审核时可以采取以下三种行动：
*   **`"approve"`**: 按原参数执行工具。
*   **`"edit"`**: 修改工具参数后执行。
*   **`"reject"`**: 拒绝执行（跳过该工具调用）。

### 3. 处理流程与恢复

处理中断分为三个步骤：调用、检测中断、恢复执行。

```typescript
import { Command } from "@langchain/langgraph";

// 1. 初始调用 (必须带 thread_id)
const config = { configurable: { thread_id: "unique-id" } };
let result = await agent.invoke({ messages: [{ role: "user", content: "Delete file.txt" }] }, config);

// 2. 检测中断
if (result.__interrupt__) {
  const interrupts = result.__interrupt__[0].value;
  const actionRequests = interrupts.actionRequests; // 待审核的操作列表

  // 3. 构建决策
  // 注意：decisions 数组的顺序必须与 actionRequests 一致
  const decisions = actionRequests.map(action => {
    // 示例逻辑：批准所有操作
    return { type: "approve" };
    
    // 如果是编辑，结构如下：
    // return {
    //   type: "edit",
    //   editedAction: { name: action.name, args: { ...newArgs } }
    // };
  });

  // 4. 恢复执行 (使用 Command)
  result = await agent.invoke(
    new Command({ resume: { decisions } }),
    config // 必须使用相同的 config
  );
}
```

### 4. 高级特性

*   **批量审批**: 如果 Agent 一次调用了多个敏感工具，中断会包含所有请求。恢复时提供的 `decisions` 数组必须与请求顺序一一对应。
*   **子智能体覆盖**: 子智能体可以拥有独立的 `interrupt_on` 配置，覆盖主智能体的设置。

### 5. 最佳实践

1.  **始终使用 Checkpointer**: 没有它无法保存暂停时的状态。
2.  **保持 Thread ID 一致**: 恢复执行时必须使用与初始调用相同的 `thread_id`。
3.  **风险分级**: 对高风险操作（如删除数据）开启完全中断，对中低风险操作可仅开启部分权限或不开启。