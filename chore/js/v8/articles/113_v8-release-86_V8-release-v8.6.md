# V8 v8.6 发布：底层算力优化与 Wasm SIMD

# V8 release v8.6: Underlying compute and Wasm SIMD

V8 v8.6 见证了性能从微观语法到宏观算力的多重演进。

### 1. 技术亮点：`Number.prototype.toString` 飞跃

针对小整数（Smi），V8 推出了基于 Torque 实现的快速路径。微基准测试显示，这一常用的高频操作性能提升了 **75%**。

### 2. 重磅特性：Liftoff SIMD

WebAssembly 的基线编译器 Liftoff 终于开始支持 SIMD 指令。这意味着 Wasm 应用在冷启动阶段即能利用现代处理器的向量运算能力，极大地缩短了需要进行复杂数学计算的应用（如视频编解码）的响应时间。

### 3. 一针见血的见解

v8.6 表明 V8 的性能优化已经深入到了“毫秒级”的基础设施中。无论是通过 Torque 重写核心库，还是为 Wasm 引入硬件加速，其核心目标都是将繁重的计算任务尽可能地推向底层和硬件，为现代 Web 带来的沉重算力负载提供更强劲的支持。
