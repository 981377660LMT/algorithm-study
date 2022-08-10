# 兔子分组
# 两只兔子在一组的话 得分为aij(只计算一次)
# 求分组后的最大得分
# n<=16 =>O(3^n)


import sys


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

n = int(input())
matrix = []
for _ in range(n):
    matrix.append(list(map(int, input().split())))


# 预处理子集和
dpSum = [0] * (1 << n)
for state in range(1 << n):
    cur = [i for i in range(n) if state & (1 << i)]
    if len(cur) >= 2:
        for i in range(len(cur)):
            for j in range(i + 1, len(cur)):
                dpSum[state] += matrix[cur[i]][cur[j]]


# dp枚举子集的子集
# dp = dpSum[:]
dp = [0] * (1 << n)
for state in range(1 << n):
    dp[state] = max(0, dpSum[state])
    g1, g2 = state, 0
    while g1:
        dp[state] = max(dp[state], dp[g1] + dp[g2])
        g1 = (g1 - 1) & state
        g2 = state ^ g1

print(dp[-1])
