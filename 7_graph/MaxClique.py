# 最大团
# https://tjkendev.github.io/procon-library/python/graph/bron-kerbosch.html
# O(3^N/3) で最大团を求める。


from typing import List, Tuple


def maxClique(n, edges: List[Tuple[int, int]]) -> List[int]:
    def dfs(V, P, X):
        if not P and not X:
            return V  # V is a maximal clique
        u = next(iter(X or P))  # a pivot vertex
        res = set()
        for v in P - adjList[u]:
            nextRes = dfs(V | {v}, P & adjList[v], X & adjList[v])
            if len(res) < len(nextRes):
                res = nextRes
            P.remove(v)
            X.add(v)
        return res

    adjList = [set() for _ in range(n)]
    for u, v in edges:
        adjList[u].add(v)
        adjList[v].add(u)

    # sort vertices in a degeneracy ordering
    V = list(range(n))
    V.sort(key=lambda x: len(adjList[x]), reverse=True)

    res = set()
    P = set(range(n))
    X = set()
    for v in V:
        cand = dfs({v}, P & adjList[v], X & adjList[v])
        if len(res) < len(cand):
            res = cand
        P.remove(v)
        X.add(v)
    return sorted(res)


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e9))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    # https://atcoder.jp/contests/abc002/tasks/abc002_4
    n, m = map(int, input().split())
    edges = []
    for _ in range(m):
        x, y = map(int, input().split())
        edges.append((x - 1, y - 1))
    res = maxClique(n, edges)
    print(len(res))
