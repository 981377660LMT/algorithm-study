from functools import lru_cache


class Solution:
    def longestPalindromeSubseq(self, s: str) -> int:
        n = len(s)

        @lru_cache(None)
        def dfs(l, r) -> int:
            if l > r:
                return 0
            if l == r:
                return 1
            if s[l] == s[r]:
                return dfs(l + 1, r - 1) + 2
            return max(dfs(l, r - 1), dfs(l + 1, r))

        return dfs(0, n - 1)


print(Solution().longestPalindromeSubseq(s="bbbab"))
