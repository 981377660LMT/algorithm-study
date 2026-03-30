MCP 是工具的 USB 接口——标准化协议，任何 Agent 都能复用社区写好的 Server
2025 年 MCP 走向事实标准——9700 万月下载、1 万+活跃 Server，并加入 Linux Foundation 旗下 AAIF 做中立治理
Shannon 用的是简化版 HTTP 调用——够用但功能不完整，适合快速集成
安全问题很重要——Prompt Injection、权限组合攻击、伪装 Server 都是真实风险
生产必备配置：域名白名单、响应大小限制、超时控制、熔断器

- Tools 是写操作，Resources 是读操作
