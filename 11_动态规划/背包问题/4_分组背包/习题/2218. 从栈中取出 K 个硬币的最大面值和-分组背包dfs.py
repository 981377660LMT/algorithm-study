# 2218. 从栈中取出 K 个硬币的最大面值和
# https://leetcode.cn/problems/maximum-value-of-k-coins-from-piles/
# 1 <= k <= sum(piles[i].length) <= 2000
# 1 <= n <= 1000
# 这道题时间复杂度为O(k*sum(piles[i].length))
# 分组背包，每组只能取一个前缀


from typing import List
from functools import lru_cache


INF = int(1e20)


class Solution:
    def maxValueOfCoins(self, piles: List[List[int]], k: int) -> int:
        """时间复杂度O(背包容量*物品个数)."""

        def max(a: int, b: int) -> int:
            return a if a > b else b

        @lru_cache(None)
        def dfs(index: int, remain: int) -> int:
            if index == n:
                return 0 if remain == 0 else -INF

            res = 0
            curGroup, curSum = piles[index], 0
            res = dfs(index + 1, remain)
            for select in range(1, min(remain, len(curGroup)) + 1):
                curSum += curGroup[select - 1]
                res = max(res, dfs(index + 1, remain - select) + curSum)
            return res

        n = len(piles)
        res = dfs(0, k)
        dfs.cache_clear()
        return res


print(Solution().maxValueOfCoins(piles=[[1, 100, 3], [7, 8, 9]], k=2))
print(
    Solution().maxValueOfCoins(
        piles=[[100], [100], [100], [100], [100], [100], [1, 1, 1, 1, 1, 1, 700]], k=7
    )
)
