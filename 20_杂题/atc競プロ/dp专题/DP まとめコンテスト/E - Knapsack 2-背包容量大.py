# n<=100
# w<=1e9 scorei<=1e3
# !容量大但是物品价值不大 需要改变dp状态定义:
# !dp[i][score] 前 i 个物品取到score时的最小体积

import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

n, w = map(int, input().split())
scoreSum = 0
goods = []
for _ in range(n):
    weight, score = map(int, input().split())
    goods.append((weight, score))
    scoreSum += score

dp = [INF] * int(scoreSum + 5)
dp[0] = 0
for weight, score in goods:
    for s in range(scoreSum, -1, -1):
        if s - score < 0:
            break
        dp[s] = min(dp[s], dp[s - score] + weight)

res = INF
for i in range(scoreSum, -1, -1):
    if dp[i] <= w:
        print(i)
        exit(0)
