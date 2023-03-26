from typing import List
from Rerooting import Rerooting


def findMinHeightTrees(n: int, edges: List[List[int]]) -> List[int]:
    """求树中每个结点为根时到各个结点的最远距离."""
    E = int

    def e(root: int) -> E:
        return 0

    def op(childRes1: E, childRes2: E) -> E:
        return max(childRes1, childRes2)

    def composition(fromRes: E, parent: int, cur: int, direction: int) -> E:
        """direction: 0: cur -> parent, 1: parent -> cur"""
        if direction == 0:  # cur -> parent
            return fromRes + 1
        return fromRes + 1  # parent -> cur

    R = Rerooting(n)
    for u, v in edges:
        R.addEdge(u, v)
    dp = R.rerooting(e=e, op=op, composition=composition, root=0)
    return dp
