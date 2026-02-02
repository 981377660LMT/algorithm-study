# 78. V8 release v6.7 | V8 v6.7 发布：BigInt 与 Wasm 性能提升

V8 v6.7 的核心亮点是正式支持了 Arbitrary-precision BigInt 以及 WebAssembly 的启动加速。

## 1. BigInt：打破 JavaScript 的数值局限

JavaScript 传统的 `Number` 只能安全地表示至 $2^{53}-1$。

- **实现细节**：BigInt 在内存中表示为一组 32 位或 64 位的 "digits"。
- **内建操作**：V8 为 BigInt 实现了完整的算术运算。
- **注意点**：BigInt 与 Number 不能直接混合进行算术运算（如 `10n + 5` 会抛错），这是为了强迫开发者显式处理精度边界。

## 2. WebAssembly 启动优化：流式编译 (Streaming Compilation)

在 v6.7 之前，WebAssembly 需要完全下载后才能开始编译。

- **新技术**：引入 **Streaming Compilation**。V8 可以在 Wasm 字节码还在下载时，边下载边在后台线程进行 Tier-1（快速启动编译器，如 Liftoff/Baseline）编译。
- **效果**：极大减少了大型 Wasm 模块的首次运行时间，往往在下载完成时，编译也几乎完成了。

## 3. 内存减少：垃圾回收器的改进

v6.7 对回收大对象（Large Objects）进行了优化。

- **并行清理 (Parallel Sweeping)**：在标记阶段结束后，清理工作可以跨多个后台线程并行执行，缩短了垃圾回收的整体时长。

## 4. 其它特性

- **WeakRef (阶段性引入)**：开始在内部支持弱引用，这是管理大型对象图且不阻止内存回收的关键工具。
- **String.prototype.trimStart / trimEnd**：这两个标准方法取代了非标准的 `trimLeft/trimRight`。
