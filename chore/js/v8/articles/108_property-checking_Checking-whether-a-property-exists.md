# 检查属性是否存在：性能深度对比

## Checking whether a property exists

- **Original Link**: [https://v8.dev/blog/property-checking](https://v8.dev/blog/property-checking)
- **Publication Date**: 2020-04-14
- **Summary**: 文章深入分析了 `in` 操作符、`hasOwnProperty` 以及 `obj.prop !== undefined` 等属性检查方法的性能差异。通过剖析 V8 内部的 Inline Caches (IC) 机制，揭示了不同场景下的最优选。

---

### 1. 三种主要的检查方式

在 JavaScript 中，判断对象是否有某个属性通常有以下写法：

1.  `'prop' in obj`：检查自身及其原型链。
2.  `obj.hasOwnProperty('prop')`：仅检查自身属性。
3.  `obj.prop !== undefined`：最常用的简写，但在属性值确实为 `undefined` 时失效。

### 2. V8 内部的优化机制

V8 使用 **Inline Cache (IC)** 来加速属性访问。

#### A. `in` 操作符的优化

对于 `in` 操作符，V8 会根据对象的 **Hidden Class (Map)** 生成专门的代码。

- 如果属性在对象的 Map 中，V8 可以直接返回 `true` 而无需查询底层的哈希表。
- 即使在原型链上，V8 也会缓存查找结果。

#### B. `hasOwnProperty` 的开销

尽管 `hasOwnProperty` 在逻辑上更直接，但在某些 V8 版本中，如果它不是作为内置函数被快速调用（例如通过 `Object.prototype.hasOwnProperty.call(obj, 'prop')`），它可能比 `in` 慢。

- **改进**：V8 后来通过内联（Inlining）和生成特定的字节码指令大大优化了 `hasOwnProperty`。

### 3. 性能排行与建议

- **最高执行效率**：直接比较 `obj.prop !== undefined`。
  - **逻辑**：V8 将其简化为一次属性加载。如果加载失败（空槽），则返回 `undefined`。
  - **限制**：不能区分“属性不存在”和“属性值为 undefined”。
- **语义最准确**：`in`。
  - 即使属性在原型链上，由于 IC 的存在，其速度几乎与直接加载一样快。
- **最安全的“自身”检查**：`Object.hasOwn()`（当时尚未完全普及时，推荐使用 `hasOwnProperty`）。

### 4. 技术深度：从 Map 到常量

如果一个对象被 V8 视为“快速对象”（Fast Object），属性信息存储在 Map 中。属性检查本质上变成了：

1.  加载对象的 Map。
2.  在 Map 的描述符数组中执行二进制或线性搜索。
3.  在 IC 预热后，步骤 2 被跳过，直接根据 Map 地址对比跳转。

### 5. 结论

在对性能要求极高的热点代码中，如果逻辑允许，优先使用 `!== undefined`。如果需要处理原型链或必须确认属性存在，`in` 操作符是性能与语义兼顾的最佳平衡点。
