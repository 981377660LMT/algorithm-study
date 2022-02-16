from collections import  deque
from typing import DefaultDict

WeightedDirectedGraph = DefaultDict[int, DefaultDict[int, int]]


def spfa(graph: WeightedDirectedGraph, start: int, target: int):
    """spfa求单源最短路"""
    n = len(graph)
    dist = [int(1e20)] * n
    dist[start] = 0

    queue = deque()
    queue.append(start)
    visited = [False] * n
    visited[start] = True
    while queue:
        cur = queue.popleft()
        for next, weight in graph[cur].items():
            if dist[cur] + weight < dist[next]:
                dist[next] = dist[cur] + weight
                if not visited[next]:
                    visited[next] = True
                    queue.append(next)

    return dist[target]
