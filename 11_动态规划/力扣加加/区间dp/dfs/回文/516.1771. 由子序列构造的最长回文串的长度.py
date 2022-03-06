# 返回可按上述方法构造的最长 回文串 的 长度 。如果无法构造回文串，返回 0 。
# 拼接后就是leetcode516题，但是有一个前提条件，从拼接的字符串中选取的子序列的起点必须小于len(word1), 终点必须大于等于len(word1)
# 1 <= word1.length, word2.length <= 1000
class Solution:
    def longestPalindrome(self, word1: str, word2: str) -> int:
        s = word1 + word2
        n = len(s)
        res = 0
        dp = [[0] * n for _ in range(n)]
        for i in range(n):
            dp[i][i] = 1

        for i in range(n - 1, -1, -1):
            for j in range(i + 1, n):
                if s[i] == s[j]:
                    dp[i][j] = dp[i + 1][j - 1] + 2
                    if i < len(word1) and j >= len(word1):
                        res = max(res, dp[i][j])
                else:
                    dp[i][j] = max(dp[i + 1][j], dp[i][j - 1])
        return res


print(Solution().longestPalindrome(word1="cacb", word2="cbba"))
# 输出：5
# 解释：从 word1 中选出 "ab" ，从 word2 中选出 "cba" ，得到回文串 "abcba" 。
print(Solution().longestPalindrome(word1="aa", word2="bb"))
