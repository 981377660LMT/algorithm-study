# 跌得稳才跑得快：去优化的性能损耗

# V8 release v3.22: Deoptimization improvements and Smi checks

### 1. 降低“坠机”惩罚

在 **Crankshaft** 引擎中，当类型预测（Speculation）失败时，引擎必须触发 **Deoptimization**，将执行流退回至解释器。
v3.22 优化了 **Deoptimization Table** 的内存布局，使得状态恢复的时间缩短。同时，增强了对 **Small Integer (Smi)** 边界检查的静态分析，减少了由于整数溢出导致的代码回退。

### 2. 寄存器分配改进

此版本细化了在热点循环中的寄存器分配逻辑，减少了在循环迭代中不必要的 Spill（将变量写回内存）动作。

### 3. 一针见血的见解

高性能引擎不仅要“跑得极快”，更要“跌得稳”。在高复杂度的 JS 应用中，减少惩罚性的去优化开销，比单纯提高 JIT 峰值吞吐量更能保证用户体感的平滑度。
