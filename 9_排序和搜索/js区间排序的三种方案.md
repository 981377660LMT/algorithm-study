js 区间排序的三种方案

1. slice + sort，但是不是原地排序，排序完后需要再赋值回去.
2. subarray + sort，Uint32 数组保存元素用于排序的 key, key 必须是 uint32.
3. 采用 golang 1.20 的 stable sort.

TODO:
内存和时间的测试对比
