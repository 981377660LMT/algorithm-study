from heapq import nlargest, nsmallest
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)
# 给你一个 n 个点的 简单有向图 （没有重复边的有向图），节点编号为 0 到 n - 1 。如果这些边是双向边，那么这个图形成一棵 树 。

# 给你一个整数 n 和一个 二维 整数数组 edges ，其中 edges[i] = [ui, vi] 表示从节点 ui 到节点 vi 有一条 有向边 。

# 边反转 指的是将一条边的方向反转，也就是说一条从节点 ui 到节点 vi 的边会变为一条从节点 vi 到节点 ui 的边。

# 对于范围 [0, n - 1] 中的每一个节点 i ，你的任务是分别 独立 计算 最少 需要多少次 边反转 ，从节点 i 出发经过 一系列有向边 ，可以到达所有的节点。

# 请你返回一个长度为 n 的整数数组 answer ，其中 answer[i]表示从节点 i 出发，可以到达所有节点的 最少边反转 次数。
from typing import Callable, Generic, List, TypeVar

T = TypeVar("T")


class Rerooting(Generic[T]):
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


# https://leetcode.cn/problems/minimum-edge-reversals-so-every-node-is-reachable/solutions/2445515/python-huan-gen-dpmo-ban-jiao-cheng-by-9-il9s/
class Solution:
    def minEdgeReversals(self, n: int, edges: List[List[int]], values: List[int]) -> List[int]:
        # def e(root: int) -> E:
        #     return 0

        # def op(childRes1: E, childRes2: E) -> E:
        #     return childRes1 + childRes2

        # def composition(fromRes: E, parent: int, cur: int, direction: int) -> E:
        #     """direction: 0: cur -> parent, 1: parent -> cur"""
        #     from_ = parent if direction == 0 else cur
        #     to = cur if direction == 0 else parent
        #     return fromRes + ((from_, to) not in ok)

        # 维护每个子树的最小值和次小值
        E = Tuple[int, int]

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
        return R.rerooting(e=e, op=op, composition=composition)


print(
    Solution().minEdgeReversals(
        n=6, edges=[[0, 1], [1, 3], [2, 3], [4, 0], [4, 5]], values=[1, 2, 3, 4, 5, 6]
    )
)
