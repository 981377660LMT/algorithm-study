https://ottffyzy.github.io/algos/gt/johnson/

Floyd-Warshall 算法是一个更适合稠密图的算法，若设 V 表示图中点数 E 为图中边数，则复杂度为 O(V3) 是一个与边无关的量。

对于稀疏图，E 与 V 同阶的情形，此时我们发现如果图是非负权图，我们可以用以`每个点为起点的方式跑 V 次 Dijkstra 算法`，从而得到一个时间复杂度更好的算法
时间复杂度为 O(V((E+V)logE))=O(N2logN)

所以此时问题来了，Floyd-Warshall 是`适用于有负权图`的，那么对于稀疏带负权的图我们能得到一个更好的复杂度的算法吗？

## 使用带势能的 Dijkstra 战胜负权图

使用 Dijkstra 战胜负权图 AtCoder Beginner Contest 237 E - 严格鸽的文章 - 知乎
https://zhuanlan.zhihu.com/p/470948347
