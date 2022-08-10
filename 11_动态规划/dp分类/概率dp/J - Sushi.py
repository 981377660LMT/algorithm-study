# 开始每个盘子里有ai个寿司
# 每次随机一个盘子编号 吃掉那个盘子的一个寿司
# 求吃完所有盘子里寿司需要的步数的期望值
# n<=300
# 1<=ai<=3

# !状态如何定义:dfs(one, two, three) 表示还有几个盘剩下一个/二个/三个寿司 时 吃完的期望次数

from collections import Counter
from functools import lru_cache
import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# 概率dp

n = int(input())
nums = list(map(float, input().split()))
counter = Counter(nums)


# memo = [-1] * (n + 1) * (n + 1) * (n + 1)
# def dfs(remain1: int, remain2: int, remain3: int) -> float:
#     if remain1 == remain2 == remain3 == 0:
#         return 0
#     hash_ = remain1 * n * n + remain2 * n + remain3
#     if memo[hash_] != -1:
#         return memo[hash_]

#     sum_ = remain1 + remain2 + remain3
#     res = n / sum_  # 吃到寿司
#     p1, p2, p3 = remain1 / sum_, remain2 / sum_, remain3 / sum_
#     if remain1:
#         res += p1 * dfs(remain1 - 1, remain2, remain3)
#     if remain2:
#         res += p2 * dfs(remain1 + 1, remain2 - 1, remain3)
#     if remain3:
#         res += p3 * dfs(remain1, remain2 + 1, remain3 - 1)
#     memo[hash_] = res
#     return res

a, b, c = counter[1], counter[2], counter[3]
dp = [[[0.0] * (n + 1) for _ in range(n + 1)] for _ in range(n + 1)]
for i in range(c + 1):
    for j in range(c + b + 1):
        for k in range(n + 1):
            if i + j + k > n:
                break
            remain = i + j + k
            if remain == 0:
                continue
            dp[i][j][k] += n / remain
            p1, p2, p3 = i / remain, j / remain, k / remain
            if i:
                dp[i][j][k] += p1 * dp[i - 1][j + 1][k]
            if j:
                dp[i][j][k] += p2 * dp[i][j - 1][k + 1]
            if k:
                dp[i][j][k] += p3 * dp[i][j][k - 1]
print(dp[c][b][a])
