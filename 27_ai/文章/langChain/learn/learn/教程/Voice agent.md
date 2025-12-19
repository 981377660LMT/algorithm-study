这份文档详细介绍了如何使用 LangChain 构建一个基于 **"三明治架构" (Sandwich Architecture)** 的语音 Agent。这种架构由三个独立但流水线式连接的组件组成：**语音转文本 (STT)** -> **LangChain Agent** -> **文本转语音 (TTS)**。

以下是基于文档内容的完整实现指南和代码总结。

### 核心架构

整个系统是一个流式处理管道（Streaming Pipeline），数据在各个阶段异步流动，以实现低延迟（目标 < 700ms）。

1.  **STT (Speech-to-Text)**: 接收音频流，实时生成文本。
2.  **Agent**: 接收文本，进行推理（可能调用工具），流式输出回复文本。
3.  **TTS (Text-to-Speech)**: 接收回复文本，实时合成音频流返回给客户端。

### 1. 环境准备

你需要安装 LangChain 核心库以及用于处理流和 WebSocket 的工具。

```bash
npm install langchain @langchain/core @langchain/langgraph @langchain/anthropic zod uuid ws
# 还需要安装你选择的 STT/TTS 提供商的 SDK 或使用 WebSocket 适配器
```

### 2. 完整代码实现

为了简化展示，这里将三个核心阶段串联起来。假设你已经有了 `AssemblyAISTT` 和 `CartesiaTTS` 的适配器类（参考文档中的实现）。

#### 第一步：定义 Agent (大脑)

Agent 负责接收文本输入，维护对话记忆，并生成回复。

```typescript
import { createAgent } from "langchain/agents"; // 假设使用高级 API
import { ChatAnthropic } from "@langchain/anthropic";
import { HumanMessage } from "@langchain/core/messages";
import { MemorySaver } from "@langchain/langgraph";
import { tool } from "@langchain/core/tools";
import { z } from "zod";
import { v4 as uuidv4 } from "uuid";
import type { VoiceAgentEvent } from "./types"; // 假设定义了事件类型

// 1. 定义工具
const addToOrder = tool(
  async ({ item, quantity }) => {
    return `Added ${quantity} x ${item} to the order.`;
  },
  {
    name: "add_to_order",
    description: "Add an item to the customer's sandwich order.",
    schema: z.object({
      item: z.string(),
      quantity: z.number(),
    }),
  }
);

// 2. 创建 Agent
const model = new ChatAnthropic({ model: "claude-3-haiku-20240307" });
const agent = createAgent({
  model,
  tools: [addToOrder],
  checkpointer: new MemorySaver(),
  systemPrompt: `You are a helpful sandwich shop assistant. 
  Be concise. Do NOT use emojis or markdown. 
  Your response is for text-to-speech.`,
});

// 3. Agent 流处理函数
export async function* agentStream(
  eventStream: AsyncIterable<VoiceAgentEvent>
): AsyncGenerator<VoiceAgentEvent> {
  const threadId = uuidv4(); // 每个会话一个 ID

  for await (const event of eventStream) {
    // 透传上游事件 (如 STT 中间结果)
    yield event;

    // 当收到最终 STT 结果时，触发 Agent
    if (event.type === "stt_output") {
      const stream = await agent.stream(
        { messages: [new HumanMessage(event.transcript)] },
        {
          configurable: { thread_id: threadId },
          streamMode: "messages", // 关键：流式输出 Token
        }
      );

      for await (const chunk of stream) {
        // 假设 chunk 结构包含 text
        if (chunk.content) {
             yield { type: "agent_chunk", text: chunk.content as string, ts: Date.now() };
        }
      }
    }
  }
}
```

#### 第二步：构建管道 (Pipeline)

将 STT、Agent 和 TTS 串联起来。这里使用异步生成器（Async Generators）来实现背压和流式传输。

```typescript
import { sttStream } from "./stt";   // 封装了 AssemblyAI
import { agentStream } from "./agent";
import { ttsStream } from "./tts";   // 封装了 Cartesia
import { writableIterator } from "./utils"; // 一个简单的可写迭代器工具

// 模拟 WebSocket 处理函数 (例如使用 Hono 或 Express WS)
async function handleWebSocketConnection(ws: WebSocket) {
  // 1. 创建输入流，用于接收客户端音频
  const inputStream = writableIterator<Uint8Array>();

  // 2. 构建处理管道：Audio -> STT -> Agent -> TTS -> Audio
  // 每一个步骤都消费上一步的输出流
  const transcriptEventStream = sttStream(inputStream);
  const agentEventStream = agentStream(transcriptEventStream);
  const outputEventStream = ttsStream(agentEventStream);

  // 3. 启动输出循环：将最终的 TTS 音频发回客户端
  const processOutput = async () => {
    for await (const event of outputEventStream) {
      if (event.type === "tts_chunk") {
        // 发送音频二进制数据给客户端播放
        ws.send(event.audio);
      } else if (event.type === "stt_output") {
        // 可选：发送识别到的文本给客户端展示
        ws.send(JSON.stringify({ type: "transcript", text: event.transcript }));
      }
    }
  };

  // 启动处理
  processOutput().catch(console.error);

  // 4. 处理客户端发来的消息 (音频块)
  ws.on("message", (data) => {
    if (data instanceof Buffer) {
      inputStream.push(new Uint8Array(data));
    }
  });

  ws.on("close", () => {
    inputStream.end(); // 关闭流
  });
}
```

### 关键技术点总结

1.  **流式优先 (Streaming First)**:
    *   **STT**: 使用 WebSocket 实时发送音频块，而不是上传整个文件。
    *   **Agent**: 使用 `streamMode="messages"`，一旦 LLM 生成第一个 Token 就立即推送到下游。
    *   **TTS**: 使用支持流式输入的 TTS 服务（如 Cartesia, ElevenLabs Streaming），一旦收到 Agent 的部分文本就开始合成音频。

2.  **并发处理 (Producer-Consumer)**:
    *   在 `sttStream` 和 `ttsStream` 中，通常需要同时做两件事：向外部服务发送数据（Producer）和从外部服务接收数据（Consumer）。使用 `Promise.all` 确保两者并发运行。

3.  **状态管理**:
    *   使用 LangGraph 的 `MemorySaver` 或其他 Checkpointer 来保存对话历史，确保 Agent 知道上下文（例如用户之前点了什么三明治）。

4.  **提示词工程 (Prompt Engineering)**:
    *   System Prompt 必须明确指示 Agent **不要使用 Emoji 或 Markdown**，因为 TTS 引擎通常无法正确朗读这些字符，或者读出来很奇怪。
    *   要求 Agent 回复**简洁**，因为语音交互的听觉带宽比视觉阅读低。