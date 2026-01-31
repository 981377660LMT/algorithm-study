# 完整的 MCP 体验：VS Code 中的全规范支持

链接：https://code.visualstudio.com/blogs/2025/06/12/full-mcp-spec-support

## 深入分析

### 1. 从“黑屏工具”到“语义服务”

VS Code 对 MCP (Model Context Protocol) 全规范的支持，标志着 AI 外部能力的标准被彻底刷新。除了基础的 Tool Calling，以下三个原语的加入至关重要：

- **Prompts (动态工作流)**：MCP Server 不再只是提供单次调用的工具，而是可以向编辑器注入标准化的、上下文感知的工作流。
- **Resources (语义资产)**：Agent 可以像操作本地文件一样，无缝地读写外部资源（如数据库记录、GitHub 讨论、甚至 Playwright 的浏览器截图）。

### 2. sampling：Agent 协同的关键

Sampling (采样) 功能解决了多 Agent 协作中的“核心密钥管理”难题。MCP Server 不再需要内置庞大的 AI SDK，而是可以反向调用主编辑器配置好的模型。这种设计不仅安全（密钥不离编辑器），还让层级化的 Agent 协同变为可能。

### 3. 安全与身份验证的基座

引入标准化的 **Authorization 规范**，解决了连接企业级、带权限数据源的痛点。通过将认证委派给现有的身份提供商，MCP 终于可以承载敏感的生产力数据模型。
