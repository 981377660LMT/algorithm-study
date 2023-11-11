# 男女配对
# n男n女配对 各有喜好 询问有多少种方法将男女完全配对。
# n <= 21

# 校园自行车分配
# 状压dp 复杂度优化
# 从 O(2^n*n^2) 优化到 O(2^n*n)
# !已经遍历过的女性visited中包含了配对的男生信息 因此优化掉index维度
# !dp[visited] 表示前popcount(visited)个男生与状态为visited的女生配对的方案数
import sys


sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(4e18)


def popcount(n: int):
    c = (n & 0x5555555555555555) + ((n >> 1) & 0x5555555555555555)
    c = (c & 0x3333333333333333) + ((c >> 2) & 0x3333333333333333)
    c = (c & 0x0F0F0F0F0F0F0F0F) + ((c >> 4) & 0x0F0F0F0F0F0F0F0F)
    c = (c & 0x00FF00FF00FF00FF) + ((c >> 8) & 0x00FF00FF00FF00FF)
    c = (c & 0x0000FFFF0000FFFF) + ((c >> 16) & 0x0000FFFF0000FFFF)
    c = (c & 0x00000000FFFFFFFF) + ((c >> 32) & 0x00000000FFFFFFFF)
    return c


########################################################################
n = int(input())
matrix = []
for _ in range(n):
    matrix.append(list(map(int, input().split())))


# !TLE
# @lru_cache(None)
# def dfs(index: int, visited: int) -> int:
#     if index == n:
#         return 1
#     res = 0
#     for next in range(n):
#         if matrix[index][next] and not (1 << next) & visited:
#             res += dfs(index + 1, visited | (1 << next))
#             res %= MOD
#     return res


# print(dfs(0, 0))

# !优化掉一个维度

# target = (1 << n) - 1


# @lru_cache(None)
# def dfs(visited: int) -> int:
#     if visited == target:
#         return 1

#     res = 0
#     matched = popcount(visited)
#     for next in range(n):
#         if not (visited >> next) & 1 and matrix[matched][next]:
#             res += dfs(visited | (1 << next))
#             res %= MOD
#     return res


# print(dfs(0))


# !优化掉一个维度
dp = [0] * (1 << n)
dp[0] = 1

for state in range(1, 1 << n):
    i = popcount(state)
    for j in range(n):  # 枚举女生
        if (state >> j) & 1 and matrix[i - 1][j]:
            dp[state] += dp[state ^ (1 << j)]
            dp[state] %= MOD
print(dp[-1])
