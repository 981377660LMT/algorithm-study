- Claude Code Agent Teams 是 Anthropic 最新推出的一项实验性功能，允许在单个 Claude Code 会话中编排多个并行的 Agent 实例。它通过 “Team Lead + Teammates” 的架构，利用共享任务列表和消息系统来协调复杂任务。
- 底层依赖文件系统，包括 Team Lead、Teammates、Task List 和 Mailbox 四个组件，有任务依赖与文件锁机制。
  | 组件 | 职责 |
  |------|------|
  | Team Lead | 主 Claude Code session,创建团队、拆分任务、分配工作、汇总结果 |
  | Teammates | 独立的 Claude Code 实例,各自有自己的上下文窗口,独立执行任务 |
  | Task List | 共享的任务列表,所有队员都能看到,可以认领和完成任务 |
  | Mailbox | 消息系统,Agent 之间的通信通道 |

  和 subagent 最本质的区别是：`teammate 之间可以直接通信`。subagent 只能把结果报告给主 Agent，是单向的上报关系；而 Agent Teams 里的队员可以互相发消息、互相质疑对方的结论，更像真实的团队协作。

- 最适合研究、审查、新功能开发等独立性较强的任务；不适合强依赖顺序的任务或频繁修改同一文件的场景。

- Agent Teams 不是简单的 “多窗口” 或 “后台线程”，它是一个完整的编排层。在这个架构中，启动会话的 Claude Code 实例自动成为 Team Lead（队长），负责创建团队、生成任务并指派工作；而新生成的实例称为 Teammates（队员），它们在独立的上下文中运行，但共享同一个项目配置。
  它们通过两个核心抽象进行协作：
  - Shared Task List（共享任务列表）：协调工作的核心，支持依赖管理（pending/blocked）。Lead 可以指派任务，Teammates 也可以 “自领”（Self-claim）。
  - Mailbox（信箱机制）：基于消息的异步通信系统，允许 Teammates 之间直接对话，或向 Lead 汇报，无需 Lead 轮询。
