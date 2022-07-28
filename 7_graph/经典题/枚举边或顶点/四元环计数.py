# https://github.dev/EndlessCheng/codeforces-go/tree/master/copypasta

from typing import List, Tuple


def countCycle4(n: int, edges: List[List[int]]) -> int:
    """无向图四元环计数 边定向 O(E^3/2)"""
    adjList1 = [[] for _ in range(n)]
    deg = [0] * n
    less = lambda u, v: deg[u] < deg[v] or (deg[u] == deg[v] and u < v)

    res = 0
    for u, v in edges:
        adjList1[u].append(v)
        adjList1[v].append(u)
        deg[u] += 1
        deg[v] += 1

    adjList2 = [[] for _ in range(n)]
    for cur, nexts in enumerate(adjList1):
        for next in nexts:
            if less(cur, next):
                adjList2[cur].append(next)

    count = [0] * n
    for cur, nexts in enumerate(adjList1):
        for next1 in nexts:
            for next2 in adjList2[next1]:
                if less(cur, next2):
                    count[cur] += 1
                    res += count[cur]
    return res


if __name__ == "__main__":
    n = 5
    edges = [[0, 1], [1, 2], [2, 3], [3, 4], [4, 0]]
    print(countCycle4(n, edges))
