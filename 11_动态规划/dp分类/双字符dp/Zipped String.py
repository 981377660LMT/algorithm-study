# eturn whether there's any way of obtaining c by merging characters in order from a and b.
# 从a b 顺序选字符 能否组成c

# n<=1000


from functools import lru_cache


class Solution:
    def solve(self, a, b, c):
        @lru_cache(None)
        def dfs(i, j):
            if i + j >= len(c):
                return True

            res = False
            if i < len(a) and a[i] == c[i + j]:
                res = res | dfs(i + 1, j)
            if j < len(b) and b[j] == c[i + j]:
                res = res | dfs(i, j + 1)

            return res

        if len(c) != len(a) + len(b):
            return False

        res = dfs(0, 0)
        dfs.cache_clear()
        return res

