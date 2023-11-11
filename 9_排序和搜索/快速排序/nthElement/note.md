https://cp-algorithms.com/sequences/k-th.html

---

median of medians 是一种中位数（近似）选取算法，常用于其他选择算法中（主要是 QuickSelect 算法中）进行 pivot 元素的选取。

在如 QuickSelect 这样的选择算法中，我们 pivot 元素的选取对于我们算法的效率有很大的影响，median of medians 算法可以帮助我们在线性时间内选出中位数附近的元素

Median of medians 算法也被称作 BFPRT 算法
[BFPRT——Top k 问题的终极解法](https://zhuanlan.zhihu.com/p/291206708)
[你所不了解的算法——BFPRT 算法](https://www.luogu.com.cn/blog/448887/ni-suo-fou-liao-xie-di-suan-fa-bfprt-suan-fa)
