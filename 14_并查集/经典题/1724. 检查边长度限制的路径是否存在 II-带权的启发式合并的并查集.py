from typing import List

# https://leetcode-cn.com/problems/checking-existence-of-edge-length-limited-paths-ii/comments/1177392
# 带权的启发式合并的并查集


# 总结：
# 想象两个子节点不断向上走 联通即父节点相同 但是在上找的过程中会收到limit限制停下


class DistanceLimitedPathsExist:
    # 以给定的无向图初始化对象。
    # edgeList[i] = [ui, vi, disi] 表示一条连接 ui 和 vi ，距离为 disi 的边
    # 注意，同一对节点间可能有多条边，且该图可能不是连通的。
    def __init__(self, n: int, edgeList: List[List[int]]):
        self.parent = list(range(n))
        self.size = [1] * n
        self.cost = [0] * n
        for u, v, w in sorted(edgeList, key=lambda x: x[-1]):
            self.__union(u, v, w)

    # 是否存在一条从 p 到 q 的路径，且路径中每条边的距离都`严格小于` limit
    # 总结：
    # 分别找到 p、q 经过小于 limit 的边能达到的最高的祖宗节点，判断是否相同即可。
    # 不能用路径压缩了，会破坏树的结构。
    def query(self, p: int, q: int, limit: int) -> bool:
        return self.__find(p, limit) == self.__find(q, limit)

    def __union(self, x, y, w):
        # 注意这里w+1，严格小于
        rx, ry = self.__find(x, w + 1), self.__find(y, w + 1)
        if rx == ry:
            return

        if self.size[rx] > self.size[ry]:
            rx, ry = ry, rx
        self.parent[rx] = ry
        self.size[ry] += self.size[rx]
        self.cost[rx] = w

    def __find(self, x, limit):
        while self.parent[x] != x and self.cost[x] < limit:
            x = self.parent[x]
        return x


d = DistanceLimitedPathsExist(6, [[0, 2, 4], [0, 3, 2], [1, 2, 3], [2, 3, 1], [4, 5, 5]])
print(d.__dict__)
print(d.query(2, 3, 2))
print(d.query(1, 3, 3))
print(d.query(2, 0, 3))
print(d.query(0, 5, 6))

#  [true, false, true, false]
