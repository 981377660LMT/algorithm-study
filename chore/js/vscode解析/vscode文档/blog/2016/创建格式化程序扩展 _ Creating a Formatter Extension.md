---
title: 创建格式化程序扩展
date: 2016/11/15
url: https://code.visualstudio.com/blogs/2016/11/15/formatters-best-practices
---

## 深入分析

这篇技术指南揭示了 VS Code “简洁核心”的设计理念，重点讨论了如何通过标准化 API 提供一致的用户体验。

1.  **API vs 自定义命令**：
    - **错误做法**：将格式化实现为独立的 Command（如 `extension.formatFoo`），这会导致用户无法通过统一的“Format Document”快捷键调用。
    - **推荐做法**：注册 `DocumentFormattingEditProvider`。这样做可以让扩展自动接入 VS Code 的核心功能，如 **Format on Save** 和 **Format on Paste**。
2.  **多格式化器竞争机制**：
    - VS Code 允许一个语言拥有多个格式化器，但也指明了最佳实践：扩展应该提供开关配置（如 `html.format.enable`），允许用户手动禁用冲突的默认格式化器。
3.  **用户体验的一致性**：
    - 通过标准化 API，用户无论在写什么语言，其操作心智模型都是统一的。这种“UI 提供骨架，插件提供灵魂”的模式是 VS Code 生态稳固的基石。

---
