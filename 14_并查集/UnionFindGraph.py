# 无向图中：
# 联通分量数(part) = 树的个数(treeCount) + 环的个数
# 边的个数: 总点数 - 树的个数
# 树的性质: 联通分量中的点数 = 边数 + 1
# 环的性质: 联通分量中的点数 = 边数

from collections import defaultdict
from typing import DefaultDict, List


class UnionFindGraphArray:
    """
    可撤销并查集，维护连通分量为树的联通分量个数.
    数组实现.
    """

    __slots__ = ("part", "treeCount", "_n", "_data", "_edge", "_history")

    def __init__(self, n: int):
        self.part = n
        self.treeCount = n  # 联通分量为树的联通分量个数(孤立点也算树)
        self._n = n
        self._data = [-1] * n
        self._edge = [0] * n
        self._history = []  # (root,data,edge)

    def union(self, u: int, v: int) -> bool:
        """添加边对(u,v)."""
        u, v = self.find(u), self.find(v)
        self._history.append((u, self._data[u], self._edge[u]))  # big
        self._history.append((v, self._data[v], self._edge[v]))  # small
        if u == v:
            if self.isTree(u):
                self.treeCount -= 1
            self._edge[u] += 1
            return False
        if self._data[u] > self._data[v]:
            u, v = v, u
        if self.isTree(u) or self.isTree(v):
            self.treeCount -= 1
        self._data[u] += self._data[v]
        self._data[v] = u
        self._edge[u] += self._edge[v] + 1
        self.part -= 1
        return True

    def find(self, u: int) -> int:
        """不能路径压缩."""
        cur = u
        while self._data[cur] >= 0:
            cur = self._data[cur]
        return cur

    def undo(self) -> bool:
        """撤销上一次合并操作，没合并成功也要撤销."""
        if not self._history:
            return False
        small, smallData, smallEdge = self._history.pop()
        big, bigData, bigEdge = self._history.pop()
        self._data[small] = smallData
        self._data[big] = bigData
        self._edge[small] = smallEdge
        self._edge[big] = bigEdge
        if big == small:
            if self.isTree(big):
                self.treeCount += 1
        else:
            if self.isTree(big) or self.isTree(small):
                self.treeCount += 1
            self.part += 1
        return True

    def solve(self) -> int:
        """
        从每条边中恰好选一个点, 最多能选出多少个不同的点.
        对每个大小为m的连通块,树的贡献为m-1,环的贡献为m.
        因此答案为`总点数-树的个数`.
        """
        return self._n - self.treeCount

    def isTree(self, u: int) -> bool:
        """u所在的联通分量是否为树."""
        root = self.find(u)
        vertex = self.getSize(root)
        return vertex == self._edge[root] + 1

    def isCycle(self, u: int) -> bool:
        """u所在的联通分量是否为环."""
        root = self.find(u)
        vertex = self.getSize(root)
        return vertex == self._edge[root]

    def countTree(self) -> int:
        """联通分量为树的联通分量个数(孤立点也算树)."""
        return self.treeCount

    def countCycle(self) -> int:
        """联通分量为环的联通分量个数(孤立点不算环)."""
        return self.part - self.treeCount

    def countEdge(self) -> int:
        return self._n - self.treeCount

    def getSize(self, x: int) -> int:
        """x所在的连通分量的点数."""
        return -self._data[self.find(x)]

    def getEdge(self, x: int) -> int:
        """x所在的连通分量的边数."""
        return self._edge[self.find(x)]

    def getGroups(self) -> DefaultDict[int, List[int]]:
        groups = defaultdict(list)
        for i in range(self._n):
            root = self.find(i)
            groups[root].append(i)
        return groups

    def __repr__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())


class UnionFindGraphMap:
    """
    可撤销并查集，维护连通分量为树的联通分量个数.
    字典实现.
    """

    __slots__ = ("part", "treeCount", "_data", "_edge", "_history")

    def __init__(self):
        self.part = 0
        self.treeCount = 0  # 联通分量为树的联通分量个数(孤立点也算树)
        self._data = dict()
        self._edge = dict()
        self._history = []  # (root,data,edge)

    def union(self, u: int, v: int) -> bool:
        """添加边对(u,v)."""
        u, v = self.find(u), self.find(v)
        self._history.append((u, self._data[u], self._edge[u]))  # big
        self._history.append((v, self._data[v], self._edge[v]))  # small
        if u == v:
            if self.isTree(u):
                self.treeCount -= 1
            self._edge[u] += 1
            return False
        if self._data[u] > self._data[v]:
            u, v = v, u
        if self.isTree(u) or self.isTree(v):
            self.treeCount -= 1
        self._data[u] += self._data[v]
        self._data[v] = u
        self._edge[u] += self._edge[v] + 1
        self.part -= 1
        return True

    def find(self, u: int) -> int:
        """不能路径压缩."""
        if u not in self._data:
            self.add(u)
            return u
        cur = u
        while self._data.get(cur, -1) >= 0:
            cur = self._data[cur]
        return cur

    def undo(self) -> bool:
        """撤销上一次合并操作，没合并成功也要撤销."""
        if not self._history:
            return False
        small, smallData, smallEdge = self._history.pop()
        big, bigData, bigEdge = self._history.pop()
        self._data[small] = smallData
        self._data[big] = bigData
        self._edge[small] = smallEdge
        self._edge[big] = bigEdge
        if big == small:
            if self.isTree(big):
                self.treeCount += 1
        else:
            if self.isTree(big) or self.isTree(small):
                self.treeCount += 1
            self.part += 1
        return True

    def solve(self) -> int:
        """
        从每条边中恰好选一个点, 最多能选出多少个不同的点.
        对每个大小为m的连通块,树的贡献为m-1,环的贡献为m.
        因此答案为`总点数-树的个数`.
        """
        return len(self._data) - self.treeCount

    def isTree(self, u: int) -> bool:
        """u所在的联通分量是否为树."""
        root = self.find(u)
        vertex = self.getSize(root)
        return vertex == self._edge[root] + 1

    def isCycle(self, u: int) -> bool:
        """u所在的联通分量是否为环."""
        root = self.find(u)
        vertex = self.getSize(root)
        return vertex == self._edge[root]

    def countTree(self) -> int:
        """联通分量为树的联通分量个数(孤立点也算树)."""
        return self.treeCount

    def countCycle(self) -> int:
        """联通分量为环的联通分量个数(孤立点不算环)."""
        return self.part - self.treeCount

    def countEdge(self) -> int:
        return len(self._data) - self.treeCount

    def getSize(self, x: int) -> int:
        return -self._data[self.find(x)]

    def getEdge(self, x: int) -> int:
        return self._edge[self.find(x)]

    def getGroups(self) -> DefaultDict[int, List[int]]:
        groups = defaultdict(list)
        for k in self._data:
            root = self.find(k)
            groups[root].append(k)
        return groups

    def add(self, u: int) -> bool:
        """添加点u."""
        if u in self._data:
            return False
        self._data[u] = -1
        self._edge[u] = 0
        self.part += 1
        self.treeCount += 1
        return True

    def __repr__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())


if __name__ == "__main__":
    uf = UnionFindGraphArray(10)

    uf.union(0, 1)
    uf.union(1, 2)
    print(uf.countEdge(), uf.countTree(), uf.part)
    uf.union(0, 2)
    print(uf.countEdge(), uf.countTree(), uf.part)
