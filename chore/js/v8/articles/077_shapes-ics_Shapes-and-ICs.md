# 77. Shapes and ICs | 图形与内联缓存：V8 对象模型核心解析

理解 V8 的核心在于理解它如何利用 **Shapes (Hidden Classes)** 和 **Inline Caches (ICs)** 来伪造静态语言的性能。

## 1. Shapes (Hidden Classes)

JavaScript 是动态类型的，属性查找通常需要极其昂贵的 `Hash Map` 搜索。

- **原理**：V8 给每个对象关联一个 "Shape"（或 Map）。Shape 记录了对象的布局（属性名称及其在内存中的偏移量）。
- **转换链 (Transition Chains)**：当你给空对象 `o` 添加属性 `x` 时，它的 Shape 从 $S_0$ 转移到 $S_1$。如果再添加 `y`，转移到 $S_2$。
- **共享 Shape**：拥有相同属性顺序的对象将共享同一个 Shape，从而节省内存。

## 2. Inline Caches (ICs)

即使有了 Shape，引擎仍然需要检查 Shape 并查找属性。为了进一步加速，V8 引入了 **IC**。

- **原理**：在执行 `o.x` 的代码位置（Call Site），V8 会记住上一次访问时 `o` 的 Shape 以及属性 `x` 的位置。
- **状态转移**：
  - **单态 (Monomorphic)**：总是看到相同的 Shape。性能最优，可直接内联偏移量。
  - **多态 (Polymorphic)**：看到少数几种 Shape。生成一段判断逻辑（if-else）。
  - **超态 (Megamorphic)**：看到大量不同的 Shape。退化到常规的慢速查找。

## 3. 性能建议

1. **初始化顺序一致**：始终以相同的顺序初始化对象属性，以便它们共享同一个 Shape。
   ```javascript
   // Bad
   const p1 = { x: 1, y: 2 }
   const p2 = { y: 2, x: 1 } // 产生不同的 Shape
   ```
2. **避免动态添加属性**：尽量在构造函数或初始字面量中定义所有属性。

## 4. 内部演进

从 v6 时代开始，V8 开始将这些元数据存储在 **FeedbackVector** 中，这使得处理逻辑更加模块化，也更容易被 TurboFan 优化。
