# 使用 VS Code 检查容器

链接：https://code.visualstudio.com/blogs/2019/10/31/inspecting-containers

## 摘要

相比于使用命令行 `docker exec -it` 进行简陋的交互，VS Code 的 Dev Containers 扩展提供了一种更强大的方式：直接将 VS Code 挂载到运行中的容器内。本文详细介绍了“Attach to Container”的工作流：VS Code 会动态地在容器内安装 VS Code Server，并以此为跳板实现源码浏览、断点调试和扩展运行。文章还解释了配置文件（如 `express-server.json`）如何持久化存储连接信息和特定容器所需的扩展列表。

## 一针见血的分析

“挂载到正在运行的容器（Attach to Container）”是 VS Code 实现 **“无感知全功能调试”** 的绝佳案例。其核心价值在于消除了传统远程调试中繁琐的“文件映射”和“源码路径转换”难题——因为 VS Code Server 就运行在容器内部的真实文件系统中。这种模式最值得称道的工程细节是**扩展的“本地/远程”自动分类分发**：主题等 UI 扩展留在本地，而 GitLens、语言服务器等核心工具则透明地部署在容器工作负载中。这使得“容器”不再是一个黑盒的运行环境，而变成了一个具备“IDE 意识”的主动开发载体，极大地加速了大规模微服务架构下的现场排障效率。
