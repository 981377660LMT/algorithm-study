- 将区间[0,n)分成任意组的最小代价
  `dp[j] = min{dp[i] + f(i,j)`
  如果 f(i,j)可以拆成关于 f2(j) 的一次函数,即
  `dp[j] = min{dp[i] + f1(i)*f2(j) + f3(i) + f4(j)`
  那么可以用 CHT 优化 dp, `O(n^2) => O(n)/O(nlogn)`

- 将区间[0,n)分成 k 组的最小代价
  `dp[k][j] = min{dp[k-1][i] + f(i,j)`
  如果 f(i,j)可以拆成关于 f2(j) 的一次函数,即
  `dp[j] = min{dp[i] + f1(i)*f2(j) + f3(i) + f4(j)`
  那么可以用 CHT 优化 dp, `O(k*n^2) => O(k*n)/O(k*nlogn)`

## 实现 CHT 的几种方法:

1. deque
   https://tjkendev.github.io/procon-library/python/convex_hull_trick/deque.html
   使用条件:追加的直线斜率单增/单减。
2. 李超线段树
   https://tjkendev.github.io/procon-library/python/convex_hull_trick/li_chao_tree.html
   追加的直线斜率不满足单增/单减的时候,可以使用李超线段树。
   且动态开点不需要预先知道查询，非常方便。

## 实现

拆成 dp[j] = min(ij 交叉项+只含 i 的项+只含 j 的项和常数项)

---

原点为顶点的最大三角形面积
https://yukicoder.me/submissions/665173

---

https://atcoder.jp/contests/abc373/editorial/11028

**Convex hull trick は本来，いくつかの一次関数の最大値を求めるアルゴリズムですが，関数が互いにたかだか 1 回しか交わらないという性質を持てば，一次関数でないクラスの関数についても同様のアルゴリズムで最大値クエリを処理することができます**

`给定一个函数集，如果这些函数两两之间最多只有一个交点，那么可以用 CHT 求最值。`
