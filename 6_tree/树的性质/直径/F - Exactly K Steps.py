"""
给定一棵n(n<=2e5)个节点的树。
有q(q≤2e5)次询问,每次询问给出两个数字(u, k),
请找到距离点u的为k的点(输出任意一个即可)。如果没有输出-1。

q 个查询 询问距离树结点u距离为k的结点是否存在,并输出一个这样的结点

k的最大值来自于u到直径的两个端点
寻找直径+倍增
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
