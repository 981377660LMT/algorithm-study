# 每个子树的最小值和次小值


from heapq import nsmallest
from typing import List, Tuple

from Rerooting import Rerooting

INF = int(1e18)


class Solution:
    def minEdgeReversals(self, n: int, edges: List[List[int]], values: List[int]) -> List[int]:
        E = Tuple[int, int]  # 维护每个子树的最小值和次小值

        def e(root: int) -> E:
            return (INF, INF)

        def op(childRes1: E, childRes2: E) -> E:
            a, b = childRes1
            c, d = childRes2
            min1, min2 = nsmallest(2, [a, b, c, d])
            return (min1, min2)

        def composition(fromRes: E, parent: int, cur: int, direction: int) -> E:
            """direction: 0: cur -> parent, 1: parent -> cur"""
            from_ = parent if direction == 0 else cur
            min1, min2 = nsmallest(2, [fromRes[0], fromRes[1], values[from_]])
            return (min1, min2)

        R = Rerooting(n)
        for u, v in edges:
            R.addEdge(u, v)
        dp = R.rerooting(e=e, op=op, composition=composition)
        return dp
