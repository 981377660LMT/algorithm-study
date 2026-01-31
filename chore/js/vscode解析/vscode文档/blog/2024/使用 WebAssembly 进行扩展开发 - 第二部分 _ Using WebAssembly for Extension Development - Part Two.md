# 使用 WebAssembly 进行扩展开发 - 第二部分

链接：https://code.visualstudio.com/blogs/2024/06/07/wasm-part2

## 深入分析

### 1. 异步化的平衡：Worker 隔离逻辑

第二部分重点解决了扩展开发的体验瓶颈：**UI 线程阻塞**。通过 meta 模型的自动生成的胶水代码，开发者可以轻松地将耗时的 WASM 逻辑（如大规模代码搜索）推送到后台 Worker 执行。VS Code 的 component model 自动处理了跨线程的消息序列化和同步状态同步。

### 2. Web 版 IDE 的救星：WASM Language Server

这是一个里程碑式的演示：在 WASM 中运行完整的 LSP (Language Server)。利用 `@vscode/wasm-wasi-lsp` 桥接层，Rust 写的分析引擎可以直接通过文件系统 API 访问工作区。这意味着即便在浏览器版本中，VS Code 也能通过 WASM 获得与原生桌面端一致的辅助能力（如 Goto Definition）。

### 3. 未来挑战：全量 Async 支持

文章诚实地指出了当前的局限性——WASM 规范目前对原生异步支持尚不完美。VS Code 团队正密切关注 WASI 0.3 的进展，这反映了他们不仅是工具使用者，更是 Web 基础技术标准的深度推动者。
