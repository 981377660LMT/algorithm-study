from typing import List
from collections import deque


DIRECTIONS = [(1, 0), (-1, 0), (0, 1), (0, -1)]


# 表示一个空的房间
EMPTY = 2147483647
# 表示一墙
WALL = -1
# 表示一扇门
GATE = 0


class Solution:
    def wallsAndGates(self, rooms: List[List[int]]) -> None:
        """
        Do not return anything, modify rooms in-place instead.
        """

        row, col = len(rooms), len(rooms[0])
        queue = deque()

        for r in range(row):
            for c in range(col):
                if rooms[r][c] == GATE:
                    queue.append((r, c, 0))

        while queue:
            cur_x, cur_y, dist = queue.popleft()
            for dx, dy in DIRECTIONS:
                next_x, next_y = cur_x + dx, cur_y + dy
                if 0 <= next_x < row and 0 <= next_y < col and rooms[next_x][next_y] == EMPTY:
                    rooms[next_x][next_y] = dist + 1
                    queue.append((next_x, next_y, dist + 1))


map = [
    [2147483647, -1, 0, 2147483647],
    [2147483647, 2147483647, 2147483647, -1],
    [2147483647, -1, 2147483647, -1],
    [0, -1, 2147483647, 2147483647],
]
Solution().wallsAndGates(map)
print(map)
