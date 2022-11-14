# 给你一个字符串 s，请你将 s 分割成一些子串，使每个子串都是回文。
# 返回符合要求的 最少分割次数

# !2472. 不重叠回文子字符串的最大数目

from 中心扩展法 import expand2


class Solution:
    def minCut(self, s: str) -> int:
        isPalindrome = expand2(s)  # !可用马拉车优化到O(n)处理+O(1)判断
        n = len(s)
        dp = [n] * (n + 1)
        dp[0] = 0
        for i in range(1, n + 1):
            for j in range(i):
                if isPalindrome[j][i - 1]:
                    dp[i] = min(dp[i], dp[j] + 1)
        return dp[-1] - 1


print(Solution().minCut("aab"))
print(Solution().minCut("ab"))
