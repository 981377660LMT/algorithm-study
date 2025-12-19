这份指南介绍了如何使用 LangChain 构建一个 **Supervisor (监督者)** 多智能体架构。在这种模式下，一个中心化的监督者智能体负责协调多个专门的子智能体（Worker），每个子智能体专注于特定的领域（如日历管理、邮件发送）。

以下是基于文档内容的完整 TypeScript 实现代码。

### 1. 安装依赖

```bash
npm install langchain @langchain/openai zod
```

### 2. 完整代码示例

这段代码实现了三层架构：
1.  **底层工具**：执行具体 API 操作（如创建日程、发送邮件）。
2.  **子智能体**：将自然语言转化为工具调用。
3.  **监督者智能体**：接收用户复杂指令，将其拆解并分发给子智能体。

```typescript
import { tool } from "@langchain/core/tools";
import { createAgent } from "langchain/agents";
import { ChatOpenAI } from "@langchain/openai";
import * as z from "zod";

// --- 1. 定义底层工具 (Mock) ---

// 日历工具：创建事件
const createCalendarEvent = tool(
  async ({ title, startTime, endTime, attendees }) => {
    return `Event created: ${title} from ${startTime} to ${endTime} with ${attendees.length} attendees`;
  },
  {
    name: "create_calendar_event",
    description: "Create a calendar event. Requires exact ISO datetime format.",
    schema: z.object({
      title: z.string(),
      startTime: z.string().describe("ISO format: '2024-01-15T14:00:00'"),
      endTime: z.string().describe("ISO format: '2024-01-15T15:00:00'"),
      attendees: z.array(z.string()).describe("email addresses"),
    }),
  }
);

// 邮件工具：发送邮件
const sendEmail = tool(
  async ({ to, subject, body }) => {
    return `Email sent to ${to.join(", ")} - Subject: ${subject}`;
  },
  {
    name: "send_email",
    description: "Send an email via email API.",
    schema: z.object({
      to: z.array(z.string()).describe("email addresses"),
      subject: z.string(),
      body: z.string(),
    }),
  }
);

// --- 2. 创建专门的子智能体 ---

// 初始化 LLM
const llm = new ChatOpenAI({ model: "gpt-4o", temperature: 0 });

// 日历子智能体：负责处理时间解析和日程安排
const calendarAgent = createAgent({
  model: llm,
  tools: [createCalendarEvent],
  systemPrompt: `
    You are a calendar scheduling assistant.
    Parse natural language scheduling requests into proper ISO datetime formats.
    Use create_calendar_event to schedule events.
    Always confirm what was scheduled in your final response.
  `.trim(),
});

// 邮件子智能体：负责撰写和发送邮件
const emailAgent = createAgent({
  model: llm,
  tools: [sendEmail],
  systemPrompt: `
    You are an email assistant.
    Compose professional emails based on natural language requests.
    Use send_email to send the message.
    Always confirm what was sent in your final response.
  `.trim(),
});

// --- 3. 将子智能体封装为工具 ---
// 这一步是 Supervisor 模式的关键：监督者看到的不是底层 API，而是"能力"

const scheduleEventTool = tool(
  async ({ request }) => {
    // 调用日历子智能体
    const result = await calendarAgent.invoke({
      messages: [{ role: "user", content: request }],
    });
    // 返回子智能体的最终回复
    const lastMessage = result.messages[result.messages.length - 1];
    return lastMessage.content;
  },
  {
    name: "schedule_event",
    description: `
      Schedule calendar events using natural language.
      Use this when the user wants to create, modify, or check calendar appointments.
      Input: Natural language scheduling request.
    `.trim(),
    schema: z.object({
      request: z.string().describe("Natural language scheduling request"),
    }),
  }
);

const manageEmailTool = tool(
  async ({ request }) => {
    // 调用邮件子智能体
    const result = await emailAgent.invoke({
      messages: [{ role: "user", content: request }],
    });
    const lastMessage = result.messages[result.messages.length - 1];
    return lastMessage.content;
  },
  {
    name: "manage_email",
    description: `
      Send emails using natural language.
      Use this when the user wants to send notifications or reminders.
      Input: Natural language email request.
    `.trim(),
    schema: z.object({
      request: z.string().describe("Natural language email request"),
    }),
  }
);

// --- 4. 创建 Supervisor (监督者) 智能体 ---

const supervisorAgent = createAgent({
  model: llm,
  tools: [scheduleEventTool, manageEmailTool], // 监督者只拥有高级工具
  systemPrompt: `
    You are a helpful personal assistant.
    You can schedule calendar events and send emails.
    Break down user requests into appropriate tool calls and coordinate the results.
    When a request involves multiple actions, use multiple tools in sequence.
  `.trim(),
});

// --- 5. 运行示例 ---

async function run() {
  // 一个包含两个不同领域任务的复杂请求
  const query =
    "Schedule a meeting with the design team next Tuesday at 2pm for 1 hour, " +
    "and send them an email reminder about reviewing the new mockups.";

  console.log(`User Query: ${query}\n`);

  const stream = await supervisorAgent.stream({
    messages: [{ role: "user", content: query }],
  });

  for await (const step of stream) {
    // 打印每一步的执行情况
    for (const update of Object.values(step)) {
      if (update && typeof update === "object" && "messages" in update) {
        // @ts-ignore
        for (const message of update.messages) {
           console.log(`[${message._getType()}]: ${message.content}`);
        }
      }
    }
  }
}

run().catch(console.error);
```

### 核心架构解析

1.  **关注点分离 (Separation of Concerns)**:
    *   **日历 Agent** 不需要知道如何写邮件，它只专注于解析 "next Tuesday" 这种模糊的时间描述并调用 API。
    *   **邮件 Agent** 不需要知道时间安排，它只专注于生成得体的邮件内容。
    *   **Supervisor** 不需要知道具体的 API 参数格式（如 ISO 时间格式），它只需要知道"这件事该交给日历 Agent 还是邮件 Agent"。

2.  **工具封装 (Wrapping Agents as Tools)**:
    *   代码中的 `scheduleEventTool` 和 `manageEmailTool` 是连接 Supervisor 和 Sub-Agents 的桥梁。
    *   它们接收自然语言字符串 (`request`) 作为输入，这使得 Supervisor 可以直接用自然语言给子智能体下达指令。

3.  **工作流**:
    *   用户输入 -> Supervisor 解析 -> 调用 `schedule_event` -> 日历 Agent 解析并执行 -> 返回结果给 Supervisor -> Supervisor 调用 `manage_email` -> 邮件 Agent 执行 -> 返回结果 -> Supervisor 汇总回复用户。