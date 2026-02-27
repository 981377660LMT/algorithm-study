多 Agent 如何协作？ DAG、Supervisor、还是 Handoff？
Token 预算如何控制？ 单次调用还是整个 workflow？
工具执行失败怎么办？ 重试、降级、还是人工介入？
如何保证安全？ 沙箱隔离、权限控制、审计日志

---

Orchestrator (Go) - 编排、预算、策略
Agent Core (Rust) - 执行、沙箱、限流
LLM Service (Python) - 推理、工具、向量

---

"The best way to learn is to build."

---

- Agent 是能自主完成任务的 AI 系统，核心在于"自己做决定"而不是"你说一句它动一句"。
- 你给它一个目标、一个工具箱和一套边界（预算 / 权限 / 审批 / 沙箱），它在循环里推进任务，直到完成或停下。
- Agent 的核心是自主决策循环

---

单 Agent 能力有限，多 Agent 协作成为主流。同时，企业开始关注：

- 成本控制（Token 预算）
- 安全性（沙箱执行）
- 可靠性（持久化、重试）
- 可观测性（监控、追踪）

这就是本书关注的重点：生产级 Agent 系统。
