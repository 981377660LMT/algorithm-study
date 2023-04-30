#
#  Transitive Reduction of DAG
#
#  Description:
#  A transitive reduction of a graph G = (V, E) is a
#  a graph H = (V, F) such that transitive closures of
#  H and G are the same.
#  There are possibly many transitive reductions with
#  the fewest edges, and finding one of them is NP-hard.
#  On he other hand, if a graph is directed acyclic,
#  its transitive reduction uniquely exists and can be
#  found in O(nm) time.
#  Note that transitive closure and reduction have the
#  same time complexity on DAG.
#
#  Algorithm:
#  For each vertex u, compute longest path distance from u
#  to v in adj[u]. Then, remove all edges (u,v) with d(u,v) > 1.
#
#  Complexity:
#  O(nm). Usually the coefficient is not so large.
#


# DAG的传递简约(Transitive Reduction)
# 删除有向无环图中的跨层冗余依赖关系，可用于依赖关系处理中对非必要依赖关系的简化。
# O(VE), 在 python 的 networkx 这个库中有这个方法
# 对每个顶点u,计算从u到adj[u]中的v的最长路径距离。然后删除所有满足d(u,v)>1的边(u,v)。

from typing import List


def dagTransitiveReduction(n: int, edges: List[List[int]]) -> List[List[int]]:
    def max(a: int, b: int) -> int:
        return a if a > b else b

    def dfs(cur: int) -> None:
        dist[cur] = 0
        for next in dag[cur]:
            if dist[next] < 0:
                dfs(next)
        order.append(cur)
        for next in order:
            dist[next] = 0
        for i in range(len(order) - 1, -1, -1):
            tmp = order[i]
            for next in dag[tmp]:
                dist[next] = max(dist[next], dist[tmp] + 1)
        dag[cur] = [v for v in dag[cur] if dist[v] <= 1]

    dag = [[] for _ in range(n)]
    for u, v in edges:
        dag[u].append(v)
    order = []
    dist = [-1] * n
    for i in range(n):
        if dist[i] < 0:
            dfs(i)
    return dag


if __name__ == "__main__":
    n = 5
    edges = [[0, 1], [1, 2], [2, 3], [3, 4]]
    print(dagTransitiveReduction(n, edges))
    edges.append([0, 2])
    print(dagTransitiveReduction(n, edges))
