Hooks 解决三个问题：看（可观测）、管（可控制）、扩（可扩展）
事件分级很重要——不是所有事件都要持久化，LLM_PARTIAL 不存，WORKFLOW_COMPLETED 必存
暂停/恢复用 Temporal Signal——不是轮询，是真正的阻塞等待
人工审批是安全护栏——基于策略触发，支持超时自动拒绝
Hook 要非阻塞——队列满了就丢，不能拖慢主流程

---

Claude Code 有一套简单但实用的 Hooks 机制，值得参考。
它把 Hooks 定义在 .claude/hooks/ 目录下，用独立脚本实现：
调用方式很简单：通过 stdin 传入事件数据，脚本处理后返回。
