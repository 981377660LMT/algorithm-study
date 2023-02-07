# dp[k][j]=min(dp[k-1][i]+f(i,j)) (0<=i<j<=n)
# !将区间[0,n)分成k组的最小代价
# !f(i,j): 左闭右开区间[i,j)的代价(0<=i<j<=n)
# 时间复杂度O(k*n^2)
# 可用分治/CHT优化到O(k*nlogn)

from typing import Callable, List

INF = int(1e18)


# !f(i,j): 左闭右开区间[i,j)的代价(0<=i<j<=n)
def offlineDp1(k: int, n: int, f: Callable[[int, int], int]) -> int:
    dp = [[INF] * (n + 1) for _ in range(k + 1)]
    dp[1][0] = 0
    for j in range(1, n + 1):
        dp[1][j] = f(0, j)  # 分成1组时的代价
    for k_ in range(2, k + 1):  # !分成k_组
        for j in range(1, n + 1):
            for i in range(j):  # 枚举转移决策点
                dp[k_][j] = min(dp[k_][j], dp[k_ - 1][i] + f(i, j))
    return dp[k][n]


# !滚动更新
def offlineDp2(k: int, n: int, f: Callable[[int, int], int]) -> int:
    dp = [0] + [f(0, j) for j in range(1, n + 1)]  # 分成1组时的代价
    for _ in range(2, k + 1):
        ndp = [INF] * (n + 1)
        for j in range(1, n + 1):
            for i in range(j):  # 枚举转移决策点
                ndp[j] = min(ndp[j], dp[i] + f(i, j))
        dp = ndp
    return dp[n]


# 1959. K 次调整数组大小浪费的最小总空间
# https://leetcode.cn/problems/minimum-total-space-wasted-with-k-resizing-operations/
class Solution:
    def minSpaceWastedKResizing(self, nums: List[int], k: int) -> int:
        def f(i, j):
            if i == j:
                return 0
            return max(nums[i:j]) * (j - i) - sum(nums[i:j])  # 可用前缀和/st表优化

        return offlineDp2(k + 1, len(nums), f)


print(Solution().minSpaceWastedKResizing(nums=[10, 20, 15, 30, 20], k=2))
