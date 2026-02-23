# C++ 的高性能垃圾回收 (CPPGC/Oilpan)

## High-performance garbage collection for C++

- **Original Link**: [https://v8.dev/blog/cppgc](https://v8.dev/blog/cppgc)
- **Publication Date**: 2020-05-26
- **Summary**: V8 引入了 `cppgc`（原 Blink 中的 Oilpan 引擎），这是一个专为 C++ 优化的垃圾回收器。它通过统一 JS 和 C++ 的堆管理（Unified Heap），解决了跨语言引用导致的内存泄漏和零散化问题，并带来了增量标记和并发回收等高性能特性。

---

### 1. 背景：为什么 C++ 需要 GC？

在高性能渲染引擎（如 Chrome 的 Blink）和 JS 运行时中，JS 对象经常与 C++ 对象（如 DOM 节点）相互引用。

- **内存管理难题**：
  - C++ 通常使用引用计数（Reference Counting, `std::shared_ptr` 等），而 JS 使用垃圾回收（GC）。
  - **循环引用**：如果 JS 对象引用 C++ 对象，而该 C++ 对象又反过来引用回 JS，传统的引用计数无法打破这种环，导致内存泄漏。
- **跨边界追踪的开销**：在旧架构中，V8 必须保守地假设所有被 JS 触达的 C++ 对象都是活跃的。

### 2. 核心机制：CPPGC (Oilpan)

`cppgc` 是从 Blink 迁入 V8 核心库的一个组件，旨在提供：

#### A. 统一堆 (Unified Heap)

`cppgc` 让 V8 的垃圾回收器能够直接“看穿” C++ 对象堆。当 V8 进行主回收时，它可以追踪到 C++ 堆中的对象，反之亦然。这允许精确地回收跨越 JS/C++ 边界的死循环。

#### B. 标记-清除 (Mark-Sweep) 策略

与传统的并行/并发 GC 类似，`cppgc` 使用位图来标记对象活跃性：

- **精确追踪**：通过名为 `Trace` 的方法，C++ 开发者可以显式告诉 GC 对象包含哪些引用。
- **写屏障 (Write Barriers)**：在修改引用时通知 GC，支持增量标记，减少停顿时间。

### 3. 三大优化特性

1.  **增量标记 (Incremental Marking)**：将标记工作分解成小块，在 JS 执行的空闲时间进行，避免长时间的 UI 卡顿。
2.  **并发清除 (Concurrent Sweeping)**：在后台线程回收死对象内存，几乎不占用主线程资源。
3.  **栈扫描优化**：对 C++ 栈进行保守扫描以识别指针，确保不会遗漏任何活跃对象。

### 4. 影响与意义

- **内存显著下降**：Chrome 中的 DOM 对象回收变得更加及时，尤其是在具有复杂生命周期的单页应用中。
- **性能更平稳**：由于统一了回收策略，减少了因为过度防御导致的内存膨胀和不必要的 Full GC。
- **嵌入者友好**：除了 Chrome，其他使用 V8 的 C++ 嵌入程序（如 Node.js 插件或桌面应用）也可以选择开启 `cppgc` 来管理其 C++ 内存。
