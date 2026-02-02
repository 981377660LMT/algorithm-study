# 进入实战：Wasm JSPI 的 Origin Trial

# WebAssembly JSPI is going to origin trial

### 1. 现实世界的测试

JSPI 正式进入 Origin Trial，意味着开发者可以在生产环境的浏览器中测试这个原本属于高深魔法的功能。
Origin Trial 带来的反馈将直接决定 V8 如何优化 **Stack Switching** 的底层实现。

### 2. 栈切换的性能压榨

为了支撑 Origin Trial，V8 团队对 Wasm 到 JS 的栈切换路径进行了深度扫描。优化了寄存器溢出（Register Spill）和恢复的过程，确保异步 Wasm 调用的延迟降到最低。

### 3. 一针见血的见解

这是一次“走出实验室”的关键飞跃。它证明了将复杂的 C++/Rust 传统异步代码无缝迁移到 Web 平台的商业前景。
