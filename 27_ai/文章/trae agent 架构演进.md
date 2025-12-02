# trae agent 架构演进

本文讨论了 Trae Agent 的架构演进、Trae Builder 模式介绍及常用场景、自定义 Agent 相关内容，还提及了未来演进路线和 Multi Agent 的优势与挑战。关键要点包括：

1. Trae Builder 模式：可帮助用户从 0 到 1 开发项目或融入已有项目，AI 助手会根据需求调用不同工具，使回答更精确有效。
2. Trae Builder 常用场景：包括 0 到 1 场景，如创建一个 web 应用；Bug Fix 场景，修复程序中不符合设计的缺陷；还支持自定义 Agent，用户可定制提示词和工具。
3. Trae Agent v1.x 架构：用户通过自然语言对话交互，Agent 按【思考 > 规划 > 执行 > 观察】循环工作，核心模块包括上下文管理、任务规划和工具调用等，还提到了代码知识图谱和上下文裁剪。
4. Trae Agent v2.x 架构演进：流程上用模型驱动替代固定工作流，取消固定的 Proposal 流程；Context 优化方面，减少预设召回流程，自动压缩历史对话，去除默认全局 Context 召回。
5. Trae Agent v3.x 未来演进路线：期望构建能处理复杂场景的通用 Agent，采用 Graph 方式支持灵活流程配置，Agent 维护独立 State。
6. Multi Agent 优势：包括 Prompt 隔离、Tools 隔离和 History 隔离，可避免信息分散、降低决策难度、减少信息噪音。
7. 技术挑战：实现 Multi Agent 高效协作存在调度协作、异构 Agent 联动和跨 Agent 信息传递等技术挑战。

---

自定义 Agent：支持用户定制`提示词`和`工具`，以完成搭建符合特定领域需求的 Agent。
Agent：眼（上下文 Context）、脑(Memory)、手(Tools(MCP))

---
