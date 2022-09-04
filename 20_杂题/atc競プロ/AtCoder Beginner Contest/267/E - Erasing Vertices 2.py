"""
贪心 每次删删除代价最低的点
删完更新一遍其它点的cost
"""

import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

n, m = map(int, input().split())
values = list(map(int, input().split()))
adjList = [[] for _ in range(n)]
for _ in range(m):
    u, v = map(int, input().split())
    u, v = u - 1, v - 1
    adjList[u].append(v)
    adjList[v].append(u)
