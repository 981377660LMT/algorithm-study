---
title: 推介 GitHub Issues 集成
date: 2020/05/06
url: https://code.visualstudio.com/blogs/2020/05/06/github-issues-integration
---

## 深入分析

GitHub Issues 与 VS Code 的官方深度集成，标志着 IDE 正在向“研发全流程工作站”演进。

1.  **打破“上下文切换”魔咒**：
    - 在集成之前，开发者需要频繁切换浏览器与编辑器来查阅 issue 详情。
    - **核心特性**：直接在代码库中识别 `#123` 形式的提及，通过悬停预览（Hover）展示 Issue 描述、评论和指派人。这让代码及其设计背景在视觉上融为一体。
2.  **打通内外环操作**：
    - **从 TODO 到 Issue**：支持将代码中的 `// TODO` 注释一键转换为 GitHub Issue 并自动关联源码链接。
    - **工作流自动化**：在 Issues 视图中点击“工作并在新分支上打开”，会自动检出新分支、关联上下文，省去了大量手动 Git 操作。
3.  **内聚性扩展开发模式**：
    - VS Code 团队通过将此功能放在 GitHub Pull Requests & Issues 扩展中，而不是内置到内核，展示了“核心精简、通过强大的插件 API 扩展垂直能力”的设计初衷。

---
