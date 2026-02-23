# VS Code 博客深入分析补充篇

## 2023

### 三月（March）

- [Visual Studio Code 与 GitHub Copilot / Visual Studio Code and GitHub Copilot](https://code.visualstudio.com/blogs/2023/03/30/vscode-copilot)

## Visual Studio Code 与 GitHub Copilot（2023 年 3 月）

**核心概念：从"代码补全"到"聊天对话"，AI 在 IDE 中的范式转变**

- **历史背景**：
  - 2021 年，GitHub Copilot 首次发布，初期功能是"幽灵文本"（ghost text）补全——在你输入代码时，根据上下文生成下一行代码建议。
  - 使用 OpenAI 的 Codex 模型，首次大规模在生产中应用 LLM 做代码生成。
  - 体验是"Tab-Tab-Tab"流畅开发：写一个注释或函数名，Copilot 自动生成块代码。

- **ChatGPT 冲击（2022 年 11 月）**：
  - ChatGPT 爆火后，VS Code 团队加速了 Chat 功能的探索。
  - 内部 hackathon 涌现大量 AI 创意：改进重命名、基于示例的代码转换、自然语言生成 glob 模式/正则等。
  - 但逐渐意识到：**聊天对话才是关键**。

- **三层 Copilot 体验**：

  **1. 内联补全（Inline Completions）**
  - 边写代码边出现建议（幽灵文本）。
  - 最轻量级交互，保持"沉浸式"开发流。

  **2. 编辑器内聊天（In-Editor Chat / Inline Chat）**
  - 快捷键 Cmd+I (macOS) / Ctrl+I (Windows/Linux)，在编辑器弹出聊天框。
  - 示例：选中 `users` 数组，问"把 username 字段拆分成 firstName 和 lastName"→ Copilot 提议改动 → 点"Inline Diff"预览 → 接受。
  - **优势**：不离开代码，纯粹交互式编辑。

  **3. Chat 视图（Dedicated Chat View）**
  - 在侧栏打开完整 Chat 视图，支持多轮对话。
  - 支持斜杠命令（slash commands）和话题范围化，如 `/vscode 怎么关闭面包屑导航`。
  - 支持 Quick Chat（Shift+Alt+Cmd+L），快速弹窗提问。

- **为什么选择 Chat 而非简单网页聊天**：
  - **上下文感知**：VS Code 知道整个工作区，可以提供完整代码库上下文给 LLM，改善回答质量。
  - **例子**：用网页 ChatGPT 很难要求它"优化跨多个文件的代码"，但 VS Code Chat 可以做到。
  - **多步流程**：编程任务往往多步骤（类似博客教程），Chat 视图非常适合。
  - **调试助手**：问"怎么配置 `launch.json` 和 `tasks.json`"，Copilot 直接生成文件内容和步骤，vs 查文档快很多。
  - **人在中心**：LLM 不完美，Chat 的双向对话让用户决策——提出澄清问题、纠正错误、投票反馈。

- **如何用好 Copilot（文中建议）**：
  1. **你是 Pilot，Copilot 是副驾**——你决定采纳什么建议，融合什么代码。
  2. **用 Copilot 做繁琐任务**——生成测试、样本数据、代码框架等。
  3. **提供更多上下文**——不说"Node Express TypeScript"，而说"用 Express.js 和 TypeScript 搭建 Node.js 网站"，再迭代。
  4. **使用话题范围化**——`/vscode` 或 `/explain` 等快捷命令。
  5. **承认 Copilot 会犯错**——提问、投票、反馈来改进。
  6. **友好问候**——"Hello Copilot"会改善心情（闲趣）。

- **责任 AI**：
  - Microsoft 和 GitHub 遵循"负责 AI"原则，有审计机制。
  - 承认对 AI 快速发展的顾虑，尊重选择不用 Copilot 的开发者。

- **未来路线**（文中预告）：
  - 自然语言搜索、代码生成。
  - 自动生成 commit message 和 PR 描述。
  - 更聪明的重命名和重构。
  - 代码转换。

- **哲学**：Copilot 是"对话式编程"的开始，而非终点。Chat 的引入把 VS Code 从"编辑工具"升级为"交互式 AI 开发环境"。

---

## 2022

### 七月（July）

- [Visual Studio Code 服务器 / The Visual Studio Code Server](https://code.visualstudio.com/blogs/2022/07/07/vscode-server)

## Visual Studio Code 服务器（2022 年 7 月）

**核心概念：从桌面客户端到浏览器云端，VS Code 架构的云原生演进**

- **演进路线**：
  - **2019**：Remote Development 扩展 — 在本地 VS Code 中开发远程机器（WSL、Docker、SSH）。
  - **2020**：GitHub Codespaces — 在云上开发，可从本地 VS Code 或浏览器访问。
  - **2022**：VS Code Server — 把 VS Code 的后端服务解耦出来，任何地方部署，通过浏览器访问。

- **VS Code 的架构优势**：
  - VS Code 设计为**多进程应用**：
    - **前端**：你输入代码的地方（UI/编辑器）。
    - **后端**：扩展、终端、调试等服务运行的地方。
  - 这种分离使得"前端在浏览器，后端在远程"成为可能。

- **VS Code Server 是什么**：
  - 一个 CLI 工具 + 后端服务，可以在任何地方部署（本地 VM、云主机等）。
  - 通过浏览器访问，使用 vscode.dev（VS Code for Web）+ 安全隧道连接。
  - 无需手工配置 SSH 或 HTTPS（但也支持）。

- **快速开始（文中示例）**：
  1. 在远程机器（如 WSL）运行安装脚本 → 启动 `code-server` CLI。
  2. CLI 会显示设备码，用 GitHub 账号认证到隧道服务。
  3. 给远程机器取个名字（如"elegant-pitta"）。
  4. 获得一个 vscode.dev URL，用任何设备打开它。
  5. 连接成功后，远程文件系统出现在 Explorer，即可编码。

- **设计特点**：
  - **安全隧道**：GitHub 账号认证，无需裸露 IP。
  - **跨设备**：从任何设备（手机、平板、异地电脑）打开生成的 URL 即可编码。
  - **扩展支持**：后端运行扩展，前端渲染，完整功能支持。
  - **一键启动**：相比 SSH 配置复杂度低很多。

- **当时的局限**：
  - 私有预览，需要申请。
  - `code-server` CLI 与 `code` CLI 分离（后来会统一）。
  - 文档和视频还在建设中。

- **战略意义**：
  - **云原生开发**：VS Code 从"安装在本地的客户端"成长为"云原生开发工具"，支持浏览器、移动等多平台。
  - **敲定事实**：Codespaces 可以用，但 Server 提供了"自托管"选项 — 企业可在自己的基础设施上运行。
  - **未来铺垫**：为后来的"AI-native VS Code"奠定基础 — LLM 可以部署在云端，客户端轻量级。
  - **竞争力**：vs JetBrains Fleet 等云 IDE 竞品，VS Code 提前布局多元化接入方式。

---

## 总结与趋势

从 2022-2026 年的 VS Code 博客演进可看出：

| 年份 | 核心主题                   | 关键词                                                           |
| ---- | -------------------------- | ---------------------------------------------------------------- |
| 2022 | 云端开发基础设施           | VS Code Server、Codespaces、浏览器访问                           |
| 2023 | Chat + Copilot 的 IDE 集成 | In-Editor Chat、Chat View、LLM at scale                          |
| 2024 | AI 功能开放 + 多模型       | Copilot Edits、扩展 API、模型选择、Copilot for Azure             |
| 2025 | 多代理协作 + 企业化        | Agent Sessions、MCP、Subagents、Open Source、Private Marketplace |
| 2026 | 代理 UI 和可视化           | MCP Apps、视觉表达、代理编程                                     |

**三个深层转变**：

1. **从工具到平台**：VS Code 不再只是编辑器，而是"开发者操作系统"。
2. **从个人到团队**：从"我一个人用 Copilot"演进到"我和多个 Agent 协作"、"团队共享指令和市场"。
3. **从黑盒到透明**：AI 功能逐步开源、可审计、可定制（责任 AI）。
