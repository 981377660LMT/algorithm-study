# 每次可以花费A[i] 来喂动物i和i+1 (取模)
# 所有动物都喂到的最小花费
# 环形分类:第一个A[0]选不选
# !2<=n<=3e5

from functools import lru_cache
import sys


sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(1e18)


# # !dfs python3.8
# @lru_cache(None)
# def dfs(index: int, hasPre: bool, root: bool) -> int:
#     """当前在index 前一个点是否选择 第一个点是否选择"""
#     if index == n:
#         return 0 if (hasPre or root) else INF
#     res = dfs(index + 1, True, root) + cost[index]
#     if hasPre:
#         res = min(res, dfs(index + 1, False, root))
#     return res


# n = int(input())
# cost = list(map(int, input().split()))
# print(min(dfs(1, True, True) + cost[0], dfs(1, False, False)))

######################################################################
# !dp
n = int(input())
cost = list(map(int, input().split()))
dp = [[[INF, INF] for _ in range(2)] for _ in range(n)]  # (index,pre,root) [不选 选]
dp[0][0][0] = 0
dp[0][1][1] = cost[0]
for i in range(1, n):
    for pre in range(2):
        for cur in range(2):
            if pre or cur:
                for root in range(2):
                    dp[i][cur][root] = min(
                        dp[i][cur][root], dp[i - 1][pre][root] + (cost[i] if cur else 0)
                    )

res = INF
for pre in range(2):
    for root in range(2):
        if pre or root:
            res = min(res, dp[-1][pre][root])  # !最少要选一个
print(res)
