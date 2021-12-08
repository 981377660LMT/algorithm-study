# 依赖于前和下 可以从下到上、从前到后 或者 从前到后、从下到上
# 也可以 长度从短到长，相当于斜着来
from typing import List


class Solution:
    def PredictTheWinner(self, ns: List[int]) -> bool:
        N = len(ns)
        dp = [[0] * N for _ in range(N + 1)]
        for l in range(N):  # 长度从小到大
            for i in range(N - l):  # 以 i 为 开头的，l 为长度
                j = i + l
                dp[i][j] = max(ns[i] - dp[i + 1][j], ns[j] - dp[i][j - 1])
        print(dp)
        return dp[0][-1] >= 0


print(Solution().PredictTheWinner([1, 5, 2]))
print(Solution().PredictTheWinner([1, 5, 233, 7]))
