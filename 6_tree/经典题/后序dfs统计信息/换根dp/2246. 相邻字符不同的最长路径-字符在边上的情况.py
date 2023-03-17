# 请你找出路径上任意`相邻边`都没有分配到相同字符的 最长路径 ，并返回该路径的长度。
# !把边权先变为点权
from typing import Callable, Generic, List, TypeVar

INF = int(1e18)


class Solution:
    def longestPath(self, parent: List[int], s: str) -> int:
        def e(root: int) -> int:
            return 0

        def op(childRes1: int, childRes2: int) -> int:
            return max(childRes1, childRes2)

        def composition(fromRes: int, parent: int, cur: int, direction: int) -> int:
            """direction: 0: cur -> parent, 1: parent -> cur"""
            if direction == 0:  # cur -> parent
                return fromRes + 1 if s[cur] != s[parent] else -INF
            return fromRes + 1 if s[cur] != s[parent] else -INF  # parent -> cur

        n = len(parent)
        R = _Rerooting(n)
        for cur, pre in enumerate(parent):
            if pre == -1:
                continue
            R.addEdge(pre, cur)
        dp = R.rerooting(e, op, composition)
        return max(dp) + 1  # 由于是路径上的点的个数，所以要加1


T = TypeVar("T")


class _Rerooting(Generic[T]):

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
        root=0,
    ) -> List["T"]:
        root -= self._decrement
        assert 0 <= root < self._n
        parents = [-1] * self._n
        order = [root]
        stack = [root]
        while stack:
            cur = stack.pop()
            for next in self.adjList[cur]:
                if next == parents[cur]:
                    continue
                parents[next] = cur
                order.append(next)
                stack.append(next)

        dp1 = [e(i) for i in range(self._n)]
        dp2 = [e(i) for i in range(self._n)]
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
