# https://atcoder.jp/contests/tdpc/editorial/720
# https://simezi-tan.hatenadiary.org/entry/20140908/1410125522
# うなぎ

# n<=1000
# k<=50
# 在树中寻找k条不相交的路径有多少种方案

# dfs[root][allPath][curPath] O(n*k^2)
# `root为根的树中 有allPath条路径 从root出发的路径有curPath条` 时的方案数
# 每个dfs里是一个背包问题

import sys
import os

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(4e18)

n, k = map(int, input().split())
adjList = [[] for _ in range(n)]
for _ in range(n - 1):
    u, v = map(int, input().split())
    u, v = u - 1, v - 1
    adjList[u].append(v)
    adjList[v].append(u)


# TODO
