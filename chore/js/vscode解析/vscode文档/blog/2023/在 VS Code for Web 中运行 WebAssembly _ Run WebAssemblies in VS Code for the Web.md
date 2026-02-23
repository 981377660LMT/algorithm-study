# 在 VS Code for Web 中运行 WebAssembly

链接：https://code.visualstudio.com/blogs/2023/06/05/vscode-wasm-wasi

## 深入分析

### 1. 攻克 Web 版 IDE 的最后堡垒：执行环境

VS Code for the Web (vscode.dev) 长期面临无法运行非 JS 代码（如 C++/Python/Rust）的挑战。团队通过引入 **WASM + WASI (WebAssembly System Interface)**，实现了在浏览器内直接运行经编译后的 Python 解释器。

### 2. 核心技术：SharedArrayBuffer 与 Atomics

这是一个高难度的架构设计：WASM 执行通常是**同步**的，而浏览器的文件系统 API 是**异步**的（无法直接在同步流中等待）。VS Code 团队利用 Worker 线程，配合 `SharedArrayBuffer` 和 `Atomics`，实现了一套“挂起 Worker -> 宿主异步处理 -> 原子化通知唤醒”的桥接机制。

### 3. 生态意义：跨平台插件系统的闭环

通过 `@vscode/wasm-wasi` 库，开发者可以使用 Rust 重写传统的 Unix coreutils（如 `ls`, `cat`），并将其提供给 Web 版终端。这让 VS Code 在浏览器中不再只是一个“UI 壳子”，而是一个真正具备计算和工具能力的轻量化 OS。
