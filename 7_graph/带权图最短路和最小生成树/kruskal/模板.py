from collections import defaultdict
from typing import DefaultDict, List, Tuple


def kruskal(n: int, edges: List[Tuple[int, int, int]]) -> Tuple[int, List[int]]:
    """Kruskal算法求无向图最小生成树

    Args:
        n (int): 节点`个数`,并查集初始化为(0,1,2,...,n-1)
        edges (List[Tuple[int, int, int]]): 边的列表,每个元素是`(u, v, w)`表示无向边u到v,权重为w

    Returns:
        Tuple[int, List[Tuple[int, int, int]]]: 最小生成树的边权和,组成最小生成树的边的索引

    - 如果不存在,则求出的是森林中的多个最小生成树
    """
    uf = UnionFindArray(n)
    cost, res = 0, []

    edgesWithIndex = sorted([(i, *edge) for i, edge in enumerate(edges)], key=lambda e: e[-1])
    for ei, u, v, w in edgesWithIndex:
        root1, root2 = uf.find(u), uf.find(v)
        if root1 != root2:
            cost += w
            uf.union(root1, root2)
            res.append(ei)
            if len(res) == n - 1:
                return cost, res

    return -1, res


class UnionFindArray:
    """元素是0-n-1的并查集写法,不支持动态添加

    初始化的连通分量个数 为 n
    """

    __slots__ = ("n", "part", "parent", "rank")

    def __init__(self, n: int):
        self.n = n
        self.part = n
        self.parent = list(range(n))
        self.rank = [1] * n

    def find(self, x: int) -> int:
        while self.parent[x] != x:
            self.parent[x] = self.parent[self.parent[x]]
            x = self.parent[x]
        return x

    def union(self, x: int, y: int) -> bool:
        """rank一样时 默认key2作为key1的父节点"""
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
