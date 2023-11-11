CDQ 分治
todo https://oi-wiki.org/misc/cdq-divide/
推荐 https://blog.nowcoder.net/n/f44d4aada5a24f619442dd6ddffa7320
推荐 https://zhuanlan.zhihu.com/p/332996578
https://www.bilibili.com/video/BV1mC4y1s7ic
[学习笔记]CDQ 分治和整体二分 https://www.luogu.com.cn/blog/Owencodeisking/post-xue-xi-bi-ji-cdq-fen-zhi-hu-zheng-ti-er-fen
https://www.luogu.com.cn/blog/ljc20020730/cdq-fen-zhi-xue-xi-bi-ji
动态逆序对 https://www.luogu.com.cn/problem/P3157 https://www.luogu.com.cn/problem/UVA11990
CDQ 优化 DP https://www.luogu.com.cn/problem/P2487

---

流程：

1. 将问题转为数据结构问题，变为时间序列上的修改与查询
2. 分治

---

CDQ 优化 DP：先分治左半段，然后用这段值，双指针去更新 mid+1 到 r，再分治右半段

https://www.cnblogs.com/ydtz/p/16603498.html

---

CDQ 分治只能算作一种方法而非一种通用的算法，对于一段操作序列，我们从中间分开，先处理左边，再处理右边，最后加上左边对右边的影响。归并排序实际上就是一个 CDQ 分治。
https://www.cnblogs.com/TianMeng-hyl/p/14978786.html

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
