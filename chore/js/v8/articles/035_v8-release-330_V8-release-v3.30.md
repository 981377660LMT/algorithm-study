# 内存的呼吸：更激进的存活空间回收

# V8 release v3.30: Scavenge optimizations and memory efficiency

### 1. 内存抖动的克星

在 v3.30 中，V8 对 **Scavenger（新生代 GC）** 进行了大幅优化。
新生代内存采用的是 Semi-space 复制算法。优化点在于，当检测到存活对象较多时，V8 会调整复制的并发策略，并更早地将“历经多次（通常是 2 次）Scavenge 仍存活”的对象晋升（Promotion）到老生代，从而腾出宝贵的 From-space 空间给新的分配请求。

### 2. 开发者工具：内存泄漏的显微镜

这一版增强了 Heap Snapshot 的精度。开发者现在可以不仅仅看到对象占用了多少字节，还能更清晰地通过 **Retainer Tree** 看到是谁不放手，尤其是由于闭包（Closure）导致的隐性引用链路。

### 3. 一针见血的见解

内存控制能力是衡量顶级 JS 引擎的第二指标。V8 在这个阶段的努力，是通过更细致的 GC 启发式算法（Heuristics），让 JS 程序在高强度的 DOM 操作下，依然能保持平滑的分配曲线。
