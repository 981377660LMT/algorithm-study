# 版本控制

1. 记录状态(state,持久化)

   - 不能全量保存,需要 `structual sharing` 的可持久化数据结构（Persistent data structure）
   - **每次修改历史版本总会返回一个新的数据结构(immutable ,不可变对象)**
   - 一般使用树型结构实现(包含链表,非线性结构)
   - 持久化数据结构可以带来许多好处，比如异常安全（Exception Safety）和并发性（Concurrency）。

2. 记录变化(action/operation)
   - mutation/action
   - undo redo 的对顶栈
   - 一些部分可持久化数据结构

## 完全可持久化栈

- 链表实现 (git)
- Path copying (路径复制,将即将被修改的点路径上的所有节点克隆出一个新的。这些修改必须通过数据结构逐级连接)

## 可持久化线段树(静态)

静态查询区间第 k 小

## 完全可持久化数组

### 完全可持久化数组有**在线离线**两种维护方法

N:数组长度 M:更新次数 Q:查询次数

- 离线:

  - 预处理查询+dfs(前提是操作可逆) `O(N+M+Q)`

- 在线:

  - full backup (Copy on Write 写入时复制,使用复制整个数据结构的方式来记录每次更改) `O(N*M+Q)`
  - 状态复元(前提是操作可逆) `O(N+M*Q)`
  - 数据结构 `O((Q+M)*logN)`

### 完全可持久化数组的 api:

1. 在某个历史版本上修改某一个位置上的值(生成一个完全一样的版本，不作任何改动;从 1 开始编号，版本 0 表示初始状态数组)
2. 访问某个历史版本上的某一位置的值

将每个 version 视为结点,那么所有的 version 连接构成了一棵树

https://www.luogu.com.cn/problem/P3919
https://qiita.com/wotsushi/items/72e7f8cdd674741ffd61#%E5%8F%82%E8%80%83%E8%A8%98%E4%BA%8B
https://www.cnblogs.com/Icys/p/Persistence_SegmentTree.html

py 库 https://github.com/tobgu/pyrsistent
js 库 https://github.com/immutable-js/immutable-js

## 完全可持久化并查集

用可持久化数组维护的并查集，本质和可持久化数组是一样的。
https://www.luogu.com.cn/blog/SSerxhs/solution-p3402

fa[a] = b，为了可持久化，我们就用可持久化数组来维护 fa[i]。注意这里不能再使用路径压缩了，道理很简单，可持久化要尽可能减少修改的次数。但是我们依然保留了一种优化方式：在维护 fa[i] 的同时维护一个 dep[i] ，表示这个节点的深度，保证在合并时是深度较小的点向深度较大的点合并即可。

## 完全可持久化平衡树

https://www.luogu.com.cn/problem/P3835

## 参考

[永続データ構造](https://qiita.com/wotsushi/items/72e7f8cdd674741ffd61#%E5%8F%82%E8%80%83%E8%A8%98%E4%BA%8B)
[持久化数据结构学习笔记——序列](https://zhuanlan.zhihu.com/p/33859991) (主要在说如何优化编码来节省空间)
[记录历史：持久化数据结构](https://quant67.com/post/algorithms/ads/persistent/persistent.html)
[可持久化数据结构](https://zh.m.wikipedia.org/zh-hans/%E5%8F%AF%E6%8C%81%E4%B9%85%E5%8C%96%E6%95%B0%E6%8D%AE%E7%BB%93%E6%9E%84)
[陈立杰：可持久化数据结构研究.pdf](https://github.com/Misaka233/algorithm/blob/master/%E9%99%88%E7%AB%8B%E6%9D%B0%EF%BC%9A%E5%8F%AF%E6%8C%81%E4%B9%85%E5%8C%96%E6%95%B0%E6%8D%AE%E7%BB%93%E6%9E%84%E7%A0%94%E7%A9%B6.pdf)
[Re 永続データ構造が分からない人のためのスライド](https://www.slideshare.net/qnighy/re-15443018)
[競技プログラミングにおける永続データ構造問題まとめ](https://blog.hamayanhamayan.com/entry/2017/05/21/001252)
[Persistent Data Structure](https://fuzhe1989.github.io/2017/11/07/persistent-data-structure/)

## 补充

什么地方需要用到持久性数据结构？

1. 函数式编程语言。它的定义就要求了不能有可变数据和可变数据结构。
2. 并发编程(eg:协同编辑)。
3. Lazy Evaluation。(如果一个数据结构是可变的，我们肯定不会放心对它使用 Lazy Evaluation。)
4. 使用 Persistent Map/HashMap 有助于简化 Prototype 的实现。
