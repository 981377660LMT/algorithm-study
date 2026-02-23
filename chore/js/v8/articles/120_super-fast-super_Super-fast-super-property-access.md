# 超快速的 super 属性访问

# Super fast super property access

V8 通过将 `super.prop` 的查找从缓慢的运行时 C++ 调用迁移到 JIT 优化路径，实现了性能的飞跃。

### 1. 技术背景：为什么 super 慢？

在 JS 中，`super.x` 不同于 `obj.x`。它的“查找起始对象”必须是当前方法所属对象的原型（`HomeObject.__proto__`），但其 `this` 接收者（Receiver）必须保持为当前的实例。这种复杂的上下文绑定以往需要回退到 C++ 运行时处理。

### 2. 核心优化：LoadSuperIC

- **指令化**：V8 引入了新的字节码 `LdaNamedPropertyFromSuper`。
- **内联缓存 (IC)**：通过 `LoadSuperIC`，V8 能够记住 `HomeObject` 原型的 Shape（Hidden Class）。一旦稳定，就可以在 IC 层直接根据偏移量定位属性，跳过原型链遍历。
- **常量化处理**：V8 将 `HomeObject` 迁移到了类上下文（Class Context）中。TurboFan 编译器可以将 `HomeObject` 及其 `__proto__` 视为常量嵌入生成的机器码中。

### 3. 一针见血的见解

`super` 访问的“平民化”意味着它从一个昂贵的动态查找特性变回了编译器可预测的静态布局。对于大量使用 ES6 类继承的现代框架来说，这项优化消除了继承层次带来的隐形性能税。
