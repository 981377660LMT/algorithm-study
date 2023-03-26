from typing import List
from Rerooting import Rerooting

INF = int(1e18)


def getDistToSpecials(n: int, rawTree: List[List[int]], specials: List[int]) -> List[int]:
    """求`原树`中每个点到特殊点的最远距离."""
    isSpecial = [False] * n
    for v in specials:
        isSpecial[v] = True

    E = int

    def e(root: int) -> E:
        return 0 if isSpecial[root] else -INF

    def op(childRes1: E, childRes2: E) -> E:
        return max(childRes1, childRes2)

    def composition(fromRes: E, parent: int, cur: int, direction: int) -> E:
        """direction: 0: cur -> parent, 1: parent -> cur"""
        return fromRes + 1

    R = Rerooting(n)
    for cur in range(n):
        for next in rawTree[cur]:
            if cur < next:
                R.addEdge(cur, next)
    dp = R.rerooting(e, op, composition)
    return dp
