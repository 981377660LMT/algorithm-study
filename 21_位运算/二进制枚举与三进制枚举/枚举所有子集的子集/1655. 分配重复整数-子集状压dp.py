# dp[i][j] 表示：数组中的前 i 个元素，能否满足顾客的子集合 j的订单需求
from collections import Counter, defaultdict
from typing import List

# 1 <= needs.length <= 10
# nums 中至多有 50 个不同的数字。


class Solution:
    def canDistribute(self, nums: List[int], needs: List[int]) -> bool:
        m = len(needs)
        subsum = [0] * (1 << m)
        for i, num in enumerate(needs):
            for preState in range(1 << i):
                subsum[preState | (1 << i)] = subsum[preState] + num

        counter = Counter(nums)
        counts = sorted(counter.values(), reverse=True)
        n = len(counts)

        dp = [False] * (1 << m)
        for state in range(1 << m):
            dp[state] = counts[0] >= subsum[state]

        for i in range(1, n):
            ndp = [False] * (1 << m)
            for state in range(1 << m):
                # 注意这里需要判断g1是空集的情况
                if dp[state]:
                    ndp[state] = True
                    continue
                g1, g2 = state, 0
                while g1:  # g1非空
                    ndp[state] = ndp[state] or (dp[g2] and counts[i] >= subsum[g1])
                    g1 = (g1 - 1) & state
                    g2 = state ^ g1
            dp = ndp

        return dp[-1]


print(Solution().canDistribute(nums=[1, 2, 3, 3], needs=[2]))
