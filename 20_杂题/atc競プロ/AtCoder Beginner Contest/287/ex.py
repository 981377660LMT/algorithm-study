"""
!bitset优化floyd n可以到2000 O(n^3/64)
"""

import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


if __name__ == "__main__":
    n, m = map(int, input().split())
    adjList = [[] for _ in range(n)]
    for _ in range(m):
        u, v = map(int, input().split())
        adjList[u - 1].append(v - 1)
