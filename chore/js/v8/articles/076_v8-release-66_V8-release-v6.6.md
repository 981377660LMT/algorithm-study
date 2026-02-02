# 76. V8 release v6.6 | V8 v6.6 发布：加速异步性能与 try-catch 优化

V8 v6.6 正式标志着对现代 JS 特性的全面高性能支持，特别是异步函数的开销大幅降低。

## 1. 异步性能优化：Promise 与 Async/Await

此版本通过 **CodeStubAssembler** 重构了 Promise 的内部实现，并优化了 `async/await` 生成的字节码。

- **减少微任务调度**：在不违反规范的前提下，精简了微任务执行的中间环节。
- **结果**：部分典型的异步基准测试性能提升了 **20% 到 50%**。

## 2. try-catch 性能提升

历史上，`try-catch` 会阻止编译器对所在函数进行许多优化（如内联）。

- **优化点**：Ignition 解释器更好地处理了异常范围。TurboFan 现在能够更智能地对含有 `try-catch` 的代码块进行分析。
- **新特性**：支持草案中的 **Optional Catch Binding** (`try { ... } catch { ... }` 省略异常变量)。

## 3. JSON ⊂ ECMAScript

之前某些特殊的 Unicode 字符（如 `U+2028` 和 `U+2029`）在 JSON 中有效但在 JS 字符串文本中无效会报错。v6.6 修复了这一不匹配，使得 JS 源代码可以包含合法的 JSON 文本。

## 4. 内建函数提速

- **Array.prototype.reduce**：通过 CSA 优化大幅提升了执性能。
- **String.prototype.trim**：也获得了显著加速。

## 5. 背景编译 (Background Compilation)

V8 开始将更多的脚本编译工作从主线程转移到后台线程。通过在网络下载的同时并发进行字节码生成，主线程被显著释放，网页交互性（FID）更佳。
