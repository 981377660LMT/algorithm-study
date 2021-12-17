from typing import List
from collections import deque

# 每一步，你可以选择一个单元格并将它反转
# 如果存在和它相邻的单元格，那么这些相邻的单元格也会被反转。

# 请你返回将矩阵 mat 转化为全零矩阵的最少反转次数，如果无法转化为全零矩阵，请返回 -1 。

# 1 <= m <= 3
# 1 <= n <= 3
# 果断bfs

# 1. 状压 起始和结束状态
# 2. 转移

# 时间复杂度2^n*m*n
class Solution:
    def minFlips(self, mat: List[List[int]]) -> int:
        m, n = len(mat), len(mat[0])
        state = sum(cell << (i * n + j) for i, row in enumerate(mat) for j, cell in enumerate(row))
        queue = deque([(state, 0)])
        visited = set([state])

        while queue:
            cur, step = queue.popleft()
            if cur == 0:
                return step

            for i in range(m):
                for j in range(n):
                    next = cur
                    for r, c in (i, j), (i, j + 1), (i, j - 1), (i + 1, j), (i - 1, j):
                        if m > r >= 0 <= c < n:
                            # 反转相邻位
                            next ^= 1 << (r * n + c)
                    if next not in visited:
                        visited.add(next)
                        queue.append((next, step + 1))

        return -1


print(Solution().minFlips([[0, 0], [0, 1]]))
