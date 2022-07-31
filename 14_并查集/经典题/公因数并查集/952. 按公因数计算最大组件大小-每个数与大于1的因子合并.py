# 只有当 A[i] 和 A[j] 共用一个大于 1 的公因数时，A[i] 和 A[j] 之间才有一条边。
# 返回图中最大连通组件的大小。
# 1 <= A.length <= 2e4
# !1 <= A[i] <= 1e5

# !并查集通过每个数的因子间接相连 枚举每个数的大于1的因子 因子与每个数相连即可


from collections import defaultdict
from functools import lru_cache
from math import floor
from typing import DefaultDict, Generic, Hashable, Iterable, List, Optional, TypeVar


T = TypeVar("T", bound=Hashable)


@lru_cache(None)
def getFactors(n: int) -> List[int]:
    """n 的所有大于1的因数"""
    if n <= 0:
        return []
    small, big = [], []
    upper = floor(n**0.5) + 1
    for i in range(2, upper):
        if n % i == 0:
            small.append(i)
            if i != n // i:
                big.append(n // i)
    return small + big[::-1]


class UnionFindMap(Generic[T]):
    """当元素不是数组index时(例如字符串)，更加通用的并查集写法，支持动态添加"""

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
        if self.rank[root1] > self.rank[root2]:
            root1, root2 = root2, root1
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


class Solution:
    def largestComponentSize(self, nums: List[int]) -> int:
        """每个数与大于1的因子合并(从上往下合并)"""
        uf = UnionFindMap()
        for num in nums:
            for factor in getFactors(num):
                uf.union(factor, num)

        group = defaultdict(int)
        for num in nums:
            root = uf.find(num)
            group[root] += 1
        return max(group.values(), default=0)


print(Solution().largestComponentSize([20, 50, 9, 63]))
