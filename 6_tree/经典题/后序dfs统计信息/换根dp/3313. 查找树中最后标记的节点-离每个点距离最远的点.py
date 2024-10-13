# https://leetcode.cn/problems/find-the-last-marked-nodes-in-tree/description/
# 3313. 查找树中最后标记的节点-离每个点距离最远的点


from typing import List, Tuple
from Rerooting import Rerooting


E = Tuple[int, int]  # (最远距离, 最远距离的节点)


class Solution:
    def lastMarkedNodes(self, edges: List[List[int]]) -> List[int]:
        def e(root: int) -> E:
            return (-1, -1)

        def op(childRes1: E, childRes2: E) -> E:
            if childRes1[0] > childRes2[0]:
                return childRes1
            return childRes2

        def composition(fromRes: E, parent: int, cur: int, direction: int) -> E:
            """direction: 0: cur -> parent, 1: parent -> cur"""
            from_ = cur if direction == 0 else parent
            if fromRes[0] == -1:
                return (1, from_)
            return (fromRes[0] + 1, fromRes[1])

        n = len(edges) + 1
        R = Rerooting(n)
        for u, v in edges:
            R.addEdge(u, v)
        dp = R.rerooting(e, op, composition)
        return [p for _, p in dp]
