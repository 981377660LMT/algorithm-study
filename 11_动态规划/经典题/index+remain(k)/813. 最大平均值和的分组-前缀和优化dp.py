from typing import List
from functools import lru_cache
from itertools import accumulate

# 我们将给定的数组 A 分成 K 个相邻的非空子数组 ，
# !我们的分数由每个子数组内的平均值的总和构成。计算我们所能得到的最大分数是多少。
# 注意我们必须使用 A 数组中的每一个数进行分组，并且分数不一定需要是整数。
# 1 <= A.length <= 100.
# 1 <= A[i] <= 10000.


class Solution:
    def largestSumOfAverages2(self, nums: List[int], k: int) -> float:
        """dp[k][i]表示前i个数分成k组的最大平均值和 O(n*k)

        dp[c][i] = max(dp[c][i], dp[c-1][j] + (preSum[i+1] - preSum[j+1]) / (i - j))
        """
        n = len(nums)
        preSum = [0] + list(accumulate(nums))
        dp = [[0.0] * (n + 1) for _ in range(k + 1)]
        # !每个数分成1组的情况
        for i in range(1, n + 1):
            dp[1][i] = preSum[i] / i

        for c in range(2, k + 1):  # !分成2组开始dp
            for i in range(1, n + 1):
                for pre in range(i):  # 前面分割的位置
                    dp[c][i] = max(dp[c][i], dp[c - 1][pre] + (preSum[i] - preSum[pre]) / (i - pre))
        return dp[k][n]

    def largestSumOfAverages(self, nums: List[int], k: int) -> float:
        """O(n^2*k)"""
        n = len(nums)
        preSum = [0] + list(accumulate(nums))

        @lru_cache(None)
        def dfs(index: int, remain: int) -> float:
            if index == n:
                return 0
            if remain == 1:
                return (preSum[-1] - preSum[index]) / (n - index)
            res = 0
            for i in range(index, n):
                res = max(
                    res, (preSum[i + 1] - preSum[index]) / (i - index + 1) + dfs(i + 1, remain - 1)
                )
            return res

        return dfs(0, k)


print(Solution().largestSumOfAverages(nums=[9, 1, 2, 3, 9], k=3))
print(Solution().largestSumOfAverages2(nums=[9, 1, 2, 3, 9], k=3))
# 输出: 20
# 解释:
# A 的最优分组是[9], [1, 2, 3], [9]. 得到的分数是 9 + (1 + 2 + 3) / 3 + 9 = 20.
# 我们也可以把 A 分成[9, 1], [2], [3, 9].
# 这样的分组得到的分数为 5 + 2 + 6 = 13, 但不是最大值.
