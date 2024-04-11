区间数颜色

[区间颜色统计问题](https://taodaling.github.io/blog/2019/04/30/%E5%8C%BA%E9%97%B4%E6%93%8D%E4%BD%9C%E9%97%AE%E9%A2%98/)

https://www.luogu.com.cn/problem/SP3267
https://www.luogu.com.cn/problem/solution/P1972?page=2

- 离线

  1. 树状数组
  2. 莫队

- 在线
  1. 主席树

---

给定一个序列 a1,…,an，之后给出 q 个查询，第 i 个查询为 li,ri，询问 ali,…,ari 中有多少不同的数。
定义一个新序列 last1,…,lastn，**lasti 表示第 i 个数之前与 ai 相同的数中最大的下标**（如果不存在则设为 −1）。那么现在查询询问的实际上是区间 **lastli,…,lastri 中有多少数小于 li**。
