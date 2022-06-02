# 请你帮忙验证 s3 是否是由 s1 和 s2 交错 组成的。
from functools import lru_cache


class Solution:
    def isInterleave(self, s1: str, s2: str, s3: str) -> bool:
        @lru_cache(None)
        def dfs(i: int, j: int) -> bool:
            if i == len(s1) and j == len(s2):
                return True
            res = False
            if i < len(s1) and s1[i] == s3[i + j]:
                res = res or dfs(i + 1, j)
            if j < len(s2) and s2[j] == s3[i + j]:
                res = res or dfs(i, j + 1)
            return res

        if len(s1) + len(s2) != len(s3):
            return False
        return dfs(0, 0)

