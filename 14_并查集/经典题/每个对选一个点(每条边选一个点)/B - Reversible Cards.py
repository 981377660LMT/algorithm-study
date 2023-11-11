# https://atcoder.jp/contests/arc111/tasks/arc111_b

# 输入 n(≤2e5) 和一个 n 行 2 列的矩阵，矩阵元素范围 [1,4e5]。
# 从每行中恰好选一个数，你最多能选出多少个不同的数？


from collections import defaultdict, deque
from typing import DefaultDict, List, Tuple


class UnionFindGraphArray:
    """并查集维护无向图每个连通块的边数和顶点数."""

    __slots__ = ("part", "_parent", "vertex", "edge")

    def __init__(self, n: int):
        self.part = n
        self.vertex = [1] * n  # 每个联通块的顶点数
        self.edge = [0] * n  # 每个联通块的边数
        self._parent = list(range(n))

    def find(self, x: int) -> int:
        while self._parent[x] != x:
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
        for i in range(len(self._parent)):
            root = self.find(i)
            groups[root].append(i)
        return groups

    def getSize(self, key: int) -> int:
        return self.vertex[self.find(key)]

    def getEdge(self, key: int) -> int:
        return self.edge[self.find(key)]

    def __repr__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())


def selectOneFromEachPair(pairs: List[Tuple[int, int]]) -> int:
    """从每个对中恰好选一个数，最多能选出多少个不同的数.
    对每个大小为m的连通块,树的贡献为m-1,环的贡献为m.
    """
    id = dict()
    for u, v in pairs:
        id.setdefault(u, len(id))
        id.setdefault(v, len(id))

    n = len(id)
    uf = UnionFindGraphArray(n)
    for u, v in pairs:
        uf.union(id[u], id[v])

    res = 0
    for root, g in uf.getGroups().items():
        isTree = uf.edge[root] == len(g) - 1
        res += len(g) - isTree
    return res


def selectOneFromEachPair2(pairs: List[Tuple[int, int]]) -> List[int]:
    """
    从每个对中恰好选一个数，最多能选出多少个不同的数.
    !返回每个pair选取的数(方案).

    环:从叶子结点开始拓扑排序,然后从每个未被访问的结点开始dfs.(这里的环包括自环)
    树:从叶子结点开始向中心拓扑排序.
    """

    def solveTree(group: List[int]) -> None:
        queue = deque(u for u in group if deg[u] == 1)
        while queue:
            cur = queue.popleft()
            visited[cur] = True
            onCycle[cur] = False
            for next, ei in adjList[cur]:
                if visited[next]:
                    continue
                select[ei] = cur
                deg[next] -= 1
                if deg[next] == 1:
                    queue.append(next)

    def solveCycle(group: List[int]) -> None:
        queue = deque(u for u in group if deg[u] == 1)
        while queue:
            cur = queue.popleft()
            visited[cur] = True
            onCycle[cur] = False
            for next, ei in adjList[cur]:
                if visited[next]:
                    continue
                select[ei] = cur
                deg[next] -= 1
                if deg[next] == 1:
                    queue.append(next)
        for v in group:
            if not visited[v]:
                _dfsCycle(v, -1)

    def _dfsCycle(cur: int, pre: int) -> None:
        if visited[cur]:
            return
        visited[cur] = True
        for next, ei in adjList[cur]:
            if next == pre or not onCycle[next]:
                continue
            select[ei] = cur
            _dfsCycle(next, cur)
            break

    id = dict()
    for u, v in pairs:
        id.setdefault(u, len(id))
        id.setdefault(v, len(id))

    n = len(id)
    uf = UnionFindGraphArray(n)
    adjList = [[] for _ in range(n)]  # (next, ei)
    deg = [0] * n
    for i, (u, v) in enumerate(pairs):
        u, v = id[u], id[v]
        uf.union(u, v)
        adjList[u].append((v, i))
        adjList[v].append((u, i))
        deg[u] += 1
        deg[v] += 1

    select = [-1] * len(pairs)  # 每个pair选取的数
    visited = [False] * n
    onCycle = [True] * n
    for root, g in uf.getGroups().items():
        if uf.edge[root] == len(g) - 1:
            solveTree(g)
        else:
            solveCycle(g)

    rid = {v: k for k, v in id.items()}
    for i, v in enumerate(select):
        if v == -1:
            select[i] = pairs[i][0]  # 任意选一个,选第一个
        else:
            select[i] = rid[v]
    return select


if __name__ == "__main__":
    # https://atcoder.jp/contests/arc111/tasks/arc111_b
    # B - Reversible Cards
    # !给定一些卡片,正反面标有数字,你需要反转一些卡片
    # !使得所有卡片的正面数字种类最多

    # 正反面连边,问题变为选取每个边的一个端点，问最多选取多少个
    # 对每个大小为n的连通块
    # 树: n-1
    # 有环: n

    import sys

    sys.setrecursionlimit(int(1e6))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    n = int(input())
    edges = [tuple(map(int, input().split())) for _ in range(n)]
    # print(selectOneFromEachPair(edges))
    print(len(set(selectOneFromEachPair2(edges))))
