from typing import List
from collections import deque

# 题目要求:求出每个连通区域(保证是矩形)的左上角/右下角坐标
# 用bfs即可,起始点是左上角,不断更新右下角,最后一个点肯定是右下角
class Solution:
    def findFarmland(self, land: List[List[int]]) -> List[List[int]]:
        m, n = len(land), len(land[0])
        res = []
        for i in range(m):
            for j in range(n):
                if land[i][j]:  # found farmland
                    self.bfs(land, m, n, res, i, j)
        return res

    def bfs(self, land, m, n, res, i, j):
        mini, minj = i, j
        maxi, maxj = i, j
        queue = deque([(i, j,)])
        land[i][j] = 0  # mark as visited
        while queue:
            i, j = queue.popleft()
            for ii, jj in (i - 1, j), (i, j - 1), (i, j + 1), (i + 1, j):
                if 0 <= ii < m and 0 <= jj < n and land[ii][jj]:
                    queue.append((ii, jj))
                    land[ii][jj] = 0
                    maxi = max(maxi, ii)
                    maxj = max(maxj, jj)
        res.append([mini, minj, maxi, maxj])


print(Solution().findFarmland(land=[[1, 0, 0], [0, 1, 1], [0, 1, 1]]))
