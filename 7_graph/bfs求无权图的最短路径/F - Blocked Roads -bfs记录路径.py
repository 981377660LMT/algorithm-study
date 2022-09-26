"""删边最短路+bfs记录路径"""
# !给定一个有向无权稠密图 删除任意一条边 求0到n-1的最短路
# 如果不可到达 则输出-1
# n<=400 m<=n*(n-1)

# 先求出0到n-1的一条最短路 最多V-1条边
# !如果删去的边在这条路上 则跑一次bfs O(V+E)
# !如果不在 则直接输出最短路长度
# !总时间复杂度O((V+E)*V)

import sys
from collections import deque
from typing import List, Tuple

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":

    def bfs(
        start: int, adjList: List[List[Tuple[int, int]]], banEdge: int
    ) -> Tuple[List[int], List[Tuple[int, int]]]:
        """bfs求最短路,并记录任意一条路径"""
        n = len(adjList)
        dist = [INF] * n
        dist[start] = 0
        queue = deque([start])
        pre = [(-1, -1)] * n  # (preNode, preEdge) 记录路径

        while queue:
            cur = queue.popleft()
            for next, edge in adjList[cur]:
                if edge == banEdge:
                    continue
                cand = dist[cur] + 1
                if cand < dist[next]:
                    pre[next] = (cur, edge)  # type: ignore
                    dist[next] = cand
                    queue.append(next)

        return dist, pre

    n, m = map(int, input().split())
    adjList = [[] for _ in range(n)]
    for i in range(m):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        adjList[u].append((v, i))  # !记录每条边的编号

    dist, pre = bfs(0, adjList, -1)
    minDist = dist[n - 1]
    if minDist == INF:
        print(*[-1] * m, sep="\n")
        exit(0)

    path = set()  # 记录最短路上的边id
    cur = n - 1
    while cur != -1:
        path.add(pre[cur][1])
        cur = pre[cur][0]

    for i in range(m):
        if i not in path:  # !不影响结果
            print(minDist)
        else:
            dist, _ = bfs(0, adjList, i)
            print(dist[n - 1] if dist[n - 1] != INF else -1)
