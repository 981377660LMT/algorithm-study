# 0表示空地可以通过，1表示墙壁无法通过，
# 2表示不仅可以通过，还可以在该格放置一个传送门，
# 3表示有且仅有的唯一的固定传送门。

# 在游戏开始时，牛牛可以选择一块2类型的格子放置传送门并且不可以中途更改，
# 在游戏过程中，
# 如果到达其中一个传送门，则可以传送往另一个传送门

# 牛牛想知道从起点走到终点最少需要走几步（使用传送门也算作一步）

# n<=500
# 两次bfs求出起点与终点到各个传送门的距离,之后枚举传送门的放置地点即可
from collections import deque
from typing import List

DIR4 = [(0, 1), (1, 0), (0, -1), (-1, 0)]


class Solution:
    def solve(self, n: int, grid: List[List[int]]) -> int:
        # 起点到各个点的dist
        dist1 = dict({(0, 0): 0})
        # 终点到各个点的dist
        dist2 = dict({(n - 1, n - 1): 0})
        doorX, doorY = -1, -1
        portals = []

        queue = deque([(0, 0, 0)])
        visited = set([(0, 0)])
        while queue:
            x, y, d = queue.popleft()
            for dx, dy in DIR4:
                nx, ny = x + dx, y + dy
                if 0 <= nx < n and 0 <= ny < n and grid[nx][ny] != 1 and (nx, ny) not in visited:
                    if grid[nx][ny] == 3:
                        doorX, doorY = nx, ny
                    elif grid[nx][ny] == 2:
                        portals.append((nx, ny))
                    visited.add((nx, ny))
                    dist1[(nx, ny)] = d + 1
                    queue.append((nx, ny, d + 1))

        queue = deque([(n - 1, n - 1, 0)])
        visited = set([(n - 1, n - 1)])
        while queue:
            x, y, d = queue.popleft()
            for dx, dy in DIR4:
                nx, ny = x + dx, y + dy
                if 0 <= nx < n and 0 <= ny < n and grid[nx][ny] != 1 and (nx, ny) not in visited:
                    visited.add((nx, ny))
                    dist2[(nx, ny)] = d + 1
                    queue.append((nx, ny, d + 1))

        res = 0x3F3F3F3F
        for x, y in portals:
            res = min(
                res,
                dist1[(x, y)] + dist2[(doorX, doorY)] + 1,
                dist1[(doorX, doorY)] + dist2[(x, y)] + 1,
            )
        return res

