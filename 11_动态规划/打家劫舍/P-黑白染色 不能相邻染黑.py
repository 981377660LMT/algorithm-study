"""黑白染色 不能相邻染黑 求方案数"""

import sys
from typing import Tuple

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(4e18)


n = int(input())
adjList = [[] for _ in range(n)]
for i in range(n - 1):
    a, b = map(int, input().split())
    adjList[a - 1].append(b - 1)
    adjList[b - 1].append(a - 1)


# @lru_cache(None)
# def dfs(cur: int, pre: int, color: int) -> int:
#     res = 1
#     for next in adjList[cur]:
#         if next != pre:
#             cand = dfs(next, cur, 0)
#             if color != 1:
#                 cand += dfs(next, cur, 1)
#                 cand %= MOD
#             res = (res * cand) % MOD
#     return res


# res = dfs(0, -1, 0) + dfs(0, -1, 1)
# res %= MOD
# print(res)


def dfs(cur: int, pre: int) -> Tuple[int, int]:
    """黑白染色 不能相邻染黑 求方案数"""
    white, black = 1, 1
    for next in adjList[cur]:
        if next != pre:
            w, b = dfs(next, cur)
            white = (white * ((w + b) % MOD)) % MOD
            black = (black * w) % MOD

    return white, black


res = sum(dfs(0, -1))
res %= MOD
print(res)
