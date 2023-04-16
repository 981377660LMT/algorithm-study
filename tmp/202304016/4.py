from functools import lru_cache
from typing import List, Tuple, Optional
from collections import defaultdict, Counter, deque
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 现有一棵无向、无根的树，树中有 n 个节点，按从 0 到 n - 1 编号。给你一个整数 n 和一个长度为 n - 1 的二维整数数组 edges ，其中 edges[i] = [ai, bi] 表示树中节点 ai 和 bi 之间存在一条边。

# 每个节点都关联一个价格。给你一个整数数组 price ，其中 price[i] 是第 i 个节点的价格。

# 给定路径的 价格总和 是该路径上所有节点的价格之和。

# 另给你一个二维整数数组 trips ，其中 trips[i] = [starti, endi] 表示您从节点 starti 开始第 i 次旅行，并通过任何你喜欢的路径前往节点 endi 。

# 在执行第一次旅行之前，你可以选择一些 非相邻节点 并将价格减半。


# 返回执行所有旅行的最小价格总和。
# 1 <= n <= 50
# edges.length == n - 1
# 0 <= ai, bi <= n - 1
# edges 表示一棵有效的树
# price.length == n
# price[i] 是一个偶数
# 1 <= price[i] <= 1000
# 1 <= trips.length <= 100
# 0 <= starti, endi <= n - 1

from typing import List, Tuple


# class Tree
class LCA_HLD:
    __slots__ = (
        "depth",
        "parent",
        "tree",
        "depthWeighted",
        "lid",
        "rid",
        "_idToNode",
        "_top",
        "_heavySon",
        "_dfn",
    )

    def __init__(self, n: int) -> None:
        self.depth = [0] * n
        self.parent = [-1] * n
        self.tree = [[] for _ in range(n)]
        self.depthWeighted = [0] * n
        self.lid = [0] * n
        self.rid = [0] * n
        self._idToNode = [0] * n
        self._top = [0] * n
        self._heavySon = [0] * n
        self._dfn = 0

    def build(self, root=-1) -> None:
        if root != -1:
            self._build(root, -1, 0, 0)
            self._markTop(root, root)
            return
        n = len(self.tree)
        for i in range(n):
            if self.parent[i] == -1:
                self._build(i, -1, 0, 0)
                self._markTop(i, i)

    def addEdge(self, from_: int, to: int, weight: int) -> None:
        self.tree[from_].append((to, weight))
        self.tree[to].append((from_, weight))

    def addDirectedEdge(self, from_: int, to: int, weight: int) -> None:
        self.tree[from_].append((to, weight))

    def lca(self, u: int, v: int) -> int:
        while True:
            if self.lid[u] > self.lid[v]:
                u, v = v, u
            if self._top[u] == self._top[v]:
                return u
            v = self.parent[self._top[v]]

    def dist(self, u: int, v: int, weighted=True) -> int:
        if weighted:
            return (
                self.depthWeighted[u]
                + self.depthWeighted[v]
                - 2 * self.depthWeighted[self.lca(u, v)]
            )
        return self.depth[u] + self.depth[v] - 2 * self.depth[self.lca(u, v)]

    def kthAncestor(self, root: int, k: int) -> int:
        """k:0-indexed;如果不存在,返回-1.
        kthAncestor(root,0) = root
        """
        if k > self.depth[root]:
            return -1
        while True:
            u = self._top[root]
            if self.lid[root] - k >= self.lid[u]:
                return self._idToNode[self.lid[root] - k]
            k -= self.lid[root] - self.lid[u] + 1
            root = self.parent[u]

    def jump(self, from_: int, to: int, step: int) -> int:
        """
        从 from 节点跳向 to 节点,跳过 step 个节点(0-indexed).
        返回跳到的节点,如果不存在这样的节点,返回-1.
        """
        if step == 1:
            if from_ == to:
                return -1
            if self.isInSubtree(to, from_):
                return self.kthAncestor(to, self.depth[to] - self.depth[from_] - 1)
            return self.parent[from_]
        c = self.lca(from_, to)
        dac = self.depth[from_] - self.depth[c]
        dbc = self.depth[to] - self.depth[c]
        if step > dac + dbc:
            return -1
        if step <= dac:
            return self.kthAncestor(from_, step)
        return self.kthAncestor(to, dac + dbc - step)

    def getPath(self, from_: int, to: int) -> List[int]:
        res = []
        composition = self.getPathDecomposition(from_, to, True)
        for a, b in composition:
            if a <= b:
                res += self._idToNode[a : b + 1]
            else:
                res += self._idToNode[b : a + 1][::-1]
        return res

    def getPathDecomposition(self, from_: int, to: int, vertex: bool) -> List[Tuple[int, int]]:
        """返回沿着`路径顺序`的 [起点,终点] 的 欧拉序 `左闭右闭` 数组.
        注意不一定沿着欧拉序, 但是沿着路径顺序.
        """
        up, down = [], []
        while True:
            if self._top[from_] == self._top[to]:
                break
            if self.lid[from_] < self.lid[to]:
                down.append((self.lid[self._top[to]], self.lid[to]))
                to = self.parent[self._top[to]]
            else:
                up.append((self.lid[from_], self.lid[self._top[from_]]))
                from_ = self.parent[self._top[from_]]
        offset = 1 ^ vertex
        if self.lid[from_] < self.lid[to]:
            down.append((self.lid[from_] + offset, self.lid[to]))
        elif self.lid[to] + offset <= self.lid[from_]:
            up.append((self.lid[from_], self.lid[to] + offset))
        return up + down[::-1]

    def isInSubtree(self, child: int, root: int) -> bool:
        """child是否在root的子树中(child和root不能相等)"""
        return self.lid[root] <= self.lid[child] < self.rid[root]

    def subtreeSize(self, v: int, root=-1) -> int:
        """以root为根时,结点v的子树大小"""
        if root == -1:
            return self.rid[v] - self.lid[v]
        if v == root:
            return len(self.tree)
        x = self.jump(v, root, 1)
        if self.isInSubtree(v, x):
            return self.rid[v] - self.lid[v]
        return len(self.tree) - self.rid[x] + self.lid[x]

    def id(self, root) -> Tuple[int, int]:
        """返回 root 的欧拉序区间, 左闭右开, 0-indexed."""
        return self.lid[root], self.rid[root]

    def eid(self, u: int, v: int) -> int:
        """返回边 u-v 对应的 欧拉序起点编号, 0-indexed."""
        id1, id2 = self.lid[u], self.lid[v]
        return id1 if id1 > id2 else id2

    def _build(self, cur: int, pre: int, dep: int, dist: int) -> int:
        subSize, heavySize, heavySon = 1, 0, -1
        for next, weight in self.tree[cur]:
            if next != pre:
                nextSize = self._build(next, cur, dep + 1, dist + weight)
                subSize += nextSize
                if nextSize > heavySize:
                    heavySize, heavySon = nextSize, next
        self.depth[cur] = dep
        self.depthWeighted[cur] = dist
        self.parent[cur] = pre
        self._heavySon[cur] = heavySon
        return subSize

    def _markTop(self, cur: int, top: int) -> None:
        self._top[cur] = top
        self.lid[cur] = self._dfn
        self._idToNode[self._dfn] = cur
        self._dfn += 1
        if self._heavySon[cur] != -1:
            self._markTop(self._heavySon[cur], top)
            for next, _ in self.tree[cur]:
                if next != self._heavySon[cur] and next != self.parent[cur]:
                    self._markTop(next, next)
        self.rid[cur] = self._dfn


