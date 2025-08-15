from bisect import bisect_left
from typing import Callable, List, Tuple


def levelCount(tree: "Tree") -> Callable[[int, int], int]:
    """查询root的子树中,`绝对深度`为depth的顶点个数."""
    n = len(tree.tree)
    groupByDepth = [[] for _ in range(n)]
    for dep, id in zip(tree.depth, tree.lid):
        groupByDepth[dep].append(id)
    for v in groupByDepth:
        v.sort()

    def f(root: int, depth: int) -> int:
        start, end = tree.id(root)
        pos = groupByDepth[depth]
        count1 = bisect_left(pos, start)
        count2 = bisect_left(pos, end)
        return count2 - count1

    return f


def spanningTreeWeightedSum(tree: "Tree", spanningTree: List[int]) -> int:
    """查询生成树 spanningTree 的边权和."""
    # 按照dfs序排序, 然后计算相邻节点的距离
    order = sorted(range(len(spanningTree)), key=lambda x: tree.lid[spanningTree[x]])
    res = 0
    for pre, cur in zip(order, order[1:] + [order[0]]):
        preNode = spanningTree[pre]
        curNode = spanningTree[cur]
        res += tree.dist(preNode, curNode)
    return res // 2


class Tree:
    __slots__ = (
        "depth",
        "parent",
        "tree",
        "depthWeighted",
        "lid",
        "rid",
        "idToNode",
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
        self.idToNode = [0] * n
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
                return self.idToNode[self.lid[root] - k]
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
                res += self.idToNode[a : b + 1]
            else:
                res += self.idToNode[b : a + 1][::-1]
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

    def subSize(self, v: int, root=-1) -> int:
        """以root为根时,结点v的子树大小"""
        if root == -1:
            return self.rid[v] - self.lid[v]
        if v == root:
            return len(self.tree)
        x = self.jump(v, root, 1)
        if self.isInSubtree(v, x):
            return self.rid[v] - self.lid[v]
        return len(self.tree) - self.rid[x] + self.lid[x]

    def rootedLca(self, u: int, v: int, w: int) -> int:
        """以任意一个点为根, 其他两个点的最近公共祖先."""
        lca1 = self.lca(w, u)
        lca2 = self.lca(w, v)
        if lca1 == lca2:
            return self.lca(u, v)
        return lca1 if self.depth[lca1] > self.depth[lca2] else lca2

    def rootedParent(self, u: int, root: int) -> int:
        return self.jump(u, root, 1)

    def id(self, root) -> Tuple[int, int]:
        """返回 root 的欧拉序区间, 左闭右开, 0-indexed."""
        return self.lid[root], self.rid[root]

    def eid(self, u: int, v: int) -> int:
        """返回返回边 u-v 对应的 边id."""
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
        self.idToNode[self._dfn] = cur
        self._dfn += 1
        heavySon = self._heavySon[cur]
        if heavySon != -1:
            self._markTop(heavySon, top)
            for next, _ in self.tree[cur]:
                if next != heavySon and next != self.parent[cur]:
                    self._markTop(next, next)
        self.rid[cur] = self._dfn


if __name__ == "__main__":

    def demo() -> None:
        tree = Tree(6)
        tree.addEdge(0, 1, 1)
        tree.addEdge(0, 2, 1)
        tree.addEdge(1, 3, 1)
        tree.addEdge(1, 4, 1)
        tree.addEdge(2, 5, 1)
        tree.build(0)
        print(tree.subSize(0, root=4))
        print(tree.rootedLca(0, 1, 4))
        print(tree.rootedLca(0, 1, 3))
        print(tree.rootedParent(0, root=5))

    # https://atcoder.jp/contests/abc202/tasks/abc202_e
    def abc202_e() -> None:
        import sys

        sys.setrecursionlimit(10**6)
        n = int(input())
        parents = [int(x) - 1 for x in input().split()]
        tree = Tree(n)
        for i, p in enumerate(parents):
            tree.addDirectedEdge(p, i + 1, 1)
        tree.build(0)

        lc = levelCount(tree)
        q = int(input())
        for _ in range(q):
            root, dep = map(int, input().split())
            print(lc(root - 1, dep))

    abc202_e()

    class Solution:
        # 3553. 包含给定路径的最小带权子树 II
        # https://leetcode.cn/problems/minimum-weighted-subgraph-with-the-required-paths-ii/solutions/3679978/python-bu-chu-yi-er-de-xie-fa-gua-he-ren-60e2/
        # 带修版本：
        # https://codeforces.com/problemset/problem/176/E
        # https://paste.ubuntu.com/p/bkz8vhyXNM/
        # !维护节点 dfs 序。我们需要知道每个节点dfs序的左端点和右端点。(dfn, idToNode)
        def minimumWeight(self, edges: List[List[int]], queries: List[List[int]]) -> List[int]:
            n = len(edges) + 1
            tree = Tree(n)
            for u, v, w in edges:
                tree.addEdge(u, v, w)
            tree.build(0)

            # res = [0] * len(queries)
            # for i, (u, v, w) in enumerate(queries):
            #     d1, d2 = tree.dist(u, w), tree.dist(v, w)
            #     meet = tree.rootedLca(u, v, w)
            #     overlap = tree.dist(meet, w)
            #     res[i] = d1 + d2 - overlap
            # return res
            return [spanningTreeWeightedSum(tree, [u, v, w]) for u, v, w in queries]
