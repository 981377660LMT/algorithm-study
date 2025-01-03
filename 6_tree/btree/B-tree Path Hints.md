https://github.com/tidwall/btree/issues/7
https://github.com/tidwall/btree/blob/master/PATH_HINT.md
**B-tree Path Hints** 是一种针对 **B-Tree**（B 树）所做的**搜索与插入优化**技巧，通过**记录并重用上一次操作中定位到目标元素的“路径”**，来加速后续对相邻或相似键的访问。可以把它理解为一种“局部性优化”：`如果连续处理的键彼此在 B-Tree 中的位置比较接近，那么直接使用“上次找到的路径”去做下一次搜索，会比每次从头在节点里做二分查找更快。`

---

## 什么是 Path Hint

1. **正常的 B-Tree 搜索**

   - 在 B-Tree 的每个节点上，都需要对节点的 `keys` 做一次二分查找，确定下一个子节点索引；
   - 当节点中的元素数较多或访问的键分布在树的不同地方时，通常可以保持 **O(log N)** 的效率。

2. **Path Hint**：预定义/缓存的搜索“路径”

   - Path Hint 就是一条从根节点到目标元素所在节点（或应插入节点）的索引序列。
   - 在**下次操作**时，如果发现**要访问的键和上一次访问的键相近**，就先沿用这个路径“提示”来快速定位到大概的位置，而不是每次都从节点中点位置做二分搜索。

3. **发生偏差时的修正**
   - 如果这个“提示路径”完全正确，就能跳过节点内部的二分搜索，直接查到目标位置，性能可能得到几倍提升。
   - 如果提示路径错了，B-Tree 会**回退**到正常的二分搜索，并在完成后**更新**这条 Path Hint，使得下次能更准确。

---

## 为什么能提升性能

- **现实场景中键的“局部性”**

  - 例如：
    1. **时间序列（time series）**：插入的一批点往往落在同一个或相邻的区间；
    2. **区间插入**：在表中连续插入一批有序数据；
    3. **Redis 风格的键**：`"user:98512:name"`, `"user:98512:email"`，对同一“user”前缀的多个键进行操作。
  - 这些情况下，后一次要访问的键常和上一次相差不大，于是之前存下的“路径”也往往还很准确。

- **加速效果**
  - 当 Path Hint **完全正确**时，有机会减少在每个节点上做二分查找的开销，获得 **1.5 倍~3 倍**的速度提升（作者实测）。
  - 如果 Path Hint 完全失效，则最坏也只比正常二分搜索慢大约 5% 左右，并不会严重拖累性能。

---

## 使用方式

- **接口**：每次 B-Tree 的插入/查找/删除等操作，都可以带一个“Path Hint”参数进去；操作完成后，如果判定 hint 有偏差，就会改写 hint 的具体值，方便下次调用时更准确。
- **单线程**：往往就保留一个全局的 Path Hint，就可在每次操作时复用。
- **多线程**：每个线程可以持有一个“私有”的 Path Hint，避免并发数据竞争；或者每个客户端连接都保有一个 Path Hint。
- **维护成本很低**：B-Tree 会对传入的 hint 做校验和修正，一旦出错也会 fallback 到常规二分搜索。

---

## 小结

**B-tree Path Hints** 就是利用“访问局部性”来减少节点内部二分搜索成本的优化手段，具体做法是在每次操作结束后**记录从根节点到达目标位置的“路径”**，并在下次相邻或相近的操作中**直接重用**这条路径，避免全局二分搜索。这样做能在很多实际应用场景下带来可观的吞吐提升，且在最坏情况下也不会拖累整体性能。
