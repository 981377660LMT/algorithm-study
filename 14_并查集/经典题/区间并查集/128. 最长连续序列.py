from collections import defaultdict
from typing import DefaultDict, Generic, Hashable, Iterable, List, Optional, TypeVar


T = TypeVar("T", bound=Hashable)


class UnionFindMap(Generic[T]):
    """当元素不是数组index时(例如字符串),更加通用的并查集写法,支持动态添加"""

    __slots__ = ("part", "parent", "rank")

    def __init__(self, iterable: Optional[Iterable[T]] = None):
        self.part = 0
        self.parent = dict()
        self.rank = dict()
        for item in iterable or []:
            self.add(item)

    def union(self, key1: T, key2: T) -> bool:
        """rank一样时 默认key2作为key1的父节点"""
        root1 = self.find(key1)
        root2 = self.find(key2)
        if root1 == root2:
            return False
        if self.rank[root1] > self.rank[root2]:
            root1, root2 = root2, root1
        self.parent[root1] = root2
        self.rank[root2] += self.rank[root1]
        self.part -= 1
        return True

    def find(self, key: T) -> T:
        if key not in self.parent:
            self.add(key)
            return key

        while self.parent.get(key, key) != key:
            self.parent[key] = self.parent[self.parent[key]]
            key = self.parent[key]
        return key

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


class Solution:
    def longestConsecutive(self, nums: List[int]) -> int:
        if not nums:
            return 0

        S = set(nums)
        uf = UnionFindMap()

        # !合并左右
        for num in nums:
            if num - 1 in S:
                uf.union(num, num - 1)
            if num + 1 in S:
                uf.union(num, num + 1)
        return max((len(g) for g in uf.getGroups().values()), default=1)
