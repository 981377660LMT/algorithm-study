# 一定在最短路径上的边


import sys

input = lambda: sys.stdin.readline().rstrip("\r\n")

from bisect import bisect_left
from heapq import heappop, heappush
from typing import List, Sequence, Tuple


INF = int(4e18)


def edgesMustOnShortestPath(
    n: int, edges: List[Tuple[int, int, int]], start: int, target: int
) -> List[bool]:
    """
    给定一个n个点m条边的有向带权图.
    对于每条边(u, v, w), 判断是否一定在从start到target的最短路径上.
    """
    adjList = [[] for _ in range(n)]
    rAdjList = [[] for _ in range(n)]
    for _, (u, v, w) in enumerate(edges):
        adjList[u].append((v, w))
        rAdjList[v].append((u, w))

    dist1, dist2 = dijkstra(n, adjList, start), dijkstra(n, rAdjList, target)
    candEdges = []  # [startTime, endTime, edgeId]
    for i, (u, v, w) in enumerate(edges):
        tmp = dist1[u] + w + dist2[v]
        if tmp == dist1[t]:
            candEdges.append([dist1[u], dist1[u] + w, i])  # 可能在最短路径上的边

    allNums = set()
    for a, b, _ in candEdges:
        allNums.add(a)
        allNums.add(b)
    allNums = sorted(allNums)
    for i in range(len(candEdges)):
        candEdges[i][0] = bisect_left(allNums, candEdges[i][0])
        candEdges[i][1] = bisect_left(allNums, candEdges[i][1])
    diff = [0] * (len(allNums) + 1)
    for a, b, _ in candEdges:
        diff[a] += 1
        diff[b] -= 1
    for i in range(1, len(diff)):
        diff[i] += diff[i - 1]

    mustOnShortestPath = [False] * m  # 每条边是否一定在最短路径上
    for i, (time1, _, eid) in enumerate(candEdges):
        mustOnShortestPath[eid] = diff[time1] == 1  # 这个范围的边权只有一条边
    return mustOnShortestPath


def dijkstra(n: int, adjList: Sequence[Sequence[Tuple[int, int]]], start: int) -> List[int]:
    dist = [INF] * n
    dist[start] = 0
    pq = [(0, start)]

    while pq:
        curDist, cur = heappop(pq)
        if dist[cur] < curDist:
            continue
        for next, weight in adjList[cur]:
            cand = dist[cur] + weight
            if cand < dist[next]:
                dist[next] = cand
                heappush(pq, (dist[next], next))
    return dist


if __name__ == "__main__":
    # President and Roads
    # https://www.luogu.com.cn/problem/CF567E
    # 给出一个有向图，从起点走到终点（必须走最短路），问一条边是否一定会被经过.
    # 如果不一定经过它，可以减小它的多少边权使得经过它（边权不能减少到 0），如果可以的话，要使减少的边权最小.
    n, m, s, t = map(int, input().split())
    s -= 1
    t -= 1
    edges = []
    adjList = [[] for _ in range(n)]
    rAdjList = [[] for _ in range(n)]
    for _ in range(m):
        a, b, c = map(int, input().split())
        a -= 1
        b -= 1
        edges.append((a, b, c))
        adjList[a].append((b, c))
        rAdjList[b].append((a, c))

    dist1, dist2 = dijkstra(n, adjList, s), dijkstra(n, rAdjList, t)
    candEdges = []  # [startTime, endTime, edgeId]
    for i, (u, v, w) in enumerate(edges):
        tmp = dist1[u] + w + dist2[v]
        if tmp == dist1[t]:
            candEdges.append([dist1[u], dist1[u] + w, i])  # 可能在最短路径上的边

    allNums = set()
    for a, b, _ in candEdges:
        allNums.add(a)
        allNums.add(b)
    allNums = sorted(allNums)
    for i in range(len(candEdges)):
        candEdges[i][0] = bisect_left(allNums, candEdges[i][0])
        candEdges[i][1] = bisect_left(allNums, candEdges[i][1])
    diff = [0] * (len(allNums) + 1)
    for a, b, _ in candEdges:
        diff[a] += 1
        diff[b] -= 1
    for i in range(1, len(diff)):
        diff[i] += diff[i - 1]

    mustOnShortestPath = [False] * m  # 每条边是否一定在最短路径上
    for i, (time1, _, eid) in enumerate(candEdges):
        mustOnShortestPath[eid] = diff[time1] == 1  # 这个范围的边权只有一条边

    for i, (u, v, w) in enumerate(edges):
        if mustOnShortestPath[i]:
            print("YES")
        else:
            tmp = dist1[u] + w + dist2[v]
            tmp -= dist1[t] - 1
            if tmp >= w:
                print("NO")
            else:
                print("CAN", tmp)
