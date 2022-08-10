# n<=100
# w<=1e5 scorei<=1e9
# !容量不大但是物品价值很大
# !dp[i][cap]  前 i 个物品总体积不超过 cap 的最大价值

import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

n, w = map(int, input().split())
goods = []
for _ in range(n):
    weight, score = map(int, input().split())
    goods.append((weight, score))

dp = [0] * (w + 1)
for weight, score in goods:
    for cap in range(w, -1, -1):
        if cap - weight < 0:
            break
        dp[cap] = max(dp[cap], dp[cap - weight] + score)

print(max(dp))
