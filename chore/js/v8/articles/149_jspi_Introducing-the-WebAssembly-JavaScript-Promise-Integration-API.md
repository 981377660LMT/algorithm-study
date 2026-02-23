# JSPI：深度解析 Wasm 与 JS 异步鸿沟的终极桥梁

## JSPI: Deep Dive into the Ultimate Bridge Between Wasm and JS Async Worlds

- **Original Link**: [https://v8.dev/blog/jspi](https://v8.dev/blog/jspi) (Intro) & [https://v8.dev/blog/jspi-newapi](https://v8.dev/blog/jspi-newapi) (Update)
- **Publication Date**: 2024-07-01 & 2024-06-04
- **Summary**: **JSPI (JavaScript Promise Integration)** 是 V8 引入的一项核心技术，旨在解决 Wasm 的顺序执行逻辑与 JS 非阻塞异步环境之间的本质矛盾。通过在 VM 层引入 **栈切换（Stack Switching）** 机制，它允许 Wasm 代码在调用异步 JS 函数时“暂停”并在 Promise 完成后“恢复”，且无需对二进制代码进行昂贵的 Asyncify 改写。

---

### 1. 核心难题：协议栈对话的本质矛盾

Wasm 与 Web 平台的对话本质上是两种执行哲学的对抗：

- **Native 逻辑 (Wasm)**：倾向于**“阻塞式顺序执行”**。在 C++/Rust 编写的传统库中，代码逻辑是线性的。例如，调用一个网络请求函数，代码会挂起直到数据返回。
- **Web 运行时 (JS)**：基于**“非阻塞异步事件循环”**。主线程绝不允许阻塞，任何 I/O 操作都必须通过 Promise 异步处理。

**为什么 Asyncify 不是长久之计？**
在 JSPI 出现前，开发者被迫使用 **Asyncify** 工具。它通过静态代码改写，在每个函数调用前后插入保存/恢复状态的逻辑。这导致：

1. **代码膨胀**：生成的 Wasm 体积通常增加 50% 以上。
2. **性能损耗**：即使是非异步路径也会受到状态检查的拖累。
3. **不可维护**：很难处理复杂的递归和动态链接。

---

### 2. JSPI 的本质：VM 级的无缝栈切换 (Stack Switching)

JSPI 并非在软件层面模拟异步，而是在 **V8 引擎级别引入了“协程（Coroutine）”式的栈管理**：

- **隔离栈空间**：V8 为 JSPI 托管的 Wasm 执行流分配了一个**独立的逻辑栈**，与 JS 的主调用栈完全隔离。
- **挂起 (Suspend)**：当 Wasm 调用一个被标记为 `Suspending` 的 JS 异步函数并收到一个未决的 Promise 时，V8 不会阻塞主线程。它会记录 Wasm 当前的寄存器状态和栈指针，然后**“切出（Switch Out）”** Wasm 栈，将控制权交还给 JS 事件循环。
- **恢复 (Resume)**：当 Promise resolve 后，事件循环触发回调。V8 将执行环境**“切回（Switch Back）”**对应的 Wasm 栈，恢复现场。对于 Wasm 代码而言，仿佛只是经历了一个稍长的同步函数调用。

这种切换是在用户态由 V8 通过修改 IP 和 SP 完成的，不涉及内核上下文切换，开销极低（实测约 1µs）。

---

### 3. API 的演进：从复杂到显式 (v1 vs v2)

V8 在 [148](https://v8.dev/blog/jspi-newapi) 中更新了 API 设计，反映了从“机制暴露”到“语义封装”的转变：

- **旧 API (v1)**：引入了显式的 `Suspender` 对象，开发者必须手动在 JS 胶水代码和 Wasm 导出中传递它。这导致了高度的耦合和代码污染。
- **新 API (v2)**：引入了更简洁的 `WebAssembly.promising` (导出) 和 `WebAssembly.suspending` (导入) 包装器。
  - **边界触发**：异步上下文的生命周期现在由 Wasm 模块的导出边界自动管理。
  - **逻辑解耦**：Wasm 内部代码完全不需要感知 `Suspender` 的存在，真正实现了二进制级别的“无感异步”。

---

### 4. 关键优化：分段栈 (Segmented Stacks)

为了应对潜在的内存压力，V8 正致力于实现**可增长的分段栈**。

- **痛点**：传统的协程如果每个都分配 1MB 固定栈，并发数千个就会耗尽内存。
- **解决方案**：JSPI 初始分配极小的栈空间，并根据需求动态增加。这使得一个 Web 应用可以支持成千上万个并发的 Wasm 执行路径，极大地扩展了复杂应用的并发能力。

---

### 5. 实战意义：Native 生态上云的“最后一块拼图”

JSPI 极大地拓宽了 Wasm 的应用边界：

1. **数据库驱动**：如 SQLite、PostgreSQL。这些库底层的同步 I/O 逻辑无需重写，即可直接对接浏览器的 `fetch` 或 `IndexedDB`。
2. **网络协议栈**：OpenSSL 等库可以作为同步中间件运行，对接 JS 异步网络层。
3. **动态加载逻辑**：允许 Wasm 逻辑在中途挂起，去 fetch 并异步编译另一个 Wasm 模块，然后再无缝恢复执行。

### 一针见血的技术洞察

JSPI 是 **VM 层面并发原语的一次彻底解放**。它将 Wasm 从“JS 的纯计算加速插件”提升到了“拥有独立调度能力的运行时实体”。通过从底层抹平“同步阻塞与异步循环”之间的鸿沟，它正式宣告了：**任何原本在桌面端运行的高性能 C++/Rust 全量代码库，现在都可以零修改（逻辑层面）地在浏览器中复活。**
