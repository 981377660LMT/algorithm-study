# https://atcoder.jp/contests/arc111/tasks/arc111_b

# 输入 n(≤2e5) 和一个 n 行 2 列的矩阵，矩阵元素范围 [1,4e5]。
# 从每行中恰好选一个数，你最多能选出多少个不同的数？


from collections import defaultdict, deque
from typing import DefaultDict, List, Tuple


class UnionFindGraphMap:
    """并查集维护无向图每个连通块的边数和顶点数"""

    __slots__ = ("part", "_parent", "vertex", "edge")

    def __init__(self):
        self._parent = dict()
        self.part = 0
        self.vertex = dict()  # 每个联通块的顶点数
        self.edge = dict()  # 每个联通块的边数

    def add(self, key: int) -> bool:
        if key in self._parent:
            return False
        self._parent[key] = key
        self.vertex[key] = 1
        self.edge[key] = 0
        self.part += 1
        return True

    def find(self, x: int) -> int:
        if x not in self._parent:
            self.add(x)
            return x
        while self._parent.get(x, x) != x:
            self._parent[x] = self._parent[self._parent[x]]
            x = self._parent[x]
        return x

    def union(self, x: int, y: int) -> bool:
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            self.edge[rootX] += 1  # !两个顶点已经在同一个连通块了，这个连通块的边数+1
            return False
        if self.vertex[rootX] > self.vertex[rootY]:
            rootX, rootY = rootY, rootX
        self._parent[rootX] = rootY
        self.vertex[rootY] += self.vertex[rootX]
        self.edge[rootY] += self.edge[rootX] + 1
        self.part -= 1
        return True

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)

    def getGroups(self) -> DefaultDict[int, List[int]]:
        groups = defaultdict(list)
        for key in self._parent:
            root = self.find(key)
            groups[root].append(key)
        return groups

    def getSize(self, key: int) -> int:
        return self.vertex[self.find(key)]

    def __repr__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())

    def __len__(self) -> int:
        return self.part


def selectOneFromEachPair(pairs: List[Tuple[int, int]]) -> int:
    """从每个对中恰好选一个数，最多能选出多少个不同的数.
    对每个大小为m的连通块,树的贡献为m-1,环的贡献为m.
    """
    uf = UnionFindGraphMap()
    for u, v in pairs:
        uf.union(u, v)
    res = 0
    for root, g in uf.getGroups().items():
        isTree = uf.edge[root] == len(g) - 1
        res += len(g) - isTree
    return res


# TODO
def selectOneFromEachPair2(pairs: List[Tuple[int, int]]) -> List[Tuple[int, int]]:
    """从每个对中恰好选一个数，最多能选出多少个不同的数.
    对每个大小为m的连通块,树的贡献为m-1,环的贡献为m.
    """
    adjMap = defaultdict(list)
    visited = dict()
    for i, (u, v) in enumerate(pairs):
        adjMap[u].append((v, i))
        adjMap[v].append((u, i))
        visited[u] = False
        visited[v] = False

    res = []
    for i in visited:
        if visited[i]:
            continue
        vertex, edge = 0, 0
        queue = deque([i])
        while queue:
            cur = queue.popleft()
            if visited[cur]:
                continue
            visited[cur] = True
            vertex += 1
            count = 0
            for next, ei in adjMap[cur]:
                if visited[next] and count == 0:
                    count += 1
                    continue
                if visited[next] and count == 1:
                    edge += 1
                    continue
                queue.append(next)
                edge += 1
                res.append((ei, next if not visited[next] else cur))
        return res


if __name__ == "__main__":
    # https://atcoder.jp/contests/arc111/tasks/arc111_b
    # B - Reversible Cards
    # !给定一些卡片,正反面标有数字,你需要反转一些卡片
    # !使得所有卡片的正面数字种类最多

    # 正反面连边,问题变为选取每个边的一个端点，问最多选取多少个
    # 对每个大小为n的连通块
    # 树: n-1
    # 有环: n
    # 如果要求方案，需要EBCC求边双连通分量.

    n = int(input())
    edges = [tuple(map(int, input().split())) for _ in range(n)]
    print(selectOneFromEachPair2(edges))
