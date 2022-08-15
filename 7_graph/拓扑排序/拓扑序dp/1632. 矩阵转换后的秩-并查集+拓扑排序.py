from collections import defaultdict, deque
from typing import List


# 给你一个 m x n 的矩阵 matrix ，请你返回一个新的矩阵 answer ，
# 其中 answer[row][col] 是 matrix[row][col] 的`秩`。

# 1 <= m, n <= 500
class Solution:
    def matrixRankTransform(self, matrix: List[List[int]]) -> List[List[int]]:
        """
        大小关系:拓扑排序的deps
        相等的点秩相等，可视为一个点：并查集
        """
        ROW, COL = len(matrix), len(matrix[0])
        uf = UnionFind(ROW * COL + 10)

        # 1.首先把所有行和列中元素相同的节点合并
        for r in range(ROW):
            rowRecord = defaultdict(list)
            for c in range(COL):
                rowRecord[matrix[r][c]].append(r * COL + c)
            for points in rowRecord.values():
                for p1, p2 in zip(points, points[1:]):
                    uf.union(p1, p2)

        for c in range(COL):
            colRecord = defaultdict(list)
            for r in range(ROW):
                colRecord[matrix[r][c]].append(r * COL + c)
            for points in colRecord.values():
                for p1, p2 in zip(points, points[1:]):
                    uf.union(p1, p2)

        # 2. 建图，因为要保证入度不重复计算，所以连边之前判断是否已经连接；所有值相等的点要缩成一个点
        adjMap = defaultdict(set)
        indegree = [0] * (ROW * COL)

        for r in range(ROW):
            row_ = sorted([(matrix[r][c], c) for c in range(COL)])
            for p1, p2 in zip(row_, row_[1:]):
                root1, root2 = uf.find(r * COL + p1[1]), uf.find(r * COL + p2[1])
                if root1 != root2 and root2 not in adjMap[root1]:
                    adjMap[root1].add(root2)
                    indegree[root2] += 1

        for c in range(COL):
            col_ = sorted([(matrix[r][c], r) for r in range(ROW)])
            for p1, p2 in zip(col_, col_[1:]):
                root1, root2 = uf.find(p1[1] * COL + c), uf.find(p2[1] * COL + c)
                if root1 != root2 and root2 not in adjMap[root1]:
                    adjMap[root1].add(root2)
                    indegree[root2] += 1

        # 3.拓扑排序，把入度等于0的节点加入队列，拓扑序dp求每个数的下界
        queue = deque([])
        for r in range(ROW):
            for c in range(COL):
                id = r * COL + c
                if indegree[id] == 0:
                    queue.append(id)

        mins = [1] * (ROW * COL)
        while queue:
            cur = queue.popleft()
            for next in adjMap[cur]:
                indegree[next] -= 1
                mins[next] = max(mins[next], mins[cur] + 1)
                if indegree[next] == 0:
                    queue.append(next)

        # 4.把结果放入矩阵
        res = [[0] * COL for _ in range(ROW)]
        for r in range(ROW):
            for c in range(COL):
                root = uf.find(r * COL + c)
                res[r][c] = mins[root]
        return res


class UnionFind:
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


print(Solution().matrixRankTransform(matrix=[[1, 2], [3, 4]]))
