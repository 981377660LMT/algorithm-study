# Debug Adapter Protocol 新主页

链接：https://code.visualstudio.com/blogs/2018/08/07/debug-adapter-protocol-website

## 一针见血的分析

1. **调试领域的 LSP 时刻**：DAP 的核心价值在于将调试逻辑与 UI 分离。在 DAP 出现之前，调试器开发者必须为每个 IDE 编写特定的适配逻辑（M 种调试器 * N 种 IDE = M*N 的复杂度）。DAP 将其简化为 M+N。
2. **协议的独立性**：将 DAP 移出 VS Code 的代码仓库并建立独立网站，是 VS Code 团队「去 VS Code 化」的战略动作。这标志着协议已经成熟到可以作为行业标准，被 Visual Studio、Emacs、Vim 等其他工具原生支持。
3. **架构的稳定性**：DAP 早期受 V8 调试协议启发，使用 JSON 格式。虽然与 LSP (JSON-RPC) 略有不同，但其定义了一套完整的生命周期（Initializing, Configuration, Launching/Attaching, Stopping），奠定了现代 IDE 通用调试架构的基础。

## 摘要

VS Code 团队宣布为**调试适配器协议 (Debug Adapter Protocol, DAP)** 建立新的官方网站和独立代码仓库。

- **解耦逻辑**：DAP 使前端 UI 与后端调试器实现分离，开发者只需编写一次调试适配器。
- **功能特性**：支持断点（源码、函数、条件等）、变量查看、堆栈跟踪、多进程调试以及交互式 REPL。
- **行业标准**：DAP 不再仅限于 VS Code，Visual Studio 等其他 IDE 也已开始支持。
- **新家地址**：文档、协议规范及 SDK 迁至 [https://microsoft.github.io/debug-adapter-protocol](https://microsoft.github.io/debug-adapter-protocol/)。
