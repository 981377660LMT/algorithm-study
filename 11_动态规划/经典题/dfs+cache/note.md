计算最优解问题 ———— **@cache,永远的神**

1. `cur + k`为参数的 dfs 模型
   直接标准模板
   `813. 最大平均值和的分组.py`
   `1959. K 次调整数组大小浪费的最小总空间.py`

   ```Python
   class Solution:
      def largestSumOfAverages(self, nums: List[int], k: int) -> float:
          n = len(nums)
          preSum = [0] + list(accumulate(nums))

          @lru_cache(None)
          def dfs(cur: int, remain: int) -> float:
              if cur == n:
                  return 0
              if remain == 1:
                  return (preSum[-1] - preSum[cur]) / (n - cur)

              res = 0
              for i in range(cur, n-remain): # 注意这里不写n,写n-remain可以减枝
                  res = max(
                      res, (preSum[i + 1] - preSum[cur]) / (i - cur + 1) + dfs(i + 1, remain - 1)
                  )

              return res

          return dfs(0, k)
   ```

2. 几种中间的更新 res 过程:

   1. 枚举分割点更新(范围条件/有限状态条件)
   2. dfs 后序更新
   3. `当前cur选或不选`更新

3. 直接模拟 bfs (记忆化+参数范围剪枝)

4. 注意事项：cache 有时要加 maxsize=状态数才能 ac
   `1770. 执行乘法运算的最大分数.py`
   关于使用 cache，lru_cache 超时的原因（Python 垃圾回收机制）
   https://leetcode-cn.com/problems/maximum-score-from-performing-multiplication-operations/solution/ben-gai-san-fen-zhong-jie-jue-by-r4c12-8kxx/
