# 将 n 个不同的元素划分为 k 个圆排列的方案数

from typing import List


MOD = int(1e9 + 7)


def getStirling1Table(n: int, k: int) -> List[List[int]]:
    dp = [[0] * (k + 1) for _ in range(n + 1)]
    dp[0][0] = 1
    for i in range(1, n + 1):
        for j in range(1, k + 1):
            dp[i][j] = (dp[i - 1][j - 1] + (i - 1) * dp[i - 1][j]) % MOD
    return dp


if __name__ == "__main__":
    table = getStirling1Table(1010, 1010)
    # https://leetcode.cn/problems/number-of-ways-to-rearrange-sticks-with-k-sticks-visible/
    # 1866. 恰有 K 根木棍可以看到的排列数目
    class Solution:
        def rearrangeSticks(self, n: int, k: int) -> int:
            """长度为从 1 到 n 的整数。请你将这些木棍排成一排，并满足从左侧 可以看到 恰好 k 根木棍

            划分为k个部分,每个部分排列种数为圆排列数(每个部分的最大值站在开头)
            """
            return table[n][k]
