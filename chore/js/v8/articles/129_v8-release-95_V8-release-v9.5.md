# V8 v9.5 发布：国际化增强及 WebAssembly 异常处理

# V8 release v9.5: Intl.DisplayNames v2 and WebAssembly Exception Handling

V8 v9.5 在性能表现和工程化解耦方面都有显著提升。

### 1. 核心特性：WebAssembly 异常处理 (Wasm Exception Handling)

该版本在 WebAssembly 生态中迈出了关键一步：原生异常处理支持。

- **深度解析**：传统的 Wasm 异常是通过 JS 模拟的，涉及昂贵的跨语言切换和堆栈展开。新方案引入了原生的 `try`、`catch` 和 `throw` 指令。
- **性能飞跃**：在测试中，基于原生 EH 的代码体积比 JS 模拟方案减小了约 **30%**，执行性能在异常密集型任务中提升了约 **30%**。

### 2. 工程化革新：API 拆分与解耦

V8 的 C++ 头文件进行了“大手术”，将庞大的 `v8.h` 拆分为多个独立的子头文件（如 `v8-isolate.h`, `v8-local.h`）。

- **编译效率**：这种物理解耦显著提升了嵌入 V8 的大型项目（如 Node.js 或 Chromium）的编译速度，减少了因修改代码导致的冗余编译时间。

### 3. 一针见血的见解

原生异常处理的引入，补齐了 Wasm 作为高效编译目标的最后一块短板。它让那些依赖异常机制的语言（如 C++、C#）在 Web 环境下不再有沉重的运行时包袱，真正实现了接近原生的性能。
