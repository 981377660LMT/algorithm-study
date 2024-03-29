# 934. 最短的桥

from collections import deque
from typing import List

# 在给定的二维二进制数组 A 中，存在两座岛。（岛是由四面相连的 1 形成的一个最大组。）
# 现在，我们可以将 0 变为 1，以使两座岛连接起来，变成一座岛。
# !返回必须翻转的 0 的最小数目。（可以保证答案至少是 1 。）

# 思路:
# 1.找起点
# 2.dfs将岛全部加入queue 原地标记-1
# 3.多源bfs最短路径 找到1就返回

DIR4 = ((0, 1), (0, -1), (1, 0), (-1, 0))


class Solution:
    def shortestBridge(self, grid: List[List[int]]) -> int:
        def floodfill(row: int, col: int) -> None:
            grid[row][col] = -1
            queue.append((row, col))
            for dr, dc in DIR4:
                nextRow, nextCol = row + dr, col + dc
                if 0 <= nextRow < n and 0 <= nextCol < n and grid[nextRow][nextCol] == 1:
                    floodfill(nextRow, nextCol)

        n, step, queue = len(grid), 0, deque()
        floodfill(*next((r, c) for r in range(n) for c in range(n) if grid[r][c] == 1))

        while queue:
            len_ = len(queue)
            for _ in range(len_):
                curRow, curCol = queue.popleft()
                for dr, dc in DIR4:
                    nextRow, nextCol = curRow + dr, curCol + dc
                    if 0 <= nextRow < n and 0 <= nextCol < n:
                        if grid[nextRow][nextCol] == 1:
                            return step
                        elif grid[nextRow][nextCol] == 0:
                            grid[nextRow][nextCol] = -1
                            queue.append((nextRow, nextCol))
            step += 1

        raise ValueError("No solution")


print(
    Solution().shortestBridge(
        [[1, 1, 1, 1, 1], [1, 0, 0, 0, 1], [1, 0, 1, 0, 1], [1, 0, 0, 0, 1], [1, 1, 1, 1, 1]]
    )
)
