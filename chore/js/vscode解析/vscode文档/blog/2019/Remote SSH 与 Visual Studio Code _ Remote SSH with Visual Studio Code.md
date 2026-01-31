# Remote SSH 与 Visual Studio Code

链接：https://code.visualstudio.com/blogs/2019/07/25/remote-ssh

## 摘要

VS Code 远程开发扩展允许开发者通过 SSH 连接到 Linux 虚拟机、容器或 WSL 并在其中进行开发。本文重点介绍了 Remote - SSH 扩展，解释了它如何打破本地硬件限制，让开发者能够利用远程高性能服务器的存储和算力，同时保持与本地一致的 IDE 体验。文章详细说明了 VS Code Server 的初始化过程、SSH 密钥配置、端口转发（让本地浏览器访问远程 Web 应用）以及远程扩展的管理。

## 一针见血的分析

Remote - SSH 的成功标志着 VS Code 从一个“本地编辑器”转型为“解耦的客户端-服务器架构工具”。其最核心的工程贡献在于成功的**关注点分离（Separation of Concerns）**：将 UI 渲染、主题、快捷键等前端逻辑保留在本地，而将文件系统操作、终端运行、语言服务分析（LSP）等重型逻辑下沉到远程生成的 VS Code Server 中。这种“极低延迟”的体验得益于其智能的扩展分类模型：UI 扩展随客户端安装，而分析类扩展随工作负载安装。此外，通过 SSH 隧道实现的双向端口转发（Port Forwarding）极大地降低了 Web 开发的调试门槛，使得“本地开发、云端运行”不再是一种妥协，而是一种增强。
