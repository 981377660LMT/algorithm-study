from collections import defaultdict
from typing import List


class Solution:
    def minTransfers(self, transactions: List[List[int]]) -> int:
        """
        给定一群人之间的交易信息列表,计算能够还清所有债务的最小次数。
        n<=8
        """
        deg = defaultdict(int)
        for u, v, w in transactions:
            deg[u] += w
            deg[v] -= w
        nums = [deg[i] for i in deg if deg[i] != 0]
        n = len(nums)

        subsum = [0] * (1 << n)
        for i in range(n):
            for pre in range(1 << i):
                subsum[pre | (1 << i)] = subsum[pre] + nums[i]

        dp = [0] * (1 << n)  # 最多能分成多少个集合 使得每个集合的和为0
        for state in range(1 << n):
            g1, g2 = state, 0
            while g1:
                dp[state] = max(dp[state], int(subsum[g1] == 0) + dp[g2])
                g1 = (g1 - 1) & state
                g2 = state ^ g1
        return n - dp[-1]


print(Solution().minTransfers([[0, 1, 10], [1, 0, 1], [1, 2, 5], [2, 0, 5]]))

"""
子问题
给定一个有 n 个元素的集合，满足：
该集合的元素之和为 0,但该集合的所有真子集的元素之和均不为0。求最小的调配次数,使得所有元素都变为 0。
答案为 (n - 1)。 
((n - 1)次调配后，我们至少有 (n - 1) 个 0)

原问题
将所有元素分成最多的集合，使得每个集合的元素之和为 0。
设最后能分成 k 个集合，根据子问题的答案可得到原问题的答案就是 (n - k)。
通过位运算 DP 解决上述问题。
"""
