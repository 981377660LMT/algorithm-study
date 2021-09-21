**平衡树(set/map)数据结构**

一般指 treeSet(Java)/sortedList(python)/multiset(C++ STL)
很多情况下，线段树的一个简单替代方案是二分法维护有序集合

解决:
**数据流区间合并问题**
`352. 将数据流变为多个不相交区间`
`352Number Stream to Intervals`
**天际线问题**
`699. 掉落的方块`
`218. 天际线问题`

Java
`218. 天际线问题`中
优先队列的 remove 操作成为了瓶颈，如何优化？

由于优先队列的 remove 操作需要先经过 O(n)的复杂度进行查找，再通过 O(logn) 的复杂度进行删除。因此整个 remove 操作的复杂度是 O(n) 的，这导致了我们算法整体复杂度为 O(n^2)。

优化方式包括：

1. 使用基于红黑树的 TreeMap 代替优先队列；
2. 或是使用「哈希表」记录「执行了删除操作的高度」及「删除次数」，在每次使用前先检查堆顶高度是否已经被标记删除，如果是则进行 poll 操作，并更新删除次数，直到遇到一个没被删除的堆顶高度。(lazy deletion)
   **「力扣」第 480 题：滑动窗口中位数 的 官方题解。**
3. 使用哈希堆

作者：AC_OIer
链接：https://leetcode-cn.com/problems/the-skyline-problem/solution/gong-shui-san-xie-sao-miao-xian-suan-fa-0z6xc/
来源：力扣（LeetCode）
著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。
