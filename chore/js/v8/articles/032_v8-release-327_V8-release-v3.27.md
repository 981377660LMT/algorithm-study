# 编译器的博弈：Crankshaft 的最后辉煌

# V8 release v3.27: Refined optimizations and property access

这是 V8 在全面转向 TurboFan 之前的平稳迭代期。

### 1. 属性访问的极速缓存

优化了常量（Literal）对象的属性查找速度。
如果你在循环中频繁访问一个对象的固定属性，Crankshaft 现在能生成更精简的机器码片段。它通过直接硬编码（Hardcoding）偏移量的方式，跳过了 Hidden Class 的动态检查（前提是通过 OSR 判定代码极其热点）。

### 2. 生成器（Generators）的进一步稳固

虽然已有支持，但此版本修复了 Generator 在复杂的 `try-finally` 块中的状态恢复 bug。这涉及到回溯执行栈时，对上下文（Context）对象的精确重建。

### 3. 一针见血的见解

v3.27 展示了 Crankshaft 的局限性：为了每一丁点性能提升，代码复杂度都在指数级上升。这种“螺丝壳里做道场”的优化，最终促成了下一代编译器 TurboFan（海之节点架构）的加速研发。这不仅是性能的竞赛，更是工程架构可维护性的决战。
