# 你可以删除 s 中任意数目的字符，使得 s 平衡 。我们称 s 平衡的 当不存在下标对 (i,j) 满足 i < j 且 s[i] = 'b' 同时 s[j]= 'a' 。
# 请你返回使 s 平衡 的 最少 删除次数。
# 。我们称 s 平衡的 当不存在下标对 (i,j) 满足 i < j 且 s[i] = 'b' 同时 s[j]= 'a' 。

# 1 <= s.length <= 10^5

# 本质上是求最长非递减子序列长度
class Solution:
    def minimumDeletions(self, s: str) -> int:
        n = len(s)
        dp = [0] * (n + 1)

        b = 0
        for i in range(1, n + 1):
            if s[i - 1] == 'b':
                dp[i] = dp[i - 1]
                b += 1
            else:
                dp[i] = min(dp[i - 1] + 1, b)

        # 删这个a或删前面所有的b
        return dp[-1]


print(Solution().minimumDeletions(s="aababbab"))
# 输出：2
# 解释：你可以选择以下任意一种方案：
# 下标从 0 开始，删除第 2 和第 6 个字符（"aababbab" -> "aaabbb"），
# 下标从 0 开始，删除第 3 和第 6 个字符（"aababbab" -> "aabbbb"）。
