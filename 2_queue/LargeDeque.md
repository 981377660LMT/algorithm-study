- 除了目前的两端批量增删、随机访问、区间求和、区间插入删除外，还能在 O(log B)（B=块数）时间内支持以下常见操作，只需在每个节点上维护额外信息或懒标记，并在 split/merge 及 update 时同步处理即可：

  1. 区间加值（Range Add）  
     – 节点加 `add` 懒标记，`sum += add*size`，下推时子节点继承。

  2. 区间赋值（Range Assign）  
     – 节点加 `assign` 懒标记，重置 value、sum=value\*size，并覆盖子懒标记。

  3. 区间翻转（Reverse）  
     – 节点加 `rev` 懒标记，交换 `left/right`，下推时子节点反转标记异或。

  4. 区间最值（Range Min/Max Query）  
     – 每节点维护 `min`/`max`，合并时 `min=min(l.min, value, r.min)` 等。

  5. 区间乘积、区间 GCD、区间哈希 等  
     – 类似维护 `prod`、`gcd`、`hash` 字段与合并逻辑。

  6. 按值分裂/合并（Split/Merge by Key）  
     – 不按位置 split，而按节点值或某个比较函数 split。

  7. 前缀和二分（MaxRight/MinLeft）  
     – 给定累积函数和阈值，二分查找最早／最晚满足条件的位置。

  8. Order‐Statistic：查询区间 k-th 最小／第几大的元素  
     – 维护子树中各值的计数或结点数，split by rank。

- 只要在 `Node` 增加对应字段＋懒标记，并在 update、`pushDown`、split、merge 里同步维护，就能统一支持上述所有高级区间操作。
