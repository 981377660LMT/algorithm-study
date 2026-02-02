# Wasm JSPI 的演进：新的 API 设计

# WebAssembly JSPI has a new API

这是对 JSPI (JavaScript Promise Integration) API 的一次重大更新。

### 1. 更加符合 Web 标准

新的 API 调整了 Promise 的交互方式，使其更符合当前的 Web 规范。主要改进在于如何处理导出函数（Exported Functions）的暂停与恢复逻辑。
V8 在内部优化了从 Wasm 栈切换到 JS 回调时的上下文保存效率，减少了创建临时 Promise 对象带来的 GC 开销。

### 2. 状态管理的细化

增强了对异步操作挂起状态的追踪。通过在 Wasm 实例级别维护更多的元数据，JSPI 现在能更稳健地处理多个并发的异步调用。

### 3. 一针见血的见解

Wasm 不再是异步编程的二等公民。这一 API 的更新标志着 JSPI 正在从一个“能用”的实验特性，进化为一个“好用”的工业标准。
