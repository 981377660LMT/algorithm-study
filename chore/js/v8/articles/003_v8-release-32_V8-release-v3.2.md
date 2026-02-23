# V8 release v3.2 发布：ES5 严谨模式与内存优化

# V8 release v3.2: ES5 Strict Mode and Memory Tuning

发布于 2011 年左右，这个版本是 V8 迈向现代化 JavaScript 标准的关键一步。

### 1. ES5 Strict Mode 的全面支持

V8 v3.2 标志着对 ECMAScript 5 严谨模式（Strict Mode）的初步完整支持。通过 `"use strict";` 指令，V8 可以对代码进行更激进的静态推断，因为它消除了 `with` 语句和不确定的全局变量绑定，这直接减少了编译器在生成机器码时的保护性检查。

### 2. Smi (Small Integers) 的深度优化

在 32 位和 64 位系统上，V8 使用指针标记（Pointer Tagging）来区分整数和引用。v3.2 对 Smi 的范围和操作指令进行了优化，特别是在加法和位运算的快速路径中，减少了溢出检查的开销。对于早期高性能计算应用（如浏览器内的物理模拟）意义重大。

### 3. 一针见血的见解

严谨模式不仅是为程序开发者设计的，更是为 JIT 编译器设计的。通过限制 JS 的极端动态性，严谨模式让 V8 能够生成更加接近原生 C 性能的指令流。

---

- **归档链接**: [V8 Blog (Archived)](https://v8.dev/blog)
- **核心关键词**: `Strict Mode`, `ES5`, `Smi Optimization`, `Pointer Tagging`
