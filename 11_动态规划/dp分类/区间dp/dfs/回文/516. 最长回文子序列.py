from functools import lru_cache


class Solution:
    def longestPalindromeSubseq(self, s: str) -> int:
        """最长回文子序列"""

        @lru_cache(None)
        def dfs(left: int, right: int) -> int:
            if left > right:
                return 0
            if left == right:
                return 1
            if s[left] == s[right]:
                return dfs(left + 1, right - 1) + 2
            return max(dfs(left, right - 1), dfs(left + 1, right))

        n = len(s)
        return dfs(0, n - 1)


print(Solution().longestPalindromeSubseq(s="bbbab"))
