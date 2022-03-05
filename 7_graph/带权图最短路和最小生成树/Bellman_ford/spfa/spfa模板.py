from collections import deque
from typing import DefaultDict

WeightedDirectedGraph = DefaultDict[int, DefaultDict[int, int]]

# O(n^2)


def spfa(graph: WeightedDirectedGraph, start: int, target: int):
    """spfa求单源最短路，适用于解决带有负权重的图，是Bellman-ford的常数优化版"""
    n = len(graph)
    dist = [int(1e20)] * n
    dist[start] = 0

    queue = deque()
    queue.append(start)
    visited = [False] * n
    visited[start] = True
    while queue:
        cur = queue.popleft()
        # 更新过谁，就拿谁去更新别人
        for next, weight in graph[cur].items():
            if dist[cur] + weight < dist[next]:
                dist[next] = dist[cur] + weight
                if not visited[next]:
                    visited[next] = True
                    queue.append(next)

    return dist[target]


# spfa可以过很多dijk的题
# 但是网格的图容易卡spfa
# 有边数限制也能用spfa，spfa本质就是让bf不去枚举到不可能会拓展的边，bf能做的spfa都能
