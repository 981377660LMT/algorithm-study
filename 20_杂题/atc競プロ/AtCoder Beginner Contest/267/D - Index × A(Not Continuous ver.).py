"""卡记忆化搜索"""

import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(1e18)

n, m = map(int, input().split())
nums = list(map(int, input().split()))


# def dfs(index: int, count: int) -> int:
#     if index == n:
#         return 0 if count == m else -INF
#     hash_ = index * n + count
#     if memo[hash_] != -INF:
#         return memo[hash_]
#     res = dfs(index + 1, count)
#     if count + 1 <= m:
#         res = max(res, dfs(index + 1, count + 1) + nums[index] * (count + 1))
#     memo[hash_] = res
#     return res


# memo = [-INF] * (n + 1) * (n + 1)
# print(dfs(0, 0))

# dp[i][j] 为前i个数中, 选择了j个数的方案数
dp = [[-INF] * (m + 1) for _ in range(n + 1)]
dp[0][0] = 0
for i in range(1, n + 1):
    for j in range(m + 1):
        dp[i][j] = max(dp[i][j], dp[i - 1][j])
        if j + 1 <= m:
            dp[i][j + 1] = max(dp[i][j + 1], dp[i - 1][j] + nums[i - 1] * (j + 1))


print(dp[n][m])
