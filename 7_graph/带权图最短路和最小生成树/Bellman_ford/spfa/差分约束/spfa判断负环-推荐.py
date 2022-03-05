# 检查是否存在负权环的方法为：记录一个点的入队次数，如果超过V-1次说明存在负权环，
# 因为最短路径上除自身外至多V-1个点，故一个点不可能被更新超过V-1次
from collections import deque
from typing import DefaultDict


WeightedDirectedGraph = DefaultDict[int, DefaultDict[int, int]]


def spfa(graph: WeightedDirectedGraph, start: int) -> bool:
    """spfa判断负权环，返回True表示存在负权环，False表示不存在负权环"""
    n = len(graph)
    dist = [int(1e20)] * n
    dist[start] = 0

    queue = deque()
    queue.append(start)
    visited = [False] * n
    visited[start] = True

    visitedTimes = [0] * n
    visitedTimes[start] = 1

    while queue:
        cur = queue.popleft()
        for next, weight in graph[cur].items():
            if dist[cur] + weight < dist[next]:
                dist[next] = dist[cur] + weight

                # 多了这一段逻辑
                visitedTimes[next] += 1
                if visitedTimes[next] >= n:
                    return True

                if not visited[next]:
                    visited[next] = True
                    queue.append(next)

    return False


if __name__ == '__main__':
    #  x0-x1>=1    x0-1>=x1
    #  x2-x3<=2    x3+2>=x2
    #  x0=x2      x0+0>=x2  x2+0>=x0
    edges = []
    edges.append([0, 1, -1])
    edges.append([3, 2, 2])
    edges.append([0, 2, 0])
    edges.append([2, 0, 0])
    print(spfa(edges, 0))

