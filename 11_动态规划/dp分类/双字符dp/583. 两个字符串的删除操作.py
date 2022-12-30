# 583. 两个字符串的删除操作
from functools import lru_cache


class Solution:
    def minDistance(self, word1: str, word2: str) -> int:
        @lru_cache(None)
        def dfs(i: int, j: int) -> int:
            if i == n1:
                return n2 - j
            if j == n2:
                return n1 - i
            if word1[i] == word2[j]:
                return dfs(i + 1, j + 1)
            return min(dfs(i + 1, j) + 1, dfs(i, j + 1) + 1)

        n1, n2 = len(word1), len(word2)
        res = dfs(0, 0)
        dfs.cache_clear()
        return res


assert Solution().minDistance("sea", "eat") == 2
