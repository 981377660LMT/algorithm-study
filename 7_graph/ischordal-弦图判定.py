#
#  Chordal graph recognition
#
#  Description:
#  !A graph is chordal if any cycle C with |C| >= 4 has a chord.
#  It is also characterized by a perfect elimination ordering (PEO,完美消除序列):
#  An ordering pi is a PEO if, for all u in V,
#  B(u) := {u} \cup { v in N(u) : pi(v) > pi(u) }
#  forms a clique. A graph is chordal if and only if it admits a PEO.
#
#  !Many problems on a chordal graph can be solved in P by using PEO.
#
#  Algorithm:
#  We find a PEO to regocnize a chordal graph.
#  The maximum cardinality search (MCS), which iterates
#  select u in V with largest |N(u) \cap selected|,
#  gives a PEO [Rose and Tarjan'75]; thus we perform MCS and
#  then verify whether it is PEO or not.
#
#  Complexity:
#  !O(n+m) time and space
#
#  References
#  D. J. Rose and R. Endre Tarjan (1975):
#  Algorithmic aspects of vertex elimination.
#  In Proceedings of the 7th Annual ACM Symposium on Theory of Computing (STOC'75),
#  pp. 245--254.
#


# 弦图判定
# 弦(chordal):连接环上不相邻两点之间的边
# 弦图:任意长度>=4的环都有至少一条弦的图

from typing import List


def isChordal(n: int, edges: List[List[int]]) -> bool:
    adjList = [[] for _ in range(n)]
    for u, v in edges:
        adjList[u].append(v)
        adjList[v].append(u)
    top = 0
    rank, score = [-1] * n, [0] * n
    R = [set() for _ in range(n)]
    bucket = [[] for _ in range(n)]  # bucket Dijkstra
    bucket[0] = list(range(n))
    i = 0
    while i < n:
        while not bucket[top]:
            top -= 1
        u = bucket[top].pop()
        if rank[u] >= 0:
            continue
        rank[u] = i
        i += 1
        p = -1
        for v in adjList[u]:
            if rank[v] >= 0:
                R[u].add(v)
                if p == -1 or rank[p] < rank[v]:
                    p = v
            else:
                score[v] += 1
                bucket[score[v]].append(v)
                if top < score[v]:
                    top = score[v]
        if p >= 0:
            for v in R[u]:
                if v != p and v not in R[p]:
                    return False
    return True


def maxCliqueChordal(n: int, edges: List[List[int]]) -> int:
    """弦图的最大团.Not verified."""
    adjList = [[] for _ in range(n)]
    for u, v in edges:
        adjList[u].append(v)
        adjList[v].append(u)
    top = 0
    pi, rank, score = [0] * n, [-1] * n, [0] * n
    bucket = [[] for _ in range(n)]
    bucket[0] = list(range(n))
    i = 0
    while i < n:
        while not bucket[top]:
            top -= 1
        u = bucket[top].pop()
        if rank[u] >= 0:
            continue
        pi[i] = u
        rank[u] = i
        i += 1
        for v in adjList[u]:
            if rank[v] >= 0:
                continue
            score[v] += 1
            bucket[score[v]].append(v)
            top = max(top, score[v])

    res = 0
    visited = set()
    R = [set() for _ in range(n)]
    for u in pi:
        p = -1
        for v in adjList[u]:
            if v not in visited:
                continue
            R[u].add(v)
            if p == -1 or rank[p] < rank[v]:
                p = v
        res = max(res, len(R[u]))
        visited.add(u)
    return res


if __name__ == "__main__":
    n = 3
    edges = [[0, 1], [1, 2], [2, 0]]
    assert isChordal(n, edges)
    n = 4
    edges = [[0, 1], [1, 2], [2, 3], [3, 0]]
    assert not isChordal(n, edges)
    n = 4
    edges = [[0, 1], [1, 2], [2, 3], [3, 0], [3, 1]]
    assert isChordal(n, edges)
    n = 5
    edges = [[0, 1], [1, 2], [2, 3], [3, 0], [4, 0], [4, 1], [4, 2], [4, 3]]
    assert not isChordal(n, edges)
    n = 5
    edges = [[0, 1], [1, 2], [2, 3], [3, 4], [4, 0], [4, 1], [4, 2]]
    assert isChordal(n, edges)

    print(maxCliqueChordal(n, edges))
