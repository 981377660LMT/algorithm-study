# https://atcoder.jp/contests/dp/tasks/dp_r
# !求有向图中长为k的路径数
# n<=50 k<=1e18

# 图的邻接矩阵幂
# !第k个dp[i][j] 表示从i到j的长为k的路径条数为dp[i][j]
# ndp[i][j] = ∑(dp[i][k]*dp[k][j])
# !注意到`dp转移类似于矩阵乘法` 所以用矩阵快速幂优化
# !时间复杂度O(n^3logk)
import sys


sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(4e18)

from matqpow import matqpow1

#############################################################

n, k = map(int, input().split())
adjMatrix = []  # 从i到j的转移方案数(路径长为1)
for _ in range(n):
    adjMatrix.append(list(map(int, input().split())))

T = matqpow1(adjMatrix, k, MOD)  # dp 矩阵
res = 0
for r in range(n):
    for c in range(n):
        res += T[r][c]
        res %= MOD
print(res)

# 输入:
# 4 2
# 0 1 0 0
# 0 0 1 1
# 0 0 0 1
# 1 0 0 0
# 输出:
# 6
