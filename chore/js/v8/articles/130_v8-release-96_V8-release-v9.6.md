# V8 v9.6 发布：WebAssembly 引用类型等

# V8 release v9.6: WebAssembly Reference Types and more

V8 v9.6 在 WebAssembly (Wasm) 与 JavaScript 的互操作性上迈出了重要一步。

### 1. 技术核心：WebAssembly Reference Types (引用类型)

本版本最重要的技术特性是正式支持 **WebAssembly Reference Types**。

- **深度解析**：引入了 `externref`（以前称为 `anyref`）类型。在此之前，Wasm 只能存储数值（i32, f64 等）。如果需要持有 JS 对象引用，必须在 JS 端维护一个对象表，并将表的索引传给 Wasm。
- **直接持有**：现在，Wasm 模块可以直接、安全地持有 JS 对象的非不透明引用，并将其存储在局部变量、全局变量或 `WebAssembly.Table` 中。

### 2. 技术洞察：打通 Wasm 与 GC 的连接

这项改进彻底打通了 Wasm 与 V8 垃圾回收器 (GC) 的连接。

- **GC 穿透**：GC 现在可以穿透 Wasm 的栈和表格，准确追踪 JS 对象的生命周期。
- **降本增效**：这为未来的 Wasm GC 提案奠定了底层基础，并显著降低了 JS 框架与 Wasm 模块频繁交互时的编组（Marshaling）开销。

### 3. 一针见血的见解

通过 `externref`，Wasm 终于不再是一个与 JS 对象世界完全隔绝的“数值孤岛”。它标志着 Wasm 正在从一个单纯的计算插件演变为能深度参与 JS 复杂对象模型的一等公民，极大地拓宽了高性能 Web 应用的边界。
