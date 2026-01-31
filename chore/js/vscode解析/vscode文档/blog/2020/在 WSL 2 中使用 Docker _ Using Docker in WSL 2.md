---
title: 在 WSL 2 中使用 Docker
date: 2020/03/02
url: https://code.visualstudio.com/blogs/2020/03/02/docker-in-wsl2
---

## 深入分析

Docker Desktop 的 WSL 2 后端是云原生开发在 Windows 上的终极解决方案。

1.  **架构演进：告别 Hyper-V 重型 VM**：
    - **旧架构**：Docker 在 Windows 上运行在一个厚重的 Hyper-V 虚拟机里，启动慢、资源消耗大且与宿主机文件系统交换缓慢。
    - **新架构**：利用 WSL 2 共享内核，Docker 守护进程（DockerD）直接运行在 WSL 2 轻量级虚拟机中。
2.  **性能优势**：
    - **文件系统同步**：在 WSL 2 环境中运行 `docker build` 或持久化存储时，文件操作速度显著提升（受益于 WSL 2 的原生 Linux 性能）。
    - **资源动态分配**：不再需要预分配固定比例的 CPU/内存，WSL 2 的动态资源调配让系统整体运行更流畅。
3.  **VS Code 的协同作用**：
    - 通过 **Docker 扩展** 与 **WSL 扩展** 的组合，开发者可以在 Windows UI 下享受 Linux 容器化的完整闭环。这为 **Dev Containers**（开发容器）在 Windows 下的大流行铺平了道路，实现了开发环境的“原子化”管理。

---
