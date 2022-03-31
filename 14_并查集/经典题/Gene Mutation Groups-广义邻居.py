from collections import defaultdict
from typing import DefaultDict, Generic, Iterable, List, Optional, TypeVar


T = TypeVar('T')


class UnionFindMap(Generic[T]):
    def __init__(self, iterable: Optional[Iterable[T]] = None):
        self.count = 0
        self.parent = dict()
        self.rank = defaultdict(lambda: 1)
        for item in iterable or []:
            self._add(item)

    def union(self, key1: T, key2: T) -> bool:
        root1 = self.find(key1)
        root2 = self.find(key2)
        if root1 == root2:
            return False
        if self.rank[root1] > self.rank[root2]:
            root1, root2 = root2, root1
        self.parent[root1] = root2
        self.rank[root2] += self.rank[root1]
        self.count -= 1
        return True

    def find(self, key: T) -> T:
        if key not in self.parent:
            self._add(key)
            return key

        while self.parent.get(key, key) != key:
            self.parent[key] = self.parent[self.parent[key]]
            key = self.parent[key]
        return key

    def isConnected(self, key1: T, key2: T) -> bool:
        return self.find(key1) == self.find(key2)

    def getRoots(self) -> List[T]:
        return list(set(self.find(key) for key in self.parent))

    def getGroup(self) -> DefaultDict[T, List[T]]:
        groups = defaultdict(list)
        for key in self.parent:
            root = self.find(key)
            groups[root].append(key)
        return groups

    def _add(self, key: T) -> bool:
        if key in self.parent:
            return False
        self.parent[key] = key
        self.rank[key] = 1
        self.count += 1
        return True


# If strings a and b are the same string except for one character, then a and b are in the same mutation group.
# If strings a and b are in a group and b and c are in a group, then a and c are in the same group.
class Solution:
    def solve(self, genes: List[str]):
        """求出group数量"""
        uf = UnionFindMap()
        for word in genes:
            for i in range(len(word)):
                next = word[:i] + '*' + word[i + 1 :]
                uf.union(word, next)
        return len(uf.getRoots())

