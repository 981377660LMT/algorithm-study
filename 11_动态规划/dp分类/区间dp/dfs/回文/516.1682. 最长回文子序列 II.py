from functools import lru_cache
from string import ascii_lowercase

# !长度为偶数。
# 除中间的两个字符外，其余任意两个连续字符不相等。
# 返回 s 的最长“好的回文子序列”的长度。
# 1 <= s.length <= 250


class Solution:
    def longestPalindromeSubseq(self, s: str) -> int:
        @lru_cache(None)
        def dfs(left: int, right: int, char: str) -> int:
            """[left,right]这一段里最长的回文子序列的长度，两端字母为char
            
            
            """
            if left >= right:
                return 0
            if s[left] != char:
                return dfs(left + 1, right, char)
            if s[right] != char:
                return dfs(left, right - 1, char)

            res = 0
            for nextChar in ascii_lowercase:
                if nextChar == char:
                    continue
                res = max(res, dfs(left + 1, right - 1, nextChar) + 2)
            return res

        n = len(s)
        return max(dfs(0, n - 1, char) for char in ascii_lowercase)


print(Solution().longestPalindromeSubseq(s="bbabab"))
# 输出: 4
# 解释: s 的最长“好的回文子序列”是 "baab"。
