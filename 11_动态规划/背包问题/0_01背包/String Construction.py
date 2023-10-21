# 不理解dp数组时 就用记忆化dfs

# n ≤ 50
# a,b<=50
# 问最多可以取几个数
from functools import lru_cache


class Solution:
    def solve(self, strings, a, b):
        @lru_cache(None)
        def dfs(index: int, remainA: int, remainB: int) -> int:
            if remainA < 0 or remainB < 0:
                return -int(1e20)
            if index == n:
                return 0

            res = dfs(index + 1, remainA, remainB)
            countA, countB = strings[index].count('A'), strings[index].count('B')
            res = max(res, 1 + dfs(index + 1, remainA - countA, remainB - countB))
            return res

        n = len(strings)
        res = dfs(0, a, b)
        dfs.cache_clear()
        return res


print(Solution().solve(["AABB", "AAAB", "A", "B"], 4, 2))
# We can take these strings using 4 "A"s and 2 "B"s ["AAAB","A","B"]
