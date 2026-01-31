---
title: vscode.dev
date: 2021/10/20
url: https://code.visualstudio.com/blogs/2021/10/20/vscode-dev
---

## 深入分析

vscode.dev 的发布实现了 VS Code “编辑器作为 Web 应用”的最终愿景。

1.  **浏览器中的桌面能力**：
    - **File System Access API**：核心突破在于利用现代浏览器的文件系统 API，允许用户直接在浏览器中打开和保存本地文件夹，无需将代码上传到云端。
    - **端到端隔离**：所有代码处理均在本地浏览器内完成，兼顾了安全性与 Web 的便携性。
2.  **技术演进：从服务端到客户端**：
    - 与 GitHub Codespaces 不同，vscode.dev **没有配套的计算实例 (Server-less)**。
    - **Web Worker 扩展**：VS Code 通过 Web Worker 运行时支持了大量的 UI 级和语法级扩展。这意味着传统的 Node.js 扩展需要重构为 Web 兼容版本。
3.  **连接万物**：
    - 通过对 GitHub 和 Azure Repos 的原生集成，vscode.dev 变成了“代码审查”和“快速修改”的最佳工具。
    - **WASM 应用场景**：为后续在浏览器中运行 Python/C++ 调试器（通过 WASM）奠定了基础。

---
