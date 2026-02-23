# V8 v8.4 发布：私有方法与 WebAssembly 引用类型

# V8 release v8.4: Private methods and WebAssembly reference types

V8 v8.4 在语言特性的深度支持上跨出了一大步，特别是在类的封装性和 Wasm 与 JS 的互操作性上。

### 1. 强力封装：私有方法与访问器

本版本正式支持了私有类方法（Private Methods）和私有 Getter/Setter。这些成员通过 `#` 符号前缀标识，由 V8 在内部通过私有符号（Private Symbols）和特定的闭包引用实现，确保了在 JS 源码层面具有完全的不可访问性，是实现高内聚低耦合类的关键。

### 2. 跨界桥梁：Wasm 引用类型 (`externref`)

引入了 WebAssembly 引用类型（Reference Types）。最显著的是 `externref`，它允许 Wasm 模块直接持有、传递和接收不透明的 JS 对象句柄（如 DOM 元素、JS 函数等）。这消除了以往通过表格/数组索引跳转获取 JS 对象的繁琐间接层。

### 3. 正则引擎优化：属性转义的支持

Unicode 属性转义（`\p{...}`）在正则表达式中获得了更广泛的支持，使得处理多语言文本时的字符分类变得更加可靠和高性能。

### 4. 一针见血的见解

v8.4 标志着 JS 离完成其“面向对象”成熟化又近了一步。而 Wasm 引用类型的加入，则宣告了 Wasm 将不再仅仅局限于一个黑盒计算引擎，它正迅速演变为 Web 环境中与 JS 平起平坐的一等公民，两者之间的边界正变得前所未有的模糊。
