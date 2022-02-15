from typing import List

# 2 <= n <= 100


class UnionFind:
    # 初始化 parent & weights
    def __init__(self, n):
        self.parent = list(range(n))
        self.weights = None

    def find(self, x):
        if x != self.parent[x]:
            self.parent[x] = self.find(self.parent[x])
        return self.parent[x]

    def union(self, x, y):
        self.parent[x] = y

    # 最小生成树-kruskal mst只需要传入排好序的edges即可
    def getMSTWeight(self, edges) -> int:
        if self.weights is not None:
            return self.weights
        else:
            self.weights = 0
            for _, x, y, w in edges:
                rx, ry = self.find(x), self.find(y)
                if rx != ry:
                    self.union(rx, ry)
                    self.weights += w
            return self.weights


# 请你找到给定图中最小生成树的所有关键边和伪关键边。
# 如果从图中删去某条边，会导致最小生成树的权值和增加，那么我们就说它是一条关键边。
# 伪关键边则是可能会出现在某些最小生成树中但不会出现在所有最小生成树中的边。


# 总结:
# 1. 遍历每条边，看是不是最小生成树的边(连接当前边再进行后续最小生成树合成)
# 2. 看是不是关键边(去掉无关紧要,伪边)


class Solution:
    def findCriticalAndPseudoCriticalEdges(self, n: int, edges: List[List[int]]) -> List[List[int]]:

        # 对边根据权重从小到大排序
        edges = [[i] + edge for i, edge in enumerate(edges)]
        edges.sort(key=lambda x: x[-1])

        # 计算最小生成树的权值weights
        uf = UnionFind(n)
        weights = uf.getMSTWeight(edges)

        # 判断每条边 in 关键边/伪关键边/与构建最小生成树无关的边
        keyEdges = []
        pseudoEdges = []
        for i, edge in enumerate(edges):

            # 创建删除该条边的集合
            _, u, v, w = edge
            tmpEdges = edges[:i] + edges[i + 1 :]

            # 连接当前边再进行后续最小生成树合成
            uf = UnionFind(n)
            # 在树中先取这条边
            uf.union(u, v)
            weightWithW = uf.getMSTWeight(tmpEdges) + w

            # 大于说明该边是与构建最小生成树无关的边
            if weightWithW > weights:
                continue

            # 去掉当前边，得到此时最小生成树权值
            uf = UnionFind(n)
            weightWithoutW = uf.getMSTWeight(tmpEdges)

            # 去掉无关紧要，伪边
            if weightWithW == weightWithoutW:
                pseudoEdges.append(edge[0])
            else:
                keyEdges.append(edge[0])

        return [keyEdges, pseudoEdges]


print(
    Solution().findCriticalAndPseudoCriticalEdges(
        n=5, edges=[[0, 1, 1], [1, 2, 1], [2, 3, 2], [0, 3, 2], [0, 4, 3], [3, 4, 3], [1, 4, 6]]
    )
)
# 输出：[[0,1],[2,3,4,5]]
