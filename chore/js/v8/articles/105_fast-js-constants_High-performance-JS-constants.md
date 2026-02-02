# 高性能 JavaScript 常量：消除“魔法值”的开销

## High-performance JS constants

- **Original Link**: [https://v8.dev/blog/fast-js-constants](https://v8.dev/blog/fast-js-constants)
- **Publication Date**: 2020-01-20
- **Summary**: 长期以来，在全局作用域定义 `const` 或访问全局常量一直比字面量慢。本文讲解了 V8 如何通过全局属性元胞（Global Property Cells）和优化编译器，使常量访问达到原生字面量的速度。

---

### 1. 痛点：为什么 `const` 不够快？

在 JS 中，全局变量（包括 `const`）是存储在全局对象（`window` 或 `global`）的属性中的。

- **查找成本**：每次访问 `MY_CONSTANT`，理论上都要去全局哈希表走一圈。
- **不可预测性**：即使是 `const`，在脚本重新加载或不同上下文中可能表现不同。

### 2. V8 的解决方案：Property Cell

V8 为全局属性引入了 **Property Cell**。

- 它是一个特殊的内部对象，不仅存储属性的值，还存储该属性的 **元数据 (Metadata)**（例如：它是只读的吗？它的类型变了吗？）。

#### A. 常量折叠与内联

当 V8 的 TurboFan 编译器看到一个被标记为 `read-only` 的 Property Cell 时，它会直接将该单元格中的值硬编码到生成的机器码中。

- **效果**：`return MY_CONSTANT` 变成了类似于 `return 42` 的指令，完全跳过了内存查找。

#### B. 依赖追踪与去优化 (Deoptimization)

既然是 JS，就总有例外。如果一个全局变量被修改（尽管 `const` 不允许修改，但全局对象的某些属性可能被删除或重新定义），V8 必须保证正确性。

- V8 使用 **依赖追踪 (Dependency Tracking)**。如果 Property Cell 的内容发生变化，所有依赖于该单元格的优化代码都会立即失效并降级（Deopt）。

### 3. 热重载与 V8 别名

对于像 `NaN`, `Infinity`, `undefined` 这样的语言内置常量，V8 现在使用特殊的底层别名（Aliases）。

- 以前：访问 `undefined` 是一次受保护的加载。
- 现在：它被视为一个根寄存器偏移量，速度等同于读取寄存器。

### 4. 开发者建议

- **拥抱 `const`**：不要为了所谓性能把常量手动替换成字面量。V8 已经能完美识别并优化它们。
- **模块作用域优于全局**：虽然全局常量很快，但模块作用域（Module Scope）的变量在 V8 中有更确定的生存期和依赖树，往往能获得更稳定的优化。

### 5. 总结

通过 Property Cell 技术，V8 抹平了代码可读性（使用命名常量）与性能（使用原始值）之间的鸿沟。
