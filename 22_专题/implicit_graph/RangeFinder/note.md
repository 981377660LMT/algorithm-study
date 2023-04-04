## Finder(RangeFinder、OnlineFinder) 这类数据结构通常配合在线算法使用。

- API:

1. prev(x): 返回数据结构中小于 x 的最大的未被删除的元素，如果不存在则返回 None.
2. next(x): 返回数据结构中大于 x 的最小的未被删除的元素，如果不存在则返回 None.
3. erase(x): 删除数据结构中的 x，如果 x 不存在则不做任何操作。

4. 可选：erase(left,right): 删除数据结构中的 [left,right) 区间的元素，如果区间为空则不做任何操作。

- 使用:
  在线算法中:
  `setUsed`中调用 erase 函数
  `findUnUsed`中调用 next/prev 函数
- 可用的数据结构
  并查集
  有序集合/set
  fastset
  链表
