# WebAssembly 类型反射：打破核心边界的“透明化”革命

# WebAssembly Type Reflection: Breaking the Core Boundary with Transparency

### 1. 背景：JS 与 Wasm 的边界黑盒 (The Boundary Black Box)

在类型反射 API 出现之前，JS 与 Wasm 之间的交互处于一种**“半盲”状态**。

- **黑盒困境**：虽然 Wasm 是强类型的，但在 JS 侧，所有导出的 Wasm 函数都被抽象为通用的 `Function` 对象。
- **自省缺失**：JS 脚本在运行时无法通过编程方式得知某个 `WebAssembly.Function` 到底需要多少个参数，或者这些参数是 `i32` 还是 `f64`。
- **依赖先验知识**：开发者必须依赖预先生成的胶水代码（如 Emscripten 生成的库）或静态协议。这种“静默契约”一旦在动态加载或插件化场景下失效，就会导致运行时错误。

### 2. 核心挑战：JS BigInt 与 Wasm i64 的类型桥梁 (The i64 Bridge)

跨越语言边界时，`i64` 是最具挑战性的类型：

- **精度鸿沟**：JS 的传统 `Number` 无法精确表示 Wasm 的 64 位整数（`i64`）。
- **BigInt 强制性**：为了兼容，所有进入 Wasm 的 `i64` 必须是 JS 的 `BigInt`。
- **转换复杂性**：
  - **动态陷阱**：如果不透明，JS 调用者可能误传一个 `Number` 给 `i64` 参数，导致抛出 `TypeError`。
  - **自适应封装**：类型反射 API 引入了 `WebAssembly.Function` 构造函数，允许 JS 函数手动绑定 Wasm 签名。这意味着 JS 可以根据反射出的签名，动态地将数值转换为 `BigInt`，确保桥接层的类型安全。

### 3. API 扩展：`type()` 方法的内省能力 (Introspection via `type()`)

该提案的核心是在 Wasm 各种对象（Function, Global, Memory, Table）上扩展了 `type()` 方法。对于函数而言：

- **自省导出**：通过 `func.type()`，你可以直接获取一个包含 `parameters` 和 `results` 数组的对象。
  ```javascript
  // 返回示例
  {
    parameters: ["i32", "i64"],
    results: ["f32"]
  }
  ```
- **对称构造**：不仅能读取，还能定义。`new WebAssembly.Function(type, jsFunc)` 允许开发者在 JS 侧为纯 JS 函数“贴上”Wasm 的类型标签。

### 4. 应用场景：动态元编程与类型安全 FFI (Dynamic Metaprogramming & FFI)

- **动态加载器与链接器**：在多模块协作场景中（如插件系统），宿主程序可以在运行时检测 Wasm 插件的签名，动态生成符合规范的适配器。
- **自动化的 FFI 层**：开发者可以编写通用的封装框架，自动处理 JS 对象到 Wasm 结构体的映射，而无需为每个函数手动编写声明。
- **调试与代理（Instrumentation）**：拦截器（Interceptors）或日志代理可以在不知道原始签名的情况下，通过反射自动记录函数的入参和出参。

### 5. 洞察：从“盲目调用”到“知己知彼”的互操作性升级

**一针见血的见解**：
类型反射（Type Reflection）打破了 Wasm 的**“静默契约”**。
在过去，Wasm 相对于 JS 像是一个“不可观测的零件”，只能按照预设的方式运转。现在的升级让 Wasm 变成了一个**“自修复且自解说的透明组件”**。

这种转变标志着 Wasm 从“嵌入式二进制包”向“完全集成化的运行时对象”演进。它让 JS 真正拥有了治理 Wasm 资产的能力：

> **不再是“我希望这个调用是对的”，而是“我知道这个调用必须怎么做”。**

这种“知己知彼”的状态，是构建复杂、高性能、且具备高度动态性的 Wasm 运行时组件（如 Wasm-based Serverless 架构）的基石。
