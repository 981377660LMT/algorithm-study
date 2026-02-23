# V8 v5.7 发布：异步性能爆发与 Wasm 默认启用

# V8 release v5.7: Async function performance and WebAssembly by default

V8 v5.7 标志着两个支撑现代 Web 的核心技术的成熟。

### 1. 异步函数 (Async Functions) 提速 4x

通过将 `async`/`await` 重新实现为基于 **Generator** 的底层抽象并由 TurboFan 进行去调度化优化，异步函数的执行效率提升了 4 倍。在此之前，async 函数通常比 Promise 慢得多，而 5.7 是两者性能对比的天平开始倾斜的时刻。

### 2. WebAssembly (Wasm) 正式启用

这是 Web 历史上的重大时刻：Chrome 57 默认开启 WebAssembly 支持。V8 团队不仅实现了规范，还引入了初步的分层编译想法，确保 Wasm 模块能够在保持近机速执行的同时，尽可能快地启动。

### 3. 正则表达式优化

正则表达式引擎开始利用新编译器的红利。通过在 TurboFan 中内联常见的正则匹配路径，匹配速度提升了约 15%，且减少了跨 JS/C++ 边界带来的开销。

### 4. 一针见血的见解

v5.7 是“赋能之作”。它解决了异步代码的性能隐患，并引入了 Wasm 这一强援。从此，Web 开发者不再被局限于单线程低效代码，异步化与二进制计算成为了对抗复杂业务逻辑的两大利器。
