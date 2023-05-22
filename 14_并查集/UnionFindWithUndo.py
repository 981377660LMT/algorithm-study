"""
UnionFindWithUndo/RevocableUnionFind/RollbackUnionFind
可撤销并查集/带撤销操作的并查集

不能使用路径压缩优化（因为路径压缩会改变结构）；
为了不超时必须使用按秩合并优化,复杂度nlogn

配合回溯的场景使用
撤销相当于弹出栈顶元素

很少用到撤销操作，因为并查集的撤销可以变成倒着合并

应用场景:
可持久化并查集的离线处理
!在树上(版本之间)dfs 递归时要union结点 回溯时候需要撤销的场合
"""

# 可撤销并查集(时间旅行)

# API:
# RollbackUnionFind(int sz)：

# Union(int x, int y)：
# Find(int k)：
# IsConnected(int x, int y)：

# Undo()：撤销上一次合并操作，没合并成功也要撤销.

# Snapshot():内部保存当前状态。
#  !Snapshot() 之后可以调用 Rollback(-1) 回滚到这个状态.
# Rollback(state = -1)：回滚到指定状态。
#   state等于-1时，会回滚到snapshot()中保存的状态。
#   否则，会回滚到指定的state次union调用时的状态。
# GetState()：
#   返回当前状态，即union调用的次数。

from collections import defaultdict
from typing import DefaultDict, Generic, Hashable, Iterable, List, Optional, TypeVar


class UnionFindArrayWithUndo:
    __slots__ = ("n", "part", "_parent", "_rank", "_optStack")

    def __init__(self, n: int):
        self.n = n
        self.part = n
        self._parent = list(range(n))
        self._rank = [1] * n
        self._optStack = []

    def find(self, x: int) -> int:
        """不能使用路径压缩优化"""
        while self._parent[x] != x:
            x = self._parent[x]
        return x

    def union(self, x: int, y: int) -> bool:
        """x所在组合并到y所在组"""
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            self._optStack.append((-1, -1, -1))
            return False
        if self._rank[rootX] > self._rank[rootY]:
            rootX, rootY = rootY, rootX
        self._parent[rootX] = rootY
        self._rank[rootY] += self._rank[rootX]
        self.part -= 1
        self._optStack.append((rootX, rootY, self._rank[rootX]))
        return True

    def undo(self) -> None:
        """
        用一个栈记录前面的合并操作，
        撤销时要依次取出栈顶元素做合并操作的逆操作.
        !没合并成功也要撤销.
        """
        if not self._optStack:
            return
        rootX, rootY, rankX = self._optStack.pop()
        if rootX == -1:
            return
        self._parent[rootX] = rootX
        self._rank[rootY] -= rankX
        self.part += 1

    def reset(self) -> None:
        while self._optStack:
            self.undo()

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)

    def getGroups(self) -> DefaultDict[int, List[int]]:
        groups = defaultdict(list)
        for key in range(self.n):
            root = self.find(key)
            groups[root].append(key)
        return groups


T = TypeVar("T", bound=Hashable)


class UnionFindMapWithUndo(Generic[T]):
    """
    带撤销操作的并查集

    不能使用路径压缩优化（因为路径压缩会改变结构）；
    为了不超时必须使用按秩合并优化,复杂度nlogn
    """

    __slots__ = ("part", "_parent", "_rank", "_optStack")

    def __init__(self, iterable: Optional[Iterable[T]] = None):
        self.part = 0
        self._parent = dict()
        self._rank = dict()
        self._optStack = []
        for item in iterable or []:
            self.add(item)

    def find(self, key: T) -> T:
        """不能使用路径压缩优化"""
        if key not in self._parent:
            self.add(key)
            return key
        while self._parent.get(key, key) != key:
            key = self._parent[key]
        return key

    def union(self, key1: T, key2: T) -> bool:
        """rank一样时 默认key2作为key1的父节点"""
        root1 = self.find(key1)
        root2 = self.find(key2)
        if root1 == root2:
            self._optStack.append((-1, -1, -1))
            return False
        if self._rank[root1] > self._rank[root2]:
            root1, root2 = root2, root1
        self._parent[root1] = root2
        self._rank[root2] += self._rank[root1]
        self.part -= 1
        self._optStack.append((root1, root2, self._rank[root1]))
        return True

    def undo(self) -> None:
        """
        用一个栈记录前面的合并操作，
        撤销时要依次取出栈顶元素做合并操作的逆操作.
        !没合并成功也要撤销.
        """
        if not self._optStack:
            return
        root1, root2, rank1 = self._optStack.pop()
        if root1 == -1:
            return
        self._parent[root1] = root1
        self._rank[root2] -= rank1
        self.part += 1

    def reset(self) -> None:
        while self._optStack:
            self.undo()

    def isConnected(self, key1: T, key2: T) -> bool:
        return self.find(key1) == self.find(key2)

    def getRoots(self) -> List[T]:
        return list(set(self.find(key) for key in self._parent))

    def getGroups(self) -> DefaultDict[T, List[T]]:
        groups = defaultdict(list)
        for key in self._parent:
            root = self.find(key)
            groups[root].append(key)
        return groups

    def add(self, key: T) -> bool:
        if key in self._parent:
            return False
        self._parent[key] = key
        self._rank[key] = 1
        self.part += 1
        return True

    def __repr__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())

    def __len__(self) -> int:
        return self.part

    def __contains__(self, key: T) -> bool:
        return key in self._parent


if __name__ == "__main__":
    uf = UnionFindArrayWithUndo(10)
    uf.union(2, 4)
    assert uf.isConnected(2, 4)
    uf.undo()
    assert not uf.isConnected(2, 4)

    uf2 = UnionFindMapWithUndo()
    uf2.union(2, 4)
    assert uf2.isConnected(2, 4)
    uf2.undo()
    assert not uf2.isConnected(2, 4)
