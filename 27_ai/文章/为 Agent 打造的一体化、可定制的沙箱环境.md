https://github.com/agent-infra/sandbox
https://sandbox.agent-infra.com/zh/guide/start/introduction
https://sandbox.agent-infra.com/zh/blog/announcing-0

AI Agent 在执行复杂任务时，需要在浏览器、代码执行、文件系统之间切换。传统多沙箱方案面临环境割裂、数据搬运、鉴权复杂等问题。
AIO Sandbox 通过一个 Docker 镜像整合 Agent Env 环境三大件（浏览器、命令行、代码执行），提供统一文件系统与鉴权，并支持镜像定制，提升了 Agent 任务执行与交付效率。

AIO Sandbox 在一个沙盒内集成浏览器、代码执行、终端、可视化接管、正反向代理、MCP、鉴权等基础功能，可根据需求进行沙盒环境定制，让不同的 Agent “在一个环境容器内中更高效地完成任务”。

---

SandBox: 受控、隔离的执行环境。用于运行浏览器、代码或命令行，控制资源与权限，降低对宿主系统的影响与风险。
