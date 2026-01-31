# 在 Visual Studio Code 中的 GitHub Pull Requests

链接：https://code.visualstudio.com/blogs/2018/09/10/introducing-github-pullrequests

## 一针见血的分析

这篇文章宣告了 VS Code 试图统治开发者“完整生命周期”的野心：

1.  **打破编辑器与 Web 的边界**：长期以来，PR 评审被视为 Web 端的任务。VS Code 改变了这一点，将评论、状态检查、甚至分支检出直接带入编辑器。这不仅提高了效率，更重要的是提供了 **Contextual Depth**（上下文深度）——你可以在全功能的语言服务（跳转定义、查找引用）支持下进行代码评审。
2.  **Commenting API 的标准化**：为了支撑这个功能，VS Code 引入了通用的 **Commenting API**。这意味着除了 GitHub，GitLab、Bitbucket 等平台也可以通过同样的 API 实现一流的集成体验。这再次体现了 VS Code “协议先行”的设计哲学。
3.  **Checkout and Run 流程**：文章强调了在本地检出 PR 并直接运行的能力。这在繁琐的 Web UI 评审中是极难实现的。VS Code 通过深度集成 Git 和终端，让“边运行边评审”成为可能，极大地提高了大型开源项目的代码质量。
4.  **微软与 GitHub 的协同效应**：这是在微软宣布收购 GitHub 后的早期重要合作成果。它展示了两个生态系统融合后，如何通过工具链路的打通来重塑整个软件工程流程。

## 摘要

2018年9月，微软发布了 GitHub Pull Requests 扩展的公开预览版。该扩展将 PR 评审流程直接嵌入 VS Code，支持在编辑器内直接发表评论、协同代码。它引入了底层的 Commenting API，提供了一种全新的“原地评审”体验，让开发者无需离开代码环境即可完成高质量的协作。
