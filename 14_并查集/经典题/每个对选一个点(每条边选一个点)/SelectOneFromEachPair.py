from collections import defaultdict
from typing import DefaultDict, List, Tuple


class SelectOneFromEachPairMap:
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

    def solve(self) -> int:
        """
        从每条边中恰好选一个点, 最多能选出多少个不同的点.
        对每个大小为m的连通块,树的贡献为m-1,环的贡献为m.
        因此答案为`总点数-树的个数`.
        """
        return len(self._data) - self.treeCount

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

    def isTree(self, u: int) -> bool:
        """u所在的联通分量是否为树."""
        root = self.find(u)
        vertex = self.getSize(root)
        return vertex == self._edge[root] + 1

    def countTree(self) -> int:
        """联通分量为树的联通分量个数(孤立点也算树)."""
        return self.treeCount

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


class SelectOneFromEachPairArray:
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

    def solve(self) -> int:
        """
        从每条边中恰好选一个点, 最多能选出多少个不同的点.
        对每个大小为m的连通块,树的贡献为m-1,环的贡献为m.
        因此答案为`总点数-树的个数`.
        """
        return self._n - self.treeCount

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

    def isTree(self, u: int) -> bool:
        """u所在的联通分量是否为树."""
        root = self.find(u)
        vertex = self.getSize(root)
        return vertex == self._edge[root] + 1

    def countTree(self) -> int:
        """联通分量为树的联通分量个数(孤立点也算树)."""
        return self.treeCount

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

    def getSize(self, x: int) -> int:
        return -self._data[self.find(x)]

    def getEdge(self, x: int) -> int:
        return self._edge[self.find(x)]

    def getGroups(self) -> DefaultDict[int, List[int]]:
        groups = defaultdict(list)
        for i in range(self._n):
            root = self.find(i)
            groups[root].append(i)
        return groups

    def __repr__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())


if __name__ == "__main__":
    # https://atcoder.jp/contests/arc111/tasks/arc111_b
    # B - Reversible Cards
    # !给定一些卡片,正反面标有数字,你需要反转一些卡片
    # !使得所有卡片的正面数字种类最多

    # 正反面连边,问题变为选取每个边的一个端点，问最多选取多少个
    # 对每个大小为n的连通块
    # 树: n-1
    # 有环: n

    def reversibleCards(cards: List[Tuple[int, int]]) -> int:
        uf = SelectOneFromEachPairMap()
        for u, v in cards:
            uf.union(u, v)
        return uf.solve()

    # https://atcoder.jp/contests/abc302/tasks/abc302_h
    # 给定一棵树，每个点有两个值。
    # 对于v=1,2,3,...,n，问从点1到点 v的最短路径途径的每个点中，
    # 各选一个数，其不同数的个数的最大值。
    def ballCollector(
        n: int, edges: List[Tuple[int, int]], pairs: List[Tuple[int, int]]
    ) -> List[int]:
        def dfs(cur: int, pre: int) -> None:
            a, b = pairs[cur]
            uf.union(a, b)
            res[cur] = uf.solve()
            for next in adjList[cur]:
                if next != pre:
                    dfs(next, cur)
            uf.undo()

        uf = SelectOneFromEachPairMap()
        adjList = [[] for _ in range(n)]
        for u, v in edges:
            adjList[u].append(v)
            adjList[v].append(u)
        res = [0] * n
        dfs(0, -1)
        return res

    import sys

    sys.setrecursionlimit(int(1e9))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    n = int(input())
    pairs = [tuple(map(int, input().split())) for _ in range(n)]
    edges = []
    for _ in range(n - 1):
        u, v = map(int, input().split())
        edges.append((u - 1, v - 1))
    res = ballCollector(n, edges, pairs)
    print(*res[1:])
