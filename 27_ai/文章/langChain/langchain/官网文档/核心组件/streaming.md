这份文档详细介绍了 **LangChain.js 中的流式传输 (Streaming)** 机制。流式传输对于提升 LLM 应用的响应速度和用户体验（UX）至关重要。

以下是文档的核心知识点分析：

### 1. 为什么需要流式传输？
*   **响应性**：LLM 生成完整回复需要时间，流式传输可以让用户在结果生成完成前就看到部分内容（如“打字机”效果）。
*   **透明度**：展示 Agent 的思考过程和工具调用状态。

### 2. 三种核心流式模式 (Stream Modes)
LangChain 通过 `agent.stream(input, { streamMode: ... })` 方法支持不同的流式粒度：

*   **`streamMode: "updates"` (Agent 进度)**
    *   **作用**：在 Agent 的每个**步骤**（Step）完成后触发事件。
    *   **粒度**：较粗。例如，模型生成完完整回复后触发一次，工具执行完触发一次。
    *   **内容**：包含该步骤产生的状态更新（如 `AIMessage` 或 `ToolMessage`）。

*   **`streamMode: "messages"` (LLM Token)**
    *   **作用**：实时流式传输 LLM 生成的每一个 **Token**。
    *   **粒度**：最细。
    *   **内容**：包含 Token 内容和元数据（如 `langgraph_node` 标识来源）。

*   **`streamMode: "custom"` (自定义更新)**
    *   **作用**：允许在**工具内部**手动发送自定义状态更新。
    *   **场景**：工具执行时间较长，需要反馈进度（如“正在查询数据库...”、“已获取 50% 数据”）。
    *   **实现**：在工具函数中使用 `config.writer("消息")` 发送数据。

### 3. 混合模式
你可以同时启用多种模式，例如 `streamMode: ["updates", "messages"]`，以便同时获取宏观的步骤更新和微观的 Token 生成。

### 代码示例：自定义流式更新
这是文档中比较独特的一个功能，允许工具与 UI 进行实时交互：

```typescript
const getWeather = tool(
  async (input, config) => {
    // 使用 writer 发送中间状态
    config.writer?.(`正在查找 ${input.city} 的数据...`);
    // ... 执行耗时操作
    config.writer?.(`数据获取成功！`);
    return `天气晴朗`;
  },
  // ... 配置
);

// 调用时启用 custom 模式
await agent.stream(..., { streamMode: "custom" });
```

### 总结
*   如果只需要显示最终结果的打字机效果，用 `"messages"`。
*   如果需要显示“正在思考”、“正在调用工具”，用 `"updates"`。
*   如果需要工具内部的详细进度条，用 `"custom"`。