# 一般矩阵乘法

- 一个古老而有名的理论算法是复杂度`O(n^log2(7))`的[Straseen](https://zhuanlan.zhihu.com/p/268392799)算法，基于分治的思想

- 一般的矩阵乘法最快是 `O(n^2.37)`的
  https://link.zhihu.com/?target=https%3A//arxiv.org/pdf/2210.10173.pdf

- 对稀疏矩阵、非方阵情况的矩阵乘法有一些更快的做法
- 矩阵乘法算法的实际速度表现依赖于硬件。在算法竞赛的场景中无法使用 GPU，所以我们主要讨论在普通 PC 上使用 CPU 的实现

- 一般用一维数组+交换变量+循环展开(loop unrolling)等方法优化

---

TODO
https://judge.yosupo.jp/problem/matrix_product

实际工程场景很少见 01 矩阵，F2 矩阵也不常见，大部分都是浮点矩阵，一般都是 gpu 算
