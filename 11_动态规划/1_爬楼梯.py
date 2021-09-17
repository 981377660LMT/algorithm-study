MOD = 10 ** 9 + 7


class Solution:
    def waysToStep(self, n: int) -> int:

        if n == 1:
            return 1
        if n == 2:
            return 2
        if n == 3:
            return 4
        dp = [0] * (n + 1)
        dp[:4] = [1, 1, 2, 4]
        for i in range(3, n + 1):
            dp[i] = (dp[i - 1] % MOD + dp[i - 2] % MOD + dp[i - 3] % MOD) % MOD
        return dp[n]
