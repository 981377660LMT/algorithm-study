#
#  Graph Isomorphism
#
#  Description:
#  Directed graphs G and H are isomprhic
#  if there is a bijection between these vertices such that
#  (u,v) in E(G) iff (pi(u),pi(v)) in E(H).
#  It is open that testing graph isomorphism is NP-complete.
#
#  Algorithm:
#  Cordella-Foggia-Sansone-Vento's algorithm (aka, VF2 algorithm)
#  with large degree heuristics.
#
#  Complexity:
#  O(d n!) time, O(n) space.
#
#  Verify:
#  SPOJ401
#
#  References:
#  L. P. Cordella, P. Foggia, C. Sansone, and M. Vento (2004):
#  A (sub)graph isomorphism algorithm for matching large graphs.
#  IEEE Transactions on Pattern Analysis and Machine Intelligence,
#  vol.28, no.10, pp.1367--1372.

#  有向图的同构

from typing import Mapping, Sequence, Union


Graph = Union[Mapping[int, Sequence[int]], Sequence[Sequence[int]]]


def isomorphicNaive(g: "Graph", h: "Graph") -> bool:
    """有向图的同构,简单版."""
    if len(g) != len(h):
        return False
    n = len(g)
    gToH = [-1] * n
    hToG = [-1] * n

    def match(s: int) -> bool:
        if s == n:
            for t in range(n):
                nbhG = [False] * n
                nbhH = [False] * n
                for p in g[t]:
                    nbhG[gToH[p]] = True
                for q in h[gToH[t]]:
                    nbhH[q] = True
                if nbhG != nbhH:
                    return False
            return True
        for u in range(n):
            if hToG[u] >= 0:
                continue
            gToH[s] = u
            hToG[u] = s
            if match(s + 1):
                return True
            gToH[s] = -1
            hToG[u] = -1
        return False

    return match(0)


if __name__ == "__main__":
    n, m = map(int, input().split())
    adjList1 = [[] for _ in range(n)]
    for _ in range(m):
        u, v = map(int, input().split())
        adjList1[u - 1].append(v - 1)
        adjList1[v - 1].append(u - 1)
    n, m = map(int, input().split())
    adjList2 = [[] for _ in range(n)]
    for _ in range(m):
        u, v = map(int, input().split())
        adjList2[u - 1].append(v - 1)
        adjList2[v - 1].append(u - 1)

    print("Yes" if isomorphicNaive(adjList1, adjList2) else "No")
