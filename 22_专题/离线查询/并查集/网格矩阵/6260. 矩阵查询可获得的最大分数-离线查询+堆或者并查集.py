# 6260. 矩阵查询可获得的最大分数

# 给你一个大小为 m x n 的整数矩阵 grid 和一个大小为 k 的数组 queries 。
# 找出一个大小为 k 的数组 answer ，且满足对于每个整数 queres[i] ，
# 你从矩阵 左上角 单元格开始，重复以下过程：
# 如果 queries[i] 严格 大于你当前所处位置单元格，如果该单元格是第一次访问，
# 则获得 1 分，并且你可以移动到所有 4 个方向（上、下、左、右）上任一 相邻 单元格。
# 否则，你不能获得任何分，并且结束这一过程。
# 在过程结束后，answer[i] 是你可以获得的最大分数。注意，对于每个查询，你可以访问同一个单元格 多次 。
# 返回结果数组 answer 。
# !ROW*COL<=10^5
# !Q<=10^5

# !注意不能bfs维护每次需要访问的点,这样是O(ROW*COL*(ROW+COL)) 会超时
# 1. 离线查询+堆 => 每次更新都看当前值最小的点,看是否能继续扩散
# 2. 离线查询+并查集 => 类似于水位上升/洪水泛滥,最后查询(0,0)所在的连通分量大小
# 因为这个题目描述是遍历 所以容易想到队列搜索 如果是一个水位上涨的描述 容易想并查集
# 事实上 floodfill 和 遍历 是一样的

from collections import defaultdict
from heapq import heappop, heappush
from typing import DefaultDict, List


DIR4 = [(0, 1), (0, -1), (1, 0), (-1, 0)]


class Solution:
    def maxPoints(self, grid: List[List[int]], queries: List[int]) -> List[int]:
        """离线查询+堆 => 每次更新都看当前值最小的点,看是否能继续扩散(接雨水II)"""
        ROW, COL = len(grid), len(grid[0])
        Q = sorted([(qv, qi) for qi, qv in enumerate(queries)], key=lambda x: x[0])
        res = [0] * len(Q)

        pq = [(grid[0][0], 0, 0)]
        visited = [[False] * COL for _ in range(ROW)]
        visited[0][0] = True
        count = 0
        for qv, qi in Q:
            while pq and pq[0][0] < qv:  # !当前最小值小于当前查询值,可以继续扩散
                _, curRow, curCol = heappop(pq)
                count += 1
                for dr, dc in DIR4:
                    newRow, newCol = curRow + dr, curCol + dc
                    if 0 <= newRow < ROW and 0 <= newCol < COL and not visited[newRow][newCol]:
                        visited[newRow][newCol] = True
                        heappush(pq, (grid[newRow][newCol], newRow, newCol))
            res[qi] = count
        return res

    def maxPoints2(self, grid: List[List[int]], queries: List[int]) -> List[int]:
        """离线查询+并查集"""
        ROW, COL = len(grid), len(grid[0])
        nums = []
        for r in range(ROW):
            for c in range(COL):
                nums.append((grid[r][c], r, c))
        nums.sort(key=lambda x: x[0])
        Q = sorted([(qv, qi) for qi, qv in enumerate(queries)], key=lambda x: x[0])

        ni = 0
        res = [0] * len(Q)
        uf = UnionFindArray(ROW * COL)
        for qv, qi in Q:
            while ni < len(nums) and nums[ni][0] < qv:  # 水位上涨,合并
                _, curRow, curCol = nums[ni]
                ni += 1
                for dr, dc in DIR4:
                    nextRow, nextCol = curRow + dr, curCol + dc
                    if 0 <= nextRow < ROW and 0 <= nextCol < COL and grid[nextRow][nextCol] < qv:
                        uf.union(curRow * COL + curCol, nextRow * COL + nextCol)
            if grid[0][0] < qv:
                res[qi] = uf.rank[uf.find(0)]
        return res


class UnionFindArray:

    __slots__ = ("n", "part", "parent", "rank")

    def __init__(self, n: int):
        self.n = n
        self.part = n
        self.parent = list(range(n))
        self.rank = [1] * n

    def find(self, x: int) -> int:
        while x != self.parent[x]:
            self.parent[x] = self.parent[self.parent[x]]
            x = self.parent[x]
        return x

    def union(self, x: int, y: int) -> bool:
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            return False
        if self.rank[rootX] > self.rank[rootY]:
            rootX, rootY = rootY, rootX
        self.parent[rootX] = rootY
        self.rank[rootY] += self.rank[rootX]
        self.part -= 1
        return True

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)

    def getGroups(self) -> DefaultDict[int, List[int]]:
        groups = defaultdict(list)
        for key in range(self.n):
            root = self.find(key)
            groups[root].append(key)
        return groups

    def getRoots(self) -> List[int]:
        return list(set(self.find(key) for key in self.parent))

    def __repr__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())

    def __len__(self) -> int:
        return self.part
