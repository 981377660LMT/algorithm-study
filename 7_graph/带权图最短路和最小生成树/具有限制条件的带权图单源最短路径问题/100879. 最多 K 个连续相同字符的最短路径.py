# 100879. 最多 K 个连续相同字符的最短路径
# https://leetcode.cn/problems/shortest-path-with-at-most-k-consecutive-identical-characters/description/
# 返回一条从节点 0 到节点 n - 1 的路径的 最小总边权 ，
# 并要求该路径上所有节点标签按顺序 拼接 后，最多包含 k 个 连续相同 字符。
# 如果不存在有效路径，返回 -1。
#
# 分层图最短路
# n<=5e4,k<=50
# 可以记录每个结点

from heapq import heappop, heappush
from typing import List

INF = int(4e18)


class Solution:
    def shortestPath(self, n: int, edges: List[List[int]], labels: str, k: int) -> int:
        adjList = [[] for _ in range(n)]
        for x, y, w in edges:
            adjList[x].append((y, w))

        dist = [[INF] * (k + 1) for _ in range(n)]
        pq = []

        def add(to: int, count: int, d: int) -> None:
            if d < dist[to][count]:
                dist[to][count] = d
                heappush(pq, (d, to, count))

        add(0, 1, 0)
        while pq:
            d, cur, count = heappop(pq)
            if cur == n - 1:
                return d
            if d > dist[cur][count]:
                continue
            for next_, weight in adjList[cur]:
                if labels[cur] != labels[next_]:
                    add(next_, 1, weight + d)
                elif count + 1 <= k:
                    add(next_, count + 1, d + weight)
        return -1

    def shortestPath2(self, n: int, edges: List[List[int]], labels: str, k: int) -> int:
        adjList = [[] for _ in range(n)]
        for u, v, w in edges:
            adjList[u].append((v, w))
        pq = [(0, 0, 1)]
        best = [INF] * n
        while pq:
            d, u, cnt = heappop(pq)
            if u == n - 1:
                return d
            if cnt >= best[u]:
                continue
            best[u] = cnt
            for v, w in adjList[u]:
                if labels[v] == labels[u]:
                    if cnt + 1 <= k:
                        heappush(pq, (d + w, v, cnt + 1))
                else:
                    heappush(pq, (d + w, v, 1))
        return -1
