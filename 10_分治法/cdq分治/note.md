CDQ 分治（基于归并排序的分治）
todo https://oi-wiki.org/misc/cdq-divide/
推荐 **https://blog.nowcoder.net/n/f44d4aada5a24f619442dd6ddffa7320**
推荐 **https://zhuanlan.zhihu.com/p/332996578**
https://www.bilibili.com/video/BV1mC4y1s7ic
[学习笔记]CDQ 分治和整体二分 https://www.luogu.com.cn/blog/Owencodeisking/post-xue-xi-bi-ji-cdq-fen-zhi-hu-zheng-ti-er-fen

---

应用场景：

- 解决和点对有关的问题。

> P4602 [CTSC2018] 混合果汁
> https://www.luogu.com.cn/problem/P4602

- 1D 动态规划的优化与转移。（注意必须中序遍历:如果将 CDQ 分治的递归树看成一颗线段树，那么 CDQ 分治就是这个线段树的 中序遍历函数，因此我们相当于按顺序处理了所有的 DP 值，只是转移顺序被拆开了而已)

> cdq 分治+单调队列优化 dp
> P4269 [USACO18FEB] Snow Boots G
> https://www.luogu.com.cn/problem/P4269

- 通过 CDQ 分治，将一些动态问题转化为静态问题(时间序列分治,按照中序遍历序进行分治才能保证每一个修改都是严格按照时间顺序执行的)。
  静态问题：只查询，或者查询在所有更改之后的问题称为静态问题，其余则是动态问题

---

> CDQ 分治的意义
> 如果一道题只用到了一层 CDQ，那这道题也一定能用其他数据结构解决，所以一般大家都见不到这个算法。
> 不过不愿意写李超树和动态凸壳的人用 CDQ 分治来解动态凸壳的斜率优化问题倒是挺常见的。
> 由于 CDQ 分治本身不使用其他任何数据结构，这表示我们还留有底牌。即使现在在做的问题上再加一维，直接把线段树等结构套进去就能做。
> 说人话就是 CDQ 分治可以避免一部分的树套树，以及树套树套树。（还有一部分可以使用整体二分的思路去避免高级数据结构堆叠）
> 个人觉得分治类算法无论是泛用性还是上限都很高，在允许暴力枚举当前分治的区间的前提下复杂度是优秀的 nlogn，如果题目复杂需要嵌套其他数据结构时，也比纯数据结构套数据结构来的更容易，但是在算法竞赛中又往往被忽视（一般遇到一个题先想用个什么数据结构去维护）。
> 程序=算法+数据结构，而不是数据结构的嵌套和堆叠。

---

TODO，很多没懂
