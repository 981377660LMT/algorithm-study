from functools import lru_cache
from itertools import accumulate
from typing import List

MOD = int(1e9 + 7)
INF = int(1e20)

# 1 <= k <= sum(piles[i].length) <= 2000
# 1 <= n <= 1000

# 这道题时间复杂度为O(k*sum(piles[i].length))


class Solution:
    def maxValueOfCoins(self, piles: List[List[int]], k: int) -> int:
        """时间复杂度O(背包容量*物品个数)"""
        n = len(piles)
        preSums = []
        for i in range(n):
            preSums.append([0] + list(accumulate(piles[i])))

        @lru_cache(None)
        def dfs(index: int, remain: int) -> int:
            if remain < 0:
                return -int(1e20)
            if index >= n:
                if remain == 0:
                    return 0
                return -int(1e20)

            res = 0
            for select in range(len(piles[index])):
                next = dfs(index + 1, remain - select)
                res = max(res, next + preSums[index][select])
            return res

        res = dfs(0, k)
        dfs.cache_clear()
        return res


print(Solution().maxValueOfCoins(piles=[[1, 100, 3], [7, 8, 9]], k=2))
print(
    Solution().maxValueOfCoins(
        piles=[[100], [100], [100], [100], [100], [100], [1, 1, 1, 1, 1, 1, 700]], k=7
    )
)

