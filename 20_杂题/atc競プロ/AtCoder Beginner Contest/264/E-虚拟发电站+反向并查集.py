"""
1. 需要一堆结点相互连通时，常用的技巧为用一个虚拟结点连通这些点 O(n^2) -> O(n)
2. 反向并查集
"""

# 虚拟发电站+反向并查集
# 每次剪断一根电线，查询多少个城市与发电站相连
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

from collections import defaultdict
from typing import DefaultDict, Generic, Hashable, Iterable, List, Optional, TypeVar


T = TypeVar("T", bound=Hashable)


class UnionFindMap(Generic[T]):
    """当元素不是数组index时(例如字符串)，更加通用的并查集写法，支持动态添加"""

    __slots__ = ("part", "parent", "rank")

    def __init__(self, iterable: Optional[Iterable[T]] = None):
        self.part = 0
        self.parent = dict()
        self.rank = defaultdict(lambda: 1)
        for item in iterable or []:
            self._add(item)

    def union(self, key1: T, key2: T) -> bool:
        """rank一样时 默认key2作为key1的父节点"""
        root1 = self.find(key1)
        root2 = self.find(key2)
        if root1 == root2:
            return False
        self.parent[root1] = root2
        self.rank[root2] += self.rank[root1]
        self.part -= 1
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

    def getGroups(self) -> DefaultDict[T, List[T]]:
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
        self.part += 1
        return True

    def __repr__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())

    def __len__(self) -> int:
        return self.part

    def __contains__(self, key: T) -> bool:
        return key in self.parent


# 城市为0-n-1
# 发电厂为n - n+m-1
# 一度切れた電線は、その後のイベントにおいても切れたままです。
n, m, e = map(int, input().split())
edges = []

uf = UnionFindMap()
for _ in range(e):
    a, b = map(int, input().split())
    edges.append((a - 1, b - 1))  # 所有电线

DUMMY = -1  # !大发电站
for i in range(n, n + m):
    uf.union(i, DUMMY)

q = int(input())
adds = [int(input()) - 1 for _ in range(q)]
adds.reverse()  # 添加的电线
bad = set(adds)

for i in range(len(edges)):
    if i in bad:
        continue
    a, b = edges[i]
    uf.union(a, b)


res = []
for i in adds:
    res.append(uf.rank[uf.find(DUMMY)] - (m + 1))  # !注意需要找到大发电站的组而不是大发电站本身
    u, v = edges[i]
    uf.union(u, v)

res.reverse()
print(*res, sep="\n")
