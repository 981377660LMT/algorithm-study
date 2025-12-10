这份文档介绍了如何在 LangChain.js 中集成 **Model Context Protocol (MCP)**。

**MCP** 是一个开放协议，旨在标准化应用程序向 LLM 提供工具和上下文的方式。通过 `@langchain/mcp-adapters` 库，LangChain Agent 可以轻松使用定义在 MCP 服务器上的工具。

以下是核心内容总结：

### 1. 安装依赖
首先需要安装适配器库：
```bash
npm install @langchain/mcp-adapters
```

### 2. 连接 MCP 服务器
使用 `MultiServerMCPClient` 可以同时连接多个 MCP 服务器。支持两种主要的传输方式（Transports）：

*   **stdio**: 启动本地子进程进行通信（适用于本地工具）。
*   **sse (Server-Sent Events) / http**: 通过 HTTP 请求通信（适用于远程服务器）。

```typescript
import { MultiServerMCPClient } from "@langchain/mcp-adapters";
import { createAgent } from "langchain";

const client = new MultiServerMCPClient({
    // 1. 本地 Math 服务器 (stdio)
    math: {
        transport: "stdio",
        command: "node",
        args: ["/path/to/math_server.js"],
    },
    // 2. 远程 Weather 服务器 (HTTP/SSE)
    weather: {
        transport: "sse",
        url: "http://localhost:8000/mcp",
    },
});

// 获取所有服务器上的工具
const tools = await client.getTools();

// 创建 Agent 并传入工具
const agent = createAgent({
    model: "claude-sonnet-4-5-20250929",
    tools, 
});
```

### 3. 创建自定义 MCP 服务器
你可以使用 `@modelcontextprotocol/sdk` 来构建自己的 MCP 服务器。

*   **定义工具**：设置工具名称、描述和输入 Schema（Zod 风格）。
*   **处理请求**：编写 `CallToolRequestSchema` 的处理逻辑。
*   **连接传输层**：使用 `StdioServerTransport` 或 `SSEServerTransport` 启动服务。

### 4. 关键注意事项
*   **无状态 (Stateless)**：`MultiServerMCPClient` 默认是无状态的。每次工具调用都会创建一个新的 MCP `ClientSession`，执行完毕后立即清理。
*   **工具转换**：LangChain 会自动将 MCP 工具转换为 LangChain 标准工具格式，使其可以直接在 Agent 中使用。