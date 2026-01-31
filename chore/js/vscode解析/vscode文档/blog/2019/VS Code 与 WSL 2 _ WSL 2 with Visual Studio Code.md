---
title: VS Code 与 WSL 2
date: 2019/09/03
url: https://code.visualstudio.com/blogs/2019/09/03/wsl2
---

## 摘要

2019 年 9 月，WSL 2 的正式发布改变了 Windows 上的 Linux 开发体验。相比于 WSL 1 使用系统调用转换（Translation），WSL 2 运行了一个包含完整内核的超轻量自研虚拟机，使得文件 IO 性能提升了 20 倍。本文解释了 VS Code 如何利用 Remote - WSL 扩展在 Windows 上渲染 UI，而将后端服务（VS Code Server）直接运行在这种高性能的 Linux 虚拟环境中，从而消除了跨系统的权限和库兼容性问题。

## 一针见血的分析

WSL 2 的发布配合 VS Code 远程开发扩展，彻底终结了“在 Windows 上开发 Linux 项目”的割裂历史。其工程美学在于**“高性能的透明性”**：WSL 2 解决了 WSL 1 无法支持复杂系统调用（如 Go 调试器、Docker 守护进程）的架构痛点，而 VS Code 巧妙地利用 Linux 挂载 Windows 文件系统的能力，实现了“代码在 Linux，编辑在 Windows”的极致体验。这是一场关于 **“内核级集成”** 的胜利——通过将一个完整的 Linux 环境作为 IDE 的运行时后台，VS Code 成功将 Windows 桌面变成了全球最强大的 Linux 开发工作站。

## 深入分析

WSL 2 的发布彻底改变了 Windows 上的开发体验，而 VS Code 是这一变革的最佳搭档。

1.  **内核重构带来的性能飞跃**：
    - **WSL 1**：通过翻译层将 Linux 系统调用转换为 Windows 调用，兼容性受限且文本 IO 性能差。
    - **WSL 2**：引入了真正的轻量级 Linux 内核（通过 Hyper-V 运行），实现了高达 **20 倍** 的 IO 性能提升。这意味着 `git clone`、`npm install` 等操作的速度终于能与原生 Linux 媲美。
2.  **VS Code 远程架构的“杀手级应用”**：
    - VS Code 通过 Remote - WSL 扩展，将 **VS Code Server** 直接运行在 WSL 2 内部。
    - **零摩擦体验**：开发者在 Windows 下操作 UI，但所有代码执行、调试、IntelliSense 都在 Linux 环境中完成。它解决了 Windows 开发中长期存在的路径转换、换行符（LF vs CRLF）以及原生模块编译失败等顽疾。
3.  **开发范式转移**：
    - 此后，Windows 不再是 Linux 的克隆或模拟器，而是成为了运行 Linux 环境的最佳“宿主机”之一。它让开发者能同时享受 Windows 优秀的桌面体验和 Linux 强大的生产力工具链。

---
