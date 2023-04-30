#
#  Claw-free graph recognition
#
#  Description:
#  A graph is claw-free if it contains no claws, i.e.,
#  o--o--o
# |
# o
#  as a subgraph. In other words, in a claw-free graph,
#  any neighbors N(v) in contains no triangles.
#
#  Algorithm:
#  We test the triangle-freeness in each N(v) of G~.
#  Here, we can use that |N(v)| <= 2 sqrt(m) due to the
#  Turan's theorem (hand-shaking type theorem).
#
#  Complexity:
#  O(n d^3). Here, d <= 2 sqrt(m).
#  Note that d^3 can be reduced to d^\omega by using
#  fast matrix multiplication.
#
# !ClawfreeGraphRecognition-无爪图判定


from typing import List


def isClawFree(n: int, edges: List[List[int]]) -> bool:
    adjList = [[] for _ in range(n)]
    for u, v in edges:
        adjList[u].append(v)
        adjList[v].append(u)
    threshold = 0
    N = [set() for _ in range(n)]
    for s in range(n):
        threshold += len(adjList[s])
        for v in adjList[s]:
            N[s].add(v)
    threshold = 2 * threshold**0.5
    for s in range(n):
        nexts = adjList[s]
        if len(nexts) > threshold:
            return False
        for i in range(len(nexts)):
            for j in range(i + 1, len(nexts)):
                if nexts[j] not in N[nexts[i]]:
                    for k in range(i + 2, len(nexts)):
                        if nexts[k] not in N[nexts[i]] and nexts[k] not in N[nexts[j]]:
                            return False
    return True