# TODO
# !getPath 修改dfsPath的模板(起点到终点)
# !维护点权的树剖


class Solution:
    def minimumTotalPrice(
        self, n: int, edges: List[List[int]], price: List[int], trips: List[List[int]]
    ) -> int:
        adjList = [[] for _ in range(n)]
        for u, v in edges:
            adjList[u].append(v)
            adjList[v].append(u)

        lca = LCA_HLD(n)
        for u, v in edges:
            lca.addEdge(u, v, 1)
        lca.build(0)

        hit = [0] * n
        for start, end in trips:
            path = lca.getPath(start, end)
            for v in path:
                hit[v] += 1
        price = [p * hit[i] for i, p in enumerate(price)]

        # !树形dp TODO (不相邻选择结点模型)
        @lru_cache(None)
        def dfs(cur: int, pre: int, preR: bool) -> int:
            res = 0
            for next in adjList[cur]:
                if next == pre:
                    continue
                cand1 = dfs(next, cur, False) + price[next]
                cand2 = INF
                if not preR:
                    cand2 = dfs(next, cur, True) + price[next] // 2
                res += min(cand1, cand2)
            return res

        res = min(dfs(0, -1, False) + price[0], dfs(0, -1, True) + price[0] // 2)
        dfs.cache_clear()
        return res


# 5
# [[0,2],[1,4],[2,3],[3,4]]
# [30,2,2,4,32]
# [[3,3],[3,2],[2,1],[2,4],[3,4],[2,1],[2,4],[4,0],[3,2],[0,4],[3,3],[1,4],[2,0],[0,4],[4,1],[1,2],[2,2],[2,4],[3,4],[1,3],[4,2],[1,0],[3,2],[0,0],[0,4],[4,4],[1,3],[2,4],[4,2],[0,4],[2,0],[4,2],[0,0],[2,1],[4,4],[3,0],[1,1],[1,2],[1,3],[1,4]]
print(
    Solution().minimumTotalPrice(
        5,
        [[0, 2], [1, 4], [2, 3], [3, 4]],
        [30, 2, 2, 4, 32],
        [
            [3, 3],
            [3, 2],
            [2, 1],
            [2, 4],
            [3, 4],
            [2, 1],
            [2, 4],
            [4, 0],
            [3, 2],
            [0, 4],
            [3, 3],
            [1, 4],
            [2, 0],
            [0, 4],
            [4, 1],
            [1, 2],
            [2, 2],
            [2, 4],
            [3, 4],
            [1, 3],
            [4, 2],
            [1, 0],
            [3, 2],
            [0, 0],
            [0, 4],
            [4, 4],
            [1, 3],
            [2, 4],
            [4, 2],
            [0, 4],
            [2, 0],
            [4, 2],
            [0, 0],
            [2, 1],
            [4, 4],
            [3, 0],
            [1, 1],
            [1, 2],
            [1, 3],
            [1, 4],
        ],
    )
)
print(
    Solution().minimumTotalPrice(
        9,
        [[2, 5], [3, 4], [4, 1], [1, 7], [6, 7], [7, 0], [0, 5], [5, 8]],
        [4, 4, 6, 4, 2, 4, 2, 14, 8],
        [
            [1, 5],
            [2, 7],
            [4, 3],
            [1, 8],
            [2, 8],
            [4, 3],
            [1, 5],
            [1, 4],
            [2, 1],
            [6, 0],
            [0, 7],
            [8, 6],
            [4, 0],
            [7, 5],
            [7, 5],
            [6, 0],
            [5, 1],
            [1, 1],
            [7, 5],
            [1, 7],
            [8, 7],
            [2, 3],
            [4, 1],
            [3, 5],
            [2, 5],
            [3, 7],
            [0, 1],
            [5, 8],
            [5, 3],
            [5, 2],
        ],
    )
)
