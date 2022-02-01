from collections import Counter, defaultdict
from typing import Generic, Iterable, List, TypeVar

T = TypeVar('T')


class UnionFindMap(Generic[T]):
    def __init__(self, iterable: Iterable[T] = None):
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

    def _add(self, key: T) -> bool:
        if key in self.parent:
            return False
        self.parent[key] = key
        self.rank[key] = 1
        self.count += 1
        return True


class Solution:
    def groupStrings(self, words: List[str]) -> List[int]:
        states: List[int] = []
        for i, word in enumerate(words):
            state = 0
            for char in word:
                state |= 1 << (ord(char) - 97)
            states.append(state)

        uf = UnionFindMap(states)
        statesSet = set(states)

        for state in states:
            for i in range(26):
                addOrRemove = state ^ (1 << i)
                if addOrRemove in statesSet:
                    uf.union(addOrRemove, state)

                if (state >> i) & 1:
                    replace = state ^ (1 << i) | (1 << 27)
                    uf.union(replace, state)

        groupCounter = Counter()
        for state in states:
            root = uf.find(state)
            groupCounter[root] += 1

        return [uf.count, max(groupCounter.values())]


# 2 1 2 1
print(Solution().groupStrings(words=["a", "b", "ab", "cde"]))
print(Solution().groupStrings(words=["a", "ab", "abc"]))
