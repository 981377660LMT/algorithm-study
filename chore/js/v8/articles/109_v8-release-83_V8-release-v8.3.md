# V8 v8.3 发布：ArrayBuffer 并行回收与 4GB Wasm 内存

## V8 release v8.3

- **Original Link**: [https://v8.dev/blog/v8-release-83](https://v8.dev/blog/v8-release-83)
- **Publication Date**: 2020-05-04
- **Summary**: V8 v8.3 重点改进了 `ArrayBuffer` 的回收效率，支持了高达 4GB 的 WebAssembly 内存，并修复了原型链上 TypedArray 的属性查找 Bug。

---

### 1. 核心性能改进

#### A. ArrayBuffer 并行追踪与释放

在 V8 堆中，`ArrayBuffer` 对象本身很小，但其背后的数据（Backing Store）可能很大且分配在堆外。

- **旧瓶颈**：以往回收 `ArrayBuffer` 时，需要在主线程完成大量的 Backing Store 释放工作，尤其是在处理成千上万个缓冲区时会造成显著停顿。
- **新技术**：引入了全新的追踪机制，允许垃圾回收器在后台线程并发地迭代和释放不再使用的 Backing Store。
- **效果**：在 `ArrayBuffer` 密集型应用中，GC 停顿时间减少了约 50%。

#### B. 更大的 WebAssembly 内存

根据 Wasm 规范的更新，V8 现在允许 Wasm 模块请求最高达 **4GB** 的内存（在此之前通常限制在 2GB 左右）。这为大型游戏、图像处理工具等内存渴求型应用在浏览器中运行铺平了道路。

### 2. 重要修复与 API 变动

#### A. 原型链上的 TypedArray 陷阱

JavaScript 规范规定，在向对象存储属性时，需要检查原型链。V8 发现了一个 Bug：如果 `TypedArray` 位于原型链上且访问越界时，V8 的快速查找处理器可能会错误地安装处理器。

- **例子**：`v = {}; v.__proto__ = new Int32Array(1); v[2] = 123;` 应该让 `v[2]` 仍然是 `undefined`（因为越界写入被忽略），但旧版 V8 会将其存入 `v`。
- **修复**：v8.3 修正了这一逻辑，确保了符合规范的行为。

#### B. WeakRefs API 调整

experimental 的 `WeakRefs` 和 `FinalizationRegistry` API 进行了重命名和重构（例如 `FinalizationGroup` 改名为 `FinalizationRegistry`）。现在的清理任务由 V8 自动调度，不再需要嵌入者（Embedder）手动介入。

### 3. 开发者启示

开发者现在可以更放心地在 JS 中管理大规模二进制数据，因为 `ArrayBuffer` 的生命周期管理变得更加轻量。同时，4GB 的 Wasm 限制提升也标志着 Web 平台计算能力的进一步进化。
