# 追求 VS Code 中的卓越智能

链接：https://code.visualstudio.com/blogs/2023/11/13/vscode-copilot-smarter

## 深入分析

### 1. 智能的根源：克服“上下文贫瘠”

AI 对代码的理解上限取决于它能看到的上下文。这篇文章标志着 **RAG (检索增强生成)** 在 IDE 中的深度应用——引进了 **@workspace 参与者**。它通过 GitHub 搜索索引与本地词法索引结合，将最相关的代码片段喂给大模型。

### 2. 交互的原语：Chat Participants 与 Slash Commands

VS Code 确立了一套标准化的 AI 扩展协议：

- **@ 参与者**：领域专家。如 `@workspace` 懂项目，`@vscode` 懂编辑器配置。
- **斜杠命令**：特定任务的快捷操作。如 `/fix` 修复报错，`/doc` 生成注释。
  这不仅统一了交互，也为第三方扩展（如 Docker, Azure）贡献自己的 AI 能力奠定了基石。

### 3. 多模态与情感化：语音输入与 Affirmation

Speech 扩展的加入让开发者可以用语音直接命令 Copilot。有趣的是，文章还展示了 DALL-E 参与者生成的“鼓励猫”，这体现了 VS Code 团队对“开发者体验 (DX)”的理解不仅在于生产力，也在于工作情绪的调节。
