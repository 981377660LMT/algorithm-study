# row,col<=1e5
# 墙壁数n<=3000
# 需要用O(n^2)的解法 而不是O(ROW*COL)的解法

# 左上角移动到右下角，一共有多少种路线
# 每一步可以向右或向下移动一格，不可以经过障碍物
# https://www.acwing.com/activity/content/code/content/2328319/
# AcWing 306. 杰拉尔德和巨型象棋
# 容斥原理，用无障碍情况下的方案数减去经过障碍的方案数即可。
# !dp[i]:原点出发走到第i个障碍物且中间不经过其余障碍物的方案数
# !到当前第i个障碍的合法方案等于到i的所有方案减去到i之前每个障碍的不合法方案
# !dp[i]=f(ri,ci)-∑dp[j]*f(ri-rj,ci-cj)
# 终点看作是最后一个障碍


import sys
from typing import List, Tuple


sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)

fac = [1]
ifac = [1]
for i in range(1, int(2e5 + 10)):
    fac.append(fac[-1] * i % MOD)
    ifac.append((ifac[-1] * pow(i, MOD - 2, MOD)) % MOD)


def C(n: int, k: int) -> int:
    if n < 0 or k < 0 or n < k:
        return 0
    return ((fac[n] * ifac[k]) % MOD * ifac[n - k]) % MOD


def cal(deltaX: int, deltaY: int) -> int:
    """没有障碍时向下移动delatX步,向右移动delatY步的方案数"""
    if deltaX < 0 or deltaY < 0:
        return 0
    return C(deltaX + deltaY, deltaX)


def solve(row: int, col: int, bad: List[Tuple[int, int]]) -> int:
    """左上角移动到右下角的方案数 bad是所有的障碍"""
    if (row - 1, col - 1) in bad or (0, 0) in bad:
        return 0
    n = len(bad)
    bad.append((row - 1, col - 1))
    bad.sort()
    dp = [0] * (n + 1)
    for i in range(n + 1):
        r, c = bad[i]
        dp[i] = cal(r, c)
        for j in range(i):
            preR, preC = bad[j]
            if preR <= r and preC <= c:
                dp[i] -= dp[j] * cal(r - preR, c - preC)
                dp[i] %= MOD
    return dp[n]


row, col, n = map(int, input().split())
bad = []
for _ in range(n):
    x, y = map(int, input().split())
    x, y = x - 1, y - 1
    bad.append((x, y))

print(solve(row, col, bad))
