以下内容是对 **“Golang 内存模型与分配机制”** 一文的一个相对系统、详细的总结梳理。文中核心关注点在于 Golang 运行时在堆上如何分配对象，以及为保证高并发、低锁争用而采用的多级缓存与不同大小规格（spanClass）相结合的思路。希望能帮助你快速梳理 Go 内存分配背后的关键概念与源码脉络。

---

## 1. 前言

- **内存管理** 与 **垃圾回收** 属于 Go 语言最复杂的运行时模块之一。本文首先聚焦于 **Golang 的内存模型与分配机制**，其后再衍生到 GC 话题。
- Go 借鉴了操作系统“分页管理”的思路，以及“多级缓存 + 不同大小规格”的设计，并通过无锁或细粒度锁机制，尽可能提高分配性能。
- 贯穿本文的主线：**Go 对内存分配的关键是如何以最小代价，在高并发情形下高效完成对象分配**。

---

## 2. Golang 内存模型：核心概念

Go 运行时从操作系统获取大块内存后，会将其作为**堆 (mheap)** 的抽象，然后再细分、切割、缓存，形成一个多级结构：

1. **mheap**（全局堆）

   - 以 **页**（8KB）为最小分配粒度。
   - 维护全局空闲页信息，通过 **pageAlloc**（基数树）快速找到连续空闲页，再把页拼装成 `mspan`。
   - 以 **heapArena** 为分配单位（64MB 一块）向操作系统申请内存。
   - 持有所有的 **mcentral**（每种 `spanClass` 对应一个）。

2. **mcentral**（中心缓存）

   - 每个 `mcentral` 专门管理**一种 `spanClass`**（对象大小 + 是否含指针）。
   - 聚合了该类所有的 `mspan`，并以“可用（partial）”与“满（full）”链表区分是否仍能分配对象。
   - 加锁粒度只限于该 `spanClass`；有利于减小锁竞争。

3. **mcache**（线程缓存，准确说是每个 P 缓存）

   - 每个 P（GMP 模型中的 P）独享一份 mcache，无锁操作。
   - 其中有 `alloc[spanClass]` 数组，存放对应 `spanClass` 的可分配 mspan。
   - 如果本地 mspan 用完，就向 mcentral 要一个新的；若 mcentral 也没有则会向 mheap 要。
   - **tiny allocator**：专门处理 16B 以下且无指针的微小对象，加速分配。

4. **mspan**（最小管理单元）

   - mspan 大小是整页倍数（8KB 的整数倍），内部被切分成等大小的小块（object 大小 = `bytes/obj`）。
   - 同一个 mspan 只管理一种大小规格（一个 `spanClass`）的对象，因此内聚性较好。
   - 内部通过 `allocCache` + `bitmap` 标记哪些块已被分配。分配时可通过 `Ctz64` 在位图上找到首个可用块。

5. **spanClass**（大小规格 + 是否含指针）
   - Go 预先定义了 **67 种常见大小规格**（从 8B ~ 32KB），再加上一位 `noscan` 标记是否包含指针，共计 136 种。
   - 用于把对象的大小映射到对应的 `mspan`，避免产生大量外部碎片。
   - 对象实际分配时会向上取整到最近的规格；这样会产生一定内部碎片。

### 2.1 多级缓存与无锁 / 细锁化

- 顶层 `mheap` 持有全局锁。频繁使用会导致严重性能问题。
- 中间层 `mcentral`：每种大小规格一个互斥锁，即**粒度为“对象规格”**。
- 最外层 `mcache`：每个 P 一份，无需锁。典型情况下，对象分配只需操作 `mcache`。
- 如果本地缓存不够，才会“逐层退化”到 `mcentral` 或 `mheap` 去获取更大块内存。

### 2.2 PageAlloc 与 HeapArena

- **pageAlloc**：Go 使用基数树（radix tree）来管理空闲页。
  - 每棵基数树对应 16GB 的空间；Go 最多支持 2^14 棵树，可索引到 256TB。
  - 树节点存储“前缀连续空闲页”、“后缀连续空闲页”、“最大连续空闲页”等信息，用于自顶向下快速查找。
- **heapArena**：每块 64MB，是向操作系统申请的一大段地址空间。
  - 里面记录 `[8192]*mspan` 数组，便于从地址反查该页隶属哪个 mspan，GC 等场景会用到。

---

## 3. 对象分配流程

Go 在编译期或运行期判断对象大小，分为：

1. **微对象**：< 16B 且无指针（`noscan`）。
2. **小对象**：≤ 32KB（大部分场景适用，包含有指针或无指针）
3. **大对象**：> 32KB。

总体分配过程一般是“多级缓存”模式，自顶向下查找空间，但只要在某一级命中就直接返回，无需继续走下去。

### 3.1 `mallocgc` 主流程概览

