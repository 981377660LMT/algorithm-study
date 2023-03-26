# 6294. 最大价值和与最小价值和的差值
# 求每个点作为根节点时，到叶子节点的最大距离

from typing import List
from Rerooting import Rerooting


class Solution:
    def maxOutput(self, n: int, edges: List[List[int]], price: List[int]) -> int:
        E = int

        def e(root: int) -> E:
            return 0

        def op(childRes1: E, childRes2: E) -> E:
            return max(childRes1, childRes2)

        def composition(fromRes: E, parent: int, cur: int, direction: int) -> E:
            """direction: 0: cur -> parent, 1: parent -> cur"""
            if direction == 0:  # cur -> parent
                return fromRes + price[cur]
            return fromRes + price[parent]  # parent -> cur

        R = Rerooting(n)
        for u, v in edges:
            R.addEdge(u, v)
        dp = R.rerooting(e, op, composition)
        return max(dp)
