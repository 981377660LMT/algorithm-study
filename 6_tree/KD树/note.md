[kd 树算法之思路篇](https://www.joinquant.com/view/community/detail/dd60bd4e89761b916fe36dc4d14bb272)
每一个节点都是按照上下或者左右进行平分的，因此如果两个点在树中的距离较近，那么它们的实际距离就是比较近的。
kd 树是一个二叉树结构，它的每一个节点记载了【特征坐标，切分轴，指向左枝的指针，指向右枝的指针】

- Python 的 scikit-learn 机器学习包提供了蛮算、kd 树和 ball 树三种 kNN 算法
- 给定一堆已有的样本数据，和一个被询问的数据点（红色五角星），我们如何找到离五角星最近的 15 个点？
  https://leetcode-cn.com/problems/vFjcfV/solution/kd-treeban-zi-ti-by-mo-yan-24-63mv/

- golang 2 维 k-d 树实现
  **KDT 的核心思想就是对矩形的水平分割。**
  每棵子树都代表的是一个矩形；如果某棵树不平衡，就暴力重构。
  https://github.dev/EndlessCheng/codeforces-go/blob/master/misc/atcoder/abc274/e
  https://www.luogu.com.cn/blog/lc-2018-Canton/solution-p4148
  kd 树维护二维平面上的点集，查询到给定点的最近的距离。

- 博客
  https://www.luogu.com.cn/blog/lc-2018-Canton/solution-p4148
  https://oi-wiki.org/ds/kdt/
  https://github.dev/EndlessCheng/codeforces-go/tree/master/copypasta

---

Binary space partitioning Tree， BSPTree
二叉空间分割

---

KD 树像线段树一样，可以维护二维幺半群
[单点修改区间查询](https://maspypy.github.io/library/ds/kdtree/kdtree_monoid.hpp)
[区间修改单点查询](https://maspypy.github.io/library/ds/kdtree/dual_kdtree_monoid.hpp)
[区间修改区间查询](https://maspypy.github.io/library/ds/kdtree/kdtree_acted_monoid.hpp)

---

https://github.com/spaghetti-source/algorithm/tree/master/data_structure

---

https://taodaling.github.io/blog/2019/04/19/KD-Tree/
KDTree 意思是 K Dimension Tree，即 K 维树，其可以用于存储 K 维空间中的点。

实现非常简单，每个结点有左右子结点，如果结点处于第 i
层，那么其左子树存储的是第 i 维值小于该结点的结点，因此每一层的比较规则都不同，k 层一循环。

KDTree 支持插入和删除，与其它二分搜索平衡树区别不大。如果用已知结 n
个结点建树，可以利用第 k 大元素算法在 O(nlog2n)的时间复杂度内建立完成。但是这样没法保证修改后依旧平衡，可以利用替罪羊树的思路，提供一个平衡因子，在**失去平衡后暴力重建子树**。这样插入和删除的时间复杂度摊还为 O(n(log2n)2)。

KDTree 支持很多操作，比较典型的就是区间查询和最近点查询。
区间查询很简单，和线段树差不多，记住要在区间无相加时剪枝。一次查询时间复杂度为 O(2kn^(1−1/k))。

最近点查询的时间复杂度与距离计算的算法有关，但是无论哪种时间复杂度都很不好算（捂脸）。有两种可行的优化，一种是每次进入左右子结点之前，先预估哪边可能最近的顶点，并选择拥有最近顶点的那个分支，这样一般效果不错。（说白了就是剪枝呗）
