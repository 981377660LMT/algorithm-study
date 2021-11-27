from typing import List
from collections import deque


# 每个单元格不是 0（空）就是 1（障碍物）
# 您 最多 可以消除 k 个障碍物，请找出从左上角 (0, 0) 到
# 右下角 (m-1, n-1) 的最短路径，并返回通过该路径所需的步数。如果找不到这样的路径，则返回 -1。

# 1 <= m, n <= 40


class Solution:
    def shortestPath(self, grid: List[List[int]], k: int) -> int:
        if not any(grid):
            return -1
        m, n = len(grid), len(grid[0])
        queue = deque()
        queue.append((0, 0, k, 0))
        # 1.剪枝:贴墙走需要的最大k
        k = min(k, m + n - 3)
        visited = set([(0, 0, k)])

        while queue:
            i, j, remain, steps = queue.popleft()
            if i == m - 1 and j == n - 1:
                return steps
            for ni, nj in [(i - 1, j), (i + 1, j), (i, j - 1), (i, j + 1)]:
                if ni >= 0 and ni < m and nj >= 0 and nj < n:
                    nk = remain - grid[ni][nj]
                    if nk < 0 or (ni, nj, nk) in visited:
                        continue
                    queue.append((ni, nj, nk, steps + 1))
                    visited.add((ni, nj, nk))
        return -1


print(Solution().shortestPath([[0, 0, 0], [1, 1, 0], [0, 0, 0], [0, 1, 1], [0, 0, 0]], 1))
