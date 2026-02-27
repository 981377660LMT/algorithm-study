MCP 服务器 工具数量 Token 消耗
GitHub 官方 93 个工具 ~55,000 tokens
Task Master 59 个工具 ~45,000 tokens
一个 Claude Code 用户报告：启用几个 MCP 后，上下文使用量达到 178k/200k（89%），其中 MCP 工具定义就占了 63.7k。还没开始干活，上下文已经快满了。
`问题的根源是：MCP 在启动时加载所有工具定义。不管你用不用，93 个 GitHub 工具的 schema 都要塞进上下文。`

---

- 2025 年 10 月，Anthropic 在 Claude Code 中引入 Skills。核心设计思路是：按需加载，而不是全量加载。

- Skill 的作者最了解任务特性，所以让 Skill 自己决定执行模式（?）。
