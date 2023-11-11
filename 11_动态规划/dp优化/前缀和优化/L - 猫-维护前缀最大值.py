# 每只猫的幸福感为与距离1米内的猫的好感度fij之和
# !把1-n猫顺序放到坐标轴上 求幸福感最大和
# n<=1000
# -1000<=fij<=1000

# !dp[i][j] 表示 `猫i左边距离1米以内的猫从j开始` 时的幸福感之和
# dp[i][j] = max(dp[i-1][k]) + (f[i][j 到 i]) (k=1,2,...,j)
# !前缀和优化dp => 维护最前缀大值

from itertools import accumulate
import sys


sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

n = int(input())
grid = []  # 好感度矩阵
preSum = [[] for _ in range(n)]  # 用于求第i只猫左边有k只猫距离1米以内时的幸福感之和
for i in range(n):
    row = list(map(int, input().split()))
    grid.append(row)
    cur = [0] + list(accumulate(row))
    preSum[i] = cur


dp = [-INF] * (n + 5)
dp[0] = 0
for i in range(1, n):
    ndp = [-INF] * (n + 5)
    preMax = -INF  # !维护前缀最大值
    for j in range(i + 1):
        preMax = max(preMax, dp[j])
        ndp[j] = max(ndp[j], preMax + preSum[i][i] - preSum[i][j])
    dp = ndp
print(max(dp) * 2)  # 注意乘二是 f[i][j] 一个方向算一次，f[j][i]还要算一次
