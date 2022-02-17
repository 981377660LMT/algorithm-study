# 返回消除所有0的最少操作次数
# 1284. 转化为全零矩阵的最少反转次数
from collections import deque
from typing import List


class Solution:
    def removeOnes(self, grid: List[List[int]]) -> int:
        row, col = len(grid), len(grid[0])
        state = sum(1 << (r * col + c) for r in range(row) for c in range(col) if grid[r][c] == 1)
        queue = deque([(state, 0)])
        visited = set([state])

        while queue:
            cur, step = queue.popleft()
            if cur == 0:
                return step

            for r in range(row):
                for c in range(col):
                    if grid[r][c] == 0:
                        continue
                    next = cur
                    for nextR in range(row):
                        next &= ~(1 << (nextR * col + c))
                    for nextC in range(col):
                        next &= ~(1 << (r * col + nextC))
                    if next not in visited:
                        visited.add(next)
                        queue.append((next, step + 1))

        return -1


print(Solution().removeOnes(grid=[[0, 1, 0], [1, 0, 1], [0, 1, 0]]))
