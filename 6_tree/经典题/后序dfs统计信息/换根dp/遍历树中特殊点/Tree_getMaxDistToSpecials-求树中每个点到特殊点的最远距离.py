from typing import List, Tuple
from Rerooting import Rerooting

INF = int(1e18)


def getDistMaxToSpecials(n: int, edges: List[Tuple[int, int]], specials: List[int]) -> List[int]:
    """求`原树`中每个点到可达的(连通)所有特殊点的距离的最大值."""
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
    for cur, next in edges:
        R.addEdge(cur, next)
    dp = R.rerooting(e, op, composition)
    return dp
