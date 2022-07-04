from functools import lru_cache
import sys
from typing import List

sys.setrecursionlimit(int(1e9))
MOD = int(1e9 + 7)


@lru_cache(None)
def fac(n: int) -> int:
    """n的阶乘"""
    if n == 0:
        return 1
    return n * fac(n - 1) % MOD


@lru_cache(None)
def ifac(n: int) -> int:
    """n的阶乘的逆元"""
    return pow(fac(n), MOD - 2, MOD)


@lru_cache(None)
def C(n: int, k: int) -> int:
    if k < 0 or k > n:
        return 0
    return ((fac(n) * ifac(k)) % MOD * ifac(n - k)) % MOD


class Solution:
    def uniquePathsWithObstacles(self, obstacleGrid: List[List[int]]) -> int:
        row, col = len(obstacleGrid), len(obstacleGrid[0])
        bad = []
        bad.append((0, 0))
        for r in range(row):
            for c in range(col):
                if obstacleGrid[r][c] == 1:
                    bad.append((r + 1, c + 1))
        bad.append((row, col))
        bad.sort()

        n = len(bad)
        dp = [0] * n
        dp[0] = 1
        for i in range(1, n):
            r, c = bad[i]
            dp[i] = C(r + c - 2, r - 1)
            for j in range(1, i):
                preR, preC = bad[j]
                if preR <= r and preC <= c:
                    dp[i] -= dp[j] * C(r + c - preR - preC, r - preR)
                    dp[i] %= MOD

        return dp[n - 1]


print(Solution().uniquePathsWithObstacles([[0, 0, 0], [0, 1, 0], [0, 0, 0]]))
