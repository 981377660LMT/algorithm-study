# 只有当 A[i] 和 A[j] 共用一个大于 1 的公因数时，A[i] 和 A[j] 之间才有一条边。
# 返回图中最大连通组件的大小。
# 1 <= A.length <= 2e4
# !1 <= A[i] <= 1e5


from collections import defaultdict
from typing import DefaultDict, List


class Solution:
    def largestComponentSize(self, nums: List[int]) -> int:
        """每个因子与倍数合并(从下往上合并) 注意合并前判断倍数在不在nums里"""
        ok = set(nums)
        max_ = max(nums)
        uf = UnionFindArray(max_ + 1)
        for i in range(2, max_ + 1):
            for j in range(i * 2, max_ + 1, i):
                if j in ok:
                    uf.union(i, j)

        group = defaultdict(int)
        for num in nums:
            root = uf.find(num)
            group[root] += 1
        return max(group.values(), default=0)


class UnionFindArray:
    """元素是0-n-1的并查集写法,不支持动态添加

    初始化的连通分量个数 为 n
    """

    def __init__(self, n: int):
        self.n = n
        self.part = n
        self.parent = list(range(n))
        self.rank = [1] * n

    def find(self, x: int) -> int:
        if x != self.parent[x]:
            self.parent[x] = self.find(self.parent[x])
        return self.parent[x]

    def union(self, x: int, y: int) -> bool:
        """rank一样时 默认key2作为key1的父节点"""
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            return False
        if self.rank[rootX] > self.rank[rootY]:
            rootX, rootY = rootY, rootX
        self.parent[rootX] = rootY
        self.rank[rootY] += self.rank[rootX]
        self.part -= 1
        return True

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)

    def getGroups(self) -> DefaultDict[int, List[int]]:
        groups = defaultdict(list)
        for key in range(self.n):
            root = self.find(key)
            groups[root].append(key)
        return groups

    def getRoots(self) -> List[int]:
        return list(set(self.find(key) for key in self.parent))

    def __repr__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())

    def __len__(self) -> int:
        return self.part


print(Solution().largestComponentSize([20, 50, 9, 63]))
