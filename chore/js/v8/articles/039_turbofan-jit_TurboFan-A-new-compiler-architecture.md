# TurboFan：下一代 JIT 编译器架构的诞生

# TurboFan: A new compiler architecture

这是 V8 历史上最重要的技术转折点之一：为了应对 ES6 的复杂性，V8 彻底重写了其编译器核心。

### 1. 架构突破：节点海 (Sea-of-Nodes)

TurboFan 抛弃了 Crankshaft 的线性中间表示（IR），采用了名为“节点海”的架构。

- **解耦控制流与数据流**：在图中，数据依赖和控制依赖被统一表示。这允许优化器进行以前不可能实现的超大规模代码重排和死代码消除（Dead code elimination）。

### 2. 分层与通用性

以往的 Crankshaft 深度绑定特定的机器架构，导致支持新平台（如 ARM64）异常痛苦。TurboFan 采用了高度分层的设计：

- **前端**：解析各种源（JS, Wasm）。
- **中端**：执行与架构无关的通用优化。
- **后端**：生成特定指令。

### 3. 一针见血的见解

TurboFan 不仅仅是为了快，它是为了“生存”。面对 ES6 引入的 `try-catch`、`eval` 等 Crankshaft 无法处理的顽疾，TurboFan 是一次降维打击。它的出现，让 V8 能够以一套统一的数学模型（Sea-of-Nodes）来解决所有复杂的动态语言优化难题。
