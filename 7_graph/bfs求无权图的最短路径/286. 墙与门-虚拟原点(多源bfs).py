# 你被给定一个 m × n 的二维网格 rooms ，网格中有以下三种可能的初始化值：

# -1 表示墙或是障碍物
# 0 表示一扇门
# INF 无限表示一个空的房间。然后，我们用 2^31 - 1 = 2147483647 代表 INF。
# 你可以认为通往门的距离总是小于 2147483647 的。
# 你要给每个空房间位上填上该房间到 最近门的距离 ，如果无法到达门，则填 INF 即可。


from collections import deque
from typing import List


INF = 2147483647
DIR4 = ((0, 1), (0, -1), (1, 0), (-1, 0))


class Solution:
    def wallsAndGates(self, rooms: List[List[int]]) -> None:
        """
        Do not return anything, modify rooms in-place instead.
        """
        ROW, COL = len(rooms), len(rooms[0])
        queue = deque()
        for r in range(ROW):
            for c in range(COL):
                if rooms[r][c] == 0:
                    queue.append((r, c, 0))

        while queue:
            curRow, curCol, curDist = queue.popleft()
            for dr, dc in DIR4:
                nextRow, nextCol = curRow + dr, curCol + dc
                if 0 <= nextRow < ROW and 0 <= nextCol < COL and rooms[nextRow][nextCol] == INF:
                    rooms[nextRow][nextCol] = curDist + 1
                    queue.append((nextRow, nextCol, curDist + 1))
