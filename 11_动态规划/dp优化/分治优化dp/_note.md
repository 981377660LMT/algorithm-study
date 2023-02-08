offlineDp -> onlineDp (オフライン・オンライン変換)

- offlineDp 可以理解为将区间[0,n)分成 k 组的最小代价
  dp[k][j]=min(dp[k-1][i]+f(i,j)) (i<j)

  - 如果代价函数 f(i,j) 是决策单调(Monotone)的，那么可以用分治 dp 优化，时间复杂度从 `O(n^2*k)` 降到 `O(n*logn*k)`
    > 备注: Monge ⇒ Totally Monotone(TM) ⇒ Monotones，即满足四边形不等式的函数一定满足决策单调性

- onlineDp 可以理解为将区间[0,n)分成`任意组`的最小代价
  dp[j]=min(dp[i]+f(i,j)) (i<j)
  - 如果 offline 问题存在复杂度 O(M(n))的解,那么 online 问题存在复杂度 O(M(n)logn)的解
  - 如果 f 可以拆成 i 的一次函数,那么可以用 CHT(ConvexHullTrick) 优化 dp
