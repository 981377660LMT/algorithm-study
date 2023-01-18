# Definition for a binary tree node.

# 124. 二叉树中的最大路径和
# https://leetcode.cn/problems/binary-tree-maximum-path-sum/
# 路径和 是路径中各节点值的总和。
# 给你一个二叉树的根节点 root ，返回其 最大路径和 。
# !-1000 <= Node.val <= 1000


from typing import Callable, Generic, List, TypeVar, Optional

INF = int(1e18)


class Solution:
    def maxPathSum(self, root: Optional["TreeNode"]) -> int:
        def e(root: int) -> int:
            return 0

        def op(childRes1: int, childRes2: int) -> int:
            return max(childRes1, childRes2)

        def composition(fromRes: int, parent: int, cur: int, direction: int) -> int:
            if direction == 0:  # cur -> parent
                return fromRes + values[cur]
            return fromRes + values[parent]  # parent -> cur

        _, edges, values = treeToGraph1(root)
        n = len(values)
        R = Rerooting(n)
        for u, v in edges:
            R.addEdge(u, v)
        res = R.rerooting(e, op, composition)
        return max([a + b for a, b in zip(values, res)])  # 加上自己的值


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


from typing import List, Optional, Tuple


class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


def treeToGraph1(
    root: Optional["TreeNode"],
) -> Tuple[List[List[int]], List[Tuple[int, int]], List[int]]:
    """二叉树转无向图,返回邻接表,边列表,结点值列表"""

    def dfs(root: Optional["TreeNode"]) -> int:
        if not root:
            return -1
        nonlocal globalId
        values.append(root.val)
        curId = globalId
        globalId += 1
        if root.left:
            childId = dfs(root.left)
            edges.append((curId, childId))
        if root.right:
            childId = dfs(root.right)
            edges.append((curId, childId))
        return curId

    globalId = 0  # 0 - n-1
    edges = []
    values = []
    dfs(root)

    n = len(values)
    adjList = [[] for _ in range(n)]
    for u, v in edges:
        adjList[u].append(v)
        adjList[v].append(u)
    return adjList, edges, values
