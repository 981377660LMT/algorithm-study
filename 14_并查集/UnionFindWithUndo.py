"""
配合回溯的场景使用
撤销相当于弹出栈顶元素

很少用到撤销操作，因为并查集的撤销可以变成倒着合并

应用场景:
可持久化并查集的离线处理
!在树上(版本之间)dfs 递归时要union结点 回溯时候需要撤销的场合
"""


from collections import defaultdict
from typing import DefaultDict, Generic, Hashable, Iterable, List, Optional, TypeVar


class RevocableUnionFindArray:
    """
    带撤销操作的并查集

    不能使用路径压缩优化（因为路径压缩会改变结构）；
    为了不超时必须使用按秩合并优化,复杂度nlogn
    """

    __slots__ = ("n", "part", "parent", "rank", "optStack")

    def __init__(self, n: int):
        self.n = n
        self.part = n
        self.parent = list(range(n))
        self.rank = [1] * n
        self.optStack = []

    def find(self, x: int) -> int:
        """不能使用路径压缩优化"""
        while self.parent[x] != x:
            x = self.parent[x]
        return x

    def union(self, x: int, y: int) -> bool:
        """x所在组合并到y所在组"""
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            self.optStack.append((-1, -1, -1))
            return False

        if self.rank[rootX] > self.rank[rootY]:
            rootX, rootY = rootY, rootX

        self.parent[rootX] = rootY
        self.rank[rootY] += self.rank[rootX]
        self.part -= 1
        self.optStack.append((rootX, rootY, self.rank[rootX]))
        return True

    def revocate(self) -> None:
        """
        用一个栈记录前面的合并操作，
        撤销时要依次取出栈顶元素做合并操作的逆操作.
        !没合并成功也要撤销.
        """
        if not self.optStack:
            raise IndexError("no union option to revocate")

        rootX, rootY, rankX = self.optStack.pop()
        if rootX == -1:
            return

        self.parent[rootX] = rootX
        self.rank[rootY] -= rankX
        self.part += 1

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)

    def getGroups(self) -> DefaultDict[int, List[int]]:
        groups = defaultdict(list)
        for key in range(self.n):
            root = self.find(key)
            groups[root].append(key)
        return groups


T = TypeVar("T", bound=Hashable)


class RevocableUnionFindMap(Generic[T]):
    """
    带撤销操作的并查集

    不能使用路径压缩优化（因为路径压缩会改变结构）；
    为了不超时必须使用按秩合并优化,复杂度nlogn
    """

    __slots__ = ("part", "parent", "rank", "optStack")

    def __init__(self, iterable: Optional[Iterable[T]] = None):
        self.part = 0
        self.parent = dict()
        self.rank = dict()
        self.optStack = []
        for item in iterable or []:
            self.add(item)

    def find(self, key: T) -> T:
        """不能使用路径压缩优化"""
        if key not in self.parent:
            self.add(key)
            return key

        while self.parent.get(key, key) != key:
            key = self.parent[key]
        return key

    def union(self, key1: T, key2: T) -> bool:
        """rank一样时 默认key2作为key1的父节点"""
        root1 = self.find(key1)
        root2 = self.find(key2)
        if root1 == root2:
            self.optStack.append((-1, -1, -1))
            return False
        if self.rank[root1] > self.rank[root2]:
            root1, root2 = root2, root1
        self.parent[root1] = root2
        self.rank[root2] += self.rank[root1]
        self.part -= 1
        self.optStack.append((root1, root2, self.rank[root1]))
        return True

    def revocate(self) -> None:
        """
        用一个栈记录前面的合并操作，
        撤销时要依次取出栈顶元素做合并操作的逆操作.
        !没合并成功也要撤销.
        """
        if not self.optStack:
            raise IndexError("no union option to revocate")

        root1, root2, rank1 = self.optStack.pop()
        if root1 == -1:
            return

        self.parent[root1] = root1
        self.rank[root2] -= rank1
        self.part += 1

    def isConnected(self, key1: T, key2: T) -> bool:
        return self.find(key1) == self.find(key2)

    def getRoots(self) -> List[T]:
        return list(set(self.find(key) for key in self.parent))

    def getGroups(self) -> DefaultDict[T, List[T]]:
        groups = defaultdict(list)
        for key in self.parent:
            root = self.find(key)
            groups[root].append(key)
        return groups

    def add(self, key: T) -> bool:
        if key in self.parent:
            return False
        self.parent[key] = key
        self.rank[key] = 1
        self.part += 1
        return True

    def __repr__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())

    def __len__(self) -> int:
        return self.part

    def __contains__(self, key: T) -> bool:
        return key in self.parent


if __name__ == "__main__":
    uf = RevocableUnionFindArray(10)
    uf.union(2, 4)
    assert uf.isConnected(2, 4)
    uf.revocate()
    assert not uf.isConnected(2, 4)

    uf2 = RevocableUnionFindMap()
    uf2.union(2, 4)
    assert uf2.isConnected(2, 4)
    uf2.revocate()
    assert not uf2.isConnected(2, 4)
