# abc419-E - Subarray Sum Divisibility-按模分组dp
# https://atcoder.jp/contests/abc419/tasks/abc419_e
# 给定一个长度为N的整数数组。每次操作可以任选其中一个数，将其加一。
# 求最小操作次数，使得每个长度为L的子段和都能被M整除。
#
# N,M<=500
#
# 当且仅当：
# 1. (A[0]+A[1]+...+A[L-1]) % M == 0
# 2. (A[i] - A[i+L]) % M == 0
#
# !我们可以按照 i % L 进行分组，进行分组，每组的组内都必须满足同余 M。
# !cost[i][j] -> 第i个分组在模 M 等于 j 的情况下，最少需要多少次操作。
# !dp[i][j] -> 前i个数之和模 M 等于 j 的最小操作数。

INF = int(1e18)

N, M, L = map(int, input().split())
A = list(map(int, input().split()))

groups = [[] for _ in range(L)]
for i, v in enumerate(A):
    groups[i % L].append(v)


cost = [[0] * M for _ in range(L)]
for g in range(L):  # 枚举每个组
    for m in range(M):
        cost[g][m] = sum((m - v) % M for v in groups[g])

dp = [[INF] * M for _ in range(L + 1)]
dp[0][0] = 0
for i in range(L):
    for s in range(M):
        for t in range(M):
            cur = (s + t) % M
            tmp = dp[i][s] + cost[i][t]
            if tmp < dp[i + 1][cur]:
                dp[i + 1][cur] = tmp
print(dp[L][0])
