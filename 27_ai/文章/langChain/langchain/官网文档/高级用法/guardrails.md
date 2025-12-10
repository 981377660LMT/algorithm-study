这两份文档详细介绍了 **LangChain.js** 中的 **Guardrails (护栏)** 和 **Custom Middleware (自定义中间件)**。

简单来说，**Guardrails** 是目的（为了安全和合规），而 **Middleware** 是实现这一目的的主要手段（通过拦截和控制执行流）。

以下是对这两部分核心内容的详细解读：

### 1. Guardrails (安全护栏)

护栏用于验证和过滤 Agent 的输入与输出，确保应用安全、合规。

#### 两种实现策略
*   **确定性护栏 (Deterministic)**：基于规则（如正则、关键词匹配）。
    *   *优点*：速度快、成本低、结果可预测。
    *   *缺点*：可能漏掉语义上的违规。
*   **模型级护栏 (Model-based)**：使用 LLM 进行语义评估。
    *   *优点*：能理解上下文，捕捉隐晦的违规。
    *   *缺点*：速度慢、成本高。

#### 内置护栏
LangChain 提供了开箱即用的中间件作为护栏：
*   **PII Detection (敏感信息检测)**：检测邮箱、信用卡号等。支持 `redact` (替换)、`mask` (掩码)、`block` (阻断) 等策略。
*   **Human-in-the-loop (人机协同)**：在执行敏感操作（如转账、删库）前暂停，等待人工批准。

#### 自定义护栏
你可以通过中间件钩子在不同阶段实施护栏：
*   **Before Agent**：在处理开始前拦截（如：检查用户输入是否包含恶意关键词）。
*   **After Agent**：在返回结果前拦截（如：检查 AI 生成的内容是否安全）。

---

### 2. Custom Middleware (自定义中间件)

中间件通过 **Hooks (钩子)** 介入 Agent 的执行流程。你可以使用 `createMiddleware` 函数来创建。

#### 两种钩子风格
1.  **Node-style Hooks (节点式)**：在特定点顺序执行。适用于日志记录、验证和状态更新。
    *   `beforeAgent` / `afterAgent`：Agent 启动前/结束后（单次调用仅触发一次）。
    *   `beforeModel` / `afterModel`：模型调用前后。
2.  **Wrap-style Hooks (包裹式)**：包裹模型或工具的调用。适用于重试、缓存、动态修改参数。
    *   `wrapModelCall`：拦截并控制模型调用。
    *   `wrapToolCall`：拦截并控制工具调用。

#### 状态与上下文 (State & Context)
中间件可以扩展 Agent 的数据结构：
*   **State (`stateSchema`)**：
    *   用于在 Agent 执行生命周期内持久化数据（如计数器）。
    *   支持私有字段（以 `_` 开头），不会包含在最终结果中。
*   **Context (`contextSchema`)**：
    *   **只读**的元数据，每次调用 `invoke` 时传入（如 `userId`, `tenantId`）。
    *   用于传递配置或用户信息，不会跨调用持久化。

#### 执行顺序
当有多个中间件时（如 `[m1, m2, m3]`）：
*   **Before 钩子**：顺序执行 (m1 -> m2 -> m3)。
*   **Wrap 钩子**：嵌套执行 (m1 包裹 m2，m2 包裹 m3)。
*   **After 钩子**：**逆序**执行 (m3 -> m2 -> m1)。

#### 流程控制 (Jumps)
中间件可以通过返回 `jumpTo` 来改变执行流：
*   `jumpTo: "end"`：立即结束 Agent 执行（常用于阻断违规请求）。
*   `jumpTo: "model"` / `jumpTo: "tools"`：跳转到特定节点。

### 代码示例：创建一个简单的关键词过滤护栏

这是一个结合了 Guardrails 理念和 Custom Middleware 技术的示例：

```typescript
import { createAgent, createMiddleware, AIMessage } from "langchain";

// 创建一个自定义中间件作为护栏
const keywordGuardrail = (bannedWords: string[]) => {
  return createMiddleware({
    name: "KeywordGuardrail",
    // 使用 beforeModel 钩子在模型调用前检查
    beforeModel: {
      canJumpTo: ["end"], // 声明该钩子有权结束执行
      hook: (state) => {
        const lastMessage = state.messages.at(-1);
        if (!lastMessage || lastMessage._getType() !== "human") return;

        const content = lastMessage.content.toString().toLowerCase();
        
        // 检查是否包含违禁词
        const hasBannedWord = bannedWords.some(word => content.includes(word));

        if (hasBannedWord) {
          return {
            // 返回拒绝消息
            messages: [new AIMessage("I cannot discuss that topic.")],
            // 立即结束，不调用 LLM
            jumpTo: "end",
          };
        }
        return;
      },
    },
  });
};

const agent = createAgent({
  model: "gpt-4o",
  tools: [],
  // 注册中间件
  middleware: [keywordGuardrail(["hack", "exploit"])],
});

// 测试
const result = await agent.invoke({
    messages: [{ role: "user", content: "How to hack a server?" }]
});

console.log(result.messages.at(-1).content); 
// 输出: "I cannot discuss that topic."
```