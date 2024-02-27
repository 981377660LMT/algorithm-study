# 不管是啥 dp，先想暴力

`状态总数x每个状态的决策数x每次状态转移所需时间`

1. 减少状态总数
   - 修改状态表示 (这种题目会有数据量提示，例如 nums[i]<=1e5 `(为什么不是 1e9?)`)
2. 减少状态决策数
   - 单调队列优化
   - 斜率优化 (单调队列维护上/下凸壳)
   - 剪枝
3. 减少状态转移时间
   - 预处理
   - 数据结构优化(线段树查询)

<!-- 插头dp 四边形优化dp 太难了 不学了 -->

## 为什么 dp 优化这类题要用 dp 数组貰う dp 而不是记忆化 dfs

因为 dp 优化往往需要获取某个前缀的信息 记忆化 dfs 无法做到这一点

- https://blog.hamayanhamayan.com/entry/2017/03/20/234711
- https://codeforces.com/blog/entry/8219
- https://codeforces.com/blog/entry/47932

## 优化方案

- 前缀和
- 线段树等数据结构
- 矩阵快速幂(dp 在相邻行间转移,且每行状态数不超过 100) `O(n^3logk)`
- CHT
- 分治
- 四边形不等式
- offlineInlineDp

---

DP 优化方法大杂烩
https://www.cnblogs.com/alex-wei/p/DP_Involution.html
https://www.cnblogs.com/alex-wei/p/DP_optimization_method_II.html
https://www.cnblogs.com/alex-wei/p/dp_tricks.html
DP 做题记录
https://www.cnblogs.com/alex-wei/p/simple_DP.html
