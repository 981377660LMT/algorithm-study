# 1 <= s.length <= 2000

# 此解法O(n^3)
class Solution(object):
    def minCut(self, s):
        """
        :type s: str
        :rtype: int
        """
        n = len(s)
        dp = [n - 1] * n
        for i in range(n):
            if self.isPalindrome(s[0 : i + 1]):
                dp[i] = 0
                continue

            # 枚举分割点
            for j in range(i):
                if self.isPalindrome(s[j + 1 : i + 1]):
                    dp[i] = min(dp[i], dp[j] + 1)

        return dp[n - 1]

    def isPalindrome(self, s):
        return s == s[::-1]


# python切片 不超时 js超时
print(Solution().minCut(s="aab"))
