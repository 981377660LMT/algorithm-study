from typing import Dict, List
from collections import Counter, defaultdict


class UnionFind:
    def __init__(self, n: int):
        self.size = n
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

    def getGroups(self) -> Dict[int, List[int]]:
        groups = defaultdict(list)
        for key in range(self.size):
            root = self.find(key)
            groups[root].append(key)
        return groups


# 1. 并查集获取帮派邻接表
# 2. 计算每个连通块对应的source元素与target的差集
class Solution:
    def minimumHammingDistance(
        self, source: List[int], target: List[int], allowedSwaps: List[List[int]]
    ) -> int:
        n = len(source)
        uf = UnionFind(n)
        for i, j in allowedSwaps:
            uf.union(i, j)

        # 获取根节点对应的连通块(可交换区域)
        groups = uf.getGroups()

        # 计算每个连通块对应的source元素与target的差集 ??
        # 即: 可交换区域作比较,看不同的个数
        common = 0
        for group in groups.values():
            c1 = Counter([source[i] for i in group])
            c2 = Counter([target[i] for i in group])
            # 区域里元素相同的对数
            common += len(list((c1 & c2).elements()))
        return n - common


print(
    Solution().minimumHammingDistance(
        source=[1, 2, 3, 4], target=[2, 1, 4, 5], allowedSwaps=[[0, 1], [2, 3]]
    )
)


# 输出：1
# 解释：source 可以按下述方式转换：
# - 交换下标 0 和 1 指向的元素：source = [2,1,3,4]
# - 交换下标 2 和 3 指向的元素：source = [2,1,4,3]
# source 和 target 间的汉明距离是 1 ，二者有 1 处元素不同，在下标 3 。

