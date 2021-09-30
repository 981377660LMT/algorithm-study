class Solution(object):
    def minCut(self, s):
        """
        :type s: str
        :rtype: int
        """
        N = len(s)
        dp = [N - 1] * N
        for i in range(N):
            if self.isPalindrome(s[0 : i + 1]):
                dp[i] = 0
                continue
            for j in range(i):
                if self.isPalindrome(s[j + 1 : i + 1]):
                    dp[i] = min(dp[i], dp[j] + 1)
        return dp[N - 1]

    def isPalindrome(self, s):
        return s == s[::-1]


# python切片 不超时 js超时
