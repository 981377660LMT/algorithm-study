from functools import lru_cache

# 长度为偶数。
# 除中间的两个字符外，其余任意两个连续字符不相等。
# 返回 s 的最长“好的回文子序列”的长度。
# 1 <= s.length <= 250


class Solution:
    def longestPalindromeSubseq(self, s: str) -> int:
        # dp(i,j,c)代表i到j之间且两端字母为c时的解，两端字母c只有26种选择，穷举就行。
        # 复杂度(状态数)：O(26n^2)
        @lru_cache(None)
        def dfs(l: int, r: int, char: int) -> int:
            if l >= r:
                return 0
            if s[l] != char:
                return dfs(l + 1, r, char)
            if s[r] != char:
                return dfs(l, r - 1, char)
            return (
                max(
                    dfs(l + 1, r - 1, chr(ord('a') + i))
                    for i in range(26)
                    if chr(ord('a') + i) != char
                )
                + 2
            )

        # 枚举
        return max(dfs(0, len(s) - 1, chr(ord('a') + i)) for i in range(26))


print(Solution().longestPalindromeSubseq(s="bbabab"))
# 输出: 4
# 解释: s 的最长“好的回文子序列”是 "baab"。
