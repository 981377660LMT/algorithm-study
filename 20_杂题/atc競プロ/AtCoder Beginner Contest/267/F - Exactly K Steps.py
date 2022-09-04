"""
q 个查询 询问距离结点u距离为k的结点个数并输出一个这样的结点
启发式合并
"""

import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


if __name__ == "__main__":

    n = int(input())
    adjList = [[] for _ in range(n)]
    for _ in range(n - 1):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        adjList[u].append(v)
        adjList[v].append(u)
    q = int(input())
    for _ in range(q):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
