<!-- https://github.dev/EndlessCheng/codeforces-go/blob/master/copypasta/leftist_tree.go#L11 -->
<!-- https://nyaannyaan.github.io/library/data-structure/skew-heap.hpp -->
<!-- https://www.cnblogs.com/flashhu/p/8324551.html -->
<!-- https://www.luogu.com.cn/blog/command-block/lct-xiao-ji -->

- 斐波那契堆,用于高速化 prim 最小生成树算法
- radixHeap,用于高速化 dijkstra 单源最短路算法
- 可持久化左偏树,用于求 k 短路
- skewHeap, 用于求有向图的最小生成树(skewHeap 可以给堆全体加上一个数)

---

模板题
https://www.luogu.com.cn/problem/P3377
https://www.luogu.com.cn/problem/P2713

启发式合并的可并堆 => O(logn^2)
skewHeap => O(logn)
