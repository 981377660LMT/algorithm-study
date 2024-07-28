# 0表示空地 1表示墙
# !现在要阻断左上角到右下角的路 问最少需要加多少墙
# https://binarysearch.com/problems/Walled-Off

# 答案是 0 1 2
from collections import defaultdict
from typing import DefaultDict, List


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


DIR8 = ((0, 1), (1, 0), (0, -1), (-1, 0), (1, 1), (1, -1), (-1, 1), (-1, -1))


class Solution:
    def solve(self, matrix: List[List[int]]) -> int:
        # !如果障碍物阻断了上右边到左下边的路径 那么此时两点就不连通
        ROW, COL = len(matrix), len(matrix[0])
        uf = UnionFindArray(ROW * COL + 10)

        P1, P2 = ROW * COL + 5, ROW * COL + 6  # !右上边 左下边 如果这两个虚拟点相连 那么就被阻断了
        for r in range(ROW):
            for c in range(COL):
                if matrix[r][c] == 1:  # 障碍物
                    cur = r * COL + c
                    for dr, dc in DIR8:
                        nr, nc = r + dr, c + dc
                        if 0 <= nr < ROW and 0 <= nc < COL:
                            next = nr * COL + nc
                            if matrix[nr][nc] == 1:
                                uf.union(cur, next)
                        else:
                            if nr < 0 or nc >= COL:
                                uf.union(cur, P1)
                            elif nr >= ROW or nc < 0:
                                uf.union(cur, P2)

        if uf.isConnected(P1, P2):
            return 0

        root1, root2 = uf.find(P1), uf.find(P2)
        # 中间如果能多一个点相连，那么就只需要1个
        for r in range(ROW):
            for c in range(COL):
                if matrix[r][c] == 0:
                    if (r, c) in [(0, 0), (ROW - 1, COL - 1)]:
                        continue

                    # !看连接以后是否能把两个组连起来(中间一个人牵着左右手)
                    # !类似827. 最大人工岛的虚拟添加的技巧 找连接后影响的组 而不是真正去连接
                    roots = set()
                    for dr, dc in DIR8:
                        nr, nc = r + dr, c + dc
                        if 0 <= nr < ROW and 0 <= nc < COL:
                            if matrix[nr][nc] == 1:
                                roots.add(uf.find(nr * COL + nc))
                        else:
                            if nr < 0 or nc >= COL:
                                roots.add(root1)
                            elif nr >= ROW or nc < 0:
                                roots.add(root2)
                    if root1 in roots and root2 in roots:
                        return 1

        return 2


print(Solution().solve(matrix=[[0, 1, 1], [0, 0, 0], [0, 0, 0], [1, 1, 0]]))
