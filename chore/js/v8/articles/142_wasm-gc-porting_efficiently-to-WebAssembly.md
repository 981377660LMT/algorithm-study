# WasmGC：开启托管语言的 Web 高性能时代

## A new way to bring garbage collected programming languages efficiently to WebAssembly

- **Original Link**: [https://v8.dev/blog/wasm-gc-porting](https://v8.dev/blog/wasm-gc-porting)
- **Publication Date**: 2023-11-01
- **Summary**: **WasmGC** 预示着 WebAssembly (Wasm) 从“模拟低级硬件”向“现代高级虚拟机”的重大转型。它允许 Java、Dart、Kotlin 等托管语言直接利用浏览器的垃圾回收器，彻底告别了以往必须自带冗余 GC 环境的“笨重”时代。

---

### 1. 背景：Wasm 的“无 GC”尴尬期

在 WasmGC 出现之前，Wasm 仅支持**线性内存 (Linear Memory)**——这本质上是一块巨大的、扁平的字节数组。

- **自带 GC 运行时的包袱**：Java、Kotlin 或 Dart 等语言依赖垃圾回收。在 WasmGC 之前，这些语言被移植到 Wasm 时，必须将自己的一套完整 GC 逻辑编译进 Wasm 二进制中。
- **副作用**：
  - **加载极慢**：二进制文件包含大量重复的 GC 运行时代码，体积臃肿。
  - **内存泄漏天坑**：Wasm 内部对象被隐藏在线性内存中，浏览器的 JS GC 无法直接感知它们。如果 Wasm 和 JS 对象之间存在循环引用，传统的 GC 算法无法回收这些内存。

---

### 2. WasmGC 的本质：语义的升维

WasmGC 的核心逻辑是：**让 Wasm 能够直接“对话”宿主环境（浏览器）的垃圾回收器。**

这意味着 Wasm 不再只是在一个“地窖”里搬运字节，而是可以像 JavaScript 一样，定义并操作**受控对象 (Managed Objects)**。Wasm 模块可以声明 Structs（结构体）和 Arrays（数组），而 V8 引擎负责这些对象的分配和生命周期管理。

---

### 3. 技术突破：互操作性与零拷贝

- **统一的 GC 视野**：由于 WasmGC 对象与 JS 对象共用 V8 内部的同一个垃圾回收器，GC 可以无缝追踪跨越 JS/Wasm 边界的引用链。这完美解决了由于循环引用导致的内存泄漏难题。
- **零拷贝交互**：Wasm 产生的对象可以直接作为引用传递给 JS 逻辑，无需进行昂贵的二进制序列化或内存拷贝。
- **高效的内联优化**：由于 Wasm 对象携带了显式的类型信息，V8 可以在 Wasm 执行期间进行**推测性内联 (Speculative Inlining)**，这让托管语言在 Web 上的性能提升了 30% 以上。

---

### 4. 移植策略：从“模拟”到“映射”

文章详细解释了如何将传统面向对象语言移植到 WasmGC 系统：

- **类型映射**：编译器（如 Kotlin/Wasm）将源代码中的类系统直接映射到 WasmGC 的 `struct` 和 `array` 系统。
- **轻量化二进制**：由于不再需要自带 GC 运行时，一个原本数 MB 级的 Kotlin 或 Dart 应用可以被压缩到数百 KB，显著加快了首屏加载速度。

---

### 5. 实战启示：全能 Web 平台的崛起

WasmGC 的成熟是 Web 平台的一次飞跃：

- **Flutter on Web**：通过 WasmGC 重新编译，Flutter 现在能在浏览器中实现极其流畅、甚至接近原生的渲染性能。
- **复杂工具应用**：Google Sheets 利用 WasmGC 大幅优化了其在浏览器内的重型计算引擎，显著降低了内存占用并提升了计算吞吐量。

### 一针见血的技术洞察

WasmGC 标志着 WebAssembly 已经从**“底层系统语言（C++/Rust）的避风港”**演变成了**“全能的高性能语言执行基座”**。通过将内存管理的主导权交还给引擎，它消除了托管语言在 Web 端长久以来的“二等公民”尴尬。从此，语言的性能不再取决于其是否“原生”，而取决于其如何优雅地与宿主 GC 协作。
