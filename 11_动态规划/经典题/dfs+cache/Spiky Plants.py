# 多刺植物

# n ≤ 100,000
from functools import lru_cache

INF = int(1e20)


class Solution:
    def solve(self, heights, costs):
        """最小成本增加植物的高度，使得相邻高度不同"""

        @lru_cache(None)
        def dfs(index: int, pre: int) -> int:
            """复杂度O(3*n)"""
            if index == n:
                return 0

            res = INF
            for i in range(3):
                if heights[index] + i != pre:
                    res = min(res, costs[index] * i + dfs(index + 1, heights[index] + i))
            return res

        n = len(heights)
        res = dfs(0, INF)
        dfs.cache_clear()
        return res


print(
    Solution().solve(
        heights=[1, 1, 2],
        costs=[3, 1, 7],
    )
)

# We can increase the second height by 2, which costs 2 * 1 = 2.
