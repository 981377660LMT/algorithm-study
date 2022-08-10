# 在树上选取若干个点，这些点之间都是连通的
# 对v从0到n-1 输出`选取顶点v的组合`有多少种
# n<=1e5


from collections import deque
from functools import lru_cache
import sys
from typing import Any, Callable, List


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
INF = int(4e18)


# n, MOD = map(int, input().split())
# adjList = [[] for _ in range(n)]
# for _ in range(n - 1):
#     a, b = map(int, input().split())
#     adjList[a - 1].append(b - 1)
#     adjList[b - 1].append(a - 1)


# TLE TODO
# @lru_cache(None)
# def dfs(cur: int, pre: int) -> int:
#     res = 1
#     for next in adjList[cur]:
#         if next == pre:
#             continue
#         nextRes = dfs(next, cur)
#         res *= nextRes + 1  # 加1表示子节点不选
#         res %= MOD
#     return res


# for cur in range(n):
#     print(dfs(cur, -1))

###############################################################
# !正解似乎是换根dp
def rerooting(
    N: int,
    T: List[List[int]],
    init: Callable[[], int],
    node: Callable[[int], int],
    merge: Callable[[int, int], int],
):
    """换根dp框架 https://atcoder.jp/contests/dp/submissions/33369486"""
    N1 = N + 1
    dq = deque([(1, 0)])
    order = []
    while dq:
        xp = dq.pop()
        x, p = xp
        order.append(xp)
        for xx in T[x]:
            if xx == p:
                continue
            dq.append((xx, x))

    dp1 = [0] * N1
    ls = [0] * N1
    rs = [0] * N1

    for x, p in reversed(order):
        tx = T[x]
        v = init()
        for xx in tx:
            if xx == p:
                continue
            ls[xx] = v
            v = merge(v, node(dp1[xx]))
        v = init()
        for xx in reversed(tx):
            if xx == p:
                continue
            rs[xx] = v
            v = merge(v, node(dp1[xx]))
        dp1[x] = v

    dpx = [0] * N1
    res = [0] * N1
    res[1] = dp1[1]
    for i in range(1, N):
        x, p = order[i]
        dpx[x] = v = merge(merge(ls[x], rs[x]), node(dpx[p]))
        res[x] = merge(node(v), dp1[x])
    return res


n, MOD = map(int, input().split())
if n == 1:
    print(1)
    exit(0)

adjList = [[] for _ in range(n + 1)]
for _ in range(n - 1):
    a, b = map(int, input().split())
    adjList[a].append(b)
    adjList[b].append(a)


res = rerooting(n, adjList, lambda: 1, lambda x: (x + 1) % MOD, lambda x, y: x * y % MOD)
print(*res, sep="\n")
