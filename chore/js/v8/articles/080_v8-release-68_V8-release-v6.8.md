# 80. V8 release v6.8 | V8 v6.8 发布：Array.flat 与 TypedArray 优化

V8 v6.8 带来了显著的性能提升和内存优化，重点在于内建函数的 CSA 化（CodeStubAssembler）以及内存结构的精简。

## 1. 数组新特性：Array.prototype.flat / flatMap

这两个方法正式从草案进入实现。V8 使用 **CodeStubAssembler (CSA)** 实现了这些方法，确保了在处理大规模嵌套数组时的性能。

- `flat`：支持可选参数 `depth`，通过递归展平数组。
- `flatMap`：结合了 `map` 和 `flat(1)`，避免了中间数组的分配，极大提高了效率。

## 2. 内存优化：SharedFunctionInfo 与 ScopeInfo 的分离

在 v6.8 之前，`SharedFunctionInfo` (SFI) 对象与其关联的 `ScopeInfo` 紧密耦合。

- **优化点**：V8 团队发现许多函数并不需要完整的 `ScopeInfo`。通过将两者分离，并使 `ScopeInfo` 仅在确实需要（如闭包访问）时分配，从而减少了每个函数实例的内存开销。
- **效果**：这使得在大型网页（特别是含有大量 JS 函数的页面）中的内存占用显著下降。

## 3. TypedArray.prototype.sort 的性能飞跃

在此版本之前，`TypedArray.sort` 使用的是通用数组的 C++ 实现。

- **改进**：V8 重新用 **CSA** 实现了针对 TypedArray 的排序算法。
- **原理**：由于 TypedArray 元素类型固定（Int32, Float64 等），CSA 能够生成高度优化的机器码，避免了普通数组排序中频繁的类型检查和属性查找开销。

## 4. BigInt 比较操作的优化

v6.8 进一步完善了 BigInt 的操作。

- **混合比较**：支持 `BigInt` 与 `Number` 的直接比较（如 `10n > 5`）。这涉及到在内部进行精密的数值转换和比较逻辑，以防止由于精度差异导致的错误判断。

## 5. 垃圾回收：Orinoco 持续推进

此版本继续优化了并发标记（Concurrent Marking）的稳定性，减少了主线程在标记阶段的停顿时间。
