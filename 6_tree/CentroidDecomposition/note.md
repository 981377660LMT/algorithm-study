## 重心分解+点分治 处理大规模的树上`路径`问题

1. 暴力遍历路径 `O(n^2)`
   枚举顶点，对每个点出发 dfs
2. 分治+递归 `O(nlogn)`
   对每个点，考虑包含这个点的路径和不包含这个点的路径

   - **包含这个点的路径**，可以用 dfs 求出
   - **不包含这个点的路径**，删除这个点，然后对每个子树递归求解

   在单链时，直接枚举点会退化成 O(n^2)
   **重心分解**：如果选取一个 mid 值(树的重心)作为子树的根，那么每个递归问题中子树的大小都不超过整棵树的一半，所以每次递归问题规模可以下降至一半或更低，从而将时间复杂度降到 O(nlogn)

   > 应用:
   > https://zhuanlan.zhihu.com/p/359209926
   > 求无根树中长度为 k 的路径数目
   > 树上距离不超过 upper 的点对数
   > 求最长的 gcd 大于 1 的路径。
   > ...

   我们只需要关注`子树中包含当前重心的路径`,通常需要 容斥原理/卷积 来计算子树的贡献

## 在使用重心分解前，想想树形 dp/换根 dp 能不能做

## 点分治的核心其实就是树的重心，如果了解了树的重心的做法其实就知道怎么做了

「部分木->オイラーツアー, パス->HL 分解, 同心円状->重心分解」

## 构造一棵树，满足点分树的两条性质

https://atcoder.jp/contests/abc291/tasks/abc291_h

重心的性质:

- 以树的重心为根时，所有子树的大小都不超过整棵树大小的一半
- 树中所有点到某个点的距离和中，到重心的距离和是最小的；如果有两个重心，那么到它们的距离和一样。
- 把两棵树通过一条边相连得到一棵新的树，那么新的树的重心在连接原来两棵树的重心的路径上。
- 在一棵树上添加或删除一个叶子，那么它的重心最多只移动一条边的距离。

## 实装

https://ei1333.github.io/library/test/verify/yukicoder-1002.test.cpp
https://twitter.com/tatyam_prime/status/1629903556856643585?s=20

业务逻辑：

- removed[i]: 顶点 i 是否被删除
- decomposition(i): 业务逻辑函数, 内部先调用 getCentroid(i) 得到重心, 然后求出通过重心的答案，再删除重心，然后对每个子树递归 decomposition 求解

求点分树(点分树)：

- getSize(i): 以 i 为根的子树的大小 (忽略 removed[i] 为 true 的顶点 dfs)
- getCentroid(i): 以 i 为根的子树的重心 (忽略 removed[i] 为 true 的顶点 dfs)
- build(i)： 构建以 i 为根的子树的点分树 (忽略 removed[i] 为 true 的顶点 dfs)

- 点分树最高 logn 层
  点分树的题一般是给定一个范围[l,r]，统计与某个结点距离[l,r]内的结点某个量的和
  https://judge.yosupo.jp/problem/vertex_add_range_contour_sum_on_tree

---

No.1002 Twotone-两种颜色的路径数
https://yukicoder.me/submissions/449141
