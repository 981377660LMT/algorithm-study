# 掷骰子n次 求正面次数多于反面的概率
# n<=3000

import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(4e18)


n = int(input())
nums = list(map(float, input().split()))

target = (n + 1) // 2

dp = [0.0] * (n + 10)
dp[0] = 1 - nums[0]
dp[1] = nums[0]
for i in range(1, n):
    ndp = [0.0] * (n + 10)
    for j in range(n + 1):
        ndp[j] += dp[j] * (1 - nums[i])
        ndp[j + 1] += dp[j] * nums[i]
    dp = ndp

print(sum(dp[i] for i in range(target, n + 1)))
