RMQ (Range Minimum/Maximum Query)：
对于长度为 n 的数组 A，回答若干询问 RMQ(A,i,j)(i,j<=n-1)，返回数组 A 中下标在 i,j 范围内的最小（大）值，也就是说，RMQ 问题是指求区间最值的问题。
关于 RMQ 问题的几种解法

1. 暴力法最简单的方法，就是遍历数组直接搜索，但是这种方式时间复杂度是 O(n)。对于数组长度较大，性能要求高的场景不适用。一般用这个算法就等着 TLE，时间复杂度最坏 O（Q\*N），也不一定超时，签到题可能就直接让你过了。

2. **ST（Sparse Table）算法**
   ST 算法是一种更加高效的算法，基于动态规划的思想，以 O(nlogn)的预处理代价，换取 O(1)的查询时性能。但是，是离线的，也就是说每次修改都是 O(nlogn)复杂度，那么用在带修的题目上就显得捉襟见肘了。

3. **线段树**是基于分治的思想来实现的，建立是 o（nlogn）查询为 O（logN),那么也就是说这个可以进行修改，单点修改维护也是 logN。

分析也就是说，我们可以抛开 1/3 不谈，当题目是`离线的时侯使用 ST 算法更快，当题目是在线的时候直接使用线段树维护即可，`

4. **查询时固定区间大小:滑动窗口的最值**可以用单调队列 O(1)

注意:`RMQ 问题可与 LCA 转化` 建立笛卡尔树 只需要 O(n)复杂度离线处理

**注意:当数组长度足够长时，此时使用 Uint32Array 才会比普通数组快(结论来源于 st 表的性能测试)**

---

O(n)构建 O(1)查询的 RMQ

https://hotman78.github.io/cpplib/data_structure/RMQ.hpp
https://noshi91.hatenablog.com/entry/2018/08/16/125415
https://hotman78.github.io/cpplib/data_structure/arg_rmq.hpp

`+-1 RMQ (Plus-Minus-One-RMQ)`
https://ei1333.github.io/library/structure/others/plus-minus-one-rmq.hpp
`用欧拉序和+-1RMQ O(n)构建 O(1)求 LCA`
https://ei1333.github.io/library/graph/tree/pmormq-lowest-common-ancestor.hpp

https://ljt12138.blog.uoj.ac/blog/4874
https://www.cnblogs.com/bzy-blog/p/14353073.html
