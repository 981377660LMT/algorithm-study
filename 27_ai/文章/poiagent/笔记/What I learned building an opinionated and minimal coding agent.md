# [What I learned building an opinionated and minimal coding agent](https://mariozechner.at/posts/2025-11-30-pi-coding-agent/)

[pi-mono](https://github.com/badlogic/pi-mono) 作者

- 简单可预测工具
- OpenClaw 选用的底层框架 pi-mono，在我看来就是这种工程哲学最极端的实践，核心只有四个工具，系统提示词不到一千个 token，完全不依赖 LangChain / LangGraph
  https://zhuanlan.zhihu.com/p/2009031121334207641
- [为什么我从 Claude Code 切换到 Pi（OpenClaw 背后的代理](https://www.reddit.com/r/ClaudeCode/comments/1r11egp/why_i_switched_from_claude_code_to_pi_the_agent/?tl=zh-hans)）

以下是对文章《What I learned building an opinionated and minimal coding agent（我是如何构建一个固执且极简的编程 Agent 的）》的精简讲解：

### 核心观点

作者（Mario Zechner）对当前市场上越来越复杂、黑盒化（如 Claude Code）的 AI 编程助手感到不满，于是自己开发了一个极其精简的开源终端代码助手 —— **pi**。结果证明，**仅靠 4 个核心工具和不到 1000 Token 的系统提示词，就能达到甚至超越复杂 Agent 的水平。**

### 痛点：现有工具的通病

1. **过度复杂与黑盒化**：像 Claude Code 增加了大量用户不需要的复杂功能（80% 的冗余），且在后台偷偷塞入隐藏的上下文，导致模型行为不可预测。
2. **缺乏可观测性**：遇到子智能体 (Sub-agents) 运行失败时，用户无法查看其内部过程。
3. **MCP（模型上下文协议）的代价**：导入 MCP 工具会一次性占用巨大的 Token（如 Playwright 占用 13.7k Token），严重挤压了真正有用的上下文空间。
4. **底层构建麻烦**：目前主流的 AI SDK（如 Vercel AI SDK）与自托管或各家不同模型的兼容性很差（例如跟踪 Token、跨提供商的思考链衔接等问题）。

### 解决方案：打造极简 Agent "pi"

作者从零构建了三层架构：统一的 AI API 层 (`pi-ai`)、极简终端 UI 框架 (`pi-tui`)、以及最终的 Agent CLI (`pi-coding-agent`)。其核心哲学是**“如无必要，坚决不建”**。

#### pi 的几个反常识（Opinionated）设计：

1. **极致精简的 Prompt**：不需要所谓的长篇大论，仅用几行字清楚说明角色和工具。由于现在的前沿模型都经过了强化学习 (RL) 调优，它们天生就知道如何扮演“程序员”。
2. **只有 4 个工具 (The "Core 4")**：
   - **`read`**: 读文件。
   - **`write`**: 写新文件。
   - **`edit`**: 精准编辑文件。
   - **`bash`**: 执行任意终端命令。
3. **YOLO 模式（默认裸奔）**：不搞所谓的“安全护栏”或“权限确认”。既然是替人写代码、跑代码的工具，安全隔离就是个伪命题，直接给予最高系统和网络权限即可。
4. **没有任何花哨功能**：
   - **无内置待办事项 (To-dos) / 计划模式**：需要待办或计划？让 Agent 自己写个 `TODO.md` 或 `PLAN.md` 文件在硬盘上，用户清晰可见还能协同编辑。
   - **不支持 MCP**：抛弃复杂的 MCP 服务器，改用最原始的模式：把工具写成简单的终端脚本 (CLI Tools)，让 Agent 用 `bash` 配合跑。
   - **不支持后台 Bash / 子智能体 (Sub-agents)**：需要后台运行任务？直接用系统的 `tmux`（终端复用器），这比让 Agent 自己管理后台进程可靠得多，且完全可见。

### 总结

作者通过跑 Terminal-Bench 2.0 评准测试证明：这个只有极简系统提示词和 4 个基础工具组合的 `pi`，在能力上不输于（甚至超过了）市场上主流的那些臃肿、复杂的编程 Agent。**放弃替大模型做复杂的决策流（如规划、开子进程），直接提供最通用的底层能力（读写文件 + 调 Bash），让大模型自己想办法解决问题，才是更好的架构，是未来最高效的发展方向。**
