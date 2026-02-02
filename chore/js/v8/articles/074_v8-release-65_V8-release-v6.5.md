# 74. V8 release v6.5 | V8 v6.5 发布：import.meta 与 Wasm 背景编译

V8 v6.5 继续深化了 WebAssembly 和现代模块化 JS 的集成。

## 1. import.meta

正式支持 `import.meta` 对象，允许模块获取其自身的元数据（如 URL）。这在构建平台无关的 JS 代码时非常有用。

## 2. WebAssembly 背景编译 (Background Compilation)

继之前的流式编译后，v6.5 进一步优化了多核利用率。

- **原理**：编译任务会被拆分为更小的单元，并分发到更多的后台工作线程中。
- **效果**：对于大型 Wasm 项目（如 AutoCAD Web），编译速度提升了数倍，主线程几乎不再被长耗时的编译任务阻塞。

## 3. TypedArray 性能提升 (CSA 重构)

此版本完成了对 `TypedArray` 关键 API（如 `subarray`, `slice`）的 **CSA 化**。

- **零拷贝优化**：在特定场景下，通过更精细的边界检查和内存管理，减少了不必要的内存拷贝。

## 4. 垃圾回收器的微调

此版本对新生代垃圾回收器（Scavenger）进行了并行化改进。现在，存活对象的存活判断和复制过程可以利用多个线程，极大地缩短了 minor GC 的时间。
