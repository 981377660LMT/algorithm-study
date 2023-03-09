# 给出一个有向图, 问对于每三个不同的点若满足
# !(若顶点A指向顶点B, 顶点B指向顶点C, 则顶点A指向顶点C) 需要添加几条边.
# 不存在不传递三元组(nontansitive triple)
# n<=2000,m<=2000
# !考虑从每个点bfs O(V+E)

# !题目条件可以转换为一个顶点和它可以到达的每一个顶点之间都要有一条边,
# 那这样的话就可以对每一个点跑一遍bfs,
# 累加每个点的边数, 最后在边的总数量中减掉一开始就有的M条边
# !可达:dist!=INF
# 边数很多的话,需要用传递闭包求出每个点的可达点


from collections import deque
from typing import List, Tuple

INF = int(1e18)


def countNontransitiveTriple1(n: int, edges: List[Tuple[int, int]]) -> int:
    """O(V*E)"""

    adjList = [[] for _ in range(n)]
    for u, v in edges:
        adjList[u].append(v)

    count = 0
    for i in range(n):
        dist = [INF] * n
        dist[i] = 0
        queue = deque([i])
        while queue:
            cur = queue.popleft()
            for next in adjList[cur]:
                if dist[next] == INF:
                    dist[next] = dist[cur] + 1
                    queue.append(next)
        count += sum(v != INF for v in dist) - 1  # -1 for itself
    return count - len(edges)


def countNontransitiveTriple2(n: int, edges: List[Tuple[int, int]]) -> int:
    """O(n^3/64)"""
    dp = [1 << i for i in range(n)]
    for u, v in edges:
        dp[v] |= 1 << u  # u -> v
    for k in range(n):
        for i in range(n):
            if dp[i] & (1 << k):
                dp[i] |= dp[k]
    return sum(v.bit_count() for v in dp) - len(edges) - n


if __name__ == "__main__":
    n, m = map(int, input().split())
    edges = []
    for _ in range(m):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        edges.append((u, v))
    # print(countNontransitiveTriple1(n, edges))
    print(countNontransitiveTriple2(n, edges))
