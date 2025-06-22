#
#  Betweenness centrality of undirected unweighted graph (Brandes)
#
#  Description:
#
#  Compute betweenness centrality, defined by
#  f(u) := \sum_{u,t \eq v} |s-t shortest paths that contains v|/|s-t shortest paths|
#
#  Algorithm:
#
#  Brandes's algorithm, O(nm) time, O(m) space.
#
#  References:
#
#  U. Brandes (2001): A faster algorithm for betweenness centrality.
#  Journal of Mathematical Sociology, vol.25, pp.163–177.
#
# O(VE)求出无向图每个点的中介中心性(用于地铁路线优化)
# 如何简单地理解中心度，什么是closeness、betweenness和degree？ - 何燕杰的回答 - 知乎
# https://www.zhihu.com/question/22610633/answer/143644471
#
# 点度中心性（degree，点的度数；微信好友数量）
# 接近中心性（closeness，到其他所有点的最短路的平均长度；去规模化）
# 中介中心性（betweenness，承担最短路桥梁(shortest path/bridge)的次数除以所有的路径数量；社交达人）
#
from collections import deque
from typing import List


def betweennessCentrality(n: int, edges: List[List[int]]) -> List[float]:
    """O(VE)求出无向图每个点的中介中心性."""
    adjList = [[] for _ in range(n)]
    for u, v in edges:
        adjList[u].append(v)
        adjList[v].append(u)

    res = [0.0] * n
    for s in range(n):
        path = []
        counter = [0] * n
        counter[s] = 1
        dist = [-1] * n
        dist[s] = 0
        queue = deque([s])
        while queue:
            cur = queue.popleft()
            path.append(cur)
            for v in adjList[cur]:
                if dist[v] < 0:
                    dist[v] = dist[cur] + 1
                    queue.append(v)
                if dist[v] == dist[cur] + 1:
                    counter[v] += counter[cur]
        cand = [0] * n
        while path:
            cur = path.pop()
            for v in adjList[cur]:
                if dist[v] == dist[cur] + 1:  # 在最短路上
                    cand[cur] += counter[cur] / counter[v] * (1 + cand[v])
            if cur != s:
                res[cur] += cand[cur]
    return res


if __name__ == "__main__":
    # 0 - 1 - 2 - 3
    #     |   |
    #     4 - 5
    n = 6
    edges = [[0, 1], [1, 2], [2, 3], [1, 4], [2, 5], [4, 5]]
    print(betweennessCentrality(n, edges))
