# 所有数字都不同 求使得数组递增的最小交换次数
# 此题类似于情侣牵手
# 把哪些冲突的放在一组，解决这些冲突需要(size-1)次交换

# 因此最后需要的交换数为 n-count

from typing import List

from collections import defaultdict
from typing import DefaultDict, Generic, Iterable, List, Optional, TypeVar


# 元素是0-n-1的并查集写法，不支持动态添加
class UnionFindArray:
    def __init__(self, n: int):
        self.n = n
        self.count = n
        self.parent = list(range(n))
        self.rank = [1] * n

    def find(self, x: int) -> int:
        if x != self.parent[x]:
            self.parent[x] = self.find(self.parent[x])
        return self.parent[x]

    def union(self, x: int, y: int) -> bool:
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            return False
        if self.rank[rootX] > self.rank[rootY]:
            rootX, rootY = rootY, rootX
        self.parent[rootX] = rootY
        self.rank[rootY] += self.rank[rootX]
        self.count -= 1
        return True

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)

    def getGroups(self) -> DefaultDict[int, List[int]]:
        groups = defaultdict(list)
        for key in range(self.n):
            root = self.find(key)
            groups[root].append(key)
        return groups


class Solution:
    def solve(self, nums: List[int]) -> int:
        n = len(nums)
        uf = UnionFindArray(n)
        pairs = [(num, i) for i, num in enumerate(nums)]
        pairs.sort()

        for i in range(n):
            # 混淆的组相连
            uf.union(pairs[i][1], i)

        return n - uf.count


print(Solution().solve(nums=[3, 2, 1, 4]))
# We can swap 3 and 1.
