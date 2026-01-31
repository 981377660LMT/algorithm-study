# 超越工具，在 VS Code 中集成 MCP

链接：https://code.visualstudio.com/blogs/2025/05/12/agent-mode-meets-mcp

## 摘要

Model Context Protocol (MCP) 是 2025 年 AI 领域的重要标准。VS Code 通过深度集成 MCP，允许 Agent 模式下的 Copilot 发现并调用各种外部工具。本文详细介绍了 MCP 的传输协议（stdio/sse）、配置方式，以及它如何帮助开发者将本地脚本、数据库查询甚至企业内部 API 无缝开放给 AI 助手，从而构建出高度定制化的“特种 Agent”。

## 一针见血的分析

MCP 对 AI 时代的意义，正如同 LSP 对云原生开发工具链的意义。VS Code 的这一举措确立了 **“工具插槽（Tool Slots）”** 的标准模型。通过支持 stdio 和 sse 传输，VS Code 让“写一个为 AI 后台服务的脚本”变得像写一个命令行程序一样简单。这极大地缓解了 LLM 模型本身的“知识滞后”和“私有环境盲区”问题。在工程实践上，这实现了**能力供给（Tool Providers）与能力消费（AI Agents）的松耦合**，使得开发者不再受限于 GitHub Copilot 官方提供的固定工具集，而是能够利用开源社区中海量的 MCP 服务器快速组装出垂直行业的生产力引擎。

## 深入分析

### 1. 协议的胜利：Agent 模式的初次会师

这篇文章预示了 VS Code 将 MCP (Model Context Protocol) 选定为 AI 通讯的标准基座。当 **Agent 模式遇见 MCP**，它意味着 Agent 真正获得了“感知外部世界”的能力。

### 2. 打破 IDE 的边界

在此之前，AI 只能“看”项目内的文件。通过集成 MCP，Agent 可以：

- 实时访问外部 API 文档。
- 查询实时运行的本地或远程数据库。
- 甚至驱动浏览器进行端到端测试。
  这种集成为 Agent 提供了除了代码之外的**动态、实时、领域特定**的上下文。

### 3. 生态开放性

微软并没有试图发明一套私有的工具调用标准，而是押注开源的 MCP 协议。这种姿态极大地降低了插件开发者为了适配不同 AI 编辑器而进行重复劳动的门槛。
