# dp[j]=min(dp[i]+f(i,j)) (0<=i<j)
# !f(i,j): 左闭右开区间[i,j)的代价(0<=i<j<=n)
# 时间复杂度O(n^2)
# 可用CHT优化到O(nlogn)/O(n) 或者 offlineOnlineDp 优化到O(nlogn^2)


from typing import Callable


def onlineDp(n: int, f: Callable[[int, int], int]) -> int:
    dp = [0] + [f(0, j) for j in range(1, n + 1)]  # 分成1组时的代价
    for j in range(1, n + 1):
        for i in range(j):  # 枚举转移决策点
            dp[j] = min(dp[j], dp[i] + f(i, j))
    return dp[n]