```go
func mallocgc(size uintptr, typ *_type, needzero bool) unsafe.Pointer {
    // 1. 获取当前 M、P、mcache
    // 2. 判断是否 noscan
    // 3. 根据对象大小分类：
    //    - 微对象 (<16B && noscan) ：尝试走 tiny allocator
    //    - 小对象 (<=32KB) ：映射到对应 spanClass，再从 mcache alloc/ nextFree
    //    - 大对象 (>32KB) ：直接从 mheap 分配
    // 4. 若 mcache 拿不到对象，则向 mcentral 拿，如果 mcentral 也满了，则向 mheap 要页来组装 mspan
    // 5. 若 mheap 也无可用页，则向操作系统申请（mmap）
    // 6. 返回分配的地址
}
```

#### (1) 微对象分配：`tiny allocator`

- 仅适用于“无指针 + size < 16B”。
- mcache 里有一个 `tiny` 块（16B），用 `tinyoffset` 线性分配；若够用直接返回，否则分配一个新的 tiny 块。

```go
if noscan && size < maxTinySize {
    off := c.tinyoffset
    ...
    if off+size <= maxTinySize && c.tiny != 0 {
        x = unsafe.Pointer(c.tiny + off)
        c.tinyoffset = off + size
        ...
        return x
    }
    // 分配新的 tiny 块 ...
}
```

#### (2) 小对象分配（≤ 32KB）

- 根据大小找对应 `spanClass`（0~66，再加 1 bit 表示 noscan），从 `mcache.alloc[spc]` 拿一个 mspan，调用 `nextFreeFast / nextFree` 找可用块。
- 若拿不到，则 mcache.refill -> 从 mcentral.cacheSpan 获取 mspan；
- mcentral 如果也无可用 mspan，就 grow -> `mheap.allocSpan`；
- mheap 若无空闲页，就 `sysAlloc` 向 OS 申请。

#### (3) 大对象分配（> 32KB）

- 跳过 tiny 与小对象逻辑，直接从 mheap 要空闲页组出 mspan 并返回。

### 3.2 mcache 分配：`nextFreeFast`

- mspan 内部维持了 `allocCache`（一个 64 位 bitmap）和 `freeindex`。
- 调用 `Ctz64` 找到最低位的空闲 block，做分配后更新 bitmap。

```go
func nextFreeFast(s *mspan) gclinkptr {
    theBit := sys.Ctz64(s.allocCache)
    if theBit < 64 {
        result := s.freeindex + uintptr(theBit)
        ...
        s.allocCache >>= (theBit + 1)
        s.freeindex = freeidx
        return gclinkptr(result*s.elemsize + s.base())
    }
    return 0
}
```

### 3.3 mcentral 分配：`cacheSpan`

- 若 mcache 没有空位，就向 `mcentral` 对应的 partial/full 链表找 mspan。
- 若 partial/full 都没可用，则调用 `grow()` 从 mheap 拿空闲页组装一个新的 mspan。

```go
func (c *mcentral) cacheSpan() *mspan {
    // 尝试从 partialUnswept / fullUnswept 获取
    // ...
    // 找不到就 c.grow() -> mheap.alloc -> allocSpan
    // 拿到新 mspan 后返回
}
```

### 3.4 mheap 分配：`allocSpan`

- 在加锁(`sched.lock`)下，通过 `pageAlloc.alloc(npages)` 找到连续空闲页基地址 `base`。
- 组装成 mspan，记录到 `heapArena.spans`。
- 如果 `pageAlloc` 也没找到足够页，就去 `sysAlloc -> mmap(...)` 从操作系统申请更多内存。

### 3.5 向操作系统申请：`sysAlloc`

- `mheap.sysAlloc` 最终调用 `mmap` 或等价系统调用获得新的地址空间，并将其纳入 `pageAlloc` 的基数树索引。

---

## 4. 小结与展望

1. **多级缓存**：Go 在分配时从“最内层”的 `mcache` 开始，无锁快速分配；若本地不够，再逐级退化到 `mcentral`（粒度锁）或 `mheap`（全局锁）。大幅减少锁争用。
2. **按大小规格划分**：通过 `spanClass` 事先把对象大小分为 ~67 个区间，内部按固定块分配，减少碎片。
3. **分层管理**：mheap 担任全局内存的“源头”，基于 8KB 页 + 基数树 `pageAlloc` 查找空闲区块；再配合 `heapArena` 维护页与 mspan 间映射。
4. **对象分配速度**：在多数场景下只需操作 `mcache`，无需加锁，非常高效。
5. **垃圾回收**：在分配基础之上，Go 再采用“三色标记 + 混合写屏障” 实现并发 GC，回收不再使用的对象所占内存。此部分将是下一步深入的话题。

这正是 Go 运行时为适应高并发和降低停顿而进行的精妙设计：**把复杂性尽量压在多级结构上，给常规分配提供最快路径**。在理解了内存分配机制后，后续再看 GC 时，“为什么 GC 要关注 mspan”“如何快速从对象地址找到 mspan” 等，就会有更直观的答案。
