# 推介 Markdown 语言服务器

链接：https://code.visualstudio.com/blogs/2022/08/16/markdown-language-server

## 摘要

通过六年的演进，VS Code 的 Markdown 支持从简单的语法高亮发展到了完整的语言服务。2022 年，VS Code 将 Markdown 逻辑迁移到了独立的语言服务器（LSP）中。本文介绍了 Markdown Language Server 的动机：支持更复杂的链接重命名、引用查找和断链检测，同时利用 LSP 的跨进程特性避免阻塞扩展宿主，并为 Monaco 等其他编辑器提供通用的 Markdown 编辑智能。

## 一针见血的分析

将 Markdown 从传统的 Regex 驱动升级到 **LSP 驱动**，是 VS Code 将“非代码资产”作为“一等公民”处理的里程碑。在这个转变中，最重要的工程实践是**渐进式迁移（Incremental Migration）**：Matt Bierner 通过在 `main` 分支上逐步迁移功能，实现了在不破坏生产环境的前提下，完成了从插件 API 到 LSP Pull Diagnostics 模型的平滑过渡。这不仅降低了渲染进程的压力（主线程不再负责链接验证），更通过开源 `vscode-markdown-languageservice` 建立了一个跨工具的 Markdown 生态标准。这再次印证了 VS Code 的核心哲学：通过标准化的协议（LSP）将通用的编辑智能从特定的编辑器实现中解耦出来。
