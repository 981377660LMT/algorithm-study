# k短路(K-Shortest Walk)
# n,m,k<=3e5
# 给定有向图
# 输出从s到t的k条最短路长度
from heapq import heappop, heappush
from typing import List, Tuple


INF = int(1e18)


def main():
    n, m, start, end, k = map(int, input().split())
    adjList = [[] for _ in range(n)]
    for _ in range(m):
        u, v, w = map(int, input().split())
        adjList[u].append((w, v))
    res = shortest_paths(adjList, start, end, k)  # !求出来可能不足k条
    for d in res:
        print(d)
    for _ in range(k - len(res)):
        print(-1)


# https://judge.yosupo.jp/submission/87055
# Eppstein's algorithm
def shortest_paths(graph: List[List[Tuple[int, int]]], start: int, end: int, k: int) -> List[int]:
    n = len(graph)
    revg = [[] for _ in range(n)]
    for u in range(n):
        for w, v in graph[u]:
            revg[v].append((w, u))
    d, p = dijkstra(revg, end)
    if d[start] == INF:
        return []

    t = [[] for _ in range(n)]
    for u in range(n):
        if p[u] != -1:
            t[p[u]].append(u)

    h = [None] * n
    q = [end]
    for u in q:
        seenp = False
        for w, v in graph[u]:
            if d[v] == INF:
                continue
            c = w + d[v] - d[u]
            if not seenp and v == p[u] and c == 0:
                seenp = True
                continue
            h[u] = EHeap.insert(h[u], c, v)
        for v in t[u]:
            h[v] = h[u]
            q.append(v)

    res = [d[start]]
    if not h[start]:
        return res

    q = [(d[start] + h[start].key, h[start])]
    while q and len(res) < k:
        cd, ch = heappop(q)
        res.append(cd)
        if h[ch.value]:
            heappush(q, (cd + h[ch.value].key, h[ch.value]))
        if ch.left:
            heappush(q, (cd + ch.left.key - ch.key, ch.left))
        if ch.right:
            heappush(q, (cd + ch.right.key - ch.key, ch.right))
    return res


def dijkstra(g, src):
    n = len(g)
    d, p = [INF] * n, [-1] * n
    d[src] = 0
    q = [(0, src)]
    while q:
        du, u = heappop(q)
        if du != d[u]:
            continue
        for w, v in g[u]:
            if du + w < d[v]:
                d[v] = du + w
                p[v] = u
                heappush(q, (d[v], v))
    return d, p


# Leftist heap
class EHeap:
    __slots__ = "rank", "key", "value", "left", "right"

    def __init__(self, rank, key, value, left, right):
        self.rank = rank
        self.key = key
        self.value = value
        self.left = left
        self.right = right

    @staticmethod
    def insert(a, k, v):
        if not a or k < a.key:
            return EHeap(1, k, v, a, None)
        l, r = a.left, EHeap.insert(a.right, k, v)
        if not l or r.rank > l.rank:
            l, r = r, l
        return EHeap((r.rank if r else 0) + 1, a.key, a.value, l, r)

    def __lt__(self, _):
        return False


if __name__ == "__main__":
    main()
