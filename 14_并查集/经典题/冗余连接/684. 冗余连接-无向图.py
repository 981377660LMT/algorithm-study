from typing import List


# 给定往一棵 n 个节点 (节点值 1～n) 的树中添加一条边后的图。
# 添加的边的两个顶点包含在 1 到 n 中间，且这条附加的边不属于树中已存在的边。
# 图的信息记录于长度为 n 的二维数组 edges ，edges[i] = [ai, bi]
# 表示图中在 ai 和 bi 之间存在一条边。

# 请找出一条可以删去的边，删除后可使得剩余部分是一个有着 n 个节点的树。
# 如果有多个答案，则返回数组 edges 中最后出现的边。


class Solution:
    def findRedundantConnection(self, edges: List[List[int]]) -> List[int]:
        uf = UnionFindArray(len(edges))
        for u, v in edges:
            u, v = u - 1, v - 1
            if uf.isConnected(u, v):
                return [u + 1, v + 1]
            uf.union(u, v)
        return []


# 1. use dfs to find the unique cycle. O(n).
# 2. Actually we can solve a more general problem in O(n+m) time:
#    given a graph, find all edges that could be a non-tree edge.
#    Use dfs to arbitrarily find a spanning tree.
#    For each non-tree edge (u,v), use O(n)-O(1) LCA to find the LCA t of u and v,
#    and mark all edges on the path from u to t and from v to t, using +1/-1 prefix sum tag.
#    Finally perform another dfs and use the tags to find all tree edges that can be contained in a cycle.


from collections import defaultdict
from typing import DefaultDict, List


class UnionFindArray:

    __slots__ = ("n", "part", "parent", "rank")

    def __init__(self, n: int):
        self.n = n
        self.part = n
        self.parent = list(range(n))
        self.rank = [1] * n

    def find(self, x: int) -> int:
        while x != self.parent[x]:
            self.parent[x] = self.parent[self.parent[x]]
            x = self.parent[x]
        return x

    def union(self, x: int, y: int) -> bool:
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            return False
        if self.rank[rootX] > self.rank[rootY]:
            rootX, rootY = rootY, rootX
        self.parent[rootX] = rootY
        self.rank[rootY] += self.rank[rootX]
        self.part -= 1
        return True

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)

    def getGroups(self) -> DefaultDict[int, List[int]]:
        groups = defaultdict(list)
        for key in range(self.n):
            root = self.find(key)
            groups[root].append(key)
        return groups

    def getRoots(self) -> List[int]:
        return list(set(self.find(key) for key in self.parent))

    def __repr__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())

    def __len__(self) -> int:
        return self.part
