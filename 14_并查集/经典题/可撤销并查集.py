"""这个基本上没有用"""


from typing import DefaultDict, List
from collections import defaultdict


# 使用带撤销操作的并查集，不能使用路径压缩优化（因为路径压缩会改变结构）；
# 为了不超时必须使用按秩合并优化，复杂度nlogn


class UnionFind:
    def __init__(self, n: int):
        self.n = n
        self.count = n
        self.parent = list(range(n))
        self.rank = [1] * n
        self.optStack = []

    def find(self, x: int) -> int:
        """不能使用路径压缩优化"""
        assert 0 <= x < self.n, 'key out of range'
        while self.parent[x] != x:
            x = self.parent[x]
        return x

    def union(self, x: int, y: int) -> bool:
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            self.optStack.append((-1, -1, -1))
            return False

        if self.rank[rootX] > self.rank[rootY]:
            rootX, rootY = rootY, rootX

        self.parent[rootX] = rootY
        self.rank[rootY] += self.rank[rootX]
        self.optStack.append((rootX, rootY, self.rank[rootX]))
        self.count -= 1
        return True

    def revocate(self) -> None:
        """用一个栈记录前面的合并操作，撤销时只需要依次取出栈顶元素做合并操作的逆操作即可"""
        if not self.optStack:
            raise IndexError('no union option to revocate')

        rootX, rootY, rankX = self.optStack.pop()
        if rootX == -1:
            return

        self.parent[rootX] = rootX
        self.rank[rootY] -= rankX
        self.count += 1

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)

    def getGroups(self) -> DefaultDict[int, List[int]]:
        groups = defaultdict(list)
        for key in range(self.n):
            root = self.find(key)
            groups[root].append(key)
        return groups


if __name__ == '__main__':
    uf = UnionFind(10)
    uf.union(2, 4)
    assert uf.isConnected(2, 4)
    uf.revocate()
    assert not uf.isConnected(2, 4)

