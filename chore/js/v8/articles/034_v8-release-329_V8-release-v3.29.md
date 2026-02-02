# 解析器的减负：并行解析初现峥嵘

# V8 release v3.29: Background parsing and script speed

### 1. 彻底解决主线程阻塞：后台解析

在 v3.29 之前，脚本的下载和解析几乎是串行的。
此版本强化了 **Background Parsing**。当脚本正在加载时，V8 会在子线程中预热解析器。这意味着当脚本下载完成时，主线程拿到的可能已经是一个近乎完整的 Abstract Syntax Tree (AST)，大大缩短了 First Meaningful Paint 的时间。

### 2. 字符串匹配的飞跃

改进了 `String.prototype.indexOf` 和相关正则匹配算法。V8 针对短字符串和重复模式做了专门的汇编优化（SIMD 前奏），利用 CPU 的向量指令加速字符对比过程。

### 3. 一针见血的见解

脚本大小随着现代 Web 应用爆炸式增长。V8 v3.29 的核心贡献在于它意识到：不仅要让代码跑得快，更要让代码“醒得快”。多线程参与编译解析，是 V8 应对现代化巨型单页应用（SPA）的第一把利刃。
