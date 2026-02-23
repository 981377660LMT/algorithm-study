# VS Code 迁移到进程沙箱

链接：https://code.visualstudio.com/blogs/2022/11/28/vscode-sandbox

## 一针见血的分析

1. **从桌面应用向 Web 架构的终极演进**：沙箱化的核心是「去 Node.js 化」。通过禁止 Renderer 进程直接调用 Node 接口，VS Code 强制自己实现了一套基于 Web 标准（Uint8Array, MessagePort, Custom Protocol）的架构。这不仅提升了安全性，也让 Electron 版与 Web 版的代码路径彻底统一。
2. **IPC 的性能博弈**：当所有系统操作都必须跨进程时，IPC 成为瓶颈。VS Code 引入了 `MessagePort` 实现进程间的直连（Renderer <-> Extension Host），绕过了繁忙的主进程（Main Process），确保了在大负荷下的 UI 响应速度。
3. **「边飞边修飞机」的工程奇迹**：这是一个跨度三年的重构工程。通过每月迭代逐步将渲染进程中的 Node 模块外迁到 Utility Process 或 Shared Process，最终在不中断用户使用的情况下完成了底层架构的跃迁。

## 摘要

本博文详细介绍了 VS Code 历时三年完成的 Electron 进程沙箱化（Sandbox）迁移过程，旨在提升安全性和架构健壮性。

- **核心目标**：禁止在渲染进程（Renderer Process）中直接使用 Node.js APIs，对齐 Web 安全模型。
- **关键技术**：
  - **Preload Scripts**：作为受信任环境与沙箱环境的桥梁。
  - **MessagePort**：构建高效的直连 IPC 通道，减少主进程负担。
  - **vscode-file 协议**：取代 `file://` 协议，增强跨域安全控制。
  - **UtilityProcess**：推出新的 Electron API 用于托管扩展宿主等复杂进程。
- **架构收益**：大幅降低了渲染进程被攻击后执行任意系统代码的风险，并显著提升了窗口重载与切换工作区的速度。

## 深入分析

### 1. 架构大手术：切除渲染进程的 Node.js

Electron 早期允许渲染进程直接使用 Node.js，这虽方便却隐患巨大。VS Code 团队耗时三年，将所有 Node 依赖从 UI 线程中剥离，转而使用 **Preload Scripts** 和 **Context Bridge**。这不仅是为了安全，更是为了使桌面版和 Web 版代码架构完全对齐。

### 2. UtilityProcess API 与进程解耦

为了安置被剥离的“扩展宿主”，VS Code 向 Electron 贡献了 **UtilityProcess API**。现在的架构中，插件运行在独立的单进程里，通过 **MessagePort** 与渲染进程直接通信（不经过主进程中转），从而在加固安全边界的同时，保持了极高的通信效率。

### 3. V8 代码缓存优化

为了不让沙箱化导致启动变慢，团队优化了 Chromium 的代码缓存策略。通过 `bypassHeatCheck` 选项，强制 V8 在启动时记录热点代码，避免了重复解析 11MB 以上的压缩后的 workbench 代码，确保了冷启动性能。
