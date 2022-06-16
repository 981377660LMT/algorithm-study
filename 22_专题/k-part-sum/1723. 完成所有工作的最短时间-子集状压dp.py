from typing import List

# !前i个孩子 分配的饼干状态为j的情况下，前 i 个孩子的不公平程度的最小值
# !时间复杂度O(k*3^n)

# 如何用 Python 的状压 DP 写法通过此题
# https://leetcode.cn/problems/find-minimum-time-to-finish-all-jobs/solution/by-endlesscheng-d2oa/
# 优化 1：把 min 和 max 拆开，改为手写，避免额外的函数调用开销。
# 优化 2：预处理每个二进制数的子集。


class Solution:
    def distributeCookies(self, cookies: List[int], k: int) -> int:
        n = len(cookies)
        subSum = [0] * (1 << n)  # 子集对应的和 O(2^n) 求出
        for i, num in enumerate(cookies):
            for preState in range(1 << i):
                subSum[preState | (1 << i)] = subSum[preState] + num

        # 给第i个孩子分配的饼干集合为state
        # 那么前 i 个孩子的不公平程度就是max(dp[i-1][j^state],subSum[state])
        dp = [[int(1e20)] * (1 << n) for _ in range(k)]  # 前i个孩子 分配的饼干状态为j的情况下，前 i 个孩子的不公平程度的最小值
        for state in range(1 << n):
            dp[0][state] = subSum[state]

        # O(n*3^n) 枚举子集的子集
        for i in range(1, k):
            for state in range(1 << n):
                g1 = state
                g2 = state ^ g1
                while g1:
                    dp[i][state] = min(dp[i][state], max(dp[i - 1][g1], subSum[g2]))
                    g1 = (g1 - 1) & state
                    g2 = state ^ g1

        return dp[-1][-1]
