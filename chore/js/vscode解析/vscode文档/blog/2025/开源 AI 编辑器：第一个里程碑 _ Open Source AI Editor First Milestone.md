# 开源 AI 编辑器：第一个里程碑

链接：https://code.visualstudio.com/blogs/2025/06/30/openSourceAIEditorFirstMilestone

## 深入分析

### 1. 行胜于言：Copilot Chat 的全栈代码公开

这个里程碑标志着微软兑现了 2025 年 5 月的开源承诺。GitHub Copilot Chat 扩展正式以 MIT 协议开源，包括了 Agent 模式的复杂状态管理逻辑，以及最为关键的：由微软团队打磨数年的**系统提示词（System Prompts）**。

### 2. 消除 AI 黑盒疑虑

通过公开 Telemetry、Context 获取逻辑以及 Prompt 构建过程，VS Code 为开源社区树立了“透明 AI”的样板。开发者可以清晰查看到底有哪些数据被发往模型，这对医疗、金融等强合规行业在评估 AI 工具安全性方面具有里程碑意义。

### 3. 下一步计划：内核合龙

文章明确指出，此番开源只是第一步，未来的目标是将这些从扩展中解耦出的通用 AI 能力深度重构成 VS Code 的核心组件，从而让“AI 原生”成为编辑器的底层属性而非外部插件。
