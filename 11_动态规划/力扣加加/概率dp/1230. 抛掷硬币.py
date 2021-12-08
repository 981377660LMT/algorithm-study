from typing import List

# 有一些不规则的硬币。在这些硬币中，prob[i] 表示第 i 枚硬币正面朝上的概率。
# 请对每一枚硬币抛掷 一次，然后返回正面朝上的硬币数等于 target 的概率。

# dp[i][j代表前 i 个硬币出现 j 次朝上的概率
# dp[i][j] = dp[i-1][j-1] * prob[i] + dp[i-1][j] * (1 - prob[i])
# 即 前面投了 j - 1 次, 再乘上这次正面.  或者  前面已经 j 次, 乘上 这次 反面几率.
class Solution:
    def probabilityOfHeads(self, prob: List[float], target: int) -> float:
        m, n = len(prob), target
        dp = [[0] * (n + 1) for _ in range(m + 1)]
        dp[0][0] = 1
        for i in range(m):
            dp[i + 1][0] = dp[i][0] * (1 - prob[i])

        for i in range(1, m + 1):
            for j in range(1, n + 1):
                dp[i][j] = dp[i - 1][j - 1] * prob[i - 1] + dp[i - 1][j] * (1 - prob[i - 1])
        return dp[-1][-1]


print(Solution().probabilityOfHeads(prob=[0.5, 0.5, 0.5, 0.5, 0.5], target=0))
