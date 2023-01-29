# Short Circuit

# 如果前一个单词尾与后一个单词首字母相同，那么就可以相连 (建边使用每个字母作为结点)
# 能否使每个单词都用上，连成环？
# words ≤ 15000

# 这个图最多26个点
# 判断一个有向图的是否存在欧拉回路
# 1. 所有点连通（并查集）
# 2. 所有点出度等于入度

from collections import defaultdict
from typing import DefaultDict, Generic, Iterable, List, Optional, TypeVar


class Solution:
    def solve(self, words):
        ind, outd = defaultdict(int), defaultdict(int)
        uf = UnionFindMap[int]()  # 不能自动union的并查集，需要add手动加入元素
        visited = set()

        for word in words:
            first, last = ord(word[0]), ord(word[-1])
            
            uf.add(first)
            uf.add(last)
            uf.union(first, last)

            outd[first] += 1
            ind[last] += 1

            visited.add(first)
            visited.add(last)

        return uf.count == 1 and all(ind[key] == outd[key] for key in visited)


T = TypeVar('T')


class UnionFindMap(Generic[T]):
    def __init__(self, iterable: Optional[Iterable[T]] = None):
        self.count = 0
        self.parent = dict()
        self.rank = defaultdict(lambda: 1)
        for item in iterable or []:
            self.add(item)

    def union(self, key1: T, key2: T) -> bool:

        if key1 not in self.parent or key2 not in self.parent:
            return False

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

    def add(self, key: T) -> bool:
        if key in self.parent:
            return False
        self.parent[key] = key
        self.rank[key] = 1
        self.count += 1
        return True


print(
    Solution().solve(
        words=[
            "gablatores",
            "Alvah",
            "narcotist",
            "hoastman",
            "halomorphic",
            "unflead",
            "bedye",
            "reclusive",
            "asself",
            "overquarter",
        ]
    )
)  # We can form the following circle: chair --> racket --> touch --> height --> tunic --> chair.
