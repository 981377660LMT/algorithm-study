# 三个字符串的LCS
from functools import lru_cache


class Solution:
    def solve(self, a, b, c):
        @lru_cache(None)
        def dfs(i, j, k):
            if i == n1 or j == n2 or k == n3:
                return 0

            if a[i] == b[j] == c[k]:
                return dfs(i + 1, j + 1, k + 1) + 1
            else:
                return max(dfs(i + 1, j, k), dfs(i, j + 1, k), dfs(i, j, k + 1))

        n1, n2, n3 = len(a), len(b), len(c)
        return dfs(0, 0, 0)
