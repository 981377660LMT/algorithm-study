import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

from itertools import accumulate


class Solution:
    def numPermsDISequence(self, s: str) -> int:
        """优化:时间复杂度O(n^2)"""
        n = len(s)
        dp = [1] * (n + 1)
        for i in range(n):
            ndp, dpSum = [0] * (n + 1), [0] + list(accumulate(dp, lambda x, y: (x + y) % MOD))
            if s[i] == "I":
                for j in range(n - i):
                    ndp[j] = (dpSum[n - i + 1] - dpSum[j + 1]) % MOD
            else:
                for j in range(n - i):
                    ndp[j] = dpSum[j + 1] % MOD

            dp = ndp
        return dp[0]


n, k = map(int, input().split())
