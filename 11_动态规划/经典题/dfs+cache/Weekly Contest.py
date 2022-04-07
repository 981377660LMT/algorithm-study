# n ≤ 1,000
# k ≤ 50

# 1. 你最多尝试k次
# 2. 如果尝试失败，chance会减少(心理负担加重)
# 3. 做出一题时，考试结束
from functools import lru_cache
from math import floor


class Solution:
    def solve(self, points, chances, k):
        """求预期获得的最大分数:先挑战难题，shoot for the moon(力爭最好)"""

        @lru_cache(None)
        def dfs(index: int, remain: int) -> int:
            if remain <= 0:
                return 0
            if index == n:
                return 0

            jump = dfs(index + 1, remain)
            # 做不对的概率，作对的概率 相加
            select = dfs(index + 1, remain - 1) * (1 - (events[index][1] / 100)) + (
                events[index][0] * events[index][1] / 100
            )

            return max(jump, select)

        n = len(points)
        events = list(zip(points, chances))  # 先做最难，且最有把握的
        events.sort(reverse=True)

        res = dfs(0, k)
        dfs.cache_clear()
        return floor(res)


print(Solution().solve(points=[1000, 300, 10000], chances=[20, 100, 1], k=2))
# There's 3 problems in this contest:

# First problem wins 1000 points and you have 20% chance of solving it. This has expected value of 200 points.
# Second problem wins 300 points and you have 100% chance of solving it. This has expected value of 300 points.
# Third problem wins 10000 points and you have 1% chance of solving it. This has expected value of 100 points.
# The optimal strategy is to attempt the first problem for a chance to get a 1000 points, and then if it's not solvable, attempt the second problem for a guaranteed 300 points. The expected value is

# 440 = 0.2 * 1000 + 0.8 * 300
