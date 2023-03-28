"""适用于森林的换根dp"""

from collections import defaultdict
from typing import Callable, Generic, List, TypeVar, Dict

T = TypeVar("T")


class RerootingForest(Generic[T]):

    __slots__ = ("adjList", "_n", "_decrement")

    def __init__(self, n: int, decrement: int = 0):
        self.adjList = [[] for _ in range(n)]
        self._n = n
        self._decrement = decrement

    def addEdge(self, u: int, v: int) -> None:
        u -= self._decrement
        v -= self._decrement
        self.adjList[u].append(v)
        self.adjList[v].append(u)

    def rerooting(
        self,
        e: Callable[[int], T],
        op: Callable[[T, T], T],
        composition: Callable[[T, int, int, int], T],
        groupRoot: int,
    ) -> Dict[int, "T"]:
        """groupRoot 是联通分量的根节点."""
        groupRoot -= self._decrement
        assert 0 <= groupRoot < self._n
        parents = defaultdict(lambda: -1)
        order = [groupRoot]

        stack = [groupRoot]
        while stack:
            cur = stack.pop()
            for next in self.adjList[cur]:
                if next == parents[cur]:
                    continue
                parents[next] = cur
                order.append(next)
                stack.append(next)

        dp1 = dict({v: e(v) for v in order})
        dp2 = dict({v: e(v) for v in order})
        for cur in order[::-1]:
            res = e(cur)
            for next in self.adjList[cur]:
                if parents[cur] == next:
                    continue
                dp2[next] = res
                res = op(res, composition(dp1[next], cur, next, 0))
            res = e(cur)
            for next in self.adjList[cur][::-1]:
                if parents[cur] == next:
                    continue
                dp2[next] = op(res, dp2[next])
                res = op(res, composition(dp1[next], cur, next, 0))
            dp1[cur] = res

        for newRoot in order[1:]:
            parent = parents[newRoot]
            dp2[newRoot] = composition(op(dp2[newRoot], dp2[parent]), parent, newRoot, 1)
            dp1[newRoot] = op(dp1[newRoot], dp2[newRoot])
        return dp1


# 310-求树上每个节点到其他节点的最远距离
# 310. 最小高度树
# 在所有可能的树中，具有最小高度的树（即，min(h)）被称为 最小高度树 。
def findMinHeightTrees(n: int, edges: List[List[int]]) -> List[int]:
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

    def dfsForSubSize(cur: int, parent: int) -> int:
        res = 1
        for next in R.adjList[cur]:
            if next != parent:
                res += dfsForSubSize(next, cur)
        subSize[cur] = res
        return res

    R = RerootingForest(n)
    for u, v in edges:
        R.addEdge(u, v)

    subSize = [0] * n
    # dfsForSubSize(0, -1)
    dp = R.rerooting(e=e, op=op, composition=composition, groupRoot=0)
    min_ = min(dp)
    return [i for i in range(n) if dp[i] == min_]
