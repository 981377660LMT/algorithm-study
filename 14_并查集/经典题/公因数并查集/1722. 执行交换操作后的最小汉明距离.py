from typing import List
from collections import Counter, defaultdict


class UnionFind:
    def __init__(self, n):
        self.parent = list(range(n))
        self.rank = [1] * n

    def find(self, p):
        if p != self.parent[p]:
            self.parent[p] = self.find(self.parent[p])  # path compression
        return self.parent[p]

    def union(self, p, q):
        prt, qrt = self.find(p), self.find(q)
        if prt == qrt:
            return False  # already connected
        if self.rank[prt] > self.rank[qrt]:
            prt, qrt = qrt, prt
        self.parent[prt] = qrt
        self.rank[qrt] += self.rank[prt]
        return True


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
        groups = defaultdict(list)
        for i in range(n):
            root = uf.find(i)
            groups[root].append(i)
        print(groups)

        # 计算每个连通块对应的source元素与target的差集 ??
        # 即: 可交换区域作比较,看不同的个数
        common = 0
        for indexes in groups.values():
            c1 = Counter([source[i] for i in indexes])
            c2 = Counter([target[i] for i in indexes])
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

