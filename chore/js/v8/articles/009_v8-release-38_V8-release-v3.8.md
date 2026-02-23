# 类型系统的深度推断

# V8 release v3.8: Advanced type inference and Hydrogen polish

### 1. 完善 Hydrogen IR

在 v3.8 中，Crankshaft 编译器的 **Hydrogen Intermediate Representation** 变得更加成熟。它能够处理更复杂的 `try-catch` 块和 ES5 的新特性。
V8 改进了在 Hydrogen 层面的类型预测算法，使得编译器在处理包含复杂逻辑分支的函数时，依然能保持很高的内联成功率。

### 2. 精确式 GC 的进一步优化

通过改进对 **Stack Roots**（栈上的根对象）扫描的路径，缩短了垃圾回收初期的标记准备时间。

### 3. 一针见血的见解

JIT 编译器的强弱不仅取决于它能优化的天花板，更取决于它能覆盖的代码地平线。v3.8 的工作让越来越多的“复杂代码”也能享受到 Crankshaft 的加速红利，而不再是仅仅针对微型基准测试起效。
