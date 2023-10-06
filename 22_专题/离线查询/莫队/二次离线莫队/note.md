## 二次离线莫队

- 解决的问题

一般莫队有 O(n√n)次端点移动，如果要用数据结构维护信息的话，就有 o(n√n)次修改和 O(n√n)次查询。
而莫队二次离线能够优化为成 O(n) 次修改和 O(n√n)次查询，从而允许使用一些`修改复杂度大而查询复杂度小`的方式来维护信息。

- 使用条件

1. 每个数 A[i]对答案的贡献是可交换群，即`f(x,start,end) = f(x,0,end) - f(x,0,start)`
2. 可以 `O(1) 查询每个数的贡献`, `O(sqrt(n))` 添加一个新元素

---

https://olddrivertree.blog.uoj.ac/blog/4656
三道经典分块题的更优复杂度解法&[Ynoi2019 模拟赛]题解

1. 强制在线的区间逆序对
   Yuno loves sqrt technology I
   https://www.luogu.com.cn/problem/P5046
2. 允许离线的区间逆序对 -> 二次离线莫队
3. 强制在线的区间众数查询
