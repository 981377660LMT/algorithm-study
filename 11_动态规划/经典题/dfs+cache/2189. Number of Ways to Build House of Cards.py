# 1 <= n <= 500
from functools import lru_cache

# 没太看懂题目，不过感觉状态只和 (当前行的三角形数，剩下的卡片数) 有关，所以考虑记忆化dfs搜索试一下。
# 时间复杂度为 O(n^3)，空间复杂度为 O(n^2)


class Solution:
    def houseOfCards(self, n: int) -> int:
        @lru_cache(None)
        def dfs(curTriangle: int, remain: int) -> int:
            if remain <= 0:
                return int(remain == 0)
            res = 0
            for nextTriangle in range(curTriangle):
                res += dfs(nextTriangle, remain - 2 * nextTriangle - (nextTriangle - 1))
            return res

        res = 0
        for baseTriangle in range(1, n):
            res += dfs(baseTriangle, n - 2 * baseTriangle - (baseTriangle - 1))
        return res


print(Solution().houseOfCards(5))
print(Solution().houseOfCards(16))
