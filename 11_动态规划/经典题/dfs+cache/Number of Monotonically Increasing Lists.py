# 单增数组的个数
# Return the number of distinct lists such that for each list:

# There are k positive numbers whose sum is n
# Every number is greater than or equal to the number on its left

from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e9))


class Solution:
    def solve(self, n, k):
        @lru_cache(None)
        def dfs(toSum: int, toSelect: int) -> int:
            """讨论开头的数是不是1，如果不是，将每个数减去k递归"""
            if toSum < toSelect:
                return 0
            if toSelect <= 0:
                return int(toSum == 0)
            return (dfs(toSum - 1, toSelect - 1) + dfs(toSum - toSelect, toSelect)) % int(1e9 + 7)

        return dfs(n, k)


print(Solution().solve(n=8, k=4))
