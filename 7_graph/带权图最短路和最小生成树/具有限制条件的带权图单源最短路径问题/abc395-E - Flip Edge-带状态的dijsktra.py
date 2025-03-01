# E - Flip Edge (带状态的dijsktra)
# https://atcoder.jp/contests/abc395/tasks/abc395_e
#
# 概述
# 题目给定一个 N 顶点 M 边的有向图，你从顶点 1 开始，想要到达顶点 N。 每一步可以执行以下两种操作之一：
#
# !沿着一条有向边移动到相连顶点，成本为 1
# !反转图中所有边的方向，成本为 X
# 要求：计算从顶点 1 到顶点 N 的最小总成本。
#
# 思路
#
# 维护两个状态：原图状态和反图状态
# 在这两个状态之间用 Dijkstra 算法寻找最短路径


from heapq import heappop, heappush
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n, m, x = map(int, input().split())

    graph, revGraph = [[] for _ in range(n)], [[] for _ in range(n)]
    for _ in range(m):
        a, b = map(int, input().split())
        a, b = a - 1, b - 1
        graph[a].append(b)
        revGraph[b].append(a)

    dist = [[INF] * 2 for _ in range(n)]
    dist[0][0] = 0
    pq = [(0, 0, 0)]  # (cost, node, reversed)
    while pq:
        cost, node, reversed = heappop(pq)
        if dist[node][reversed] < cost:
            continue

        curGraph = graph if reversed == 0 else revGraph
        for nextNode in curGraph[node]:
            nextCost = cost + 1
            if nextCost < dist[nextNode][reversed]:
                dist[nextNode][reversed] = nextCost
                heappush(pq, (nextCost, nextNode, reversed))

        nextReversed = 1 ^ reversed
        nextCost = cost + x
        if nextCost < dist[node][nextReversed]:
            dist[node][nextReversed] = nextCost
            heappush(pq, (nextCost, node, nextReversed))

    res = min(dist[n - 1])
    print(res)
