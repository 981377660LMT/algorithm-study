**应用于离线查询**

https://zhuanlan.zhihu.com/p/25017840
莫队的精髓就在于，离线得到了一堆需要处理的区间后，`合理的安排这些区间计算的次序`以得到一个较优的复杂度。

原来：两点之间的转移开销为两点之间的曼哈顿距离
连接所有点的最优方案为一棵树，那么整体的时间复杂度就是这棵树上所有曼哈顿距离之和。
于是乎最优的复杂度肯定是这棵树是最小生成树的时候，也就是曼哈顿距离最小生成树。但这么打貌似代码复杂度有点大。。而且在实际的转移中肯定会出现分支，需要建边

解决：以`询问左端点所在的分块的序号为第一关键字，右端点的大小为第二关键字进行排序`，按照排序好的顺序计算，复杂度就会大大降低

**什么样的离线查询可以用莫队算法?**
如果已知区间[left,right]的答案 并且由此可以推出区间[left±1,right±1]的答案 那么就可以用莫队算法

https://ei1333.hateblo.jp/entry/2017/09/11/211011

---

二次离线莫队 将 `O(nsqrt(n)logn)` 优化到 `O(nsqrt(n))`
https://www.cnblogs.com/jz-597/p/13598510.html
