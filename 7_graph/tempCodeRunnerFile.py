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

from typing import Mapping, Sequence, Union, Iterable


Graph = Union[Mapping[int, Iterable[int]], Sequence[Iterable[int]]]


def isomorphic(g: "Graph", h: "Graph") -> bool:
    """有向图的同构.O(d*n!)."""
    if len(g) != len(h):
        return False

    n = len(g)
    rG, rH = [[] for _ in range(n)], [[] for _ in range(n)]
    for i in range(n):
        for j in g[i]:
            rG[j].append(i)
        for j in h[i]:
            rH[j].append(i)

    hToG = [-1] * n
    gToH = [-1] * n
    inG, outG, inH, outH = [0] * n, [0] * n, [0] * n, [0] * n

    def match(k: int) -> bool:
        if k == n:
            return True
        s = -1
        for t in range(n):
            if gToH[t] >= 0:
                continue
            if s == -1 or inG[s] + outG[s] < inG[t] + outG[t]:
                s = t

        for u in range(n):

            def check(s: int, u: int) -> bool:
                if hToG[u] >= 0:
                    return False
                if inG[s] != inH[u] or outG[s] != outH[u]:
                    return False
                tInG, tOutG, tNewG, tInH, tOutH, tNewH = 0, 0, 0, 0, 0, 0
                for t in g[s]:
                    if gToH[t] >= 0:
                        if gToH[t] not in h[u]:
                            return False
                    else:
                        if inG[t]:
                            tInG += 1
                        if outG[t]:
                            tOutG += 1
                        if not inG[t] and not outG[t]:
                            tNewG += 1

                for v in h[u]:
                    if hToG[v] >= 0:
                        if hToG[v] not in g[s]:
                            return False
                    else:
                        if inH[v]:
                            tInH += 1
                        if outH[v]:
                            tOutH += 1
                        if not inH[v] and not outH[v]:
                            tNewH += 1

                if tInG != tInH or tOutG != tOutH or tNewG != tNewH:
                    return False

                for t in rG[s]:
                    if gToH[t] >= 0:
                        if gToH[t] not in rH[u]:
                            return False
                    else:
                        if inG[t]:
                            tInG += 1
                        if outG[t]:
                            tOutG += 1
                        if not inG[t] and not outG[t]:
                            tNewG += 1

                for v in rH[u]:
                    if hToG[v] >= 0:
                        if hToG[v] not in rG[s]:
                            return False
                    else:
                        if inH[v]:
                            tInH += 1
                        if outH[v]:
                            tOutH += 1
                        if not inH[v] and not outH[v]:
                            tNewH += 1

                if tInG != tInH or tOutG != tOutH or tNewG != tNewH:
                    return False
                return True

            if not check(s, u):
                continue

            inG[s] += 1
            outG[s] += 1
            for t in rG[s]:
                inG[t] += 1
            for t in g[s]:
                outG[t] += 1
            inH[u] += 1
            outH[u] += 1
            for v in rH[u]:
                inH[v] += 1
            for v in h[u]:
                outH[v] += 1

            hToG[u] = s
            gToH[s] = u
            if match(k + 1):
                return True
            hToG[u] = -1
            gToH[s] = -1

            inG[s] -= 1
            outG[s] -= 1
            for t in rG[s]:
                inG[t] -= 1
            for t in g[s]:
                outG[t] -= 1
            inH[u] -= 1
            outH[u] -= 1
            for v in rH[u]:
                inH[v] -= 1
            for v in h[u]:
                outH[v] -= 1

        return False

    return match(0)


def isomorphicNaive(g: "Graph", h: "Graph") -> bool:
    """有向图的同构,简单版.
    n<=8.
    """
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
    # https://atcoder.jp/contests/abc232/tasks/abc232_c
    n, m = map(int, input().split())
    adjList1 = [[] for _ in range(n)]
    for _ in range(m):
        u, v = map(int, input().split())
        adjList1[u - 1].append(v - 1)
        adjList1[v - 1].append(u - 1)

    adjList2 = [[] for _ in range(n)]
    for _ in range(m):
        u, v = map(int, input().split())
        adjList2[u - 1].append(v - 1)
        adjList2[v - 1].append(u - 1)

    print("Yes" if isomorphic(adjList1, adjList2) else "No")
