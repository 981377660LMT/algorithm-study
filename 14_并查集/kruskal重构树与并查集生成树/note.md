# 从树的角度分析并查集的合并过程

`Kruskal重构树（ProcessOfMergingTree)`
定义：第 i 次操作 a_i 和 b_i 分属于两个不同子树，那么 Kruskal 会新建一个结点 u，然后让 a_i 所在子树的根和 b_i 所在子树的根分别连向 u，作为 u 的两个儿子。

- 每个点对应子树里都是边长小于等于其的点权的联通块
- 每个节点的权值肯定大于等于其子树中任意一个点的权值，因为构造最小生成树的时候越大的边出现的越晚，对应建立的点深度也越浅

https://www.luogu.com.cn/blog/user9012/ke-lu-si-ka-er-zhong-gou-shu-lve-xie

---

https://atcoder.jp/contests/agc002/tasks/agc002_d

https://www.mathenachia.blog/agc002d-usereditorial/#toc1

解法:

1. 整体二分
2. 部分可持久化并查集
3. 操作分块+重构(或者二进制分组)
4. disjoint set union (DSU) のマージ過程を表す木

---

Kruskal 重构树 与 并查集生成树(森林)本质一样
