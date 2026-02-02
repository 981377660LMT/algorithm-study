# WebAssembly 的推测优化：去优化与内联的技术飞跃

## Speculative Optimizations for WebAssembly using Deopts and Inlining

- **Original Link**: [https://v8.dev/blog/wasm-speculative-optimizations](https://v8.dev/blog/wasm-speculative-optimizations)
- **Publication Date**: 2025-06-24
- **Summary**: V8 为 WebAssembly 引入了原本在 JS 中成熟的推测性优化（Speculative Optimizations）。通过收集运行时反馈实现函数内联，并结合去优化（Deoptimization）机制，平衡了 WasmGC 等动态分发特性带来的性能开销。

---

### 1. 背景：WebAssembly 的性能瓶颈与 AOT 的局限性

在最初的设计中，WebAssembly (Wasm) 被视为一种静态类型的字节码，其执行高度依赖于 AOT（提前编译）的确定性。然而，随着 **WasmGC**（面向托管语言的垃圾回收提案）的引入，Java、Kotlin 和 Dart 等语言大量进入 Wasm 领域，带来了新的挑战：

- **动态分发 (Dynamic Dispatch) 的开销**：这些语言高度依赖虚函数和接口。在 Wasm 中，这表现为 `call_indirect` 或 `call_ref` 指令。
- **AOT 的防御性编译**：由于底层 AOT 编译器无法在编译时确定这些间接调用的目标，必须生成通用的、包含多重检查的防御性机器码，这不仅指令冗长，更阻碍了**函数内联**这一核心优化手段。
- **性能“天花板”**：缺乏运行时的真实反馈，使得 Wasm 编译器在面对高度动态的托管语言代码时，性能提升遇到了瓶颈。

---

### 2. 核心改进：推测优化 (Speculative Optimizations)

V8 将其在 JavaScript 引擎中积累多年的**运行时反馈（Type Feedback）**机制移植到了 Wasm 中，打破了“Wasm 必须是纯静态”的固有思维。

#### A. 反馈向量 (Feedback Vector) 的引入

V8 的基准编译器 **Liftoff** 现在会在生成热点函数代码时，插入一段轻量级的监控逻辑：

- 记录 `call_indirect` 或 `call_ref` 在运行时实际调用的函数索引（Function Index）。
- 统计调用的目标分布及其频率。

#### B. 分层编译与推测性内联 (Speculative Inlining)

当函数被标记为“热点”并转交给 **TurboFan**（V8 的高级优化编译器）时，TurboFan 会读取反馈向量：

- **模式识别**：如果反馈显示某个间接调用在 99% 的情况下都指向同一个函数 `A`。
- **激进优化**：TurboFan 会假设下一次调用仍然是 `A`，并直接将函数 `A` 的机器码**内联**到当前函数中。内联后，编译器可以进行跨函数的常量折叠、局部变量重命名和死代码消除。

---

### 3. 关键机制：去优化 (Deoptimization/Deopts)

推测性优化本质上是一种“赌博”，既然有赌博，就必须有“认输”的退路。

- **目标检查 (Target Check)**：在内联代码的最前端，V8 会插入一个极简的检查（Check），验证当前的调用目标是否确实是推测的那个。
- **去优化入口 (Deopt Entry)**：一旦检查失败（例如，一个接口调用突然指向了另一个实现类），执行流会立即跳转到 `DeoptimizationEntry`。
- **状态恢复**：V8 会将当前优化的、可能已被高度“扭曲”的寄存器状态和堆栈记录序列化，重新映射（Remap）回 **Liftoff**（基准编译器）能够理解的原始状态。
- **平滑回退**：程序会在未优化的代码版本中继续执行，直到 V8 下一次根据新的反馈重新进行优化。

**为何不直接用慢速路径？**
如果只使用标准的条件分支（Slow Path），优化编译器仍需考虑分支出口后的所有可能状态，这极大限制了指令级并行和重排。而 **Deopt** 物理上移除了热点路径以外的逻辑，使得剩下的机器码几乎与原生手写的 C++ 一样紧凑。

---

### 4. 针对 WasmGC 的特化处理

- **`call_ref` 优化**：WasmGC 中的函数引用包含了实例信息。V8 通过推测当前实例（Instance）的一致性，使得跨实例的函数调用也能被高效内联。
- **多样性支持**：支持了 `call`、`call_ref`、`call_indirect` 及其对应的尾调用 (Tail-call) 变体的全方位内联。

---

### 5. 实验数据与性能提升

- **Dart 微基准测试**：在结合内联与去优化后，性能提升超过 **50%**。
- **Google Sheets (Wasm 版)**：在处理大规模公式计算引擎时，受益于去优化机制，性能提升约 **7%**。
- **SQLite 3 (Wasm)**：作为计算密集型应用，整体也获得了约 **1%** 的性能增益。

---

### 6. 实战意义与启示

1.  **动态特性的“原生级”性能**：这一改进意味着像 Java/Kotlin 这样高度依赖多态和动态分发的语言，在 V8 上的运行效率正迅速逼近 C++/Rust 这种静态分发的语言。
2.  **打破 AOT 的迷思**：运行时反馈（PGO - Profile-Guided Optimization）在复杂应用面前往往比纯粹的离线编译更有效，V8 证明了 JIT 动态优化的力量。
3.  **开发者建议**：虽然 V8 能够处理多态调用，但维持调用目标的稳定性仍然有助于避免频繁的 **Deopt**。
