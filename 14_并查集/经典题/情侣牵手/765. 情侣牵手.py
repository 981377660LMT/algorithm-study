from typing import List

# 所有数字都不同 求使得数组递增的最小交换次数
# 此题类似于情侣牵手
# 把哪些冲突的放在一组，解决这些冲突需要(size-1)次交换

# 因此最后需要的交换数为 n-count


from collections import defaultdict
from typing import DefaultDict, List


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
    def minSwapsCouples(self, nums: List[int]) -> int:
        """
        如果我们有 k 对情侣形成了错误环，需要交换 k - 1 次才能让情侣牵手。
        问题转化成 n / 2 对情侣中，有多少个这样的环
        """

        n = len(nums)
        uf = UnionFindArray(n // 2)
        for i in range(0, n, 2):
            # 这两组连在一起
            uf.union(nums[i] // 2, nums[i + 1] // 2)  # 除以2表示对应的分组
        return n // 2 - uf.count

