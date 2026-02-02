# V8 v5.6 发布：Ignition 与 TurboFan 的实战演练

# V8 release v5.6: Testing the new pipeline on ES.next

V8 v5.6 标志着新架构开始大规模接手现代 JavaScript 特性的执行。

### 1. 现代语法的“新家”：Ignition + TurboFan

在此版本之前，新特性（如 Generator）是通过旧架构缓慢模拟的。从 v5.6 开始，所有的 ES2015+ 特性以及 Crankshaft 无法优化的代码段（如 `try-catch`）都开始正式运行在 **Ignition (解释器) + TurboFan (优化编译器)** 组合上。
结果立竿见影：Generator 性能提升了近 3 倍。

### 2. Orinoco 项目：并发过滤记忆集

引入了 **Concurrent remembered set filtering**。记忆集（Remembered Sets）记录了老生代对新生代对象的引用。通过在后台线程并行过滤这些引用，主线程在执行 Scavenge（新生代回收）时的停顿时间大幅缩减。

### 3. 一针见血的见解

v5.6 是 V8 在“修飞机时换引擎”的阶段。通过让新架构优先负责复杂的 ES6 特性，V8 在不影响老代码稳定性的前提下，成功让现代 JS 语法摆脱了“慢”的标签。
