# 开源 AI 编辑器：第二个里程碑

链接：https://code.visualstudio.com/blogs/2025/11/04/openSourceAIEditorSecondMilestone

## 深入分析

### 1. 整合与提效：Inline Suggestions 的彻底开源

过去几年来，VS Code 中的 GitHub Copilot 一直分为两个扩展：GitHub Copilot 扩展（用于幽灵文本建议）和 GitHub Copilot Chat 扩展（用于聊天和下一编辑建议）。我们正在努力在单个 VS Code 扩展中提供所有 Copilot 功能：Copilot Chat。
为了实现这一目标，我们现在正在测试禁用 Copilot 扩展，并从 Copilot Chat 提供所有内联建议。

第二个里程碑的重点在于将 AI 产生最高频的环节——**行内实时代码补全（Inline Suggestions）**——推向开源。通过将 `GitHub Copilot` 扩展的功能合并进 `Copilot Chat`，用户体验变得更加简洁，同时也降低了网络请求的冗余延迟。

### 2. 深入工程实现：揭秘补全链路

文章对补全逻辑进行了教科书级的拆解：

- **Typing-as-suggested**：精准判断用户是否正在按原建议输入，避免不必要的 LLM 请求。
- **Caching 与请求复用**：在流式输出过程中，利用上一个字符生成的中间状态，实现毫秒级的响应回滚。
- **Block Trimmer**：智能裁剪模型输出的多行内容，确保生成的代码符合当前上下文的语法开闭闭环。

### 3. 命名规范的统一：迈向 2026

官方正式确立了 **"Inline Suggestions"** 这一统一称呼，涵盖了传统的 Ghost Text 和最新的 Next Edit Suggestions。伴随着旧版扩展的废弃计划，VS Code 正朝着一个更整洁、更统一的 AI 开发环境全速前进。
