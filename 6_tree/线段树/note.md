分享 | 线段树的模板总结 (typescript) (更新中...）
https://leetcode.cn/circle/discuss/4rJDBt/#typescript-%E6%A8%A1%E6%9D%BF%E5%8F%8A%E4%BD%BF%E7%94%A8%E6%A1%88%E4%BE%8B

## 我对线段树的理解:

0. 用来维护幺半群(monoid)的区间信息
1. 2 个类型、5 个操作

- Data 线段树结点维护的信息
- Lazy 更新的值的类型

- e 区间信息的幺元
- id 更新的幺元(懒标记幺元)
- op 区间信息的合并函数 (从叶子结点往上 **pushUp** 时调用)
- mapping 懒标记更新区间信息 (从根结点往下 **pushDown** 时调用)
- composition 懒标记合并函数 (从根结点往下 **pushDown** 时调用)

一般会从 **pushDown** 中抽离出 **propagate** 函数复用

2. 递归线段树的写法
   [我对线段树的理解](./%E6%88%91%E5%AF%B9%E7%BA%BF%E6%AE%B5%E6%A0%91%E7%9A%84%E7%90%86%E8%A7%A3.ts)

## 碎碎念

- `如果没有区间修改，没有必要用pushDown+懒标记，因为每次都会修改到叶子结点`
- `如果只有单点查询，没有必要用pushUp(op合并结点信息)，因为信息全在叶子节点`
-
- 写线段树的话，推荐使用 `ts+数组实现(节点的写法很容易MLE)+考虑大数`
- 注意一般的线段树下标(索引)从 1 开始(同树状数组)，动态开点线段树下标(值域)可以从 0 开始

---

单点更新非递归线段树的简易写法
https://hackmd.io/@tatyam-prime/rkA5wJMdo
