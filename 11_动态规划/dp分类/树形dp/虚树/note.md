https://oi-wiki.org/graph/virtual-tree/
https://beet-aizu.github.io/library/tree/auxiliarytree.cpp
https://www.cnblogs.com/ydtz/p/16276275.html
https://smijake3.hatenablog.com/entry/2019/09/15/200200
https://blog.sengxian.com/algorithms/virtual-tree

- 有的时候，题目给你一棵树，然后有 q 组询问，每组询问为与树上有关的 k 个点的相关问题，
  明显的暗示: ∑k<=5e5
  每次查询不能 dp 整棵树，因为复杂度太高，所以要用虚树来优化。
  也就是用给定的点集和他们的 lca 来构造一棵新的树，然后在这棵新树上 dp。
  这颗新树的结点个数最多为 2n-1，所以所有查询总复杂度为 O(nlogn)。

- 常用技巧
  1. LCA + 倍增查询压缩后边权最大/最小值
  2. 每次查询前将点集打标记 visited,查询完后再清除标记
  3. 将根节点加入点集后,建立出来的虚树一定包含根节点，方便 dp

---

Auxiliary Tree 虚树
https://oi-wiki.org/graph/virtual-tree/
https://cmwqf.github.io/2020/04/17/%E6%B5%85%E8%B0%88%E8%99%9A%E6%A0%91/
https://tjkendev.github.io/procon-library/python/graph/auxiliary_tree.html
指定された頂点たちの最小共通祖先関係を保って木を圧縮してできる補助的な木
!有的时候，题目给你一棵树，然后有 q 组询问，每组询问为与树上有关的 k 个点的相关问题，
而这个问题如果只有一个询问的话，可以用树形 dp 轻松解决
那么我们可以考虑每次只在这 k 个点及相关的点构成的树上进行 dp
往往需要虚树上进行树形 DP
!把一些有用的点给拿出来，然后通过最少的 lca 把这些节点给穿到一起，在新的树上做树形 dp。
点集个数为 k 时,最多 2\*k-1 个顶点
