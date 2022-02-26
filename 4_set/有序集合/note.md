python 有序集合
http://www.grantjenks.com/docs/sortedcontainers/index.html

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

**python 的切片相比其他语言，是很快很快很快的**
https://stackoverflow.com/questions/12537716/why-is-slice-assignment-faster-than-list-insert
python3 数据结构
c++--unordered_map<> ---------- python3--dict()
c++--unordered_set<> ---------- python3--set()
c++--map<> --------------pytho3--sortedcontainers.SortedDict()------ collections.OrderedDict()
c++--set<> ---------------pytho3--sortedcontainers.SortedSet()
c++--multiset<>----------pytho3--sortedcontainers.SortedList()
对于区间判断是否重叠，我们可以反向判断，也可以正向判断。 暴力的方法是每次对所有的课程进行判断是否重叠，这种解法可以 AC。我们也可以进一步优化，
使用二叉查找树**二分查找**/来简化时间复杂度。
最后我们介绍了一种 **Count-Map 方法来通用解决所有的问题**，不仅可以完美解决这三道题，还可以扩展到《会议室》系列的两道题。

区间重叠判断的思路

1. 是否与一个区间交叉:[a,b]与[c,d] 等价于 a<=d&&b>=c
2. 是否与多个区间交叉:设 i = bisect.bisect_right(self.intervals, start)
   j = bisect.bisect_left(self.intervals, end)
   如果 i%2 为奇数 则说明 **start is in some stored interval**
   **新插入的区间起点与前面老的区间重叠了**
   如果 i!==j 则说明新的区间内包含了后面的老区间端点
   **新插入的区间终点与后面老的区间重叠了**
   `715. Range 模块`

区间技巧
使用两个 map 记录左右端点对应区间
`352. 将数据流变为多个不相交区间`
`352Number Stream to Intervals`
`855. 考场就座`

3. python 的 sortedList 本质是只有两层的 B-tree
   https://grantjenks.com/docs/sortedcontainers/implementation.html
