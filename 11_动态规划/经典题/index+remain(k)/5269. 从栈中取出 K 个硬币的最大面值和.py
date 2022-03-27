from functools import lru_cache
from itertools import accumulate
from typing import List

MOD = int(1e9 + 7)
INF = int(1e20)


# 1 <= k <= sum(piles[i].length) <= 2000
# 1 <= n <= 1000

# 这道题时间复杂度为O(k*sum(piles[i].length))

# 换句话说就是暴力解法，对每个栈都讨论了当前能取到的最大值
class Solution:
    def maxValueOfCoins(self, piles: List[List[int]], k: int) -> int:
        """时间复杂度O(n*k)"""
        n = len(piles)
        preSums = []
        for i in range(n):
            preSums.append([0] + list(accumulate(piles[i])))

        @lru_cache(None)
        def dfs(index: int, remain: int) -> int:
            if index == n or remain < 0:
                return 0

            res = 0
            for select in range(min(remain + 1, len(piles[index]) + 1)):
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

