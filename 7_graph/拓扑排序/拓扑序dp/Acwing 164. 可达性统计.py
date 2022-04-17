# 给定一张 N 个点 M 条边的有向无环图，分别统计从每个点出发能够到达的点的数量。
# 1≤N,M≤30000


from collections import defaultdict
from functools import lru_cache
from typing import Set


n, m = map(int, input().split())
adjMap = defaultdict(set)


for _ in range(m):
    u, v = map(int, input().split())
    adjMap[u - 1].add(v - 1)


dp = {i: set([i]) for i in range(n)}


@lru_cache(None)
def dfs(cur: int) -> Set[int]:
    for next in adjMap[cur]:
        dp[cur] |= dfs(next)
    return dp[cur]


for key in range(n):
    print(len(dfs(key)))
dfs.cache_clear()


# from collections import defaultdict
# from functools import lru_cache


# n, m = map(int, input().split())
# adjMap = defaultdict(set)


# for _ in range(m):
#     u, v = map(int, input().split())
#     adjMap[u - 1].add(v - 1)


# res = [1 << i for i in range(n)]


# @lru_cache(None)
# def dfs(cur: int) -> int:
#     for next in adjMap[cur]:
#         res[cur] |= dfs(next)
#     return res[cur]


# for key in range(n):
#     print(bin(dfs(key))[2:].count('1'))

# dfs.cache_clear()

