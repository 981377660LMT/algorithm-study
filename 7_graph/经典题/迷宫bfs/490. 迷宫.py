from typing import List
from collections import deque

# 由空地（用 0 表示）和墙（用 1 表示）组成的迷宫 maze 中有一个球
# 球可以途经空地向 上、下、左、右 四个方向滚动，且在遇到墙壁前不会停止滚动
# 球可以途经空地向 上、下、左、右 四个方向滚动，且在遇到墙壁前不会停止滚动


class Solution:
    def hasPath(self, maze: List[List[int]], start: List[int], destination: List[int]) -> bool:
        if not any(maze):
            return False

        m, n = len(maze), len(maze[0])
        queue = deque([(start[0], start[1])])
        visited = [[False for _ in range(n)] for _ in range(m)]
        visited[start[0]][start[1]] = True

        while queue:
            r, c = queue.popleft()
            if [r, c] == destination:
                return True
            for dr, dc in ((0, 1), (1, 0), (0, -1), (-1, 0)):
                nr = r + dr
                nc = c + dc
                # 沿着这个方向一直进行
                while 0 <= nr < m and 0 <= nc < n and maze[nr][nc] == 0:
                    nr += dr
                    nc += dc
                # 碰壁后退
                nr -= dr
                nc -= dc
                if not visited[nr][nc]:
                    visited[nr][nc] = True
                    queue.append((nr, nc))

        return False


print(
    Solution().hasPath(
        [[0, 0, 1, 0, 0], [0, 0, 0, 0, 0], [0, 0, 0, 1, 0], [1, 1, 0, 1, 1], [0, 0, 0, 0, 0]],
        [0, 4],
        [4, 4],
    )
)

