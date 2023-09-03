# 区间并查集+重链剖分合并路径上的所有点
# G - Minimum Reachable City-有向树加边后抵达的最小编号
# 给定一颗特殊的有向树,第i条边连接p和i+1且p<i+1.(沿着树边走,编号递增)
# 现在给定q个操作
# !1 u v: 在u和v之间加一条边,保证连边之前可以在最开始的树上从v到达u
# !2 u: 询问从u出发,能到达的最小编号的点是多少.

# 解:
# !操作1可以等价于连接v到u之间的所有`反向边` => 无向图，可以并查集维护连通块内最小编号
# 操作1可以等价于连接v到u之间的所有反向边，然后就构成了一个无向图，
# 需要把这些结点用并查集全部合并，自然可以用支持区间合并的并查集；
# !但问题是这些点在区间上不连续，可以用重链剖分把路径分割成在欧拉序上连续的几段，
# 然后用区间并查集把每一段合并，段和段之间的端点合并，合并后回调更新每个组的最小值


from RangeUnionFind import UnionFindRange


from typing import List, Tuple, Union


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

    def getHeavyChild(self, v: int) -> int:
        """返回结点v的重儿子.如果没有重儿子,返回-1."""
        k = self.lid[v] + 1
        if k == len(self.tree):
            return -1
        w = self._idToNode[k]
        if self.parent[w] == v:
            return w
        return -1

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

    def rootedLca(self, u: int, v: int, root: int) -> int:
        lca1 = self.lca(root, u)
        lca2 = self.lca(root, v)
        lca3 = self.lca(u, v)
        return lca1 ^ lca2 ^ lca3

    def id(self, root) -> Tuple[int, int]:
        """返回 root 的欧拉序区间, 左闭右开, 0-indexed."""
        return self.lid[root], self.rid[root]

    def eid(self, u: int, v: int) -> int:
        """返回返回边 u-v 对应的 欧拉序起点编号, 1 <= eid <= n-1., 0-indexed."""
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
        heavySon = self._heavySon[cur]
        if heavySon != -1:
            self._markTop(heavySon, top)
            for next, _ in self.tree[cur]:
                if next != heavySon and next != self.parent[cur]:
                    self._markTop(next, next)
        self.rid[cur] = self._dfn


def minimumRechableCity(
    n: int,
    edges: List[Tuple[int, int]],
    queries: List[Union[Tuple[int, int, int], Tuple[int, int]]],
) -> List[int]:
    def callback(big: int, small: int) -> None:
        min_ = min(eulerMin[big], eulerMin[small])
        eulerMin[big] = min_

    tree = LCA_HLD(n)
    for u, v in edges:
        tree.addDirectedEdge(u, v, 1)
    tree.build(0)

    lid = tree.lid
    eulerMin = [0] * n  # 欧拉序为i的点所在组里的最小编号
    for i in range(n):
        eulerMin[lid[i]] = i

    res = []
    uf = UnionFindRange(n)
    for op, *args in queries:
        if op == 1:
            u, v = args
            path = tree.getPathDecomposition(u, v, True)
            for a, b in path:  # [a,b]欧拉序区间
                if a > b:
                    a, b = b, a
                uf.unionRange(a, b, callback)
            for (_, pre), (cur, _) in zip(path, path[1:]):
                uf.union(pre, cur, callback)
        else:
            v = args[0]
            euler = lid[v]
            res.append(eulerMin[uf.find(euler)])
    return res


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e9))
    input = lambda: sys.stdin.readline().rstrip("\r\n")
    n = int(input())
    parents = list(map(int, input().split()))
    edges = []
    for i, p in enumerate(parents):
        edges.append((p - 1, i + 1))
    q = int(input())
    queries = []
    for _ in range(q):
        op, *args = map(int, input().split())
        if op == 1:
            queries.append((op, args[0] - 1, args[1] - 1))
        else:
            queries.append((op, args[0] - 1))
    res = minimumRechableCity(n, edges, queries)
    for r in res:
        print(r + 1)
