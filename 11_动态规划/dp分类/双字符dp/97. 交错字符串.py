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
                res |= dfs(i + 1, j)
            if j < len(s2) and s2[j] == s3[i + j]:
                res |= dfs(i, j + 1)
            return res

        if len(s1) + len(s2) != len(s3):
            return False
        return dfs(0, 0)

    def isInterleave2(self, s1: str, s2: str, s3: str) -> bool:

        n1, n2, n3 = len(s1), len(s2), len(s3)
        if n1 + n2 != n3:
            return False

        dp = [[False] * (n2 + 1) for _ in range(n1 + 1)]
        dp[0][0] = True
        for i in range(n1 + 1):  # !注意取0个字符的情况
            for j in range(n2 + 1):
                if i > 0:
                    if s1[i - 1] == s3[i + j - 1]:
                        dp[i][j] |= dp[i - 1][j]
                if j > 0:
                    if s2[j - 1] == s3[i + j - 1]:
                        dp[i][j] |= dp[i][j - 1]

        return dp[-1][-1]


print(Solution().isInterleave("aabcc", "dbbca", "aadbbcbcac"))
