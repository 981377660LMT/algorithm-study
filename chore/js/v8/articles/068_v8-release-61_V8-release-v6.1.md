# V8 v6.1 发布：再见，Crankshaft！

# V8 release v6.1: Saying goodbye to Crankshaft

Crankshaft 作为 V8 曾经的王牌优化器，正式被 TurboFan 取代并移除。

### 1. 技术更迭：TurboFan 的全面胜利

长期以来，Crankshaft 虽然快速，但它无法处理 ES6 的高级特性（如 `try-catch`、`eval`、`with`、`generators`）。TurboFan 凭借其创新的 **Sea-of-Nodes** 设计，不仅支持所有 ES 语法，还在代码生成的紧凑度上全面超越了前辈。

### 2. 迭代器性能暴走

`Map` 和 `Set` 的迭代性能在 v6.1 中提升了惊人的 11 倍。优化的核心在于：不再通过 C++ 运行时中转，而是利用 TurboFan 的内建函数编译器直接生成原生的寄存器操作逻辑（Inlined Fast Paths）。

### 3. Wasm 渐入佳境

WebAssembly 开始支持 Web 平台的基本特性，V8 对 asm.js 的支持也开始转入“转译为 Wasm”的路径，这意味着 JS 社区曾经提倡的“近机速环境”终于有了标准化的统一目标。

### 4. 一针见血的见解

v6.1 的意义在于“统一”。TurboFan 统一了 JS 各个子集（ES5/ES6/asm.js）的优化入口，其强大的通用优化算法（如逃逸分析、调度算法）开始全局生效，让原本各行其是的优化路径汇聚成了一股合力。
