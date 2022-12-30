from functools import lru_cache

# 给定两个字符串s1 和 s2，返回 使两个字符串相等所需删除字符的 ASCII 值的最小和 。


class Solution:
    def minimumDeleteSum(self, s1: str, s2: str) -> int:
        @lru_cache(None)
        def dfs(i: int, j: int) -> int:
            if i == n1:
                return sufSum2[n2 - j]
            if j == n2:
                return sufSum1[n1 - i]
            if s1[i] == s2[j]:
                return dfs(i + 1, j + 1)
            return min(dfs(i + 1, j) + ord(s1[i]), dfs(i, j + 1) + ord(s2[j]))

        n1, n2 = len(s1), len(s2)
        sufSum1, sufSum2 = [0] * (n1 + 1), [0] * (n2 + 1)
        for i in range(n1 - 1, -1, -1):
            sufSum1[i] = sufSum1[i + 1] + ord(s1[i])
        for i in range(n2 - 1, -1, -1):
            sufSum2[i] = sufSum2[i + 1] + ord(s2[i])
        sufSum1, sufSum2 = sufSum1[::-1], sufSum2[::-1]

        res = dfs(0, 0)
        dfs.cache_clear()
        return res


assert Solution().minimumDeleteSum("sea", "eat") == 231
