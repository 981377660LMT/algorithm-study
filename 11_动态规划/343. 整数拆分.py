# 给定一个正整数 n ，将其拆分为 k 个 正整数 的和（ k >= 2 ），并使这些整数的乘积最大化。


class Solution:
    def integerBreak(self, n: int) -> int:
        if n < 4:
            return n - 1
        res = 1
        while n > 4:
            res *= 3
            n -= 3
        return res * n

    def integerBreak2(self, n: int) -> int:
        dp = [0] * (n + 1)
        for i in range(2, n + 1):
            max_ = 0
            for j in range(1, i):
                # (i-j) 继续拆 or 不拆
                max_ = max(max_, j * (i - j), j * dp[i - j])
            dp[i] = max_
        return dp[n]
