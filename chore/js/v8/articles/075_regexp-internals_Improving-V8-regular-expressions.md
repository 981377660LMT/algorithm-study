# 75. Improving V8 regular expressions: magic numbers and Unicode | 提升 V8 正则表达式：魔术数字与 Unicode 支持

正则表达式引擎 Irregexp 在此版本中经历了一次重大重构，旨在提供更符合标准的 Unicode 支持和更高的执行效率。

## 1. Unicode 属性转义 (Unicode Property Escapes)

支持了 `\p{...}` 语法，允许根据 Unicode 属性进行匹配（如 `\p{Script=Greek}`）。

- **实现原理**：由于 Unicode 属性涵盖了数千个字符，Irregexp 并不是简单地生成巨大的 if-else，而是将这些属性编译成紧凑的字符集合表示（BitMap 或 Range List）。

## 2. 具名捕获组 (Named Capture Groups)

支持了 `(?<name>...)` 语法。

- **元数据扩展**：V8 扩展了正则表达式执行结果的元数据存储。在内部，具名捕获组会被映射到数字索引，但在 `exec` 的结果中通过 `groups` 对象暴露。

## 3. 后行断言 (Lookbehind Assertions)

这也是极具挑战的特性。

- **挑战**：正则引擎通常是向前扫描的。为了支持 `(?<=...)`，Irregexp 实现了反向匹配逻辑或在匹配成功后进行回溯校验。

## 4. 内部优化：从 JS 到 CSA

正则表达式的很多周边逻辑（如 `replace` 中带有全局标志 `/g` 的路径）从自托管的 JavaScript 迁移到了 **CodeStubAssembler**。

- **收益**：以前这一过程涉及大量的 C++ 与 JS 之间的切换开销。迁移到 CSA 后，整个匹配和结果构建流程都在优化的机器码中完成，性能提升显著。

## 5. 魔术数字与 Unicode 范围

V8 通过内部算法自动寻找最优的 Unicode 判定“魔术数字”。例如，对于特定的字符集，选择特定的位操作掩码，利用现代 CPU 的 SIMD 或位运算能力快速过滤掉非匹配字符，这在处理多字节 Unicode 时尤为关键。
