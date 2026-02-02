# 低延迟、高吞吐量的垃圾回收：Orinoco 项目

# Low-latency, high-throughput garbage collection

Orinoco 项目是 V8 现代垃圾回收器的基石，旨在通过并行（Parallel）和并发（Concurrent）技术消灭 Jank（卡顿）。

### 1. 并发标记 (Concurrent Marking)

最沉重的任务是标记数以百万计的对象。Orinoco 让标记程序在后台线程中运行。当 JS 主线程修改对象图时，它会触发“写屏障”（Write Barriers）来通知后台标记器，从而在不停止执行的情况下完成大部分标记工作。

### 2. 并行整理 (Parallel Compaction)

在清理阶段，传统的整理算法只有主线程在工作。Orinoco 引入了并行化，让多个线程同时负责搬运对象。为了解决竞争，V8 巧妙地使用了基于页面的 Bitmap 来追踪指针变更，使得数以万计的对象搬迁能够安全高效地并发进行。

### 3. 一针见血的见解

Orinoco 的核心理念是“让主线程只做 JS 的事”。通过将耗时的标记和移动操作异步化，V8 把 10ms 以上的长停顿（STW）切分成了无数个 <1ms 的极小片段，这是现代流畅 Web 体验的生命线。
