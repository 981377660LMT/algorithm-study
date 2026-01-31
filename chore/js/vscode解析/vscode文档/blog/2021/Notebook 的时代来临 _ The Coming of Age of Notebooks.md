# Notebook 的时代来临

链接：https://code.visualstudio.com/blogs/2021/08/05/notebooks

## 摘要

本文宣布了 VS Code Notebooks 的“成熟期”，标志着 Notebook 支持从插件形式演变为原生内核功能。文章讨论了将 Notebook 体验与 VS Code 核心生态系统（如快捷键、扩展支持、差异对比）深度融合的努力。新的 Notebook 架构不仅支持 Jupyter，还通过通用 API 支持了如 REST Book、GitHub Issues Notebooks 等自定义领域的交互式文档。

## 一针见血的分析

Notebooks 的原生化是 VS Code 历史上最重要的 UI 范式迁移之一。其工程核心在于**对“文档模型”的重新定义**：将传统的线性文本编辑器抽象为由多个单元（Cells）组成的有序列表。这种架构允许 VS Code 将强大的代码编辑能力（IntelliSense, Diagnostics）注入到每一个 Cell 中，而不仅仅是作为一个大型文本缓冲区。此外，通过引入“控制器（Controller）”和“渲染器（Renderer）”的解耦模型，VS Code 解决了交互式输出的安全性和性能问题（利用 Webview 隔离）。这是一场关于“可发现性”的胜利——让数据科学家和开发者在同一个工具中，以最高标准享受 IDE 对交互式计算的支持。
