# 即断开一个点,让联通区域增加最多


from typing import List
from collections import Counter


class UnionFind:
    def __init__(self, n):
        self.n = n
        self.part = n
        self.parent = list(range(n))
        self.size = [1] * n

    def find(self, x):
        if x != self.parent[x]:
            self.parent[x] = self.find(self.parent[x])
        return self.parent[x]

    def union(self, x, y):
        root_x = self.find(x)
        root_y = self.find(y)
        if root_x == root_y:
            return False
        if self.size[root_x] > self.size[root_y]:
            root_x, root_y = root_y, root_x
        self.parent[root_x] = root_y
        self.size[root_y] += self.size[root_x]
        self.part -= 1
        return True

    def isconnected(self, p, q):
        return self.find(p) == self.find(q)


# 1.构建并查集，做连通的时候，保存长度(长度也是有多少个节点的根是同一个，包括自身)
# 2.对 initaial 每个元素做find，找到根
# 3.对根计数，多个病毒指向同一个根的，删了也没用，另一个还是会感染所有，也就是删除其中任何一个，最后感染的还是这条连通分量的所有
# 4.选择感染最大的进行删除


class Solution:
    def minMalwareSpread(self, graph: List[List[int]], initial: List[int]) -> int:
        n = len(graph)
        uf = UnionFind(n)

        # 遍历
        for i in range(len(graph)):
            for j in range(i, len(graph)):
                if graph[i][j] == 1:
                    uf.union(i, j)

        # 看是哪些流派
        counter = Counter(uf.find(u) for u in initial)
        res = (-1, min(initial))
        for node in initial:
            root = uf.find(node)
            if counter[root] == 1:  # 初始感染者里这个流派只有一个人感染了,哪个流派人多就拯救哪个流派,否则会感染更多人
                if uf.size[root] > res[0]:
                    res = uf.size[root], node
                elif uf.size[root] == res[0] and node < res[1]:
                    res = uf.size[root], node

        return res[1]


print(Solution().minMalwareSpread(graph=[[1, 1, 0], [1, 1, 0], [0, 0, 1]], initial=[0, 1]))

