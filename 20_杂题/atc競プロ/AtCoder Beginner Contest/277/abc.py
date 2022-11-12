import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 10
# 9
#   階建てのビルがあり、N 本のはしごがかかっています。
# ビルの 1 階にいる高橋君ははしごを繰り返し使って（0 回でもよい）できるだけ高い階へ上りたいと考えています。
# はしごには 1 から N までの番号がついており、はしご i は A
# i
# ​
#   階と B
# i
# ​
#   階を結んでいます。はしご i を利用すると A
# i
# ​
#   階から B
# i
# ​
#   階へ、または B
# i
# ​
#   階から A
# i
# ​
#   階へ双方向に移動することができますが、それ以外の階の間の移動は行うことはできません。
# また、高橋君は同じ階での移動は自由に行うことができますが、はしご以外の方法で他の階へ移動することはできません。
# 高橋君は最高で何階へ上ることができますか？
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

    def getGroups(self):
        groups = defaultdict(set)
        for key in self.parent:
            root = self.find(key)
            groups[root].add(key)
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
    n = int(input())
    uf = UnionFindMap()
    for i in range(n):
        a, b = map(int, input().split())
        uf.union(a, b)
    groups = uf.getGroups()
    for group in groups.values():
        if 1 in group:
            print(max(group))
            exit(0)
    print(1)
