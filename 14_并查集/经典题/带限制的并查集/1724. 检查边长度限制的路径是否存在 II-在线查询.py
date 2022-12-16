from typing import List

# https://leetcode.cn/problems/checking-existence-of-edge-length-limited-paths-ii/solution/ke-chi-jiu-hua-bing-cha-ji-by-megurine-l1m6/# 带权的启发式合并的并查集
# 带限制的并查集
# 1.所有边升序排列，从小到大依次合并并查集，每次合并时，记录连接的边的权值（时间戳）
# 2.查找公共祖先时，时间戳必须严格小于 limit 时，才能向祖先移动
# 显然，由于连接边带有时间戳信息，不能进行路径压缩，只能使用按秩合并来进行优化

# 带限制的并查集


class UnionFindWithLimit:
    def __init__(self, n: int):
        self.parent = list(range(n))
        self.rank = [1] * n
        self.time = [0] * n

    def union(self, x: int, y: int, time: int) -> None:
        rootX, rootY = self.find(x), self.find(y)
        if rootX == rootY:
            return
        if self.rank[rootX] < self.rank[rootY]:
            rootX, rootY = rootY, rootX
        self.parent[rootY] = rootX
        self.rank[rootX] += self.rank[rootY]
        self.time[rootY] = time  # 向root的距离

    def find(self, x: int, limit=int(1e20)) -> int:
        while self.parent[x] != x and self.time[x] < limit:
            x = self.parent[x]
        return x


class DistanceLimitedPathsExist:
    def __init__(self, n: int, edgeList: List[List[int]]):
        """
        以给定的无向图初始化对象。
        edgeList[i] = [ui, vi, disti] 表示一条连接 ui 和 vi ，距离为 disti 的边
        注意，同一对节点间可能有多条边，且该图可能不是连通的。
        """
        self.uf = UnionFindWithLimit(n)
        for u, v, w in sorted(edgeList, key=lambda x: x[-1]):
            self.uf.union(u, v, w)

    def query(self, p: int, q: int, limit: int) -> bool:
        """
        是否存在一条从 p 到 q 的路径，且路径中每条边的距离都`严格小于` limit

        分别找到 p、q 经过小于 limit 的边能达到的最高的祖宗节点，判断是否相同即可。
        不能用路径压缩了，会破坏树的结构。
        """
        return self.uf.find(p, limit) == self.uf.find(q, limit)


d = DistanceLimitedPathsExist(6, [[0, 2, 4], [0, 3, 2], [1, 2, 3], [2, 3, 1], [4, 5, 5]])
print(d.__dict__)
print(d.query(2, 3, 2))
print(d.query(1, 3, 3))
print(d.query(2, 0, 3))
print(d.query(0, 5, 6))

#  [true, false, true, false]
