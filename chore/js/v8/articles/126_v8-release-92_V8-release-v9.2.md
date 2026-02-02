# V8 v9.2 发布：Atomics.waitAsync 与共享指针压缩

# V8 release v9.2: Atomics.waitAsync and Shared Pointer Compression Cage

本版本的核心架构改进是“共享指针压缩笼”。

### 1. 深度解析：Shared Pointer Compression Cage (共享指针压缩笼)

V8 的指针压缩技术（Pointer Compression）通过只存储 32 位偏移量来节省内存。

- **旧模式**：此前，每个 Isolate（独立实例）都有自己独立的 4GB 内存布局（Cage）。
- **新架构**：从 9.2 开始，同一个进程内的所有 Isolate 可以共享同一个 4GB 的 Cage。
- **意义**：这不仅减少了虚拟内存占用的碎片化，更为 **Worker 间通信** 和 **Shared Memory** 提供了更高效的基础，因为 32 位压缩指针现在可以在不同线程间直接通用而无需昂贵的地址补码转换。

### 2. 标准落地：`Atomics.waitAsync`

引入了异步等待的原生支持，解决了在主线程中无法同步等待原子变量的尴尬局面，极大地增强了 JS 多线程协作的灵活性。

### 3. 一针见血的见解

共享指针压缩笼（Shared Cage）是 V8 向现代多线程 JS runtime 演进的必然结果。它打破了“隔离实例”之间的内存物理屏障，通过架构层面的统一，让多核心、多 Isolate 环境下的内存足迹（Footprint）得到了质的优化。
